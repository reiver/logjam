export function isWebRTCSupported(){
    if (window.navigator.userAgent.includes("Edge")) return false;

    return navigator.getUserMedia ||
        navigator.webkitGetUserMedia ||
        navigator.mozGetUserMedia ||
        navigator.msGetUserMedia ||
        window.RTCPeerConnection;
}

