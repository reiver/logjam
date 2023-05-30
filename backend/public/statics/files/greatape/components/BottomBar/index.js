import { html } from 'htm'
import { Icon, Container, Controllers, toggleAttendees, attendees } from 'components'

export const BottomBar = () => {
  return html`<${Container}>
    <div class="w-full grid grid-cols-12 dark:bg-secondary-1-a py-3 dark:text-gray-0 text-gray-2">
      <div class="col-span-3">
        <div class="h-full flex items-center">
          <div class="flex items-center gap-2 dark:bg-gray-2  bg-gray-0 rounded-full px-4 py-1">
            <${Icon} icon="Link" />
            <span>https://great.ape/r/thu-kiu</span>
            <div class="ml-3">
              <${Icon} icon="Copy" />
            </div>
          </div>
        </div>
      </div>
      <div class="col-span-6 flex items-center justify-center">
        <${Controllers} icon="Camera" />
      </div>
      <div class="col-span-3 text-right">
        <div class="h-full flex items-center justify-end">
          ${attendees.value.length
            ? html`<div
                onClick="${toggleAttendees}"
                class="transition-all select-none cursor-pointer flex items-center gap-2 rounded-md hover:bg-gray-0 hover:bg-opacity-10 hover:dark:bg-gray-2 hover:dark:bg-opacity-20 py-1 px-3"
              >
                <${Icon} icon="Avatar" />
                <span>${attendees.value.length} attendee${attendees.value.length > 1 ? 's' : ''}</span>
              </div>`
            : null}
        </div>
      </div>
    </div>
  <//>`
}

export default BottomBar
