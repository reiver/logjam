import {useCallback, useEffect, useState} from "react";
import {useSocket} from "../hooks/useSocket";
import Main from "./Main";
import {useParams} from "react-router-dom";
import {PEER_CONNECTION_CONFIG, turnStatus} from "../config/myPeerConnectionConfig";
import {useMessenger} from "../hooks/useMessenger";
import {usePeerConnectionMap} from "../hooks/usePeerConnectionMap";
import {useLocalStream} from "../hooks/useLocalStream";

export const Socket = ({myName}: { myName: string }) => {
    console.log('[Render] Socket');

    // const [myUsername, setMyUsername] = useState('');
    let myUsername: string;
    let {myRole} = useParams() || "audience";

    let {peerConnectionMap, setPeerConnectionMap} = usePeerConnectionMap();
    let {localStream, setLocalStream} = useLocalStream();

    const messenger = useMessenger();
    const socket = useSocket();
    // console.log(socket);

    useEffect(() => {
        // start
        messenger.send({
                type: "start",
                data: myName
            }
        );

        messenger.send({
                type: "role",
                data: myRole
            }
        );

        messenger.send({
                type: "turn_status",
                data: turnStatus()
            }
        );
    }, []);

    function connectUser(targetUsername: string) {
        let myPeerConnection = addPeerConnection(targetUsername);
        if (myPeerConnection && localStream) {
            console.log('Adding tracks... len=', localStream.getTracks().length);
            localStream.getTracks().forEach((track) => {
                myPeerConnection?.addTrack(track, localStream);
            });
        }
    }

    async function sendVideoAnswer(msg: any) {
        console.log('[Socket] sending video answer');
        console.log(msg);
        let targetUsername = msg.name ? msg.name : '';
        let targetName = msg.username;
        const desc = new RTCSessionDescription(msg.sdp);
        let myPeerConnection = addPeerConnection(targetUsername);
        console.log(myPeerConnection);
        if (myPeerConnection){
            await myPeerConnection.setRemoteDescription(desc);

            let answer = await myPeerConnection.createAnswer();

            await myPeerConnection.setLocalDescription(answer);

            messenger.send({
                name: myUsername,
                target: targetUsername,
                type: "video-answer",
                sdp: myPeerConnection.localDescription,
            });
        }
    }

    function addPeerConnection(targetUsername: string) {
        if (peerConnectionMap.get(targetUsername)) {
            console.log('peerConnection already exists.');
            return peerConnectionMap.get(targetUsername);
        }
        let newPeerConnection = createPeerConnection(targetUsername);
        setPeerConnectionMap((
            prev: Map<string, RTCPeerConnection>) => new Map(prev).set(targetUsername, newPeerConnection));
        return newPeerConnection;
    }

    function createPeerConnection(targetUsername: string) {
        let peerConnection = new RTCPeerConnection(PEER_CONNECTION_CONFIG);

        peerConnection.onicecandidate = function (event) {
            if (event.candidate) {
                messenger.send({
                    type: "new-ice-candidate",
                    target: 'targetUsername',
                    candidate: event.candidate,
                });
            }
        };

        peerConnection.onnegotiationneeded = function (event) {
            console.log("onnegotiationneeded", event);
            peerConnection
                .createOffer()
                .then(function (offer) {
                    return peerConnection.setLocalDescription(offer);
                })
                .then(function () {
                    messenger.send({
                        name: myUsername,
                        target: targetUsername,
                        type: "video-offer",
                        sdp: peerConnection.localDescription,
                    });

                })
                .catch((err) => {
                    messenger.send({
                        type: "log",
                        data: "onnegotiationneeded audience error:" + JSON.stringify(err),
                    });
                });
        };
        console.log('peerConnection created: ', peerConnection);

        return peerConnection;
    }

    useEffect(() => {
        if (myRole === 'audience') {

        }
    }, [myRole])

    function getStateDescription(readyState: number) {
        switch (readyState) {
            case 0:
                return "CONNECTING	Socket has been created. The connection is not yet open.";
            case 1:
                return "OPEN	The connection is open and ready to communicate."
            case 2:
                return "CLOSING	The connection is in the process of closing."
            case 3:
                return "CLOSED	The connection is closed or couldn't be opened."
            default:
                return "Unknown readyState"
        }
    }

    // function createPeerConnection(targetUsername: string) {
    //     myPeerConnection = new RTCPeerConnection(myPeerConnectionConfig);
    //     myPeerConnection.onicecandidate = function (event: RTCPeerConnectionIceEvent) {
    //         if (event.candidate) {
    //             message.send(socket, {
    //                 type: "new-ice-candidate",
    //                 target: targetUsername,
    //                 candidate: event.candidate,
    //             });
    //         }
    //     };
    //
    //     myPeerConnection.ontrack = function (event: RTCTrackEvent) {
    //         console.log("TRACK", event, myRole);
    //     };
    //
    //     myPeerConnection.onnegotiationneeded = function (event: Event) {
    //         console.log("onnegotiationneeded", event);
    //     };
    //
    //     myPeerConnection.onremovetrack = function (event: any) {
    //         console.log('onremovetrack', event)
    //     };
    //
    //     myPeerConnection.oniceconnectionstatechange = function (event: Event) {
    //         console.log("oniceconnectionstatechange", event);
    //     };
    //
    //     myPeerConnection.onicegatheringstatechange = function (event: Event) {
    //         console.log("onicegatheringstatechange", event);
    //     };
    //
    //     myPeerConnection.onsignalingstatechange = function (event: Event) {
    //         console.log("onsignalingstatechange", event);
    //     };
    //
    //     myPeerConnection.onicecandidateerror = (event: Event) => {
    //         console.log("onicecandidateerror", event);
    //     };
    //     return myPeerConnection;
    // }
    //
    // function connectUser(targetUsername: string) {
    //     // console.log("Connect User", targetUsername);
    //     if (myPeerConnections[targetUsername]) {
    //         // console.log("Already Connected");
    //         return true;
    //     }
    //     let myPeerConnection = createPeerConnection(targetUsername);
    //
    //     // console.log("localStream", localStream);
    //     if (localStream) {
    //         localStream.getTracks().forEach((track) => {
    //             // track["mmcomp"] = "VIDEO";
    //             myPeerConnection.addTrack(track, localStream);
    //         });
    //         // console.log("Added tracks to connection");
    //     }
    //     for (const tuser in altStreams) {
    //         const tstreams = altStreams[tuser];
    //         tstreams[0].getTracks().forEach((track) => {
    //             myPeerConnection.addTrack(track, tstreams[0]);
    //         });
    //     }
    //     // if (shareScreenStream) {
    //     //   shareScreenStream.getTracks().forEach((track) => {
    //     //     // track["mmcomp"] = "VIDEO";
    //     //     myPeerConnection.addTrack(track, shareScreenStream);
    //     //   });
    //     //   // console.log("Added tracks to connection");
    //     // }
    //     // if (remoteStream) {
    //     //   remoteStream.getTracks().forEach((track) => {
    //     //     // track["mmcomp"] = "VIDEO";
    //     //     myPeerConnection.addTrack(track, remoteStream);
    //     //   });
    //     //   // console.log("Added tracks to connection");
    //     // }
    // }

    const onMessage = useCallback((message) => {
        let msg = messenger.receive(message);
        msg.data = (msg.Data && !msg.data) ? msg.Data : msg.data;
        msg.type = (msg.Type && !msg.type) ? msg.Type : msg.type;
        switch (msg.type as string) {
            case "alt-video-offer":
                // await altVideoOffer();
                break;
            case "new-ice-candidate":
                // newIceCandidate();
                break;
            case "alt-new-ice-candidate":
                // newAltIceCandidate();
                break;
            case "video-offer":
                console.log('###########');
                sendVideoAnswer(msg.data).then();
                break;
            case "video-answer":
                // videoAnswer();
                break;
            case "alt-video-answer":
                // newAltVideoAnswer();
                break;
            case "start":
                // setMyUsername(msg.Data);
                myUsername = msg.data;
                // start();
                break;
            case "role":
                // role();
                break;
            case "add_audience":
                // console.log("add_audience", msg.data);
                connectUser(msg.data);
                break;
            case "add_broadcast_audience":
                // console.log("add_audience", msg.data);
                // connectUser();
                break;
            case "tree":
                // tree();
                break;
            case "alt-_broadcast":
                // altBroadcast();
                break;
            default:
                console.log(msg);
        }
    }, [localStream]);

    useEffect(() => {
        socket.addEventListener("message", onMessage);
        // console.log('added socket onMessage listener');

        return () => {
            socket.removeEventListener("message", onMessage);
            // console.log('socket onMessage listener removed');
        };
    }, [socket, onMessage]);

    return (
        <div>
            <button onClick={(e) => console.log(localStream)}>Local</button>
            <Main myName={myName} myRole={myRole}/>

        </div>
        // <LocalStream/>
    )
}
