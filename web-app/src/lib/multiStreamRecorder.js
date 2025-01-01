import MultiStreamsMixer from 'multistreamsmixer';


class MultiStreamRecorder {
    constructor() {
        this.resetRecording()
        this.videoType = "video/mp4"; // Ensure browser compatibility
    }

    resetRecording() {
        this.streams = [];
        this.recordedBlobs = [];
        this.mediaRecorder = null;
        this.canvas = null;
        this.canvasStream = null;
        this.videos = [];
        this.roomname = null;
        this.audioMixer = null;
        this.audioMixerStreams = [];
    }

    addStreams(streams) {
        if (!this.canvas) {
            return
        }
        console.log("Add Stream called: ", streams)

        // Remove duplicates, keeping the last occurrence
        const seenStreamIds = new Set(this.streams.map((s) => s.id));
        streams.forEach((stream) => {
            if (!seenStreamIds.has(stream.id)) {
                seenStreamIds.add(stream.id);
                this.streams.push(stream);
                const video = document.createElement("video");
                video.srcObject = new MediaStream(stream.getVideoTracks());
                video.muted = true; // Mute to avoid feedback
                video.play();
                this.videos.push(video);
            }
        });

        console.log("Total Num of Streams After Addition: ", this.streams);

        /**
         * Dynamic Mixing of Audio
         */

        //initalize audio mixer
        // if (this.audioMixer == null) {
        //     console.log("Creating new Mixer")
        //     this.audioMixer = new MultiStreamsMixer(this.streams);
        // } else {
        //     console.log("Adding stream to existing mixer")
        //     // Add new streams only to audioMixerStreams
        //     this.streams.forEach(stream => {
        //         if (!this.audioMixerStreams.includes(stream)) {
        //             this.audioMixerStreams.push(stream);
        //             this.audioMixer.appendStreams([stream])
        //         }
        //     });
        // }
    }

    removeStream(streamId) {
        if (!this.canvas) {
            return
        }

        const streamIndex = this.streams.findIndex((stream) => stream.id === streamId);
        if (streamIndex !== -1) {
            this.streams.splice(streamIndex, 1);
            this.videos.splice(streamIndex, 1);
        }
        console.log("Total Num of Streams After Removal: ", this.streams);
    }


    // Initialize canvas if not already done
    initCanvas() {
        if (!this.canvas) {
            const canvas = document.createElement("canvas");
            document.body.appendChild(canvas); // Optionally attach to the DOM
            canvas.width = 1280; // Set canvas width
            canvas.height = 720; // Set canvas height
            this.canvas = canvas;
        }
        return this.canvas;
    }

    // Start Recording
    async startRecording(roomname, streams) {
        this.roomname = roomname;
        if (!this.canvas) {
            this.initCanvas();
        }
        const ctx = this.canvas.getContext("2d");

        this.addStreams(streams)

        const drawFrames = () => {
            ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

            const numStreams = this.videos.length;
            const shareScreenIndex = this.streams.findIndex(stream => stream.isShareScreen);

            const shareStream = this.videos[shareScreenIndex];
            const otherStreams = this.videos.filter((_, index) => index !== shareScreenIndex);

            // Utility function to draw rounded rectangles
            const drawRoundedRect = (x, y, width, height, radius) => {
                ctx.beginPath();
                ctx.moveTo(x + radius, y);
                ctx.lineTo(x + width - radius, y);
                ctx.quadraticCurveTo(x + width, y, x + width, y + radius);
                ctx.lineTo(x + width, y + height - radius);
                ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height);
                ctx.lineTo(x + radius, y + height);
                ctx.quadraticCurveTo(x, y + height, x, y + height - radius);
                ctx.lineTo(x, y + radius);
                ctx.quadraticCurveTo(x, y, x + radius, y);
                ctx.closePath();
            };

            if (shareScreenIndex !== -1) {
                // Screen share is available

                // Draw screen share stream (70% width, full height on the left)
                const shareWidth = this.canvas.width * 0.7;
                drawRoundedRect(0, 0, shareWidth, this.canvas.height, 10);
                ctx.fillStyle = "black";
                ctx.fill();
                ctx.drawImage(shareStream, 0, 0, shareWidth, this.canvas.height);

                // Draw other streams (30% width, stacked vertically on the right)
                const spacing = 10; // Space between video containers
                const streamHeight = (this.canvas.height - spacing * (otherStreams.length - 1)) / otherStreams.length;
                otherStreams.forEach((video, index) => {
                    let yPosition;
                    let height;

                    if (otherStreams.length === 1) {
                        // Center the single video vertically on the right side
                        height = this.canvas.height / 2; // Half the canvas height
                        yPosition = (this.canvas.height - height) / 2; // Center vertically
                    } else {
                        // Stack videos vertically with spacing
                        height = streamHeight;
                        yPosition = index * (streamHeight + spacing);
                    }

                    const xPosition = this.canvas.width * 0.7;

                    // Draw rounded container
                    drawRoundedRect(xPosition, yPosition, this.canvas.width * 0.3, height, 10);
                    ctx.fillStyle = "black";
                    ctx.fill();

                    // Draw video
                    ctx.drawImage(video, xPosition, yPosition, this.canvas.width * 0.3, height);
                });

            }
            else {
                // Default layout if no screen share is present
                const spacing = 10; // Space between video containers
                if (numStreams === 1) {
                    // Full screen for 1 stream
                    drawRoundedRect(0, 0, this.canvas.width, this.canvas.height, 10);
                    ctx.fillStyle = "black";
                    ctx.fill();
                    ctx.drawImage(this.videos[0], 0, 0, this.canvas.width, this.canvas.height);
                } else if (numStreams === 2) {
                    // Split screen for 2 streams
                    const halfWidth = this.canvas.width / 2;
                    drawRoundedRect(0, 0, halfWidth - spacing / 2, this.canvas.height, 10);
                    drawRoundedRect(halfWidth + spacing / 2, 0, halfWidth - spacing / 2, this.canvas.height, 10);
                    ctx.fillStyle = "black";
                    ctx.fill();
                    ctx.drawImage(this.videos[0], 0, 0, halfWidth - spacing / 2, this.canvas.height);
                    ctx.drawImage(this.videos[1], halfWidth + spacing / 2, 0, halfWidth - spacing / 2, this.canvas.height);
                } else if (numStreams === 3) {
                    // Top row with 2 streams, bottom row with 1 centered
                    const halfHeight = this.canvas.height / 2;
                    drawRoundedRect(0, 0, this.canvas.width / 2 - spacing / 2, halfHeight, 10);
                    drawRoundedRect(this.canvas.width / 2 + spacing / 2, 0, this.canvas.width / 2 - spacing / 2, halfHeight, 10);
                    drawRoundedRect(this.canvas.width / 4, halfHeight + spacing, this.canvas.width / 2, halfHeight - spacing, 10);
                    ctx.fillStyle = "black";
                    ctx.fill();
                    ctx.drawImage(this.videos[0], 0, 0, this.canvas.width / 2 - spacing / 2, halfHeight);
                    ctx.drawImage(this.videos[1], this.canvas.width / 2 + spacing / 2, 0, this.canvas.width / 2 - spacing / 2, halfHeight);
                    ctx.drawImage(this.videos[2], this.canvas.width / 4, halfHeight + spacing, this.canvas.width / 2, halfHeight - spacing);
                } else if (numStreams === 4) {
                    // 2 rows of 2 streams each
                    const halfWidth = this.canvas.width / 2;
                    const halfHeight = this.canvas.height / 2;
                    drawRoundedRect(0, 0, halfWidth - spacing / 2, halfHeight, 10);
                    drawRoundedRect(halfWidth + spacing / 2, 0, halfWidth - spacing / 2, halfHeight, 10);
                    drawRoundedRect(0, halfHeight + spacing, halfWidth - spacing / 2, halfHeight - spacing, 10);
                    drawRoundedRect(halfWidth + spacing / 2, halfHeight + spacing, halfWidth - spacing / 2, halfHeight - spacing, 10);
                    ctx.fillStyle = "black";
                    ctx.fill();
                    ctx.drawImage(this.videos[0], 0, 0, halfWidth - spacing / 2, halfHeight);
                    ctx.drawImage(this.videos[1], halfWidth + spacing / 2, 0, halfWidth - spacing / 2, halfHeight);
                    ctx.drawImage(this.videos[2], 0, halfHeight + spacing, halfWidth - spacing / 2, halfHeight - spacing);
                    ctx.drawImage(this.videos[3], halfWidth + spacing / 2, halfHeight + spacing, halfWidth - spacing / 2, halfHeight - spacing);
                }
            }

            // Draw stream names on the canvas
            this.videos.forEach((video, index) => {
                const streamName = `${this.streams[index]?.name || `Stream ${index + 1}`}${this.streams[index]?.role === "broadcast" ? " (Host)" : ""}`;
                let x, y;

                if (shareScreenIndex !== -1 && index === shareScreenIndex) {
                    // Position for the screen share stream
                    x = 10; // Left side margin
                    y = 30; // Top margin
                } else if (shareScreenIndex !== -1) {
                    // Position for other streams on the right
                    const otherIndex = otherStreams.indexOf(video);
                    if (otherStreams.length === 1) {
                        // Adjust position for a single stream on the right
                        x = this.canvas.width * 0.7 + 10;
                        y = (this.canvas.height - this.canvas.height / 2) / 2 + 30; // Center vertically + top margin
                    } else {
                        // Stacked streams
                        x = this.canvas.width * 0.7 + 10;
                        y = otherIndex * (this.canvas.height / otherStreams.length) + 30;
                    }
                } else {
                    // Default position calculation
                    x = index % 2 === 0 ? 10 : this.canvas.width / 2 + 10; // Position X
                    y = Math.floor(index / 2) * (this.canvas.height / 2) + 20; // Position Y
                }

                ctx.font = "20px Arial";
                ctx.fillStyle = "white";
                ctx.fillText(streamName, x, y);
            });

            requestAnimationFrame(drawFrames); // Continuously update canvas
        };

        drawFrames();

        // Combine canvas video and audio tracks
        const canvasStream = this.canvas.captureStream(30); // Capture 30 FPS video from canvas
        this.canvasStream = canvasStream;

        //static recoridng of Audio streams (need all stream to be available)
        this.audioMixer = new MultiStreamsMixer(this.streams);

        console.log("Audio Mixer: ", this.audioMixer)
        console.log("Audio Mixer Streams: ", this.audioMixer.getMixedStream().getAudioTracks())

        // Combine video from the canvas and the combined audio stream
        const combinedStream = new MediaStream([...canvasStream.getVideoTracks(), ...this.audioMixer.getMixedStream().getAudioTracks()]);

        console.log("Combined stream: ", combinedStream)
        console.log("Combined stream Videos: ", combinedStream.getVideoTracks())
        console.log("Combined stream Audios: ", combinedStream.getAudioTracks())


        this.mediaRecorder = new MediaRecorder(combinedStream, {
            mimeType: this.videoType,
        });

        this.mediaRecorder.ondataavailable = (event) => {
            if (event.data && event.data.size > 0) {
                this.recordedBlobs.push(event.data);
            }
        };

        this.mediaRecorder.start();
        console.log("Recording started");
    }

    // Stop Recording
    stopRecording() {
        console.log("Stop Recording");
        this.mediaRecorder.onstop = () => {
            const blob = new Blob(this.recordedBlobs, { type: this.videoType });
            const url = URL.createObjectURL(blob);

            const currentTime = new Date();

            // Create download link
            const a = document.createElement("a");
            a.href = url;
            a.download = `${this.roomname}__Recording__${currentTime}__.mp4`;
            a.click();

            console.log("Recording saved");

            if (this.audioMixer != null) {
                this.audioMixer.releaseStreams()
            }
            this.resetRecording();
        };

        this.mediaRecorder.stop();
    }
}

export default MultiStreamRecorder;
