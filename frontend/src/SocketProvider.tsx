import React, {useEffect, useState, createContext, ReactChild, useCallback} from "react";

const SOCKET_RECONNECTION_TIMEOUT = 30;

interface ISocketProvider {
    myName: string,
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
    console.log('myName: ', props.myName);
    const [ws, setWs] = useState<WebSocket>(webSocket);

    // onOpen Handler
    const onOpen = useCallback((event) => {
        let data = {
            type: "start",
            data: props.myName
        };
        console.log('[onOpen] sent: ', data);
        ws.send(JSON.stringify(data));
    }, []);

    // onClose Handler
    const onClose = useCallback((event) => {
        console.log(
            `[socket closed] Connection ${event.wasClean ?
                "closed cleanly" : "died"}, code=${event.code} reason=${event.reason}`
        );
        setTimeout(() => {
            setWs(new WebSocket(getSocketUrl()));
        }, SOCKET_RECONNECTION_TIMEOUT);
    }, []);

    // onError Handler
    const onError = useCallback((error) => {
        console.log('[socket error] ', error);
    }, []);


    useEffect(() => {
        ws.addEventListener("open", onOpen);
        console.log('added socket onOpen listener');

        return () => {
            ws.removeEventListener("open", onOpen);
            console.log('socket onOpen listener removed');
        };
    }, [ws]);

    useEffect(() => {
        ws.addEventListener("close", onClose);
        console.log('added socket onClose listener');

        return () => {
            ws.removeEventListener("close", onClose);
            console.log('socket onClose listener removed');
        };
    }, [ws]);

    useEffect(() => {
        ws.addEventListener("error", onError);
        console.log('added socket onError listener');

        return () => {
            ws.removeEventListener("error", onError);
            console.log('socket onError listener removed');
        };
    }, [ws]);

    return (
        <SocketContext.Provider value={ws}>
            {props.children}
        </SocketContext.Provider>
    )
}
