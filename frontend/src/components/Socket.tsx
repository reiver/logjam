import {useCallback, useEffect} from "react";
import {useSocket} from "../hooks/useSocket";
import Main from "./Main";
import LocalVideo from "./LocalVideo";
import {useParams} from "react-router-dom";
import {receiveMessage, sendMessage} from "../helpers/message";
import myPeerConnectionConfig from "../config/myPeerConnectionConfig";

export const Socket = ({myName}: { myName: string }) => {
    let {myRole} = useParams() || "audience";
    let myPeerConnection: any;
    let myPeerConnections: any = {};

    console.log('Socket component');

    const socket = useSocket();
    console.log(socket);

    useEffect(() => {
        sendMessage(socket, {
                type: "start",
                data: myName
            }
        );
        console.log('start sent');

        sendMessage(socket, {
                type: "role",
                data: myRole
            }
        );
        console.log('role sent');


    }, []);

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
    //             sendMessage(socket, {
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
        let msg = receiveMessage(message);
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
            case "alt-broadcast":
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
        // <LocalVideo/>
    )
}
