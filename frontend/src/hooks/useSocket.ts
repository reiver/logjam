import {SocketContext} from "../providers/SocketProvider";
import {useContext} from "react";

export const useSocket = () => {
    console.log('useSocket');
    return useContext(SocketContext);
}
