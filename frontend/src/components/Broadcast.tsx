import React, {useEffect} from 'react';
import {PEER_CONNECTION_CONFIG} from "../config/myPeerConnectionConfig";
import Main from "./Main";
import {useStreamMap} from "../hooks/useStreamMap";

export default function Broadcast({myName}: { myName: string }) {
    let socket: WebSocket;
    let myPeerConnection: RTCPeerConnection | undefined;
    let myUsername = 'NoUsername';
    let myPeerConnectionArray: any = {};
    let iceCandidates: any[] = [];

    let {streamMap, setStreamMap, enableLocalStream} = useStreamMap();

    useEffect(() => {
        setupSignalingSocket().then(e =>
            startBroadcasting().then()
        )
    })

    const handleMessage = async (event: any) => {
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
            case 'video-answer':
                console.log('Got answer.', msg);
                audiencePeerConnection = createOrGetPeerConnection(msg.data);
                try {
                    await audiencePeerConnection.setRemoteDescription(new RTCSessionDescription(msg.sdp));
                } catch (e: any) {
                    console.log('setRemoteDescription failed with exception: ' + e.message);
                    console.log(audiencePeerConnection);
                    console.log(msg.sdp);
                }
                break;
            case 'new-ice-candidate':
                console.log('Got ICE candidate.', msg);
                audiencePeerConnection = createOrGetPeerConnection(msg.data);
                iceCandidates.push(new RTCIceCandidate(msg.candidate));
                if (audiencePeerConnection && audiencePeerConnection.remoteDescription) {
                    audiencePeerConnection.addIceCandidate(iceCandidates.pop());
                }
                break;
            case 'role':
                if (msg.data === "no:broadcast") {
                    alert("You are not a broadcaster anymore!");
                    socket.close();
                } else if (msg.data === "yes:broadcast") {
                    // localVideoTag.srcObject = localStream;
                    // document.getElementById('signalConnectBtn').disabled = true;
                    // document.getElementById('myName').disabled = true;
                    // document.getElementById('myName').value = myName;
                } else {
                    // localVideoTag.srcObject = null;
                }
                break;
            case 'start':
                if (msg.error) {
                    alert(msg.error);
                    return;
                }

                myUsername = msg.data;
                break;
            case 'add_audience':
            case 'add_broadcast_audience':
                connectToAudience(msg.data);
                break;
            default:
                break;
        }
    };

    const setupSignalingSocket = async () => {
        const baseUrl = window.location.href.split("//")[1].split("/")[0];
        let protocol = "wss";
        if (window.location.href.split("//")[0] === "http:") {
            protocol = "ws";
        }
        socket = new WebSocket(`${protocol}://${baseUrl}/ws`);
        socket.onmessage = handleMessage;
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
    }

    const startBroadcasting = async () => {
        if (!streamMap.get('localStream')) {
            console.log('# 3 #');
            console.log('==(1)==> The caller captures local Media via MediaDevices.getUserMedia')
            enableLocalStream().then(_ => {
                    socket.send(
                        JSON.stringify({
                                type: "role",
                                data: "broadcast"
                                // data: myRole==='broadcast' ? 'alt-broadcast' : myRole
                            }
                        )
                    );
                }
            );
        }
    }
    const newPeerConnectionInstance = (target: any, addLocalStream = true) => {
        const peerConnection = new RTCPeerConnection(PEER_CONNECTION_CONFIG);

        peerConnection.onicecandidate = (event) => {
            if (event.candidate) {
                socket.send(
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
                socket.send(
                    JSON.stringify({
                        type: "video-offer",
                        sdp: peerConnection.localDescription,
                        target,
                        name: myUsername,
                    })
                );
                console.log('onnegotiationneeded done');
            } catch (e) {
                console.log('onnegotiationneeded failed:', e);
            }
        };

        // if (addLocalStream) peerConnection.addStream(localStream);
        if (addLocalStream) {
            let localStream = streamMap.get('localStream');
            if (localStream) {
                console.log('Adding tracks... len=', localStream.getTracks().length);
                localStream.getTracks().forEach((track) => {
                    if (myPeerConnection && localStream) {
                        myPeerConnection.addTrack(track, localStream);
                        console.log('==(2)==> The caller creates RTCPeerConnection and calls RTCPeerConnection.addTrack()')
                    } else {
                        console.error('unable to add tracks')
                    }
                });
            }
        }
        return peerConnection;
    };
    const createOrGetPeerConnection = (audienceName: any) => {
        if (myPeerConnectionArray[audienceName]) return myPeerConnectionArray[audienceName];

        myPeerConnectionArray[audienceName] = newPeerConnectionInstance(audienceName, false);

        return myPeerConnectionArray[audienceName];
    }
    const connectToAudience = (audienceName: any) => {
        console.log('connecting to', audienceName);
        if (!streamMap.get('localStream')) return;
        if (myPeerConnectionArray[audienceName]) return;

        myPeerConnectionArray[audienceName] = newPeerConnectionInstance(audienceName);
    };
    return (
        <Main myName={myName} myRole={"broadcast"}/>
    )
}