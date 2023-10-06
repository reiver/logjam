import Meeting from 'pages/Meeting'
import { lazy } from 'preact-iso'
import { useState } from 'preact/compat'
const PageNotFound = lazy(() => import('../_404'))

export const HostPage = ({ params: { displayName } }: { params?: { displayName?: string } }) => {
  const [room, setRome] = useState('')
  if (displayName && room) {
    if (displayName[0] !== '@') return <PageNotFound />
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    if (e.target.elements[0].value) setRome(e.target.elements[0].value)
  }
  if (!room)
    return (
      <form class="flex flex-col gap-3 max-w-[250px]" onSubmit={handleSubmit}>
        <label htmlFor="roomName">Room Name</label>
        <input name="roomName" />
        <button>start the meeting</button>
      </form>
    )

  if (room)
    return (
      <Meeting
        params={{
          displayName,
          room,
          name: displayName.replace('@', ''),
        }}
      />
    )
}
export default HostPage
