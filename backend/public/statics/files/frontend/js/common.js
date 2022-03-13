function getWsUrl() {
    const baseUrl = window.location.href.split("//")[1].split("/")[0];
    const protocol = window.location.href.split("//")[0] === "http:" ? "ws" :"wss";
    return `${protocol}://${baseUrl}/ws`
}

function addRemoteVideoTag(stream) {
    const tagId = 'remoteVideo-' + stream.id;
    if (document.getElementById(tagId)) return;

    const videoTag = document.createElement('video');
    videoTag.className = 'video-container remote-video';
    videoTag.playsInline = true;
    videoTag.autoplay = true;
    videoTag.id = tagId;
    videoTag.srcObject = stream;
    document.getElementById('screen').appendChild(videoTag);
    arrangeVideoContainers();
};

function clearRemoteVideos() {
    document.getElementsByClassName('remote-video').remove();
}

function createSparkRTC() {
    if (myRole === 'broadcast'){
        return new SparkRTC('broadcast', {
            localStreamChangeCallback: (stream) => {
                console.log('got stream', stream);
                document.getElementById('localVideo').srcObject = stream;
                console.log('Set Local Stream');
            },
            remoteStreamCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (document.getElementById(tagId)) return;

                addRemoteVideoTag(stream);
            },
            signalingDisconnectedCallback: () => {
                clearRemoteVideos()
            },
        });
    }
    document.getElementById('localVideo').remove();
    document.getElementById('localScreen').remove();
    return new SparkRTC('audience', {
        remoteStreamCallback: (stream) => {
            const tagId = 'remoteVideo-' + stream.id;
            if (document.getElementById(tagId)) return;

            addRemoteVideoTag(stream);
        },
        signalingDisconnectedCallback: () => {
            clearRemoteVideos();
        },
    });
}
