const CAMERA_ON = "images/cam-on.png";
const CAMERA_OFF = "images/cam-off.png";
const MIC_ON = "images/mic-on.png";
const MIC_OFF = "images/mic-off.png";
const SPARK_LOGO = "images/spark-logo.png";

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
    }else{
        img.dataset.status = 'on';
        img.src = CAMERA_ON;
    }
}


function onMicButtonClick() {
    const img = document.getElementById("mic");
    if (img.dataset.status === 'on') {
        img.dataset.status = 'off';
        img.src = MIC_OFF;
    }else{
        img.dataset.status = 'on';
        img.src = MIC_ON;
    }
}

function onLoad() {
    arrangeVideoContainers();
    let graph = new GraphVisualization();
    graph.draw(DATA);
}