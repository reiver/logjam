class SparkRTC {
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

    role = 'broadcast';
    localStream;
    remoteStreamNotified = false;
    remoteStreams = [];
    socket;
    myName = 'NoName';
    myUsername = 'NoUsername';
    myPeerConnectionArray = {};
    iceCandidates = [];
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
        let msg;
        try {
            msg = JSON.parse(event.data);
        } catch (e) {
            return;
        }
        msg.data = (msg.Data && !msg.data) ? msg.Data : msg.data;
        msg.type = (msg.Type && !msg.type) ? msg.Type : msg.type;

        if (msg.type !== 'new-ice-candidate') console.log(msg);
        let audiencePeerConnection;
        switch (msg.type) {
            case 'video-offer':
            case 'alt-video-offer':
                this.handleVideoOfferMsg(msg);
                break;
            case 'video-answer':
            case 'alt-video-answer':
                // console.log('Got answer.', msg);
                audiencePeerConnection = this.createOrGetPeerConnection(msg.data, true);
                try {
                    await audiencePeerConnection.setRemoteDescription(new RTCSessionDescription(msg.sdp));
                } catch (e) {
                    console.log('setRemoteDescription failed with exception: ' + e.message);
                    console.log(audiencePeerConnection);
                    console.log(msg.sdp);
                }
                break;
            case 'new-ice-candidate':
            case 'alt-new-ice-candidate':
                // console.log('Got ICE candidate.', msg);
                audiencePeerConnection = this.createOrGetPeerConnection(msg.data);
                this.iceCandidates.push(new RTCIceCandidate(msg.candidate));
                if (audiencePeerConnection && audiencePeerConnection.remoteDescription) {
                    console.log('Adding Candidate');
                    audiencePeerConnection.addIceCandidate(this.iceCandidates.pop());
                }
                break;
            case 'role':
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
                console.log('Start connecting to broadcaster');
                this.sendStreamTo(msg.data, this.localStream);
                break;
            case 'alt-broadcast':
                if (this.role === 'broadcast' && confirm(`${msg.name} wants to broadcast, do you approve?`)) {
                    this.socket.send(
                        JSON.stringify({
                            type: "alt-broadcast-approve",
                            target: msg.data,
                        })
                    );
                }
                break;
            default:
                break;
        }
    };
    setupSignalingSocket = (url, myName) => {
        if (myName)
            this.myName = myName;
        const socket = new WebSocket(url);
        socket.onmessage = this.handleMessage;
        socket.onopen = () => {
            console.log("WebSocket connection opened");
            socket.send(
                JSON.stringify({
                    type: "start",
                    data: myName,
                })
            );
        };
        socket.onclose = () => {
            console.log("WebSocket connection closed");
            this.remoteStreamNotified = false;
            this.myPeerConnectionArray = {};
            if (this.signalingDisconnectedCallback) this.signalingDisconnectedCallback;
            this.setupSignalingSocket(url, myName);
        };
        this.socket = socket;
    };
    startShareScreen = async () => {
        try {
            const shareStream = await navigator.mediaDevices
                .getDisplayMedia({
                    audio: true,
                    video: true,
                });
            this.remoteStreams.push(shareStream);
            console.log('Sending to all peers');
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
                data: "audience",
            })
        );
    };
    raiseHand = () => {
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
                console.log('onnegotiationneeded done');
            } catch (e) {
                console.log(e);
                alert('onnegotiationneeded failed:', e);
            }
        };

        peerConnection.ontrack = (event) => {
            const stream = event.streams[0];
            console.log({stream});
            if (this.remoteStreams.indexOf(event.streams[0]) !== -1) return;
            if (this.remoteStreamCallback)
                this.remoteStreamCallback(stream);
            if (this.remoteStreams.indexOf(stream) === -1)
                this.remoteStreams.push(stream);
            console.log("onTrack");
            if (!this.remoteStreamNotified) {
                this.remoteStreamNotified = true;
                this.socket.send(
                    JSON.stringify({
                        type: "stream",
                        data: "true",
                    })
                );
            }

            console.log('Sending to all peers');
            for (const userId in this.myPeerConnectionArray) {
                if (userId === target) continue;
                const apeerConnection = this.myPeerConnectionArray[userId];
                if (!apeerConnection.isAdience) return;

                stream.getTracks().forEach((track) => {
                    apeerConnection.addTrack(track, stream);
                });
            }
        };

        peerConnection.oniceconnectionstatechange = (event) => {
            if (peerConnection.iceConnectionState == 'disconnected') {
                console.log('Disconnected', peerConnection);
                this.socket.close();
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
        console.log('connecting to', audienceName, !!this.localStream, this.remoteStreams.length);
        if (!this.localStream && this.remoteStreams.length === 0) return;
        if (!this.myPeerConnectionArray[audienceName]) {
            this.myPeerConnectionArray[audienceName] = this.newPeerConnectionInstance(audienceName, this.localStream || this.remoteStreams, true);
        }

        if (this.remoteStreams.length > 0) {
            this.remoteStreams.forEach((astream) => {
                console.log('Adding remote stream to peer connection', astream.id);
                astream.getTracks().forEach((track) => {
                    this.myPeerConnectionArray[audienceName].addTrack(track, astream);
                });
            });
            console.log('remoteStreams', this.remoteStreams);
        }
    };
    sendStreamTo = (target, stream) => {
        const peerConnection = this.createOrGetPeerConnection(target, false);
        stream.getTracks().forEach((track) => {
            peerConnection.addTrack(track, stream);
        });
    };
    start = () => {
        if (this.role === 'broadcast') {
            return this.startBroadcasting();
        }

        return this.startReadingBroadcast();
    };

    constructor(role, options = {}) {
        this.role = role;
        this.localStreamChangeCallback = options.localStreamChangeCallback;
        this.remoteStreamCallback = options.remoteStreamCallback;
        this.signalingDisconnectedCallback = options.signalingDisconnectedCallback;
    }
}