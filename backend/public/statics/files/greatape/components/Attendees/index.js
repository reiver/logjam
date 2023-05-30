import { signal, computed } from '@preact/signals'
import { Icon } from 'components'
import { html } from 'htm'
import clsx from 'clsx'

export const attendees = signal([
  // {
  //   name: 'Alex Suprun',
  //   isHost: true,
  //   avatar: 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1480&q=80',
  //   raisedHand: false,
  //   hasCamera: true,
  // }
])

export const isAttendeesOpen = signal(false)
export const toggleAttendees = () => {
  isAttendeesOpen.value = !isAttendeesOpen.value
}

export const attendeesWidth = computed(() => {
  if (!isAttendeesOpen.value) return 0
  return 350
})

export const Participant = ({ participant }) => {
  return html` <div
    class=${clsx('flex w-full justify-between items-center rounded-md px-2 py-1', 'cursor-pointer hover:dark:bg-white hover:dark:bg-opacity-10 hover:bg-gray-500 hover:bg-opacity-10 transition-all')}
  >
    <div class="flex items-center gap-2">
      ${participant.avatar
        ? html`<img src="${participant.avatar}" class="w-9 h-9 rounded-full object-cover" />`
        : html`<div class="dark:bg-gray-300 dark:bg-opacity-30 bg-opacity-30 bg-gray-400 rounded-full w-9 h-9 flex justify-center items-center">
            <${Icon} icon="Avatar" width="20px" height="20px" />
          </div>`}

      <div class="flex flex-col justify-center">
        <span class="text-gray-3 dark:text-white-f-9 text-bold-12">${participant.name}</span>
        ${participant.isHost && html`<span class="text-gray-1 dark:text-gray-0 text-regular-12">Host</span>`}
      </div>
    </div>
    ${(participant.raisedHand || participant.hasCamera) &&
    html`
      <div>
        <${Icon} icon=${participant.raisedHand ? 'Hand' : 'Camera'} width="25px" height="25px" />
      </div>
    `}
  </div>`
}

export const Attendees = () => {
  return html` <div
    class="${clsx(
      'h-full min-w-[350px] border rounded-lg p-2',
      'bg-white-f border-gray-0 text-secondary-1-a',
      'dark:bg-gray-3 dark:border-0 dark:text-white-f-9',
      'absolute top-0',
      'transition-all ease-in-out',
      'lg:right-10 right-4',
      {
        'translate-x-[100%] lg:-mr-10 -mr-4': !isAttendeesOpen.value,
        'translate-x-[100%]': !isAttendeesOpen.value,
      }
    )}"
  >
    <div class="flex flex-col py-2  gap-2">
      <div class="flex justify-center items-center w-full gap-2">
        <${Icon} icon="Avatar" />
        <span>Attendees List (${attendees.value.length} ${attendees.value.length > 1 ? 'people' : 'person'})</span>
      </div>
      <div class="flex flex-col gap-2 w-full mt-6">
        ${attendees.value
          .sort((attendee) => {
            if (attendee.isHost) return -1
            else 0
          })
          .map((attendee, i) => {
            return html`<${Participant} key=${i} participant=${attendee} />`
          })}
      </div>
    </div>
  </div>`
}
