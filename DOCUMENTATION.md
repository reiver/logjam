
# **LOGJAM**

### Code Documentation 

In this file I will be explaining the code of **logjam**, I will try to use simple language and explain the usage, working & why code is writen the way its written.

<!-- Starting with UI.JS File from Here -->
.  

.  

# Let's start with `ui.js` File

`ui.js` is the main file interacting the `UI Layer`, infact it is the part of `UI Layer`. let's look into each of it's fucntion in depth.

---


<!-- f1 -->
## `makeId` Function

### Signature
`makeId(length: number): string`


### Parameters
- `length`: A `number` value representing the desired length of the returned string.

### Return Value
- The function returns a `string` value representing a random combination of characters with the length specified in the `length` parameter.

### Description
The `makeId` function takes in a single `length` parameter, which is used to determine the length of the returned string. 

- The function uses a `for` loop to generate a random combination of characters by selecting a random character from a predefined set of characters stored in the `CHARACTERS` constant and appending it to the `result` string. 
- The loop continues until the length of `result` is equal to the `length` specified in the `length` parameter. 
- Finally, the `result` string is returned as the result of the `makeId` function.

It should be noted that the `CHARACTERS` constant is defined in the code as a Global Variable, and its value would need to be specified in order for the `makeId` function to work as intended.

 
---
<!-- f2 -->
## `arrangeVideoContainers` Function

### Signature
`arrangeVideoContainers(): void`


### Parameters
None

### Return Value
None

### Description
The `arrangeVideoContainers` function is used to arrange multiple video containers on the screen by adjusting their height accordingly. 

- The function starts by using the `document.getElementById` and `getElementsByClassName` methods to retrieve all elements with the class `video-container` within an element with the id `screen`.
- The function then calculates the total number of video containers and uses this value to calculate the `flexRatio` and `maxHeight` values, which are used to set the `flex` and `max-height` CSS properties for each video container.
- Finally, the `arrangeVideoContainers` function uses the `Array.from` method to iterate over the video containers, setting the `flex` and `max-height` properties for each container using the `style.setProperty` method.

---
<!-- f3 -->
## `onCameraButtonClick` Function

### Signature
`function onCameraButtonClick()`


### Parameters
None

### Return Value
None

### Description
This function takes action on camera button click.

- It changes the status of the camera, based on the status of the camera the source of the image is changed to either `CAMERA_OFF` or `CAMERA_ON`. 
- It also calls `sparkRTC.disableVideo` to disable the video stream.

---
<!-- f4 -->
## `onMicButtonClick` Function

### Signature
`function onMicButtonClick()`

### Parameters
None

### Return Value
None

### Description
This function takes certain actions on mic button click.

- It mute or unmutes the mic based on its status. 
- It changes the status of the mic and updates the source of the image to either `MIC_OFF` or `MIC_ON`. 
- It also calls `sparkRTC.disableAudio` to disable the audio stream.

---
<!-- f5 -->
## `createVideoElement` Function
  
### Signature
`function createVideoElement(videoId, muted = false)`

### Parameters
- `videoId`: Id of the video element to be created.
- `muted`: A boolean value indicating whether the video should be muted or not.

### Return Value
- It returns The created video element.

### Description
This function creates a new video element to display a video stream (local or remote). 

- It creates a container div with class `video-container` and appends a `video` element to it. 
- It sets the properties of the video element, such as id, autoplay, playsInline, and mute status. 
- Finally, it appends the container to the element with id `screen` and returns the created video element.

---
<!-- f6 -->
## `getVideoElement(videoId)` Function

### Signature
`function getVideoElement(videoId)`


### Parameters
- `videoId: integer` - The id of the video element to be retrieved

### Return value
- Returns a `video` element with the given ID, if not found `creates` a new video element and returns that.

### Description
This function is used to retrieve a video element from the screen using its ID. 

- If a video element with the given ID is not found, the function creates a new video element with the given ID and returns it.

---
<!-- f7 -->
## `removeVideoElement(videoId)` Function

### Signature
`function removeVideoElement(videoId)`


### Parameters
- `videoId: Integer` - The id of the video element to be removed

### Return value
None

### Description
This function is used to `remove` the video element from the screen using its ID.

---
<!-- f8 -->
## `onNetworkIsSlow(downlink)` Function

### Signature
`function onNetworkIsSlow(downlink)`


### Parameters
- `downlink: Inetger` - The downlink value indicating the network speed

### Return value
None

### Description
This function is used to `display` the network status as `slow`. 

- If the downlink value is greater than `0`.
- It sets the status as very `slow` and displays an `alert` message indicating the network status is very slow. 
- If the downlink value is `0`, it sets the status as `disconnected` and displays an `alert` message indicating the network is disconnected.

---
<!-- f9 -->
## `onNetworkIsNormal` Function 

### Signature 
`function onNetworkIsNormal()`

### Parameters
None

### Return Value
None

### Description:
This function is used to hide the element with ID `net` from the HTML document. 

- The display style of the element is set to `none`. 
- This means that the element will not be `visible` on the page. 
- The purpose of this function is to `remove` the notification of a network issue when the network is back to normal.

---
<!-- f10 -->
## `onShareScreen` Function

### Signature
`async function onShareScreen()`

### Parameters
None

### Return Value 
None

### Description
This function is used to `enable` or `disable` screen sharing in a real-time communication (RTC) application using WebRTC API. 

- The function uses the `sparkRTC` object to start or stop the screen sharing. 
- When the function is called for the first time, it starts the screen sharing by calling the `sparkRTC.startShareScreen()` method. 
- The method returns a stream of the shared screen, which is assigned to the `shareScreenStream` variable. 
- The function then updates the source of the local screen video element with the `new stream`. 

When the function is called again.

- It stops the screen sharing by stopping all the tracks of the `shareScreenStream` and setting it to `null`.
- The local screen video element is also updated to have a `null` source. 
- The `removeVideoElement('localScreen')` method is used to remove the video element from the page.

The function also updates the status of the **share screen** button by setting the `dataset.status` attribute and changing the source of the image with the `SCREEN_ON` or `SCREEN_OFF` image path.

---
<!-- f11 -->
## `setMyName` Function

### Signature
`function setMyName()`

### Parameters
None

### Return Value
None

### Description:
This function is used to set the user's name in a Video application. 

- The function first tries to retrieve the name from the **local storage** using the key `logjam_myName`. 
- The retrieved name is then set as the value of the input field with ID `inputName`.

If the name cannot be retrieved from the **local storage**.

- the function generates a new name using the `makeId(20)` function and sets it as the value of the `myName` variable. 
- The generated name is then stored in the local storage using the key `logjam_myName`.

The purpose of the function is to provide a unique name for each user in the Video application. 

- The name is stored in the **local storage** so that the user can access it across multiple sessions.

--- 
<!-- f12 -->
## `handleClick` Function

### Signature
`async function handleClick(turn = true)`

### Parameters
- `turn: boolean` (optional, default is `true`)

### Return Value
- `False`

### Description
This function is used to handle the click event of the submit button in a Video application. 

- The function first retrieves the name entered in the input field with ID `inputName` and stores it in the local storage using the key `logjam_myName`. 
- The function then updates the visibility and display styles of the page and the `getName` section.

The function also calls the `start(turn)` method, which is used to start the Video session.

- The `turn` parameter is used to specify the type ice server, TURN or STUN.

Finally, the function returns `false`, which is used to prevent the form from submitting and reloading the page.

---
<!-- f13 -->
## `handleResize` Function

### Signature
`function handleResize()`

### Parameters
None

### Return Value
None

### Description
This function is used to handle the resize event of the window in a Video application. 

- The function first clears the previous timeout using the `clearTimeout(window.resizedFinished)` method.\
- The function then sets a new timeout using the `setTimeout` method. 
- The timeout function calls two methods, `graph.draw(graph.treeData)` and `arrangeVideoContainers()`, to redraw the graph and rearrange the video containers.

The purpose of the function is to update the layout of the video window whenever the window is resized, ensuring that the graph and video containers are displayed correctly. 

- The timeout is used to **prevent** the function from being called **too frequently** during rapid resizing.

---
<!-- f14 -->
## `getMyRole` Function

### Signature
`function getMyRole()`

### Parameters
None

### Return Value
- `String` (either `"broadcast"` or `"audience"`)

### Description
This function is used to retrieve the `role` of the user in a video application. 

- The function first retrieves the `query` string from the current URL using the `window.location.search` property.
- The function then creates a new `URLSearchParams` object using the query string, allowing it to access the URL parameters. 
- The function uses the `urlParams.get('role')` method to retrieve the value of the `role` parameter.
- The function returns either `"broadcast"` or `"audience"` based on the value of the `role` parameter. 
- If the parameter is not present or its value is not `"broadcast"`, the function returns `"audience"`.

The purpose of the function is to determine the user's role in the video session, either as a broadcaster or as an audience member.

---
<!-- f15 -->
## `getRoomName` Function

### Signature
`function getRoomName()`

### Parameters
None

### Return Value
- `String` (room name)

### Description
This function is used to retrieve the `name` of the `room` in a video application. 

- The function first retrieves the query string from the current URL using the `window.location.search` property.
- The function then creates a new `URLSearchParams` object using the query string, allowing it to access the URL parameters. 
- The function uses the `urlParams.get('room')` method to retrieve the value of the `room` parameter.
- The function returns the value of the `room` parameter, which is the name of the room.

The purpose of the function is to determine the name of the room the user is joining in the video session.

---
<!-- f16 -->
## `getDebug` Function

### Signature
`function getDebug()`

### Parameters
None

### Return Value
- `Boolean` (true or false)

### Description
This function is used to retrieve the value of the `debug` parameter in a video application. 

- The function first retrieves the query string from the current URL using the `window.location.search` property.
- The function then creates a new `URLSearchParams` object using the query string, allowing it to access the URL parameters. 
- The function uses the `urlParams.get('debug')` method to retrieve the value of the `debug` parameter.
- The function returns the value of the `debug` parameter converted to a `boolean` value, which indicates whether or not the debug mode is **enabled**.

The purpose of the function is to determine if the debug mode is **enabled** in the video session, allowing for the display of additional information for debugging purposes.

---
<!-- f17 -->
## `setupSignalingSocket` Function

### Signature
`function setupSignalingSocket()`

### Parameters
None

### Return Value
None

### Description
This function is used to setup the signaling socket in a video application. 

- The function uses the `sparkRTC.setupSignalingSocket()` method and passes the following parameters to it:
    - `getWsUrl()`: a function that returns the URL for the WebSockets server.
    - `myName`: a string representing the user's name.
    - `roomName`: a string representing the name of the video room.

The purpose of the function is to initiate the signaling socket connection and pass the necessary information to it, allowing the application to establish communication between the users.

---
<!-- f18 -->
## `start` Function

### Signature
`async function start(turn = true)`

### Parameters
- `turn` (optional): a boolean value indicating whether or not the user should act as a broadcaster or an audience. The default value is "true".

### Return Value
None

### Description
This function is used to start the broadcasting or listening to the broadcast in a chat application. 

- The function calls the `setupSignalingSocket()` function to initiate the signaling socket and then uses the `sparkRTC.start()` method and passes the "turn" parameter to it.

The purpose of the function is to start the communication between the users by setting up the signaling socket and determining whether the user should act as a `broadcaster` or an `audience`.

---
<!-- f19 -->
## `onLoad` Function

### Signature
`function onLoad()`

### Parameters
None

### Return Value
None

### Description
The function `onLoad` is used to execute some statements when the page loads. 

- It sets the role of the user and room name based on the URL. It creates a new instance of the SparkRTC class and sets the `myName` property. 
- The logs are hidden if the `debug` query string parameter is not present. 
- It also calls `handleResize` when the window is resized and initializes the `Graph` object.

---
<!-- f20 -->
## `onRaiseHand` Function

### Signature
`async function onRaiseHand()`

### Parameters
None

### Return Value
None

### Description
The function `onRaiseHand` is used to handle the click event on the `raise_hand` button. 

- If the `raise_hand` button status is `on`, the function raises the hand by calling `sparkRTC.raiseHand` and the microphone and camera icons are displayed.

---
<!-- f21 -->
## `addLog` Function

### Signature
`function addLog(log)`

### Parameters
- log: `string` A string that needs to be added to the logs

### Return Value
None

### Description
The `addLog` function takes in a `string` log as a parameter and adds it to the logs. 

- It does this by getting the logs element from the document using `document.getElementById('logs')`, creating a new paragraph element using `document.createElement('p')`, setting the innerText of the paragraph element to the log parameter, and finally, appending the paragraph element to the logs element.

---
<!-- f22 -->
## `enableAudioVideoControls` Function

### Signature
`function enableAudioVideoControls()`

### Parameters
None

### Return Value
None

### Description
The function `enableAudioVideoControls` sets the visibility of three different icons to display. 

- It sets the display property of the elements with ids `mic`, `camera`, and `share_screen` to an empty string. 
- This makes the elements visible on the page.

---
<!-- f23 -->
## `disableAudioVideoControls` Function

### Signature
`function disableAudioVideoControls()`

### Parameters
None

### Return Value
None

### Description
This function is used to `hide` the audio and video controls from the UI. 

- It gets the elements with the id `mic`, `camera`, and `share_screen` from the document and sets their display property to `none`. 

---

.

.  

# Next we will look at `SparkRTC.js` File

<!-- from here we will start SPARKRTC.JS -->

`SparkRTC.js` is the Main Class interacting with the `Signaling` Server via Web Sockets. It controls all the application logic. let's look into each of it's function in depth.

---


<!-- f1 -->
## `handleVideoOfferMsg` Function

### Signature
`handleVideoOfferMsg = async (msg) => {}`

### Parameters
- `msg`: an object containing data for the message

### Return Value
None

### Description
This function is called when the client receives a message of type `video-offer` or `alt-video-offer`. 

- This function will ***create*** a `peer connection` or ***retrieve*** an existing one for the message sender and set the `remote description` with the `sdp` provided in the message. 
- Then, it will create an `answer` for the offer and send it back to the message sender.

---
<!-- f2 -->
## `handleMessage` Function
### Signature
`handleMessage = async (event) => {}`

### Parameters
- `event`: an event object with the data of the message received on the web socket.

### Return Value
None

### Description
This function is the main socket handler for messages received on the web socket. 

- It parses the message and based on the message type, makes decisions on what to do next.
- The function first tries to parse the message into a JSON object, and if that fails it returns. 
- Then, it sets the `data` and `type` properties of the message to the right values, as they may be in different cases in the message.
- Next, the function has a switch statement to handle different types of messages, such as `video-offer`, `video-answer`, `new-ice-candidate`, `role`, `start`, `add_audience`, and `alt-broadcast-approve`. 
- The actions taken for each type of message may include creating or retrieving a peer connection, updating the status, connecting to an audience, sending a stream to an audience, and more.

---
<!-- f3 -->
## `broadcasterUserId` Function

### Signature
`broadcasterUserId : () => string`

### Parameters
None

### Return Value
- `string`

### Description
This function returns the `user id` of the broadcaster by iterating through the `myPeerConnectionArray` object and checking for the user with `isAudience` property set to `false`. 

- If no such user is found, it returns `null`.

---
  
## `ping` Function

### Signature
`ping(): void`

### Parameters
None

### Return Value
None

### Description
The `ping` function is an arrow function that sends a message over a socket connection using the `socket.send()` method. 

- The message is a JSON string that contains a `type` property whose value is either `"tree"` or `"ping"`, depending on whether the `treeCallback` property is truthy or falsy. 
- If `treeCallback` is truthy, the `type` property value will be `"tree"`, and if it is falsy, the `type` property value will be `"ping"`. 
- The `send()` method sends the message to the server through the socket connection.

This function is useful for testing the connection status of the socket and can be used to keep the connection alive by sending periodic pings to the server.

---

## `setupSignalingSocket` Function

### Signature
`setupSignalingSocket(url: string, myName: string, roomName: string): Promise<WebSocket>`

### Parameters
- `url`: A `string` value representing the URL to connect the WebSocket.
- `myName`: A `string` value representing the name of the client.
- `roomName`: A `string` value representing the name of the room the client wants to join.

### Return Value
- The function returns a `Promise` that resolves to a `WebSocket` object.

### Description
The `setupSignalingSocket` function is an arrow function that takes in three parameters - `url`, `myName`, and `roomName`. 

- It returns a `Promise` that resolves to a `WebSocket` object. 
- The function first logs the input parameters to the console. 
- If a `pingInterval` property exists, it is cleared. Then, if `myName` and `roomName` parameters are provided, they are assigned to the `myName` and `roomName` properties of the function. 
- The function then creates a new `WebSocket` object using the specified `url` and `roomName`.

The function sets up event handlers for the `onmessage`, `onopen`, `onclose`, and `onerror` events of the `WebSocket` object. 

- When the `onopen` event is triggered, the function sends a `"start"` message to the server, containing the `myName` parameter in its `data` property, using the `send()` method. 

    - It also sets up a `pingInterval` that calls the `ping` function every 5 seconds to keep the connection alive. 
    - Finally, it resolves the `Promise` with the `WebSocket` object.

- When the `onclose` event is triggered, the function sets various properties to their default values and calls the `signalingDisconnectedCallback` and `startProcedure` callbacks if they exist.

- When the `onerror` event is triggered, the function logs the error message to the console, rejects the `Promise`, and displays an alert to the user.

This function is useful for setting up a `WebSocket` connection to a signaling server, which can then be used to exchange information between clients in a WebRTC session.

---

## `startBroadcasting` Function

### Signature
`startBroadcasting(data: string): Promise<MediaStream | undefined>`

### Parameters
- `data`: A `string` value representing the role of the user.

### Return Value
- The function returns a `Promise` that resolves with a `MediaStream` object representing the user's local media stream.

### Description
The `startBroadcasting` function takes in a single optional `data` parameter, which is used to determine the role of the user. 

- If the user's `localStream` is not already set, the function attempts to retrieve the local stream using the `getUserMedia` method from the `navigator.mediaDevices` API. 
- If successful, the `localStream` is stored in the `this.localStream` property and added to the `this.remoteStreams` array. 
- The function then sends a message to the signaling server with the `data` and type `role`. Finally, the function returns a Promise that resolves with the `localStream`.
- If an error occurs while attempting to retrieve the local stream, the function updates the status with an error message and displays an alert message. 
- The function also logs the error message to the console and returns `undefined`.

---

## `startReadingBroadcast` Function

### Signature
`startReadingBroadcast(): void`

### Parameters
None

### Return Value
None

### Description
The `startReadingBroadcast` function sends a request to the server to start the broadcasting role for the client as an audience member. 

- It first updates the status to indicate that it is requesting the audience role. 
- Then, it sends a JSON object to the server using the `socket.send` method with the `type` property set to "role" and the `data` property set to "audience". 

This JSON object indicates to the server that the client is requesting the audience role. Finally, the function logs the success of the message sending operation to the console using the `log` method.

---

## `raiseHand` Function

### Signature
`raiseHand(): void | MediaStream`

### Parameters
None

### Return Value
- The function returns `void` if `startedRaiseHand` is true, otherwise, it returns a `MediaStream` by calling the `startBroadcasting` function with `'alt-broadcast'` as the `data` parameter.

### Description
The `raiseHand` function checks if the `startedRaiseHand` flag is set to `true`. 

- If it is, the function returns `void` and exits. Otherwise, it sets `startedRaiseHand` to `true` and calls the `startBroadcasting` function with `'alt-broadcast'` as the `data` parameter. 
- The `startBroadcasting` function is expected to return a `MediaStream` object which is then returned by the `raiseHand` function.

It uses the `startBroadcasting` function which I explaiend above.

---

## `onDataChannelOpened` Function

### Signature
`onDataChannelOpened(dc: RTCDataChannel, target: string, pc: RTCPeerConnection): void`

### Parameters
- `dc`: An `RTCDataChannel` object representing the data channel that has been opened.
- `target`: A `string` value representing the target of the data channel.
- `pc`: An `RTCPeerConnection` object representing the peer connection used to open the data channel.

### Description
The `onDataChannelOpened` function is called when a data channel is opened between two peers. 

- The function takes in three parameters, `dc`, `target`, and `pc`. 
- The `dc` parameter represents the data channel that has been opened, `target` represents the target of the data channel, and `pc` represents the peer connection used to open the data channel.
- The function uses a `setInterval` loop to check the `readyState` of the data channel. 
- If the `readyState` is `open`, the function sends a message through the data channel using the `dc.send` method. If the `readyState` is `connecting`, `closing`, or `closed`, the function logs the corresponding message to the console. 
- The loop continues until the `readyState` is `closed`, at which point the loop is exited using the `clearInterval` method.

---

## `restartEverything` Function

### Signature
`restartEverything(peerConnection: RTCPeerConnection, target: string): void`

### Parameters
- `peerConnection`: An `RTCPeerConnection` instance representing the connection with the remote peer.
- `target`: A `string` value representing the ID of the target peer.

### Return Value
None

### Description
The `restartEverything` function is used to restart the peer connection when a disconnection is detected. 

- The function takes in two parameters, `peerConnection` which represents the connection with the remote peer and `target` which represents the ID of the target peer. 

- The function first sets the `remoteStreamNotified` flag to `false`. 

    - If the length of the remote streams array is zero, the function immediately returns.

- The function then gets the IDs of all the tracks associated with the remote stream and iterates through them. 

    - For each track ID, it checks all the peer connections in the `myPeerConnectionArray` object, except the target peer connection. 
    - If a peer connection is an audience peer connection, it removes the sender associated with the track ID.

After that, the function removes all the remote streams from the `remoteStreams` array.

- If the `parentStreamId` is set and it is included in the remote streams array, the function updates the status with a message that the parent stream is disconnected. 

    - It then calls the `remoteStreamDCCallback` for each remote stream. 
    - The `parentStreamId` is then set to undefined, and the `remoteStreams` array is reset to an empty array.

- The function then tries to call the `remoteStreamDCCallback` function for the first remote stream of the `peerConnection` object.

- Finally, if the `parentDC` or `startedRaiseHand` flag is set, the `startProcedure` function is called.

---

## `checkParentDisconnection` Function

### Signature
`checkParentDisconnection(pc: RTCPeerConnection, target: string): void`

### Parameters
- `pc`: An instance of `RTCPeerConnection` representing the peer connection to check for disconnection of the parent.
- `target`: A `string` value representing the ID of the target to which the connection belongs.

### Return Value
None

### Description
The `checkParentDisconnection` function checks for the disconnection of the parent. 

- It takes in two parameters, `pc` and `target`. 
- The `pc` parameter is an instance of `RTCPeerConnection` representing the peer connection to check for disconnection of the parent. 
- The `target` parameter is a `string` value representing the ID of the target to which the connection belongs.

The function creates an interval that checks if the `pc` parameter is not an audience connection and if the connection is not alive. 
- If the connection is not alive, the function sets the `parentDC` property of the class to `true`. 
- The `restartEverything` function is then called with the `pc` and `target` parameters as arguments. 
- The `setInterval` is then cleared.

It should be noted that the `parentDisconnectionTimeOut` property of the class is used to specify the interval timeout value.

---

## `newPeerConnectionInstance` Function

### Signature
`newPeerConnectionInstance(target: String, theStream, isAdience = false: Boolean): RTCPeerConnection`

### Parameters
None

### Return Value
- The function returns a new `RTCPeerConnection` instance.

### Description
This function `newPeerConnectionInstance()` creates a new `RTCPeerConnection` instance by providing its configuration object `myPeerConnectionConfig` (previously defined as the default configuration object with the required settings), along with the target, the stream to be transmitted, and a boolean `isAdience` which is false by default.

- First, the function initializes some variables including `intervalId`, `peerConnection`, and ``dataChannel``. 
- The newly created `RTCPeerConnection` object then has the property `isAdience` set to the value of the `isAdience` parameter provided by the function call. 
- The property `alive` is also set to `true`.

A data channel is created for the connection and is named "chat" using the `createDataChannel()` method. 

- An `onopen` event handler is then added to `dataChannel` that will invoke the `onDataChannelOpened()` method passing in `dataChannel`, `target`, and `peerConnection` as arguments.
- In addition to the data channel, an `ondatachannel` event handler is added to `peerConnection` to handle data received on the connection. 
- When a data message is received, the function checks whether it was sent by the parent connection, and if so, sets `peerConnection.alive` to `true`. 
- The function then calls `checkParentDisconnection()` with `peerConnection` and `target` as arguments to handle disconnections.

Next, `onerror`, `onbufferedamountlow`, and `onclose` event handlers are added to `receive`, which is the receiving end of the data channel, to handle errors, low buffer thresholds, and the closing of the data channel.

`peerConnection` is then set up with several other event handlers, including `onconnectionstatechange`, `onicecandidate`, `onnegotiationneeded`, and `ontrack`.

- `onconnectionstatechange` is a general handler that logs the current connection state.
- `onicecandidate` sends a JSON string to the server, containing information about the ice candidate which is received from the server.
- `onnegotiationneeded` is responsible for preparing the video offer and sending it to the server.
- `ontrack` is the handler that is invoked whenever a new track is added to `peerConnection`. The handler receives an `event` object which includes a `stream` array that includes details about the stream received from the connection. This handler first checks whether the `localStream` has been added previously to the stream, and whether the `newTrackCallback` is already defined for this stream. If the `remoteStreamCallback` function has been defined, it will be invoked, passing the `stream` object as an argument. 
- The `remoteStreams` array will be updated with the new `stream` object. 


Finally, the function sends a message to the server to notify it that the stream has been received.

---

## `createOrGetPeerConnection` Function

### Signature
`createOrGetPeerConnection(audienceName: string, isAudience?: boolean): RTCPeerConnection`

### Parameters
- `audienceName`: A `string` value representing the name of the audience for the peer connection.
- `isAudience`: An optional `boolean` value representing whether the connection is for an audience or not. Defaults to `false`.

### Return Value
The function returns an `RTCPeerConnection` object representing the peer connection for the specified `audienceName`.

### Description
The `createOrGetPeerConnection` function takes in two parameters, `audienceName` and `isAudience`, representing the name of the audience and whether the connection is for an audience or not. 

- If a peer connection with the same `audienceName` already exists in `myPeerConnectionArray`, the function returns it.
- Otherwise, it calls the `newPeerConnectionInstance` function to create a new peer connection and stores it in `myPeerConnectionArray` with `audienceName` as the key. 
- Finally, the function returns the newly created peer connection object.

It should be noted that the `newPeerConnectionInstance` function needs to be defined in order for `createOrGetPeerConnection` to work as intended.

---

## `connectToAudience` Function

### Signature
`connectToAudience(audienceName: string): void`

### Parameters
- `audienceName`: A `string` value representing the name of the audience to connect to.

### Return Value
None

### Description
The `connectToAudience` function takes in a single `audienceName` parameter, which is used to specify the name of the audience to connect to. 

- The function first updates the status to indicate that it is connecting to the specified audience. 
- If there is no local stream and no remote streams available, the function returns without doing anything.
- If there is no existing peer connection to the specified audience, the function creates a new peer connection instance using the `newPeerConnectionInstance` function and adds it to the `myPeerConnectionArray` object using the `audienceName` as the key. 
- If a local stream or remote streams are available, the function adds the appropriate tracks to the peer connection instance using the `addTrack` method.

It uses the `newPeerConnectionInstance` function which is explained earlier.

---

## `sendStreamTo` Function

### Signature
`sendStreamTo(target: string, stream: MediaStream): void`

### Parameters
- `target`: A `string` representing the target audience name.
- `stream`: A `MediaStream` object representing the stream to be sent.

### Return Value
None

### Description
The `sendStreamTo` function takes in two parameters: `target` and `stream`. 

- The `target` parameter is a `string` representing the target audience name, while the `stream` parameter is a `MediaStream` object representing the stream to be sent. 
- The function first calls the `createOrGetPeerConnection` method to get or create a new `RTCPeerConnection` object associated with the specified target audience. 
- If a new `RTCPeerConnection` object is created, it is assigned to the `peerConnection` constant. 
- The `stream` tracks are then added to the `peerConnection` using the `addTrack` method.

If the `createOrGetPeerConnection` method returns an existing `RTCPeerConnection` object associated with the specified target audience, the `addTrack` method is called on the existing `RTCPeerConnection` object instead.

---

## `start` Function

### Signature
`async start(turn: boolean): Promise<void>`

### Parameters
- `turn`: A `boolean` value specifying whether to use a TURN server for the peer connections or not.

### Return Value
- The function returns a `Promise` that resolves to `void`.

### Description
The `start` function is an asynchronous function that takes in a single `turn` parameter to specify whether to use a TURN server for the peer connections or not. 

- If `turn` is `false`, the `myPeerConnectionConfig` is updated to remove the TURN servers from the list of `iceServers`. 
- The function then proceeds to start the appropriate process based on the `role` and `constraints` of the instance. 
- If the instance is in `broadcast` role, the function starts the broadcasting process, else it starts reading the broadcast. 

If there are no `audio` or `video` constraints specified, the `raise hand` element is removed from the HTML. 

The function updates the status with the current process at each step of the way. The function returns a promise that resolves to `void`. 

---

## `disableVideo` Function

### Signature
`disableVideo(enabled: boolean): void`

### Parameters
- `enabled`: A `boolean` value representing whether the video should be enabled or disabled.

### Return Value
None

### Description
The `disableVideo` function takes in a single `enabled` parameter, which is used to enable or disable the video tracks. 

- It does so by using a `forEach` loop to iterate through all tracks of the `localStream` and check whether the kind of the track is 'video'. 
- If it is, the `enabled` parameter is used to set the value of the `enabled` property of the track. 
- If `enabled` is `true`, the track is enabled, and if `enabled` is `false`, the track is disabled. 

This function can be used to enable or disable the video tracks of a `localStream`. 

---

## `disableAudio` Function

### Signature
`disableAudio(enabled: boolean): void`

### Parameters
- `enabled`: A `boolean` value representing the desired state of the audio track. If `true`, the audio track is enabled. If `false`, the audio track is disabled.

### Return Value
None

### Description
The `disableAudio` function takes in a single `enabled` parameter, which is used to enable or disable the audio track. 

- The function uses the `getTracks` method to get all the tracks in the `localStream`. 
- It then iterates through each track and checks if the track kind is `audio`. 
- If it is `audio`, then the `enabled` value is set to the `track.enabled` property. 
- If `enabled` is `true`, then the audio track is enabled. If `enabled` is `false`, then the audio track is disabled. 

---

## `wait` Function

### Signature
`wait(mil: number): Promise<void>`

### Parameters
- `mil`: A `number` value representing the amount of time to wait in milliseconds.

### Return Value
- The function returns a `Promise` object that resolves to `void` after the specified amount of time has passed.

### Description
The `wait` function takes in a single `mil` parameter, which is used to determine the amount of time to wait before resolving the `Promise`. 

- The function creates a new `Promise` object that waits for the specified amount of time using `setTimeout` and then resolves the `Promise`. 
- This function can be used to pause the execution of a program for a certain amount of time.


---

## `getBroadcasterStatus` Function

### Signature
`getBroadcasterStatus(): Promise<string>`

### Parameters
None

### Return Value
The function returns a `Promise` that resolves to a `string` value representing the status of the broadcaster.

### Description
The `getBroadcasterStatus` function sends a message to the server requesting the status of the broadcaster. 

- It then enters a loop that waits for a response from the server with a maximum number of iterations specified by the `max` constant, or until it receives a response. 
- If no response is received, the function will reject the `Promise` with an error.
- Once the response is received, the function resolves the `Promise` with the value of the `broadcasterStatus` property, which was set in the `handleMessage` function.
- If the `reconnect` parameter is `true`, the `startProcedure` function is called after resolving the `Promise`.

It uses the `wait` function which I explained above.

---

## `getSupportedConstraints` Function

### Signature
`getSupportedConstraints(): Promise<void>`

### Parameters
None

### Return Value
- The function returns a `Promise` which resolves to `void`.

### Description
The `getSupportedConstraints` function is an `async` function that retrieves the user's media devices and checks if there are any available audio and/or video input devices. 

- If an audio and/or video input device is not found, the corresponding property in the `this.constraints` object is set to `false`. 
- If the `this.constraintResults` callback function is defined, it is called with the `this.constraints` object as its argument. 
- If the user denies permission to access the media devices, an error will be thrown.

Note that the `navigator.mediaDevices.enumerateDevices()` method is used to retrieve the user's media devices.

- This method returns a promise that resolves to an array of `MediaDeviceInfo` objects, which contain information about each of the user's media devices. 
- The `kind` property of each `MediaDeviceInfo` object is used to determine if the device is an audio or video input device.

---

## `updateTheStatus` Function

### Signature
`updateTheStatus(status: string): void`

### Parameters
- `status`: A `string` value representing the status to be updated.

### Return Value
None

### Description
The `updateTheStatus` function takes in a single `status` parameter, which represents the status to be updated. 

- The function calls the `updateStatus` function, if it is defined, passing in the `status` parameter. 
- If the `updateStatus` function is not defined or fails to execute, the function does not throw any errors.

---

## `lowerHand` Function

### Signature
`lowerHand(): Promise<void>`

### Parameters
None

### Return Value
- The function returns a `Promise` that resolves to `void`.

### Description
The `lowerHand` function is an asynchronous function that lowers the hand of the local user. 

- It first checks if `localStream` is present, and then retrieves the first peer connection object from the `myPeerConnectionArray` object. 
- It then gets the IDs of all tracks in `localStream` and all senders in the peer connection. 
- After that, it loops through each track ID and for each track ID it loops through all senders to check if the sender's track ID is the same as the current track ID. 
- If the track IDs match, the sender is removed from the peer connection. 
- Then, it stops all tracks in `localStream` and sets `localStream` to `null`. 

Finally, it sets `startedRaiseHand` to `false` to indicate that the local user has lowered their hand.

---

## `spreadLocalStream` Function

### Signature
`spreadLocalStream(): void`

### Return Value
None

### Parameters
None

### Description
The `spreadLocalStream` function sends the local stream to each adience connected to the broadcaster.

- It loops through each connection in `myPeerConnectionArray` and checks if the connection is an audience connection, and if so, it calls the `sendStreamTo` function, passing in the target and local stream as arguments.

---

## `constructor` Function

### Signature
`constructor(role: string, options: object)`

### Parameters
- `role`: A `string` value representing the role of the user (`broadcaster` OR `audience`).
- `options`: An optional `object` containing additional options to configure the object.

### Return Value
- It's a `constructor` function for `SparkRTC` Class, does not return anything.

### Description
The `constructor` function is the constructor of a class and is used to initialize the object's properties. 

- The function takes in two parameters, `role` and `options`, which are used to configure the object. 
- The `role` parameter is a `string` that represents the role of the user (e.g. `publisher`, `subscriber`, etc.), and the `options` parameter is an optional `object` that contains additional options to configure the object.

The function sets the following properties based on the `options` parameter:

- `localStreamChangeCallback`
- `remoteStreamCallback`
- `remoteStreamDCCallback`
- `signalingDisconnectedCallback`
- `treeCallback`
- `raiseHandConfirmation`
- `newTrackCallback`
- `startProcedure`
- `log`
- `constraintResults`
- `updateStatus`

If a property is not provided in the `options` parameter, it is set to a default value. The function also logs a message indicating the role of the user.

---

.

.  

# Next we will look at `common.js` File

<!-- from here we will start common.JS -->

`Common.js` is the file contains some genenric fucntions, helper functions used by `sparkRTC.js` and `Ui.js`. let's dive into each function for futher understanding.

---

## `getWsUrl` Function

### Signature
`getWsUrl(): string`

### Return Value
- The function returns a `string` value representing the WebSocket URL based on the current window location.

### Paramters
None

### Description
The `getWsUrl` function is used to return the WebSocket URL based on the current window location. 

- The function uses the `window.location` object to get the current URL and then splits it to get the base URL and protocol. 
- It then constructs the WebSocket URL using the `baseUrl` and `protocol` variables and returns it.

---

## `clearScreen` Function

### Signature
`clearScreen(): void`

### Return Value
None

### Parameters
None

### Description
The `clearScreen` function is used to clear the screen by removing all the child nodes of the `screen` element. 

- It also makes the `raise_hand` and `share_screen` elements visible again by setting their `display` style property to `'block'`.
- The `clearScreen` function first makes the `raise_hand` and `share_screen` elements visible again by setting their `display` style property to `'block'`. 
- It then gets a reference to the `screen` element and checks if it exists. 
- If it does not exist, the function returns immediately. 
- If it does exist, the function removes all its child nodes by calling `removeChild` on each child node until there are no more child nodes left.

---


## `createSparkRTC` Function

### Signature
`createSparkRTC(): SparkRTC`

### Return Value
- The function returns a `SparkRTC` instance.

### Parameters
None

### Description
The `createSparkRTC` function is used to create a `SparkRTC` instance. 

- The function checks the myRole variable and creates the instance accordingly. 
- If the myRole is '`broadcast`', it creates the instance with certain options and returns it. 
- If the `myRole` is not '`broadcast`', it creates the instance with different options and returns it.

The options for the '`broadcast`' role includes following callbacks.

#### `localStreamChangeCallback`

- This callback function is called when the local stream changes. It takes a single parameter `stream`, which is an object representing the local media stream. 
- It sets the source of the video element with the id '`localVideo`' to the new stream and enables the audio and video controls.

#### `remoteStreamCallback`

- This callback function is called when a new remote stream is received. It takes a single parameter `stream`, which is an object representing the remote media stream. 
- It first creates a unique id for the video element using the remote `stream id` and then creates a new `video` element with the generated id. 
- It then sets the source of the video element to the received stream and plays the video.

#### `remoteStreamDCCallback`

- This callback function is called when a remote stream is `disconnected`. 
- It takes a single parameter stream, which is an object representing the remote media stream. 
- It first tries to find the video element with the id 'remoteVideo' and if it exists, it removes the element. 
- If it does not exist, it tries to find the video element with the id 'localVideo' and removes it.

#### `signalingDisconnectedCallback`

- This callback function is called when the signaling connection is `disconnected`. It clears the screen.

#### `treeCallback`

- This callback function is called when a `tree` is received. It takes a single parameter tree, which is a JSON string representing the tree data. 
- It first tries to parse the JSON string to extract the tree data. 
- It then checks if the tree data exists and if it does, it calls the `draw` function of the `graph` object to draw the tree.

#### `raiseHandConfirmation`

- This callback function is called when the user raises their hand. It takes a single parameter msg, which is a string representing the message to be displayed. 
- It currently returns true which will always confirm the raise hand action.

#### `startProcedure`

- This callback function is called when the user starts the procedure. It is an asynchronous function that first calls the handleClick function.

#### `log`

- This callback function is called when the `log` function is called. It takes a single parameter log, which is a string representing the log message. It then calls the `addLog` function to add the log to the log list.

#### `constraintResults`

- This callback function is called when the media `constraints` are set. It takes a single parameter constraints, which is an object representing the media constraints. 
- It then checks if the audio constraint is false and removes the '`mic`' element from the page.

#### `updateStatus`
- This callback function is called when the `status` is updated. It takes a single parameter status, which is a string representing the updated status message. It sets the text of the element with the id '`status`' to the new status message.


The options for the audience role include only following callbacks:  

- `remoteStreamCallback`, a `remoteStreamDCCallback`, a `signalingDisconnectedCallback`, a `startProcedure`, a `log`, and an `updateStatus` callback.

---

## `registerNetworkEvent` Function

### Signature
`registerNetworkEvent(): void`

### Return Value
None

### Parameters
None

### Description
The `registerNetworkEvent` function is used to monitor the network status of the browser.

- If the `navigator.connection` object is available, the function calls the `handleNetworkStatus` function to set the initial network status. 
- It also adds an `onchange` event listener to the `navigator.connection` object to handle changes in the network status.


---

## `handleNetworkStatus` Function

### Signature
`handleNetworkStatus(event: Event): void`

### Parameters
- `event`: An optional `Event` object representing the event that triggered the function.

### Return Value
None

### Description
The `handleNetworkStatus` function takes in an optional `event` parameter, which represents the event that triggered the function. 

- It first checks if the `currentTarget` property of the `event` object is defined, and if not, it uses the `navigator.connection` object to get the current network connection status. 
- It then checks if the `downlink` property of the network connection is less than or equal to `1`, which indicates a slow network. If the network is slow, the function calls the `onNetworkIsSlow` function. 
- If the network is normal, the function calls the `onNetworkIsNormal` function.

---

.

.  

# Next we will look at `Graph.js` File

<!-- from here we will start Graph.JS -->

Graph class is a being used to Draw a Graph in form of `Binary Tree`, to display Names of all the Audienecs and Broadcasters.

## `setup` Function

### Signature
`setup(): void`

### Parameters
None

### Return Value
None

### Description
The `setup` function is a method of a class that sets the width and height of the SVG element, and declares a tree layout with the specified size. 

- It uses the D3 library to select the SVG element and set its attributes for width and height. 
- It then appends a group element with a `transform` attribute set to translate by 30 on the x-axis. 
- The `treemap` variable is assigned a tree layout with the specified height and width. 

---

## `update` Function

### Signature
`update(source): void`

### Return Value
None

### Parameters
`source`: A tree `node` object representing the root of the tree structure to be updated.

### Description
- The `update` function is used to update the `tree` structure represented in the SVG element.
- It takes a source node representing the root of the tree and computes the new layout of the tree. The `nodes` and `links` of the tree are then updated and transitioned to their new positions using `D3` transitions.

#### Nodes

- The positions of the nodes are computed by normalizing for fixed depth and setting their `x` and `y` attributes based on their `depth`.
- The nodes are updated by selecting all `g` elements with class `node` and data binding them to the nodes using their unique `id` attribute.
- Any new nodes are then created and `positioned` at the parent's previous `position`.
- The `circles` representing the nodes are updated with their new attributes and `style`. 
- The `radius` is determined by whether the node has children or not. 
- The fill `color` is set to either `light steel blue` or `white` based on whether the node has hidden children or not.
- The `text` representing the node's name is updated with its new value.

#### Links

- The `links` are updated by selecting all `path` elements with class link and data binding them to the links using their unique id attribute. 
- Any new links are then created and positioned at the parent's previous `position`.
- The links are `transitioned` to their new positions based on their parent and `child` nodes.
- Any `exiting` nodes or links are removed.

#### Transitions

- All transitions are defined with a `duration` specified by the duration property.
- The transitions for updating nodes and links use the `D3` merge method to transition the updated elements to their new positions.

#### Source

- The `source` parameter represents the root of the tree to be updated.
- The source parameter is used to `position` any new nodes or links at the parent's previous position. 
- It is also used to `transition` any `exiting` nodes or links back to the parent's previous position.

---

## `draw` Function

### Signature
`draw(data: object): void`

### Parameters
- `data`: An `object` representing the data to be visualized.

### Return Value
None

### Description
The `draw` function is responsible for rendering the tree structure on the screen. 

- It takes in a single `data` parameter, which is an `object` representing the data to be visualized. 
- The function first calls the `setup` method to set up the initial layout of the tree. 
- It then uses the `d3.hierarchy` method to assign the parent, children, height, and depth to each node in the data. 
- The root node is set to be the topmost node of the tree. 
- The `treeData` is calculated using the `d3.tree()` method and the `root` node, and the `update` method is called to render the tree on the screen.

It should be noted that the `height` and `width` values used in the `d3.tree()` method are not defined in the code snippet provided and would need to be specified in order for the `draw` function to work as intended. 

---