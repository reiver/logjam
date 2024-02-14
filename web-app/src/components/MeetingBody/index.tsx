import { Attendees, AttendeesBottomSheet, BottomBarBottomSheet, Container } from 'components'
import { Stage, streamers } from './Stage.js'
import { useState, useEffect } from 'preact/compat'

export { streamers }

export const MeetingBody = ({ customStyles }) => {

  useEffect(() => {
    if (customStyles) {
      // Create a style element and append it to the head of the document
      const styleElement = document.createElement('style');
      document.head.appendChild(styleElement);

      // Set the CSS content of the style element
      styleElement.textContent = customStyles;
    }
  }, [])

  return (
    <>
      <Container class="relative h-full flex-grow overflow-hidden py-4 flex items-center"  id="meeting-body">
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
