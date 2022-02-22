import {useContext} from "react";
import {PeerConnectionArrayContext} from "../providers/PeerConnectionArrayProvider";

export const usePeerConnectionArray = () => {
    return useContext(PeerConnectionArrayContext);
}
