import React, { useEffect, useState, createContext, ReactChild } from "react";

const SOCKET_URL = 'ws://localhost:8080/ws';
const SOCKET_RECONNECTION_TIMEOUT = 30;
const webSocket = new WebSocket(SOCKET_URL);

export const SocketContext = createContext(webSocket);

interface ISocketProvider {
    children: React.ReactChild;
}

export const SocketProvider = (props: ISocketProvider) => {
    console.log('SocketProvider');
    const [ws, setWs] = useState<WebSocket> (webSocket);

    useEffect(() => {
        const onClose = () => {
            console.log('closed');
            setTimeout(() => {
                setWs(new WebSocket(SOCKET_URL));
            }, SOCKET_RECONNECTION_TIMEOUT);
        };

        ws.addEventListener("close", onClose);

        return () => {
            console.log('clean')
            ws.removeEventListener("close", onClose);
        };
    }, [ws, setWs]);

    return (
        <SocketContext.Provider value={ws}>{props.children}</SocketContext.Provider>
    )
}
