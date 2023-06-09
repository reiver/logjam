import { Icon, IconButton, makeDialog, Tooltip } from 'components';
import { html } from 'htm';
import {
    currentUser,
    onStartShareScreen,
    onStopShareScreen,
    raiseHandMaxLimitReached,
    sparkRTC,
    updateUser,
} from '../../pages/meeting.js';

export const Controllers = () => {
    const {
        isHost,
        hasCamera,
        hasMic,
        ableToRaiseHand,
        sharingScreenStream,
        isStreamming,
        isCameraOn,
        isMicrophoneOn,
        isMeetingMuted,
    } = currentUser.value;
    const toggleMuteMeeting = () => {
        updateUser({
            isMeetingMuted: !isMeetingMuted,
        });
    };

    const handleShareScreen = async () => {
        if (!sharingScreenStream) {
            const stream = await sparkRTC.value.startShareScreen();
            onStartShareScreen(stream);
            updateUser({
                sharingScreenStream: stream,
            });
        } else {
            await sparkRTC.value.stopShareScreen(sharingScreenStream);
            onStopShareScreen(sharingScreenStream);
        }
    };
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
    const onRaiseHand = async () => {
        if (isStreamming) {
            updateUser({
                isStreamming: false,
            });
            sparkRTC.value.lowerHand();
        } else {
            updateUser({
                isRaisingHand: true,
            });
            sparkRTC.value.raiseHand();
            makeDialog('info', {
                message: 'Raise hand request has been sent.',
                icon: 'Check',
            });
        }
    };
    const handleReload = () => {
        sparkRTC.value.startProcedure(true);
    };

    return html`<div class="flex gap-5 py-5">
        <${Tooltip} label=${isMeetingMuted ? 'Listen' : 'Deafen'}>
            <${IconButton}
                variant=${isMeetingMuted && 'danger'}
                onClick=${toggleMuteMeeting}
            >
                <${Icon} icon="Volume${isMeetingMuted ? 'Off' : ''}" />
            <//>
        <//>

        ${!isStreamming &&
        html` <${Tooltip} label="Reconnect">
            <${IconButton} onClick=${handleReload}>
                <${Icon} icon="Reconnect" />
            <//>
        <//>`}
        ${isStreamming &&
        isHost &&
        html` <${Tooltip}
            label="${!sharingScreenStream
                ? 'Share Screen'
                : 'Stop Sharing Screen'}"
        >
            <${IconButton}
                variant=${sharingScreenStream && 'danger'}
                onClick=${handleShareScreen}
            >
                <${Icon} icon="Share${sharingScreenStream ? 'Off' : ''}" />
            <//>
        <//>`}
        ${((!raiseHandMaxLimitReached.value &&
            !isStreamming &&
            ableToRaiseHand) ||
            (isStreamming && !isHost && ableToRaiseHand)) &&
        html`<${Tooltip} label=${isStreamming ? 'Lower Hand' : 'Raise Hand'}>
            <${IconButton}
                onClick=${onRaiseHand}
                variant=${isStreamming && 'danger'}
            >
                <${Icon} icon="Hand" /> <//
        ><//>`}
        ${hasCamera &&
        isStreamming &&
        html` <${Tooltip}
            label=${!isCameraOn ? 'Turn Camera On' : 'Turn Camera Off'}
        >
            <${IconButton}
                variant=${!isCameraOn && 'danger'}
                onClick=${toggleCamera}
            >
                <${Icon} icon="Camera${!isCameraOn ? 'Off' : ''}" /> <//
        ><//>`}
        ${hasMic &&
        isStreamming &&
        html` <${Tooltip}
            label=${!isMicrophoneOn
                ? 'Turn Microphone On'
                : 'Turn Microphone Off'}
        >
            <${IconButton}
                variant=${!isMicrophoneOn && 'danger'}
                onClick=${toggleMicrophone}
            >
                <${Icon} icon="Microphone${!isMicrophoneOn ? 'Off' : ''}" />
            <//>
        <//>`}
    </div>`;
};

export default Controllers;
