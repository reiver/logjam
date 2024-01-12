import { Attendees, AttendeesBottomSheet, BottomBarBottomSheet, Container } from 'components'
import { Stage, streamers } from './Stage.js'
import { useState, useEffect } from 'preact/compat'

export { streamers }

export const MeetingBody = ({ customStyles }) => {
  const [backgroundColor, setBackgroundColor] = useState(customStyles.background.color); // Set initial color
  const [backgroundImage, setBackgroundImage] = useState(customStyles.background.image); // Set initial color

  useEffect(() => {
    console.log("image: ", customStyles.background.background),
    setBackgroundColor(customStyles.background.color)
    setBackgroundImage(customStyles.background.background)
    //background-color: ${backgroundColor}; 
  }, [])
  return (
    <>
      <Container class="relative h-full flex-grow overflow-hidden py-4 flex items-center" style={`background-image: ${backgroundImage}`} id="meeting-body">
        <Stage customStyles={customStyles} />
        <Attendees />
      </Container>
      <div class="sm:hidden">
        <BottomBarBottomSheet />
        <AttendeesBottomSheet />
      </div>
    </>
  )
}
