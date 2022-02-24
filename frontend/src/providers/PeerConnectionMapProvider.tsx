import React, {useEffect, useState, createContext, ReactChild} from "react";

const defaultPeerConnectionMap =  new Map<string, RTCPeerConnection>();

export const PeerConnectionMapContext = createContext({
    peerConnectionMap: defaultPeerConnectionMap,
    setPeerConnectionMap: (peerConnectionMap:
                               (prev: Map<string, RTCPeerConnection>) => Map<string, RTCPeerConnection>) => {}
})

export const PeerConnectionMapProvider = (props: {
    children: ReactChild;
}) => {
    // console.log('PeerConnectionMapProvider');
    const [peerConnectionMap, setPeerConnectionMap] = useState(defaultPeerConnectionMap);


    useEffect(() => {
        console.log('peerConnections:', peerConnectionMap);
    }, [peerConnectionMap]);

    return (
        <PeerConnectionMapContext.Provider value={{peerConnectionMap,  setPeerConnectionMap}}>
            {props.children}
        </PeerConnectionMapContext.Provider>
    )
}
