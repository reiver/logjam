/**
 *  This file is main JS file to intract with index.html file
 * 
 *  It contains fucntions to respond when user intratct with UI elements.
 * 
 *  It contains fucntions to initate the Broadcaster or Audiance
 */

const CAMERA_ON = "images/cam-on.png";
const CAMERA_OFF = "images/cam-off.png";
const MIC_ON = "images/mic-on.png";
const MIC_OFF = "images/mic-off.png";
const SCREEN_ON = "images/screen-on.png";
const SCREEN_OFF = "images/screen-off.png";
// const SPARK_LOGO = "images/spark-logo.png";
const RAISE_HAND_ON = "images/hand.png";
const RAISE_HAND_OFF = "images/hand-off.png";
const CHARACTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
const verySlowColor = 'invert(64%) sepia(66%) saturate(4174%) hue-rotate(334deg) brightness(100%) contrast(92%)';
const DCColor = 'invert(13%) sepia(99%) saturate(4967%) hue-rotate(350deg) brightness(92%) contrast(96%)';

let graph;
let sparkRTC;
let myName;
let myRole;
let shareScreenStream;
let roomName;

/**
 * Description
 * Function to generate Random UUID of any length
 * @param {any} length any number
 * @returns {any}
 */
function makeId(length) {
    let result = '';
    for (let i = 0; i < length; i++) {
        result += CHARACTERS.charAt(Math.floor(Math.random() * CHARACTERS.length));
    }
    return result;
}


/** 
 * Function to arrange multiple video containers on screen
 * by adjusting their height accordingly 
 * 
 */
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


/**
 * Description
 * Function to perfom take appropriate actions upon camera button click,
 * Enable or Disable Local Video Track
 * @returns {any}
 */
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


/**
 * Description
 * Funnction to take certain actions on mic button click,
 * Mute or Unmute the Mic
 * @returns {any}
 */
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

/** Function to create new Vido element
 * 
 * To display, Video stream (local or remote)
*/
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

/**
 * Function to get ref to Video element on Screen
 * @param {any} videoId
 * @returns {any}
 */
function getVideoElement(videoId) {
    let video = document.getElementById(videoId);
    return video ? video : createVideoElement(videoId, true);
}

/**
 * Description
 * Function to remove the video element from the screen by using it's ID
 * @param {any} videoId
 * @returns {any}
 */
function removeVideoElement(videoId) {
    let video = document.getElementById(videoId);
    if (!video) return;
    let videoContainer = video.parentNode;
    if (!videoContainer) return;
    document.getElementById('screen').removeChild(videoContainer);
    arrangeVideoContainers();
}

/**
 * Description
 * Function to display Slow network status
 * @param {any} downlink
 * @returns {any} Nothing
 */
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

/**
 * Description Function to display normal network status
 * @returns {any}
 */
function onNetworkIsNormal() {
    document.getElementById('net').style.display = 'none';
}

/**
 * Description
 * Function to start or stop screen share and update UI accordingly
 * @returns {any}
 */
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


/**
 * Description
 * Function to save current user name to local storage, if not name provided generate random name
 * @returns {any}
 */
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


/**
 * 
 * Function to handle Click on Enter Button on Mian Page
 * 
 * @param {any} turn=true
 * @returns {any}
 */
async function handleClick(turn = true) {
    let newName = document.getElementById("inputName").value;
    if (newName) {
        myName = newName;
        localStorage.setItem('logjam_myName', myName);
    }
    document.getElementById("page").style.visibility = "visible";
    document.getElementById("getName").style.display = "none";

    await start(turn);

    return false;
}


function handleResize() {
    clearTimeout(window.resizedFinished);
    window.resizedFinished = setTimeout(function () {
        graph.draw(graph.treeData);
        arrangeVideoContainers();
    }, 250);

}


/**
 * 
 * Function to get User role from URL
 * 
 * @returns {any}
 */
function getMyRole() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return urlParams.get('role') === 'broadcast' ? "broadcast" : "audience";
}

/**
 * 
 * Function to get Room name from URL
 * @returns {any}
 */
function getRoomName() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return urlParams.get('room');
}

/**
 * 
 * Function to get debug status from URL
 * @returns {any}
 */
function getDebug() {
    const queryString = window.location.search;
    const urlParams = new URLSearchParams(queryString);
    return Boolean(urlParams.get('debug'));
}

/**
 * 
 * Function to setUp Signaling socket
 * @returns {any}
 */
function setupSignalingSocket() {
    return sparkRTC.setupSignalingSocket(getWsUrl(), myName, roomName);
}


/**
 * 
 * Function to setup signaling socket and start broadcasting or 
 * start listening to broadcast if role is Audiance.
 * @param {any} turn=true
 * @returns {any}
 */
async function start(turn = true) {
    await setupSignalingSocket();
    return sparkRTC.start(turn);
}


/**
 *
 * The very first fucntion to call when page Loads
 * 
 * @returns {any}
 */
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
    //console.log("DATA: ",DATA);
    //graph.draw(DATA);

    arrangeVideoContainers();
}


/**
 * 
 * Function to handle click on handRaise button
 * 
 * @returns {any}
 */
async function onRaiseHand() {
    const img = document.getElementById("raise_hand");
    console.log(img.dataset.status);
    if (img.dataset.status === 'on') {
        if (sparkRTC.localStream) {
            // if (confirm(`Do you want to stop streaming?`)) {
            //     console.log('stopping...');
            //     removeVideoElement('localVideo-' + sparkRTC.localStream.id);
            //     disableAudioVideoControls();
            //     sparkRTC.lowerHand();
            // }
            return;
        }
        const stream = await sparkRTC.raiseHand();
        // const tagId = 'localVideo-' + stream.id;
        // const video = createVideoElement(tagId, true);
        // video.srcObject = stream;
        document.getElementById('mic').style.display = '';
        document.getElementById('camera').style.display = '';
    }
}

/**
 * 
 * Function to add log to LOGS list
 * 
 * @param {any} log it could be any string as a log
 * @returns {any}
 */
function addLog(log) {
    const logs = document.getElementById('logs');
    const p = document.createElement('p');
    p.innerText = log;
    logs.appendChild(p);
}

/**
 * 
 * Function to display Audio Video Controls on Screen
 * 
 * @returns {any}
 */
function enableAudioVideoControls() {
    document.getElementById('mic').style.display = '';
    document.getElementById('camera').style.display = '';
    document.getElementById('share_screen').style.display = '';
}

/**
 * 
 * Function to hide Audio Video controls
 * 
 * @returns {any}
 */
function disableAudioVideoControls() {
    document.getElementById('mic').style.display = 'none';
    document.getElementById('camera').style.display = 'none';    
    document.getElementById('share_screen').style.display = 'none';
}
