import React, {useEffect, useState, createContext, ReactChild} from "react";

const SOCKET_RECONNECTION_TIMEOUT = 30;

interface ISocketProvider {
    children: ReactChild;
}

function getSocketUrl() {
    return 'ws://localhost:8080/ws'
    // const url = window.location.href.split("//");
    // const protocol = url[0] === "http:" ? "wss" : "ws";
    // const baseUrl = url[1].split("/")[0];
    // return `${protocol}://${baseUrl}/ws`;
}

const webSocket = new WebSocket(getSocketUrl());
export const SocketContext = createContext(webSocket);

export const SocketProvider = (props: ISocketProvider) => {
    console.log('SocketProvider');
    const [ws, setWs] = useState<WebSocket>(webSocket);

    useEffect(() => {
        const onClose = (event: { wasClean: boolean; code: number; reason: string; }) => {
            console.log(
                `[socket closed] Connection ${event.wasClean ?
                    "closed cleanly" : "died"}, code=${event.code} reason=${event.reason}`
            );
            setTimeout(() => {
                setWs(new WebSocket(getSocketUrl()));
            }, SOCKET_RECONNECTION_TIMEOUT);
        };

        const onError = (error: any) =>{
            console.log('[socket error] ', error);
        }

        ws.addEventListener("close", onClose);
        ws.addEventListener("error", onError);

        return () => {
            console.log('socket listeners removed')
            ws.removeEventListener("close", onClose);
            ws.removeEventListener("error", onError);
        };
    }, [ws, setWs]);

    return (
        <SocketContext.Provider value={ws}>
            {props.children}
        </SocketContext.Provider>
    )
}
