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
    socket;
    myName = 'NoName';
    myUsername = 'NoUsername';
    myPeerConnectionArray = {};
    iceCandidates = [];
    handleVideoOfferMsg = async (msg) => {
        const broadcasterPeerConnection = this.myPeerConnectionArray[msg.name] || this.newPeerConnectionInstance(msg.name);
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
    }
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
                this.handleVideoOfferMsg(msg);
                break;
            case 'video-answer':
                console.log('Got answer.', msg);
                audiencePeerConnection = this.createOrGetPeerConnection(msg.data);
                try {
                    await audiencePeerConnection.setRemoteDescription(new RTCSessionDescription(msg.sdp));
                } catch (e) {
                    console.log('setRemoteDescription failed with exception: ' + e.message);
                    console.log(audiencePeerConnection);
                    console.log(msg.sdp);
                }
                break;
            case 'new-ice-candidate':
                console.log('Got ICE candidate.', msg);
                audiencePeerConnection = this.createOrGetPeerConnection(msg.data);
                this.iceCandidates.push(new RTCIceCandidate(msg.candidate));
                if (audiencePeerConnection && audiencePeerConnection.remoteDescription) {
                    audiencePeerConnection.addIceCandidate(this.iceCandidates.pop());
                }
                break;
            case 'role':
                if (this.role === 'broadcast') {
                    if (msg.data === "no:broadcast") {
                        alert("You are not a broadcaster anymore!");
                        this.socket.close();
                    } else if (msg.data === "yes:broadcast") {
                        this.localStreamChangeCallback(this.localStream);
                    } else {
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
            myPeerConnection = undefined;
            myPeerConnectionArray = [];
            setupSignalingSocket();
        };
        this.socket = socket;
    }
    startBroadcasting = async () => {
        try {
            this.localStream = await navigator.mediaDevices.getUserMedia({
                audio: true,
                video: true,
            });
            this.socket.send(
                JSON.stringify({
                    type: "role",
                    data: "broadcast",
                })
            );
        } catch (e) {
            console.log(e);
            alert('Unable to get access to your webcam and microphone.');
        }
    }
    startReadingBroadcast = async () => {
        this.socket.send(
            JSON.stringify({
                type: "role",
                data: "audience",
            })
        );
    }
    newPeerConnectionInstance = (target, addLocalStream = true) => {
        const peerConnection = new RTCPeerConnection(this.myPeerConnectionConfig);

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
            console.log("onTrack");
            this.remoteStreamCallback(event.streams[0]);
            this.socket.send(
                JSON.stringify({
                    type: "stream",
                    data: "true",
                })
            );
        };

        if (addLocalStream && this.localStream)
            this.localStream.getTracks().forEach((track) => {
                peerConnection.addTrack(track, this.localStream);
            });

        return peerConnection;
    };
    createOrGetPeerConnection = (audienceName) => {
        if (this.myPeerConnectionArray[audienceName]) return this.myPeerConnectionArray[audienceName];

        this.myPeerConnectionArray[audienceName] = this.newPeerConnectionInstance(audienceName, false);

        return this.myPeerConnectionArray[audienceName];
    }
    connectToAudience = (audienceName) => {
        console.log('connecting to', audienceName);
        if (!this.localStream) return;
        if (this.myPeerConnectionArray[audienceName]) return;

        this.myPeerConnectionArray[audienceName] = this.newPeerConnectionInstance(audienceName);
    };
    start = () => {
        if (this.role === 'broadcast') {
            return this.startBroadcasting();
        }

        return this.startReadingBroadcast();
    }

    constructor(role, localStreamChangeCallback, remoteStreamCallback) {
        this.role = role;
        this.localStreamChangeCallback = localStreamChangeCallback;
        this.remoteStreamCallback = remoteStreamCallback;
    }
}