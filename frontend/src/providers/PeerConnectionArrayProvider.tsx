import React, {useEffect, useState, createContext, ReactChild} from "react";

interface PeerConnectionArrayContextValue {
    peerConnectionArray: RTCPeerConnection[],
    setPeerConnectionArray: (peerConnectionArray: RTCPeerConnection[]) => void
}

const ContextValueInitials: PeerConnectionArrayContextValue = {
    peerConnectionArray: [],
    setPeerConnectionArray: () => {}
}

export const PeerConnectionArrayContext = createContext<PeerConnectionArrayContextValue>(ContextValueInitials);

export const PeerConnectionArrayProvider = (props: {
    children: ReactChild;
}) => {
    console.log('PeerConnectionArrayProvider');
    const [peerConnectionArray, setPeerConnectionArray] = useState<RTCPeerConnection[]>([]);


    useEffect(() => {
        console.log('peerConnections:', peerConnectionArray);
    }, [peerConnectionArray]);

    return (
        <PeerConnectionArrayContext.Provider value={
            {peerConnectionArray, setPeerConnectionArray}
        }>
            {props.children}
        </PeerConnectionArrayContext.Provider>
    )
}
