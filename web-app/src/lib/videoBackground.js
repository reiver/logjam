export class VideoBackground {
  setBackVideoBackground = async (image, videoTrack, audioTrack) => {
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
    // transform function
    const transformer = new TransformStream({
      async transform(videoFrame, controller) {
        // we send the video frame to MediaPipe
        videoFrame.width = videoFrame.displayWidth;
        videoFrame.height = videoFrame.displayHeight;
        await selfieSegmentation.send({ image: videoFrame });

        // we create a new videoFrame
        const timestamp = videoFrame.timestamp;
        const newFrame = new VideoFrame(_canv, { timestamp });

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
