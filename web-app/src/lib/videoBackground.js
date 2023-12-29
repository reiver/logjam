export class VideoBackground {
  setBackVideoBackground = async (
    image,
    videoTrack,
    audioTrack,
    blur = false
  ) => {
    this.bgImage.src = image;

    // instance of SelfieSegmentation object
    const selfieSegmentation = new SelfieSegmentation({
      locateFile: (file) =>
        `https://cdn.jsdelivr.net/npm/@mediapipe/selfie_segmentation/${file}`,
    });

    // set the model and mode
    selfieSegmentation.setOptions({
      modelSelection: 1,
      selfieMode: true,
    });

    // set the callback function for when it finishes segmenting
    selfieSegmentation.onResults(this.onResults.bind(this));

    // definition of track processor and generator
    const trackProcessor = new MediaStreamTrackProcessor({ track: videoTrack });
    const trackGenerator = new MediaStreamTrackGenerator({ kind: "video" });

    const _canv = this.canvas;
    const _this = this; //save ref to this

    // transform function
    const transformer = new TransformStream({
      async transform(videoFrame, controller) {
        const timestamp = videoFrame.timestamp; //save the time stamp
        var newFrame = null;
        var blurFrame = null;

        if (blur === true) {
          // Apply a blur effect to the video frame
          blurFrame = await _this.blurBackground(videoFrame, 5); // Adjust blur amount as needed
          // console.log("Blured Video Frame: ", blurFrame);
        } else {
          // we send the video frame to MediaPipe
          videoFrame.width = videoFrame.displayWidth;
          videoFrame.height = videoFrame.displayHeight;
          await selfieSegmentation.send({ image: videoFrame });
        }
        if (blurFrame != null) {
          newFrame = new VideoFrame(blurFrame, { timestamp });
        } else {
          newFrame = new VideoFrame(_canv, { timestamp });
        }
        // we close the current videoFrame and queue the new one
        videoFrame.close();
        controller.enqueue(newFrame);
      },
    });

    // we pipe the stream through the transform function
    trackProcessor.readable
      .pipeThrough(transformer)
      .pipeTo(trackGenerator.writable);

    // add the new mediastream to video element
    const processedStream = new MediaStream();
    processedStream.addTrack(audioTrack);
    processedStream.addTrack(trackGenerator);

    return processedStream;
  };

  blurBackground = async (videoFrame, blurAmount) => {
    const offscreenCanvas = new OffscreenCanvas(
      videoFrame.displayWidth,
      videoFrame.displayHeight
    );
    const ctx = offscreenCanvas.getContext("2d");
    ctx.drawImage(videoFrame, 0, 0);

    // Apply blur effect
    ctx.filter = `blur(${blurAmount}px)`;
    ctx.drawImage(offscreenCanvas, 0, 0);

    // // Reset the filter to remove blur for subsequent drawings
    // ctx.filter = "none";

    return offscreenCanvas.transferToImageBitmap();
  };

  onResults(results) {
    if (this.ctx != undefined) {
      this.ctx.save();
      this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
      this.ctx.drawImage(
        results.segmentationMask,
        0,
        0,
        this.canvas.width,
        this.canvas.height
      );

      this.ctx.globalCompositeOperation = "source-out";
      const pat = this.ctx.createPattern(this.bgImage, "no-repeat");
      this.ctx.fillStyle = pat;
      this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

      // Only overwrite missing pixels.
      this.ctx.globalCompositeOperation = "destination-atop";
      this.ctx.drawImage(
        results.image,
        0,
        0,
        this.canvas.width,
        this.canvas.height
      );

      this.ctx.restore();
    } else {
      console.log("ctx is undefined");
    }
  }

  constructor() {
    this._height = 1920;
    this._width = 1200;
    // the background image
    this.bgImage = new Image(this._height, this._width);

    // an OffscreenCanvas that combines background and human pixels
    this.canvas = new OffscreenCanvas(this._height, this._width);
    this.ctx = this.canvas.getContext("2d");
    console.log("ctx: ", this.ctx);
  }
}
