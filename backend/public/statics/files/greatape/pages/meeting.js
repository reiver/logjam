import { computed, signal } from '@preact/signals';
import {
    BottomBar,
    Button,
    MeetingBody,
    TopBar,
    attendees,
    attendeesBadge,
    isMoreOptionsOpen,
    makeDialog,
    streamers,
} from 'components';
import { html } from 'htm';
import { useEffect } from 'preact';
import { isAttendeesOpen } from '../components/Attendees/index.js';
import { createSparkRTC, getWsUrl } from '../lib/common.js';

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

    isMoreOptionsOpen.value = false;

    streamers.value = {
        ...streamers.value,
        [stream.id]: {
            name: stream.name,
            isHost: true,
            avatar: '',
            raisedHand: false,
            hasCamera: false,
            stream,
            isShareScreen: true,
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

const log = (data) => {
    const date = new Date().toLocaleTimeString();
    console.log('[', date, '] ', data);
};

const startSocket = async (name, room, host = null) => {
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
            isStreamming: role === 'broadcast',
            isHost: role === 'broadcast',
        });
        const setupSparkRTC = async () => {
            sparkRTC.value = createSparkRTC(role, {
                onUserInitialized: (userId) => {
                    currentUser.userId = userId;
                },
                localStreamChangeCallback: (stream) => {
                    log('[Local Stream Callback]', stream);
                    console.log(stream);
                    streamers.value = {
                        ...streamers.value,
                        [stream.id]: {
                            name: stream.name,
                            isHost: role === 'broadcast',
                            avatar: '',
                            raisedHand: false,
                            hasCamera: false,
                            stream,
                            isShareScreen: stream.isShareScreen || false,
                        },
                    };
                },
                remoteStreamCallback: (stream) => {
                    sparkRTC.value.getLatestUserList();
                    log(`[Remote Stream Callback] ${stream}`);
                    console.log(stream);
                    log(`NameCallback: ${stream.name}`);

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
                            isShareScreen: stream.isShareScreen || false,
                        },
                    };

                    if (!sparkRTC.value.broadcasterDC && role === 'audience') {
                        broadcastIsInTheMeeting.value = true;
                    }
                },
                remoteStreamDCCallback: (stream) => {
                    sparkRTC.value.getLatestUserList();

                    onStopStream(stream);

                    if (role === 'audience') {
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
                onRaiseHand: (message, user, ...rest) => {
                    log(`[On Raise Hand Request] ${message}`, user, rest);

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
                    onUserRaisedHand(
                        user.userId,
                        new Date(),
                        raiseHandCallback
                    );

                    return handler;
                },
                onStart: async () => {
                    if (meetingStatus.value) {
                        if (role === 'audience') {
                            let idList = [];
                            for (const id in sparkRTC.value
                                .myPeerConnectionArray) {
                                const peerConn =
                                    sparkRTC.value.myPeerConnectionArray[id];
                                await peerConn.close();
                                idList.push(id);
                            }
                            idList.forEach(
                                (id) =>
                                    delete sparkRTC.value.myPeerConnectionArray[
                                        id
                                    ]
                            );
                            sparkRTC.value.remoteStreams = [];
                            sparkRTC.value.localStream
                                ?.getTracks()
                                ?.forEach((track) => track.stop());
                            sparkRTC.value.localStream = null;
                        }

                        await startSocket(name, room, host);
                    }
                },
                altBroadcastApprove: (isStreamming) => {
                    updateUser({ isStreamming, ableToRaiseHand: true });
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
                            isHost: role === 'broadcaster',
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
                updateStatus: (status) => {
                    log(status);
                },
                treeCallback: (tree) => {},
                signalingDisconnectedCallback: () => {},
            });

            log(`Setup SparkRTC`);
            await sparkRTC.value.setupSignalingSocket(
                getWsUrl(host),
                JSON.stringify({ name, email: '' }),
                room
            );
            await startSocket(name, room, host);
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
        <${TopBar} />
        ${meetingStatus.value
            ? html` <${MeetingBody} />
                  <${BottomBar} />`
            : html`<div
                  class="flex flex-col justify-center items-center sm:p-10 rounded-md gap-4 h-full flex-grow"
              >
                  <span class="text-bold-18 text-center"
                      >You Left The Live Show</span
                  >
                  <div
                      class="flex w-full justify-center items-center gap-4 flex-row max-w-[85%] sm:max-w-[400px]"
                  >
                      <${Button}
                          onClick=${rejoinMeeting}
                          variant="outline"
                          class="flex-1 w-full px-0"
                          >Go To Home Feed<//
                      >
                      <${Button}
                          onClick=${rejoinMeeting}
                          variant="primary"
                          class="flex-1 w-full px-0"
                          >Rejoin<//
                      >
                  </div>
              </div>`}
    </div>`;
};

export default Meeting;
