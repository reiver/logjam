import {useCallback, useEffect} from "react";
import {useSocket} from "../hooks/useSocket";
import Main from "./Main";
import LocalVideo from "./LocalVideo";
import {useParams} from "react-router-dom";
import {receiveMessage, sendMessage} from "../helpers/message";

export const Socket = ({myName}: {myName: string}) => {
    let { myRole } = useParams() || "audience";

    console.log('Socket component');

    const socket = useSocket();
    console.log(socket);

    useEffect(() => {
        sendMessage(socket,{
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


    const onMessage = useCallback((message) => {
        let msg = receiveMessage(message);
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
