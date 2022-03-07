import React, {useEffect, useState} from 'react';
import {PEER_CONNECTION_CONFIG} from "../config/myPeerConnectionConfig";
import Main from "./Main";
import {useStreamMap} from "../hooks/useStreamMap";
import {useSocket} from "../hooks/useSocket";
import {usePeerConnectionMap} from "../hooks/usePeerConnectionMap";


export default function Broadcast({myName}: { myName: string }) {
    let socket: WebSocket = useSocket();
    let myPeerConnection: RTCPeerConnection | undefined;
    const [myUsername, setMyUsername] = useState();
    // let myPeerConnectionArray: any = {};
    let {peerConnectionMap, setPeerConnectionMap} = usePeerConnectionMap();

    let iceCandidates: any[] = [];

    let {streamMap, setStreamMap, enableLocalStream} = useStreamMap();

    useEffect(() => {
        if (socket.readyState != socket.OPEN) return;
        if (!myUsername) {
            socket.send(
                JSON.stringify({
                    type: "start",
                    data: myName,
                })
            );
            console.log('[Sent] start : myName = ', myName);
        } else {
            if (streamMap.get('localStream')) return;
            enableLocalStream().then(_ =>
                socket.send(
                    JSON.stringify({
                            type: "role",
                            data: "broadcast"
                            // data: myRole==='broadcast' ? 'alt-broadcast' : myRole
                        }
                    )
                )
            )
        }
    }, [myName, myUsername, socket]);

    useEffect(() => {
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

            peerConnection.onnegotiationneeded = () => {
                try {
                    peerConnection.setLocalDescription().then(e => {
                            console.log('peerConnection.setLocalDescription() DONE')
                            peerConnection.createOffer().then(e => {
                                    console.log('peerConnection.createOffer() DONE')
                                    socket.send(
                                        JSON.stringify({
                                            type: "video-offer",
                                            sdp: peerConnection.localDescription,
                                            target,
                                            name: myUsername,
                                        })
                                    );
                                    console.log('video-offer sent')
                                }
                            )
                        }
                    )

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
            if (peerConnectionMap.get(audienceName)) return peerConnectionMap.get(audienceName);
            let newPeerConnection = newPeerConnectionInstance(audienceName, false);
            setPeerConnectionMap((
                prev: Map<string, RTCPeerConnection>) => new Map(prev).set(audienceName, newPeerConnection));
            console.log('[PeerConnection] new peerConnection added to peerConnectionMap');

            return newPeerConnection;
        }

        const connectToAudience = (audienceName: any) => {
            console.log('connecting to', audienceName);
            if (!streamMap.get('localStream')) {
                console.log('streamMap.get("localStream") = null')
                return;
            }
            if (peerConnectionMap.get(audienceName)) {
                console.log('peerConnection already exists')
                return;
            }
            console.log('newPeerConnectionInstance')
            let newPeerConnection = newPeerConnectionInstance(audienceName);
            setPeerConnectionMap((
                prev: Map<string, RTCPeerConnection>) => new Map(prev).set(audienceName, newPeerConnection));
            console.log('[PeerConnection] new peerConnection added to peerConnectionMap');
        };

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
                case 'start':
                    if (msg.error) {
                        alert(msg.error);
                        return;
                    }
                    setMyUsername(msg.data);
                    // myUsername = msg.data;
                    break;
                case 'video-answer':
                    console.log('Got answer.', msg);
                    audiencePeerConnection = createOrGetPeerConnection(msg.data);
                    if (audiencePeerConnection) {
                        try {
                            await audiencePeerConnection.setRemoteDescription(new RTCSessionDescription(msg.sdp));
                        } catch (e: any) {
                            console.log('setRemoteDescription failed with exception: ' + e.message);
                            console.log(audiencePeerConnection);
                            console.log(msg.sdp);
                        }
                    } else {
                        console.error('audiencePeerConnection is null')
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
                case 'add_audience':
                case 'add_broadcast_audience':
                    connectToAudience(msg.data);
                    break;
                default:
                    break;
            }
        };

        socket.addEventListener("message", handleMessage);
        // console.log('added socket onMessage listener');

        return () => {
            socket.removeEventListener("message", handleMessage);
            // console.log('socket onMessage listener removed');
        };
    }, [iceCandidates, myPeerConnection, myUsername, peerConnectionMap, setPeerConnectionMap, socket, streamMap]);

    // const setupSignalingSocket = () => {
    // const baseUrl = window.location.href.split("//")[1].split("/")[0];
    // let protocol = "wss";
    // if (window.location.href.split("//")[0] === "http:") {
    //     protocol = "ws";
    // }
    // socket = new WebSocket(`${protocol}://${baseUrl}/ws`);

    // socket = new WebSocket('ws://localhost:8080/ws');
    // socket.onmessage = handleMessage;
    // socket.onopen = () => {
    //     console.log("WebSocket connection opened");
    //     socket.send(
    //         JSON.stringify({
    //             type: "start",
    //             data: myName,
    //         })
    //     );
    // };
    // socket.onclose = () => {
    //     console.log("WebSocket connection closed");
    //     myPeerConnection = undefined;
    //     myPeerConnectionArray = [];
    //     setupSignalingSocket();
    // };
    // }

    if (!myUsername) {
        return null;
    }

    console.log('[Render] Broadcast')
    return (
        <div>
            <h1>Broadcast</h1>
            <Main myName={myName} myRole={"broadcast"}/>

        </div>
    )
}