function getWsUrl() {
    const baseUrl = window.location.href.split("//")[1].split("/")[0];
    const protocol = window.location.href.split("//")[0] === "http:" ? "ws" :"wss";
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
    if (myRole === 'broadcast'){
        document.getElementById('raise_hand').style.display = 'none';
        return new SparkRTC('broadcast', {
            localStreamChangeCallback: (stream) => {
                console.log('got stream', stream);
                getVideoElement('localVideo').srcObject = stream;
                // getVideoElement('localVideo2').srcObject = stream;
                // getVideoElement('localVideo3').srcObject = stream;
                // getVideoElement('localVideo4').srcObject = stream;
                // getVideoElement('localVideo5').srcObject = stream;
                // getVideoElement('localVideo6').srcObject = stream;
                // getVideoElement('localVideo7').srcObject = stream;
                console.log('Set Local Stream');
            },
            remoteStreamCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (document.getElementById(tagId)) return;
                const video = createVideoElement(tagId);
                video.srcObject = stream;
            },
            remoteStreamDCCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (!document.getElementById(tagId)) return;
                console.log('disconnected remote stream', tagId);
                removeVideoElement(tagId);
            },
            signalingDisconnectedCallback: () => {
                clearScreen();
            },
            treeCallback: (tree) => {
                try {
                    const treeData = JSON.parse(tree);
                    graph.draw(treeData[0]);
                } catch (e) {
                    console.error(e);
                }
            },
            raiseHandConfirmation: (msg) => {
                console.log('raise hand confirmation', msg);
                return true;
            }
        });
    }else{
        document.getElementById('share_screen').style.display = 'none';
        return new SparkRTC('audience', {
            remoteStreamCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (document.getElementById(tagId)) return;
                const video = createVideoElement(tagId);
                video.srcObject = stream;
            },
            remoteStreamDCCallback: (stream) => {
                const tagId = 'remoteVideo-' + stream.id;
                if (!document.getElementById(tagId)) return;
                console.log('disconnected remote stream', tagId);
                removeVideoElement(tagId);
            },
            signalingDisconnectedCallback: () => {
                clearScreen();
            },
            treeCallback: (tree) => {
                try {
                    const treeData = JSON.parse(tree);
                    graph.draw(treeData[0]);
                } catch (e) {
                    console.error(e);
                }
            },
            startProcedure: async () => {
                await handleClick();
            },
        });

    }
}
