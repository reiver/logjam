import React, {useEffect, useRef} from 'react';
import {useLogger} from "../hooks/useLogger";

export default function Stream({streamId, mediaStream}: { streamId: string, mediaStream: MediaStream | undefined }) {
    const logger = useLogger();
    logger.log('Render', 'Stream');
    logger.log('Stream', streamId, mediaStream);

    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        if (videoRef && videoRef.current && mediaStream) {
            let video = videoRef.current as HTMLVideoElement;
            video.srcObject = mediaStream;
            console.log('---------------srcObject:', streamId, video.srcObject)
            console.dir(video);
        }
        return () => {
            if (mediaStream) {
                console.log('---------- CLEANUP')
                mediaStream.getTracks().forEach(track => {
                    track.stop()
                })
            }
        }
    }, [mediaStream]);

    if (!mediaStream) {
        return null;
    }

    return (
        <div>
            <video id={streamId} style={{border: "5px solid blue", backgroundColor: "white"}}
                   ref={videoRef} autoPlay playsInline muted/>

        </div>
    )

}