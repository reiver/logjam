import React, {useEffect, useRef} from 'react';

export default function Stream({streamId, mediaStream}: { streamId: string, mediaStream: MediaStream | undefined }) {
    console.log('[Render] Stream');
    // console.log('[Stream]', streamId, mediaStream);

    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        if (videoRef && videoRef.current && mediaStream) {
            let video = videoRef.current as HTMLVideoElement;
            video.srcObject = mediaStream;
            console.log('[html]', streamId + ' added as srcObject')
        }
        return () => {
            if (mediaStream) {
                console.log('[CLEANUP]')
                // mediaStream.getTracks().forEach(track => {
                //     track.stop()
                // })
            }
        }
    }, [mediaStream, streamId]);

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