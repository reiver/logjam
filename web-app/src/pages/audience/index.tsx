import { zodResolver } from '@hookform/resolvers/zod'
import { Button, FormControl, TextField } from '@mui/material'
import Meeting from 'pages/Meeting'
import { lazy } from 'preact-iso'
import { useState } from 'preact/compat'
import { useForm } from 'react-hook-form'
import z from 'zod'

const PageNotFound = lazy(() => import('../_404'))

const schema = z.object({
  room: z.string().min(1, 'This field is required'),
  name: z.string().min(1, 'This field is required'),
})

export const AudiencePage = ({ params: { room } }: { params?: { room?: string } }) => {
  const [started, setStarted] = useState(false)
  const form = useForm({
    defaultValues: {
      room: room,
      name: '',
    },
    resolver: zodResolver(schema),
  })

  const onSubmit = () => {
    setStarted(true)
  }

  if (!started)
    return (
      <div class="w-full flex justify-center items-center max-w-[500px] mx-auto mt-10 border rounded-md border-gray-300">
        <form class="flex flex-col w-full " onSubmit={form.handleSubmit(onSubmit)}>
          <span className="text-bold-12 text-black block text-center pt-5">New Live Room</span>
          <hr className="my-3" />
          <div className="p-5 flex flex-col gap-5">
            <FormControl className="w-full">
              <TextField
                label="Display name"
                variant="outlined"
                size="small"
                {...form.register('name')}
                error={!!form.formState.errors.displayName}
                helperText={form.formState.errors.displayName?.message}
              />
            </FormControl>

            <div class="flex gap-2 w-full">
              <Button variant="outlined" className="w-full normal-case" sx={{ textTransform: 'none' }}>
                Attend Live Show
              </Button>
              <Button type="submit" variant="contained" className="w-full normal-case" sx={{ textTransform: 'none' }} color="primary">
                Not now
              </Button>
            </div>
          </div>
        </form>
      </div>
    )

  if (started) {
    return (
      <Meeting
        params={{
          ...form.getValues(),
          room: room,
        }}
      />
    )
  }
}

export default AudiencePage
