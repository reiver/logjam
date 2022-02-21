import React, {useEffect, useRef} from 'react';
import {MEDIA_CONSTRAINTS} from "../config/mediaConstraints";

export default function LocalVideo() {
    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        navigator.mediaDevices.getUserMedia(MEDIA_CONSTRAINTS)
            .then((mediaStream) => {
                    if (videoRef && videoRef.current) {
                        videoRef.current.srcObject = mediaStream;
                        videoRef.current.onloadedmetadata = function (e) {
                            videoRef?.current?.play().then(r => console.log(r));
                        };
                    }
                }
            )
    }, []);

    return (
        <video ref={videoRef} autoPlay/>
    )

}