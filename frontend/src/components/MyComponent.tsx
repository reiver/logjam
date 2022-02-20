import {useCallback, useEffect} from "react";
import {useSocket} from "../hooks/useSocket";

export const MyComponent = () => {
    const socket = useSocket();
    console.log(socket);
    const onMessage = useCallback((message) => {
        const data = JSON.parse(message?.data);
        console.log('[message received]', data);
        // ... Do something with the data
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
        >Send</button>
    )
}
