/** Your class description
 *
 * SparkRTC class is main class to setUP RTC client
 *
 */
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
    /**@type {{[key:string]:RTCPeerConnection}}*/
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
    parentDC = true;
    broadcasterDC = true;

    parentDisconnectionTimeOut = 2000; //2 second timeout to check parent is alive or not
    sendMessageInterval = 10; //send message to child after every 10 ms
    metaData = {};
    userStreamData = {};
    users = [];

    /**@type {{[trackId:string]: string}}*/
    trackToStreamMap = {};
    /**@type {"Enabled" | "Disabled"}*/
    lastVideoState = "Enabled"
    /**@type {"Enabled" | "Disabled"}*/
    lastAudioState = "Enabled"

    /**
     * Function to handle Peer Connection Offer, received from Other Peer
     *
     * and Return Peer Connection Answer to Other Peer
     *
     * @param {{name: string, sdp: RTCSessionDescriptionInit }} msg
     */
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

    /**
     * A socket handler to receive, message on webSocket,
     *
     * It parses message and based on message Type make decisions
     *
     * @param {*} event
     * @returns
     */
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
                    await audiencePeerConnection.addIceCandidate(this.iceCandidates.pop());
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
                        } catch (e) {
                            console.log(e)
                        }
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
                        if (this.raiseHandConfirmation) {
                            try {
                                const result = this.raiseHandConfirmation(`${msg.name} wants to broadcast, do you approve?`)
                                console.log(`[handleMessage] alt-broadcast result`, result);
                                if (result !== true) return;
                            } catch {
                            }
                        }
                        this.socket.send(
                            JSON.stringify({
                                type: "alt-broadcast-approve",
                                target: msg.data,
                            })
                        );
                        this.raiseHands.push(msg.data);
                        this.log(`[handleMessage] ${msg.type} approving raised hand ${msg.data}`);
                        this.getMetadata();
                        setTimeout(() => {
                            const metaData = this.metaData;
                            metaData.raiseHands = JSON.stringify(this.raiseHands);
                            this.setMetadata(metaData);
                        }, 1000);
                    }
                } else {
                    // this.spreadLocalStream();
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
            case 'event-reconnect':
            case 'event-broadcaster-disconnected':
                console.log('broadcaster dc', msg.type);
                this.broadcasterDC = true;
                const broadcasterId = this.broadcasterUserId();

                for (const u in this.myPeerConnectionArray) {
                    this.myPeerConnectionArray[u].close();
                }
                this.myPeerConnectionArray = {};
                this.remoteStreams = [];
                try {
                    if (this.remoteStreamDCCallback) this.remoteStreamDCCallback('no-stream');
                } catch {
                }
                this.localStream?.getTracks()?.forEach(track => track.stop());
                this.localStream = null;
                this.startedRaiseHand = false;
                break;
            case 'event-parent-dc':
                console.log('parentDC', msg.type);
                this.parentDC = true;
                break;
            case 'metadata-get':
            case 'metadata-set':
                this.log(`[handleMessage] ${msg.type}`);
                this.metaData = JSON.parse(msg.data);
                break;
            case 'user-by-stream':
                this.log(`[handleMessage] ${msg.type}`, msg.data);
                const [userId, userName, streamId, userRole] = msg.Data.split(',');
                this.userStreamData[streamId] = {
                    userId,
                    userName,
                    userRole,
                };
                break;
            case 'user-event':
                this.log(`[handleMessage] ${msg.type}`, msg.data);
                this.getMetadata();
                setTimeout(() => {
                    const users = JSON.parse(msg.data).map((u) => {
                        const video = u.streamId !== '' ? this.streamById(u.streamId) : null;
                        return {
                            id: u.id,
                            name: u.name,
                            role: u.role,
                            video,
                        };
                    });
                    this.users = users;

                    if (this.userListUpdated) {
                        try {
                            this.userListUpdated(users);
                        } catch {
                        }
                    }
                }, 1000);
                break;
            default:
                this.log(`[handleMessage] default ${JSON.stringify(msg)}`);
                break;
        }
    };

    /**
     * Function to get Broadcaster UserID from Array of PeerConnections
     *
     * @returns UserID
     */
    broadcasterUserId = () => {
        for (const userId in this.myPeerConnectionArray) {
            if (!this.myPeerConnectionArray[userId].isAdience) return userId;
        }
        return null;
    };

    /**
     * Ping function to, request Tree
     */
    ping = () => {
        this.socket.send(JSON.stringify({
            type: this.treeCallback ? "tree" : "ping",
        }));
    };

    /**
     * Function to setup Signaling WebSocket with backend
     *
     * @param {String} url - baseurl
     * @param {String} myName
     * @param {String} roomName - room identifier
     * @returns
     */
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
                this.pingInterval = setInterval(this.ping, 5000);
                this.log(`[setupSignalingSocket] socket onopen and sent start`);
                resolve(socket);
            };
            socket.onclose = () => {
                this.remoteStreamNotified = false;
                this.myPeerConnectionArray = {};
                if (this.signalingDisconnectedCallback) this.signalingDisconnectedCallback();
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

    stopShareScreen = async (stream) =>{
        if(stream){
             for (const userId in this.myPeerConnectionArray) {
                const apeerConnection = this.myPeerConnectionArray[userId];
                if (!apeerConnection.isAdience) return false;

                stream.getTracks().forEach((track) => {

                    const sender = apeerConnection.getSenders().find(sender => sender.track && sender.track.id === track.id);
                    
                    if (sender) {
                        apeerConnection.removeTrack(sender);
                        return true;
                    }
                });
            }    
        }

        return false;
    }

    /**
     * Function to initiate Screen Share track
     *
     * @returns
     */
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

    /**
     * Function to initiate Video Broadcasting
     *
     * @param {'broadcast' | 'alt-broadcast'} data
     * @returns
     */
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
                    streamId: this.localStream.id,
                })
            );
            this.log(`[startBroadcasting] send role`);
            return this.localStream;
        } catch (e) {
            this.updateTheStatus(`Error Start Broadcasting`);
            console.log(e);
            this.log(`[startBroadcasting] ${e}`);
            alert('Unable to get access to your webcam nor microphone.');
        }
    };

    /**
     * Function to intiate Listening to / Receiving of
     *
     * Video broadcast from broadcaster
     *
     */
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

    /**
     * Function to request to broadcast video
     *
     * and Immediately starts broadcasting
     *
     * @returns initiate Broadcasting
     */
    raiseHand = () => {
        if (this.startedRaiseHand) return;
        this.startedRaiseHand = true;
        return this.startBroadcasting('alt-broadcast');
    };

    /**
     * Function to handle Data Channel Status
     * @param {RTCDataChannel} dc
     * @param {String} target
     * @param {RTCPeerConnection} pc
     * And send messages via Data Channel
     */
    onDataChannelOpened(dc, target, pc) {

        console.log("DataChannel opened:", dc);

        let intervalId = setInterval(() => {

            if (dc.readyState === "open") {
                dc.send(`Hello from ${this.myName}`);

            } else if (dc.readyState === "connecting") {
                console.log("DataChannel is in the process of connecting.");
            } else if (dc.readyState === "closing") {
                console.log("DataChannel is in the process of closing.");
            } else if (dc.readyState === "closed") {
                console.log("DataChannel is closed and no longer able to send or receive data.");

                clearInterval(intervalId); //if closed leave the loop
            }


        }, this.sendMessageInterval);

    }


    /**
     * Function to restart the Negotiation and finding a new Parent
     *
     * @param {RTCPeerConnection} peerConnection - disconnected peer RTCPeerConnection Object
     * @param {String} target
     * @returns
     */
    restartEverything(peerConnection, target) {
        this.remoteStreamNotified = false;
        if (peerConnection.getRemoteStreams().length === 0) return;
        const trackIds = peerConnection.getReceivers().map((receiver) => receiver.track.id);
        trackIds.forEach((trackId) => {
            console.log('[peerConnection.oniceconnectionstatechange] DC trackId', trackId);
            for (const userId in this.myPeerConnectionArray) {
                if (userId === target) continue;
                console.log('[peerConnection.oniceconnectionstatechange] DC userId', userId);
                const apeerConnection = this.myPeerConnectionArray[userId];
                // if (!apeerConnection.isAdience) return;
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
        for (let i = 0; i < allStreams.length; i++) {
            while (this.remoteStreams.indexOf(allStreams[i]) >= 0) {
                this.remoteStreams.splice(this.remoteStreams.indexOf(allStreams[i]), 1);
            }
        }
        if (this.parentStreamId && allStreams.map((s) => s.id).includes(this.parentStreamId)) {
            this.updateTheStatus(`Parent stream is disconnected`);
            if (this.remoteStreamDCCallback) {
                this.remoteStreams.forEach((strm) => {
                    try {
                        this.remoteStreamDCCallback(strm);
                    } catch {
                    }
                });
            }
            this.parentStreamId = undefined;
            this.remoteStreams = [];
        }

        try {
            if (this.remoteStreamDCCallback) this.remoteStreamDCCallback(peerConnection.getRemoteStreams()[0]);
        } catch {
        }


        if (this.parentDC || this.startedRaiseHand) {
            setTimeout(() => {
                this.startProcedure();
            }, 1000);
        }
    }

    /**
     * Function to check Parent's status
     *
     * whether its connected or disconnected
     *
     * @param {RTCPeerConnection & {alive?: boolean}} pc
     * @param {String} target
     */
    checkParentDisconnection(pc, target) {
        // Check for disconnection of Parent
        let id = setInterval(() => {
            if (!pc.isAdience) {

                if (pc.alive != undefined) {
                    console.log("parent alive: ", pc.alive, "state: ", pc.connectionState);

                    if (!pc.alive) { //not connected and not alive
                        console.log("Parent disconnected");

                        this.parentDC = true;

                        //restart negotiation again
                        this.restartEverything(pc, target);

                        clearInterval(id); //if disconnected leave the loop
                    }
                    pc.alive = false;
                } else {
                    console.log("Undefined: ", pc.alive);
                }

            }
        }, this.parentDisconnectionTimeOut);
    }

    /**
     * Function to create new Peer connection
     *
     * And Data Channel with each peer connection
     *
     * @param {String} target
     * @param {Array<MediaStream>} theStream
     * @param {boolean} isAudience
     * @returns
     */
    newPeerConnectionInstance = (target, theStream, isAudience = false) => {
        this.log(`[newPeerConnectionInstance] target='${target}' theStream='${theStream}' isAudience='${isAudience}'`);
        /** @type {RTCPeerConnection & {_iceIsConnected?: boolean}} */
        const peerConnection = new RTCPeerConnection(this.myPeerConnectionConfig);
        let intervalId;

        peerConnection.isAdience = isAudience;
        peerConnection.alive = true;
        /*
                // Create DataChannel
                const dataChannel = peerConnection.createDataChannel("chat");


                // Handle open event for DataChannel
                dataChannel.onopen = (e) => {
                    this.onDataChannelOpened(dataChannel, target, peerConnection);
                }

                //callback for datachannel
                peerConnection.ondatachannel = event => {
                    let receive = event.channel;

                    receive.onmessage = e => {
                        //check if message came from Only My Parent
                        if (!peerConnection.isAdience) {
                            peerConnection.alive = true;
                        }
                    }


                    this.checkParentDisconnection(peerConnection, target);


                    //handle error event
                    receive.onerror = e => {
                        console.error("DataChannel error: ", e);
                    }

                    //handle beffer amount low event
                    receive.onbufferedamountlow = () => {
                        console.log("bufferedAmount dropped below threshold.");
                    }

                    // Handle close event for DataChannel
                    receive.onclose = e => {
                        console.log("DataChannel closed:", e);
                    };
                }
        */

        // Handle connectionstatechange event
        peerConnection.onconnectionstatechange = event => {
            console.log("Connection state:", peerConnection.connectionState);
        };

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
            }
        };

        peerConnection.ontrack = (event) => {
            this.updateStatus(`Peer Connection track received for ${target} stream ids [${event.streams.map((s) => s.id).join(',')}]`);
            this.parentDC = false;
            this.broadcasterDC = false;
            this.log(`[newPeerConnectionInstance] ontrack ${JSON.stringify(event.streams)}`);
            const stream = event.streams[0];
            console.log('user-by-stream', stream.id);
            this.socket.send(JSON.stringify({
                type: 'user-by-stream',
                data: stream.id,
            }));
            if (this.remoteStreams.length === 0) {
                this.parentStreamId = stream.id;
            }
            let ended = false;
            stream.getTracks().forEach((t) => {
                if (t.readyState === "ended") {
                    ended = true;
                }
            })
            if (ended) {
                console.log("stream tracks was ended", stream.id)
                return
            }
            stream.oninactive = (event) => {
                this.log(`[newPeerConnectionInstance] stream.oninactive ${JSON.stringify(event)}`);
                console.log('[stream.oninactive] event', event, event.currentTarget.getTracks(), event.target.getTracks());
                this.remoteStreamNotified = false;


                //func to remove RTP Sender
                const removeStream = (pc) => {
                    pc.getSenders().forEach(sender => {
                        const track = sender.track;
                        if(track){
                            if (track.kind === 'video' && track.muted === true) {
                                pc.removeTrack(sender);
                            }
                        } 
                    });
                };


                //Loop through peer connection array and find audinece PC
                for(const userid in this.myPeerConnectionArray){
                    if(this.myPeerConnectionArray[userid].isAdience){
                        removeStream(this.myPeerConnectionArray[userid]);
                    }
                }


                const theEventStream = event.currentTarget;
                const trackIds = theEventStream.getTracks().map((t) => t.id);

                for (const userId in this.myPeerConnectionArray) {
                    const apeerConnection = this.myPeerConnectionArray[userId];
                    if (!apeerConnection.isAdience) continue;
                    const allSenders = apeerConnection.getSenders();
                    for (const sender of allSenders) {
                        if (!sender.track) continue;
                        console.log("the streamId", this.trackToStreamMap[sender.track.id]);
                        if (this.trackToStreamMap[sender.track.id] === theEventStream.id) {
                            try {
                                console.log("track removed: ", sender.track.id,)
                                apeerConnection.removeTrack(sender);
                                // delete this.trackToStreamMap[sender.track.id];
                            } catch (e) {
                                console.log(e);
                            }
                        }
                    }
                }

                console.log('indx', this.remoteStreams.indexOf(theEventStream));
                while (this.remoteStreams.indexOf(theEventStream) >= 0) {
                    this.remoteStreams.splice(this.remoteStreams.indexOf(theEventStream), 1);
                }
                if (this.parentStreamId && this.parentStreamId === theEventStream.id) {
                    if (this.remoteStreamDCCallback) {
                        this.remoteStreams.forEach((strm) => {
                            this.remoteStreamDCCallback(strm);
                        });
                    }
                    this.parentStreamId = undefined;
                    this.parentDC = true;
                }
                if (this.remoteStreamDCCallback) {
                    try {
                        this.remoteStreamDCCallback(event.target);
                    } catch {
                    }
                }
            };
            try {
                if (this.remoteStreamCallback)
                    this.remoteStreamCallback(stream);
            } catch {
            }
            this.remoteStreams.push(stream);
            stream.getTracks().forEach((t) => {
                this.trackToStreamMap[t.id] = stream.id;
            })
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

                const apeerConnection = this.myPeerConnectionArray[userId];
                this.updateStatus(`check Sending the stream [${stream.id}] tracks to ${userId} ${apeerConnection.isAdience.toString()}`);
                if (!apeerConnection.isAdience) continue;

                this.updateStatus(`Sending the stream [${stream.id}] tracks to ${userId}`);
                stream.getTracks().forEach((track) => {
                    try {
                        track.streamId = stream.id;
                        apeerConnection.addTrack(track, stream);
                    } catch {
                    }
                });
            }

            if (!this.started) {
                this.started = true;
            }
        };


        peerConnection.oniceconnectionstatechange = (event) => {
            this.log(`[newPeerConnectionInstance] oniceconnectionstatechange peerConnection.iceConnectionState = ${peerConnection.iceConnectionState} event = ${JSON.stringify(event)}`);
            console.log(`new ice connection state => ${peerConnection.iceConnectionState}`);
            switch (peerConnection.iceConnectionState) {
                case "connected": {
                    peerConnection._iceIsConnected = true;
                    break;
                }
                default:
                    peerConnection._iceIsConnected = false;
                    break;
            }
            if (peerConnection.iceConnectionState === 'disconnected' || peerConnection.iceConnectionState === 'failed' || peerConnection.iceConnectionState === 'closed') {
                if (!this.parentDC)
                    setTimeout(() => {
                        console.log("restarting ice");
                        peerConnection.restartIce();
                    }, 0);
                this.restartEverything(peerConnection, target);
            }
        };

        setTimeout(() => {
            if (!peerConnection._iceIsConnected) {
                peerConnection.restartIce();
            }
        }, 4000);

        return peerConnection;
    };

    /**
     * Helper fucntion to iniiate select
     *
     * Whether to create a new peer connection with [peerName] Or to Get the existing one with [peerName]
     *
     * @param {String} audienceName
     * @param {boolean} isAudience
     * @returns
     */
    createOrGetPeerConnection = (audienceName, isAudience = false) => {
        this.log(`[createOrGetPeerConnection] audienceName = ${audienceName}, isAudience = ${isAudience}`);
        if (this.myPeerConnectionArray[audienceName]) return this.myPeerConnectionArray[audienceName];

        this.myPeerConnectionArray[audienceName] = this.newPeerConnectionInstance(audienceName, true, isAudience);
        this.log(`[createOrGetPeerConnection] generate newPeerConnectionInstance`);

        return this.myPeerConnectionArray[audienceName];
    };

    /**
     * Function to add new Audience as Current Node's Children
     * @param {String} audienceName
     */
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
                    try {
                        this.myPeerConnectionArray[audienceName].addTrack(track, astream);
                    } catch {
                    }
                });
            });
        }
    };

    /**
     * Function to spread the local stream to Target Audiance peers.
     *
     * It peer connection exists it send to it, if not it creat a new one and send stream to it
     *
     * @param {String} target
     * @param {MediaStream} stream
     */
    sendStreamTo = (target, stream) => {
        this.log(`[handleMessage] sendStreamTo ${target}`);
        const peerConnection = this.createOrGetPeerConnection(target, false);
        stream.getTracks().forEach((track) => {
            if (this.lastVideoState === "Disabled") {
                this.disableVideo();
                let img = document.getElementById("camera");
                img.dataset.status = 'off';
                img.src = CAMERA_OFF;
            }
            if (this.lastAudioState === "Disabled") {
                this.disableAudio();
                let img = document.getElementById("mic");
                img.dataset.status = 'off';
                img.src = MIC_OFF;
            }

            peerConnection.addTrack(track, stream);
        });
    };

    /**
     * Function to initiate the client depending on its role
     *
     * If role is broadcaster start Broadcasting
     *
     * otherwise start listening to Broadcast
     *
     * @param {boolean} turn
     * @returns
     */
    start = async (turn = true) => {
        if (!turn) {
            this.myPeerConnectionConfig.iceServers = iceServers.filter((i) => i.url.indexOf('turn') < 0);
        }
        if (this.startedRaiseHand) {
            console.log("its true, calling it")
            setTimeout(() => {
                this.startedRaiseHand = false;
                this.raiseHand();
                document.getElementById('camera').style.display = '';
                document.getElementById('mic').style.display = '';
            }, 2000);
        }
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

    /**
     * Function to enable / disable Video track
     *
     * @param {boolean} enabled
     */
    disableVideo = (enabled = false) => {
        this.lastVideoState = enabled === true ? "Enabled" : "Disabled";
        this.localStream.getTracks().forEach((track) => {
            if (track.kind === 'video')
                track.enabled = enabled;
        });
    };

    /**
     * Function to enable / disable Audio track
     *
     * @param {boolean} enabled
     */
    disableAudio = (enabled = false) => {
        this.lastAudioState = enabled === true ? "Enabled" : "Disabled";
        this.localStream.getTracks().forEach((track) => {
            if (track.kind === 'audio')
                track.enabled = enabled;
        });
    };

    /**
     * Function to WAIT for 1 second
     *
     * @param {number} mil
     * @returns
     */
    wait = async (mil = 1000) => {
        return new Promise((res) => {
            setTimeout(() => {
                res();
            }, mil);
        });
    };
    /**
     * Function to get Broadcaster status from backend
     *
     * @returns
     */
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

    /**
     * Function to set the Peer Connection Constraints
     *
     * depending upon presence Audio and Video Devices
     */
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

    /**
     * Function to update the Status of current Client
     */
    updateTheStatus = (status) => {
        if (this.updateStatus) {
            try {
                this.updateStatus(status);
            } catch {
            }
        }
    };

    /**
     * Function to lower hand and take request(to broadcast) back, if sharing already stop sharing
     *
     * @returns
     */
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
        this.localStream.getTracks().forEach((track) => {
            track.stop();
        });
        this.localStream = null;
        this.startedRaiseHand = false;
    };

    /**
     * Function to broadcast local stream to all peer audiance
     */
    spreadLocalStream = () => {
        for (const target in this.myPeerConnectionArray) {
            if (this.myPeerConnectionArray[target].isAdience)
                this.sendStreamTo(target, this.localStream);
        }
    };
    getStreamDetails = (streamId) => {
        if (this.localStream && this.localStream.id === streamId) {
            return {
                userId: this.myUsername,
                userName: this.myName,
                userRole: this.role,
            };
        }

        if (this.userStreamData[streamId]) {
            return this.userStreamData[streamId];
        }

        return null;
    };
    setMetadata = (metadata) => {
        this.socket.send(JSON.stringify({
            type: 'metadata-set',
            data: JSON.stringify(metadata)
        }));
    };
    getMetadata = () => {
        this.socket.send(JSON.stringify({
            type: 'metadata-get',
        }));
    };
    streamById = (streamId) => {
        return this.remoteStreams.find((s) => s.id === streamId);
    };
    stopSignaling = () => {
        if (this.pingInterval) {
            clearInterval(this.pingInterval);
            this.pingInterval = null;
        }
        this.socket.onclose = () => {
        };
        this.socket.close();
    };

    /**
     * Construcor Function for Class SparkRTC
     *
     * @param {*} role
     * @param {*} options
     */
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
        this.log = options.log || ((log) => {
        });
        this.constraintResults = options.constraintResults;
        this.log(`[constructor] ${this.role}`);
        this.updateStatus = options.updateStatus;
        this.userListUpdated = options.userListUpdated;
    }
}
