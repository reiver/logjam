import {useContext} from "react";
import {PeerConnectionArrayContext} from "../providers/PeerConnectionArrayProvider";

export const usePeerConnectionArray = () => {
    console.log('usePeerConnectionArray');
    return useContext(PeerConnectionArrayContext);
}
