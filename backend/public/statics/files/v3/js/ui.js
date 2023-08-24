const CAMERA_ON = "images/cam-on.png";
const CAMERA_OFF = "images/cam-off.png";
const MIC_ON = "images/mic-on.png";
const MIC_OFF = "images/mic-off.png";
const SCREEN_ON = "images/screen-on.png";
const SCREEN_OFF = "images/screen-off.png";
// const SPARK_LOGO = "images/spark-logo.png";
const CHARACTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
const verySlowColor = 'invert(64%) sepia(66%) saturate(4174%) hue-rotate(334deg) brightness(100%) contrast(92%)';
const DCColor = 'invert(13%) sepia(99%) saturate(4967%) hue-rotate(350deg) brightness(92%) contrast(96%)';

let graph;
let sparkRTC;
let myName;
let myRole;
let shareScreenStream;
let roomName;

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
}

function createVideoElement(videoId, muted = false) {
    let container = document.createElement('div');
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
    return video ? video : createVideoElement(videoId, true);
}

function removeVideoElement(videoId) {
    let video = document.getElementById(videoId);
    if (!video) return;
    let videoContainer = video.parentNode;
    if (!videoContainer) return;
    document.getElementById('screen').removeChild(videoContainer);
    arrangeVideoContainers();
}

function onNetworkIsSlow(downlink) {
    let msg = '';
    if (downlink > 0) {
        document.getElementById('net').style.filter = verySlowColor;
        document.getElementById('net').title = 'Network Status is Very Slow!';
        msg = 'You network speed is lower than normal, therefor you may experience some difficulties.';
    } else {
        document.getElementById('net').style.filter = DCColor;
        document.getElementById('net').title = 'Network Status is Disconnected!'
        msg = 'You are DISCONNECTED!';
    }
    document.getElementById('net').onclick = () => { alert(msg); };
    document.getElementById('net').style.display = '';
}

function onNetworkIsNormal() {
    document.getElementById('net').style.display = 'none';
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
}


function setMyName() {
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
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return urlParams.get('role') === 'broadcast' ? "broadcast" : "audience";
}

function getRoomName() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return urlParams.get('room');
}

function getDebug() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return Boolean(urlParams.get('debug'));
}

function setupSignalingSocket() {
    return sparkRTC.setupSignalingSocket(getWsUrl(), myName, roomName);
}


async function start() {
    await setupSignalingSocket();
    return sparkRTC.start();
}


function onLoad() {
    // registerNetworkEvent();
    myRole = getMyRole();
    roomName = getRoomName();
    sparkRTC = createSparkRTC();
    if (!getDebug()) {
        document.getElementById('logs').style.display = 'none';
    }

    setMyName();
    graph = new Graph();
    window.onresize = handleResize;
    graph.draw(DATA);

    arrangeVideoContainers();
}


async function onRaiseHand() {
    const stream = await sparkRTC.raiseHand();
    const tagId = 'localVideo-' + stream.id;
    if (document.getElementById(tagId)) return;
    const video = createVideoElement(tagId, true);
    video.srcObject = stream;
}

function addLog(log) {
    const logs = document.getElementById('logs');
    const p = document.createElement('p');
    p.innerText = log;
    logs.appendChild(p);
}