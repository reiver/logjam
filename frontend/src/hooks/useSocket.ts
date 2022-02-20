import {SocketContext} from "../SocketProvider";
import {useContext} from "react";

export const useSocket = () => {
    console.log('useSocket');
    return useContext(SocketContext);
}
