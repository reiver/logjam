import myPeerConnectionConfig from "../config/myPeerConnectionConfig";

export default function _usePeerConnection() {

    let socket,
        altStreams,
        myRole,
        myUsername,
        targetName,
        myPeerConnection,
        myPeerConnections = [],
        myAltPeerConnections = [];

    function createPeerConnection(targetUsername, view = "audience") {
        if (view === true) {
            console.log('createPeerConnection[alt]', targetUsername);
            if (myAltPeerConnections[targetUsername]) {
                console.log('Exists', targetUsername)
                return myAltPeerConnections[targetUsername];
            }
            console.log('Not Exists', targetUsername)
            const myAltPeerConnection = new RTCPeerConnection(myPeerConnectionConfig);
            myAltPeerConnection.onicecandidate = function (event) {
                if (event.candidate) {
                    socket.send(
                        JSON.stringify({
                            type: "alt-new-ice-candidate",
                            target: targetUsername,
                            candidate: event.candidate,
                        })
                    );
                }
            };
            myAltPeerConnection.ontrack = function (event) {
                console.log("TRACK", event, myRole);
                altStreams[targetUsername] = event.streams;
                // if (!document.getElementById("alt_video").srcObject) {
                //   document.getElementById("alt_video").srcObject =
                //     event.streams[0];
                //   remoteStream = event.streams[0];
                //   myPeerConnections.map(amyPeerConnection => {
                //     remoteStream.getTracks().forEach((track) => {
                //       amyPeerConnection.addTrack(track, remoteStream);
                //     });
                //   });
                // }
                if (document.getElementById("alt_video_" + targetUsername).length === 0) {
                    document.getElementById("streams").append(`<video id="alt_video_${targetUsername}" autoplay playsinline style="width: 100%"></video><br/>`);
                    document.getElementById("alt_video_" + targetUsername).srcObject = event.streams[0];
                    myPeerConnections.map(amyPeerConnection => {
                        event.streams[0].getTracks().forEach((track) => {
                            track['mehrdad'] = 12;
                            event.streams[0]['mehrdad'] = 12;
                            amyPeerConnection.addTrack(track, event.streams[0]);
                        });
                    });
                }
            };
            myAltPeerConnection.onnegotiationneeded = function (event) {
                myAltPeerConnection
                    .createOffer()
                    .then(function (offer) {
                        return myAltPeerConnection.setLocalDescription(offer);
                    })
                    .then(function () {
                        socket.send(
                            JSON.stringify({
                                name: myUsername,
                                target: targetUsername,
                                type: "alt-video-offer",
                                sdp: myAltPeerConnection.localDescription,
                            })
                        );
                    })
                    .catch(function (e) {
                        console.log('[myAltPeerConnection] createOffer error', e);
                    });
            };
            myAltPeerConnection.onremovetrack = function (event) {
            };
            myAltPeerConnection.oniceconnectionstatechange = function (event) {
                socket.send(
                    JSON.stringify({
                        type: "log",
                        data: "oniceconnectionstatechange :" + myAltPeerConnection.iceConnectionState,
                    })
                );
                if (myAltPeerConnection.iceConnectionState === 'disconnected') {
                    console.log('Alt Disconnected', myAltPeerConnection);
                }
            };
            myAltPeerConnection.onicegatheringstatechange = function (event) {
            };
            myAltPeerConnection.onsignalingstatechange = function (event) {
            };
            myAltPeerConnection.onicecandidateerror = (event) => {
                socket.send(
                    JSON.stringify({
                        type: "log",
                        data: "onicecandidateerror :" + JSON.stringify(event),
                    })
                );
            };
            myAltPeerConnections[targetUsername] = myAltPeerConnection;
            return myAltPeerConnection;
        }
        if (myPeerConnections[targetUsername]) {
            return myPeerConnections[targetUsername];
        }
        myPeerConnection = new RTCPeerConnection(myPeerConnectionConfig);
        myPeerConnection.view = view;
        myPeerConnection.onicecandidate = function (event) {
            if (event.candidate) {
                socket.send(
                    JSON.stringify({
                        type: "new-ice-candidate",
                        target: targetUsername,
                        candidate: event.candidate,
                    })
                );
            }
        };
        myPeerConnection.ontrack = function (event) {
            console.log("TRACK", event, myRole);
            if (myRole !== "broadcast") {
                // if(event.transceiver.sender.track.mmcomp && event.transceiver.sender.track.mmcomp=="SHARE")
                if (!document.getElementById("local_video").srcObject) {
                    document.getElementById("local_video").srcObject =
                        event.streams[0];
                    let localStream = event.streams[0];
                    document.getElementById("#source-stream").text("Read " + targetName); //+ "[" + targetUsername + "]");
                    socket.send(
                        JSON.stringify({
                            type: "stream",
                            data: "true",
                        })
                    );
                    socket.send(
                        JSON.stringify({
                            type: "log",
                            data: "ontrack from " + targetName,
                        })
                    );
                }
            } else {
                // if (!document.getElementById("alt_video").srcObject) {
                //   document.getElementById("alt_video").srcObject =
                //     event.streams[0];
                //   remoteStream = event.streams[0];
                // }
            }
            // document.getElementById("received_video_" + targetUsername).srcObject = event.streams[0];
        };
        myPeerConnection.onnegotiationneeded = function (event) {
            // console.log("onnegotiationneeded", veiw, event);
            //   if (veiw === "audience") {
            myPeerConnection
                .createOffer()
                .then(function (offer) {
                    return myPeerConnection.setLocalDescription(offer);
                })
                .then(function () {
                    socket.send(
                        JSON.stringify({
                            name: myUsername,
                            target: targetUsername,
                            type: "video-offer",
                            sdp: myPeerConnection.localDescription,
                        })
                    );
                })
                .catch(function (e) {
                    // console.log("onnegotiationneeded audience error:", e);
                    socket.send(
                        JSON.stringify({
                            type: "log",
                            data: "onnegotiationneeded audience error:" + JSON.stringify(e),
                        })
                    );
                });
            //   } else if (veiw === "replicate") {
            //     myPeerConnection
            //       .createOffer()
            //       .then(function (offer) {
            //         return myPeerConnection.setLocalDescription(offer);
            //       })
            //       .then(function () {
            //         socket.send(
            //           JSON.stringify({
            //             name: myUsername,
            //             target: targetUsername,
            //             type: "replicate-offer",
            //             sdp: myPeerConnection.localDescription,
            //           })
            //         );
            //       })
            //       .catch(function (e) {
            //         console.log("onnegotiationneeded replicate error:", e);
            //       });
            //   }
        };
        myPeerConnection.onremovetrack = function (event) {
        };
        myPeerConnection.oniceconnectionstatechange = function (event) {
            socket.send(
                JSON.stringify({
                    type: "log",
                    data: "oniceconnectionstatechange :" + myPeerConnection.iceConnectionState,
                })
            );
            // console.log("ICE Change : ", event);
            if (myPeerConnection.iceConnectionState === 'disconnected') {
                console.log('Disconnected', myPeerConnections);
                //              setTimeout(() => {
                //                login();
                //              }, 1000);
                // $("#audience-btn").prop("disabled", false);
                // $("#_broadcast-btn").prop("disabled", false);
                // if ($("#audience-btn").prop("disabled")) {
                // setTimeout(() => {
                //   audience();
                // }, 1000);
                // }
            }
        };
        myPeerConnection.onicegatheringstatechange = function (event) {
        };
        myPeerConnection.onsignalingstatechange = function (event) {
        };
        myPeerConnection.onicecandidateerror = (event) => {
            // console.log("Ice Error", event);
            socket.send(
                JSON.stringify({
                    type: "log",
                    data: "onicecandidateerror :" + JSON.stringify(event),
                })
            );
            //socket.close();
        };
        myPeerConnections[targetUsername] = myPeerConnection;
        return myPeerConnection;
    }

}
