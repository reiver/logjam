import PEER_CONNECTION_CONFIG from "../config/myPeerConnectionConfig";

const MEDIA_CONSTRAINTS = (
    window.constraints = {
        audio: true,
        video: true,
    }
);

export function _broadcast() {
    // let localStream;
    // navigator.mediaDevices
    //     .getUserMedia(MEDIA_CONSTRAINTS)
    //     .then((ls) => {
    //         localStream = ls;
    //         // console.log("Local Stream", localStream);
    //         socket.send(
    //             JSON.stringify({
    //                 type: "role",
    //                 data: "_broadcast",
    //             })
    //         );
    //         const turnStatus = m
    //         socket.send(
    //             JSON.stringify({
    //                 type: "turn_status",
    //                 data: turnStatus ? "on" : "off",
    //             })
    //         );
    //     })
    //     .catch((e) => {
    //         console.log("Local St error", e);
    //     });
}