class MultiStreamRecorder {
    constructor() {
        this.resetRecording()
        this.videoType = "video/webm"; // Ensure browser compatibility
    }

    resetRecording() {
        this.streams = [];
        this.recordedBlobs = [];
        this.mediaRecorder = null;
        this.canvas = null;
        this.canvasStream = null;
        this.videos = [];
    }

    addStreams(streams) {
        // Array of media streams
        console.log("Total Num of Streams Before Deduplication: ", streams);

        // Remove duplicates, keeping the last occurrence
        const uniqueStreams = [];
        const seenStreamIds = new Set();

        // Iterate in reverse to ensure the last occurrence is kept
        for (let i = streams.length - 1; i >= 0; i--) {
            const stream = streams[i];
            if (!seenStreamIds.has(stream.id)) {
                seenStreamIds.add(stream.id);
                uniqueStreams.unshift(stream); // Add to the front of the array
            }
        }

        this.streams = uniqueStreams;
        console.log("Total Num of Streams After Deduplication: ", this.streams);
    }

    // Initialize Canvas
    initCanvas() {
        const canvas = document.createElement("canvas");
        document.body.appendChild(canvas); // Optionally attach it to the DOM for debugging
        canvas.width = 1280; // Set canvas width
        canvas.height = 720; // Set canvas height
        this.canvas = canvas;

        return canvas;
    }

    // Start Recording
    async startRecording() {
        if (!this.canvas) {
            this.initCanvas();
        }
        const ctx = this.canvas.getContext("2d");

        // Create video elements for each stream
        this.videos = this.streams.map((stream) => {
            const video = document.createElement("video");
            video.srcObject = new MediaStream(stream.getVideoTracks());
            video.muted = true; // Mute video to prevent echo
            video.play();
            return video;
        });

        const drawFrames = () => {
            ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

            this.videos.forEach((video, index) => {
                const x = (index % 2) * (this.canvas.width / 2); // X position
                const y = Math.floor(index / 2) * (this.canvas.height / 2); // Y position
                const width = this.canvas.width / 2;
                const height = this.canvas.height / 2;

                // Draw video on canvas
                if (!video.paused && !video.ended) {
                    ctx.drawImage(video, x, y, width, height);
                }

                // Draw stream name on top-left corner of the video frame
                const streamName = this.streams[index]?.name || `Stream ${index + 1}`;
                ctx.font = "20px Arial";
                ctx.fillStyle = "white";
                ctx.fillText(streamName, x + 10, y + 20); // Adjust position if needed
            });

            requestAnimationFrame(drawFrames); // Continuously update canvas
        };
        drawFrames();

        // Combine canvas video and audio tracks
        const canvasStream = this.canvas.captureStream(30); // Capture 30 FPS video from canvas
        this.canvasStream = canvasStream;

        const audioTracks = this.streams
            .flatMap((stream) => stream.getAudioTracks())
            .filter(Boolean); // Ensure no empty tracks

        const combinedStream = new MediaStream([...canvasStream.getVideoTracks(), ...audioTracks]);

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

            // Create download link
            const a = document.createElement("a");
            a.href = url;
            a.download = "multi_stream_recording.webm";
            a.click();

            console.log("Recording saved");

            this.resetRecording();
        };

        this.mediaRecorder.stop();
    }
}

export default MultiStreamRecorder;
