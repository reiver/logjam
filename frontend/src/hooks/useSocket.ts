import {useEffect, useState} from "react";
import {SocketMessage} from "../types/SocketMessage";

export default function useSocket({myName}: any) {
    const [mySocket, setMySocket] = useState<WebSocket>();
    // let socket: (WebSocket | undefined) = undefined;

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
        if (myName && !mySocket) {
            createSocket();
        }
    });

    function getSocketUrl() {
        const url = window.location.href.split("//");
        const protocol = url[0] === "http:" ? "wss" : "ws";
        const baseUrl = url[1].split("/")[0];
        return `${protocol}://${baseUrl}/ws`;
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
        let socket = new WebSocket(getSocketUrl());

        socket.onopen = function (event) {
            // console.log("[open] Connection established");
            // console.log("Sending start to server");
            if (!socket) return;
            socket.send(
                JSON.stringify({
                    type: "start",
                    data: myName,
                })
            );
        };

        socket.onclose = function (event) {
            if (event.wasClean) {
                console.log(
                    `[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`
                );
            } else {
                console.log(
                    `[close] Connection died, , code=${event.code} reason=${event.reason}`
                );
            }
            // myUsername = "";
            // myPeerConnections = [];
            // myAltPeerConnections = [];
            // myPeerConnection = null;
            // login();
        };

        socket.onerror = function (error) {
            // console.log(`[error] ${error.message}`);
        };

        socket.onmessage = async function (event) {
            let msg: SocketMessage = parseSocketMessage(event.data);
            if (!msg) return;
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
        setMySocket(socket);
    }
    return mySocket;
}


