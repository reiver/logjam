import React, {useEffect, useState, createContext, ReactChild, useCallback} from "react";

const SOCKET_RECONNECTION_TIMEOUT = 5000;
const defaultWebSocket = new WebSocket(getSocketUrl());

function getSocketUrl() {
    return 'ws://localhost:8080/ws'
    // const url = window.location.href.split("//");
    // const protocol = url[0] === "http:" ? "wss" : "ws";
    // const baseUrl = url[1].split("/")[0];
    // return `${protocol}://${baseUrl}/ws`;
}

export const SocketContext = createContext(defaultWebSocket);

export const SocketProvider = (props: {
    children: ReactChild;
}) => {
    // console.log('SocketProvider');

    const [webSocket, setWebSocket] = useState<WebSocket>(defaultWebSocket);

    // onOpen Handler
    const onOpen = useCallback((event) => {
        console.log('[socket opened]');
    }, []);

    // onClose Handler
    const onClose = useCallback((event) => {
        console.log(
            `[socket closed] Connection ${event.wasClean ?
                "closed cleanly" : "died"}, code=${event.code} reason=${event.reason}`
        );
        setTimeout(() => {
            setWebSocket(new WebSocket(getSocketUrl()));
        }, SOCKET_RECONNECTION_TIMEOUT);
    }, []);

    // onError Handler
    const onError = useCallback((error) => {
        console.log('[socket error] ', error);
    }, []);


    useEffect(() => {
        webSocket.addEventListener("open", onOpen);
        // console.log('added socket onOpen listener');

        return () => {
            webSocket.removeEventListener("open", onOpen);
            // console.log('socket onOpen listener removed');
        };
    }, [webSocket, onOpen]);

    useEffect(() => {
        webSocket.addEventListener("close", onClose);
        // console.log('added socket onClose listener');

        return () => {
            webSocket.removeEventListener("close", onClose);
            // console.log('socket onClose listener removed');
        };
    }, [webSocket, onClose]);

    useEffect(() => {
        webSocket.addEventListener("error", onError);
        // console.log('added socket onError listener');

        return () => {
            webSocket.removeEventListener("error", onError);
            // console.log('socket onError listener removed');
        };
    }, [webSocket, onError]);

    return (
        <SocketContext.Provider value={webSocket}>
            {props.children}
        </SocketContext.Provider>
    )
}
