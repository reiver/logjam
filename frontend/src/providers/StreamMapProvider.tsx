import React, {useEffect, useState, createContext, ReactChild} from "react";
import {useLogger} from "../hooks/useLogger";

const defaultStreamMap =  new Map<string, MediaStream>();

export const StreamMapContext = createContext({
    streamMap: defaultStreamMap,
    setStreamMap: (_: (prev: Map<string, MediaStream>) => Map<string, MediaStream>) => {}})

export const StreamMapProvider = (props: {
    children: ReactChild;
}) => {
    const logger = useLogger();

    const [streamMap, setStreamMap] = useState(defaultStreamMap);


    useEffect(() => {
        logger.log('StreamsMap', streamMap);
    }, [streamMap]);

    return (
        <StreamMapContext.Provider value={{streamMap,  setStreamMap}}>
            {props.children}
        </StreamMapContext.Provider>
    )
}
