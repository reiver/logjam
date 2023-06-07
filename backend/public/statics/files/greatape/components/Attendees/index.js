import { computed, signal } from '@preact/signals';
import clsx from 'clsx';
import { BottomSheet, Icon, makeDialog } from 'components';
import { html } from 'htm';
import {
    currentUser,
    onUserRaisedHand,
    raiseHandMaxLimitReached,
} from '../../pages/meeting.js';
import { deviceSize } from '../MeetingBody/Stage.js';
export const attendees = signal(
    {}
    // {
    //   name: 'Alex Suprun',
    //   isHost: true,
    //   avatar: 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1480&q=80',
    //   raisedHand: false,
    //   hasCamera: true,
    // }
);

export const attendeesCount = computed(
    () => Object.values(attendees.value).length
);

export const attendeesBadge = signal(false);

export const isAttendeesOpen = signal(false);
export const toggleAttendees = () => {
    attendeesBadge.value = false;
    isAttendeesOpen.value = !isAttendeesOpen.value;
};

export const attendeesWidth = computed(() => {
    if (!isAttendeesOpen.value || deviceSize.value === 'xs') return 0;
    return 350 + 40;
});

export const Participant = ({ participant }) => {
    const handleRaiseHand = () => {
        makeDialog(
            'confirm',
            {
                message: `"<strong>${participant.name}</strong>" has raised their hand, do you want to add them to the stage?`,
                title: 'Accept Raised Hand',
            },
            () => {
                participant.acceptRaiseHand(true);
                onUserRaisedHand(participant.userId, false);
            },
            () => {},
            {
                onReject: () => {
                    participant.acceptRaiseHand(false);
                    onUserRaisedHand(participant.userId, false);
                },
            }
        );
    };

    const raisedHand =
        participant.raisedHand && !raiseHandMaxLimitReached.value;
    return html` <div
        class=${clsx(
            'flex w-full justify-between items-center rounded-md px-2 py-1',
            'cursor-pointer hover:dark:bg-white hover:dark:bg-opacity-10 hover:bg-gray-500 hover:bg-opacity-10 transition-all'
        )}
    >
        <div class="flex items-center gap-2">
            ${participant.avatar
                ? html`<img
                      src="${participant.avatar}"
                      class="w-9 h-9 rounded-full object-cover"
                  />`
                : html`<div
                      class="dark:bg-gray-300 dark:bg-opacity-30 bg-opacity-30 bg-gray-400 rounded-full w-9 h-9 flex justify-center items-center"
                  >
                      <${Icon} icon="Avatar" width="20px" height="20px" />
                  </div>`}

            <div class="flex flex-col justify-center">
                <span class="text-gray-1 dark:text-gray-0 "
                    ><span class="text-bold-12 text-gray-3 dark:text-white-f-9 "
                        >${participant.name}</span
                    >
                    ${participant.userId == currentUser.userId ? ' (You)' : ''}
                </span>
                ${participant.isHost &&
                html`<span class="text-gray-1 dark:text-gray-0 text-regular-12"
                    >Host</span
                >`}
            </div>
        </div>
        ${(raisedHand || participant.hasCamera) &&
        html`
            <div>
                <${Icon}
                    icon=${raisedHand ? 'Hand' : 'Camera'}
                    width="25"
                    height="25px"
                    class="dark:text-gray-0 text-gray-1"
                    onClick=${raisedHand ? handleRaiseHand : null}
                />
            </div>
        `}
    </div>`;
};

export const Attendees = () => {
    return html`
        <div
            class="${clsx(
                'h-full min-w-[350px] border rounded-lg p-2 pb-0',
                'bg-white-f border-gray-0 text-secondary-1-a',
                'dark:bg-gray-3 dark:border-0 dark:text-white-f-9',
                'absolute top-0',
                'transition-all ease-in-out',
                'lg:right-10 right-4',
                {
                    'translate-x-[100%] lg:-mr-10 -mr-4':
                        !isAttendeesOpen.value,
                    'translate-x-[100%]': !isAttendeesOpen.value,
                },
                'hidden sm:block'
            )}"
            onClick=${() => (attendeesBadge.value = false)}
        >
            <div class="flex flex-col pt-2 gap-2 max-h-full">
                <div class="flex justify-center items-center w-full gap-2">
                    <${Icon} icon="Avatar" />
                    <span
                        >Attendees List (${attendeesCount}${' '}
                        ${attendeesCount > 1 ? 'people' : 'person'})</span
                    >
                </div>
                <div class="flex flex-col gap-2 w-full mt-4 pt-2 overflow-auto">
                    ${Object.values(attendees.value)
                        .sort((a, b) => {
                            let aScore = 0;
                            let bScore = 0;

                            if (a.isHost) aScore += 10;
                            if (a.hasCamera) aScore += 5;
                            if (a.raisedHand) aScore += 1;

                            if (b.isHost) bScore += 10;
                            if (b.hasCamera) bScore += 5;
                            if (b.raisedHand) bScore += 1;

                            return bScore - aScore;
                        })
                        .map((attendee, i) => {
                            return html`<${Participant}
                                key=${attendee.userId}
                                participant=${attendee}
                            />`;
                        })}
                </div>
            </div>
        </div>
    `;
};

export const AttendeesBottomSheet = () => {
    return html`<${BottomSheet}
        open=${isAttendeesOpen.value}
        onClose=${toggleAttendees}
        class="block sm:hidden"
        title="Attendees List (${Object.values(attendees.value).length})"
    >
        <div class="w-full h-full flex gap-3 pb-6 flex-col">
            <div class="flex flex-col gap-2 w-full mt-4 pt-2 overflow-auto">
                ${Object.values(attendees.value)
                    .sort((a, b) => {
                        let aScore = 0;
                        let bScore = 0;

                        if (a.isHost) aScore += 10;
                        if (a.hasCamera) aScore += 5;
                        if (a.raisedHand) aScore += 1;

                        if (b.isHost) bScore += 10;
                        if (b.hasCamera) bScore += 5;
                        if (b.raisedHand) bScore += 1;

                        return bScore - aScore;
                    })
                    .map((attendee, i) => {
                        return html`<${Participant}
                            key=${attendee.userId}
                            participant=${attendee}
                        />`;
                    })}
            </div>
        </div>
    <//>`;
};
