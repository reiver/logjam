import React, {useEffect, useState, createContext, ReactChild} from "react";

export const LocalStreamContext = createContext({
    localStream: new MediaStream(),
    setLocalStream: (localStream: MediaStream)=>{}
});

export const LocalStreamProvider = (props: {
    children: ReactChild;
}) => {
    console.log('LocalStreamProvider');
    const [localStream, setLocalStream] = useState<MediaStream>(new MediaStream());


    useEffect(() => {
        console.log('localStream:', localStream);
    }, [localStream]);

    return (
        <LocalStreamContext.Provider value={{localStream, setLocalStream}}>
            {props.children}
        </LocalStreamContext.Provider>
    )
}
