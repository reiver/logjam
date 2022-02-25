import React, {useEffect, useRef, useState} from 'react';
import {useLocalStream} from "../hooks/useLocalStream";
import {useLogger} from "../hooks/useLogger";

export default function RemoteStream({mediaStream}: { mediaStream: MediaStream }) {
    const logger = useLogger();
    logger.log('Render', 'RemoteStream');

    // let {remoteStream, setRemoteStream} = useState<MediaStream>(new MediaStream);
    // let {localStream, setLocalStream} = useLocalStream();

    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        if (videoRef && videoRef.current) {
            let video = videoRef.current as HTMLVideoElement;
            video.srcObject = mediaStream;
            // console.log('=====>localStream: ', localStream);
            console.log('=====>mediaStream received by remote stream: ', mediaStream);
        }
    },[mediaStream]);

    // useEffect(() => {
    //     function play() {
    //
    //         if (videoRef && videoRef.current && mediaStream) {
    //             let video = videoRef.current as HTMLVideoElement;
    //             setTimeout(() => {
    //                 console.log("wait 1 sec");
    //                 video.srcObject = mediaStream;
    //                 setTimeout(() => {
    //                     console.log(video.srcObject)
    //                     console.log("wait 2 sec");
    //                     video.play().then().catch(e => {
    //                         console.log(e)
    //                     });
    //                 }, 2000);
    //             }, 1000);
    //         }
    //     }
    //
    //     play();
    //     return () => {
    //
    //         if (videoRef && videoRef.current && mediaStream) {
    //             let video = videoRef.current as HTMLVideoElement;
    //             video.pause();
    //         }
    //
    //
    //     }
    // }, [videoRef, mediaStream]);


    // useEffect(() => {
    //     logger.log('Inside RemoteStream', mediaStream);
    //     if (videoRef && videoRef.current && mediaStream) {
    //         let video = videoRef.current as HTMLVideoElement;
    //         video.srcObject = mediaStream;
    //         console.log('-------Received srcObject:', video.srcObject);
    //
    //         // On video playing toggle values
    //         // video.onplaying = function () {
    //         //     isPlaying = true;
    //         // };
    //         //
    //         // // On video pause toggle values
    //         // video.onpause = function () {
    //         //     isPlaying = false;
    //         // };
    //         //
    //         //
    //         //
    //         // if (video.paused && !isPlaying) {
    //         //     console.log('----------play')
    //         //     video
    //         //         .play()
    //         //         .then(() => console.log(`--------video play`))
    //         //         .catch((error) =>
    //         //             console.log(`-------videoElem.play() failed: ${error}`))
    //         // }else{
    //         //     console.log('video.paused=', video.paused);
    //         //     console.log('isPlaying:', isPlaying);
    //         // }
    //
    //
    //     }
    // }, [videoRef, mediaStream]);

    return (
        <div>
            {/*<button onClick={play}>Play</button>*/}
            <video id="remote" style={{border: "5px solid yellow", backgroundColor: "white"}} ref={videoRef} autoPlay playsInline  muted/>

        </div>
    )

}