import { zodResolver } from '@hookform/resolvers/zod'
import { Button, FormControl, Modal, TextField } from '@mui/material'
import CloseIcon from 'assets/icons/Close.svg?react'
import CopyIcon from 'assets/icons/Copy.svg?react'
import LinkIcon from 'assets/icons/Link.svg?react'
import copy from 'clipboard-copy'
import { Icon, IconButton, Tooltip } from 'components'
import Meeting from 'pages/Meeting'
import { lazy } from 'preact-iso'
import { useState } from 'preact/compat'
import { useForm } from 'react-hook-form'
import z from 'zod'

const PageNotFound = lazy(() => import('../_404'))

const schema = z.object({
  room: z.string().min(1, 'This field is required'),
  displayName: z.string().min(1, 'This field is required'),
  description: z.string(),
})

const generateHostUrl = (displayName: string) => {
  return `${window.location.origin}/${displayName}`
}

const generateAudienceUrl = (roomName: string) => {
  return `${window.location.origin}/log/${roomName}`
}

export const HostPage = ({ params: { displayName } }: { params?: { displayName?: string } }) => {
  const [started, setStarted] = useState(false)
  const [showModal, setShowModal] = useState(false)
  const form = useForm({
    defaultValues: {
      room: '',
      displayName: displayName.replace('@', ''),
      description: '',
    },
    resolver: zodResolver(schema),
  })
  if (displayName) {
    if (displayName[0] !== '@') return <PageNotFound />
  }

  const onSubmit = () => {
    setStarted(true)
  }

  const handleCreateLink = () => {
    form.trigger().then((v) => {
      if (v) {
        setShowModal(v)
      }
    })
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
                label="Display Name"
                variant="outlined"
                size="small"
                {...form.register('displayName')}
                error={!!form.formState.errors.displayName}
                helperText={form.formState.errors.displayName?.message}
              />
            </FormControl>
            <FormControl className="w-full">
              <TextField label="Room Name" variant="outlined" size="small" {...form.register('room')} error={!!form.formState.errors.room} helperText={form.formState.errors.room?.message} />
            </FormControl>
            <FormControl className="w-full">
              <TextField
                multiline
                rows={4}
                label="Room Description"
                variant="outlined"
                size="small"
                {...form.register('description')}
                error={!!form.formState.errors.description}
                helperText={form.formState.errors.description?.message}
              />
            </FormControl>
            <div class="flex gap-2 w-full">
              <Button onClick={handleCreateLink} variant="outlined" className="w-full normal-case" sx={{ textTransform: 'none' }}>
                Create Link
              </Button>
              <Button type="submit" variant="contained" className="w-full normal-case" sx={{ textTransform: 'none' }} color="primary">
                Start Now
              </Button>
            </div>
          </div>
        </form>
        <Modal open={showModal} onClose={setShowModal.bind(null, false)}>
          <div class="max-w-[400px] rounded-xl bg-white absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 text-black relative">
            <span className="text-bold-12 text-black block text-center pt-5">Room Links</span>
            <IconButton onClick={setShowModal.bind(null, false)} variant="nothing" className="absolute top-4 right-4">
              <Icon icon={CloseIcon} width="24px" className="text-gray-1" />
            </IconButton>
            <hr className="mt-4 mb-1" />
            <div className="p-5 flex flex-col gap-5 pb-6">
              <span class="text-bold-12 text-gray-2">Copy and use host’s link for yourself, and audience link for sending to others:</span>
              <LinkCopyComponent title="Host's Link:" link={generateHostUrl('@' + form.getValues('displayName'))} />
              <LinkCopyComponent title="Audience’s Link:" link={generateAudienceUrl(form.getValues('room'))} />
            </div>
          </div>
        </Modal>
      </div>
    )

  if (started) {
    return (
      <Meeting
        params={{
          ...form.getValues(),
          displayName: `@${form.getValues('displayName')}`,
          name: displayName.replace('@', ''),
        }}
      />
    )
  }
}
export default HostPage

const LinkCopyComponent = ({ title, link }) => {
  const [copyTooltipTitle, setCopyTooltipTitle] = useState('Copy Link')
  const onCopy = () => {
    copy(link).then(() => {
      console.log('here')
      setCopyTooltipTitle('Coppied')
      setTimeout(() => {
        setCopyTooltipTitle('Copy Link')
      }, 2000)
    })
  }
  return (
    <div class="flex flex-col gap-1">
      {title && <span class="text-bold-12 text-gray-3">{title}</span>}
      <div className="bg-gray-0 px-4 py-2 text-gray-2 flex justify-between rounded-full items-center">
        <div className="flex gap-2 items-center">
          <Icon icon={LinkIcon} />
          <span class="text-medium-12">{link}</span>
        </div>
        <Tooltip label={copyTooltipTitle} hideOnClick={false}>
          <button class="cursor-pointer" onClick={onCopy}>
            <Icon icon={CopyIcon} />
          </button>
        </Tooltip>
      </div>
    </div>
  )
}
