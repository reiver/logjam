import { signal } from '@preact/signals'
import Camera from 'assets/icons/Camera.svg?react'
import CameraLight from 'assets/icons/CameraLight.svg?react'
import CameraOff from 'assets/icons/CameraOff.svg?react'
import Close from 'assets/icons/Close.svg?react'
import Headphone from 'assets/icons/Headphone.svg?react'
import Microphone from 'assets/icons/Microphone.svg?react'
import MicrophoneLight from 'assets/icons/MicrophoneLight.svg?react'
import MicrophoneOff from 'assets/icons/MicrophoneOff.svg?react'
import Settings from 'assets/icons/Settings.svg?react'
import Smartphone from 'assets/icons/Smartphone.svg?react'
import { clsx } from 'clsx'
import { Button, Icon, IconButton, Tooltip } from 'components'
import { currentUser, sparkRTC, updateUser } from 'pages/log'
import { Fragment } from 'preact'
import { useEffect, useRef } from 'preact/compat'
import { v4 as uuidv4 } from 'uuid'
import { IODevices } from '../../lib/io-devices.js'

const dialogs = signal([])
var selectedMic = signal(null)
var selectedSpeaker = signal(null)
var selectedCamera = signal(null)
const builtInLabel = 'Built-in'
const builtInThisDevice = 'Built-in (This Device)'

export const IOSettingsDialog = ({
  onOk,
  onClose,
  message: { message, title },
  okText = 'Sounds Good',
  cancelText = 'Cancel',
  okButtonVariant = 'solid',
  onReject = onClose,
  showButtons = true,
  className,
}) => {
  const selectAudioOutputDevice = async () => {
    const io = new IODevices()
    await io.initDevices()
    const devices = io.getAudioOutputDevices()
    console.log('Audio Output Devices: ', devices)
    makeIODevicesDialog(
      'io-devices',
      {
        message: 'Please choose your "Audio output":',
        title: 'Audio',
      },
      devices,
      'speaker',
      (device) => {
        setTimeout(() => {
          if (device && device.label) {
            selectedSpeaker.value = device
          }
        }, 100)
      } //on close
    )
  }

  const selectAudioInputDevice = async () => {
    const io = new IODevices()
    await io.initDevices()
    const devices = io.getAudioInputDevices()
    console.log('Audio Input Devices: ', devices)

    makeIODevicesDialog(
      'io-devices',
      {
        message: 'Please choose your "Audio input"-Microphone:',
        title: 'Microphone',
      },
      devices,
      'microphone',
      (device) => {
        setTimeout(() => {
          if (device && device.label) {
            selectedMic.value = device
          }
        }, 100)
      } //on close
    )
  }

  const selectVideoInputDevice = async () => {
    const io = new IODevices()
    await io.initDevices()
    const devices = io.getVideoInputDevices()
    console.log('Video Input Devices: ', devices)

    makeIODevicesDialog(
      'io-devices',
      {
        message: 'Please choose your "Video input":',
        title: 'Video',
      },
      devices,
      'camera',
      (device) => {
        setTimeout(() => {
          if (device && device.label) {
            selectedCamera.value = device
          }
        }, 100)
      } //on close
    )
  }

  const isIphone = () => {
    const userAgent = navigator.userAgent
    if (userAgent.match(/iPhone|iPad|iPod/i)) {
      return true
    }
    return false
  }

  console.log('Resetting..')

  return (
    <div class="absolute top-0 left-0 w-full h-full">
      <div class="z-20 absolute w-full h-full bg-black bg-opacity-60" />
      <div
        class={clsx(
          className,
          'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full sm:max-w-[400px] sm:border dark:border-gray-1 border-gray-0'
        )}
      >
        <div class="flex justify-center items-center p-5 relative">
          <span class="dark:text-white text-black text-bold-12">{title}</span>
          <Icon icon={Close} class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer" onClick={onClose} />
        </div>
        <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />

        <div class="sm:pt-8 pt-4 sm:pb-8 pb-4">
          {!isIphone() && (
            <div class="sm:py-4 py-2 flex rounded-md mx-2 cursor-pointer" onClick={selectAudioOutputDevice}>
              <div class="text-left text-bold-12 px-5 flex-1">Audio Output</div>
              <div id="selectedSpeaker" class="text-right text-bold-12 px-5 flex-1 text-gray-1 cursor-pointer">
                {selectedSpeaker.value ? selectedSpeaker.value.label : builtInLabel}
              </div>
            </div>
          )}

          <div class="sm:py-4 py-2 rounded-md mx-2 flex cursor-pointer" onClick={selectAudioInputDevice}>
            <div class="text-left text-bold-12 px-5 flex-1">Microphone</div>
            <div id="selectedMic" class="text-right text-bold-12 px-5 flex-1 text-gray-1">
              {selectedMic.value ? selectedMic.value.label : builtInLabel}
            </div>
          </div>

          <div class="sm:py-4 py-2 rounded-md mx-2 flex cursor-pointer" onClick={selectVideoInputDevice}>
            <div class="text-left text-bold-12 px-5 flex-1">Video Input</div>
            <div id="selectedCamera" class="text-right text-bold-12 px-5 flex-1 text-gray-1">
              {selectedCamera.value ? selectedCamera.value.label : builtInLabel}
            </div>
          </div>
        </div>

        {showButtons && (
          <div class="flex justify-end gap-2 p-5 pt-0">
            <Button
              size="lg"
              variant="outline"
              class="w-full flex-grow-1"
              onClick={() => {
                onReject && onReject()
                onClose()
              }}
            >
              {cancelText}
            </Button>
            <Button
              size="lg"
              variant={okButtonVariant}
              class="w-full flex-grow-1"
              onClick={() => {
                onOk(selectedMic.value, selectedCamera.value, selectedSpeaker.value)
              }}
            >
              {okText}
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}

export const IODevicesDialog = ({ onClose, message: { message, title }, devices, deviceType, className, contentClassName }) => {
  let selectedDeviceIndex = -1
  const handleDeviceClick = (index = -1, vanish = true) => {
    //mark built-in / default devices checked by defualt on initial load
    if (index === -1) {
      devices.forEach((elem, _index) => {
        if (deviceType === 'microphone' && elem.kind === 'audioinput' && (elem.label.toLowerCase().includes('default') || elem.label.toLowerCase().includes('iphone microphone'))) {
          index = _index
        } else if (deviceType === 'speaker' && elem.kind === 'audiooutput' && elem.label.toLowerCase().includes('default')) {
          index = _index
        } else if (
          deviceType === 'camera' &&
          elem.kind === 'videoinput' &&
          (elem.label.toLowerCase().includes('default') || elem.label.toLowerCase().includes('(') || elem.label.toLowerCase().includes('front'))
        ) {
          index = _index
        }
      })

      console.log('Updated Index: ', index, ' ,value: ', devices[index])
    }

    // Now, the 'index' variable will contain the index of the matching device (or -1 if none found).

    // Check if the clicked device is already selected
    if (selectedDeviceIndex === index) {
      // If it's already selected, deselect it by setting the selectedDeviceIndex to -1
      selectedDeviceIndex = -1
    } else {
      // If it's not selected, select it by setting the selectedDeviceIndex to the clicked index
      selectedDeviceIndex = index
    }

    console.log('handleDeviceClick: ', index)
    // Update the radio button's checked attribute based on the selectedDeviceIndex
    const radioInput = document.getElementById(`device${index}`) as HTMLInputElement
    console.log('radioInput: ', radioInput)
    if (radioInput) {
      radioInput.checked = selectedDeviceIndex === index
      radioInput.style.accentColor = 'black'
      if (vanish) {
        setTimeout(() => {
          console.log('Selected Device: ', devices[selectedDeviceIndex])
          onClose(devices[selectedDeviceIndex])
        }, 200)
      }
    }
  }

  // mark default device selected on inital display
  setTimeout(() => {
    if (deviceType === 'speaker' && !selectedSpeaker.value) {
      handleDeviceClick(-1, false)
    } else if (deviceType === 'microphone' && !selectedMic.value) {
      handleDeviceClick(-1, false)
    } else if (deviceType === 'camera' && !selectedCamera.value) {
      handleDeviceClick(-1, false)
    }
  }, 50)

  //check selected devices and make it selected on radio button too
  if (devices && devices.length > 0) {
    devices.forEach((value, index) => {
      console.log('Device: ', value, ' index: ', index)
      if (value.kind === 'audioinput') {
        if (selectedMic.value && selectedMic.value.deviceId === value.deviceId) {
          console.log('selectedMicis: ', selectedMic.value)
          setTimeout(() => {
            handleDeviceClick(index, false)
          }, 250)
        }
      } else if (value.kind === 'videoinput') {
        if (selectedCamera.value && selectedCamera.value.deviceId === value.deviceId) {
          console.log('selectedCamis: ', selectedCamera.value)

          setTimeout(() => {
            handleDeviceClick(index, false)
          }, 250)
        }
      } else if (value.kind === 'audiooutput') {
        if (selectedSpeaker.value && selectedSpeaker.value.deviceId === value.deviceId) {
          console.log('selectedSpeakeris: ', selectedSpeaker.value)

          setTimeout(() => {
            handleDeviceClick(index, false)
          }, 250)
        }
      }
    })
  }

  const isBuiltInDevice = (deviceType, device) => {
    switch (deviceType) {
      case 'speaker':
        if (device.label.toLowerCase().includes('default')) {
          return true
        }
        break
      case 'microphone':
        if (device.label.toLowerCase().includes('default') || device.label.toLowerCase().includes('iphone microphone')) {
          return true
        }
        break
      case 'camera':
        if (device.label.toLowerCase().includes('(') || device.label.toLowerCase().includes('front')) {
          return true
        }
        break
      default:
        // Handle other cases if needed
        break
    }

    return false
  }
  console.log(devices)

  return (
    <div class="absolute top-0 left-0 w-full h-full">
      <div class="z-20 absolute w-full h-full bg-black bg-opacity-60" />
      <div
        class={clsx(
          className,
          'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full sm:max-w-[400px] sm:border dark:border-gray-1 border-gray-0'
        )}
      >
        <div class="flex justify-center items-center p-5 relative">
          <span class="dark:text-white text-black text-bold-12">{title}</span>
          <Icon icon={Close} class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer" onClick={onClose} />
        </div>
        <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />
        <div class={clsx(contentClassName, 'text-left text-bold-12 sm:pt-8 pt-5 p-5')} dangerouslySetInnerHTML={{ __html: message }}></div>

        <form>
          <div class="sm:pb-4 pb-2">
            {devices.map((device, index) => (
              <div class="sm:py-4 py-2 rounded-md mx-2 flex items-center cursor-pointer" onClick={() => handleDeviceClick(index)}>
                <Icon
                  icon={
                    isBuiltInDevice(deviceType, device)
                      ? Smartphone
                      : deviceType === 'microphone'
                      ? MicrophoneLight
                      : deviceType === 'camera'
                      ? CameraLight
                      : deviceType === 'speaker'
                      ? Headphone
                      : Fragment
                  }
                  class="ml-5"
                  width="20px"
                  height="20px"
                />
                <div class="text-left px-2 text-bold-12 flex-1">{isBuiltInDevice(deviceType, device) ? builtInThisDevice : device.label}</div>
                <label class="flex items-right px-5 flex-0">
                  <input type="radio" name="devices" id={`device${index}`} />
                </label>
              </div>
            ))}
          </div>
        </form>
      </div>
    </div>
  )
}

export const PreviewDialog = ({
  onOk,
  onClose,
  videoStream,
  message: { message, title, yesButton, noButton },
  okText = yesButton ? yesButton : 'Sounds Good',
  cancelText = noButton ? noButton : 'Cancel',
  okButtonVariant = 'solid',
  onReject = onClose,
  showButtons = true,
  className,
  contentClassName,
}) => {
  selectedCamera.value = null
  selectedMic.value = null
  selectedSpeaker.value = null

  const videoRef = useRef<HTMLVideoElement>()
  const { hasCamera, hasMic, isCameraOn, isMicrophoneOn } = currentUser.value

  useEffect(() => {
    videoRef.current.srcObject = videoStream
  }, [videoStream])

  useEffect(() => {
    videoRef.current.playsInline = true

    const isMobile = window.parent.outerWidth <= 400 && window.parent.outerHeight <= 850
    if (isMobile) {
      //set style
      // videoRef.current.style = 'width: 100%; height: 56.25vw; max-height: 100%;'
      videoRef.current.setAttribute('style', 'width: 100%; height: 56.25vw; max-height: 100%;')
    }
  }, [])

  const toggleCamera = () => {
    sparkRTC.value.disableVideo(!isCameraOn)
    updateUser({
      isCameraOn: !isCameraOn,
    })
  }

  const toggleMicrophone = () => {
    sparkRTC.value.disableAudio(!isMicrophoneOn)
    updateUser({
      isMicrophoneOn: !isMicrophoneOn,
    })
  }

  const openDeviceSettings = () => {
    makeIOSettingsDialog(
      'io-settings',
      {
        message: '',
        title: 'Input and Output Settings',
      },
      async (mic, cam, speaker) => {
        console.log('mic: ', mic, 'cam: ', cam, 'speaker: ', speaker)

        //now change the Audio, Video and Speaker devices
        const stream = await sparkRTC.value.changeIODevices(mic, cam, speaker)

        console.log('New Stream: ', stream)
        videoStream = stream
        if (videoRef.current) {
          videoRef.current.srcObject = stream
        } else {
          console.log('No video ref')
        }
      }, //ok
      () => {} //close
    )
  }

  return (
    <div class="absolute top-0 left-0 w-full h-full">
      <div class="z-10 absolute w-full h-full bg-black bg-opacity-60" />
      <div
        class={clsx(
          className,
          'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full sm:max-w-[40%] sm:border dark:border-gray-1 border-gray-0'
        )}
      >
        <div class="flex justify-center items-center p-5 relative">
          <span class="dark:text-white text-black text-bold-12">{title}</span>
          <Icon icon={Close} class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer" onClick={onClose} />
        </div>
        <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />
        <div class={clsx(contentClassName, 'text-left text-bold-12 sm:py-8 py-5 p-5')} dangerouslySetInnerHTML={{ __html: message }}></div>
        <div class="px-5 relative">
          <video ref={videoRef} autoplay playsinline muted={true} className="aspect-video object-cover rounded-lg" />
          <div class={clsx('h-[48px] absolute top-1 right-6 gap-2 flex justify-center items-center')}>
            {!isMicrophoneOn && (
              <div className="pr-2">
                <Icon icon={MicrophoneOff} width="20px" height="20px" />
              </div>
            )}
          </div>
        </div>
        <div class="py-4 px-5 flex gap-3 justify-center">
          {hasCamera && (
            <Tooltip label={!isCameraOn ? 'Turn Camera On' : 'Turn Camera Off'}>
              <IconButton variant={!isCameraOn && 'danger'} onClick={toggleCamera}>
                {' '}
                <Icon icon={!isCameraOn ? CameraOff : Camera} />{' '}
              </IconButton>
            </Tooltip>
          )}
          {hasMic && (
            <Tooltip label={!isMicrophoneOn ? 'Turn Microphone On' : 'Turn Microphone Off'}>
              <IconButton variant={!isMicrophoneOn && 'danger'} onClick={toggleMicrophone}>
                <Icon icon={!isMicrophoneOn ? MicrophoneOff : Microphone} />
              </IconButton>
            </Tooltip>
          )}
          {hasCamera && hasMic && (
            <Tooltip key="Settings" label="Settings">
              <IconButton onClick={openDeviceSettings}>
                <Icon icon={Settings} />
              </IconButton>
            </Tooltip>
          )}
        </div>

        {showButtons && (
          <div class="flex justify-end gap-2 p-5 pt-0">
            <Button
              size="lg"
              variant="outline"
              class="w-full flex-grow-1"
              onClick={() => {
                onReject && onReject()
              }}
            >
              {cancelText}
            </Button>
            <Button size="lg" variant={okButtonVariant} class="w-full flex-grow-1" onClick={onOk}>
              {okText}
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}

export const InviteDialog = ({
  onOk,
  onClose,
  message: { message, title },
  okText = 'Send Request',
  cancelText = 'Not Now!',
  okButtonVariant = 'solid',
  onReject = onClose,
  showButtons = true,
  className,
  contentClassName,
}) => {
  return (
    <div class="absolute top-0 left-0 w-full h-full">
      <div class="z-10 absolute w-full h-full bg-black bg-opacity-60" onClick={onClose} />
      <div
        class={clsx(
          className,
          'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full sm:max-w-[400px] sm:border dark:border-gray-1 border-gray-0'
        )}
      >
        <div class="flex justify-center items-center p-5 relative">
          <span class="dark:text-white text-black text-bold-12">{title}</span>
          <Icon icon={Close} class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer" onClick={onClose} />
        </div>
        <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />
        <div class={clsx(contentClassName, 'text-left text-bold-12 sm:py-8 py-5 p-5')} dangerouslySetInnerHTML={{ __html: message }}></div>
        {showButtons && (
          <div class="flex justify-end gap-2 p-5 pt-0">
            <Button
              size="lg"
              variant="outline"
              class="w-full flex-grow-1"
              onClick={() => {
                onReject && onReject()
                onClose()
              }}
            >
              {cancelText}
            </Button>
            <Button size="lg" variant={okButtonVariant} class="w-full flex-grow-1" onClick={onOk}>
              {okText}
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}

export const ConfirmDialog = ({
  onOk,
  onClose,
  message: { message, title },
  okText = 'Accept',
  cancelText = 'Reject',
  okButtonVariant = 'solid',
  onReject = onClose,
  showButtons = true,
  className,
  contentClassName,
}) => {
  return (
    <div class="absolute top-0 left-0 w-full h-full">
      <div class="z-10 absolute w-full h-full bg-black bg-opacity-60" onClick={onClose} />
      <div
        class={clsx(
          className,
          'absolute -translate-y-full z-20 top-full left-0 right-0 sm:right-unset sm:top-1/2 sm:left-1/2 transform sm:-translate-x-1/2 sm:-translate-y-1/2 dark:bg-gray-3 dark:text-gray-0 bg-white text-gray-2 sm:rounded-lg rounded-t-lg w-full sm:max-w-[400px] sm:border dark:border-gray-1 border-gray-0'
        )}
      >
        <div class="flex justify-center items-center p-5 relative">
          <span class="dark:text-white text-black text-bold-12">{title}</span>
          <Icon icon={Close} class="absolute top-1/2 sm:right-5 right-[unset] left-5 sm:left-[unset] transform -translate-y-1/2 cursor-pointer" onClick={onClose} />
        </div>
        <hr class="dark:border-gray-2 border-gray-0 sm:block hidden" />
        <div class={clsx(contentClassName, 'text-left text-bold-12 sm:py-8 py-5 p-5')} dangerouslySetInnerHTML={{ __html: message }}></div>
        {showButtons && (
          <div class="flex justify-end gap-2 p-5 pt-0">
            <Button
              size="lg"
              variant="outline"
              class="w-full flex-grow-1"
              onClick={() => {
                onReject && onReject()
                onClose()
              }}
            >
              {cancelText}
            </Button>
            <Button size="lg" variant={okButtonVariant} class="w-full flex-grow-1" onClick={onOk}>
              {okText}
            </Button>
          </div>
        )}
      </div>
    </div>
  )
}

export const InfoDialog = ({ onOk, onClose, message: { message, icon, variant }, pointer }) => {
  return (
    <div
      class={clsx('select-none py-4 px-6 flex justify-between items-center text-medium-12 min-w-full sm:min-w-[350px] rounded-md', {
        'cursor-pointer': pointer,
        'bg-red-distructive text-white-f-9': variant === 'danger',
        'dark:bg-white-f-9 dark:text-gray-3 bg-gray-3 text-white-f-9 border dark:border-gray-1 border-gray-0': !variant,
      })}
      onClick={onClose}
    >
      <div class="text-left" dangerouslySetInnerHTML={{ __html: message }}></div>
      <Icon icon={icon} width="24" height="24" />
    </div>
  )
}

export const DialogPool = () => {
  return (
    <>
      <div className="absolute right-0 left-0 md:left-[unset] md:right-10 bottom-[5.5rem] flex flex-col justify-end gap-2 px-4 sm:px-0">
        {Object.values(dialogs.value).map((dialog) => {
          if (dialog.type === 'info') return <InfoDialog {...dialog} />
        })}
      </div>

      {Object.values(dialogs.value).map((dialog) => {
        if (dialog.type === 'confirm') return <ConfirmDialog {...dialog} />
        else if (dialog.type === 'preview') return <PreviewDialog {...dialog} />
        else if (dialog.type === 'io-settings') return <IOSettingsDialog {...dialog} />
        else if (dialog.type === 'io-devices') return <IODevicesDialog {...dialog} />
        else if (dialog.type === 'invite') return <InviteDialog {...dialog} />
      })}
    </>
  )
}

export const makeInviteDialog = (type, message, onOk, onClose, options = {}) => {
  const id = uuidv4()
  const destroy = () => {
    const dialogsTmp = { ...dialogs.value }
    delete dialogsTmp[id]
    dialogs.value = dialogsTmp
  }

  dialogs.value = {
    ...dialogs.value,
    [id]: {
      id,
      type,
      message,
      pointer: !!onClose,
      onOk: () => {
        onOk && onOk()
        destroy()
      },
      onClose: () => {
        onClose && onClose()
        destroy()
      },
      ...options,
    },
  }
}

export const makeDialog = (type, message, onOk = undefined, onClose = undefined, options = {}) => {
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
      pointer: !!onClose,
      onOk: () => {
        onOk && onOk()
        destroy()
      },
      onClose: () => {
        onClose && onClose()
        destroy()
      },
      ...options,
    },
  }
}

export const destroyDialog = (id) => {
  const dialogsTmp = { ...dialogs.value }
  delete dialogsTmp[id]
  dialogs.value = dialogsTmp
}

export const makePreviewDialog = (type, videoStream, message, onOk, onClose, options = {}) => {
  const id = uuidv4()
  const destroy = () => {
    const dialogsTmp = { ...dialogs.value }
    delete dialogsTmp[id]
    dialogs.value = dialogsTmp
  }

  dialogs.value = {
    ...dialogs.value,
    [id]: {
      id,
      type,
      videoStream,
      message,
      pointer: !!onClose,
      onOk: () => {
        onOk && onOk()
        destroy()
      },
      onClose: () => {
        onClose && onClose()
        destroy()
      },
      ...options,
    },
  }

  return id
}

export const makeIODevicesDialog = (type, message, devices, deviceType, onClose, options = {}) => {
  const id = uuidv4()
  const destroy = () => {
    const dialogsTmp = { ...dialogs.value }
    delete dialogsTmp[id]
    dialogs.value = dialogsTmp
  }

  dialogs.value = {
    ...dialogs.value,
    [id]: {
      id,
      type,
      message,
      devices,
      deviceType,
      onOk: () => {
        destroy()
      },
      onClose: (device) => {
        onClose && onClose(device)
        destroy()
      },
      ...options,
    },
  }

  return id
}

export const makeIOSettingsDialog = (type, message, onOk, onClose, options = {}) => {
  const id = uuidv4()
  const destroy = () => {
    const dialogsTmp = { ...dialogs.value }
    delete dialogsTmp[id]
    dialogs.value = dialogsTmp
  }

  dialogs.value = {
    ...dialogs.value,
    [id]: {
      id,
      type,
      message,
      pointer: !!onClose,
      onOk: (mic, cam, speaker) => {
        onOk && onOk(mic, cam, speaker)
        destroy()
      },
      onClose: () => {
        onClose && onClose()
        destroy()
      },
      ...options,
    },
  }

  return id
}

export const ToastProvider = () => {
  return (
    <div id="toast-provider">
      <DialogPool />
    </div>
  )
}
