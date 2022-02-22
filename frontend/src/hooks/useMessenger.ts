import {useSocket} from "./useSocket";

export const useMessenger = () => {
    console.log('useMessage rendered');

    let socket = useSocket();

    function send(message: any) {
        try {
            const msg = JSON.stringify(message);
            console.log('[Info] message sent: ', msg);
            socket.send(msg);
        } catch (e) {
            console.error('[Error] unable to send message. ', e);
            return;
        }
    }

    function receive(message: any) {
        if (!message) {
            console.error('[Error] received a message but it is null');
            return;
        }
        try {
            const msg = JSON.parse(message?.data);
            console.log('[Info] message received: ', msg);
            return msg;
        } catch (e) {
            console.error('[Error] unable to parse received message.', e);
            return;
        }
    }

    return {send, receive};
}
