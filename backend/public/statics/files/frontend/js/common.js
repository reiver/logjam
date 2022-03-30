function getWsUrl() {
    const baseUrl = window.location.href.split("//")[1].split("/")[0];
    const protocol = window.location.href.split("//")[0] === "http:" ? "ws" :"wss";
    return `${protocol}://${baseUrl}/ws`
}


function clearScreen() {
    let screen = document.getElementById('screen');
    if (!screen) return;
    while (screen.hasChildNodes()) {
        screen.removeChild(screen.firstChild);
    }
}

function createSparkRTC() {
    clearScreen();
    if (myRole === 'broadcast'){
        return new SparkRTC('broadcast', {
            localStreamChangeCallback: (stream) => {
                console.log('got stream', stream);
                getVideoElement('localVideo').srcObject = stream;
                console.log('Set Local Stream');
            },
            remoteStreamCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (document.getElementById(tagId)) return;
                const video = createVideoElement(tagId);
                video.srcObject = stream;
            },
            signalingDisconnectedCallback: () => {
                clearScreen();
            },
        });
    }else{
        return new SparkRTC('audience', {
            remoteStreamCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (document.getElementById(tagId)) return;
                const video = createVideoElement(tagId);
                video.srcObject = stream;
            },
            signalingDisconnectedCallback: () => {
                clearScreen();
            },
        });

    }
}
