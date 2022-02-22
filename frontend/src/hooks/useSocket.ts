import {SocketContext} from "../providers/SocketProvider";
import {useContext} from "react";

export const useSocket = () => {
    return useContext(SocketContext);
}
