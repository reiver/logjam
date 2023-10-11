import { Attendees, AttendeesBottomSheet, BottomBarBottomSheet, Container } from 'components'
import { Stage, streamers } from './Stage.js'

export { streamers }

export const MeetingBody = () => {
  return (
    <>
      <Container class="relative h-full flex-grow overflow-hidden py-4 flex items-center" id="meeting-body">
        <Stage />
        <Attendees />
      </Container>
      <div class="sm:hidden">
        <BottomBarBottomSheet />
        <AttendeesBottomSheet />
      </div>
    </>
  )
}
