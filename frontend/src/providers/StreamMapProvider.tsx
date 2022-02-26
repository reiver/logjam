import React, {useEffect, useState, createContext, ReactChild} from "react";

const defaultStreamMap =  new Map<string, MediaStream>();

export const StreamMapContext = createContext({
    streamMap: defaultStreamMap,
    setStreamMap: (_: (prev: Map<string, MediaStream>) => Map<string, MediaStream>) => {}})

export const StreamMapProvider = (props: {
    children: ReactChild;
}) => {

    const [streamMap, setStreamMap] = useState(defaultStreamMap);

    useEffect(() => {
        console.log('[StreamsMap]', streamMap);
    }, [streamMap]);

    return (
        <StreamMapContext.Provider value={{streamMap,  setStreamMap}}>
            {props.children}
        </StreamMapContext.Provider>
    )
}
