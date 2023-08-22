import { signal } from '@preact/signals';
import { clsx } from 'clsx';
import { Button, Icon, IconButton, Tooltip } from 'components';
import { html } from 'htm';
import { useEffect, useRef } from 'preact';
import { v4 as uuidv4 } from 'uuid';

import { currentUser, sparkRTC, updateUser } from '../../pages/meeting.js';

const dialogs = signal([]);

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
    const videoRef = useRef();
    const { hasCamera, hasMic, isCameraOn, isMicrophoneOn } = currentUser.value;

    useEffect(() => {
        videoRef.current.srcObject = videoStream;
    }, [videoStream]);

    useEffect(() => {
        videoRef.current.playsInline = true;
        // videoRef.current.play();
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

    return html` <div class="absolute top-0 left-0 w-full h-full">
        <div
            class="z-10 absolute w-full h-full bg-black bg-opacity-60"
        />
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
    // if (type !== 'confirm') {
    //     setTimeout(destroy, 4000);
    // }
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
};

export const ToastProvider = () => {
    return html`<div id="toast-provider">
        <${DialogPool} />
    </div>`;
};
