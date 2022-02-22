import React, {useEffect, useRef} from 'react';
import {useLocalStream} from "../hooks/useLocalStream";

export default function LocalStream() {
    let {localStream, setLocalStream} = useLocalStream();
    const videoRef = useRef<HTMLVideoElement>(null);

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
                        const video = videoRef.current as HTMLVideoElement;
                        video.srcObject = mediaStream;
                        setLocalStream(mediaStream);
                        // video.onloadedmetadata = function (e) {
                        //     console.log('video loaded')
                        //     videoRef?.current?.play().then();
                        // };
                    }
                }
            )
    }, [videoRef]);

    return (
        <video id={"video 1"} ref={videoRef} autoPlay playsInline muted/>
    )

}