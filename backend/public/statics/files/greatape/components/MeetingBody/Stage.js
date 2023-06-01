import { computed, signal } from '@preact/signals';
import clsx from 'clsx';
import { Icon, IconButton, attendeesWidth, makeDialog } from 'components';
import { html } from 'htm';
import { useEffect, useRef, useState } from 'preact';
import { memo } from 'preact/compat';
import { userInteractedWithDom } from '../../index.js';
import {
    broadcastIsInTheMeeting,
    currentUser,
    sparkRTC,
} from '../../pages/meeting.js';

export const streamers = signal([]);
export const streamersLength = computed(
    () => Object.keys(streamers.value).length
);

const topBarBottomBarHeight = 58 + 108;
const windowWidth = signal(window.innerWidth);
const windowHeight = signal(window.innerHeight);
const stageWidth = computed(() => windowWidth.value - attendeesWidth - 130);
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
        windowHeight.value - topBarBottomBarHeight - gapHeight;

    if (availableHeight < lines * height) {
        height =
            availableHeight / Math.ceil(streamersLength.value / eachPerLine);

        width = (height * 16) / 9;
    }

    return width;
});

export const Stage = () => {
    useEffect(() => {
        const onResize = () => {
            windowWidth.value = window.innerWidth;
            windowHeight.value = window.innerHeight;
        };
        window.addEventListener('resize', onResize);
        return () => window.removeEventListener('resize', onResize);
    }, []);

    return html`<div
        class="transition-all relative h-full px-4 lg:px-10"
        style="width: calc(100% - ${attendeesWidth + 20}px);"
    >
        ${broadcastIsInTheMeeting.value
            ? html`<div
                  class="flex gap-4 flex-wrap justify-center items-center h-full"
              >
                  ${Object.values(streamers.value).map((attendee, i) => {
                      return html`<div
                          key=${i}
                          style="width: ${itemsWidth.value}px"
                          class=${clsx(
                              'transition-all aspect-video relative max-w-full min-w-[100%] sm:min-w-[unset] text-white-f-9',
                              'bg-gray-1 rounded-lg min-w-10',
                              'dark:bg-gray-3 dark:'
                          )}
                      >
                          <${Video}
                              stream=${attendee.stream}
                              userId=${attendee.userId}
                              isMuted=${currentUser.value.isMeetingMuted}
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
    const handleOpenMenu = setMenuOpen.bind(null, !menuOpen);

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

    return html`<video
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
        ${isHost &&
        !isHostStream &&
        html` <${IconButton} class="absolute top-3 right-3" variant="ghost" onClick=${handleOpenMenu} ref=${menu}>
      <${Icon} icon="verticalDots" width="20px" height="20px" />

      ${
          menuOpen &&
          html`<div class="relative top-0 right-0 h-full w-full">
              <ul
                  class="bg-white absolute top-0 right-0 mt-3 -ml-2 text-black rounded-sm p-1"
              >
                  <li
                      class="w-full whitespace-nowrap px-4 py-1 rounded-sm bg-black bg-opacity-0 hover:bg-opacity-10"
                      onClick=${handleRemoveStream}
                  >
                      Stop broadcast
                  </li>
              </ul>
          </div>`
      }
      <div>
				<//>
				</${IconButton}>`} `;
});
