import React, {useEffect, useState, createContext, ReactChild} from "react";
import {MEDIA_CONSTRAINTS} from "../config/mediaConstraints";

const defaultStreamMap =  new Map<string, MediaStream>();

export const StreamMapContext = createContext({
    streamMap: defaultStreamMap,
    setStreamMap: (_: (prev: Map<string, MediaStream>) => Map<string, MediaStream>) => {},
    enableLocalStream: () => new Promise(function(p1: (value: (PromiseLike<any>)) => void,p2: (reason?: any) => void){})
})

export const StreamMapProvider = (props: {
    children: ReactChild;
}) => {

    const [streamMap, setStreamMap] = useState(defaultStreamMap);

    const enableLocalStream = async () => {
        console.log('enableLocalStream');
        try {
            const stream = await navigator.mediaDevices.getUserMedia(MEDIA_CONSTRAINTS)
            setStreamMap((
                prev: Map<string, MediaStream>) => new Map(prev).set('localStream', stream));
        } catch (err) {
            console.error(err)
        }
    }

    useEffect(() => {
        console.log('[StreamsMap] ', streamMap);
    }, [streamMap]);

    return (
        <StreamMapContext.Provider value={{streamMap,  setStreamMap, enableLocalStream}}>
            {props.children}
        </StreamMapContext.Provider>
    )
}
