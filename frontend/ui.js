function arrangeVideoContainers() {
    const videoContainers = document.getElementById('screen').getElementsByClassName('video-container');
    const videoCount = videoContainers.length;
    const flexGap = 1;
    let flexRatio = 100 / Math.ceil(Math.sqrt(videoCount));
    let flex = "0 0 calc(" + flexRatio + "% - " + flexGap + "px)";
    videoContainers.forEach(div => {
            div.style.setProperty('flex', flex);
        }
    )
}


function onLoad(){
    arrangeVideoContainers();
}