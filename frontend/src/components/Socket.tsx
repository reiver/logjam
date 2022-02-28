import React, {useCallback, useEffect, useState} from "react";
import {useSocket} from "../hooks/useSocket";
import Main from "./Main";
import {useParams} from "react-router-dom";
import {PEER_CONNECTION_CONFIG, turnStatus} from "../config/myPeerConnectionConfig";
import {useMessenger} from "../hooks/useMessenger";
import {usePeerConnectionMap} from "../hooks/usePeerConnectionMap";
// import {useLocalStream} from "../hooks/useLocalStream";
// import Data from "../data";
// import RemoteStream from "./RemoteStream";
// import {useLogger} from "../hooks/useLogger";
import {useStreamMap} from "../hooks/useStreamMap";
import {useUserMedia} from "../hooks/useUserMedia";
import Stream from "./Stream";

export const Socket = ({myName}: { myName: string }) => {
    // const logger = useLogger();
    console.log('[Render] Socket. myName=', myName);

    const [myUsername, setMyUsername] = useState('');
    // let myUsername: string;
    let {myRole} = useParams() || "audience";

    let {peerConnectionMap, setPeerConnectionMap} = usePeerConnectionMap();
    let {streamMap, setStreamMap} = useStreamMap();
    let localStream = useUserMedia();

    const [remoteStream, setRemoteStream] = useState<MediaStream>();

    const messenger = useMessenger();
    const socket = useSocket();

    useEffect(() => {
        // start
        messenger.send({
                type: "start",
                data: myName
            }
        );

    });


    const createPeerConnection = useCallback((targetUsername: string) =>{
        console.log('>>>>>> createPeerConnection: targetUsername=', targetUsername);
        let peerConnection = new RTCPeerConnection(PEER_CONNECTION_CONFIG);

        // onIceCandidate
        peerConnection.onicecandidate = function (event) {
            console.log('[PeerConnection] onIceCandidate')
            if (event.candidate) {
                messenger.send({
                    type: "new-ice-candidate",
                    target: targetUsername,
                    candidate: event.candidate,
                });
            }
        };

        // onIceCandidateError
        peerConnection.onicecandidateerror = (event: Event) => {
            // console.error('[PeerConnection] onIceCandidateError')
            messenger.send({
                type: "log",
                data: "onicecandidateerror :" + JSON.stringify(event),
            });
        };

        // onIceCandidateError
        peerConnection.oniceconnectionstatechange = (event: Event) => {
            console.error('[PeerConnection] ################oniceconnectionstatechange')
        };

        // onNegotiationNeeded
        peerConnection.onnegotiationneeded = function (event) {
            console.log("[PeerConnection] onNegotiationNeeded", event);
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

        peerConnection.ontrack = function (event: RTCTrackEvent) {
            console.log('[PeerConnection] onTrack stream length:', event.streams.length);
            event.streams.map((stream) => {
                console.log('stream id', stream.id);
                console.log('stream track length', stream.getTracks().length);
                messenger.send({
                    type: "stream",
                    data: "true",
                });
            });

            // if (!hasStream) {
            //     hasStream = true;
            //     socket.send(
            //         JSON.stringify({
            //             type: "stream",
            //             data: "true",
            //         })
            //     );
            // }
            let stream = event.streams[0];
            console.log('ontrack', stream);
            console.dir('------------tracks', stream.getTracks());
            // event.track.onended = e => setRemoteStream(event.streams[0]);
            setRemoteStream(event.streams[0]);
            // setRemoteStream(localStream);

            messenger.send({
                type: "log",
                data: "ontrack from " + 'targetName',
            });
        };

        console.log('PeerConnection', 'new peerConnection created:', peerConnection);

        return peerConnection;
    },[messenger, myUsername]);

    const addPeerConnection = useCallback((targetUsername: string) => {
        if (peerConnectionMap.get(targetUsername)) {
            console.log('[PeerConnection] connection already exists.');
            return peerConnectionMap.get(targetUsername);
        }
        let newPeerConnection = createPeerConnection(targetUsername);
        setPeerConnectionMap((
            prev: Map<string, RTCPeerConnection>) => new Map(prev).set(targetUsername, newPeerConnection));
        console.log('[PeerConnection] new peerConnection added to peerConnectionMap');
        return newPeerConnection;
    }, [createPeerConnection, peerConnectionMap, setPeerConnectionMap]);

    const connectUser = useCallback((targetUsername: string) => {
        console.log('##### connectUser');
        let myPeerConnection = addPeerConnection(targetUsername);
        if (localStream) {
            console.log('Adding tracks... len=', localStream.getTracks().length);
            localStream.getTracks().forEach((track) => {
                if (myPeerConnection && localStream) {
                    myPeerConnection.addTrack(track, localStream);
                }
            });
        }
    }, [addPeerConnection, localStream]);

    const sendVideoAnswer = useCallback(async (msg: any) => {
        console.log('[Socket] sendVideoAnswer');
        console.log(msg);
        let targetUsername = msg.name ? msg.name : '';
        let targetName = msg.username;
        const desc = new RTCSessionDescription(msg.sdp);
        let myPeerConnection = addPeerConnection(targetUsername);
        console.log(myPeerConnection);
        if (myPeerConnection) {
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
    }, [addPeerConnection, messenger, myUsername]);

    const videoAnswerReceived = useCallback((msg: any) => {
        console.log('[Socket] videoAnswerReceived')
        const desc = new RTCSessionDescription(msg.sdp);
        let myPeerConnection = peerConnectionMap.get(myUsername);
        myPeerConnection?.setRemoteDescription(desc).catch((e) => {
            console.log("Error", e);
        });
    }, [myUsername, peerConnectionMap]);


    // useEffect(() => {
    //     if (myRole === 'audience') {
    //
    //     }
    // }, [myRole])

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
                sendVideoAnswer(msg).then();
                break;
            case "video-answer":
                videoAnswerReceived(msg);
                break;
            case "alt-video-answer":
                // newAltVideoAnswer();
                break;
            case "start":
                console.log('[SocketOnMessage] start');
                setMyUsername(msg.data);
                console.log('[SetState] myUsername');
                // myUsername = msg.Data;
                messenger.send({
                        type: "role",
                        data: myRole
                        // data: myRole==='broadcast' ? 'alt-broadcast' : myRole
                    }
                );

                messenger.send({
                        type: "turn_status",
                        data: turnStatus()
                    }
                );
                // myUsername = msg.data;
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
    }, [connectUser, messenger, myRole, sendVideoAnswer, videoAnswerReceived]);

    useEffect(() => {
        socket.addEventListener("message", onMessage);
        // console.log('added socket onMessage listener');

        return () => {
            socket.removeEventListener("message", onMessage);
            // console.log('socket onMessage listener removed');
        };
    }, [socket, onMessage]);

    // {myRole === 'audience' ?
    //     <div>
    //         <h3 style={{color: "white"}}>Remote Video:</h3>
    //         <RemoteStream mediaStream={remoteStream}/>
    //     </div>
    //     : null
    // }

    return (
        <div>
            {/*<button onClick={e => setRemoteStream(remoteStream)}>Set Remote</button>*/}
            {/*<div>*/}
            {/*    <h3 style={{color: "white"}}>Remote Video:</h3>*/}
            {/*    <RemoteStream mediaStream={remoteStream}/>*/}
            {/*</div>*/}
            <h1 style={{color: "white"}}>Local</h1>
            <Stream streamId={"localStream"} mediaStream={localStream}/>
            <h1 style={{color: "white"}}>Remote</h1>
            <Stream streamId={"remoteStream"} mediaStream={remoteStream}/>
            {/*<Main myName={myName} myRole={myRole}/>*/}

        </div>
    )
}
