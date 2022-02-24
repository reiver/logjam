import React, {useEffect, useRef} from 'react';
import {useLocalStream} from "../hooks/useLocalStream";
import {useLogger} from "../hooks/useLogger";

export default function RemoteStream({mediaStream}: { mediaStream: MediaStream }) {
    const logger = useLogger();
    logger.log('Render', 'RemoteStream');

    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        logger.log('Inside RemoteStream', mediaStream);
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