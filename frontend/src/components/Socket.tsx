import {useCallback, useEffect} from "react";
import {useSocket} from "../hooks/useSocket";

export const Socket = () => {
    console.log('Socket component');

    const socket = useSocket();
    console.log(socket);

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

    const onMessage = useCallback((message) => {
        if (!message) {
            console.error('[onMessage] Socket message is null');
            return;
        }

        const msg = JSON.parse(message?.data);
        console.log('[message received]', msg);
        switch (msg.Type) {
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

        return () => {
            socket.removeEventListener("message", onMessage);
        };
    }, [socket, onMessage]);

    return (
        <button
            onClick={() => {
                console.log('message sent')
                socket.send(
                    JSON.stringify({
                        type: "start",
                        data: 'myName'
                    })
                );
            }}
        >Start</button>
    )
}
