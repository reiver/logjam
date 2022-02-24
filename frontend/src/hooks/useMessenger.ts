import {useSocket} from "./useSocket";
import {useLogger} from "./useLogger";

export const useMessenger = () => {

    let socket = useSocket();
    let logger = useLogger();

    function send(message: any) {
        try {
            logger.log('Message', 'sent: ', message);
            const msg = JSON.stringify(message);
            socket.send(msg);
        } catch (e) {
            logger.error('Message',  'unable to send message: ', e);
            return;
        }
    }

    function receive(message: any) {
        if (!message) {
            logger.error('Message', 'received a message but it is null');
            return;
        }
        try {
            let msg = JSON.parse(message?.data);
            logger.log('Message',  'received: ', msg);
            return msg;
        } catch (e) {
            logger.error('Message', 'unable to parse received message. ', e);
            return;
        }
    }

    return {send, receive};
}
