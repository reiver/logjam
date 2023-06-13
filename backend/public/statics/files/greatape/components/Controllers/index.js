import { signal } from '@preact/signals';
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

export const isMoreOptionsOpen = signal(false);
export const toggleMoreOptions = () =>
    (isMoreOptionsOpen.value = !isMoreOptionsOpen.value);
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
                ableToRaiseHand: true,
            });
            sparkRTC.value.lowerHand();
        } else {
            updateUser({
                isRaisingHand: true,
                ableToRaiseHand: false,
            });
            sparkRTC.value.raiseHand();
            makeDialog('info', {
                message: 'Raise hand request has been sent.',
                icon: 'Check',
            });
        }
    };
    const handleReload = () => {
        sparkRTC.value.startProcedure();
    };

    const toggleBottomSheet = () => {};

    return html`<div class="flex gap-5 py-3 pt-0">
        <${Tooltip} label=${isMeetingMuted ? 'Listen' : 'Deafen'}>
            <${IconButton}
                variant=${isMeetingMuted && 'danger'}
                onClick=${toggleMuteMeeting}
                class="hidden sm:flex"
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
                class="hidden sm:flex"
            >
                <${Icon} icon="Share${sharingScreenStream ? 'Off' : ''}" />
            <//>
        <//>`}
        ${((!raiseHandMaxLimitReached.value && !isStreamming) ||
            (isStreamming && !isHost)) &&
        html`<${Tooltip} label=${
            isStreamming
                ? 'Put Hand Down'
                : ableToRaiseHand
                ? 'Raise Hand'
                : 'Raise hand request has been sent'
        }>
            <div>
						<${IconButton}
                onClick=${onRaiseHand}
                variant=${isStreamming && 'danger'}
                disabled=${!ableToRaiseHand}
            >
                <${Icon} icon="Hand" /> <//>
								<//>
						</div>`}
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
        html`
            <${Tooltip}
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
            <//>
        `}
        <${Tooltip} label=${'Menu'}>
            <${IconButton}
                onClick=${toggleBottomSheet}
                onClick=${toggleMoreOptions}
                class="flex sm:hidden"
            >
                <${Icon} icon="KebabMenuVertical" />
            <//>
        <//>
    </div>`;
};

export const MoreControllers = () => {
    const { isHost, sharingScreenStream, isStreamming, isMeetingMuted } =
        currentUser.value;
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
    return html`<div class="flex gap-5 py-5 justify-center">
        <${Tooltip} label=${isMeetingMuted ? 'Listen' : 'Deafen'}>
            <${IconButton}
                variant=${isMeetingMuted && 'danger'}
                onClick=${toggleMuteMeeting}
            >
                <${Icon} icon="Volume${isMeetingMuted ? 'Off' : ''}" />
            <//>
        <//>
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
    </div>`;
};

export default Controllers;
