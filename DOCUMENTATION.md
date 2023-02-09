
# **LOGJAM**

### Code Documentation 

In this file I will be explaining the code of **logjam**, I will try to use simple language and explain the usage, working & why code is writen the way its written.

---


# `..\frontend\ui.js` File

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

## `arrangeVideoContainers` Function

### Signature
`arrangeVideoContainers(): void`


### Parameters
N/A

### Return Value
N/A

### Description
The `arrangeVideoContainers` function is used to arrange multiple video containers on the screen by adjusting their height accordingly. The function starts by using the `document.getElementById` and `getElementsByClassName` methods to retrieve all elements with the class `video-container` within an element with the id `screen`.

The function then calculates the total number of video containers and uses this value to calculate the `flexRatio` and `maxHeight` values, which are used to set the `flex` and `max-height` CSS properties for each video container.

Finally, the `arrangeVideoContainers` function uses the `Array.from` method to iterate over the video containers, setting the `flex` and `max-height` properties for each container using the `style.setProperty` method.

---

## `onCameraButtonClick` Function

### Signature
`function onCameraButtonClick()`


### Parameters
N/A

### Return Value
N/A

### Description
This function takes action on camera button click and changes the status of the camera, 
based on the status of the camera the source of the image is changed to either CAMERA_OFF or CAMERA_ON.
It also calls sparkRTC.disableVideo to disable the video stream.

---

## `onMicButtonClick` Function

### Signature
`function onMicButtonClick()`


### Parameters
N/A

### Return Value
N/A

### Description
This function takes certain actions on mic button click, it mute or unmutes the mic based on its status.
It changes the status of the mic and updates the source of the image to either MIC_OFF or MIC_ON.
It also calls sparkRTC.disableAudio to disable the audio stream.

---

## `createVideoElement` Function

### Signature
`function createVideoElement(videoId, muted = false)`


### Parameters
- `videoId`: Id of the video element to be created.
- `muted`: A boolean value indicating whether the video should be muted or not.

### Return Value
The created video element.

### Description
This function creates a new video element to display a video stream (local or remote).
It creates a container div with class `video-container` and appends a `video` element to it.
It sets the properties of the video element, such as id, autoplay, playsInline, and mute status.
Finally, it appends the container to the element with id `screen` and returns the created video element.

---

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

## `onNetworkIsSlow(downlink)` Function

### Signature
`function onNetworkIsSlow(downlink)`


### Parameters
- **downlink**: `any` - The downlink value indicating the network speed

### Return value
- Returns nothing.

### Description
This function is used to display the network status as slow. If the downlink value is greater than 0, it sets the status as very slow and displays an alert message indicating the network status is very slow. If the downlink value is 0, it sets the status as disconnected and displays an alert message indicating the network is disconnected.
