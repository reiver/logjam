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
// import Stream from "./Stream";

export const _Audience = ({myName}: { myName: string }) => {
    // const logger = useLogger();
    console.log('[Render] _Audience. myName=', myName);

    const [myUsername, setMyUsername] = useState('');
    // let myUsername: string = '';
    let myRole = "audience";

    let {peerConnectionMap, setPeerConnectionMap} = usePeerConnectionMap();
    let {streamMap, setStreamMap, enableLocalStream} = useStreamMap();
    const [ok, setOk] = useState(false);

    // const [remoteStream, setRemoteStream] = useState<MediaStream>();

    const messenger = useMessenger();
    const socket = useSocket();

    useEffect(() => {
        // start
        console.log('# 1 #');
        if (!myUsername) {
            messenger.send({
                    type: "start",
                    data: myName
                }
            );

        }

    });


    // useEffect(() => {
    // },[enableLocalStream, myRole, streamMap]);

    const createPeerConnection = useCallback(async (targetUsername: string) => {
        console.log('>>>>>> _Audience createPeerConnection: targetUsername=', targetUsername);
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
        // peerConnection.oniceconnectionstatechange = (_: Event) => {
        //     console.error('[PeerConnection] ################oniceconnectionstatechange')
        // };

        // onNegotiationNeeded
        peerConnection.onnegotiationneeded = function (event) {
            console.log("[PeerConnection] onNegotiationNeeded", event);

        };

        peerConnection.ontrack = function (event: RTCTrackEvent) {
            console.log('[PeerConnection] onTrack stream length:', event.streams.length);
            event.streams.forEach((stream) => {
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
            // setRemoteStream(event.streams[0]);
            // setRemoteStream(localStream);
            try {
                console.log('==(14)==> add remoteStream to setStreamMap')
                setStreamMap((
                    prev: Map<string, MediaStream>) => new Map(prev).set('remoteStream', stream));
            } catch (err) {
                console.error(err)
            }


            // messenger.send({
            //     type: "log",
            //     data: "ontrack from " + targetName,
            // });
        };
        console.log('PeerConnection', 'new peerConnection created:', peerConnection);

        return peerConnection;
    }, [messenger, setStreamMap]);

    // const addPeerConnection = useCallback((targetUsername: string) => {
    //     let newPeerConnection = createPeerConnection(targetUsername);
    //     setPeerConnectionMap((
    //         prev: Map<string, RTCPeerConnection>) => new Map(prev).set(targetUsername, newPeerConnection));
    //     console.log('[PeerConnection] new peerConnection added to peerConnectionMap');
    //     return newPeerConnection;
    // }, [createPeerConnection, peerConnectionMap, setPeerConnectionMap]);

    const connectUser = useCallback(async (targetUsername: string) => {
        console.log('##### createPeerConnection to ', targetUsername);
        if (peerConnectionMap.get(targetUsername)) {
            console.log('[PeerConnection] connection already exists.');
            return true;
        }

        let myPeerConnection = await createPeerConnection(targetUsername);
        setPeerConnectionMap((
            prev: Map<string, RTCPeerConnection>) => new Map(prev).set(targetUsername, myPeerConnection));
        myPeerConnection
            .createOffer()
            .then(function (offer) {
                return myPeerConnection.setLocalDescription(offer);
            })
            .then(function () {
                messenger.send({
                    name: myUsername,
                    target: targetUsername,
                    type: "video-offer",
                    sdp: myPeerConnection.localDescription,
                });
                console.log("[PeerConnection] _Audience sent video request (offer)");
            })
            .catch((err) => {
                messenger.send({
                    type: "log",
                    data: "onnegotiationneeded audience error:" + JSON.stringify(err),
                });
            });

        console.log('[PeerConnection] new peerConnection added to peerConnectionMap');


    }, [createPeerConnection, messenger, myUsername, peerConnectionMap, setPeerConnectionMap]);

    const sendVideoAnswer = useCallback(async (msg: any) => {
        console.log('[Socket] sendVideoAnswer');
        console.log(msg);
        let targetUsername = msg.name ? msg.name : '';
        // let targetName = msg.username;
        let myPeerConnection = await createPeerConnection(targetUsername);
        const desc = new RTCSessionDescription(msg.sdp);
        console.log('==(7)==> The recipient receives the offer and calls RTCPeerConnection.setRemoteDescription() to record it as the remote description')
        myPeerConnection.setRemoteDescription(desc).then( e=>
            {
                console.log('==(9)==> The recipient then creates an answer by calling RTCPeerConnection.createAnswer().')
                myPeerConnection
                .createAnswer()
                    .then((answer)=>{
                        console.log('==(10)==> The recipient calls RTCPeerConnection.setLocalDescription(), passing in the created answer, to set the answer as its local description. The recipient now knows the configuration of both ends of the connection.')
                        myPeerConnection.setLocalDescription(answer).then(e=>
                            {
                                console.log('==(11)==> The recipient uses the signaling server to send the answer to the caller.')
                                messenger.send({
                                    name: myUsername,
                                    target: targetUsername,
                                    type: "video-answer",
                                    sdp: myPeerConnection.remoteDescription,
                                });
                            }
                        )
                    })
            }
        )

    }, [createPeerConnection, messenger, myUsername]);


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
            case "new-ice-candidate":
                // newIceCandidate();
                break;
            case "video-offer":
                console.log('received offer from the broadcaster')
                sendVideoAnswer(msg).then();
                break;
            case "start":
                console.log('# 2 #');

                console.log('[SocketOnMessage] audience start');
                setMyUsername(msg.data);
                console.log(' myUsername:', msg.data);
                // console.log('[SetState] myUsername:', msg.data);


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


                // myUsername = msg.Data;
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
            case "abc":
                setOk(true);
                console.log('&&&&&&&&&&&&&')
                // tree();
                break;
            case "alt-_broadcast":
                // altBroadcast();
                break;
            default:
                console.log(msg);
        }
    }, [connectUser, messenger, myRole, sendVideoAnswer]);

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
            {/*<h1 style={{color: "white"}}>Local</h1>*/}
            {/*<Stream streamId={"localStream"} mediaStream={localStream}/>*/}
            {/*<h1 style={{color: "white"}}>Remote</h1>*/}
            {/*<Stream streamId={"remoteStream"} mediaStream={remoteStream}/>*/}
            <h1>Audience</h1>

            <Main myName={myName} myRole={myRole}/>

        </div>
    )
}
