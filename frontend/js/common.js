function getWsUrl() {
    // const baseUrl = window.location.href.split("//")[1].split("/")[0];
    // const protocol = window.location.href.split("//")[0] === "http:" ? "ws" :"wss";
    // return `${protocol}://${baseUrl}/ws`
    return 'ws://localhost:8080/ws'
}

function createSparkRTC() {
    if (myRole === 'broadcast'){
        return new SparkRTC('broadcast', (stream) => {
            console.log('got stream', stream);
            document.getElementById('localVideo').srcObject = stream;
        }, null);
    }
    return new SparkRTC('audience', null, (stream) => {
        if (document.getElementById('remoteVideo').srcObject !== stream)
            document.getElementById('remoteVideo').srcObject = stream;
    });
}
