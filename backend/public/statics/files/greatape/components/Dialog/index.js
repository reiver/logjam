import { html, render } from 'htm'
import { Button } from 'components'
import { signal } from '@preact/signals'
import { v4 as uuidv4 } from 'uuid'

const dialogs = signal([])

export const ConfirmDialog = ({ onOk, onClose, message }) => {
  return html`<div class="py-4 px-6 dark:bg-black bg-white rounded-md dark:text-[#fefefe] light:text-[#1a1a1]">
    <div class="text-center" dangerouslySetInnerHTML=${{ __html: message }}></div>
    <div class="flex justify-end mt-4 gap-2">
      <${Button} variant="red" onClick=${onClose}>Deny<//>
      <${Button} onClick=${onOk}>Accept<//>
    </div>
  </div>`
}

export const ErrorDialog = ({ onOk, onClose, message }) => {
  return html`<div class="py-4 px-6 dark:bg-red-600 bg-red-300 rounded-md dark:text-[#fefefe] light:text-[#1a1a1]">
    <div class="text-center" dangerouslySetInnerHTML=${{ __html: message }}></div>
  </div>`
}

export const SuccessDialog = ({ onOk, onClose, message }) => {
  return html`<div class="py-4 px-6 dark:bg-green-600 bg-green-300 rounded-md dark:text-[#fefefe] light:text-[#1a1a1]">
    <div class="text-center" dangerouslySetInnerHTML=${{ __html: message }}></div>
  </div>`
}

export const InfoDialog = ({ onOk, onClose, message }) => {
  return html`<div class="py-4 px-6 dark:bg-black bg-white rounded-md dark:text-[#fefefe] light:text-[#1a1a1]">
    <div class="text-center" dangerouslySetInnerHTML=${{ __html: message }}></div>
  </div>`
}

export const DialogPool = () => {
  return html`<div class="absolute top-20 left-1/2 transform -translate-x-1/2 flex flex-col gap-3">
    ${Object.values(dialogs.value).map((dialog) => {
      if (dialog.type === 'confirm') return html`<${ConfirmDialog} ...${dialog} />`
      if (dialog.type === 'success') return html`<${SuccessDialog} ...${dialog} />`
      if (dialog.type === 'error') return html`<${ErrorDialog} ...${dialog} />`
      if (dialog.type === 'info') return html`<${InfoDialog} ...${dialog} />`
    })}
  </div>`
}

export const makeDialog = (type, message, onOk, onClose) => {
  const id = uuidv4()
  const destroy = () => {
    const dialogsTmp = { ...dialogs.value }
    delete dialogsTmp[id]
    dialogs.value = dialogsTmp
  }
  if (type !== 'confirm') {
    setTimeout(destroy, 4000)
  }
  dialogs.value = {
    ...dialogs.value,
    [id]: {
      id,
      type,
      message,
      onOk: () => {
        onOk()
        destroy()
      },
      onClose: () => {
        onClose()
        destroy()
      },
    },
  }
}

export const ToastProvider = () => {
  return html`<div id="toast-provider">
    <${DialogPool} />
  </div>`
}
