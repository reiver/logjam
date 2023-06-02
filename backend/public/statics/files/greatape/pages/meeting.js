import { html } from 'htm';
import {
    TopBar,
    BottomBar,
    MeetingBody,
    attendees,
    streamers,
    makeDialog,
    Button,
} from 'components';
import { createSparkRTC, getWsUrl } from '../lib/common.js';
import { signal } from '@preact/signals';
import { useEffect } from 'preact';

export const sparkRTC = signal(null);
export const meetingStatus = signal(true);
export const broadcastIsInTheMeeting = signal(true);
export const currentUser = signal({
    isHost: false,
    isMicrophoneOn: true,
    isCameraOn: true,
    isMeetingMuted: false,
    sharingScreenStream: null,
    ableToRaiseHand: true,
    hasMic: true,
    hasCamera: true,
});

export const updateUser = (props) => {
    currentUser.value = {
        ...currentUser.value,
        ...props,
    };
};

export const onStartShareScreen = (stream) => {
    log(`ScreenShareStram: ${stream}`);

    stream.getTracks()[0].onended = async () => {
        await sparkRTC.value.stopShareScreen(stream);
        updateUser({
            sharingScreenStream: null,
        });
        onStopStream(stream);
    };

    streamers.value = {
        ...streamers.value,
        [stream.id]: {
            name: stream.name,
            isHost: true,
            avatar: '',
            raisedHand: false,
            hasCamera: false,
            stream,
        },
    };
};

export const onStopStream = (stream) => {
    const streamersTmp = { ...streamers.value };
    delete streamersTmp[stream.id];
    streamers.value = streamersTmp;
};

export const onStopShareScreen = (stream) => {
    stream.getTracks().forEach((track) => track.stop());
    updateUser({
        sharingScreenStream: null,
    });
    const streamersTmp = { ...streamers.value };
    delete streamersTmp[stream.id];
    streamers.value = streamersTmp;
};

const log = (tag, data) => {
    const date = new Date().toLocaleTimeString();

    if (data) {
        console.log('[', date, '] ', tag, ' | ', data);
    } else {
        console.log('[', date, '] ', tag);
    }
};
const setupSignalingSocket = async (host, name, room) => {
    await sparkRTC.value.setupSignalingSocket(
        getWsUrl(host),
        JSON.stringify({ name, email: '' }),
        room
    );
};
const start = async () => {
    return sparkRTC.value.start();
};

export const leaveMeeting = () => {
    if (sparkRTC.value) {
        sparkRTC.value.leaveMeeting();
        meetingStatus.value = false;
        streamers.value = [];
    }
};

const Meeting = () => {
    useEffect(() => {
        const queryParams = new URLSearchParams(window.location.search);
        const name = queryParams.get('name');
        const role = queryParams.get('role');
        const room = queryParams.get('room');
        const host = queryParams.get('host');

        updateUser({
            name,
            role,
            isStreamming: role === 'broadcast',
            isHost: role === 'broadcast',
        });

        const setupSparkRTC = async () => {
            log(`Setup SparkRTC`);

            sparkRTC.value = createSparkRTC(role, {
                localStreamChangeCallback: (stream) => {
                    log('[Local Stream Callback]', stream);
                    streamers.value = {
                        ...streamers.value,
                        [stream.id]: {
                            name: stream.name,
                            isHost: role === 'broadcast',
                            avatar: '',
                            raisedHand: false,
                            hasCamera: false,
                            stream,
                        },
                    };
                },
                remoteStreamCallback: (stream) => {
                    sparkRTC.value.getLatestUserList();

                    log(`remoteStreamCallback`, stream);
                    log(`remoteStreamCallback-Name`, stream.name);

                    streamers.value = {
                        ...streamers.value,
                        [stream.id]: {
                            name: stream.name,
                            userId: stream.userId,
                            isHost: stream.role === 'broadcast',
                            avatar: '',
                            raisedHand: false,
                            hasCamera: false,
                            stream,
                        },
                    };

                    if (!sparkRTC.value.broadcasterDC && role === 'audience') {
                        broadcastIsInTheMeeting.value = true;
                    }
                },
                remoteStreamDCCallback: (stream) => {
                    sparkRTC.value.getLatestUserList();

                    log(`remoteStreamDCCallback`, stream);

                    onStopStream(stream);

                    if (role === 'audience') {
                        if (
                            sparkRTC.value.broadcasterDC ||
                            stream === 'no-stream'
                        ) {
                            broadcastIsInTheMeeting.value = false;
                            log(`broadcasterDC...`);
                        }
                    }
                },
                onRaiseHand: (...props) => {
                    log(`[On Raise Hand Request] ${props}`);
                    return new Promise((resolve, reject) => {
                        makeDialog(
                            'confirm',
                            props[0],
                            resolve.bind(null, true),
                            resolve.bind(null, false)
                        );
                    });
                },
                onStart: async (closeSocket = true) => {
                    if (meetingStatus.value) {
                        if (role === 'audience') {
                            await sparkRTC.value.restart(closeSocket);
                        }

                        //if socket is closed, repoen again
                        if (
                            sparkRTC.value.socket &&
                            (sparkRTC.value.socket.readyState ===
                                WebSocket.CLOSED ||
                                sparkRTC.value.socket.readyState ===
                                    WebSocket.CLOSING)
                        ) {
                            await setupSignalingSocket(host, name, room);
                        }
                        
                        //start sparkRTC
                        await start();
                    }
                },
                altBroadcastApprove: (isStreamming) => {
                    updateUser({ isStreamming });
                    if (!isStreamming) {
                        sparkRTC.value.onRaiseHandRejected();
                        makeDialog(
                            'error',
                            'Your raise hand request rejected!'
                        );
                    } else {
                        makeDialog(
                            'success',
                            'Your raise hand request Approved!'
                        );
                    }
                },
                disableBroadcasting: () => {
                    updateUser({ isStreamming: false });
                    makeDialog('info', 'You just removed from broadcasting');
                    sparkRTC.value.onRaiseHandRejected();
                },
                maxLimitReached: (message) => {
                    makeDialog('error', message);
                },
                onUserListUpdate: (users) => {
                    log(`[On Users List Update], ${users}`);
                    attendees.value = users.map(
                        ({ name: userInfo, role, video }) => {
                            const { name, email } = JSON.parse(userInfo);
                            return {
                                name,
                                isHost: role === 'broadcaster',
                                avatar: '',
                                raisedHand: false,
                                hasCamera: !!video,
                            };
                        }
                    );
                },
                constraintResults: (constraints) => {
                    if (!constraints.audio) {
                        //remove mic button
                        updateUser({ hasMic: false });
                    }

                    if (!constraints.video) {
                        //remove video button
                        updateUser({ hasCamera: false });
                    }
                },
                updateStatus: (tag, data) => {
                    log(tag, data);
                },
                treeCallback: (tree) => {},
                connectionStatus: (status) => {
                    log(`Connection Status: `, status);
                    if (role === 'audience') {
                        if (status === 'failed') {
                            window.location.reload();
                        }
                    }
                },
            });

            if (sparkRTC.value) {
                //Init socket and start sparkRTC
                await setupSignalingSocket(host, name, room);
                await start();
            }
        };

        if (meetingStatus.value) {
            setupSparkRTC();
        }
    }, [meetingStatus.value]);

    const rejoinMeeting = () => {
        window.location.reload();
    };

    return html` <div
        class="flex flex-col justify-between min-h-screen dark:bg-secondary-1-a bg-white-f-9 text-medium-12 text-gray-800 dark:text-gray-200"
    >
        ${meetingStatus.value
            ? html`<${TopBar} />
                  <${MeetingBody} />
                  <${BottomBar} />`
            : html`<div
                  class="flex flex-col justify-center items-center p-10 rounded-md gap-2"
              >
                  <span>You left the meeting</span>
                  <${Button} onClick=${rejoinMeeting}>Rejoin<//>
              </div>`}
    </div>`;
};

export default Meeting;
