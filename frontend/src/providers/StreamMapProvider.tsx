import React, {useEffect, useState, createContext, ReactChild} from "react";
import {MEDIA_CONSTRAINTS} from "../config/mediaConstraints";

const defaultStreamMap =  new Map<string, MediaStream>();

export const StreamMapContext = createContext({
    streamMap: defaultStreamMap,
    setStreamMap: (_: (prev: Map<string, MediaStream>) => Map<string, MediaStream>) => {},
    enableLocalStream: () => {}
})

export const StreamMapProvider = (props: {
    children: ReactChild;
}) => {

    const [streamMap, setStreamMap] = useState(defaultStreamMap);

    const enableLocalStream = async () => {
        try {
            const stream = await navigator.mediaDevices.getUserMedia(MEDIA_CONSTRAINTS)
            setStreamMap((
                prev: Map<string, MediaStream>) => new Map(prev).set('localStream', stream));
        } catch (err) {
            console.error(err)
        }
    }

    useEffect(() => {
        console.log('[StreamsMap]', streamMap);
    }, [streamMap]);

    return (
        <StreamMapContext.Provider value={{streamMap,  setStreamMap, enableLocalStream}}>
            {props.children}
        </StreamMapContext.Provider>
    )
}
