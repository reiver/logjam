import React, {useCallback, useEffect, useState} from "react";
import {useSocket} from "../hooks/useSocket";
import Main from "./Main";
// import {useParams} from "react-router-dom";
import {PEER_CONNECTION_CONFIG, turnStatus} from "../config/myPeerConnectionConfig";
import {useMessenger} from "../hooks/useMessenger";
import {usePeerConnectionMap} from "../hooks/usePeerConnectionMap";
// import {useLocalStream} from "../hooks/useLocalStream";
// import Data from "../data";
// import RemoteStream from "./RemoteStream";
// import {useLogger} from "../hooks/useLogger";
import {useStreamMap} from "../hooks/useStreamMap";
// import Stream from "./Stream";

export const Broadcast = ({myName}: { myName: string }) => {
    // const logger = useLogger();
    console.log('[Render] Broadcast. myName=', myName);

    const [myUsername, setMyUsername] = useState('');
    // let myUsername: string = '';
    let myRole = "broadcast";

    let {peerConnectionMap, setPeerConnectionMap} = usePeerConnectionMap();
    let {streamMap, setStreamMap, enableLocalStream} = useStreamMap();

    // const [remoteStream, setRemoteStream] = useState<MediaStream>();

    const messenger = useMessenger();
    const socket = useSocket();

    useEffect(() => {
        // start
        console.log('# 1 #');
        if (!myUsername){
            messenger.send({
                    type: "start",
                    data: myName
                }
            );

        }

    });


    // useEffect(() => {
    // },[enableLocalStream, myRole, streamMap]);

    const createPeerConnection = useCallback((targetUsername: string) => {
        console.log('>>>>>> createPeerConnection: targetUsername=', targetUsername);
        let peerConnection = new RTCPeerConnection(PEER_CONNECTION_CONFIG);

        // onIceCandidate
        peerConnection.onicecandidate = function (event) {
            // console.log('[PeerConnection] onIceCandidate')
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
            // messenger.send({
            //     type: "log",
            //     data: "onicecandidateerror :" + JSON.stringify(event),
            // });
        };

        // onIceCandidateError
        peerConnection.oniceconnectionstatechange = (_: Event) => {
            // console.error('[PeerConnection] ################oniceconnectionstatechange')
        };

        // onNegotiationNeeded
        peerConnection.onnegotiationneeded = function (event) {
            console.log("[PeerConnection] onNegotiationNeeded", event);
            peerConnection
                .createOffer()
                .then(function (offer) {
                    console.log('==(3)==> The caller calls RTCPeerConnection.createOffer() to create an offer.')
                    console.log('==(4)==> The caller calls RTCPeerConnection.setLocalDescription() to set that offer as the local description (that is, the description of the local end of the connection).')
                    return peerConnection.setLocalDescription(offer);
                })
                .then(function () {
                    console.log('==(6)==> The caller uses the signaling server to transmit the offer to the intended receiver of the call.')
                    messenger.send({
                        name: myUsername,
                        target: targetUsername,
                        type: "video-offer",
                        sdp: peerConnection.localDescription,
                    });

                })
                .catch((err) => {
                    console.error('error in createOffer()')
                    messenger.send({
                        type: "log",
                        data: "onnegotiationneeded audience error:" + JSON.stringify(err),
                    });
                });
        };



        console.log('PeerConnection', 'new peerConnection created:', peerConnection);

        return peerConnection;
    }, [messenger, myUsername]);

    const videoAnswerReceived = useCallback((msg: any) => {
        console.log('==(12)==> The caller receives the answer.')
        const desc = new RTCSessionDescription(msg.sdp);
        let myPeerConnection = peerConnectionMap.get(msg.data);
        try{
            if (myPeerConnection){
                myPeerConnection.setRemoteDescription(desc).then(e=>{
                    console.log('==(13)==> The caller calls RTCPeerConnection.setRemoteDescription() to set the answer as the remote description for its end of the call. It now knows the configuration of both peers. Media begins to flow as configured.')
                } )

            }else{
                console.error('myPeerConnection is undefined')
            }
        }catch (e) {
            console.log("Error in setRemoteDescription", e);

        }
    }, [peerConnectionMap]);

    // const addPeerConnection = useCallback((targetUsername: string) => {
    //     let newPeerConnection = createPeerConnection(targetUsername);
    //     setPeerConnectionMap((
    //         prev: Map<string, RTCPeerConnection>) => new Map(prev).set(targetUsername, newPeerConnection));
    //     console.log('[PeerConnection] new peerConnection added to peerConnectionMap');
    //     return newPeerConnection;
    // }, [createPeerConnection, peerConnectionMap, setPeerConnectionMap]);

    const connectUser = useCallback((targetUsername: string) => {
        console.log('##### connectUser');
        if (peerConnectionMap.get(targetUsername)) {
            console.log('[PeerConnection] connection already exists.');
            return true;
        }

        let myPeerConnection = createPeerConnection(targetUsername);
        setPeerConnectionMap((
            prev: Map<string, RTCPeerConnection>) => new Map(prev).set(targetUsername, myPeerConnection));
        console.log('[PeerConnection] new peerConnection added to peerConnectionMap');

        let localStream = streamMap.get('localStream');
        if (localStream) {
            console.log('Adding tracks... len=', localStream.getTracks().length);
            localStream.getTracks().forEach((track) => {
                if (myPeerConnection && localStream) {
                    myPeerConnection.addTrack(track, localStream);
                    console.log('==(2)==> The caller creates RTCPeerConnection and calls RTCPeerConnection.addTrack()')
                }else{
                    console.error('unable to add tracks')
                }
            });
        }
    }, [createPeerConnection, peerConnectionMap, setPeerConnectionMap, streamMap]);



    const onMessage = useCallback((message) => {
        let msg = messenger.receive(message);
        msg.data = (msg.Data && !msg.data) ? msg.Data : msg.data;
        msg.type = (msg.Type && !msg.type) ? msg.Type : msg.type;
        switch (msg.type as string) {
            case "new-ice-candidate":
                // newIceCandidate();
                break;

            case "video-answer":
                console.log('~~~~~ received video answer from the audience')
                videoAnswerReceived(msg);
                break;
            case "start":
                console.log('# 2 #');

                console.log('[SocketOnMessage] start');
                setMyUsername(msg.data);
                // myUsername = msg.data;
                console.log(' myUsername:', msg.data);
                // console.log('[SetState] myUsername:', msg.data);

                if (!streamMap.get('localStream')){
                    console.log('# 3 #');
                    console.log('==(1)==> The caller captures local Media via MediaDevices.getUserMedia')
                    enableLocalStream().then(_=>
                        {
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

                        }
                    );
                }
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
    }, [connectUser, enableLocalStream, messenger, myRole, streamMap, videoAnswerReceived]);

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
            <h1>Broadcast</h1>
            <Main myName={myName} myRole={myRole}/>
        </div>
    )
}
