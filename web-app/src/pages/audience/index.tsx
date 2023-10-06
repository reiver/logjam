import Meeting from 'pages/Meeting'
import { useState } from 'preact/compat'

export const AudiencePage = ({ params: { room } }: { params?: { room?: string } }) => {
  const [name, setName] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    if (e.target.elements[0].value) setName(e.target.elements[0].value)
  }
  if (!name)
    return (
      <form class="flex flex-col gap-3 max-w-[250px]" onSubmit={handleSubmit}>
        <label htmlFor="roomName">Your Name</label>
        <input name="roomName" />
        <button>Join the meeting</button>
      </form>
    )

  if (name)
    return (
      <Meeting
        params={{
          room,
          name,
        }}
      />
    )
}

export default AudiencePage
