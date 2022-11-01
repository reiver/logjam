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
                getVideoElement('localVideo').srcObject = stream;
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
                    tagId = 'localVideo-' + stream.id ;
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
                return true;
            },
            startProcedure: async () => {
                await handleClick();
            },
            log: (log) => {
                addLog(log);
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
                video.play();
                document.getElementById('dc-place-holder').remove();
            },
            remoteStreamDCCallback: (stream) => {
                let tagId = 'remoteVideo-' + stream.id;
                if (!document.getElementById(tagId)) {
                    tagId = 'localVideo-' + stream.id ;
                    if (!document.getElementById(tagId)) return;
                }
                removeVideoElement(tagId);
                document.getElementById('screen').innerHTML = `<div id="dc-place-holder" style="display: block;">
                <img style="width: 100%;" src="images/broken-link-mistake-error-disconnect-svgrepo-com.svg" />
                <h1>Broadcaster is disconnected now, wait for it...</h1>
                </div>`;
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
            startProcedure: async () => {
                await handleClick();
            },
            log: (log) => {
                addLog(log);
            },
        });

    }
}
