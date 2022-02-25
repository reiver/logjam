import React, {useEffect, useState} from "react";
import {MEDIA_CONSTRAINTS} from "../config/mediaConstraints";

export const useUserMedia = () => {
    const [localStream, setLocalStream] = useState<MediaStream>();

    useEffect(() => {
        const enableStream = async () => {
            try {
                const stream = await navigator.mediaDevices.getUserMedia(MEDIA_CONSTRAINTS)
                setLocalStream(stream)
            } catch (err) {
                console.error(err)
            }
        }

        if (!localStream) {
            enableStream().then();
        } else {
            return () => {
                if (localStream) {
                    localStream.getTracks().forEach(track => {
                        track.stop()
                    })
                }
            }
        }
    }, [localStream])

    return localStream;
}
