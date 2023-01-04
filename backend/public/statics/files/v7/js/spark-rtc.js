class SparkRTC {
    started = false;
    myPeerConnectionConfig = {
        iceServers,
    };
    role = 'broadcast';
    localStream;
    remoteStreamNotified = false;
    remoteStreams = [];
    socket;
    myName = 'NoName';
    roomName = 'SparkRTC';
    myUsername = 'NoUsername';
    myPeerConnectionArray = {};
    iceCandidates = [];
    pingInterval;
    raiseHands = [];
    startedRaiseHand = false;
    targetStreams = {};
    parentStreamId;
    broadcasterStatus = '';
    constraints = {
        audio: true,
        video: true,
    };
    parentDC = false;
    handleVideoOfferMsg = async (msg) => {
        this.log(`[handleVideoOfferMsg] ${msg.name}`);
        const broadcasterPeerConnection = this.createOrGetPeerConnection(msg.name);
        await broadcasterPeerConnection.setRemoteDescription(new RTCSessionDescription(msg.sdp));
        await broadcasterPeerConnection.setLocalDescription(await broadcasterPeerConnection.createAnswer());

        this.socket.send(
            JSON.stringify({
                name: this.myUsername,
                target: msg.name,
                type: "video-answer",
                sdp: broadcasterPeerConnection.localDescription,
            })
        );
        this.log(`[handleVideoOfferMsg] send video-answer to ${msg.name} from ${this.myUsername}`);
    };
    handleMessage = async (event) => {
        this.log(`[handleMessage] ${event.data}`);
        let msg;
        try {
            msg = JSON.parse(event.data);
        } catch (e) {
            return;
        }
        msg.data = (msg.Data && !msg.data) ? msg.Data : msg.data;
        msg.type = (msg.Type && !msg.type) ? msg.Type : msg.type;

        let audiencePeerConnection;
        switch (msg.type) {
            case 'video-offer':
            case 'alt-video-offer':
                this.log(`[handleMessage] handleVideoOfferMsg ${msg.type}`);
                this.handleVideoOfferMsg(msg);
                break;
            case 'video-answer':
            case 'alt-video-answer':
                this.log(`[handleMessage] setRemoteDescription ${msg.type}`);
                audiencePeerConnection = this.createOrGetPeerConnection(msg.data, true);
                try {
                    await audiencePeerConnection.setRemoteDescription(new RTCSessionDescription(msg.sdp));
                } catch (e) {
                    console.log('setRemoteDescription failed with exception: ' + e.message);
                }
                break;
            case 'new-ice-candidate':
            case 'alt-new-ice-candidate':
                this.log(`[handleMessage] addIceCandidate ${msg.type}`);
                audiencePeerConnection = this.createOrGetPeerConnection(msg.data);
                this.iceCandidates.push(new RTCIceCandidate(msg.candidate));
                if (audiencePeerConnection && audiencePeerConnection.remoteDescription) {
                    audiencePeerConnection.addIceCandidate(this.iceCandidates.pop());
                }
                break;
            case 'role':
                this.log(`[handleMessage] ${msg.type}`);
                if (this.role === 'broadcast') {
                    if (msg.data === "no:broadcast") {
                        alert("You are not a broadcaster anymore!");
                        this.socket.close();
                    } else if (msg.data === "yes:broadcast") {
                        if (this.localStreamChangeCallback)
                            this.localStreamChangeCallback(this.localStream);
                    } else {
                        if (this.localStreamChangeCallback)
                            this.localStreamChangeCallback(null);
                    }
                } else if (msg.data === 'no:audience') {
                    if (this.remoteStreamDCCallback) {
                        try {
                            this.remoteStreamDCCallback('no-stream');
                        } catch (e) { console.log(e) }
                    }
                }
                break;
            case 'start':
                this.log(`[handleMessage] start ${JSON.stringify(msg)}`);
                if (msg.error) {
                    alert(msg.error);
                    return;
                }

                this.myUsername = msg.data;
                break;
            case 'add_audience':
            case 'add_broadcast_audience':
                this.log(`[handleMessage] add audience ${msg}`);
                this.updateTheStatus(`New Audience arrived ${msg.data}`);
                this.connectToAudience(msg.data);
                break;
            case 'alt-broadcast-approve':
                this.log(`[handleMessage] alt-broadcast-approve ${msg}`);
                this.sendStreamTo(msg.data, this.localStream);
                break;
            case 'alt-broadcast':
                this.log(`[handleMessage] ${msg.type}`);
                if (this.role === 'broadcast') {
                    if (this.raiseHands.indexOf(msg.data) === -1) {
                        this.raiseHands.push(msg.data);
                        if (!this.raiseHandConfirmation(`${msg.name} wants to broadcast, do you approve?`)) return;
                        this.socket.send(
                            JSON.stringify({
                                type: "alt-broadcast-approve",
                                target: msg.data,
                            })
                        );
                        this.log(`[handleMessage] ${msg.type} approving raised hand ${msg.data}`);
                    }
                } else {
                    this.spreadLocalStream();
                }
                break;
            case 'tree':
                this.log(`[handleMessage] ${msg.type}`);
                if (this.treeCallback) this.treeCallback(msg.data);
                break;
            case 'broadcasting':
                if (this.role === 'broadcast') return;
                this.log(`[handleMessage] ${msg.type}`);
                this.startProcedure();
                break;
            case 'event-broadcaster-disconnected':
                // console.log('event-broadcaster-disconnected');
                this.parentDC = true;
                // if (this.role === 'broadcast') return;
                // setTimeout(() => {
                //     console.log('Reconnecting ...');
                //     this.startProcedure();    
                // }, 1000);
                break;
            default:
                this.log(`[handleMessage] default ${JSON.stringify(msg)}`);
                break;
        }
    };
    ping = () => {
        this.socket.send(JSON.stringify({
            type: this.treeCallback ? "tree" : "ping",
        }));
    };
    setupSignalingSocket = (url, myName, roomName) => {
        this.log(`[setupSignalingSocket] url='${url}' myName='${myName}' roomName='${roomName}'`);
        return new Promise((resolve, reject) => {
            if (this.pingInterval) {
                clearInterval(this.pingInterval);
                this.pingInterval = null;
            }
            if (myName)
                this.myName = myName;
            if (roomName)
                this.roomName = roomName;

            this.log(`[setupSignalingSocket] installing socket`);
            const socket = new WebSocket(url + '?room=' + this.roomName);
            socket.onmessage = this.handleMessage;
            socket.onopen = () => {
                socket.send(
                    JSON.stringify({
                        type: "start",
                        data: myName,
                    })
                );
                this.pingInterval = setInterval(this.ping, 60000);
                this.log(`[setupSignalingSocket] socket onopen and sent start`);
                resolve(socket);
            };
            socket.onclose = () => {
                this.remoteStreamNotified = false;
                this.myPeerConnectionArray = {};
                if (this.signalingDisconnectedCallback) this.signalingDisconnectedCallback;
                this.log(`[setupSignalingSocket] socket onclose`);
                this.started = false;
                if (this.startProcedure) this.startProcedure();
            };
            socket.onerror = (error) => {
                console.log("WebSocket error: ", error);
                reject(error);
                this.log(`[setupSignalingSocket] socket onerror`);
                alert('Can not connect to server');
                window.location.reload();
            };
            this.socket = socket;

        })
    };
    startShareScreen = async () => {
        this.log(`[handleMessage] startShareScreen`);
        try {
            const shareStream = await navigator.mediaDevices
                .getDisplayMedia({
                    audio: true,
                    video: true,
                });
            this.remoteStreams.push(shareStream);
            for (const userId in this.myPeerConnectionArray) {
                const apeerConnection = this.myPeerConnectionArray[userId];
                if (!apeerConnection.isAdience) return;

                shareStream.getTracks().forEach((track) => {
                    apeerConnection.addTrack(track, shareStream);
                });
            }
            return shareStream;
        } catch (e) {
            console.log(e);
            this.log(`[handleMessage] startShareScreen error ${e}`);
            alert('Unable to get access to screenshare.');
        }
    };
    startBroadcasting = async (data = 'broadcast') => {
        this.log(`[startBroadcasting] ${data}`);
        try {
            if (!this.localStream) {
                this.updateTheStatus(`Trying to get local stream`);
                if (!this.constraints.audio && !this.constraints.video) {
                    this.updateTheStatus(`No media device available`);
                    throw new Error('No media device available');
                }
                this.localStream = await navigator.mediaDevices.getUserMedia(this.constraints);
                this.updateTheStatus(`Local stream loaded`);
                this.log(`[startBroadcasting] local stream loaded`);
                this.remoteStreams.push(this.localStream);
            }
            this.updateTheStatus(`Request Broadcast Role`);
            this.socket.send(
                JSON.stringify({
                    type: "role",
                    data,
                })
            );
            this.log(`[startBroadcasting] send role`);
            return this.localStream;
        } catch (e) {
            this.updateTheStatus(`Error Start Broadcasting`);
            console.log(e);
            this.log(`[startBroadcasting] ${e}`);
            alert('Unable to get access to your webcam and microphone.');
        }
    };
    startReadingBroadcast = async () => {
        this.log(`[startReadingBroadcast]`);
        this.updateTheStatus(`Request Audience Role`);
        this.socket.send(
            JSON.stringify({
                type: "role",
                data: "audience",
            })
        );
        this.log(`[startReadingBroadcast] send role audience`);
    };
    raiseHand = () => {
        if (this.startedRaiseHand) return;
        this.startedRaiseHand = true;
        return this.startBroadcasting('alt-broadcast');
    };
    newPeerConnectionInstance = (target, theStream, isAdience = false) => {
        this.log(`[newPeerConnectionInstance] target='${target}' theStream='${theStream}' isAdience='${isAdience}'`);
        const peerConnection = new RTCPeerConnection(this.myPeerConnectionConfig);
        peerConnection.isAdience = isAdience;

        peerConnection.onconnectionstatechange = (ev) => {
            // console.log(`[newPeerConnectionInstance] peerConnection.onconnectionstatechange `, ev);
        }

        peerConnection.onicecandidate = (event) => {
            this.updateStatus(`Peer Connection ice candidate arrived for ${target}: [${event.candidate}]`);
            this.log(`[newPeerConnectionInstance] onicecandidate event.candidate='${JSON.stringify(event.candidate)}'`);
            if (event.candidate) {
                this.socket.send(
                    JSON.stringify({
                        type: "new-ice-candidate",
                        candidate: event.candidate,
                        target,
                    })
                );
            }
        };

        peerConnection.onnegotiationneeded = async () => {
            this.updateStatus(`Peer Connection negotiation needed for ${target} preparing video offer`);
            this.log(`[newPeerConnectionInstance] onnegotiationneeded`);
            try {
                await peerConnection.setLocalDescription(
                    await peerConnection.createOffer()
                );
                this.socket.send(
                    JSON.stringify({
                        type: "video-offer",
                        sdp: peerConnection.localDescription,
                        target,
                        name: this.myUsername,
                    })
                );
            } catch (e) {
                console.log(e);
                this.log(`[newPeerConnectionInstance] failed ${e}`);
                // alert('onnegotiationneeded failed:', e);
            }
        };

        peerConnection.ontrack = (event) => {
            this.updateStatus(`Peer Connection track received for ${target} stream ids [${event.streams.map((s) => s.id).join(',')}]`);
            this.parentDC = false;
            this.log(`[newPeerConnectionInstance] ontrack ${JSON.stringify(event.streams)}`);
            const stream = event.streams[0];
            if (this.localStream && this.localStream.id === stream.id) return;
            if (this.newTrackCallback && !this.newTrackCallback(stream)) return;
            if (this.remoteStreams.indexOf(stream) !== -1) return;
            if (this.remoteStreams.length === 0) {
                this.parentStreamId = stream.id;
            }
            stream.oninactive = (event) => {
                this.log(`[newPeerConnectionInstance] stream.oninactive ${JSON.stringify(event)}`);
                console.log('[stream.oninactive] event', event);
                this.remoteStreamNotified = false;
                const theEventStream = event.currentTarget;
                const trackIds = peerConnection.getReceivers().map((receiver) => receiver.track.id);
                trackIds.forEach((trackId) => {
                    for (const userId in this.myPeerConnectionArray) {
                        if (userId === target) continue;
                        const apeerConnection = this.myPeerConnectionArray[userId];
                        if (!apeerConnection.isAdience) return;
                        const allSenders = apeerConnection.getSenders();
                        for (const sender of allSenders) {
                            if (!sender.track) continue;
                            if (sender.track.id === trackId) {
                                try {
                                    apeerConnection.removeTrack(sender);
                                } catch (e) {
                                    console.log(e);
                                }
                            }
                        }
                    }
                });
                this.remoteStreams.splice(this.remoteStreams.indexOf(theEventStream), 1);
                if (this.parentStreamId && this.parentStreamId === theEventStream.id) {
                    if (this.remoteStreamDCCallback) {
                        this.remoteStreams.forEach((strm) => {
                            this.remoteStreamDCCallback(strm);
                        });
                    }
                    this.parentStreamId = undefined;
                }
                if (this.remoteStreamDCCallback) {
                    try {
                        this.remoteStreamDCCallback(event.target);
                    } catch {}
                }
            };
            try {
                if (this.remoteStreamCallback)
                    this.remoteStreamCallback(stream);
            } catch { }
            this.remoteStreams.push(stream);
            if (!this.remoteStreamNotified) {
                this.remoteStreamNotified = true;
                this.log(`[newPeerConnectionInstance] A7`);
                this.socket.send(
                    JSON.stringify({
                        type: "stream",
                        data: "true",
                    })
                );
                this.log(`[newPeerConnectionInstance] stream message`);
            }
            this.targetStreams[target] = stream.id;

            for (const userId in this.myPeerConnectionArray) {
                if (userId === target) continue;
                const apeerConnection = this.myPeerConnectionArray[userId];
                if (!apeerConnection.isAdience) return;

                stream.getTracks().forEach((track) => {
                    apeerConnection.addTrack(track, stream);
                });
            }

            if (!this.started) {
                this.started = true;
            }
        };

        peerConnection.oniceconnectionstatechange = (event) => {
            this.log(`[newPeerConnectionInstance] oniceconnectionstatechange peerConnection.iceConnectionState = ${peerConnection.iceConnectionState} event = ${JSON.stringify(event)}`);
            if (peerConnection.iceConnectionState == 'disconnected') {
                this.remoteStreamNotified = false;
                console.log('[peerConnection.oniceconnectionstatechange] DC event', event);
                if (peerConnection.getRemoteStreams().length === 0) return;
                const trackIds = peerConnection.getReceivers().map((receiver) => receiver.track.id);
                trackIds.forEach((trackId) => {
                    console.log('[peerConnection.oniceconnectionstatechange] DC trackId', trackId);
                    for (const userId in this.myPeerConnectionArray) {
                        if (userId === target) continue;
                        console.log('[peerConnection.oniceconnectionstatechange] DC userId', userId);
                        const apeerConnection = this.myPeerConnectionArray[userId];
                        if (!apeerConnection.isAdience) return;
                        const allSenders = apeerConnection.getSenders();
                        for (const sender of allSenders) {
                            if (!sender.track) continue;
                            if (sender.track.id === trackId) {
                                console.log('[peerConnection.oniceconnectionstatechange] DC sender');
                                try {
                                    apeerConnection.removeTrack(sender);
                                } catch (e) {
                                    console.log(e);
                                }
                            }
                        }
                    }
                });
                const allStreams = peerConnection.getRemoteStreams();
                console.log({allStreams});
                for (let i = 0; i < allStreams.length; i++)
                    this.remoteStreams.splice(this.remoteStreams.indexOf(allStreams[i]), 1);
                
                if (this.parentStreamId && allStreams.map((s) => s.id).includes(this.parentStreamId)) {
                    this.updateTheStatus(`Parent stream is disconnected`);
                    if (this.remoteStreamDCCallback) {
                        this.remoteStreams.forEach((strm) => {
                            try {
                                this.remoteStreamDCCallback(strm);
                            } catch {}
                        });
                    }
                    this.parentStreamId = undefined;
                    this.remoteStreams = [];
                }
                this.startedRaiseHand = false;
                // if (this.role !== 'broadcast') this.getBroadcasterStatus();
                try {
                    if (this.remoteStreamDCCallback) this.remoteStreamDCCallback(peerConnection.getRemoteStreams()[0]);
                } catch { }
                if (this.parentDC) this.startProcedure();

            }
        };

        return peerConnection;
    };
    createOrGetPeerConnection = (audienceName, isAdience = false) => {
        this.log(`[createOrGetPeerConnection] audienceName = ${audienceName}, isAdience = ${isAdience}`);
        if (this.myPeerConnectionArray[audienceName]) return this.myPeerConnectionArray[audienceName];

        this.myPeerConnectionArray[audienceName] = this.newPeerConnectionInstance(audienceName, true, isAdience);
        this.log(`[createOrGetPeerConnection] generate newPeerConnectionInstance`);

        return this.myPeerConnectionArray[audienceName];
    };
    connectToAudience = (audienceName) => {
        this.updateStatus(`Connecting to ${audienceName}`);
        this.log(`[handleMessage] connectToAudience ${audienceName}`);
        if (!this.localStream && this.remoteStreams.length === 0) return;
        if (!this.myPeerConnectionArray[audienceName]) {
            this.updateStatus(`Creating peer connection to ${audienceName}`);
            this.myPeerConnectionArray[audienceName] = this.newPeerConnectionInstance(audienceName, this.localStream || this.remoteStreams, true);
        }
        this.log(`[handleMessage] generate newPeerConnectionInstance`);

        if (this.remoteStreams.length > 0) {
            this.updateStatus(`publishing stream/s to ${audienceName}`);
            this.remoteStreams.forEach((astream) => {
                astream.getTracks().forEach((track) => {
                    this.myPeerConnectionArray[audienceName].addTrack(track, astream);
                });
            });
        }
    };
    sendStreamTo = (target, stream) => {
        this.log(`[handleMessage] sendStreamTo ${target}`);
        const peerConnection = this.createOrGetPeerConnection(target, false);
        stream.getTracks().forEach((track) => {
            peerConnection.addTrack(track, stream);
        });
    };
    start = async () => {
        this.updateTheStatus(`Starting`);
        this.log(`[start] ${this.role}`);
        this.updateTheStatus(`Getting media capabilities`);
        await this.getSupportedConstraints();
        if (this.role === 'broadcast') {
            this.updateTheStatus(`Start broadcasting`);
            return this.startBroadcasting();
        } else if (!this.constraints.audio && !this.constraints.video) {
            this.updateTheStatus(`No media removing raise hand`);
            document.getElementById('raise_hand').remove();
        }

        this.updateTheStatus(`Start as audience`);
        return this.startReadingBroadcast();
    };
    disableVideo = (enabled = false) => {
        this.localStream.getTracks().forEach((track) => {
            if (track.kind === 'video')
                track.enabled = enabled;
        });
    };
    disableAudio = (enabled = false) => {
        this.localStream.getTracks().forEach((track) => {
            if (track.kind === 'audio')
                track.enabled = enabled;
        });
    };
    wait = async (mil = 1000) => {
        return new Promise((res) => {
            setTimeout(() => {
                res();
            }, mil);
        });
    };
    getBroadcasterStatus = async () => {
        const max = 5;
        const reconnect = true;
        return new Promise((resolve, reject) => {
            this.socket.send(JSON.stringify({
                type: "broadcaster-status"
            }));

            let i = 0
            while (this.broadcasterStatus === '' && i < max) {
                this.wait();
                i++;
            }

            if (this.broadcasterStatus === '') {
                return reject(new Error('No response'));
            }

            if (reconnect) this.startProcedure();
            resolve(this.broadcasterStatus);
        });

    };
    getSupportedConstraints = async () => {
        const res = await navigator.mediaDevices.enumerateDevices();
        if (!res.find((r) => r.kind === 'audioinput')) {
            this.constraints.audio = false;
        }
        if (!res.find((r) => r.kind === 'videoinput')) {
            this.constraints.video = false;
        }
        if (this.constraintResults) this.constraintResults(this.constraints);
    };
    updateTheStatus = (status) => {
        if (this.updateStatus) {
            try {
                this.updateStatus(status);
            } catch { }
        }
    };
    lowerHand = async () => {
        console.log('[lowerHand] start');
        if (!this.localStream) return;
        let apeerConnection;
        for (const id in this.myPeerConnectionArray) {
            apeerConnection = this.myPeerConnectionArray[id];
            break;
        }
        const trackIds = this.localStream.getTracks().map((receiver) => receiver.id);
        console.log('[lowerHand] trackIds', trackIds);
        const allSenders = apeerConnection.getSenders();
        console.log('[lowerHand] allSenders', allSenders);
        for (const trackId of trackIds)
            for (const sender of allSenders) {
                console.log('[lowerHand] sender', sender);
                if (!sender.track) continue;
                if (sender.track.id === trackId) {
                    console.log('[lowerHand] DC sender');
                    try {
                        apeerConnection.removeTrack(sender);
                    } catch (e) {
                        console.log(e);
                    }
                }
            }
        this.localStream = null;
        this.startedRaiseHand = false;
    };
    spreadLocalStream = () => {
        for (const target in this.myPeerConnectionArray) {
            this.sendStreamTo(target, this.localStream);
        }
    };
    constructor(role, options = {}) {
        this.role = role;
        this.localStreamChangeCallback = options.localStreamChangeCallback;
        this.remoteStreamCallback = options.remoteStreamCallback;
        this.remoteStreamDCCallback = options.remoteStreamDCCallback;
        this.signalingDisconnectedCallback = options.signalingDisconnectedCallback;
        this.treeCallback = options.treeCallback;
        this.raiseHandConfirmation = options.raiseHandConfirmation || ((msg) => {
            return window.confirm(msg);
        });
        this.newTrackCallback = options.newTrackCallback;
        this.startProcedure = options.startProcedure;
        this.log = options.log || ((log) => { });
        this.constraintResults = options.constraintResults;
        this.log(`[constructor] ${this.role}`);
        this.updateStatus = options.updateStatus;
    }
}
