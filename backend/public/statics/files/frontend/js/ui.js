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
    let flex = "0 0 calc(" + flexRatio + "% - " + flexGap + "px)";

    Array.from(videoContainers).forEach(div => {
            div.style.setProperty('flex', flex);
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

async function onShareScreen() {
    const img = document.getElementById("share_screen");
    if (!shareScreenStream) {
        img.dataset.status = 'on';
        img.src = SCREEN_ON;
        shareScreenStream = await sparkRTC.startShareScreen();
        document.getElementById('localScreen').srcObject = shareScreenStream;
    } else {
        img.dataset.status = 'off';
        img.src = SCREEN_OFF;
        shareScreenStream.getTracks().forEach(track => track.stop());
        shareScreenStream = null;
        document.getElementById('localScreen').srcObject = null;
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

function handleClick() {
    let newName = document.getElementById("inputName").value;
    if (newName) {
        myName = newName;
        localStorage.setItem('logjam_myName', myName);
        console.log('myName=', myName);
    }
    document.getElementById("page").style.visibility = "visible";
    document.getElementById("getName").style.display = "none";
    document.getElementById("myName").innerText = newName;

    return false;
}

function handleResize() {
    clearTimeout(window.resizedFinished);
    window.resizedFinished = setTimeout(function () {
        graph.draw(graph.treeData);
    }, 250);

}

function getMyRole() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return urlParams.get('role') === 'broadcast' ? "broadcast" : "audience";
}

function setupSignalingSocket() {
    return sparkRTC.setupSignalingSocket(getWsUrl(), myName);
}

function start(){
    return sparkRTC.start();
}

function onLoad() {
    myRole = getMyRole();
    console.log('myRole=', myRole);
    sparkRTC = createSparkRTC();

    setMyName();
    graph = new Graph();
    window.onresize = handleResize;
    graph.draw(DATA);

    arrangeVideoContainers();
}