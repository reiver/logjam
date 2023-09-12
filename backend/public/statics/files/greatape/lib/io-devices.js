export class IODevices {
    devices = null;

    //get all devices
    initDevices = async () => {
        try {
            this.devices = await navigator.mediaDevices.enumerateDevices();
        } catch (error) {
            console.error('Error enumerating devices:', error);
        }
    };

    //list input audio devices
    getAudioInputDevices = () => {
        var audioInputDevices = [];
        this.devices.forEach((device) => {
            if (device.kind === 'audioinput') {
                audioInputDevices.push(device);
            }
        });

        return audioInputDevices;
    };

    //list input video devices
    getVideoInputDevices = () => {
        var videoInputDevices = [];
        this.devices.forEach((device) => {
            if (device.kind === 'videoinput') {
                videoInputDevices.push(device);
            }
        });

        return videoInputDevices;
    };

    //list output video dveices
    getAudioOutputDevices = () => {
        var audioOutputDevices = [];
        this.devices.forEach((device) => {
            if (device.kind === 'audiooutput') {
                audioOutputDevices.push(device);
            }
        });

        return audioOutputDevices;
    };
}

//attach audio sink with video element
function attachSinkId(element, sinkId) {
    if (typeof element.sinkId !== 'undefined') {
        element
            .setSinkId(sinkId)
            .then(() => {
                console.log(`Success, audio output device attached: ${sinkId}`);
            })
            .catch((error) => {
                let errorMessage = error;
                if (error.name === 'SecurityError') {
                    errorMessage = `You need to use HTTPS for selecting audio output device: ${error}`;
                }
                console.error(errorMessage);
                // Jump back to first output device in the list as it's the default.
                audioOutputSelect.selectedIndex = 0;
            });
    } else {
        console.warn('Browser does not support output device selection.');
    }
}