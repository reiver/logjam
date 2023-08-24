function getWsUrl() {
    const baseUrl = window.location.href.split("//")[1].split("/")[0];
    const protocol = window.location.href.split("//")[0] === "http:" ? "ws" : "wss";
    return `${protocol}://${baseUrl}/ws`
}


function clearScreen() {
    document.getElementById('raise_hand').style.display = 'block';
    document.getElementById('share_screen').style.display = 'block';
    let screen = document.getElementById('screen');
    if (!screen) return;
    while (screen.hasChildNodes()) {
        screen.removeChild(screen.firstChild);
    }
}

function createSparkRTC() {
    clearScreen();
    document.getElementById('raise_hand').style.display = 'none';
    const sparkRTC = new SparkRTC('unknown', {
        localStreamChangeCallback: (stream) => {
            console.log('Local Stream', stream);
            if (stream) {
                getVideoElement('localVideo').srcObject = stream;
                if (document.getElementById("mic").dataset.status === 'off')
                {
                    sparkRTC.disableAudio();
                }
                if (document.getElementById("camera").dataset.status === 'off')
                {
                    sparkRTC.disableVideo();
                }
            }
        },
        remoteStreamCallback: (stream, target) => {
            console.log('remote stream', stream);
            const tagId = 'remoteVideo-' + stream.id;
            if (document.getElementById(tagId)) return;
            const video = createVideoElement(tagId, stream.target ||  target);
            video.srcObject = stream;
            video.play();
            if (sparkRTC.role === 'audience' && !sparkRTC.localStream && adminAccess) {
                // if (!sparkRTC.remoteStreams.includes(re => rs.id === sparkRTC.localStreamId))
                // if (confirm('?'))
                    sparkRTC.raiseHand();
            }
        },
        remoteStreamDCCallback: (stream) => {
            console.log('remoteStreamDCCallback', stream);
            let tagId = 'remoteVideo-' + stream.id;
            if (!document.getElementById(tagId)) {
                tagId = 'localVideo-' + stream.id;
                if (!document.getElementById(tagId)) {
                  console.log('Stream not found');
                  return;
                }
            }
            removeVideoElement(tagId);
        },
        signalingDisconnectedCallback: () => {
            clearScreen();
        },
        startProcedure: () => {
            // return //new Promise((resolve) => {
                // setTimeout(() => {
                    onLoad();
                    return handleClick();
                    // resolve();
                // }, 100);
            // });
        },
        raiseHandConfirmation: (msg) => {
            return true;
        }
    });

    return sparkRTC;
}