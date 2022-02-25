import React, {useEffect, useRef} from 'react';
import {useLocalStream} from "../hooks/useLocalStream";
import {useLogger} from "../hooks/useLogger";
import RemoteStream from "./RemoteStream";

export default function LocalStream() {
    const logger = useLogger();
    logger.log('Render', 'LocalStream');

    let {localStream, setLocalStream} = useLocalStream();
    const videoRef = useRef<HTMLVideoElement>(null);
    const videoRef2 = useRef<HTMLVideoElement>(null);

    const MEDIA_CONSTRAINTS =  {
        audio: false,
        video: {
            width: 300,
            height: 200
        },
    }

    useEffect(() => {
        navigator.mediaDevices.getUserMedia(MEDIA_CONSTRAINTS)
            .then((mediaStream) => {
                    if (videoRef && videoRef.current) {
                        let video = videoRef.current as HTMLVideoElement;
                        video.srcObject = mediaStream;
                        setLocalStream(mediaStream);
                        console.log('#### localStream from getUserMedia:', mediaStream)
                        // video.onloadedmetadata = function (e) {
                        //     console.log('video loaded')
                        //     videoRef?.current?.play().then();
                        // };
                    }
                }
            )
    }, [videoRef]);

    useEffect(() => {
        if (videoRef2 && videoRef2.current) {
            let video2 = videoRef2.current as HTMLVideoElement;
            video2.srcObject = localStream;
        }
    },[localStream]);


    return (
        <div>
            <video id={"video 1"} ref={videoRef} autoPlay playsInline muted/>
            <RemoteStream mediaStream={localStream}/>

        </div>
    )

}