import React, {useEffect} from 'react';
import {PEER_CONNECTION_CONFIG} from "../config/myPeerConnectionConfig";
import {useStreamMap} from "../hooks/useStreamMap";
import Main from "./Main";
import {useSocket} from "../hooks/useSocket";

export default function Audience({myName}: { myName: string }) {
    let socket: WebSocket = useSocket();
    let myPeerConnection: RTCPeerConnection | undefined;
    let myUsername = 'NoUsername';
    let myPeerConnectionArray: any = {};
    let iceCandidates: any[] = [];

    let {streamMap, setStreamMap} = useStreamMap();

    useEffect(() => {
        console.log('setupSignalingSocket');
        setupSignalingSocket().then(e =>
            startReadingBroadcast().then()
        )
    })

    const handleVideoOfferMsg = async (msg: any) => {
        const broadcasterPeerConnection = myPeerConnectionArray[msg.name] || newPeerConnectionInstance(msg.name);
        await broadcasterPeerConnection.setRemoteDescription(new RTCSessionDescription(msg.sdp));
        await broadcasterPeerConnection.setLocalDescription(await broadcasterPeerConnection.createAnswer());

        socket.send(
            JSON.stringify({
                name: myUsername,
                target: msg.name,
                type: "video-answer",
                sdp: broadcasterPeerConnection.localDescription,
            })
        );
    }
    const handleMessage = (event: any) => {
        let msg;
        try {
            msg = JSON.parse(event.data);
        } catch (e) {
            return;
        }
        msg.data = (msg.Data && !msg.data) ? msg.Data : msg.data;
        msg.type = (msg.Type && !msg.type) ? msg.Type : msg.type;

        if (msg.type !== 'new-ice-candidate') console.log(msg);
        switch (msg.type) {
            case 'video-offer':
                handleVideoOfferMsg(msg);
                break;
            case 'new-ice-candidate':
                iceCandidates.push(new RTCIceCandidate(msg.candidate));
                if (myPeerConnection && myPeerConnection.remoteDescription) {
                    myPeerConnection.addIceCandidate(iceCandidates.pop());
                }
                break;
            case 'role':
                // document.getElementById('signalConnectBtn').disabled = true;
                // document.getElementById('myName').disabled = true;
                // document.getElementById('myName').value = myName;
                break;
            case 'start':
                if (msg.error) {
                    alert(msg.error);
                    return;
                }

                myUsername = msg.data;
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
    const startReadingBroadcast = async () => {
        socket.send(
            JSON.stringify({
                type: "role",
                data: "audience",
            })
        );
    }
    const newPeerConnectionInstance = (target: any) => {
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

        peerConnection.ontrack = (event) => {
            console.log("onTrack");
            // if (remoteVideoTag.srcObject === event.streams[0]) return;
            // remoteVideoTag.srcObject = event.streams[0];
            if (streamMap.get('remoteStream')) return;
            setStreamMap((
                prev: Map<string, MediaStream>) => new Map(prev).set('remoteStream', event.streams[0]));

            socket.send(
                JSON.stringify({
                    type: "stream",
                    data: "true",
                })
            );
        };

        peerConnection.oniceconnectionstatechange = function (event) {
            if (peerConnection.iceConnectionState === 'disconnected') {
                console.log('Disconnected', peerConnection);
                socket.close();
            }
        };

        return peerConnection;
    };
    const connectToAudience = (audienceName: string) => {
        console.log('connecting to', audienceName);
        if (!streamMap.get('localStream')) return;
        if (myPeerConnectionArray[audienceName]) return;

        myPeerConnectionArray[audienceName] = newPeerConnectionInstance(audienceName);
    };
    return (
        <div>
            <h1>Audience</h1>
            <Main myName={myName} myRole={"audience"}/>

        </div>
    )
}