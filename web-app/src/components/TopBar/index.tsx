import { Button, Container, Logo } from 'components'

import { leaveMeeting, meetingStatus,sparkRTC } from 'pages/Meeting'

export const TopBar = () => {
  const handleLeaveMeeting = leaveMeeting
  return (
    <div class="w-full bg-white dark:bg-black py-3" id="top-bar">
      <Container>
        <div class="grid grid-cols-12">
          <div class="col-span-3 flex items-center">
            <Logo />
          </div>
          <div class="col-span-6 flex items-center justify-center">
            {meetingStatus.value && <span class="text-black dark:text-white text-center text-bold-14 hidden sm:block">{sparkRTC.value && sparkRTC.value.role===sparkRTC.value.Roles.AUDIENCE ? 'GreatApe':'Is Your Future Distributed? Welcome to the Fediverse!'}</span>}
          </div>
          <div class="col-span-3 text-right">
            {meetingStatus.value && (
              <Button variant="red" onClick={handleLeaveMeeting}>
                Leave
              </Button>
            )}
          </div>
        </div>
      </Container>
    </div>
  )
}

export default TopBar
