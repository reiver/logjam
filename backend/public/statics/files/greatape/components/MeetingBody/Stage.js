import { computed, signal } from '@preact/signals';
import clsx from 'clsx';
import { Icon, IconButton, attendeesWidth, makeDialog } from 'components';
import { html } from 'htm';
import throttle from 'lodash.throttle';
import { useEffect, useRef, useState } from 'preact';
import { memo } from 'preact/compat';
import { getDeviceConfig } from '../../hooks/use-breakpoint.js';
import { userInteractedWithDom } from '../../index.js';
import {
    broadcastIsInTheMeeting,
    currentUser,
    sparkRTC,
} from '../../pages/meeting.js';
let timeOut;
export const bottomBarVisible = signal(true);
export const fullScreenedStream = signal(null);
export const hasFullScreenedStream = computed(() => !!fullScreenedStream.value);
export const streamers = signal({});
export const streamersLength = computed(
    () => Object.keys(streamers.value).length
);
export const deviceSize = signal(getDeviceConfig(window.innerWidth));
const topBarBottomBarHeight = () =>
    document.getElementById('top-bar').offsetHeight +
    (bottomBarVisible.value
        ? document.getElementById('bottom-bar').offsetHeight
        : 0) +
    32;
const windowWidth = signal(window.innerWidth);
const windowHeight = signal(window.innerHeight);
const stageWidth = computed(
    () =>
        windowWidth.value -
        attendeesWidth.value -
        (deviceSize.value !== 'xs' ? 140 : 32)
);
const itemsWidth = computed(() => {
    let width = Math.max(
        stageWidth.value / streamersLength.value,
        stageWidth.value / 2
    );
    let height = (width * 9) / 16;
    let eachPerLine = width == stageWidth / 2 ? 2 : 1;
    const lines = Math.ceil(streamersLength.value / eachPerLine);
    const gapHeight = (lines - 1) * 16 + 16;

    const availableHeight =
        windowHeight.value - topBarBottomBarHeight() - gapHeight;

    if (availableHeight < lines * height) {
        height =
            availableHeight / Math.ceil(streamersLength.value / eachPerLine);

        width = (height * 16) / 9;
    }

    return width;
});

export const getVideoWidth = (attendee, index) => {
    if (deviceSize.value === 'xs') {
        let availableHeight = windowHeight.value - topBarBottomBarHeight();
        if (hasFullScreenedStream.value) {
            if (attendee.stream.id === fullScreenedStream.value) {
                return `100%; height: ${availableHeight}px`;
            } else return `0px; height: 0px;`;
        }
        const lines = Math.ceil(streamersLength.value / 2);
        const gapHeight = (lines - 1) * 16 + 16;
        availableHeight -= gapHeight;
        if (index == 0) {
            if (streamersLength.value === 1) {
                return `calc(100%); height: ${availableHeight}px`;
            }
            return `calc(100%); height: ${availableHeight / 2}px`;
        } else {
            const lines = Math.ceil((streamersLength.value - 1) / 2);
            availableHeight = availableHeight / 2;
            let rowHeight = availableHeight / lines;
            const columns = streamersLength.value - 1 > 1 && lines >= 1 ? 2 : 1;

            return `calc(${100 / columns}% - ${
                columns > 1 ? '8px' : '0px'
            }); height: ${rowHeight}px`;
        }
    }
    let availableHeight = windowHeight.value - topBarBottomBarHeight();
    if (hasFullScreenedStream.value) {
        if (attendee.stream.id === fullScreenedStream.value) {
            return `100%; height: ${availableHeight}px`;
        } else return `0px; height: 0px;`;
    }
    let height = (itemsWidth.value * 9) / 16;
    return `${itemsWidth.value}px;height: ${height}px;`;
};

export const Stage = () => {
    useEffect(() => {
        const onResize = throttle(() => {
            windowWidth.value = window.innerWidth;
            windowHeight.value = window.innerHeight;
            deviceSize.value = getDeviceConfig(window.innerWidth);
        }, 200);
        window.addEventListener('resize', onResize);
        return () => window.removeEventListener('resize', onResize);
    }, []);

    const documentClick = () => {
        if (timeOut) clearTimeout(timeOut);
        if (hasFullScreenedStream.value) {
            bottomBarVisible.value = true;
            handleMaximize();
        }
    };
    const handleMaximize = () => {
        if (timeOut) clearTimeout(timeOut);
        timeOut = setTimeout(() => {
            if (bottomBarVisible.value) {
                bottomBarVisible.value = false;
            }
        }, 2000);
    };
    useEffect(() => {
        if (hasFullScreenedStream.value) {
            handleMaximize();
            document
                .getElementsByTagName('body')[0]
                .addEventListener('click', documentClick);
        } else {
            document
                .getElementsByTagName('body')[0]
                .removeEventListener('click', documentClick);
            bottomBarVisible.value = true;
            if (timeOut) clearTimeout(timeOut);
        }
    }, [hasFullScreenedStream.value]);
    const handleOnClick = (e, streamId) => {
        if (streamId === fullScreenedStream.value) {
            bottomBarVisible.value = !bottomBarVisible.value;
            handleMaximize();
            e.stopPropagation();
        }
    };
    return html`<div
        class="transition-all h-full lg:px-0 relative"
        style="width: calc(100% - ${attendeesWidth}px);"
    >
        ${broadcastIsInTheMeeting.value
            ? html`<div
                  class=${clsx(
                      'flex flex-wrap justify-start sm:justify-center items-center h-full transition-all',
                      {
                          'gap-4': !hasFullScreenedStream.value,
                          'gap-0': hasFullScreenedStream.value,
                      }
                  )}
              >
                  ${Object.values(streamers.value)
                      .sort((a, b) => {
                          let aScore = 0;
                          let bScore = 0;
                          if (a.isHost) aScore += 10;
                          if (a.isShareScreen) aScore += 20;
                          if (b.isHost) bScore += 10;
                          if (b.isShareScreen) bScore += 20;
                          return bScore - aScore;
                      })
                      .map((attendee, i) => {
                          let muted = false;

                          //mute the stream if it's my local stream
                          if (attendee.isLocalStream === true) {
                              muted = true;
                          } else {
                              //mute it based on meeting status
                              muted = currentUser.value.isMeetingMuted;
                          }

                          return html`<div
                              key=${i}
                              style="width: ${getVideoWidth(attendee, i)}"
                              class=${clsx(
                                  'group transition-all aspect-video relative max-w-full text-white-f-9',
                                  'bg-gray-1 rounded-lg min-w-10',
                                  'dark:bg-gray-3 overflow-hidden'
                              )}
                              onClick=${(e) =>
                                  handleOnClick(e, attendee.stream.id)}
                          >
                              <${Video}
                                  stream=${attendee.stream}
                                  userId=${attendee.userId}
                                  isMuted=${muted}
                                  name=${attendee.name}
                                  isHostStream=${attendee.isHost}
                              />
                          </div>`;
                      })}
              </div>`
            : html`<span class="inline-block w-full text-center">
                  The broadcaster is not in the meeting, please wait until the
                  broadcaster joins
              </span>`}
    </div>`;
};

export const Video = memo(({ stream, isMuted, isHostStream, name, userId }) => {
    const [muted, setMuted] = useState(true);
    const { isHost } = currentUser.value;
    const menu = useRef();
    const videoRef = useRef();
    const [menuOpen, setMenuOpen] = useState(false);
    const toggleFullScreen = (e) => {
        if (fullScreenedStream.value === stream.id) {
            fullScreenedStream.value = null;
        } else fullScreenedStream.value = stream.id;

        e.stopPropagation();
    };
    useEffect(() => {
        videoRef.current.srcObject = stream;
    }, [stream]);

    useEffect(() => {
        if (userInteractedWithDom.value) setMuted(isMuted);
    }, [userInteractedWithDom.value, isMuted]);

    const handleRemoveStream = () => {
        makeDialog(
            'confirm',
            {
                message: `Are you sure you want to kick "<strong>${name}</strong>" off the stage?`,
                title: 'Kick Audience Off The Stage',
            },
            () => {
                sparkRTC.value.disableAudienceBroadcast(String(userId));
            },
            () => {},
            {
                okText: 'Kick',
                okButtonVariant: 'red',
                cancelText: 'Let them stay!',
            }
        );
    };
    const handleOpenMenu = (e) => {
        e.stopPropagation();
        setMenuOpen(!menuOpen);
    };
    const [isHover, setHover] = useState(false);
    console.log('init', isHover, bottomBarVisible.value);
    const handleOnClick = () => {
        setHover(!isHover);
    };
    useEffect(() => {
        if (
            (!bottomBarVisible.value && isHover) ||
            (!hasFullScreenedStream.value && isHover)
        ) {
            setHover(false);
        }
    }, [bottomBarVisible.value, hasFullScreenedStream.value]);
    useEffect(() => {
        function handleClickOutside(event) {
            if (
                menuOpen &&
                menu.current &&
                !menu.current.base.contains(event.target)
            ) {
                setMenuOpen(false);
            }
        }
        if (menuOpen)
            document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [menu, menuOpen]);

    return html`
        <div onClick=${handleOnClick} className="w-full h-full rounded-lg">
            <video
                ref=${videoRef}
                autoplay
                muted="${muted}"
                className="w-full h-full object-cover rounded-lg"
            />
            <div
                class="px-4 py-1 bg-black bg-opacity-50 text-white rounded-full absolute top-3 left-3 text-medium-12"
            >
                ${name} ${isHostStream && '(Host)'}
            </div>
            <div
                class=${clsx(
                    'absolute top-3 sm:group-hover:flex right-3 sm:hidden gap-2',
                    {
                        'group-hover:flex': isHover && bottomBarVisible.value,
                        hidden: !(isHover && bottomBarVisible.value),
                        flex: menuOpen || isHover,
                    }
                )}
            >
                <${IconButton} variant="ghost" onClick=${toggleFullScreen}>
                    <${Icon}
                        icon=${fullScreenedStream.value === stream.id
                            ? 'ScreenNormal'
                            : 'ScreenFull'}
                        width="20px"
                        height="20px"
                    />
                <//>
                ${isHost &&
                !isHostStream &&
                html`
                    <${IconButton}
                        variant="ghost"
                        onClick=${handleOpenMenu}
                        ref=${menu}
                    >
                        <${Icon}
                            icon="verticalDots"
                            width="20px"
                            height="20px"
                        />

                        ${menuOpen &&
                        html`<div
                            class="absolute top-full right-0 h-full w-full"
                        >
                            <ul
                                class="bg-white absolute top-0 right-0 mt-1 -ml-2 text-black rounded-sm p-1"
                            >
                                <li
                                    class="w-full whitespace-nowrap px-4 py-1 rounded-sm bg-black bg-opacity-0 hover:bg-opacity-10"
                                    onClick=${handleRemoveStream}
                                >
                                    Stop broadcast
                                </li>
                            </ul>
                        </div>`}
                    <//>
                `}
            </div>
        </div>
    `;
});
