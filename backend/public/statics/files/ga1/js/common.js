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
    if (myRole === 'broadcast') {
        document.getElementById('raise_hand').style.display = 'none';
        disableAudioVideoControls();
        return new SparkRTC('broadcast', {
            localStreamChangeCallback: (stream) => {
                getVideoElement('localVideo').srcObject = stream;
                enableAudioVideoControls();
            },
            remoteStreamCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (document.getElementById(tagId)) return;
                const video = createVideoElement(tagId);
                video.srcObject = stream;
                video.play();
            },
            remoteStreamDCCallback: (stream) => {
                let tagId = 'remoteVideo-' + stream.id;
                if (!document.getElementById(tagId)) {
                    tagId = 'localVideo-' + stream.id;
                    if (!document.getElementById(tagId)) return;
                }
                removeVideoElement(tagId);
            },
            signalingDisconnectedCallback: () => {
                clearScreen();
            },
            treeCallback: (tree) => {
                try {
                    const treeData = JSON.parse(tree);
                    if (!treeData) return;
                    graph.draw(treeData[0]);
                } catch (e) {
                    console.error(e);
                }
            },
            raiseHandConfirmation: (msg) => {
                console.log(`[raiseHandConfirmation] msg`, msg);
                return true;
                // if (confirm(msg)) {
                //     return true;
                // }
                // return false;
            },
            startProcedure: async () => {
                await handleClick();
            },
            log: (log) => {
                addLog(log);
            },
            constraintResults: (constraints) => {
                if (!constraints.audio) {
                    document.getElementById('mic').remove();
                }
            },
            updateStatus: (status) => {
                console.log(`[updateStatus] ${status}`);
                document.getElementById('status').innerText = status;
            },
            userListUpdated: (users) => {
                console.log('User List is updated', {users});
            },
        });
    } else {
        document.getElementById('share_screen').style.display = 'none';
        document.getElementById('mic').style.display = 'none';
        document.getElementById('camera').style.display = 'none';
        document.getElementById('sidebar-wrapper').style.display = 'none';
        const img = document.getElementById("raise_hand");
        return new SparkRTC('audience', {
            remoteStreamCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (document.getElementById(tagId)) return;
                const video = createVideoElement(tagId);
                video.srcObject = stream;
                video.play();
                document.getElementById('dc-place-holder')?.remove();
                img.dataset.status = 'on';
                img.src = RAISE_HAND_ON;
            },
            remoteStreamDCCallback: (stream) => {
                console.log(`[remoteStreamDCCallback]`, stream);
                if (stream !== 'no-stream') {
                    let tagId = 'remoteVideo-' + stream.id;
                    if (!document.getElementById(tagId)) {
                        tagId = 'localVideo-' + stream.id;
                        if (!document.getElementById(tagId)) return;
                    }
                    removeVideoElement(tagId);
                }
                if (sparkRTC.broadcasterDC || stream === 'no-stream') {
                    document.getElementById('screen').innerHTML = `<div id="dc-place-holder" style="display: block;">
                    <img style="width: 100%;" src="images/broken-link-mistake-error-disconnect-svgrepo-com.svg" />
                    <h1>Broadcaster is disconnected now, please stand by</h1>
                    </div>`;
                    img.dataset.status = 'off';
                    img.src = RAISE_HAND_OFF;
                    document.getElementById('mic').style.display = 'none';
                    document.getElementById('mic').src = MIC_ON;
                    document.getElementById('mic').dataset.status = 'on';
                    document.getElementById('camera').style.display = 'none';
                    document.getElementById('camera').src = CAMERA_ON;
                    document.getElementById('camera').dataset.status = 'on';
                }
            },
            signalingDisconnectedCallback: () => {
                clearScreen();
            },
            startProcedure: async () => {
                console.log('startProcedure');
                sparkRTC.stopSignaling();
                clearScreen();
                let idList = [];
                for (const id in sparkRTC.myPeerConnectionArray) {
                    const peerConn = sparkRTC.myPeerConnectionArray[id];
                    await peerConn.close();
                    idList.push(id)
                }
                idList.forEach((id) => delete sparkRTC.myPeerConnectionArray[id])
                sparkRTC.remoteStreams = [];
                sparkRTC.localStream?.getTracks()?.forEach(track => track.stop());
                sparkRTC.localStream = null;
                if (sparkRTC.startedRaiseHand) {
                    // sparkRTC.startedRaiseHand = false;
                    img.dataset.status = 'off';
                    img.src = RAISE_HAND_OFF;
                    document.getElementById('mic').style.display = 'none';
                    document.getElementById('mic').src = MIC_ON;
                    document.getElementById('mic').dataset.status = 'on';
                    document.getElementById('camera').style.display = 'none';
                    document.getElementById('camera').src = CAMERA_ON;
                    document.getElementById('camera').dataset.status = 'on';
                }
                await handleClick();
            },
            log: (log) => {
                addLog(log);
            },
            updateStatus: (status) => {
                console.log(`[updateStatus] ${status}`);
                document.getElementById('status').innerText = status;
            },
            userListUpdated: (users) => {
                console.log('User List is updated', {users});
            },
        });

    }
}

function registerNetworkEvent() {
    if (!navigator?.connection) {
        return alert('The browser is not a standard one so we can not monitor network status.');
    }
    handleNetworkStatus();
    navigator.connection.onchange = handleNetworkStatus;
}

function handleNetworkStatus(event) {
    const net = event?.currentTarget || navigator.connection;
    if (net.downlink <= 1) {
        console.log('Network is slow.');
        return onNetworkIsSlow(net.downlink);
    }
    onNetworkIsNormal();
}
