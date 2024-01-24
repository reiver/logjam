import { zodResolver } from '@hookform/resolvers/zod'
import { Button, FormControl, TextField } from '@mui/material'
import CopyIcon from 'assets/icons/Copy.svg?react'
import LinkIcon from 'assets/icons/Link.svg?react'
import copy from 'clipboard-copy'
import clsx from 'clsx'
import { Icon, ResponsiveModal, Tooltip } from 'components'
import Meeting from 'pages/Meeting'
import { lazy } from 'preact-iso'
import { useState } from 'preact/compat'
import { useForm } from 'react-hook-form'
import { makePreviewDialog } from 'components/Dialog'
import z from 'zod'
import { parse } from 'postcss'

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

const setCustomCssContent = (event, setContentCallback) => {
  const file = event.target.files[0];


  if (file) {
    const reader = new FileReader();

    reader.onload = (e) => {
      const content = e.target.result;
      setContentCallback(content);
    };

    reader.readAsText(file);
  } else {
    // Handle the case where no file is selected
    setContentCallback(null);
  }
};

var customStyles = null;

const handleCssFileUpload = async (event) => {
  setCustomCssContent(event, (content) => {
    console.log("Content: ", content)
    if (content) {
      customStyles = content
      return
    }
  });
};




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
      <div class="w-full flex justify-center items-center px-4 min-h-full">
        <div class="w-full max-w-[500px] mx-auto mt-10 border rounded-md border-gray-300">
          <form class="flex flex-col w-full " onSubmit={form.handleSubmit(onSubmit)}>
            <span className="text-bold-12 text-black block text-center pt-5">New Live Room</span>
            <hr className="my-3" />
            <div className="p-5 flex flex-col gap-5">
              <div className="flex flex-col gap-3">
                <span class="text-bold-12 text-gray-2">Please enter your display name and room info:</span>
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
              </div>
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
              <FormControl className="w-full">
                <TextField
                  variant="outlined"
                  size="small"
                  type="file"
                  onChange={(event) => handleCssFileUpload(event)}
                />
              </FormControl>
              <div class="flex gap-2 w-full flex-col-reverse md:flex-row">
                <Button onClick={handleCreateLink} variant="outlined" className="w-full normal-case" sx={{ textTransform: 'none' }}>
                  Create Link
                </Button>
                <Button type="submit" variant="contained" className="w-full normal-case" sx={{ textTransform: 'none' }} color="primary">
                  Start Now
                </Button>
              </div>
            </div>
          </form>
          <ResponsiveModal open={showModal} onClose={setShowModal.bind(null, false)}>
            <span className="text-bold-12 text-black block text-center pt-5">Room Links</span>
            <hr className="mt-4 mb-1 border-white md:border-gray-0" />
            <div className="p-5 flex flex-col gap-5 pb-6">
              <span class="text-bold-12 text-gray-2">Copy and use host’s link for yourself, and audience link for sending to others:</span>
              <LinkCopyComponent title="Host's Link:" link={generateHostUrl('@' + form.getValues('displayName'))} />
              <LinkCopyComponent title="Audience’s Link:" link={generateAudienceUrl(form.getValues('room'))} />
            </div>
          </ResponsiveModal>
        </div>
      </div>
    )

  if (started) {
    return (
      <Meeting
        params={{
          ...form.getValues(),
          displayName: `@${form.getValues('displayName')}`,
          name: `${form.getValues('displayName')}`,
          customStyles: customStyles
        }}
      />
    )
  }
}
export default HostPage

export const LinkCopyComponent = ({ title, link, className }) => {
  const [copyTooltipTitle, setCopyTooltipTitle] = useState('Copy Link')
  const onCopy = () => {
    copy(link).then(() => {
      setCopyTooltipTitle('Coppied')
      setTimeout(() => {
        setCopyTooltipTitle('Copy Link')
      }, 2000)
    })
  }
  return (
    <div class={clsx('flex flex-col gap-1 w-full', className)}>
      {title && <span class="text-bold-12 text-gray-3">{title}</span>}
      <div className="greatape-meeting-link-background dark:bg-gray-2 dark:text-gray-0 w-full bg-gray-0 px-4 py-2 text-gray-2 flex justify-between rounded-full items-center">
        <div className="flex gap-2 items-center overflow-hidden">
          <Icon icon={LinkIcon} class="greatape-meeting-link" />
          <span class="text-medium-12 truncate greatape-meeting-link">{link}</span>
        </div>
        <Tooltip label={copyTooltipTitle} hideOnClick={false}>
          <button class="cursor-pointer" onClick={onCopy}>
            <Icon icon={CopyIcon} class="greatape-meeting-link" />
          </button>
        </Tooltip>
      </div>
    </div>
  )
}


