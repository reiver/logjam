import { html } from 'htm'
import { Stage, streamers } from './Stage.js'
import { Attendees, Container } from 'components'

export { streamers }

export const MeetingBody = () => {
  return html`<${Container} class="relative h-full flex-grow overflow-hidden mt-4 flex items-center">
    <${Stage} />
    <$${Attendees} />
  <//>`
}
