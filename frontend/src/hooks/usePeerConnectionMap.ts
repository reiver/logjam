import {useContext} from "react";
import {PeerConnectionMapContext} from "../providers/PeerConnectionMapProvider";

export const usePeerConnectionMap = () => {
    return useContext(PeerConnectionMapContext);
}
