const CAMERA_ON = "images/cam-on.png";
const CAMERA_OFF = "images/cam-off.png";
const MIC_ON = "images/mic-on.png";
const MIC_OFF = "images/mic-off.png";
const SCREEN_ON = "images/screen-on.png";
const SCREEN_OFF = "images/screen-off.png";
// const SPARK_LOGO = "images/spark-logo.png";
const CHARACTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

let graph;
let sparkRTC;
let myName;
let myRole;
let shareScreenStream;
let roomName;
let adminAccess = false;

function makeId(length) {
    let result = '';
    for (let i = 0; i < length; i++) {
        result += CHARACTERS.charAt(Math.floor(Math.random() * CHARACTERS.length));
    }
    return result;
}


function arrangeVideoContainers() {
    const videoContainers = document.getElementById('screen')
        .getElementsByClassName('video-container');
    const videoCount = videoContainers.length;
    const flexGap = 1;
    let flexRatio = 100 / Math.ceil(Math.sqrt(videoCount));
    let flex = "0 0 " + flexRatio + "%";
    let maxHeight = 100 / Math.ceil(videoCount / Math.ceil(Math.sqrt(videoCount)));
    Array.from(videoContainers).forEach(div => {
        div.style.setProperty('flex', flex);
        div.style.setProperty('max-height', maxHeight + "%");
    }
    )
}


function onCameraButtonClick() {
    const img = document.getElementById("camera");
    if (img.dataset.status === 'on') {
        img.dataset.status = 'off';
        img.src = CAMERA_OFF;
        sparkRTC.disableVideo();
    } else {
        img.dataset.status = 'on';
        img.src = CAMERA_ON;
        sparkRTC.disableVideo(true);
    }
    callStatusUpdate();
}


function onMicButtonClick() {
    const img = document.getElementById("mic");
    if (img.dataset.status === 'on') {
        img.dataset.status = 'off';
        img.src = MIC_OFF;
        sparkRTC.disableAudio();
    } else {
        img.dataset.status = 'on';
        img.src = MIC_ON;
        sparkRTC.disableAudio(true);
    }
    callStatusUpdate();
}

function createVideoElement(videoId, infoText = '', muted = false) {
    let container = document.createElement('div');
    if (infoText !== '') {
        let info = document.createElement('div');
        info.innerHTML = infoText;
        info.className = 'video-info';
        container.appendChild(info);
    }
    container.className = 'video-container';
    let video = document.createElement('video');
    video.id = videoId;
    video.autoplay = true;
    video.playsInline = true;
    video.muted = muted;
    container.appendChild(video);
    document.getElementById('screen').appendChild(container);
    arrangeVideoContainers();
    return video;
}

function getVideoElement(videoId) {
    let video = document.getElementById(videoId);
    return video ? video : createVideoElement(videoId, '', true);
}

function removeVideoElement(videoId) {
    let video = document.getElementById(videoId);
    if (!video) return;
    let videoContainer = video.parentNode;
    if (!videoContainer) return;
    document.getElementById('screen').removeChild(videoContainer);
    arrangeVideoContainers();
}


async function onShareScreen() {
    const img = document.getElementById("share_screen");
    if (!shareScreenStream) {
        shareScreenStream = await sparkRTC.startShareScreen();
        if (shareScreenStream) {
            img.dataset.status = 'on';
            img.src = SCREEN_ON;
            const localScreen = getVideoElement('localScreen');
            localScreen.srcObject = shareScreenStream;
        }
    } else {
        img.dataset.status = 'off';
        img.src = SCREEN_OFF;
        shareScreenStream.getTracks().forEach(track => track.stop());
        shareScreenStream = null;
        const localScreen = getVideoElement('localScreen');
        localScreen.srcObject = null;
        removeVideoElement('localScreen');
    }
    callStatusUpdate();
}


function setMyName() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const nameParam = urlParams.get('name');
    if (nameParam) {
        document.getElementById('inputName').value = nameParam;
        try {
            localStorage.setItem('logjam_myName', nameParam);
        } catch (e) {
            console.log(e);
        }

        return true;
    }
    try {
        myName = localStorage.getItem('logjam_myName');
        document.getElementById('inputName').value = myName;
    } catch (e) {
        console.log(e);
    }
    if (myName === '' || !myName) {
        myName = makeId(20);
        try {
            localStorage.setItem('logjam_myName', myName);
        } catch (e) {
            console.log(e);
        }
    }
    return false;
}

function showHideControls() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    const hideParam = urlParams.get('hide_controls');
    if (hideParam) {
        document.getElementById("control-buttons-container").style.display = "none";
    }
}

async function handleClick() {
    let newName = document.getElementById("inputName").value;
    if (newName) {
        myName = newName;
        localStorage.setItem('logjam_myName', myName);
    }
    document.getElementById("page").style.visibility = "visible";
    document.getElementById("getName").style.display = "none";

    await start();

    return false;
}


function handleResize() {
    clearTimeout(window.resizedFinished);
    window.resizedFinished = setTimeout(function () {
        graph.draw(graph.treeData);
        arrangeVideoContainers();
    }, 250);

}

function getMyRole() {
    return "unknown";
}

function getRoomName() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return urlParams.get('room');
}

function setupSignalingSocket() {
    return sparkRTC.setupSignalingSocket(getWsUrl(), myName, roomName);
}

function isAdmin() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return !Boolean(urlParams.get('is_call'));
}

async function start() {
    await setupSignalingSocket();
    return sparkRTC.start();
}


function onLoad() {
    showHideControls();
    myRole = getMyRole();
    roomName = getRoomName();
    sparkRTC = createSparkRTC();
    adminAccess = isAdmin();
    if (!adminAccess) {
        document.getElementById('raise_hand').style.display = '';
    }

    setMyName();
    graph = new Graph();
    window.onresize = handleResize;
    graph.draw(DATA);

    arrangeVideoContainers();

    window.addEventListener("message", async (event) => {
        console.log("Inside Received message: ", event);
        try {
            const msg = JSON.parse(event.data);
            const micImg = document.getElementById("mic");
            const camImg = document.getElementById("camera");
            const scImg = document.getElementById("share_screen");

            switch (msg.type) {
                case 'GET_AUDIO_VIDEO_STATUS':
                    const { data } = msg;
                    data.isMyAudioMuted = micImg.dataset.status === 'off';
                    data.isMyVideoHidden = camImg.dataset.status === 'off';
                    const newMsg = { data, type: 'GOT_AUDIO_VIDEO_STATUS' };
                    parent.postMessage(JSON.stringify(newMsg), "*");
                    break;
                case 'SCREENSHARE_STATUS':
                    const scNewMsg = {
                        data: {
                            commandId: msg.data.commandId,
                            status: scImg.dataset.status === 'on' ? 'MY_SCREENSHARE_ACTIVATED' : 'MY_SCREENSHARE_AVAILABLE',
                        },
                        type: 'GOT_SCREENSHARE_STATUS',
                    };
                    console.log(scNewMsg);
                    parent.postMessage(JSON.stringify(scNewMsg), "*");
                    break;
                case 'MUTE_AUDIO':
                    if (micImg.dataset.status === 'on') {
                        micImg.dataset.status = 'off';
                        micImg.src = MIC_OFF;
                        sparkRTC.disableAudio();
                        parent.postMessage(
                            JSON.stringify({
                                type: 'MY_AUDIO_MUTED', data: {
                                    isMyAudioMuted: document.getElementById("mic").dataset.status === 'off',
                                    isMyVideoHidden: document.getElementById("camera").dataset.status === 'off',
                                    screenShareState: document.getElementById("share_screen").dataset.status === 'on' ? 'MY_SCREENSHARE_ACTIVATED' : 'MY_SCREENSHARE_AVAILABLE',
                                }
                            })
                            , "*");
                    }
                    break;
                case 'UNMUTE_AUDIO':
                    if (micImg.dataset.status === 'off') {
                        micImg.dataset.status = 'on';
                        micImg.src = MIC_ON;
                        sparkRTC.disableAudio(true);
                        parent.postMessage(
                            JSON.stringify({
                                type: 'MY_AUDIO_UNMUTED', data: {
                                    isMyAudioMuted: document.getElementById("mic").dataset.status === 'off',
                                    isMyVideoHidden: document.getElementById("camera").dataset.status === 'off',
                                    screenShareState: document.getElementById("share_screen").dataset.status === 'on' ? 'MY_SCREENSHARE_ACTIVATED' : 'MY_SCREENSHARE_AVAILABLE',
                                }
                            })
                            , "*");
                    }
                    break;
                case 'HIDE_VIDEO':
                    if (camImg.dataset.status === 'on') {
                        camImg.dataset.status = 'off';
                        camImg.src = CAMERA_OFF;
                        sparkRTC.disableVideo();
                        parent.postMessage(
                            JSON.stringify({
                                type: 'MY_VIDEO_HIDDEN', data: {
                                    isMyAudioMuted: document.getElementById("mic").dataset.status === 'off',
                                    isMyVideoHidden: document.getElementById("camera").dataset.status === 'off',
                                    screenShareState: document.getElementById("share_screen").dataset.status === 'on' ? 'MY_SCREENSHARE_ACTIVATED' : 'MY_SCREENSHARE_AVAILABLE',
                                }
                            })
                            , "*");
                    }
                    break;
                case 'UNHIDE_VIDEO':
                    if (camImg.dataset.status === 'off') {
                        camImg.dataset.status = 'on';
                        camImg.src = CAMERA_ON;
                        sparkRTC.disableVideo(true);
                        parent.postMessage(
                            JSON.stringify({
                                type: 'MY_VIDEO_UNHIDDEN', data: {
                                    isMyAudioMuted: document.getElementById("mic").dataset.status === 'off',
                                    isMyVideoHidden: document.getElementById("camera").dataset.status === 'off',
                                    screenShareState: document.getElementById("share_screen").dataset.status === 'on' ? 'MY_SCREENSHARE_ACTIVATED' : 'MY_SCREENSHARE_AVAILABLE',
                                }
                            })
                            , "*");
                    }
                    break;
                case 'SHARE_SCREEN':
                    shareScreenStream = await sparkRTC.startShareScreen();
                    if (shareScreenStream) {
                        scImg.dataset.status = 'on';
                        scImg.src = SCREEN_ON;
                        let video;
                        if (document.getElementById('screen-share')) {
                            video = document.getElementById('screen-share');
                        } else {
                            video = createVideoElement('screen-share');
                        }
                        video.srcObject = shareScreenStream;
                        video.play();
                        parent.postMessage(
                            JSON.stringify({
                                type: 'MY_SCREENSHARE_ACTIVATED', data: {
                                    isMyAudioMuted: document.getElementById("mic").dataset.status === 'off',
                                    isMyVideoHidden: document.getElementById("camera").dataset.status === 'off',
                                    screenShareState: document.getElementById("share_screen").dataset.status === 'on' ? 'MY_SCREENSHARE_ACTIVATED' : 'MY_SCREENSHARE_AVAILABLE',
                                }
                            })
                            , "*");
                    }
                    break;
                case 'STOP_SCREENSHARE':
                    scImg.dataset.status = 'off';
                    scImg.src = SCREEN_OFF;
                    if (document.getElementById('screen-share')) {
                        removeVideoElement('screen-share');
                    }
                    parent.postMessage(
                        JSON.stringify({
                            type: 'MY_SCREENSHARE_AVAILABLE', data: {
                                isMyAudioMuted: document.getElementById("mic").dataset.status === 'off',
                                isMyVideoHidden: document.getElementById("camera").dataset.status === 'off',
                                screenShareState: document.getElementById("share_screen").dataset.status === 'on' ? 'MY_SCREENSHARE_ACTIVATED' : 'MY_SCREENSHARE_AVAILABLE',
                            }
                        })
                        , "*");
                    break;
            }
        } catch (e) {
        }
        // callStatusUpdate();
    });
}


async function onRaiseHand() {
    const stream = await sparkRTC.raiseHand();
    const tagId = 'localVideo-' + stream.id;
    if (document.getElementById(tagId)) return;
    const video = createVideoElement(tagId, '', true);
    video.srcObject = stream;
    callStatusUpdate();
}

function callStatusUpdate() {
    parent.postMessage(
        JSON.stringify({
            type: 'CALL_STATUS_UPDATED', data: {
                isMyAudioMuted: document.getElementById("mic").dataset.status === 'off',
                isMyVideoHidden: document.getElementById("camera").dataset.status === 'off',
                screenShareState: document.getElementById("share_screen").dataset.status === 'on' ? 'MY_SCREENSHARE_ACTIVATED' : 'MY_SCREENSHARE_AVAILABLE',
            }
        })
        , "*");
}