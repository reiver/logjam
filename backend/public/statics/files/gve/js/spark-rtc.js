class SparkRTC {
    started = false;
    myPeerConnectionConfig = {
        iceServers: [
            {
                url: 'stun:stun.l.google.com:19302'
            },
            {
                url: 'stun:stun1.l.google.com:19302'
            },
            {

                url: 'stun:stun2.l.google.com:19302'
            },
            {

                url: 'stun:stun3.l.google.com:19302'
            },
            {

                url: 'stun:stun4.l.google.com:19302'
            },
            {
                url: "turn:turn1.turn.group.video:3478",
                username: "turnuser",
                credential: "dJ4kP05PHcKN8Ubu",
            },
            {
                url: "turn:turn2.turn.group.video:3478",
                username: "turnuser",
                credential: "XzfVP8cpNEy17hws",
            },
            {
                url: "turns:turn1.turn.group.video:443",
                username: "turnuser",
                credential: "dJ4kP05PHcKN8Ubu",
            },
            {
                url: "turns:turn2.turn.group.video:443",
                username: "turnuser",
                credential: "XzfVP8cpNEy17hws",
            },
        ],
    };
    role = 'unknown';
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
    checkStatusHandle;
    handleVideoOfferMsg = async (msg) => {
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
    };
    handleMessage = async (event) => {
        if (event.data.indexOf('pong') < 0)
            console.log('onmessage event', event);
        let msg;
        try {
            msg = JSON.parse(event.data);
        } catch (e) {
            return;
        }
        msg.data = (msg.Data && !msg.data) ? msg.Data : msg.data;
        msg.type = (msg.Type && !msg.type) ? msg.Type : msg.type;

        let audiencePeerConnection;
        if (msg.type !== 'pong')
            console.log('onmessage object', msg);
        switch (msg.type) {
            case 'video-offer':
            case 'alt-video-offer':
                this.handleVideoOfferMsg(msg);
                break;
            case 'video-answer':
            case 'alt-video-answer':
                audiencePeerConnection = this.createOrGetPeerConnection(msg.data, true);
                try {
                    await audiencePeerConnection.setRemoteDescription(new RTCSessionDescription(msg.sdp));
                } catch (e) {
                    console.log('setRemoteDescription failed with exception: ' + e.message);
                }
                break;
            case 'new-ice-candidate':
            case 'alt-new-ice-candidate':
                audiencePeerConnection = this.createOrGetPeerConnection(msg.data);
                this.iceCandidates.push(new RTCIceCandidate(msg.candidate));
                if (audiencePeerConnection && audiencePeerConnection.remoteDescription) {
                    audiencePeerConnection.addIceCandidate(this.iceCandidates.pop());
                }
                break;
            case 'role':
                this.role = 'audience';
                if (msg.data === "no:broadcast") {
                    alert("You are not a broadcaster anymore!");
                    this.socket.close();
                } else if (msg.data === "yes:broadcast") {
                    this.role = 'broadcast';
                    this.localStream = await navigator.mediaDevices.getUserMedia({
                        audio: true,
                        video: true,
                    });
                    this.remoteStreams.push(this.localStream);
                    if (this.localStreamChangeCallback)
                        this.localStreamChangeCallback(this.localStream);
                } else {
                    if (this.localStreamChangeCallback)
                        this.localStreamChangeCallback(null);
                }
                break;
            case 'start':
                if (msg.error) {
                    alert(msg.error);
                    return;
                }

                this.myUsername = msg.data;
                break;
            case 'add_audience':
            case 'add_broadcast_audience':
                this.connectToAudience(msg.data);
                break;
            case 'alt-broadcast-approve':
                this.sendStreamTo(msg.data, this.localStream);
                if (this.localStreamChangeCallback)
                    this.localStreamChangeCallback(this.localStream);
                break;
            case 'alt-broadcast':
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
                    }
                }
                break;
            case 'tree':
                if (this.treeCallback) this.treeCallback(msg.data);
                break;
            case 'event-broadcaster-disconnected':
                // onLoad();
                // start();
                // setTimeout(handleClick, 1000);
                // window.location.reload();
                if (this.startProcedure)
                    this.startProcedure();
                break;
            default:
                break;
        }
    };
    ping = () => {
        console.log('ping');
        this.socket.send(JSON.stringify({
            type: this.treeCallback ? "tree" : "ping",
        }));
    };
    setupSignalingSocket = (url, myName, roomName) => {
        return new Promise((resolve, reject) => {
            if (this.pingInterval) {
                clearInterval(this.pingInterval);
                this.pingInterval = null;
            }
            if (myName)
                this.myName = myName;
            if (roomName)
                this.roomName = roomName;

            const socket = new WebSocket(url + '?room=' + this.roomName);
            socket.onmessage = this.handleMessage;
            socket.onopen = () => {
                socket.send(
                    JSON.stringify({
                        type: "start",
                        data: myName,
                    })
                );
                this.pingInterval = setInterval(this.ping, 2000);
                resolve(socket);
            };
            socket.onclose = () => {
                this.remoteStreamNotified = false;
                this.myPeerConnectionArray = {};
                if (this.signalingDisconnectedCallback) this.signalingDisconnectedCallback;
            };
            socket.onerror = (error) => {
                console.log("WebSocket error: ", error);
                reject(error);
            };
            this.socket = socket;
        })
    };
    startShareScreen = async () => {
        try {
            const shareStream = await navigator.mediaDevices
                .getDisplayMedia({
                    audio: true,
                    video: true,
                });
            this.remoteStreams.push(shareStream);
            for (const userId in this.myPeerConnectionArray) {
                const apeerConnection = this.myPeerConnectionArray[userId];
                // if (!apeerConnection.isAdience) return;

                shareStream.getTracks().forEach((track) => {
                    apeerConnection.addTrack(track, shareStream);
                });
            }
            return shareStream;
        } catch (e) {
            console.log(e);
            alert('Unable to get access to screenshare.');
        }
    };
    startBroadcasting = async (data = 'broadcast') => {
        try {
            this.localStream = await navigator.mediaDevices.getUserMedia({
                audio: true,
                video: true,
            });
            this.remoteStreams.push(this.localStream);
            this.socket.send(
                JSON.stringify({
                    type: "role",
                    data,
                })
            );
            return this.localStream;
        } catch (e) {
            console.log(e);
            alert('Unable to get access to your webcam and microphone.');
        }
    };
    startReadingBroadcast = async () => {
        this.socket.send(
            JSON.stringify({
                type: "role",
                data: this.role,
            })
        );
    };
    raiseHand = () => {
        if (this.startedRaiseHand) return;
        this.startedRaiseHand = true;
        return this.startBroadcasting('alt-broadcast');
    };
    newPeerConnectionInstance = (target, theStream, isAdience = false) => {
        const peerConnection = new RTCPeerConnection(this.myPeerConnectionConfig);
        peerConnection.isAdience = isAdience;

        peerConnection.onicecandidate = (event) => {
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
                alert('onnegotiationneeded failed:', e);
            }
        };

        peerConnection.ontrack = (event) => {
            const stream = event.streams[0];
            if (this.localStream && this.localStream.id === stream.id) return;
            if (this.newTrackCallback && !this.newTrackCallback(stream)) return;
            if (this.remoteStreams.indexOf(stream) !== -1) return;
            if (this.remoteStreams.length === 0) {
                this.parentStreamId = stream.id;
            }
            stream.oninactive = (event) => {
                console.log('stream inactive', event);
                if (this.remoteStreamDCCallback) this.remoteStreamDCCallback(event.target);
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
                this.remoteStreams.splice(this.remoteStreams.indexOf(peerConnection.getRemoteStreams()[0]), 1);
                if (this.parentStreamId && this.parentStreamId === peerConnection.getRemoteStreams()[0].id) {
                    if (this.remoteStreamDCCallback) {
                        this.remoteStreams.forEach((strm) => {
                            this.remoteStreamDCCallback(strm);
                        });
                    }
                    this.parentStreamId = undefined;
                }
            };
            if (this.remoteStreamCallback)
                this.remoteStreamCallback(stream);
            this.remoteStreams.push(stream);
            if (!this.remoteStreamNotified) {
                this.remoteStreamNotified = true;
                this.socket.send(
                    JSON.stringify({
                        type: "stream",
                        data: "true",
                    })
                );
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
                clearTimeout(this.checkStatusHandle);
                this.checkState().then(res => this.checkStatusHandle = res);
            }
        };

        peerConnection.oniceconnectionstatechange = (event) => {
            if (peerConnection.iceConnectionState == 'disconnected') {
                if (peerConnection.getRemoteStreams().length === 0) return;
                console.log('disconnected', event);
                if (this.remoteStreamDCCallback) this.remoteStreamDCCallback(peerConnection.getRemoteStreams()[0]);
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
                this.remoteStreams.splice(this.remoteStreams.indexOf(peerConnection.getRemoteStreams()[0]), 1);
                if (this.parentStreamId && this.parentStreamId === peerConnection.getRemoteStreams()[0].id) {
                    if (this.remoteStreamDCCallback) {
                        this.remoteStreams.forEach((strm) => {
                            this.remoteStreamDCCallback(strm);
                        });
                    }
                    this.parentStreamId = undefined;
                }
            }
        };

        return peerConnection;
    };
    createOrGetPeerConnection = (audienceName, isAdience = false) => {
        if (this.myPeerConnectionArray[audienceName]) return this.myPeerConnectionArray[audienceName];

        this.myPeerConnectionArray[audienceName] = this.newPeerConnectionInstance(audienceName, true, isAdience);

        return this.myPeerConnectionArray[audienceName];
    };
    connectToAudience = (audienceName) => {
        if (!this.localStream && this.remoteStreams.length === 0) return;
        if (!this.myPeerConnectionArray[audienceName]) {
            this.myPeerConnectionArray[audienceName] = this.newPeerConnectionInstance(audienceName, this.localStream || this.remoteStreams, true);
        }

        if (this.remoteStreams.length > 0) {
            this.remoteStreams.forEach((astream) => {
                astream.getTracks().forEach((track) => {
                    this.myPeerConnectionArray[audienceName].addTrack(track, astream);
                });
            });
        }
    };
    sendStreamTo = (target, stream) => {
        const peerConnection = this.createOrGetPeerConnection(target, false);
        stream.getTracks().forEach((track) => {
            peerConnection.addTrack(track, stream);
        });
    };
    start = () => {
        // if (this.role === 'broadcast') {
        //     return this.startBroadcasting();
        // }

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

    checkState = async () => {
        console.log('[checkState] startProcedure', this.started, this.role, this.parentStreamId);
        if (!this.started || !this.startProcedure || this.role === 'broadcast') return;

        if (!this.parentStreamId) {
            //await this.startProcedure();
        } 
        this.checkStatusHandle = setTimeout(this.checkState, 1000);
        return this.checkStatusHandle;
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
    }
}