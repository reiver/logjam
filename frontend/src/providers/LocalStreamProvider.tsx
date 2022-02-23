import React, {useEffect, useState, createContext, ReactChild} from "react";

const defaultLocalStream = new MediaStream();

export const LocalStreamContext = createContext({
    localStream: defaultLocalStream,
    setLocalStream: (localStream: MediaStream)=>{}
});

export const LocalStreamProvider = (props: {
    children: ReactChild;
}) => {
    console.log('LocalStreamProvider');
    const [localStream, setLocalStream] = useState<MediaStream>(defaultLocalStream);


    useEffect(() => {
        console.log('*** localStream (provider)', localStream);
    }, [localStream]);

    return (
        <LocalStreamContext.Provider value={{localStream, setLocalStream}}>
            {props.children}
        </LocalStreamContext.Provider>
    )
}
