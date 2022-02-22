import {useCallback, useEffect} from "react";
import {useSocket} from "../hooks/useSocket";
import Main from "./Main";
import LocalStream from "./LocalStream";
import {useParams} from "react-router-dom";
import {PEER_CONNECTION_CONFIG, turnStatus} from "../config/myPeerConnectionConfig";
import {useMessenger} from "../hooks/useMessenger";

export const Socket = ({myName}: { myName: string }) => {
    console.log('Socket component rendered');

    let {myRole} = useParams() || "audience";
    let myPeerConnection: any;
    let myPeerConnections: any = {};

    const messenger = useMessenger();
    const socket = useSocket();
    console.log(socket);

    useEffect(() => {
        messenger.send({
                type: "start",
                data: myName
            }
        );
        console.log('start sent');

        messenger.send({
                type: "role",
                data: myRole
            }
        );
        console.log('role sent');

        messenger.send({
                type: "turn_status",
                data: turnStatus()
            }
        );
        console.log('role sent');

    }, []);

    useEffect(() => {
        if (myRole === 'audience') {
            let myPeerConnection = new RTCPeerConnection(PEER_CONNECTION_CONFIG);
            console.log(myPeerConnection);

            myPeerConnection.onicecandidate = function (event) {
                if (event.candidate) {
                    messenger.send({
                        type: "new-ice-candidate",
                        target: 'targetUsername',
                        candidate: event.candidate,
                    });
                }
            };

            myPeerConnection.onnegotiationneeded = function (event) {
                console.log("onnegotiationneeded", event);
                myPeerConnection
                    .createOffer()
                    .then(function (offer) {
                        return myPeerConnection.setLocalDescription(offer);
                    })
                    .then(function () {
                        messenger.send({
                            name: 'myUsername',
                            target: 'targetUsername',
                            type: "video-offer",
                            sdp: myPeerConnection.localDescription,
                        });

                    })
                    .catch(function (e) {
                        // console.log("onnegotiationneeded audience error:", e);
                        messenger.send({
                            type: "log",
                            data: "onnegotiationneeded audience error:" + JSON.stringify(e),
                        });
                    });
            };
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
        switch (msg.Type as string) {
            case "alt-video-offer":
                // await altVideoOffer();
                break;
            case "new-ice-candidate":
                // newIceCandidate();
                break;
            case "alt-new-ice-candidate":
                // newAltIceCandidate();
                break;
            case "video-answer":
                // videoAnswer();
                break;
            case "alt-video-answer":
                // newAltVideoAnswer();
                break;
            case "start":
                // start();
                break;
            case "role":
                // role();
                break;
            case "add_audience":
                // console.log("add_audience", msg.data);
                // connectUser();
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
    }, []);

    useEffect(() => {
        socket.addEventListener("message", onMessage);
        console.log('added socket onMessage listener');

        return () => {
            socket.removeEventListener("message", onMessage);
            console.log('socket onMessage listener removed');
        };
    }, [socket, onMessage]);

    return (
        <Main myName={myName} myRole={myRole}/>
        // <LocalStream/>
    )
}
