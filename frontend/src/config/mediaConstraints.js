export const MEDIA_CONSTRAINTS = (window.constraints = {
    // audio: {'echoCancellation': true},
    // video: {
    //     'width':  {'min': 300},
    //     'height': {'min': 200}
    // }
    audio: false,
    // video: true
    video: {
        width: 360,
        height: 240
    },
});
