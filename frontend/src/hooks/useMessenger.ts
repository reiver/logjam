import {useSocket} from "./useSocket";

export const useMessenger = () => {

    let socket = useSocket();

    function send(message: any) {
        try {
            console.log('[Message] sent: ', message);
            const msg = JSON.stringify(message);
            socket.send(msg);
        } catch (e) {
            console.error('[MessageError] unable to send message. ', e);
            return;
        }
    }

    function receive(message: any) {
        if (!message) {
            console.error('[MessageError] received a message but it is null');
            return;
        }
        try {
            let msg = JSON.parse(message?.data);
            console.log('[Message] received: ', msg);
            return msg;
        } catch (e) {
            console.error('[MessageError] unable to parse received message.', e);
            return;
        }
    }

    return {send, receive};
}
