const CAMERA_ON = "images/cam-on.png";
const CAMERA_OFF = "images/cam-off.png";
const MIC_ON = "images/mic-on.png";
const MIC_OFF = "images/mic-off.png";
const SCREEN_ON = "images/screen-on.png";
const SCREEN_OFF = "images/screen-off.png";
// const SPARK_LOGO = "images/spark-logo.png";
const CHARACTERS =
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

let graph;
let sparkRTC;
let myName;
let myRole;
let shareScreenStream;
let roomName;

function makeId(length) {
	let result = "";
	for (let i = 0; i < length; i++) {
		result += CHARACTERS.charAt(Math.floor(Math.random() * CHARACTERS.length));
	}
	return result;
}

function mapping(positionsAndSizes, participants) {
	// We have an ordered list of positions and sizes
	//
	// And, we have an ordered list of participants.
	//
	// In the extreme case, the count of pos and sizes are smaller than the
	// count of participants. This is what we will code to.
	//
	// Now, not only that, but we also have a series of DOM elements. These DOM
	// elements are not ordered according to the list of participants, and will
	// thus need to be mapped.
	//
	// So the trick would be to grab the list of DOM elements, and then generate
	// a map

	// First generate a map of DOM elements and map them to participants list

	const videoContainers = document
		.getElementById("screen")
		.getElementsByClassName("video-container");

	const videosMap = new Map();

	for (const div of Array.from(videoContainers)) {
		const id = div
			.querySelector("video")
			.id.replace("localVideo-", "")
			.replace("remoteVideo-", "");

        console.log(id);

		videosMap.set(
			id === "localVideo" ? div.querySelector("video").srcObject?.id : id,
			div
		);
	}

	console.log(videosMap);

	for (const [i, participantId] of [...participants].entries()) {
		const div =
			videosMap.get(participantId) ||
			getVideoElement("localVideo").parentElement;

		div.style.setProperty("position", "absolute");

		const posSz = positionsAndSizes[i];
		if (!!posSz) {
			const { pos, size } = posSz;
			console.log("Have", i, participantId);
			div.style.setProperty("width", `${size[0]}px`);
			div.style.setProperty("height", `${size[1]}px`);
			div.style.setProperty("top", `${pos[1]}px`);
			div.style.setProperty("left", `${pos[0]}px`);
		} else {
			console.log("Don't have", i, participantId);
			div.style.setProperty("width", "200px");
			div.style.setProperty("height", "200px");
			div.style.setProperty("top", "0");
			div.style.setProperty("left", "0");
			q;
		}
	}
}

function arrangeVideoContainers() {
	console.log("Video layout", sparkRTC.metaData);

	const videoLayout = JSON.parse(sparkRTC.metaData.videoLayout ?? "null");
	console.log(videoLayout);

	switch (videoLayout?.type) {
		case "silly-frame":
			console.log("Silly frame");

			const videoContainers = document
				.getElementById("screen")
				.getElementsByClassName("video-container");

			for (const [i, div] of Array.from(videoContainers).entries()) {
				div.style.setProperty("position", "absolute");

				if (i === 0) {
					div.style.setProperty("width", `${videoLayout.meta.adminSize[0]}px`);
					div.style.setProperty("height", `${videoLayout.meta.adminSize[1]}px`);
					div.style.setProperty(
						"top",
						`${videoLayout.meta.adminLocation[1]}px`
					);
					div.style.setProperty(
						"left",
						`${videoLayout.meta.adminLocation[0]}px`
					);
				}
			}
			break;
		case "blue":
			{
				const positionsAndSizes = [
					{ pos: [300, 10], size: [800, 500] },
					{ pos: [10, 510], size: [120, 120] },
					{ pos: [480, 310], size: [120, 120] },
				];

				mapping(positionsAndSizes, videoLayout.meta.participants);

				page.style.backgroundColor = "blue";
			}
			break;
		case "red":
			{
				const positionsAndSizes = [
					{ pos: [10, 10], size: [800, 500] },
					{ pos: [10, 510], size: [120, 120] },
					{ pos: [480, 510], size: [120, 120] },
				];

				mapping(positionsAndSizes, videoLayout.meta.participants);

				page.style.backgroundColor = "red";
			}
			break;
		case "yellow":
			{
				const positionsAndSizes = [
					{ pos: [300, 300], size: [480, 300] },
					{ pos: [10, 610], size: [120, 120] },
					{ pos: [480, 610], size: [120, 120] },
				];

				mapping(positionsAndSizes, videoLayout.meta.participants);

				page.style.backgroundColor = "#aaaa00";
			}
			break;
		case "purple":
			{
				const positionsAndSizes = [
					{ pos: [10, 10], size: [480, 300] },
					{ pos: [10, 310], size: [120, 120] },
					{ pos: [480, 310], size: [120, 120] },
				];

				mapping(positionsAndSizes, videoLayout.meta.participants);

				page.style.backgroundColor = "#9900ff";
			}
			break;
		case "custom":
			{
				const positionsAndSizes = videoLayout.meta.placements;

				mapping(positionsAndSizes, videoLayout.meta.participants);

				page.style.backgroundImage = `url(${videoLayout.meta.background})`;
			}
			break;
		case "control-freak":
			{
				console.log(videoLayout.meta);
				console.log("Arranging video here");
				const videoContainers = document
					.getElementById("screen")
					.getElementsByClassName("video-container");

				const streamsMap = new Map(
					videoLayout.meta.participants.map(({ key, value }) => [key, value])
				);
				for (const [i, div] of Array.from(videoContainers).entries()) {
					const id = div
						.querySelector("video")
						.id.replace("localVideo-", "")
						.replace("remoteVideo-", "");

					div.style.setProperty("position", "absolute");
					console.error(id, streamsMap);

					if (streamsMap.has(id)) {
						const properties = streamsMap.get(id);
						div.style.setProperty("width", `${properties.size[0]}px`);
						div.style.setProperty("height", `${properties.size[1]}px`);
						div.style.setProperty("top", `${properties.location[0]}px`);
						div.style.setProperty("left", `${properties.location[1]}px`);
					} else {
						div.style.setProperty("width", "200px");
						div.style.setProperty("height", "200px");
						div.style.setProperty("top", "0");
						div.style.setProperty("left", "0");
					}
				}
			}
			break;
		case "tiled":
		default: {
			console.log("Arranging video here");
			const videoContainers = document
				.getElementById("screen")
				.getElementsByClassName("video-container");
			const videoCount = videoContainers.length;
			let flexRatio = 100 / Math.ceil(Math.sqrt(videoCount));
			let flex = "0 0 " + flexRatio + "%";
			let maxHeight =
				100 / Math.ceil(videoCount / Math.ceil(Math.sqrt(videoCount)));
			Array.from(videoContainers).forEach((div) => {
				div.style = {};
				div.style.setProperty("flex", flex);
				div.style.setProperty("max-height", maxHeight + "%");
			});
		}
	}
}

function onCameraButtonClick() {
	const img = document.getElementById("camera");
	if (img.dataset.status === "on") {
		img.dataset.status = "off";
		img.src = CAMERA_OFF;
		sparkRTC.disableVideo();
	} else {
		img.dataset.status = "on";
		img.src = CAMERA_ON;
		sparkRTC.disableVideo(true);
	}
}

function onMicButtonClick() {
	const img = document.getElementById("mic");
	if (img.dataset.status === "on") {
		img.dataset.status = "off";
		img.src = MIC_OFF;
		sparkRTC.disableAudio();
	} else {
		img.dataset.status = "on";
		img.src = MIC_ON;
		sparkRTC.disableAudio(true);
	}
}

function createVideoElement(videoId, userData, muted = false) {
	let container = document.createElement("div");
	container.className = "video-container";
	let video = document.createElement("video");
	video.id = videoId;
	video.autoplay = true;
	video.playsInline = true;
	video.muted = muted;
	container.appendChild(video);
	if (userData) {
		let details = document.createElement("div");
		details.innerText = userData.userName + "[" + userData.userRole + "]";
		container.appendChild(details);
	}
	document.getElementById("screen").appendChild(container);
	arrangeVideoContainers();
	return video;
}

function getVideoElement(videoId) {
	let video = document.getElementById(videoId);
	const userData = sparkRTC.getStreamDetails(videoId);
	return video ? video : createVideoElement(videoId, userData, true);
}

function removeVideoElement(videoId) {
	let video = document.getElementById(videoId);
	if (!video) return;
	let videoContainer = video.parentNode;
	if (!videoContainer) return;
	document.getElementById("screen").removeChild(videoContainer);
	arrangeVideoContainers();
}

async function onShareScreen() {
	const img = document.getElementById("share_screen");
	if (!shareScreenStream) {
		shareScreenStream = await sparkRTC.startShareScreen();
		if (shareScreenStream) {
			img.dataset.status = "on";
			img.src = SCREEN_ON;
			const localScreen = getVideoElement("localScreen");
			localScreen.srcObject = shareScreenStream;
		}
	} else {
		img.dataset.status = "off";
		img.src = SCREEN_OFF;
		shareScreenStream.getTracks().forEach((track) => track.stop());
		shareScreenStream = null;
		const localScreen = getVideoElement("localScreen");
		localScreen.srcObject = null;
		removeVideoElement("localScreen");
	}
}

function onRequestChangeBackground() {
	const input = document.createElement("input");
	input.type = "file";
	input.onchange = (_) => {
		const file = input.files[0];

		const formData = new FormData();
		formData.append("file", file);
		fetch("https://upload.logjam.server.group.video/file", {
			method: "POST",
			body: formData,
		}).then(async (res) => {
			const { path } = await res.json();
			// sparkRtc.metaData.backgroundUrl = `https://upload.logjam.server.group.video${path}`
			sparkRTC.socket.send(
				JSON.stringify({
					// This will set user specific MetaData
					type: "user-metadata-set",
					// This will set room specific MetaData
					// type: "metadata-set",
					data: JSON.stringify({
						...getMeta(),
						backgroundUrl: `https://upload.logjam.server.group.video${path}`,
					}),
				})
			);
		});
	};
	input.click();
}

// All background layout values:
//
// - contain
// - cover
// - tiled
// let currentBackgroundLayout = 'contain';

let currentBackgroundIndex = 0;
const possibleBackgrounds = ["contain", "cover", "tiled"];

function currentBackgroundLayout() {
	return possibleBackgrounds[currentBackgroundIndex];
}

let currentLayoutIndex = 0;
const possibleLayouts = [
	{ type: "tiled", meta: null },
	{
		type: "silly-frame",
		meta: { adminLocation: [171, 200], adminSize: [400, 400] },
	},
	{ type: "control-freak", meta: { participants: [] } },
	{ type: "blue", meta: { participants: [] } },
	{ type: "red", meta: { participants: [] } },
	{ type: "yellow", meta: { participants: [] } },
	{ type: "purple", meta: { participants: [] } },
	{
		type: "custom",
		meta: { participants: [], background: "", placements: [] },
	},
];

function currentLayoutJSON() {
	return JSON.stringify(currentLayout());
}

function currentLayout() {
	return possibleLayouts[currentLayoutIndex];
}

function getMeta() {
	return {
		backgroundUrl: sparkRTC.metaData.backgroundUrl,
		backgroundLayout: sparkRTC.metaData.backgroundLayout,
		videoLayout: sparkRTC.metaData.videoLayout,
	};
}

function getRandomInt(min, max) {
	// The maximum is exclusive and the minimum is inclusive
	min = Math.ceil(min);
	max = Math.floor(max);
	return Math.floor(Math.random() * (max - min)) + min;
}

function handleControlFreak(videoLayout) {
	videoLayout.meta.participants = Object.entries(sparkRTC.userStreamData).map(
		([key]) => ({
			key,
			value: {
				location: [getRandomInt(10, 1024), getRandomInt(10, 768)],
				size: [getRandomInt(200, 800), getRandomInt(150, 600)],
			},
		})
	);
}

function handleParticipants(videoLayout) {
	videoLayout.meta.participants = [
		getVideoElement("localVideo").srcObject.id,
		...Object.entries(sparkRTC.userStreamData).map(([key]) => key),
	];
}

function handleCustom(videoLayout) {
	videoLayout.meta.participants = [
		getVideoElement("localVideo").srcObject.id,
		...Object.entries(sparkRTC.userStreamData).map(([key]) => key),
	];
}

function setAndSendLayout() {
	const videoLayout = currentLayout();
	switch (videoLayout.type) {
		case "blue":
		case "red":
		case "yellow":
		case "purple":
			{
				handleParticipants(videoLayout);
			}
			break;
		case "control-freak":
			{
				handleControlFreak(videoLayout);
			}
			break;
		case "custom":
			{
				handleCustom(videoLayout);
			}
			break;
	}

	sparkRTC.socket.send(
		JSON.stringify({
			type: "metadata-set",
			data: JSON.stringify({
				...getMeta(),
				videoLayout: JSON.stringify(videoLayout),
			}),
		})
	);
}

setTimeout(() => {
	if (getMyRole() === "broadcast") {
		setAndSendLayout();
	}
}, [300]);

document.addEventListener("keydown", (event) => {
	if (!sparkRTC || !sparkRTC.socket) {
		return;
	}

	if (getMyRole() === "broadcast") {
		if (event.key === "b") {
			// TODO: filter out the event if not admin
			currentBackgroundIndex =
				(currentBackgroundIndex + 1) % possibleBackgrounds.length;

			sparkRTC.socket.send(
				JSON.stringify({
					type: "metadata-set",
					data: JSON.stringify({
						...getMeta(),
						backgroundLayout: currentBackgroundLayout(),
					}),
				})
			);
		}

		if (event.key === "l") {
			console.log("Setting layout");
			currentLayoutIndex = (currentLayoutIndex + 1) % possibleLayouts.length;

			setAndSendLayout();
		}
	}
});

document.addEventListener("keydown", (event) => {
	if (!sparkRTC || !sparkRTC.socket) {
		return;
	}
	currentBackgroundIndex =
		(currentBackgroundIndex + 1) % possibleBackgrounds.length;

	if (event.key === "b") {
		// TODO: filter out the event

		console.log(currentBackgroundLayout());

		sparkRTC.socket.send(
			JSON.stringify({
				type: "metadata-set",
				data: JSON.stringify({
					...getMeta(),
					backgroundLayout: currentBackgroundLayout(),
				}),
			})
		);
	}
});

window.addEventListener("message", (event) => {
	const parsed = JSON.parse(event.data);
    

	if (parsed.type === "set_layout") {
		console.log("Parsed layout", parsed);
		currentLayoutIndex = possibleLayouts.findIndex(
			(l) => l.type === parsed.data.type
		);

		if (parsed.data.type === "custom") {
			possibleLayouts[currentLayoutIndex].meta.background =
				parsed.data.layout.background;
			possibleLayouts[currentLayoutIndex].meta.placements =
				parsed.data.layout.placements;
		}

		setAndSendLayout();
	}
});

setInterval(() => {
	sparkRTC.socket.send(
		JSON.stringify({
			type: "metadata-get",
			// data: JSON.stringify({"backgroundUrl": `https://upload.logjam.server.group.video${path}`})
		})
	);
	console.log(sparkRTC.metaData);

	const page = document.getElementById("page");

	page.style.backgroundPosition = "center";

	switch (sparkRTC.metaData.backgroundLayout) {
		case "contain":
			page.style.backgroundSize = "contain";
			page.style.backgroundRepeat = "no-repeat";
			break;
		case "cover":
			page.style.backgroundSize = "cover";
			break;
		case "tiled":
		default:
			page.style.backgroundSize = "contain";
			page.style.backgroundRepeat = "repeat";
	}

	page.style.backgroundImage = `url(${sparkRTC.metaData.backgroundUrl})`;
	arrangeVideoContainers();
}, 300);

function setMyName() {
	try {
		myName = localStorage.getItem("logjam_myName");
		document.getElementById("inputName").value = myName;
	} catch (e) {
		console.log(e);
	}
	if (myName === "" || !myName) {
		myName = makeId(20);
		try {
			localStorage.setItem("logjam_myName", myName);
		} catch (e) {
			console.log(e);
		}
	}
}

async function handleClick() {
	let newName = document.getElementById("inputName").value;
	if (newName) {
		myName = newName;
		localStorage.setItem("logjam_myName", myName);
	}
	document.getElementById("page").style.visibility = "visible";
	document.getElementById("getName").style.display = "none";

	await start();

	return false;
}

function handleResize() {
	clearTimeout(window.resizedFinished);
	window.resizedFinished = setTimeout(function () {
		graph.draw(graph.treeData);
		arrangeVideoContainers();
	}, 250);
}

function getMyRole() {
	const queryString = window.location.search;
	const urlParams = new URLSearchParams(queryString);
	return urlParams.get("role") === "broadcast" ? "broadcast" : "audience";
}

function getRoomName() {
	const queryString = window.location.search;
	const urlParams = new URLSearchParams(queryString);
	return urlParams.get("room");
}

function getDebug() {
	const queryString = window.location.search;
	const urlParams = new URLSearchParams(queryString);
	return Boolean(urlParams.get("debug"));
}

function setupSignalingSocket() {
	return sparkRTC.setupSignalingSocket(getWsUrl(), myName, roomName);
}

async function start() {
	await setupSignalingSocket();
	return sparkRTC.start();
}

function onLoad() {
	myRole = getMyRole();
	roomName = getRoomName();
	sparkRTC = createSparkRTC();
	if (!getDebug()) {
		document.getElementById("logs").style.display = "none";
	}

	setMyName();
	graph = new Graph();
	window.onresize = handleResize;
	graph.draw(DATA);

	arrangeVideoContainers();
}

async function onRaiseHand() {
	const stream = await sparkRTC.raiseHand();
	const tagId = "localVideo-" + stream.id;
	if (document.getElementById(tagId)) return;
	const userData = sparkRTC.getStreamDetails(stream.id);
	const video = createVideoElement(tagId, userData, true);
	video.srcObject = stream;
}

function addLog(log) {
	const logs = document.getElementById("logs");
	const p = document.createElement("p");
	p.innerText = log;
	logs.appendChild(p);
}
