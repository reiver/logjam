import React, {useEffect, useState, createContext, ReactChild} from "react";
import {useLogger} from "../hooks/useLogger";

const defaultPeerConnectionMap =  new Map<string, RTCPeerConnection>();

export const PeerConnectionMapContext = createContext({
    peerConnectionMap: defaultPeerConnectionMap,
    setPeerConnectionMap: (peerConnectionMap:
                               (prev: Map<string, RTCPeerConnection>) => Map<string, RTCPeerConnection>) => {}
})

export const PeerConnectionMapProvider = (props: {
    children: ReactChild;
}) => {
    const logger = useLogger();

    const [peerConnectionMap, setPeerConnectionMap] = useState(defaultPeerConnectionMap);

    useEffect(() => {
        logger.log('PeerConnectionMap', peerConnectionMap);
    }, [peerConnectionMap]);

    return (
        <PeerConnectionMapContext.Provider value={{peerConnectionMap,  setPeerConnectionMap}}>
            {props.children}
        </PeerConnectionMapContext.Provider>
    )
}
