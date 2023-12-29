export class VideoBackground {
  setBackVideoBackground = async (image, videoStream, blur = false) => {
    this.bgImage.src = image;
    this.blur = blur;

    const videoTrack = videoStream.getVideoTracks()[0];
    const audioTrack = videoStream.getAudioTracks()[0];

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

    const _this = this; //save ref to this

    // transform function
    const transformer = new TransformStream({
      async transform(videoFrame, controller) {
        const timestamp = videoFrame.timestamp; //save the time stamp
        var newFrame = null;

        // we send the video frame to MediaPipe
        videoFrame.width = videoFrame.displayWidth;
        videoFrame.height = videoFrame.displayHeight;
        await selfieSegmentation.send({ image: videoFrame });
        newFrame = new VideoFrame(_this.canvas, { timestamp });

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

      //set blur filter
      if (this.blur) {
        this.ctx.filter = "blur(0)";
      }

      this.ctx.drawImage(
        results.segmentationMask,
        0,
        0,
        this.canvas.width,
        this.canvas.height
      );

      //if blur enabled, blur background
      if (this.blur) {
        this.ctx.globalCompositeOperation = "source-in";
        this.ctx.drawImage(
          results.image,
          0,
          0,
          this.canvas.width,
          this.canvas.height
        );
      } else {
        this.ctx.globalCompositeOperation = "source-out";
        const pat = this.ctx.createPattern(this.bgImage, "no-repeat");
        this.ctx.fillStyle = pat;
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
      }

      // Only overwrite missing pixels.
      this.ctx.globalCompositeOperation = "destination-atop";

      //set blur amount the higher value, more blur is
      if (this.blur) {
        this.ctx.filter = "blur(70px)";
      }
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
    this.blur = false;
  }
}
