import clsx from 'clsx';
import {
    BottomSheet,
    Container,
    Controllers,
    Icon,
    MoreControllers,
    Tooltip,
    attendeesBadge,
    attendeesCount,
    isAttendeesOpen,
    isMoreOptionsOpen,
    toggleAttendees,
    toggleMoreOptions,
} from 'components';
import { html } from 'htm';
import { broadcastIsInTheMeeting, isDebugMode } from '../../pages/meeting.js';

export const BottomBar = () => {
    return html`<${Container} class=${clsx('transition-all', {})}>
        <div
            class="w-full grid grid-cols-12 dark:bg-secondary-1-a py-3 pt-0 dark:text-gray-0 text-gray-2"
            id="bottom-bar"
        >
            <div class="col-span-3 sm:block hidden">
                <div class="h-full flex items-center">
                    <div
                        class="truncate max-w-full flex items-center gap-2 dark:bg-gray-2  bg-gray-0 rounded-full px-4 py-1"
                    >
                        <${Icon} icon="Link" class="min-w-[24px]" />
                        <span class="truncate flex-grow"
                            >https://great.ape/r/thu-kiu</span
                        >
                        <div class="ml-3 min-w-[24px]">
                            <${Icon} icon="Copy" />
                        </div>
                    </div>
                </div>
            </div>
            <div
                class="col-span-12 sm:col-span-6 flex items-center justify-center"
            >
                ${broadcastIsInTheMeeting.value
                    ? html`<${Controllers} />`
                    : null}
            </div>
            <div class="col-span-3 text-right sm:block hidden">
                <div class="h-full flex items-center justify-end">
                    ${attendeesCount > 0
                        ? html` <${Tooltip}
                              label=${isAttendeesOpen.value
                                  ? 'Hide Attendees'
                                  : 'Show Attendees'}
                          >
                              <div
                                  onClick="${toggleAttendees}"
                                  class="transition-all select-none cursor-pointer flex items-center gap-2 rounded-md hover:bg-gray-0 hover:bg-opacity-10 hover:dark:bg-gray-2 hover:dark:bg-opacity-20 py-1 px-3"
                              >
                                  <div class="relative">
                                      <${Icon} icon="Avatar" />

                                      ${attendeesBadge.value &&
                                      html`<span
                                          class="absolute top-0 -right-1 w-2 h-2 rounded-full bg-red-distructive"
                                      ></span>`}
                                  </div>
                                  <span
                                      >${attendeesCount}${' '}
                                      attendee${attendeesCount > 1
                                          ? 's'
                                          : ''}</span
                                  >
                              </div>
                          <//>`
                        : null}
                </div>
            </div>
        </div>
    <//>`;
};

export const BottomBarBottomSheet = () => {
    const handleAttendeesOpen = () => {
        toggleAttendees();
    };
    return html`<${BottomSheet}
        open=${isMoreOptionsOpen.value}
        onClose=${toggleMoreOptions}
        class="block sm:hidden"
        title="More"
    >
        <div class="w-full h-full flex gap-3 py-6 flex-col pb-0">
            <span class="text-bold-14"
                >Is Your Future Distributed? Welcome to the Fediverse!</span
            >
            <div
                class="mb-2 truncate text-gray-2 dark:text-gray-0 max-w-full flex items-center gap-2 dark:bg-gray-2  bg-gray-0 rounded-full px-4 py-1"
            >
                <${Icon} icon="Link" class="min-w-[24px]" />
                <span class="truncate flex-grow"
                    >https://great.ape/r/thu-kiu</span
                >
                <div class="ml-3 min-w-[24px]">
                    <${Icon} icon="Copy" />
                </div>
            </div>
            <${Tooltip}
                label=${isAttendeesOpen.value
                    ? 'Hide Attendees'
                    : 'Show Attendees'}
            >
                <div
                    onClick="${handleAttendeesOpen}"
                    class="w-full transition-all select-none cursor-pointer flex items-center gap-2 rounded-md hover:bg-gray-0 hover:bg-opacity-10 hover:dark:bg-gray-2 hover:dark:bg-opacity-20 py-1 px-3"
                >
                    <div class="relative">
                        <${Icon} icon="Avatar" />

                        ${attendeesBadge.value &&
                        html`<span
                            class="absolute top-0 -right-1 w-2 h-2 rounded-full bg-red-distructive"
                        ></span>`}
                    </div>
                    <span
                        >${attendeesCount}${' '}
                        attendee${attendeesCount > 1 ? 's' : ''}</span
                    >
                </div>
            <//>
            ${isDebugMode.value &&
            html`<${Tooltip} label="Troubleshoot">
                <div
                    class="w-full transition-all select-none cursor-pointer flex items-center gap-2 rounded-md hover:bg-gray-0 hover:bg-opacity-10 hover:dark:bg-gray-2 hover:dark:bg-opacity-20 py-1 px-3"
                >
                    <div class="relative">
                        <${Icon} icon="Troubleshoot" />
                    </div>
                    <span>Troubleshoot</span>
                </div>
            <//>`}
            <${MoreControllers} />
        </div>
    <//>`;
};

export default BottomBar;
