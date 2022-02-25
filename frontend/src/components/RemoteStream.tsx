import React, {useEffect, useRef} from 'react';
import {useLocalStream} from "../hooks/useLocalStream";
import {useLogger} from "../hooks/useLogger";

export default function RemoteStream({mediaStream}: { mediaStream: MediaStream }) {
    const logger = useLogger();
    logger.log('Render', 'RemoteStream');

    const divRef = useRef<HTMLDivElement>(null);
    // const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        logger.log('Inside RemoteStream', mediaStream);
        let isPlaying = true;
        // if (videoRef && videoRef.current && mediaStream) {
        if (divRef && divRef.current && mediaStream && mediaStream.active) {
            let v = document.createElement('video')
            v.srcObject = mediaStream;
            v.pause();
            v.load();
            v.style.setProperty('border','5px solid yellow');
            setTimeout( function() {
                ///video.load();
                v.play().then(r=>console.log('-----play') );
            }, 1000);
            divRef.current.appendChild(v);
            console.log('---------video added')
            // let video = videoRef.current as HTMLVideoElement;
            // video.srcObject = mediaStream;

            // On video playing toggle values
            // video.onplaying = function () {
            //     isPlaying = true;
            // };
            //
            // // On video pause toggle values
            // video.onpause = function () {
            //     isPlaying = false;
            // };
            //
            //
            // console.log('-------srcObject:', video.srcObject);
            //
            // if (video.paused && !isPlaying) {
            //     console.log('----------play')
            //     video
            //         .play()
            //         .then(() => console.log(`--------video play`))
            //         .catch((error) =>
            //             console.log(`-------videoElem.play() failed: ${error}`))
            // }else{
            //     console.log('video.paused=', video.paused);
            //     console.log('isPlaying:', isPlaying);
            // }


        }
    }, [divRef, mediaStream]);

    // <video id="remote" style={{border: "5px solid yellow", backgroundColor: "white"}} ref={videoRef} controls muted/>
    return (
        <div ref={divRef}/>
    )

}