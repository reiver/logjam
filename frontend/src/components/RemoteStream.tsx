import React, {useEffect, useRef} from 'react';
import {useLocalStream} from "../hooks/useLocalStream";

export default function RemoteStream({mediaStream}: { mediaStream: MediaStream }) {
    console.log('[Render] RemoteStream');

    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        console.log('>>>>>', mediaStream);
        if (videoRef && videoRef.current && mediaStream) {
            const video = videoRef.current as HTMLVideoElement;
            video.srcObject = mediaStream;
            video.onloadedmetadata = function(e) {
                video.play().then();
            };
        }
    }, [videoRef, mediaStream]);

    return (
        <video id="remote" style={{border: "5px solid yellow", backgroundColor: "white"}} ref={videoRef} autoPlay playsInline muted/>
    )

}