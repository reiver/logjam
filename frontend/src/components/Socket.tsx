import {useEffect, useState} from "react";
import {SocketMessage} from "../types/SocketMessage";

export default function Socket({myName, mySocket, setMySocket}: any) {

    const [message, setMessage] = useState<SocketMessage>(
        {
            type: '',
            data: null,
            error: null,
            candidate: null,
            name: '',
            username: '',
            sdp: null,
            target: null
        }
    );

    useEffect(() => {
        console.log("myName=", myName);
        if (!myName) {
            console.log('myName is not set yet');
            return;
        }
        if (mySocket) {
            console.log('socket already exists');
            return;
        }
        createSocket();
    });

    function getStateDescription(state: number) {
        switch (state) {
            case 0:
                return "CONNECTING	Socket has been created. The connection is not yet open.";
            case 1:
                return "OPEN	The connection is open and ready to communicate."
            case 2:
                return "CLOSING	The connection is in the process of closing."
            case 3:
                return "CLOSED	The connection is closed or couldn't be opened."
            default:
                return ""
        }
    }

    function getSocketUrl() {
        return 'ws://localhost:8080/ws'
        // const url = window.location.href.split("//");
        // const protocol = url[0] === "http:" ? "wss" : "ws";
        // const baseUrl = url[1].split("/")[0];
        // return `${protocol}://${baseUrl}/ws`;
    }

    function parseSocketMessage(data: string) {
        let msg;
        try {
            msg = JSON.parse(data);
            if (msg.type !== 'tree') {
                console.log("Message", msg);
            }
        } catch (e) {
            console.error('Unable to parse socket message: ', e);
            return;
        }
        if (!msg) {
            console.error('Socket message is null');
            return;
        }
        return msg;
    }

    function createSocket() {
        console.log('createSocket()');
        let socket: WebSocket;
        try {
            socket = new WebSocket(getSocketUrl());
            console.log('socket created:', socket);
        } catch (e) {
            console.log('could not create socket:', e);
            return;
        }

        socket.onopen = function (event) {
            console.log("[open] Connection established");
            console.log("Sending start to server");
            socket.send(
                JSON.stringify({
                    type: "start",
                    data: myName,
                })
            );
        };

        socket.onclose = function (event) {
            console.log(
                `[close] Connection ${event.wasClean ? "closed": "died"}, code=${event.code} reason=${event.reason}`
            );
            // myUsername = "";
            // myPeerConnections = [];
            // myAltPeerConnections = [];
            // myPeerConnection = null;
            // login();
        };

        socket.onerror = function (error) {
            console.log('error:');
            console.dir(error);
        };

        socket.onmessage = async function (event) {
            let msg: SocketMessage = parseSocketMessage(event.data);
            if (!msg) {
                console.log('message is null');
                return;
            }

            setMessage(prevState => ({
                ...prevState, msg
            }));

            switch (msg.type) {
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
        }
        console.log('mySocket is set')
        setMySocket(socket);
    }

    return (
        <p style={{color: "white"}}>Socket {mySocket ? getStateDescription(mySocket.readyState) : "-"}</p>
    );
}


