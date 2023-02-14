
# **LOGJAM**

### Code Documentation 

In this file I will be explaining the code of **logjam**, I will try to use simple language and explain the usage, working & why code is writen the way its written.

<!-- Starting with UI.JS File from Here -->
.  

.  

# Let's start with `ui.js` File

<!-- f1 -->
## `makeId` Function

### Signature
`makeId(length: number): string`


### Parameters
- `length`: A `number` value representing the desired length of the returned string.

### Return Value
The function returns a `string` value representing a random combination of characters with the length specified in the `length` parameter.

### Description
The `makeId` function takes in a single `length` parameter, which is used to determine the length of the returned string. The function uses a `for` loop to generate a random combination of characters by selecting a random character from a predefined set of characters stored in the `CHARACTERS` constant and appending it to the `result` string. The loop continues until the length of `result` is equal to the `length` specified in the `length` parameter. Finally, the `result` string is returned as the result of the `makeId` function.

It should be noted that the `CHARACTERS` constant is not defined in the code snippet provided, and its value would need to be specified in order for the `makeId` function to work as intended.

 
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
The `arrangeVideoContainers` function is used to arrange multiple video containers on the screen by adjusting their height accordingly. The function starts by using the `document.getElementById` and `getElementsByClassName` methods to retrieve all elements with the class `video-container` within an element with the id `screen`.

The function then calculates the total number of video containers and uses this value to calculate the `flexRatio` and `maxHeight` values, which are used to set the `flex` and `max-height` CSS properties for each video container.

Finally, the `arrangeVideoContainers` function uses the `Array.from` method to iterate over the video containers, setting the `flex` and `max-height` properties for each container using the `style.setProperty` method.

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
This function takes action on camera button click and changes the status of the camera, based on the status of the camera the source of the image is changed to either `CAMERA_OFF` or `CAMERA_ON`. It also calls `sparkRTC.disableVideo` to disable the video stream.

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
This function takes certain actions on mic button click, it mute or unmutes the mic based on its status. It changes the status of the mic and updates the source of the image to either `MIC_OFF` or `MIC_ON`. It also calls `sparkRTC.disableAudio` to disable the audio stream.

---
<!-- f5 -->
## `createVideoElement` Function
  
### Signature
`function createVideoElement(videoId, muted = false)`


### Parameters
- `videoId`: Id of the video element to be created.
- `muted`: A boolean value indicating whether the video should be muted or not.

### Return Value
The created video element.

### Description
This function creates a new video element to display a video stream (local or remote). It creates a container div with class `video-container` and appends a `video` element to it. It sets the properties of the video element, such as id, autoplay, playsInline, and mute status. Finally, it appends the container to the element with id `screen` and returns the created video element.

---
<!-- f6 -->
## `getVideoElement(videoId)` Function

### Signature
`function getVideoElement(videoId)`


### Parameters
- **videoId**: `any` - The id of the video element to be retrieved

### Return value
- Returns a video element with the given ID, if not found creates a new video element and returns that.

### Description
This function is used to retrieve a video element from the screen using its ID. If a video element with the given ID is not found, the function creates a new video element with the given ID and returns it.

---
<!-- f7 -->
## `removeVideoElement(videoId)` Function

### Signature
`function removeVideoElement(videoId)`


### Parameters
- **videoId**: `any` - The id of the video element to be removed

### Return value
- Returns nothing.

### Description
This function is used to remove the video element from the screen using its ID.

---
<!-- f8 -->
## `onNetworkIsSlow(downlink)` Function

### Signature
`function onNetworkIsSlow(downlink)`


### Parameters
- **downlink**: `any` - The downlink value indicating the network speed

### Return value
- Returns nothing.

### Description
This function is used to display the network status as slow. If the downlink value is greater than 0, it sets the status as very slow and displays an alert message indicating the network status is very slow. If the downlink value is 0, it sets the status as disconnected and displays an alert message indicating the network is disconnected.

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
This function is used to hide the element with ID `net` from the HTML document. The display style of the element is set to `none`. This means that the element will not be `visible` on the page. The purpose of this function is to `remove` the notification of a network issue when the network is back to normal.

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
This function is used to `enable` or `disable` screen sharing in a real-time communication (RTC) application using WebRTC API. The function uses the `sparkRTC` object to start or stop the screen sharing. 

When the function is called for the first time, it starts the screen sharing by calling the `sparkRTC.startShareScreen()` method. The method returns a stream of the shared screen, which is assigned to the `shareScreenStream` variable. The function then updates the source of the local screen video element with the `new stream`. 

When the function is called again, it stops the screen sharing by stopping all the tracks of the `shareScreenStream` and setting it to `null`. The local screen video element is also updated to have a `null` source. The `removeVideoElement('localScreen')` method is used to remove the video element from the page.

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
This function is used to set the user's name in a Video application. The function first tries to retrieve the name from the **local storage** using the key `logjam_myName`. The retrieved name is then set as the value of the input field with ID `inputName`.

If the name cannot be retrieved from the **local storage**, the function generates a new name using the `makeId(20)` function and sets it as the value of the `myName` variable. The generated name is then stored in the local storage using the key `logjam_myName`.

The purpose of the function is to provide a unique name for each user in the Video application. The name is stored in the **local storage** so that the user can access it across multiple sessions.

--- 
<!-- f12 -->
## `handleClick` Function

### Signature
`async function handleClick(turn = true)`

### Parameters
`turn: boolean` (optional, default is `true`)

### Return Value
`False`

### Description
This function is used to handle the click event of the submit button in a Video application. The function first retrieves the name entered in the input field with ID `inputName` and stores it in the local storage using the key `logjam_myName`. The function then updates the visibility and display styles of the page and the `getName` section.

The function also calls the `start(turn)` method, which is used to start the Video session. The `turn` parameter is used to specify the type ice server, TURN or STUN.

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
This function is used to handle the resize event of the window in a Video application. The function first clears the previous timeout using the `clearTimeout(window.resizedFinished)` method.

The function then sets a new timeout using the `setTimeout` method. The timeout function calls two methods, `graph.draw(graph.treeData)` and `arrangeVideoContainers()`, to redraw the graph and rearrange the video containers.

The purpose of the function is to update the layout of the video window whenever the window is resized, ensuring that the graph and video containers are displayed correctly. The timeout is used to **prevent** the function from being called **too frequently** during rapid resizing.

---
<!-- f14 -->
## `getMyRole` Function

### Signature
`function getMyRole()`

### Parameters
None

### Return Value
`String` (either `"broadcast"` or `"audience"`)

### Description
This function is used to retrieve the `role` of the user in a video application. The function first retrieves the `query` string from the current URL using the `window.location.search` property.

The function then creates a new `URLSearchParams` object using the query string, allowing it to access the URL parameters. The function uses the `urlParams.get('role')` method to retrieve the value of the `role` parameter.

The function returns either `"broadcast"` or `"audience"` based on the value of the `role` parameter. If the parameter is not present or its value is not `"broadcast"`, the function returns `"audience"`.

The purpose of the function is to determine the user's role in the video session, either as a broadcaster or as an audience member.

---
<!-- f15 -->
## `getRoomName` Function

### Signature
`function getRoomName()`

### Parameters
None

### Return Value
`String` (room name)

### Description
This function is used to retrieve the `name` of the `room` in a video application. The function first retrieves the query string from the current URL using the `window.location.search` property.

The function then creates a new `URLSearchParams` object using the query string, allowing it to access the URL parameters. The function uses the `urlParams.get('room')` method to retrieve the value of the `room` parameter.

The function returns the value of the `room` parameter, which is the name of the room.

The purpose of the function is to determine the name of the room the user is joining in the video session.

---
<!-- f16 -->
## `getDebug` Function

### Signature
`function getDebug()`

### Parameters
None

### Return Value
`Boolean` (true or false)

### Description
This function is used to retrieve the value of the `debug` parameter in a video application. The function first retrieves the query string from the current URL using the `window.location.search` property.

The function then creates a new `URLSearchParams` object using the query string, allowing it to access the URL parameters. The function uses the `urlParams.get('debug')` method to retrieve the value of the `debug` parameter.

The function returns the value of the `debug` parameter converted to a `boolean` value, which indicates whether or not the debug mode is **enabled**.

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
This function is used to setup the signaling socket in a video application. The function uses the `sparkRTC.setupSignalingSocket()` method and passes the following parameters to it:
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
This function is used to start the broadcasting or listening to the broadcast in a chat application. The function calls the `setupSignalingSocket()` function to initiate the signaling socket and then uses the `sparkRTC.start()` method and passes the "turn" parameter to it.

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
The function `onLoad` is used to execute some statements when the page loads. It sets the role of the user and room name based on the URL. It creates a new instance of the SparkRTC class and sets the `myName` property. The logs are hidden if the `debug` query string parameter is not present. It also calls `handleResize` when the window is resized and initializes the `Graph` object.

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
The function `onRaiseHand` is used to handle the click event on the `raise_hand` button. If the `raise_hand` button status is `on`, the function raises the hand by calling `sparkRTC.raiseHand` and the microphone and camera icons are displayed.

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
The `addLog` function takes in a `string` log as a parameter and adds it to the logs. It does this by getting the logs element from the document using `document.getElementById('logs')`, creating a new paragraph element using `document.createElement('p')`, setting the innerText of the paragraph element to the log parameter, and finally, appending the paragraph element to the logs element.

---
<!-- f22 -->
## `enableAudioVideoControls` Function

### Signature
`function enableAudioVideoControls()`

### Parameters
- None

### Return Value
- None

### Description
The function `enableAudioVideoControls` sets the visibility of three different icons to display. It sets the display property of the elements with ids `mic`, `camera`, and `share_screen` to an empty string. This makes the elements visible on the page.

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
This function is used to `hide` the audio and video controls from the UI. It gets the elements with the id `mic`, `camera`, and `share_screen` from the document and sets their display property to `none`. 

---

.

.  

# Next we will look at `SparkRTC.js` File

<!-- from here we will start SPARKRTC.JS -->
<!-- # `..\SparkRTC.js` File -->

<!-- f1 -->
## `handleVideoOfferMsg` Function

### Signature
`handleVideoOfferMsg = async (msg) => {}`

### Parameters
- msg: an object containing data for the message

### Return Value
None

### Description
This function is called when the client receives a message of type `video-offer` or `alt-video-offer`. This function will ***create*** a `peer connection` or ***retrieve*** an existing one for the message sender and set the `remote description` with the `sdp` provided in the message. Then, it will create an answer for the offer and send it back to the message sender.

---
<!-- f2 -->
## `handleMessage` Function
### Signature
`handleMessage = async (event) => {}`

### Parameters
- event: an event object with the data of the message received on the web socket.

### Return Value
None

### Description
This function is the main socket handler for messages received on the web socket. It parses the message and based on the message type, makes decisions on what to do next.

The function first tries to parse the message into a JSON object, and if that fails it returns. Then, it sets the `data` and `type` properties of the message to the right values, as they may be in different cases in the message.

Next, the function has a switch statement to handle different types of messages, such as `video-offer`, `video-answer`, `new-ice-candidate`, `role`, `start`, `add_audience`, and `alt-broadcast-approve`. The actions taken for each type of message may include creating or retrieving a peer connection, updating the status, connecting to an audience, sending a stream to an audience, and more.

---
<!-- f3 -->
## `broadcasterUserId` Function

### Signature
`broadcasterUserId : () => string`

### Parameters
None

### Return Value
`string`

### Description
This function returns the user id of the broadcaster by iterating through the `myPeerConnectionArray` object and checking for the user with `isAudience` property set to `false`. If no such user is found, it returns `null`.

---
  
## `ping` Function

### Signature
`ping(): void`

### Parameters
None

### Return Value
None

### Description
The `ping` function is an arrow function that sends a message over a socket connection using the `socket.send()` method. The message is a JSON string that contains a `type` property whose value is either `"tree"` or `"ping"`, depending on whether the `treeCallback` property is truthy or falsy. 

If `treeCallback` is truthy, the `type` property value will be `"tree"`, and if it is falsy, the `type` property value will be `"ping"`. The `send()` method sends the message to the server through the socket connection.

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
The function returns a `Promise` that resolves to a `WebSocket` object.

### Description
The `setupSignalingSocket` function is an arrow function that takes in three parameters - `url`, `myName`, and `roomName`. It returns a `Promise` that resolves to a `WebSocket` object. 

The function first logs the input parameters to the console. If a `pingInterval` property exists, it is cleared. Then, if `myName` and `roomName` parameters are provided, they are assigned to the `myName` and `roomName` properties of the function. The function then creates a new `WebSocket` object using the specified `url` and `roomName`.

The function sets up event handlers for the `onmessage`, `onopen`, `onclose`, and `onerror` events of the `WebSocket` object. 

When the `onopen` event is triggered, the function sends a `"start"` message to the server, containing the `myName` parameter in its `data` property, using the `send()` method. It also sets up a `pingInterval` that calls the `ping` function every 5 seconds to keep the connection alive. Finally, it resolves the `Promise` with the `WebSocket` object.

When the `onclose` event is triggered, the function sets various properties to their default values and calls the `signalingDisconnectedCallback` and `startProcedure` callbacks if they exist.

When the `onerror` event is triggered, the function logs the error message to the console, rejects the `Promise`, and displays an alert to the user.

This function is useful for setting up a `WebSocket` connection to a signaling server, which can then be used to exchange information between clients in a WebRTC session.

---

## `startBroadcasting` Function

### Signature
`startBroadcasting(data: string): Promise<MediaStream | undefined>`

### Parameters
- `data`: A `string` value representing the role of the user.

### Return Value
The function returns a `Promise` that resolves with a `MediaStream` object representing the user's local media stream.

### Description
The `startBroadcasting` function takes in a single optional `data` parameter, which is used to determine the role of the user. If the user's `localStream` is not already set, the function attempts to retrieve the local stream using the `getUserMedia` method from the `navigator.mediaDevices` API. If successful, the `localStream` is stored in the `this.localStream` property and added to the `this.remoteStreams` array. The function then sends a message to the signaling server with the `data` and type `role`. Finally, the function returns a Promise that resolves with the `localStream`.

If an error occurs while attempting to retrieve the local stream, the function updates the status with an error message and displays an alert message. The function also logs the error message to the console and returns `undefined`.

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


