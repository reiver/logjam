import { signal } from '@preact/signals';
import { clsx } from 'clsx';
import { Button, Icon, IconButton, Tooltip } from 'components';
import { html } from 'htm';
import { useEffect, useRef } from 'preact';
import { v4 as uuidv4 } from 'uuid';

import { currentUser, sparkRTC, updateUser } from '../../pages/meeting.js';
import { IODevices } from '../../lib/io-devices.js';

const dialogs = signal([]);
var selectedMic = null;
var selectedSpeaker = null;
var selectedCamera = null;
const builtInLabel = 'Built-in';
const builtInThisDevice = 'Built-in (This Device)';

export const IOSettingsDialog = ({
    onOk,
    onClose,
    message: { message, title },
    okText = 'Sounds Good',
    cancelText = 'Cancel',
    okButtonVariant = 'solid',
    onReject = onClose,
    showButtons = true,
    className,
}) => {
    const selectAudioOutputDevice = async () => {
        const io = new IODevices();
        await io.initDevices();
        const devices = io.getAudioOutputDevices();
        console.log('Devices: ', devices);
        makeIODevicesDialog(
            'io-devices',
            {
                message: 'Please choose your "Audio output":',
                title: 'Audio',
            },
            devices,
            'speaker',
            (device) => {
                setTimeout(() => {
                    if (device && device.label) {
                        console.log('Audio Device', device);
                        var elem = document.getElementById('selectedSpeaker');
                        elem.innerHTML = '';
                        console.log('selectedSpeaker elem: ', elem);
                        if (device.label.toLowerCase().includes('default')) {
                            console.log('Setting inner HTML Built in');
                            elem.innerHTML = builtInLabel;
                        } else {
                            console.log('Setting inner HTML Devuice lable');

                            elem.innerHTML = device.label;
                        }

                        selectedSpeaker = device;
                    }
                }, 100);
            } //on close
        );
    };

    const selectAudioInputDevice = async () => {
        const io = new IODevices();
        await io.initDevices();
        const devices = io.getAudioInputDevices();
        makeIODevicesDialog(
            'io-devices',
            {
                message: 'Please choose your "Audio input"-Microphone:',
                title: 'Microphone',
            },
            devices,
            'microphone',
            (device) => {
                setTimeout(() => {
                    if (device && device.label) {
                        console.log('Mic Device', device);
                        const elem = document.getElementById('selectedMic');
                        console.log('selectedCamera elem: ', elem);
                        elem.innerHTML = '';
                        if (
                            device.label.toLowerCase().includes('default') ||
                            device.label
                                .toLowerCase()
                                .includes('iphone microphone')
                        ) {
                            elem.innerHTML = builtInLabel;
                        } else {
                            elem.innerHTML = device.label;
                        }
                        selectedMic = device;
                    }
                }, 100);
            } //on close
        );
    };

    const selectVideoInputDevice = async () => {
        const io = new IODevices();
        await io.initDevices();
        const devices = io.getVideoInputDevices();
        makeIODevicesDialog(
            'io-devices',
            {
                message: 'Please choose your "Video input":',
                title: 'Video',
            },
            devices,
            'camera',
            (device) => {
                setTimeout(() => {
                    if (device && device.label) {
                        console.log('Video Device', device);
                        const elem = document.getElementById('selectedCamera');
                        console.log('selectedCamera elem: ', elem);
                        elem.innerHTML = '';
                        if (
                            device.label.toLowerCase().includes('default') ||
                            device.label.toLowerCase().includes('(') ||
                            device.label.toLowerCase().includes('front')
                        ) {
                            elem.innerHTML = builtInLabel;
                        } else {
                            elem.innerHTML = device.label;
                        }

                        selectedCamera = device;
                    }
                }, 100);
            } //on close
        );
    };

    const isIphone = () => {
        const userAgent = navigator.userAgent;
        if (userAgent.match(/iPhone|iPad|iPod/i)) {
            return true;
        }
        return false;
    };

    console.log('Resetting..');

    return html` <div class="absolute top-0 left-0 w-full h-full">
        <div class="z-20 absolute w-full h-full bg-black bg-opacity-60" />
        <div
            class=${clsx(
                className,
                'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full w-full sm:max-w-[400px] sm:border dark:border-gray-1 border-gray-0'
            )}
        >
            <div class="flex justify-center items-center p-5 relative">
                <span class="dark:text-white text-black text-bold-12"
                    >${title}</span
                >
                <${Icon}
                    icon="Close"
                    class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer"
                    onClick=${onClose}
                />
            </div>
            <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />
            ${!isIphone() &&
            html`
                <div
                    class="sm:py-4 sm:pt-8 pt:4 py-2 flex pt-4 cursor-pointer"
                    onClick=${selectAudioOutputDevice}
                >
                    <div class="text-left text-bold-12 px-5 flex-1">
                        Audio Output
                    </div>
                    <div
                        id="selectedSpeaker"
                        class="text-right text-bold-12 px-5 flex-1 text-gray-1 cursor-pointer"
                    >
                        ${(() => {
                            const elem =
                                document.getElementById('selectedSpeaker');

                            if (elem && elem.innerHTML !== '') {
                                // innerHTML exists and is not empty or just whitespace
                                console.log(
                                    'Speaker innerHTML exists:',
                                    elem.innerHTML
                                );
                                if (
                                    elem.innerHTML
                                        .toLowerCase()
                                        .includes(builtInLabel.toLowerCase())
                                ) {
                                    elem.innerHTML = '';
                                    return builtInLabel;
                                }
                            } else {
                                // innerHTML does not exist or is empty/whitespace
                                console.log('Speaker innerHTML not exists');

                                console.log(
                                    'selectedSpeaker-default: ',
                                    selectedSpeaker
                                );

                                //return device name
                                if (selectedSpeaker && selectedSpeaker.label) {
                                    const labelLowerCase =
                                        selectedSpeaker.label.toLowerCase();
                                    if (labelLowerCase.includes('default')) {
                                        console.log('returning built in 1');
                                        return builtInLabel;
                                    } else {
                                        return selectedSpeaker.label;
                                    }
                                }
                                console.log('returning built in 2');
                                return builtInLabel;
                            }
                        })()}
                    </div>
                </div>
            `}

            <div
                class="sm:py-4 py-2 flex cursor-pointer"
                onClick=${selectAudioInputDevice}
            >
                <div class="text-left text-bold-12 px-5 flex-1">Microphone</div>
                <div
                    id="selectedMic"
                    class="text-right text-bold-12 px-5 flex-1 text-gray-1"
                >
                    ${(() => {
                        const elem = document.getElementById('selectedMic');

                        if (elem && elem.innerHTML !== '') {
                            // innerHTML exists and is not empty or just whitespace
                            console.log(
                                'mic innerHTML exists:',
                                elem.innerHTML
                            );
                            if (
                                elem.innerHTML
                                    .toLowerCase()
                                    .includes(builtInLabel.toLowerCase())
                            ) {
                                elem.innerHTML = '';
                                return builtInLabel;
                            }
                        } else {
                            // innerHTML does not exist or is empty/whitespace
                            console.log('mic innerHTML not exists');

                            console.log('selectedMic-default: ', selectedMic);

                            //return device name
                            if (selectedMic && selectedMic.label) {
                                const labelLowerCase =
                                    selectedMic.label.toLowerCase();
                                if (
                                    labelLowerCase.includes('default') ||
                                    labelLowerCase.includes('iphone microphone')
                                ) {
                                    return builtInLabel;
                                } else {
                                    return selectedMic.label;
                                }
                            }
                            return builtInLabel;
                        }
                    })()}
                </div>
            </div>

            <div
                class="sm:py-4 sm:pb-8 py-2 flex pb-4 cursor-pointer"
                onClick=${selectVideoInputDevice}
            >
                <div class="text-left text-bold-12 px-5 flex-1">
                    Video Input
                </div>
                <div
                    id="selectedCamera"
                    class="text-right text-bold-12 px-5 flex-1 text-gray-1"
                >
                    ${(() => {
                        const elem = document.getElementById('selectedCamera');

                        if (elem && elem.innerHTML !== '') {
                            // innerHTML exists and is not empty or just whitespace
                            console.log(
                                'Cam innerHTML exists:',
                                elem.innerHTML
                            );
                            if (
                                elem.innerHTML
                                    .toLowerCase()
                                    .includes(builtInLabel.toLowerCase())
                            ) {
                                elem.innerHTML = '';
                                return builtInLabel;
                            }
                        } else {
                            // innerHTML does not exist or is empty/whitespace
                            console.log('Cam innerHTML not exists');

                            console.log(
                                'selectedCamera-default: ',
                                selectedCamera
                            );

                            //return device name
                            if (selectedCamera && selectedCamera.label) {
                                const labelLowerCase =
                                    selectedCamera.label.toLowerCase();
                                if (
                                    labelLowerCase.includes('default') ||
                                    labelLowerCase.includes('(') ||
                                    labelLowerCase.includes('front')
                                ) {
                                    return builtInLabel;
                                } else {
                                    return selectedCamera.label;
                                }
                            }
                            return builtInLabel;
                        }
                    })()}
                </div>
            </div>

            ${showButtons &&
            html`<div class="flex justify-end gap-2 p-5 pt-0">
                <${Button}
                    size="lg"
                    variant="outline"
                    class="w-full flex-grow-1"
                    onClick=${() => {
                        onReject && onReject();
                        onClose();
                    }}
                    >${cancelText}<//
                >
                <${Button}
                    size="lg"
                    variant="${okButtonVariant}"
                    class="w-full flex-grow-1"
                    onClick=${() => {
                        onOk(selectedMic, selectedCamera, selectedSpeaker);
                    }}
                    >${okText}<//
                >
            </div>`}
        </div>
    </div>`;
};

export const IODevicesDialog = ({
    onClose,
    message: { message, title },
    devices,
    deviceType,
    className,
    contentClassName,
}) => {
    let selectedDeviceIndex = -1;
    const handleDeviceClick = (index = -1, vanish = true) => {
        //mark built-in / default devices checked already

        // if (index === -1) {
        //     devices.forEach((elem, _index) => {
        //         if (
        //             deviceType === 'microphone' &&
        //             elem.kind === 'audioinput' &&
        //             elem.label.toLowerCase().includes('default')
        //         ) {
        //             index = _index;
        //         } else if (
        //             deviceType === 'speaker' &&
        //             elem.kind === 'audiooutput' &&
        //             elem.label.toLowerCase().includes('default')
        //         ) {
        //             index = _index;
        //         } else if (
        //             deviceType === 'camera' &&
        //             elem.kind === 'videoinput' &&
        //             (elem.label.toLowerCase().includes('default') ||
        //                 elem.label.toLowerCase().includes('('))
        //         ) {
        //             index = _index;
        //         }
        //     });
        // }

        // Now, the 'index' variable will contain the index of the matching device (or -1 if none found).

        // Check if the clicked device is already selected
        if (selectedDeviceIndex === index) {
            // If it's already selected, deselect it by setting the selectedDeviceIndex to -1
            selectedDeviceIndex = -1;
        } else {
            // If it's not selected, select it by setting the selectedDeviceIndex to the clicked index
            selectedDeviceIndex = index;
        }

        console.log('handleDeviceClick: ', index);
        // Update the radio button's checked attribute based on the selectedDeviceIndex
        const radioInput = document.getElementById(`device${index}`);
        console.log('radioInput: ', radioInput);
        if (radioInput) {
            radioInput.checked = selectedDeviceIndex === index;
            radioInput.style.accentColor = 'black';
            if (vanish) {
                setTimeout(() => {
                    console.log(
                        'Selected Device: ',
                        devices[selectedDeviceIndex]
                    );
                    onClose(devices[selectedDeviceIndex]);
                }, 200);
            }
        }
    };

    //check selected devices and make it selected on radio button too
    if (devices && devices.length > 0) {
        devices.forEach((value, index) => {
            console.log('Device: ', value, ' index: ', index);
            if (value.kind === 'audioinput') {
                if (selectedMic && selectedMic.deviceId === value.deviceId) {
                    console.log('selectedMicis: ', selectedMic);
                    setTimeout(() => {
                        handleDeviceClick(index, false);
                    }, 250);
                }
            } else if (value.kind === 'videoinput') {
                if (
                    selectedCamera &&
                    selectedCamera.deviceId === value.deviceId
                ) {
                    console.log('selectedCamis: ', selectedCamera);

                    setTimeout(() => {
                        handleDeviceClick(index, false);
                    }, 250);
                }
            } else if (value.kind === 'audiooutput') {
                if (
                    selectedSpeaker &&
                    selectedSpeaker.deviceId === value.deviceId
                ) {
                    console.log('selectedSpeakeris: ', selectedSpeaker);

                    setTimeout(() => {
                        handleDeviceClick(index, false);
                    }, 250);
                }
            }
        });
    }

    const isBuiltInDevice = (deviceType, device) => {
        switch (deviceType) {
            case 'speaker':
                if (device.label.toLowerCase().includes('default')) {
                    return true;
                }
                break;
            case 'microphone':
                if (
                    device.label.toLowerCase().includes('default') ||
                    device.label.toLowerCase().includes('iphone microphone')
                ) {
                    return true;
                }
                break;
            case 'camera':
                if (
                    device.label.toLowerCase().includes('(') ||
                    device.label.toLowerCase().includes('front')
                ) {
                    return true;
                }
                break;
            default:
                // Handle other cases if needed
                break;
        }

        return false;
    };

    return html` <div class="absolute top-0 left-0 w-full h-full">
        <div class="z-20 absolute w-full h-full bg-black bg-opacity-60" />
        <div
            class=${clsx(
                className,
                'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full w-full sm:max-w-[400px] sm:border dark:border-gray-1 border-gray-0'
            )}
        >
            <div class="flex justify-center items-center p-5 relative">
                <span class="dark:text-white text-black text-bold-12"
                    >${title}</span
                >
                <${Icon}
                    icon="Close"
                    class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer"
                    onClick=${onClose}
                />
            </div>
            <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />
            <div
                class=${clsx(
                    contentClassName,
                    'text-left text-bold-12 sm:pt-8 pt-5 p-5'
                )}
                dangerouslySetInnerHTML=${{ __html: message }}
            ></div>

            <form>
                <div class="sm:pb-4 pb-2">
                    ${devices.map(
                        (device, index) => html`
                            <div
                                class="sm:py-4 py-2 flex items-center cursor-pointer"
                                onClick=${() => handleDeviceClick(index)}
                            >
                                <${Icon}
                                    icon="${(() => {
                                        //return defualt builtin

                                        if (
                                            isBuiltInDevice(deviceType, device)
                                        ) {
                                            return 'Smartphone';
                                        }

                                        //select other then default
                                        if (deviceType === 'microphone') {
                                            return 'MicrophoneLight';
                                        } else if (deviceType === 'camera') {
                                            return 'CameraLight';
                                        } else if (deviceType === 'speaker') {
                                            return 'Headphone';
                                        }
                                    })()}"
                                    class="ml-5"
                                    width="20px"
                                    height="20px"
                                />
                                <div class="text-left px-2 text-bold-12 flex-1">
                                    ${(() => {
                                        var deviceName = '';

                                        if (
                                            isBuiltInDevice(deviceType, device)
                                        ) {
                                            deviceName = builtInThisDevice;
                                        } else {
                                            deviceName = device.label;
                                        }

                                        return deviceName;
                                    })()}
                                </div>
                                <label class="flex items-right px-5 flex-0">
                                    <input
                                        type="radio"
                                        name="devices"
                                        id=${`device${index}`}
                                        device=${device}
                                    />
                                </label>
                            </div>
                        `
                    )}
                </div>
            </form>
        </div>
    </div>`;
};

export const PreviewDialog = ({
    onOk,
    onClose,
    videoStream,
    message: { message, title },
    okText = 'Sounds Good',
    cancelText = 'Cancel',
    okButtonVariant = 'solid',
    onReject = onClose,
    showButtons = true,
    className,
    contentClassName,
}) => {
    selectedCamera = null;
    selectedMic = null;
    selectedSpeaker = null;
    const videoRef = useRef();
    const { hasCamera, hasMic, isCameraOn, isMicrophoneOn } = currentUser.value;

    useEffect(() => {
        videoRef.current.srcObject = videoStream;
    }, [videoStream]);

    useEffect(() => {
        videoRef.current.playsInline = true;

        const isMobile =
            window.parent.outerWidth <= 400 && window.parent.outerHeight <= 850;
        if (isMobile) {
            //set style
            videoRef.current.style =
                'width: 100%; height: 56.25vw; max-height: 100%;';
        }
    }, []);

    const toggleCamera = () => {
        sparkRTC.value.disableVideo(!isCameraOn);
        updateUser({
            isCameraOn: !isCameraOn,
        });
    };

    const toggleMicrophone = () => {
        sparkRTC.value.disableAudio(!isMicrophoneOn);
        updateUser({
            isMicrophoneOn: !isMicrophoneOn,
        });
    };

    const openDeviceSettings = () => {
        makeIOSettingsDialog(
            'io-settings',
            {
                message: '',
                title: 'Input and Output Settings',
            },
            async (mic, cam, speaker) => {
                console.log('mic: ', mic, 'cam: ', cam, 'speaker: ', speaker);

                //now change the Audio, Video and Speaker devices
                const stream = await sparkRTC.value.changeIODevices(
                    mic,
                    cam,
                    speaker
                );

                console.log('New Stream: ', stream);
                videoStream = stream;
                if (videoRef.current) {
                    videoRef.current.srcObject = stream;
                } else {
                    console.log('No video ref');
                }
            }, //ok
            () => {} //close
        );
    };

    return html` <div class="absolute top-0 left-0 w-full h-full">
        <div class="z-10 absolute w-full h-full bg-black bg-opacity-60" />
        <div
            class=${clsx(
                className,
                'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full w-full sm:max-w-[40%] sm:border dark:border-gray-1 border-gray-0'
            )}
        >
            <div class="flex justify-center items-center p-5 relative">
                <span class="dark:text-white text-black text-bold-12"
                    >${title}</span
                >
                <${Icon}
                    icon="Close"
                    class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer"
                    onClick=${onClose}
                />
            </div>
            <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />
            <div
                class=${clsx(
                    contentClassName,
                    'text-left text-bold-12 sm:py-8 py-5 p-5'
                )}
                dangerouslySetInnerHTML=${{ __html: message }}
            ></div>
            <div class="px-5 relative">
                <video
                    ref=${videoRef}
                    autoplay
                    playsinline
                    muted="true"
                    className="aspect-video object-cover rounded-lg"
                />
                <div
                    class=${clsx(
                        'h-[48px] absolute top-1 right-6 gap-2 flex justify-center items-center'
                    )}
                >
                    ${!isMicrophoneOn &&
                    html` <div className="pr-2">
                        <${Icon}
                            icon="MicrophoneOff"
                            width="20px"
                            height="20px"
                        />
                    </div>`}
                </div>
            </div>
            <div class="py-4 px-5 flex gap-3 justify-center">
                ${hasCamera &&
                html` <${Tooltip}
                    key="Camera${!isCameraOn ? 'Off' : ''}"
                    label=${!isCameraOn ? 'Turn Camera On' : 'Turn Camera Off'}
                >
                    <${IconButton}
                        variant=${!isCameraOn && 'danger'}
                        onClick=${toggleCamera}
                    >
                        <${Icon} icon="Camera${!isCameraOn ? 'Off' : ''}" /> <//
                ><//>`}
                ${hasMic &&
                html`
                    <${Tooltip}
                        key="Microphone${!isMicrophoneOn ? 'Off' : ''}"
                        label=${!isMicrophoneOn
                            ? 'Turn Microphone On'
                            : 'Turn Microphone Off'}
                    >
                        <${IconButton}
                            variant=${!isMicrophoneOn && 'danger'}
                            onClick=${toggleMicrophone}
                        >
                            <${Icon}
                                icon="Microphone${!isMicrophoneOn ? 'Off' : ''}"
                            />
                        <//>
                    <//>
                `}
                ${hasCamera &&
                hasMic &&
                html`
                    <${Tooltip} key="Settings" label="Settings">
                        <${IconButton} onClick=${openDeviceSettings}>
                            <${Icon} icon="Settings" />
                        <//>
                    <//>
                `}
            </div>

            ${showButtons &&
            html`<div class="flex justify-end gap-2 p-5 pt-0">
                <${Button}
                    size="lg"
                    variant="outline"
                    class="w-full flex-grow-1"
                    onClick=${() => {
                        onReject && onReject();
                        onClose();
                    }}
                    >${cancelText}<//
                >
                <${Button}
                    size="lg"
                    variant="${okButtonVariant}"
                    class="w-full flex-grow-1"
                    onClick=${onOk}
                    >${okText}<//
                >
            </div>`}
        </div>
    </div>`;
};

export const ConfirmDialog = ({
    onOk,
    onClose,
    message: { message, title },
    okText = 'Accept',
    cancelText = 'Reject',
    okButtonVariant = 'solid',
    onReject = onClose,
    showButtons = true,
    className,
    contentClassName,
}) => {
    return html` <div class="absolute top-0 left-0 w-full h-full">
        <div
            class="z-10 absolute w-full h-full bg-black bg-opacity-60"
            onClick=${onClose}
        />
        <div
            class=${clsx(
                className,
                'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full w-full sm:max-w-[400px] sm:border dark:border-gray-1 border-gray-0'
            )}
        >
            <div class="flex justify-center items-center p-5 relative">
                <span class="dark:text-white text-black text-bold-12"
                    >${title}</span
                >
                <${Icon}
                    icon="Close"
                    class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer"
                    onClick=${onClose}
                />
            </div>
            <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />
            <div
                class=${clsx(
                    contentClassName,
                    'text-left text-bold-12 sm:py-8 py-5 p-5'
                )}
                dangerouslySetInnerHTML=${{ __html: message }}
            ></div>
            ${showButtons &&
            html`<div class="flex justify-end gap-2 p-5 pt-0">
                <${Button}
                    size="lg"
                    variant="outline"
                    class="w-full flex-grow-1"
                    onClick=${() => {
                        onReject && onReject();
                        onClose();
                    }}
                    >${cancelText}<//
                >
                <${Button}
                    size="lg"
                    variant="${okButtonVariant}"
                    class="w-full flex-grow-1"
                    onClick=${onOk}
                    >${okText}<//
                >
            </div>`}
        </div>
    </div>`;
};

export const InfoDialog = ({
    onOk,
    onClose,
    message: { message, icon, variant },
    pointer,
}) => {
    return html`<div
        class=${clsx(
            'select-none py-4 px-6 flex justify-between items-center text-medium-12 min-w-full sm:min-w-[350px] rounded-md',
            {
                'cursor-pointer': pointer,
                'bg-red-distructive text-white-f-9': variant === 'danger',
                'dark:bg-white-f-9 dark:text-gray-3 bg-gray-3 text-white-f-9 border dark:border-gray-1 border-gray-0':
                    !variant,
            }
        )}
        onClick=${onClose}
    >
        <div
            class="text-left"
            dangerouslySetInnerHTML=${{ __html: message }}
        ></div>
        <${Icon} icon=${icon} width="24" height="24" />
    </div>`;
};

export const DialogPool = () => {
    return html`<div
            className="absolute right-0 left-0 md:left-[unset] md:right-10 bottom-[5.5rem] flex flex-col justify-end gap-2 px-4 sm:px-0"
        >
            ${Object.values(dialogs.value).map((dialog) => {
                if (dialog.type === 'info')
                    return html`<${InfoDialog} ...${dialog} />`;
            })}
        </div>

        ${Object.values(dialogs.value).map((dialog) => {
            if (dialog.type === 'confirm')
                return html`<${ConfirmDialog} ...${dialog} />`;
            else if (dialog.type === 'preview')
                return html`<${PreviewDialog} ...${dialog} />`;
            else if (dialog.type === 'io-settings')
                return html`<${IOSettingsDialog} ...${dialog} />`;
            else if (dialog.type === 'io-devices')
                return html`<${IODevicesDialog} ...${dialog} />`;
        })}`;
};

export const makeDialog = (type, message, onOk, onClose, options = {}) => {
    const id = uuidv4();
    const destroy = () => {
        const dialogsTmp = { ...dialogs.value };
        delete dialogsTmp[id];
        dialogs.value = dialogsTmp;
    };
    if (type !== 'confirm') {
        setTimeout(destroy, 4000);
    }
    dialogs.value = {
        ...dialogs.value,
        [id]: {
            id,
            type,
            message,
            pointer: !!onClose,
            onOk: () => {
                onOk && onOk();
                destroy();
            },
            onClose: () => {
                onClose && onClose();
                destroy();
            },
            ...options,
        },
    };
};

export const destroyDialog = (id) => {
    const dialogsTmp = { ...dialogs.value };
    delete dialogsTmp[id];
    dialogs.value = dialogsTmp;
};

export const makePreviewDialog = (
    type,
    videoStream,
    message,
    onOk,
    onClose,
    options = {}
) => {
    const id = uuidv4();
    const destroy = () => {
        const dialogsTmp = { ...dialogs.value };
        delete dialogsTmp[id];
        dialogs.value = dialogsTmp;
    };

    dialogs.value = {
        ...dialogs.value,
        [id]: {
            id,
            type,
            videoStream,
            message,
            pointer: !!onClose,
            onOk: () => {
                onOk && onOk();
                destroy();
            },
            onClose: () => {
                onClose && onClose();
                destroy();
            },
            ...options,
        },
    };

    return id;
};

export const makeIODevicesDialog = (
    type,
    message,
    devices,
    deviceType,
    onClose,
    options = {}
) => {
    const id = uuidv4();
    const destroy = () => {
        const dialogsTmp = { ...dialogs.value };
        delete dialogsTmp[id];
        dialogs.value = dialogsTmp;
    };

    dialogs.value = {
        ...dialogs.value,
        [id]: {
            id,
            type,
            message,
            devices,
            deviceType,
            onOk: () => {
                destroy();
            },
            onClose: (device) => {
                onClose && onClose(device);
                destroy();
            },
            ...options,
        },
    };

    return id;
};

export const makeIOSettingsDialog = (
    type,
    message,
    onOk,
    onClose,
    options = {}
) => {
    const id = uuidv4();
    const destroy = () => {
        const dialogsTmp = { ...dialogs.value };
        delete dialogsTmp[id];
        dialogs.value = dialogsTmp;
    };

    dialogs.value = {
        ...dialogs.value,
        [id]: {
            id,
            type,
            message,
            pointer: !!onClose,
            onOk: (mic, cam, speaker) => {
                onOk && onOk(mic, cam, speaker);
                destroy();
            },
            onClose: () => {
                // selectedCamera = null;
                // selectedMic = null;
                // selectedSpeaker = null;
                onClose && onClose();
                destroy();
            },
            ...options,
        },
    };

    return id;
};

export const ToastProvider = () => {
    return html`<div id="toast-provider">
        <${DialogPool} />
    </div>`;
};
