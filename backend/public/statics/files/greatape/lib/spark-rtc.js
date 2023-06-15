/** Your class description
 *
 * SparkRTC class is main class to setUP RTC client
 *
 */
export class SparkRTC {
    started = false;
    maxRaisedHands = 2;
    myPeerConnectionConfig = {
        iceServers,
    };
    role = 'broadcast';
    localStream;
    shareStream;
    socketURL = '';
    remoteStreamNotified = false;
    remoteStreams = [];
    socket;
    myName = 'NoName';
    roomName = 'SparkRTC';
    myUsername = 'NoUsername';
    lastBroadcasterId = '';
    broadcastingApproved = false;
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
    leftMeeting = false;

    userListCallback = null;
    remoteStreamsQueue = [];

    parentDisconnectionTimeOut = 2000; // 2 second timeout to check parent is alive or not
    sendMessageInterval = 1; // send message to child after every 10 ms
    metaData = {};
    userStreamData = {};
    users = [];

    /**@type {{[trackId:string]: string}}*/
    trackToStreamMap = {};
    /**@type {"Enabled" | "Disabled"}*/
    lastVideoState = 'Enabled';
    /**@type {"Enabled" | "Disabled"}*/
    lastAudioState = 'Enabled';

    broadcastersMessage = null;

    //enum for Roles
    Roles = {
        BROADCAST: 'broadcast',
        AUDIENCE: 'audience',
        BROADCASTER: 'broadcaster',
    };

    //enum for stream kind
    StreamType = {
        SCREEN: 'screen',
        CAMERA: 'camera',
    };

    enqueue(queue, data) {
        queue.push(data);
    }

    dequeue(queue) {
        if (queue.length === 0) {
            return null;
        }
        return queue.shift();
    }

    /**
     * Function to handle Peer Connection Offer, received from Other Peer
     *
     * and Return Peer Connection Answer to Other Peer
     *
     * @param {{name: string, sdp: RTCSessionDescriptionInit }} msg
     */
    handleVideoOfferMsg = async (msg) => {
        this.updateTheStatus(`[handleVideoOfferMsg]`, msg.name);
        const broadcasterPeerConnection = this.createOrGetPeerConnection(
            msg.name
        );
        await broadcasterPeerConnection.setRemoteDescription(
            new RTCSessionDescription(msg.sdp)
        );
        await broadcasterPeerConnection.setLocalDescription(
            await broadcasterPeerConnection.createAnswer()
        );

        if (await this.checkSocketStatus())
            this.socket.send(
                JSON.stringify({
                    name: this.myUsername,
                    target: msg.name,
                    type: 'video-answer',
                    sdp: broadcasterPeerConnection.localDescription,
                })
            );
        this.updateTheStatus(
            `[handleVideoOfferMsg] send video-answer to ${msg.name} from ${this.myUsername}`
        );
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
        let msg;
        try {
            msg = JSON.parse(event.data);
        } catch (e) {
            return;
        }
        msg.data = msg.Data && !msg.data ? msg.Data : msg.data;
        msg.type = msg.Type && !msg.type ? msg.Type : msg.type;

        let audiencePeerConnection;
        switch (msg.type) {
            case 'video-offer':
            case 'alt-video-offer':
                this.updateTheStatus(
                    `[handleMessage] handleVideoOfferMsg ${msg.type}`
                );
                this.handleVideoOfferMsg(msg);
                break;
            case 'video-answer':
            case 'alt-video-answer':
                this.updateTheStatus(
                    `[handleMessage] setRemoteDescription ${msg.type}`
                );
                audiencePeerConnection = this.createOrGetPeerConnection(
                    msg.data,
                    true
                );
                try {
                    await audiencePeerConnection.setRemoteDescription(
                        new RTCSessionDescription(msg.sdp)
                    );
                } catch (e) {
                    this.updateTheStatus(
                        `setRemoteDescription failed with exception: ${e.message}`
                    );
                }

                //check ice candidates list
                if (
                    this.iceCandidates &&
                    this.iceCandidates.length > 0 &&
                    audiencePeerConnection &&
                    audiencePeerConnection.remoteDescription
                ) {
                    try {
                        while (this.iceCandidates.length) {
                            await audiencePeerConnection.addIceCandidate(
                                this.iceCandidates.pop()
                            );
                        }
                    } catch (e) {
                        this.updateTheStatus(
                            `iceCandidateError-alt-video-answer`,
                            e
                        );
                    }
                }
                break;
            case 'new-ice-candidate':
            case 'alt-new-ice-candidate':
                this.updateTheStatus(
                    `[handleMessage] addIceCandidate ${msg.type}`
                );
                audiencePeerConnection = this.createOrGetPeerConnection(
                    msg.data
                );
                if (msg.candidate) {
                    this.iceCandidates.push(new RTCIceCandidate(msg.candidate));
                }
                if (
                    audiencePeerConnection &&
                    audiencePeerConnection.remoteDescription &&
                    this.iceCandidates &&
                    this.iceCandidates.length > 0
                ) {
                    try {
                        while (this.iceCandidates.length) {
                            await audiencePeerConnection.addIceCandidate(
                                this.iceCandidates.pop()
                            );
                        }
                    } catch (e) {
                        this.updateTheStatus(
                            `iceCandidateError-alt-new-ice-candidate`,
                            e
                        );
                    }
                }
                break;
            case 'role':
                this.updateTheStatus(`[handleMessage] role:`, msg);
                if (this.role === this.Roles.BROADCAST) {
                    if (msg.data === 'no:broadcast') {
                        alert('You are not a broadcaster anymore!');
                        this.socket.close();
                    } else if (msg.data === 'yes:broadcast') {
                        this.updateTheStatus(`myName:`, this.myName);

                        const data = JSON.parse(this.myName);
                        this.localStream.name = data.name;

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
                            this.updateTheStatus(e);
                        }
                    }
                }
                break;
            case 'start':
                this.updateTheStatus(
                    `[handleMessage] start`,
                    JSON.stringify(msg)
                );
                if (msg.error) {
                    alert(msg.error);
                    return;
                }

                this.myUsername = msg.data;

                if (this.userInitialized) this.userInitialized(msg.data);

                break;
            case 'add_audience':
            case 'add_broadcast_audience':
                this.updateTheStatus(`[handleMessage] add audience ${msg}`);
                this.updateTheStatus(`New Audience arrived ${msg.data}`);
                this.connectToAudience(msg.data);
                break;
            case 'alt-broadcast-approve':
                this.updateTheStatus(
                    `[handleMessage] alt-broadcast-approve`,
                    msg
                );

                if (msg.maxLimitReached) {
                    this.localStream
                        ?.getTracks()
                        ?.forEach((track) => track.stop());
                    this.localStream = null;
                    this.startedRaiseHand = false;
                    this.broadcastingApproved = false;

                    if (this.maxLimitReached) {
                        this.maxLimitReached(
                            'Max limit of 2 Broadcasting Audiences is Reached'
                        );
                    }
                } else {
                    this.broadcastingApproved = true;

                    if (msg.result == true) {
                        this.localStream = null;
                        await this.startBroadcasting('alt-broadcast');
                        this.lastBroadcasterId = msg.data;
                        if (this.localStream) {
                            this.sendStreamTo(msg.data, this.localStream);
                        }
                    }

                    if (this.altBroadcastApprove) {
                        this.altBroadcastApprove(msg.result);
                    }
                }
                break;
            case 'alt-broadcast':
                this.updateTheStatus(`[handleMessage] ${msg.type}`);
                if (this.role === this.Roles.BROADCAST) {
                    var limitReached = false;

                    if (this.raiseHands.length >= this.maxRaisedHands) {
                        limitReached = true;
                    }
                    this.updateTheStatus(`My ID: ${this.myUsername}`);

                    if (this.raiseHands.indexOf(msg.data) === -1) {
                        var result = false;
                        if (this.raiseHandConfirmation && !limitReached) {
                            try {
                                const data = JSON.parse(msg.name);
                                const name = data.name;
                                const email = data.email;

                                result = await this.raiseHandConfirmation({
                                    name,
                                    email,
                                    userId: msg.Data,
                                });
                                this.updateTheStatus(
                                    `[handleMessage] alt-broadcast result ${result}`
                                );
                            } catch (e) {
                                console.error(e);
                                return;
                            }
                        }

                        if (await this.checkSocketStatus())
                            this.socket.send(
                                JSON.stringify({
                                    type: 'alt-broadcast-approve',
                                    target: msg.data,
                                    result,
                                    maxLimitReached: limitReached,
                                })
                            );

                        if (result !== true) return;

                        this.raiseHands.push(msg.data);
                        this.updateTheStatus(
                            `[handleMessage] ${msg.type} approving raised hand`,
                            msg.data
                        );
                        this.getMetadata();
                        setTimeout(() => {
                            const metaData = this.metaData;
                            metaData.raiseHands = JSON.stringify(
                                this.raiseHands
                            );
                            this.setMetadata(metaData);
                        }, 1000);
                    } else {
                        this.updateTheStatus(`else of this.raiseHands`);
                    }
                } else {
                    this.updateTheStatus(`else of role check`);
                    // this.spreadLocalStream();
                }
                break;
            case 'tree':
                this.updateTheStatus(`[handleMessage] ${msg.type}`);
                if (this.treeCallback) this.treeCallback(msg.data);
                break;
            case 'broadcasting':
                if (this.role === this.Roles.BROADCAST) return;
                this.broadcasterDC = false;
                this.updateTheStatus(`[handleMessage] ${msg.type}`);
                this.startProcedure(true);
                break;
            case 'event-reconnect':
            case 'event-broadcaster-disconnected':
                this.updateTheStatus(`broadcaster dc ${msg.type}`);
                this.broadcasterDC = true;
                const broadcasterId = this.broadcasterUserId();

                this.socket.send(
                    JSON.stringify({
                        type: 'stream',
                        data: 'false',
                    })
                );
                for (const u in this.myPeerConnectionArray) {
                    this.myPeerConnectionArray[u].close();
                }
                this.myPeerConnectionArray = {};
                this.remoteStreams = [];
                try {
                    if (this.remoteStreamDCCallback)
                        this.remoteStreamDCCallback('no-stream');
                } catch {}
                this.localStream?.getTracks()?.forEach((track) => track.stop());
                this.localStream = null;
                this.startedRaiseHand = false;
                break;
            case 'event-parent-dc':
                this.updateTheStatus(`parentDC ${msg.type}`);
                this.parentDC = true;

                if (this.startedRaiseHand) {
                    this.updateTheStatus(`Reloading..`);
                    if (this.parentDcMessage) this.parentDcMessage(); //show parent dc message on Audience side
                    await this.wait();
                    window.location.reload();
                } else {
                    this.updateTheStatus(`Only Starting Procedure..`);
                    this.startProcedure(true);
                }
                break;
            case 'metadata-get':
            case 'metadata-set':
                this.updateTheStatus(`[handleMessage] ${msg.type}`);
                this.metaData = JSON.parse(msg.data);
                if (this.metaData.raiseHands) {
                    //this.raiseHands=JSON.parse(this.metaData.raiseHands);
                }
                break;
            case 'user-by-stream':
                this.updateTheStatus(`[handleMessage] ${msg.type}`, msg.data);
                const [userId, userName, streamId, userRole] =
                    msg.Data.split(',');
                this.userStreamData[streamId] = {
                    userId,
                    userName,
                    userRole,
                };
                break;
            case 'user-event':
                this.updateTheStatus(`[handleMessage] ${msg.type}`);
                this.getMetadata();
                setTimeout(() => {
                    const users = JSON.parse(msg.data).map((u) => {
                        this.updateTheStatus(u);
                        const video =
                            u.streamId !== ''
                                ? this.streamById(u.streamId)
                                : null;
                        return {
                            id: u.id,
                            name: u.name,
                            role: u.role,
                            video,
                        };
                    });
                    this.users = users;

                    if (this.userListCallback) {
                        this.userListCallback(users);
                    } else {
                        this.updateTheStatus(`No Callback registered`);
                    }

                    if (this.userListUpdated) {
                        try {
                            this.userListUpdated(users);
                        } catch {}
                    }
                }, 1000);
                break;

            case 'disable-audience':
                this.updateTheStatus(`disable-audience:`, msg);

                if (this.disableBroadcasting) {
                    this.disableBroadcasting();
                }
                break;

            default:
                // this.updateTheStatus(
                //     `[handleMessage] default ${JSON.stringify(msg)}`
                // );
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
    ping = async () => {
        if (await this.checkSocketStatus())
            this.socket.send(
                JSON.stringify({
                    type: this.treeCallback ? 'tree' : 'ping',
                })
            );
    };

    disableAudienceBroadcast = async (target) => {
        this.updateTheStatus(`disableAudienceBroadcast:`, target);

        //find media stream id of that Target
        for (const id in this.myPeerConnectionArray) {
            if (this.myPeerConnectionArray[id].isAdience) {
                if (id.toString() === target.toString()) {
                    var pc = this.myPeerConnectionArray[id];
                    this.updateTheStatus(pc);
                }
            }
        }

        if (await this.checkSocketStatus())
            this.socket.send(
                JSON.stringify({
                    type: 'disable-audience',
                    target: target,
                })
            );
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
        this.updateTheStatus(
            `[setupSignalingSocket] url=${url} myName=${myName} roomName=${roomName}`
        );
        return new Promise((resolve, reject) => {
            if (this.pingInterval) {
                clearInterval(this.pingInterval);
                this.pingInterval = null;
            }
            if (myName) this.myName = myName;
            if (roomName) this.roomName = roomName;

            this.updateTheStatus(`[setupSignalingSocket] installing socket`);
            this.socketURL = url + '?room=' + this.roomName;
            const socket = new WebSocket(this.socketURL);
            socket.onmessage = this.handleMessage;
            socket.onopen = () => {
                socket.send(
                    JSON.stringify({
                        type: 'start',
                        data: myName,
                    })
                );
                this.pingInterval = setInterval(this.ping, 5000);
                this.updateTheStatus(
                    `[setupSignalingSocket] socket onopen and sent start`
                );
                resolve(socket);
            };
            socket.onclose = async () => {
                this.updateTheStatus(
                    `socket is closed in setupSignalingSocket`
                );
                this.remoteStreamNotified = false;
                this.myPeerConnectionArray = {};
                this.started = false;
                if (this.startProcedure && !this.leftMeeting) {
                    this.startProcedure();
                }
            };
            socket.onerror = (error) => {
                this.updateTheStatus(
                    `WebSocket error in setupSignalingSocket`,
                    error
                );
                reject(error);
                if (!this.leftMeeting) {
                    window.location.reload(); //reload before, alert because alert blocks the reload
                    alert('Can not connect to server');
                }
            };

            this.socket = socket;
        });
    };

    /**
     * To check the web socket's status
     * @returns
     */
    async checkSocketStatus() {
        return true;

        //todo need to improve this logic
        if (
            this.socket.readyState === WebSocket.CLOSED ||
            this.socket.readyState === WebSocket.CLOSING
        ) {
            this.startProcedure(); //restart
            return false;
        }

        if (this.socket.readyState === WebSocket.OPEN) {
            return true;
        }

        if (this.socket.readyState === WebSocket.CONNECTING) {
            //sleep for few seconds
            await this.wait(2000);
            return true;
        }

        return false;
    }

    stopShareScreen = async (stream) => {
        if (stream) {
            for (const userId in this.myPeerConnectionArray) {
                const apeerConnection = this.myPeerConnectionArray[userId];

                stream.getTracks().forEach((track) => {
                    const sender = apeerConnection
                        .getSenders()
                        .find(
                            (sender) =>
                                sender.track && sender.track.id === track.id
                        );

                    if (sender) {
                        apeerConnection.removeTrack(sender);
                        track.stop();
                    }
                });
            }

            //remove shared screen stream from remotestreams list
            this.remoteStreams.forEach((STR) => {
                if (STR.id === stream.id) {
                    this.updateTheStatus(`ids matched`);

                    let newArray = this.remoteStreams.filter(
                        (str) => str !== STR
                    );

                    this.remoteStreams = newArray;
                }
            });
        }
    };

    /**
     * Function to initiate Screen Share track
     *
     * @returns
     */
    startShareScreen = async () => {
        this.updateTheStatus(`[handleMessage] startShareScreen`);
        try {
            this.shareStream = await navigator.mediaDevices.getDisplayMedia({
                audio: true,
                video: true,
            });

            this.remoteStreams.push(this.shareStream);

            for (const userId in this.myPeerConnectionArray) {
                const apeerConnection = this.myPeerConnectionArray[userId];
                this.shareStream.getTracks().forEach((track) => {
                    apeerConnection.addTrack(track, this.shareStream);
                });
            }

            //add name to stream
            const data = JSON.parse(this.myName);
            this.shareStream.name = data.name;

            return this.shareStream;
        } catch (e) {
            console.error(e);
            this.updateTheStatus(`[handleMessage] startShareScreen error ${e}`);
            alert('Unable to get access to screenshare.');
        }
    };

    /**
     * Function to initiate Video Broadcasting
     *
     * @param {'broadcast' | 'alt-broadcast'} data
     * @returns
     */
    startBroadcasting = async (data = this.Roles.BROADCAST) => {
        this.updateTheStatus(`[startBroadcasting]`, data);
        try {
            if (!this.localStream) {
                this.updateTheStatus(`Trying to get local stream`);
                if (!this.constraints.audio && !this.constraints.video) {
                    this.updateTheStatus(`No media device available`);
                    throw new Error('No media device available');
                }
                this.localStream = await navigator.mediaDevices.getUserMedia(
                    this.constraints
                );
                this.updateTheStatus(`Local stream loaded`);
                this.updateTheStatus(`[startBroadcasting] local stream loaded`);
                this.remoteStreams.push(this.localStream);
            }
            this.updateTheStatus(`Request Broadcast Role`);

            if (await this.checkSocketStatus())
                this.socket.send(
                    JSON.stringify({
                        type: 'role',
                        data,
                        streamId: this.localStream.id,
                    })
                );
            this.updateTheStatus(`[startBroadcasting] send role`);
            return this.localStream;
        } catch (e) {
            this.updateTheStatus(`Error Start Broadcasting`);
            this.updateTheStatus(e);
            this.updateTheStatus(`[startBroadcasting] ${e}`);
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
        this.updateTheStatus(`[startReadingBroadcast]`);
        this.updateTheStatus(`Request Audience Role`);
        if (await this.checkSocketStatus())
            this.socket.send(
                JSON.stringify({
                    type: 'role',
                    data: this.Roles.AUDIENCE,
                })
            );
        this.updateTheStatus(`[startReadingBroadcast] send role audience`);
    };

    /**
     * Function to request to broadcast video
     *
     * and Immediately starts broadcasting
     *
     * @returns initiate Broadcasting
     */
    raiseHand = async () => {
        if (this.startedRaiseHand) return;
        this.startedRaiseHand = true;
        // send a raise hand request with empty streamId, no local stream at that Moment
        if (await this.checkSocketStatus())
            this.socket.send(
                JSON.stringify({
                    type: 'role',
                    data: 'alt-broadcast',
                    streamId: '',
                })
            );
    };

    async getLatestUserList() {
        this.updateTheStatus(`Request to fetch Latest UserList`);
        if (await this.checkSocketStatus())
            this.socket.send(
                JSON.stringify({
                    type: 'get-latest-user-list',
                })
            );
    }

    onRaiseHandRejected = () => {
        this.startedRaiseHand = false;
        this.broadcastingApproved = false;

        const pc = this.myPeerConnectionArray[this.lastBroadcasterId];
        if (this.localStream) {
            //remove local stream from list of remote streams
            this.remoteStreams.forEach((stream) => {
                if (stream.id === this.localStream.id) {
                    this.updateTheStatus(`ids matched`);

                    let newArray = this.remoteStreams.filter(
                        (str) => str !== stream
                    );
                    this.remoteStreams = newArray;
                }
            });

            this.localStream.getTracks().forEach((track) => {
                if (pc && pc.getSenders) {
                    pc.getSenders().forEach((sender) => {
                        if (track.id === sender?.track?.id) {
                            pc.removeTrack(sender);
                        }
                    });
                }
            });

            this.localStream.getTracks().forEach(function (track) {
                track.stop();
            });

            this.localStream = null;
        }
    };

    /**
     * Function to handle Data Channel Status
     * @param {RTCDataChannel} dc
     * @param {String} target
     * @param {RTCPeerConnection} pc
     * And send messages via Data Channel
     */
    onDataChannelOpened(dc, pc) {
        this.updateTheStatus(`DataChannel opened: ${dc}`);

        let intervalId = setInterval(() => {
            if (dc.readyState === 'open') {
                if (pc.isAdience) {
                    if (this.role === this.Roles.BROADCAST) {
                        if (this.shareStream) {
                            let message = {
                                id: this.shareStream.id,
                                type: this.StreamType.SCREEN,
                                name: this.myName,
                            };

                            dc.send(JSON.stringify(message));
                        }
                    } else {
                        dc.send(this.broadcastersMessage);
                    }
                }
            } else if (dc.readyState === 'connecting') {
                this.updateTheStatus(
                    'DataChannel is in the process of connecting.'
                );
            } else if (dc.readyState === 'closing') {
                this.updateTheStatus(
                    'DataChannel is in the process of closing.'
                );
            } else if (dc.readyState === 'closed') {
                this.updateTheStatus(
                    'DataChannel is closed and no longer able to send or receive data.'
                );

                clearInterval(intervalId); //if closed leave the loop
            }
        }, this.sendMessageInterval);
    }

    /**
     * Function to restart the Negotiation and finding a new Parent
     *
     * @param {RTCPeerConnection} peerConnection - disconnected peer RTCPeerConnection Object
     * @param {String} target
     * @param {Boolean} isAudience
     * @returns
     */
    restartEverything(peerConnection, target, isAudience) {
        this.remoteStreamNotified = false;
        //if (peerConnection.getRemoteStreams().length === 0) return;
        const trackIds = peerConnection
            .getReceivers()
            .map((receiver) => receiver.track.id);
        trackIds.forEach((trackId) => {
            this.updateTheStatus(
                `[peerConnection.oniceconnectionstatechange] DC trackId ${trackId}`
            );
            for (const userId in this.myPeerConnectionArray) {
                if (userId === target) continue;
                this.updateTheStatus(
                    `[peerConnection.oniceconnectionstatechange] DC userId ${userId}`
                );
                const apeerConnection = this.myPeerConnectionArray[userId];
                //if (!apeerConnection.isAdience) return;
                const allSenders = apeerConnection.getSenders();
                for (const sender of allSenders) {
                    if (!sender.track) continue;
                    if (sender.track.id === trackId) {
                        this.updateTheStatus(
                            `[peerConnection.oniceconnectionstatechange] DC sender`
                        );
                        try {
                            apeerConnection.removeTrack(sender);
                        } catch (e) {
                            this.updateTheStatus(e);
                        }
                    }
                }
            }
        });
        const allStreams = peerConnection.getRemoteStreams();
        this.updateTheStatus(`All Remote streams of PC`, { allStreams });
        this.updateTheStatus(
            `All Remotestreams from List before`,
            this.remoteStreams
        );

        for (let i = 0; i < allStreams.length; i++) {
            var id = allStreams[i].id;
            var index = this.remoteStreams.indexOf(allStreams[i]); // Find the index of the element

            //remove all remote streams of peer connection using stream id, to avoid empty video container
            this.remoteStreams.forEach((stream) => {
                if (stream.id === id) {
                    this.updateTheStatus(`ids matched`);

                    let newArray = this.remoteStreams.filter(
                        (str) => str !== stream
                    );

                    this.remoteStreams = newArray;
                }
            });

            this.updateTheStatus(
                `remoteStreamsLength: `,
                this.remoteStreams.length
            );
        }
        this.updateTheStatus(`All Remotestreams from List`, this.remoteStreams);

        if (
            this.parentStreamId &&
            allStreams.map((s) => s.id).includes(this.parentStreamId)
        ) {
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

        try {
            if (this.remoteStreamDCCallback)
                this.remoteStreamDCCallback(
                    peerConnection.getRemoteStreams()[0]
                );
        } catch {}

        if (
            (this.parentDC || !isAudience) &&
            this.role !== this.Roles.BROADCAST
        ) {
            this.updateTheStatus(`Waiting to restart..`);
            setTimeout(() => {
                this.startProcedure(true);
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
                    this.updateTheStatus(
                        `parent alive: ${pc.alive}, state: ${pc.connectionState}`
                    );

                    if (!pc.alive) {
                        //not connected and not alive
                        this.updateTheStatus('Parent disconnected');

                        this.parentDC = true;

                        //restart negotiation again
                        this.restartEverything(pc, target);

                        clearInterval(id); //if disconnected leave the loop
                    }
                    pc.alive = false;
                } else {
                    this.updateTheStatus(`Undefined:  ${pc.alive}`);
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
        this.updateTheStatus(
            `[newPeerConnectionInstance] target='${target}' theStream='${theStream}' isAudience='${isAudience}'`
        );
        /** @type {RTCPeerConnection & {_iceIsConnected?: boolean}} */
        const peerConnection = new RTCPeerConnection(
            this.myPeerConnectionConfig
        );
        let intervalId;

        peerConnection.isAdience = isAudience;
        peerConnection.alive = true;

        // Create DataChannel
        const dataChannel = peerConnection.createDataChannel('chat');

        // Handle open event for DataChannel
        dataChannel.onopen = (e) => {
            this.onDataChannelOpened(dataChannel, peerConnection);
        };

        //callback for datachannel
        peerConnection.ondatachannel = (event) => {
            let receive = event.channel;

            receive.onmessage = (e) => {
                this.broadcastersMessage = e.data;
            };

            //handle error event
            receive.onerror = (e) => {
                console.error('DataChannel error: ', e);
            };

            //handle beffer amount low event
            receive.onbufferedamountlow = () => {
                this.updateTheStatus('bufferedAmount dropped below threshold.');
            };

            // Handle close event for DataChannel
            receive.onclose = (e) => {
                this.updateTheStatus('DataChannel closed:', e);
            };
        };

        // Handle connectionstatechange event
        peerConnection.onconnectionstatechange = (event) => {
            this.updateTheStatus(
                `Connection state: ${peerConnection.connectionState}`
            );

            if (this.connectionStatus) {
                this.connectionStatus(peerConnection.connectionState);
            }
        };

        peerConnection.onicecandidate = async (event) => {
            this.updateTheStatus(
                `Peer Connection ice candidate arrived for ${target}: event.candidate='${JSON.stringify(
                    event.candidate
                )}'`
            );
            if (event.candidate) {
                if (await this.checkSocketStatus())
                    this.socket.send(
                        JSON.stringify({
                            type: 'new-ice-candidate',
                            candidate: event.candidate,
                            target,
                        })
                    );
            }
        };

        peerConnection.onnegotiationneeded = async () => {
            this.updateTheStatus(
                `Peer Connection negotiation needed for ${target} preparing video offer`
            );
            try {
                await peerConnection.setLocalDescription(
                    await peerConnection.createOffer()
                );
                if (await this.checkSocketStatus())
                    this.socket.send(
                        JSON.stringify({
                            type: 'video-offer',
                            sdp: peerConnection.localDescription,
                            target,
                            name: this.myUsername,
                        })
                    );
            } catch (e) {
                this.updateTheStatus(`[newPeerConnectionInstance] failed ${e}`);
            }
        };

        peerConnection.ontrack = async (event) => {
            this.updateTheStatus(
                `Peer Connection track received for ${target} stream ids [${event.streams
                    .map((s) => s.id)
                    .join(',')}]`
            );
            this.parentDC = false;
            this.broadcasterDC = false;
            this.updateTheStatus(
                `[newPeerConnectionInstance] ontrack ${JSON.stringify(
                    event.streams
                )}`
            );
            const stream = event.streams[0];

            this.updateTheStatus(`user-by-stream ${stream.id}`);
            if (await this.checkSocketStatus())
                this.socket.send(
                    JSON.stringify({
                        type: 'user-by-stream',
                        data: stream.id,
                    })
                );
            if (this.remoteStreams.length === 0) {
                this.parentStreamId = stream.id;
            }
            let ended = false;
            stream.getTracks().forEach((t) => {
                if (t.readyState === 'ended') {
                    ended = true;
                }
            });
            if (ended) {
                this.updateTheStatus(`stream tracks was ended ${stream.id}`);
                return;
            }
            stream.oninactive = (event) => {
                this.updateTheStatus(`inactiveStream `, stream);
                this.updateTheStatus(`currentTarget `, event.currentTarget);

                this.updateTheStatus(
                    `[newPeerConnectionInstance] stream.oninactive ${JSON.stringify(
                        event
                    )}`
                );
                this.updateTheStatus(`[stream.oninactive] event `, event);
                this.updateTheStatus(`targetTracks`, event.target.getTracks());

                this.remoteStreamNotified = false;

                const theEventStream = event.currentTarget;
                const trackIds = theEventStream.getTracks().map((t) => t.id);

                for (const userId in this.myPeerConnectionArray) {
                    const apeerConnection = this.myPeerConnectionArray[userId];
                    //if (!apeerConnection.isAdience) continue;
                    const allSenders = apeerConnection.getSenders();
                    for (const sender of allSenders) {
                        if (!sender.track) continue;
                        this.updateTheStatus(
                            `the streamId`,
                            this.trackToStreamMap[sender.track.id]
                        );
                        if (
                            this.trackToStreamMap[sender.track.id] ===
                            theEventStream.id
                        ) {
                            try {
                                apeerConnection.removeTrack(sender);
                                // delete this.trackToStreamMap[sender.track.id];
                            } catch (e) {
                                this.updateTheStatus(e);
                            }
                        }
                    }
                }

                this.updateTheStatus(
                    `indx`,
                    this.remoteStreams.indexOf(theEventStream)
                );

                //remove the event stream from remotestreamslist
                this.remoteStreams.forEach((stream) => {
                    if (stream.id === theEventStream.id) {
                        this.updateTheStatus(`ids matched`);

                        let newArray = this.remoteStreams.filter(
                            (str) => str !== stream
                        );

                        this.remoteStreams = newArray;
                    }
                });

                //print remote stream array
                if (this.remoteStreams.length > 0) {
                    for (var i = 0; i < this.remoteStreams.length; i++) {
                        this.updateTheStatus(
                            `RemoteStreamsList-1`,
                            this.remoteStreams[i]
                        );
                    }
                }

                if (
                    this.parentStreamId &&
                    this.parentStreamId === theEventStream.id
                ) {
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
                    } catch {}
                }
                if (
                    this.role === this.Roles.BROADCAST &&
                    this.raiseHands.includes(target)
                ) {
                    var index = this.raiseHands.indexOf(target);
                    if (index > -1) {
                        this.raiseHands.splice(index, 1);
                    }
                }

                //check meeting status and close socket
                if (this.leftMeeting) {
                    //close websocket
                    if (this.socket) {
                        this.socket.onclose = () => {
                            this.updateTheStatus(
                                `socket is closed after leaveMeeting`
                            );
                            this.resetVariables(true);
                        }; //empty on close callback
                        this.socket.close();
                        this.socket = null;

                        return;
                    }
                }
            }; //end of on Inactive

            stream.name = ''; // currently we don't know name so it's empty

            this.updateTheStatus(`ReceivedStream:`, stream);

            //print remote stream array
            if (this.remoteStreams.length > 0) {
                for (var i = 0; i < this.remoteStreams.length; i++) {
                    this.updateTheStatus(
                        `RemoteStreamsList-2:`,
                        this.remoteStreams[i]
                    );
                }
            }

            //remote Stream callback
            if (this.remoteStreamCallback) this.remoteStreamCallback(stream);

            this.enqueue(this.remoteStreamsQueue, stream); //save stream into queue for Getting it's user name

            this.registerUserListCallback(); //calback to get user list with streams, to identify username of stream

            this.remoteStreams.push(stream);
            stream.getTracks().forEach((t) => {
                this.trackToStreamMap[t.id] = stream.id;
            });
            if (!this.remoteStreamNotified) {
                this.remoteStreamNotified = true;
                this.updateTheStatus(`[newPeerConnectionInstance] A7`);

                if (await this.checkSocketStatus())
                    this.socket.send(
                        JSON.stringify({
                            type: 'stream',
                            data: 'true',
                        })
                    );
                this.updateTheStatus(
                    `[newPeerConnectionInstance] stream message`
                );
            }
            this.targetStreams[target] = stream.id;

            for (const userId in this.myPeerConnectionArray) {
                const apeerConnection = this.myPeerConnectionArray[userId];
                this.updateTheStatus(
                    `check Sending the stream [${
                        stream.id
                    }] tracks to ${userId} ${apeerConnection.isAdience.toString()}`
                );
                if (!apeerConnection.isAdience) continue;

                this.updateTheStatus(
                    `Sending the stream [${stream.id}] tracks to ${userId}`
                );
                stream.getTracks().forEach((track) => {
                    try {
                        track.streamId = stream.id;
                        apeerConnection.addTrack(track, stream);
                    } catch {}
                });
            }

            if (!this.started) {
                this.started = true;
            }
        };

        let connectedOnce = false;
        peerConnection.oniceconnectionstatechange = (event) => {
            this.updateTheStatus(
                `[newPeerConnectionInstance] oniceconnectionstatechange peerConnection.iceConnectionState = ${
                    peerConnection.iceConnectionState
                } event = ${JSON.stringify(event)}`
            );
            switch (peerConnection.iceConnectionState) {
                case 'connected': {
                    connectedOnce = true;
                    peerConnection._iceIsConnected = true;
                    break;
                }
                default:
                    peerConnection._iceIsConnected = false;
                    break;
            }
            if (
                peerConnection.iceConnectionState === 'disconnected' ||
                peerConnection.iceConnectionState === 'failed' ||
                peerConnection.iceConnectionState === 'closed'
            ) {
                if (
                    this.role === this.Roles.BROADCAST &&
                    this.raiseHands.includes(target)
                ) {
                    var index = this.raiseHands.indexOf(target);
                    if (index > -1) {
                        this.raiseHands.splice(index, 1);
                    }
                }
                if (!this.parentDC && !connectedOnce) {
                    setTimeout(() => {
                        this.updateTheStatus('restarting ice');
                        peerConnection.restartIce();
                    }, 0);
                } else {
                    this.updateTheStatus(
                        'closing the peer connection: ' + target
                    );
                    peerConnection.close();
                    delete this.myPeerConnectionArray[target];
                }
                this.restartEverything(peerConnection, target, isAudience);
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
     * Get list of user and match their name with respective stream
     *
     */
    registerUserListCallback() {
        //here's logic to get name of stream owener

        //set userlist callback to receive list of all the users in the meeting with thier streams

        this.userListCallback = (users) => {
            this.updateTheStatus(`userList:`, users);

            const stream = this.dequeue(this.remoteStreamsQueue);

            if (stream) {
                let broadcasterName = '';

                // //check if user list contain broadcaster and save it's name
                let hasBroadcaster = false;
                users.forEach((user) => {
                    if (user.role === this.Roles.BROADCASTER) {
                        hasBroadcaster = true;
                        const data = JSON.parse(user.name);
                        broadcasterName = data.name;
                    }
                });

                //if not then, enqueue curret stream again to the Queue
                if (hasBroadcaster === false) {
                    this.updateTheStatus(`hasBroadcaster`, hasBroadcaster);
                    this.enqueue(this.remoteStreamsQueue, stream);
                    this.getLatestUserList();
                    return;
                }

                if (hasBroadcaster) {
                    //iterate over each user
                    users.forEach((user) => {
                        let role = this.Roles.AUDIENCE;
                        let userName = '';

                        if (user) {
                            //if video is undefined, it means user list not updated yet, fetch again
                            //It must be null or have MediaStream
                            if (user.video === undefined) {
                                this.getLatestUserList();
                                this.updateTheStatus(
                                    `going to fetch latestuserlist`
                                );
                                return;
                            }

                            //save broadcaster role
                            if (user.role === this.Roles.BROADCASTER) {
                                role = this.Roles.BROADCAST;
                            }

                            //if video not null nor undefined, get it's name
                            if (
                                user.video !== null &&
                                user.video !== undefined
                            ) {
                                if (user.video.id === stream.id) {
                                    const data = JSON.parse(user.name);
                                    userName = data.name;
                                }
                            }

                            //user name is not null nor empty
                            if (userName != null && userName != '') {
                                if (this.remoteStreamCallback) {
                                    this.updateTheStatus(`userName:`, userName);
                                    stream.name = userName;
                                    stream.role = role;
                                    stream.userId = user.id;

                                    //stream from Broadcatser
                                    if (
                                        role === this.Roles.BROADCAST &&
                                        userName === broadcasterName
                                    ) {
                                        this.updateTheStatus(
                                            `stream from Broadcast`
                                        );
                                        this.remoteStreamCallback(stream);
                                    }

                                    //stream from Audience
                                    if (
                                        role === this.Roles.AUDIENCE &&
                                        userName != broadcasterName
                                    ) {
                                        this.updateTheStatus(
                                            `stream from Audience`
                                        );
                                        this.remoteStreamCallback(stream);
                                    }

                                    return;
                                }
                            }
                        }
                    });

                    //screen share video name
                    if (this.broadcastersMessage != null) {
                        let message = JSON.parse(this.broadcastersMessage);
                        this.updateTheStatus('message: ', message);

                        if (
                            message.id === stream.id &&
                            message.type === this.StreamType.SCREEN
                        ) {
                            if (this.remoteStreamCallback) {
                                stream.name = broadcasterName;
                                stream.role = this.Roles.BROADCAST;
                                stream.isShareScreen = true;
                                this.remoteStreamCallback(stream);
                            }
                        }
                    }
                }
            }
        };

        this.getLatestUserList();
    }

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
        this.updateTheStatus(
            `[createOrGetPeerConnection] audienceName = ${audienceName}, isAudience = ${isAudience}`
        );
        if (this.myPeerConnectionArray[audienceName])
            return this.myPeerConnectionArray[audienceName];

        this.myPeerConnectionArray[audienceName] =
            this.newPeerConnectionInstance(audienceName, true, isAudience);
        this.updateTheStatus(
            `[createOrGetPeerConnection] generate newPeerConnectionInstance`
        );

        return this.myPeerConnectionArray[audienceName];
    };

    /**
     * Function to add new Audience as Current Node's Children
     * @param {String} audienceName
     */
    connectToAudience = (audienceName) => {
        this.updateTheStatus(`Connecting to ${audienceName}`);
        this.updateTheStatus(
            `[handleMessage] connectToAudience ${audienceName}`
        );
        if (!this.localStream && this.remoteStreams.length === 0) {
            return;
        }
        if (!this.myPeerConnectionArray[audienceName]) {
            this.updateTheStatus(`Creating peer connection to ${audienceName}`);
            this.myPeerConnectionArray[audienceName] =
                this.newPeerConnectionInstance(
                    audienceName,
                    this.localStream || this.remoteStreams,
                    true
                );
        }
        this.updateTheStatus(
            `[handleMessage] generate newPeerConnectionInstance`
        );

        if (this.remoteStreams.length > 0) {
            this.updateTheStatus(`publishing stream/s to ${audienceName}`);
            this.remoteStreams.forEach((astream) => {
                this.updateTheStatus(`streamToPublish:`, astream);
                astream.getTracks().forEach((track) => {
                    try {
                        this.myPeerConnectionArray[audienceName].addTrack(
                            track,
                            astream
                        );
                    } catch {}
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
        this.updateTheStatus(`[handleMessage] sendStreamTo ${target}`);

        const peerConnection = this.createOrGetPeerConnection(target, false);
        stream.getTracks().forEach((track) => {
            if (this.lastVideoState === 'Disabled') {
                this.disableVideo();
            }
            if (this.lastAudioState === 'Disabled') {
                this.disableAudio();
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
            this.myPeerConnectionConfig.iceServers = iceServers.filter(
                (i) => i.url.indexOf('turn') < 0
            );
        }
        this.updateTheStatus(`[start] ${this.role}`);
        this.updateTheStatus(`Getting media capabilities`);
        await this.getSupportedConstraints();
        if (this.role === this.Roles.BROADCAST) {
            this.updateTheStatus(`Start broadcasting`);
            return this.startBroadcasting();
        } else if (!this.constraints.audio && !this.constraints.video) {
            this.updateTheStatus(`No media removing raise hand`);
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
        if (this.localStream) {
            this.lastVideoState = enabled === true ? 'Enabled' : 'Disabled';
            this.localStream.getTracks().forEach((track) => {
                if (track.kind === 'video') track.enabled = enabled;
            });
        }
    };

    /**
     * Function to enable / disable Audio track
     *
     * @param {boolean} enabled
     */
    disableAudio = (enabled = false) => {
        if (this.localStream) {
            this.lastAudioState = enabled === true ? 'Enabled' : 'Disabled';
            this.localStream.getTracks().forEach((track) => {
                if (track.kind === 'audio') track.enabled = enabled;
            });
        }
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
        return new Promise(async (resolve, reject) => {
            if (await this.checkSocketStatus())
                this.socket.send(
                    JSON.stringify({
                        type: 'broadcaster-status',
                    })
                );

            let i = 0;
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
    updateTheStatus = (tag, data = null) => {
        if (this.updateStatus) {
            try {
                this.updateStatus(tag, data);
            } catch (e) {
                console.log('Failed to Update the Status: ', e);
            }
        }
    };

    /**
     * Function to lower hand and take request(to broadcast) back, if sharing already stop sharing
     *
     * @returns
     */
    lowerHand = async () => {
        this.updateTheStatus('[lowerHand] start');
        if (!this.localStream) return;
        let apeerConnection;
        for (const id in this.myPeerConnectionArray) {
            apeerConnection = this.myPeerConnectionArray[id];
            break;
        }
        const trackIds = this.localStream
            .getTracks()
            .map((receiver) => receiver.id);
        this.updateTheStatus('[lowerHand] trackIds', trackIds);
        const allSenders = apeerConnection.getSenders();
        this.updateTheStatus(`[lowerHand] allSenders`, allSenders);
        for (const trackId of trackIds)
            for (const sender of allSenders) {
                this.updateTheStatus(`[lowerHand] sender`, sender);
                if (!sender.track) continue;
                if (sender.track.id === trackId) {
                    this.updateTheStatus(`[lowerHand] DC sender`);
                    try {
                        apeerConnection.removeTrack(sender);
                    } catch (e) {
                        this.updateTheStatus(e);
                    }
                }
            }
        this.localStream.getTracks().forEach((track) => {
            track.stop();
        });
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
    setMetadata = async (metadata) => {
        if (await this.checkSocketStatus())
            this.socket.send(
                JSON.stringify({
                    type: 'metadata-set',
                    data: JSON.stringify(metadata),
                })
            );
    };
    getMetadata = async () => {
        if (await this.checkSocketStatus())
            this.socket.send(
                JSON.stringify({
                    type: 'metadata-get',
                })
            );
    };
    streamById = (streamId) => {
        return this.remoteStreams.find((s) => s.id === streamId);
    };
    stopSignaling = () => {
        if (this.pingInterval) {
            clearInterval(this.pingInterval);
            this.pingInterval = null;
        }
        if (this.socket) {
            this.socket.close();
        }
    };

    /**
     * To restart again and connect to new parent
     * we don't need to close the socket this time
     */
    restart = async (closeSocket) => {
        this.updateTheStatus(`restarting.. close socket: `, closeSocket);
        //check for local stream and stop tracks
        if (this.localStream) {
            this.localStream.getTracks().forEach(function (track) {
                track.stop();
            });
            this.localStream = null;
        }

        //close all the peer connections
        let idList = [];
        if (
            this.myPeerConnectionArray &&
            this.myPeerConnectionArray.length > 0
        ) {
            for (const u in this.myPeerConnectionArray) {
                this.myPeerConnectionArray[u].close();
                idList.push(u);
            }
        }

        idList.forEach((id) => delete sparkRTC.value.myPeerConnectionArray[id]);

        //reset few variables
        this.resetVariables(false);

        //close the web socket
        if (closeSocket && this.socket) {
            this.socket.close();
            this.socket.onclose = async () => {
                this.updateTheStatus(`socket is closed in restart`);
                this.socket = null;

                //waiting to websocket to close then repoen again
                if (this.startAgain) {
                    this.startAgain();
                }
            }; //on close callback
        } else {
            this.updateTheStatus(`socket closing is not required`);
        }

        //update the UI (Controllers)
        if (this.updateUi) this.updateUi();
    };

    /**
     * To leave the meeting
     */
    leaveMeeting = async () => {
        //check for local stream and stop tracks
        //stop all the sender tracks
        if (this.localStream) {
            this.leftMeeting = true;

            await this.lowerHand();
        } else {
            //close websocket if not streaming anything
            if (this.socket) {
                this.socket.onclose = () => {
                    this.updateTheStatus(`socket is closed after leaveMeeting`);
                    this.resetVariables();
                }; //empty on close callback
                this.socket.close();
                this.socket = null;
            }
        }

        if (this.role === this.Roles.BROADCAST) {
            //stop screen share
            if (this.shareStream) {
                await this.stopShareScreen(this.shareStream);
                this.shareStream = null;
            }

            //close websocket
            if (this.socket) {
                this.socket.onclose = () => {
                    this.updateTheStatus(`socket is closed after leaveMeeting`);
                    this.resetVariables();
                }; //empty on close callback
                this.socket.close();
                this.socket = null;
            }
        }

        this.updateTheStatus(`left meeting`);
    };

    //Reset all the variables
    resetVariables = (resetAll = false) => {
        if (resetAll) {
            this.socketURL = '';
            this.socket = null;
        }
        this.myName = 'NoName';
        this.roomName = 'SparkRTC';
        this.myUsername = 'NoUsername';
        this.role = '';
        this.started = false;
        this.myPeerConnectionConfig = {
            iceServers,
        };
        this.localStream = null;
        this.remoteStreamNotified = false;
        this.remoteStreams = [];
        this.lastBroadcasterId = '';
        this.broadcastingApproved = false;
        this.myPeerConnectionArray = {};
        this.iceCandidates = [];
        this.pingInterval = null;
        this.raiseHands = [];
        this.startedRaiseHand = false;
        this.targetStreams = {};
        this.parentStreamId = null;
        this.broadcasterStatus = '';
        this.parentDC = true;
        this.broadcasterDC = true;
        this.userListCallback = null;
        this.remoteStreamsQueue = [];
        this.metaData = {};
        this.userStreamData = {};
        this.users = [];
        this.trackToStreamMap = {};
        this.lastVideoState = 'Enabled';
        this.lastAudioState = 'Enabled';
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
        this.treeCallback = options.treeCallback;
        this.raiseHandConfirmation = options.raiseHandConfirmation;
        this.altBroadcastApprove = options.altBroadcastApprove;
        this.newTrackCallback = options.newTrackCallback;
        this.startProcedure = options.startProcedure;
        this.constraintResults = options.constraintResults;
        this.updateStatus = options.updateStatus;
        this.userListUpdated = options.userListUpdated;
        this.maxLimitReached = options.maxLimitReached;
        this.disableBroadcasting = options.disableBroadcasting;
        this.connectionStatus = options.connectionStatus;
        this.userInitialized = options.userInitialized;
        this.startAgain = options.startAgain;
        this.updateUi = options.updateUi;
        this.parentDcMessage = options.parentDcMessage;
    }
}
