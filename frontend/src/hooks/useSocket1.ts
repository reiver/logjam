import React, {useState} from 'react';

export {}

// // import myPeerConnectionConfig from "./myPeerConnectionConfig";
// // import {SocketMessage} from "../types/SocketMessage";
//
// import {useEffect, useState} from "react";
//
// export default function useSocket() {
//     let socket: WebSocket;
//     let localStream: any;
//     let lastTree: any;
//     let loggedIn = false;
//     let myUsername: any;
//     let iceCandidates: any = [];
//     let altIceCandidates: any = [];
//     let myRole = "broadcast";
//     let altStreams = {};
//     let targetName: any;
//
//     let myPeerConnection: any;
//     let myPeerConnections: any = [];
//     let myAltPeerConnections: any = [];
//

//
//

//
//     function start() {
//         // console.log("START!!");
//         if (message.error) {
//             alert(message.error);
//             return;
//         }
//         loggedIn = true;
//         myUsername = msg.data
//         $("#user-name").text("You are " + myName);// + "[" + myUsername + "]");
//         if ($("#audience-btn").prop("disabled")) {
//             // console.log("reconnecting as audience");
//             audience();
//         } else if ($("#broadcast-btn").prop("disabled")) {
//             // console.log("reconnecting as broadcast");
//             broadcast();
//         }
//     }
//
//     function audience() {
//         socket.send(
//             JSON.stringify({
//                 type: "role",
//                 data: "audience",
//             })
//         );
//         const turnStatus = myPeerConnectionConfig.iceServers.findIndex(s => s.url.indexOf('turn:') === 0) >= 0;
//         socket.send(
//             JSON.stringify({
//                 type: "turn_status",
//                 data: turnStatus ? "on" : "off",
//             })
//         );
//         $("#audience-btn").prop("disabled", true);
//         $("#broadcast-btn").prop("disabled", false);
//         document.getElementById("local_video").srcObject = null;
//     }
//
//     function newIceCandidate() {
//         if (!message) return;
//
//         // console.log("ICE Canadidate Recieved!");
//         iceCandidates.push(new RTCIceCandidate(message.candidate));
//         // console.log(iceCandidates);
//         if (myPeerConnection && myPeerConnection.remoteDescription) {
//             // console.log("Adding canadidate");
//             let candidate = iceCandidates.pop();
//             myPeerConnection.addIceCandidate(candidate).catch(function (e: any) {
//                 // console.log("new-ice-candidate error:", e, candidate);
//             });
//         }
//     }
//
//     function newAltIceCandidate() {
//         if (!message) return;
//         let targetUsername = message.data;
//         let myPeerConnection = createPeerConnection(targetUsername, true);
//         // console.log("ICE Canadidate Recieved!");
//         altIceCandidates.push(new RTCIceCandidate(message.candidate));
//         // console.log(iceCandidates);
//         if (myPeerConnection && myPeerConnection.remoteDescription) {
//             // console.log("Adding canadidate");
//             let candidate = altIceCandidates.pop();
//             myPeerConnection.addIceCandidate(candidate).catch(function (e: any) {
//                 // console.log("new-ice-candidate error:", e, candidate);
//             });
//         }
//     }
//
//     function altBroadcast() {
//         if (!message) return;
//
//         if (confirm(`${message.name} wants to broadcast, do you approve?`)) {
//             socket.send(
//                 JSON.stringify({
//                     type: "alt-broadcast-approve",
//                     target: message.data,
//                 })
//             );
//         }
//     }
//
//     async function altVideoOffer() {
//         if (!message) return;
//
//         // console.log("Video Offer Rec!");
//         let targetUsername = message.name;
//         targetName = message.username;
//         let desc = new RTCSessionDescription(message.sdp);
//         let myPeerConnection = createPeerConnection(targetUsername, true);
//
//         await myPeerConnection.setRemoteDescription(desc);
//
//         let answer = await myPeerConnection.createAnswer();
//
//         await myPeerConnection.setLocalDescription(answer);
//
//         socket.send(
//             JSON.stringify({
//                 name: myUsername,
//                 target: targetUsername,
//                 type: "alt-video-answer",
//                 sdp: myPeerConnection.localDescription,
//             })
//         );
//
//         myPeerConnection.view = 'alt';
//     }
//
//     function videoAnswer() {
//         if (!message) return;
//         let desc = new RTCSessionDescription(message.sdp);
//         myPeerConnection.setRemoteDescription(desc).catch((e) => {
//             console.log("Error", e);
//         });
//     }
//
//     function newAltVideoAnswer() {
//         if (!message) return;
//         let desc = new RTCSessionDescription(message.sdp);
//         let targetUsername = message.target;
//         let myPeerConnection = createPeerConnection(targetUsername, true);
//         myPeerConnection.setRemoteDescription(desc).catch((e) => {
//             console.log("Error", e);
//         });
//     }
//
//     function connectUser() {
//         if (!message) return;
//         let targetUsername = message.data;
//         // console.log("Connect User", targetUsername);
//         if (myPeerConnections[targetUsername]) {
//             // console.log("Already Connected");
//             return true;
//         }
//         let myPeerConnection = createPeerConnection(targetUsername);
//
//         // console.log("localStream", localStream);
//         if (localStream) {
//             localStream.getTracks().forEach((track) => {
//                 // track["mmcomp"] = "VIDEO";
//                 myPeerConnection.addTrack(track, localStream);
//             });
//             // console.log("Added tracks to connection");
//         }
//         for (const tuser in altStreams) {
//             const tstreams = altStreams[tuser];
//             tstreams[0].getTracks().forEach((track) => {
//                 myPeerConnection.addTrack(track, tstreams[0]);
//             });
//         }
//         // if (shareScreenStream) {
//         //   shareScreenStream.getTracks().forEach((track) => {
//         //     // track["mmcomp"] = "VIDEO";
//         //     myPeerConnection.addTrack(track, shareScreenStream);
//         //   });
//         //   // console.log("Added tracks to connection");
//         // }
//         // if (remoteStream) {
//         //   remoteStream.getTracks().forEach((track) => {
//         //     // track["mmcomp"] = "VIDEO";
//         //     myPeerConnection.addTrack(track, remoteStream);
//         //   });
//         //   // console.log("Added tracks to connection");
//         // }
//     }
//
//     function tree() {
//         if (!message) return;
//         const treeData = JSON.parse(message.data);
//         // console.log({treeData});
//         if (lastTree && lastTree === message.data) {
//             return;
//         }
//         lastTree = msg.data;
//         if (treeData && treeData.length > 0) {
//             drawGraph(treeData);
//         }
//     }
//
//     function createSocket() {
//         socket = new WebSocket(getSocketUrl());
//

// }
