import { computed, signal } from '@preact/signals';
import {
    BottomBar,
    Button,
    MeetingBody,
    TopBar,
    attendees,
    attendeesBadge,
    makeDialog,
    streamers,
} from 'components';
import { html } from 'htm';
import { useEffect } from 'preact';
import { isAttendeesOpen } from '../components/Attendees/index.js';
import { createSparkRTC, getWsUrl, Roles } from '../lib/common.js';

export const sparkRTC = signal(null);
export const meetingStatus = signal(true);
export const broadcastIsInTheMeeting = signal(true);
export const raisedHandsCount = signal(0);
export const raiseHandMaxLimitReached = computed(() => {
    return (
        sparkRTC.value &&
        raisedHandsCount.value === sparkRTC.value.maxRaisedHands
    );
});
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
        streamers.value = {};
    }
};

export const onUserRaisedHand = (userId, raisedHand, acceptRaiseHand) => {
    attendees.value = {
        ...attendees.value,
        [userId]: {
            ...attendees.value[userId],
            raisedHand,
            acceptRaiseHand,
        },
    };
    console.log('LOWER HAND', userId, raisedHand);
};

export const getUserRaiseHandStatus = (userId) => {
    return attendees.value[userId]?.raisedHand || false;
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
            isStreamming: role === Roles.BROADCAST,
            isHost: role === Roles.BROADCAST,
        });

        const setupSparkRTC = async () => {
            log(`Setup SparkRTC`);

            sparkRTC.value = createSparkRTC(role, {
                onUserInitialized: (userId) => {
                    currentUser.userId = userId;
                },
                localStreamChangeCallback: (stream) => {
                    log('[Local Stream Callback]', stream);
                    streamers.value = {
                        ...streamers.value,
                        [stream.id]: {
                            name: stream.name,
                            isHost: role === Roles.BROADCAST,
                            avatar: '',
                            raisedHand: false,
                            hasCamera: false,
                            stream,
                            isLocalStream: true,
                        },
                    };
                },
                remoteStreamCallback: (stream) => {
                    sparkRTC.value.getLatestUserList();

                    log(`remoteStreamCallback`, stream);
                    log(`remoteStreamCallback-Name`, stream.name);

                    //if receive my localStream then mute it for me
                    let local = false;
                    if (sparkRTC.value.localStream) {
                        if (sparkRTC.value.localStream.id === stream.id) {
                            local = true;
                        }
                    }

                    streamers.value = {
                        ...streamers.value,
                        [stream.id]: {
                            name: stream.name,
                            userId: stream.userId,
                            isHost: stream.role === Roles.BROADCAST,
                            avatar: '',
                            raisedHand: false,
                            hasCamera: false,
                            stream,
                            isLocalStream: local,
                        },
                    };

                    if (
                        !sparkRTC.value.broadcasterDC &&
                        role === Roles.AUDIENCE
                    ) {
                        broadcastIsInTheMeeting.value = true;
                    }
                },
                remoteStreamDCCallback: (stream) => {
                    sparkRTC.value.getLatestUserList();

                    log(`remoteStreamDCCallback`, stream);

                    onStopStream(stream);

                    if (role === Roles.AUDIENCE) {
                        if (
                            sparkRTC.value.broadcasterDC ||
                            stream === 'no-stream'
                        ) {
                            broadcastIsInTheMeeting.value = false;
                            updateUser({
                                isStreamming: false,
                                ableToRaiseHand: true,
                            });
                            log(`broadcasterDC...`);
                        }
                    }
                },
                onRaiseHand: (user) => {
                    log(`[On Raise Hand Request]`, user);

                    let raiseHandCallback = () => {};
                    const handler = new Promise((resolve, reject) => {
                        raiseHandCallback = resolve;
                    });

                    attendeesBadge.value = true;

                    makeDialog(
                        'info',
                        {
                            message: 'Someone has raised their hand!',
                            icon: 'Clock',
                        },
                        null,
                        () => {
                            isAttendeesOpen.value = true;
                            attendeesBadge.value = false;
                        }
                    );
                    onUserRaisedHand(user.userId, true, raiseHandCallback);

                    return handler;
                },
                onStart: async (closeSocket = false) => {
                    if (meetingStatus.value) {
                        if (role === Roles.AUDIENCE) {
                            await sparkRTC.value.restart(closeSocket);
                        }

                        if (!closeSocket) {
                            //start sparkRTC
                            await start();
                        }
                    }
                },
                startAgain: async () => {
                    if (sparkRTC.value) {
                        //Init socket and start sparkRTC
                        await setupSignalingSocket(host, name, room);
                        await start();
                    }
                },
                altBroadcastApprove: (isStreamming) => {
                    updateUser({ isStreamming });
                    if (!isStreamming) {
                        sparkRTC.value.onRaiseHandRejected();
                        makeDialog('info', {
                            message: 'Your raise hand request rejected!',
                            icon: 'Close',
                        });
                    } else {
                        makeDialog('info', {
                            message: 'You’ve been added to the stage',
                            icon: 'Check',
                        });
                    }
                },
                disableBroadcasting: () => {
                    updateUser({ isStreamming: false });
                    makeDialog('info', {
                        message: 'You just removed from stage',
                        icon: 'Close',
                    });
                    sparkRTC.value.onRaiseHandRejected();
                },
                maxLimitReached: (message) => {
                    makeDialog('info', { message, icon: 'Close' });
                },
                onUserListUpdate: (users) => {
                    log(`[On Users List Update], ${users}`);
                    const usersTmp = {};
                    for (const {
                        name: userInfo,
                        role,
                        video,
                        id: userId,
                    } of users) {
                        const { name, email } = JSON.parse(userInfo);
                        usersTmp[userId] = {
                            ...(attendees.value[userId] || {}),
                            name,
                            email,
                            isHost: role === Roles.BROADCASTER,
                            avatar: '',
                            raisedHand: getUserRaiseHandStatus(userId),
                            hasCamera: !!video,
                            userId,
                            video,
                        };
                    }
                    raisedHandsCount.value = Object.values(usersTmp).reduce(
                        (prev, user) => {
                            if (!user.isHost && user.video) return prev + 1;
                            return prev;
                        },
                        0
                    );

                    attendees.value = usersTmp;
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
                treeCallback: (tree) => {
                    log(`tree`, tree);
                },
                connectionStatus: (status) => {
                    log(`Connection Status: `, status);
                    if (role === Roles.AUDIENCE) {
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
