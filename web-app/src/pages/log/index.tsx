import { computed, signal } from '@preact/signals'
import { BottomBar, Button, MeetingBody, TopBar, attendees, attendeesBadge, isMoreOptionsOpen, makeDialog, streamers } from 'components'
import { isAttendeesOpen } from 'components/Attendees'
import { ToastProvider, destroyDialog, makePreviewDialog } from 'components/Dialog'
import { Roles, createSparkRTC, getWsUrl } from 'lib/common.js'
import { useEffect } from 'preact/compat'

export const isDebugMode = signal((new URLSearchParams(window.location.search).get('debug') || '').toLowerCase() === 'true')
export const statsDataOpen = signal(false)
export const statsData = signal('')
export const sparkRTC = signal(null)
export const meetingStatus = signal(true)
export const broadcastIsInTheMeeting = signal(true)
export const raisedHandsCount = signal(0)
export const raiseHandMaxLimitReached = computed(() => {
  return sparkRTC.value && raisedHandsCount.value === sparkRTC.value.maxRaisedHands
})
export const setUserActionLoading = (userId, actionLoading) => {
  attendees.value = {
    ...attendees.value,
    [userId]: {
      ...attendees.value[userId],
      actionLoading,
    },
  }
}
var Loading = false

// const url = `stats/index.html`;
// var targetWindow = window.open(url, '_blank');
export const currentUser = signal({
  showControllers: true,
  isHost: false,
  isMicrophoneOn: true,
  isCameraOn: true,
  isMeetingMuted: false,
  sharingScreenStream: null,
  ableToRaiseHand: true,
  hasMic: true,
  hasCamera: true,
  userId: null,
  isStreamming: false,
})

export const updateUser = (props) => {
  currentUser.value = {
    ...currentUser.value,
    ...props,
  }
}

export const onStartShareScreen = (stream) => {
  log(`ScreenShareStram: ${stream}`)

  stream.getTracks()[0].onended = async () => {
    await sparkRTC.value.stopShareScreen(stream)
    updateUser({
      sharingScreenStream: null,
    })
    onStopStream(stream)
  }

  isMoreOptionsOpen.value = false

  streamers.value = {
    ...streamers.value,
    [stream.id]: {
      name: stream.name,
      isHost: true,
      avatar: '',
      raisedHand: false,
      hasCamera: false,
      stream,
      isShareScreen: true,
    },
  }
}

const displayStream = async (stream, toggleFull = false) => {
  let local = false
  if (sparkRTC.value.localStream) {
    if (sparkRTC.value.localStream.id === stream.id) {
      local = true
    }
  }

  setUserActionLoading(stream.userId, false)

  streamers.value = {
    ...streamers.value,
    [stream.id]: {
      name: stream.name,
      userId: stream.userId,
      isHost: stream.role === Roles.BROADCAST,
      avatar: '',
      raisedHand: false,
      hasCamera: false,
      muted: streamers.value[stream.id] ? streamers.value[stream.id].muted : undefined,
      stream,
      isLocalStream: local,
      isShareScreen: stream.isShareScreen || false,
      toggleScreenId: toggleFull ? stream.id : null,
    },
  }
}

const toggleFullScreen = async (stream) => {
  await displayStream(stream, true)
}

export const onStopStream = async (stream) => {
  await toggleFullScreen(stream)

  const streamersTmp = { ...streamers.value }
  delete streamersTmp[stream.id]
  streamers.value = streamersTmp
}

export const onStopShareScreen = async (stream) => {
  await onStopStream(stream)

  stream.getTracks().forEach((track) => track.stop())
  updateUser({
    sharingScreenStream: null,
  })
  const streamersTmp = { ...streamers.value }
  delete streamersTmp[stream.id]
  streamers.value = streamersTmp
}

const log = (tag, data?: any) => {
  const date = new Date().toLocaleTimeString()

  if (data) {
    console.log('[', date, '] ', tag, ' | ', data)
  } else {
    console.log('[', date, '] ', tag)
  }
}
const setupSignalingSocket = async (host, name, room, debug) => {
  await sparkRTC.value.setupSignalingSocket(getWsUrl(host), JSON.stringify({ name, email: '' }), room, debug)
}
const start = async () => {
  return sparkRTC.value.start()
}

export const leaveMeeting = () => {
  if (sparkRTC.value) {
    sparkRTC.value.leaveMeeting()
    meetingStatus.value = false
    streamers.value = {}
  }
}

export const onUserRaisedHand = (userId, raisedHand, actionLoading: boolean, acceptRaiseHand?: any) => {
  Loading = actionLoading
  attendees.value = {
    ...attendees.value,
    [userId]: {
      ...attendees.value[userId],
      raisedHand,
      actionLoading,
      acceptRaiseHand,
    },
  }
  console.log('LOWER HAND', userId, raisedHand)
  sparkRTC.value.getLatestUserList('onUserRaiseHand')
}

export const onInviteToStage = (participant) => {
  if (sparkRTC.value.raiseHands.length >= sparkRTC.value.maxRaisedHands) {
    makeDialog('info', {
      message: 'The stage is already full. try again later.',
      icon: 'Close',
      variant: 'danger',
    })
  } else {
    //send invite
    console.log('Send Intitation to ', participant)
    sparkRTC.value.inviteToStage(participant.userId)
    setUserActionLoading(participant.userId, true)

    makeDialog('info', {
      message: `You've sent the request.`,
      icon: 'Check',
    })
  }
}

export const getUserRaiseHandStatus = (userId) => {
  return attendees.value[userId]?.raisedHand || false
}

const Meeting = () => {
  useEffect(() => {
    const queryParams = new URLSearchParams(window.location.search)
    const name = queryParams.get('name')
    var role = queryParams.get('role')
    const room = queryParams.get('room')
    const host = queryParams.get('host')

    var previewDialogId = null

    if (role === null || role === '') {
      role = Roles.AUDIENCE //by default set role to Audience
    }

    updateUser({
      name,
      role,
      isStreamming: role === Roles.BROADCAST,
      isHost: role === Roles.BROADCAST,
    })

    const setupSparkRTC = async () => {
      log(`Setup SparkRTC`)

      sparkRTC.value = createSparkRTC(role, {
        onAudioStatusChange: (message) => {
          log('audioStatus: ', message)
          if (message.stream != undefined && message.type != undefined && streamers.value != undefined) {
            if (streamers.value[message.stream] != undefined) {
              streamers.value[message.stream][message.type] = message.value
              streamers.value = { ...streamers.value }
            }
          }
        },
        onUserInitialized: (userId) => {
          //@ts-ignore
          currentUser.value.userId = userId
        },
        localStreamChangeCallback: (stream) => {
          log('[Local Stream Callback]', stream)
          streamers.value = {
            ...streamers.value,
            [stream.id]: {
              name: stream.name,
              isHost: role === Roles.BROADCAST,
              avatar: '',
              raisedHand: false,
              hasCamera: false,
              stream,
              isLocalStream: true,
              isShareScreen: stream.isShareScreen || false,
            },
          }
        },
        remoteStreamCallback: async (stream) => {
          log(`remoteStreamCallback`, stream)
          log(`remoteStreamCallback-Name`, stream.name)

          await displayStream(stream)

          if (!sparkRTC.value.broadcasterDC && role === Roles.AUDIENCE) {
            broadcastIsInTheMeeting.value = true
          }
        },
        remoteStreamDCCallback: async (stream) => {
          sparkRTC.value.getLatestUserList(`remote stream DC`)

          log(`remoteStreamDCCallback`, stream)

          if (stream != 'no-stream') {
            onStopStream(stream)
          } else {
            //get all remote streams and stop them
            const streams = sparkRTC.value.remoteStreams
            streams.forEach((str) => {
              onStopStream(str)
            })

            sparkRTC.value.remoteStreams = []
          }

          //display broadcaster not in the meeting message after 1 sec, to avoid any issues
          setTimeout(() => {
            if (role === Roles.AUDIENCE) {
              if (sparkRTC.value.broadcasterDC || stream === 'no-stream') {
                //destroy preview Dialog
                if (previewDialogId !== null) {
                  destroyDialog(previewDialogId)
                }

                broadcastIsInTheMeeting.value = false
                updateUser({
                  isStreamming: false,
                  ableToRaiseHand: true,
                  isMicrophoneOn: true,
                  isCameraOn: true,
                })
                sparkRTC.value.resetAudioVideoState()
                log(`broadcasterDC...`)
              }
            }
          }, 1000)
        },
        onRaiseHand: (user) => {
          log(`[On Raise Hand Request]`, user)

          let raiseHandCallback: (values: any) => void = () => {}
          const handler = new Promise((resolve, reject) => {
            raiseHandCallback = (value) => {
              setUserActionLoading(user.userId, true)
              resolve(value)
            }
          })

          //only show message when limit is not reached
          if (sparkRTC.value.raiseHands.length < sparkRTC.value.maxRaisedHands) {
            attendeesBadge.value = true

            makeDialog(
              'info',
              {
                message: 'Someone has raised their hand!',
                icon: 'Clock',
              },
              null,
              () => {
                isAttendeesOpen.value = true
                attendeesBadge.value = false
              }
            )
          }

          onUserRaisedHand(user.userId, new Date(), false, raiseHandCallback)

          return handler
        },
        onStart: async (closeSocket = false) => {
          if (meetingStatus.value) {
            if (role === Roles.AUDIENCE) {
              await sparkRTC.value.restart(closeSocket)
            }

            if (!closeSocket) {
              //start sparkRTC
              await start()
            }
          }
        },
        startAgain: async () => {
          if (sparkRTC.value) {
            //Init socket and start sparkRTC
            await setupSignalingSocket(host, name, room, isDebugMode.value)
            await start()
          }
        },
        altBroadcastApprove: async (isStreamming, data) => {
          //@ts-ignore
          setUserActionLoading(currentUser.value.userId, false)

          console.log('altBroadcastApprove: data: ', data)

          if (!isStreamming) {
            sparkRTC.value.onRaiseHandRejected()
            makeDialog('info', {
              message: 'You’ve been Rejected',
              icon: 'Close',
              variant: 'danger',
            })
            updateUser({
              ableToRaiseHand: true,
            })
          } else {
            const localStream = await sparkRTC.value.getAccessToLocalStream()

            previewDialogId = makePreviewDialog(
              'preview',
              localStream,
              {
                message: 'Set the default state of your “Video” and “Audio” before joining the stage please',
                title: 'Join The Stage',
              },
              () => {
                //onOk
                updateUser({
                  isStreamming,
                  ableToRaiseHand: true,
                })
                sparkRTC.value.joinStage(data)
                makeDialog('info', {
                  message: 'You’ve been added to the stage',
                  icon: 'Check',
                })

                //send user mute status to everyone to update the Ui
                setTimeout(() => {
                  if (sparkRTC.value.lastAudioState === sparkRTC.value.LastState.DISABLED) {
                    sparkRTC.value.sendAudioStatus(false)
                  } else {
                    sparkRTC.value.sendAudioStatus(true)
                  }
                }, 2000)
              },
              () => {
                //onClose
                updateUser({
                  ableToRaiseHand: true,
                  isMicrophoneOn: true,
                  isCameraOn: true,
                })

                sparkRTC.value.resetAudioVideoState()
                sparkRTC.value.cancelJoinStage(data)
                sparkRTC.value.onRaiseHandRejected()
              }
            )
          }
        },
        disableBroadcasting: () => {
          updateUser({
            isStreamming: false,
            ableToRaiseHand: true,
            isMicrophoneOn: true,
            isCameraOn: true,
          })
          makeDialog('info', {
            message: 'You just removed from stage',
            icon: 'Close',
            variant: 'danger',
          })
          sparkRTC.value.resetAudioVideoState()
          sparkRTC.value.onRaiseHandRejected()
        },
        maxLimitReached: (message) => {
          makeDialog('info', { message, icon: 'Close' })
          updateUser({
            isStreamming: false,
            ableToRaiseHand: true,
            isMicrophoneOn: true,
            isCameraOn: true,
          })
          sparkRTC.value.resetAudioVideoState()
        },
        onUserListUpdate: (users) => {
          // log(`[On Users List Update]`, users);
          const usersTmp = {}
          for (const { name: userInfo, role, video, id: userId } of users) {
            const { name, email } = JSON.parse(userInfo)
            usersTmp[userId] = {
              ...(attendees.value[userId] || {}),
              name,
              email,
              isHost: role === Roles.BROADCASTER,
              avatar: '',
              raisedHand: getUserRaiseHandStatus(userId),
              hasCamera: !!video,
              userId,
              video,
            }
          }
          //get latest raise hand count from sparkRTC
          if (sparkRTC.value.role === sparkRTC.value.Roles.BROADCAST) {
            raisedHandsCount.value = sparkRTC.value.raiseHands.length
          } else {
            //@ts-ignore
            raisedHandsCount.value = Object.values(usersTmp).reduce((prev, user) => {
              //@ts-ignore
              if (!user.isHost && user.video) return prev + 1
              return prev
            }, 0)
          }

          attendees.value = usersTmp
        },
        constraintResults: (constraints) => {
          if (!constraints.audio) {
            //remove mic button
            updateUser({ hasMic: false })
          }

          if (!constraints.video) {
            //remove video button
            updateUser({ hasCamera: false })
          }
        },
        updateStatus: (tag, data) => {
          log(tag, data)
        },
        treeCallback: (tree) => {
          log(`tree`, tree)
        },
        connectionStatus: (status) => {
          log(`Connection Status: `, status)
        },
        updateUi: () => {
          updateUser({
            showControllers: true,
            isStreamming: false,
            ableToRaiseHand: true,
          })
        },
        parentDcMessage: () => {
          makeDialog('info', {
            message: `For some unexpected reason, you've gotten disconnected. please wait some seconds to reconnect again.`,
            icon: 'Close',
            variant: 'danger',
          })
        },
        userLoweredHand: (data, name) => {
          //@ts-ignore
          onUserRaisedHand(data, false, false)
          log('userLoweredHand: ', data)
          sparkRTC.value.getLatestUserList('UserLowerHand')

          //get raise hand count from attendees list
          const rC = Object.values(attendees.value).reduce((prev, user) => {
            //@ts-ignore
            if (!user.isHost && user.raisedHand) return prev + 1
            return prev
          }, 0)
          console.log('raiseHandCount:', rC)

          if (rC === 0) {
            attendeesBadge.value = false
          }

          if (name) {
            makeDialog('info', {
              message: `${name} has rejected your request to join stage.`,
              icon: 'Close',
              variant: 'danger',
            })
          }
        },

        invitationToJoinStage: async (msg) => {
          console.log('invitationToJoinStage: ', msg)
          //show preview dialog to Join stage
          const localStream = await sparkRTC.value.getAccessToLocalStream()

          previewDialogId = makePreviewDialog(
            'preview',
            localStream,
            {
              message: 'The host has requested you to come on stage. Set the default state of your “Video” and “Audio” before joining please.',
              title: 'Join The Stage',
              yesButton: 'Join The Stage!',
              noButton: 'Not Now!',
            },
            () => {
              //onOk
              updateUser({
                isStreamming: true,
                ableToRaiseHand: true,
              })
              sparkRTC.value.joinStage(msg.data)
              makeDialog('info', {
                message: 'You’ve been added to the stage',
                icon: 'Check',
              })

              //send user mute status to everyone to update the Ui
              setTimeout(() => {
                if (sparkRTC.value.lastAudioState === sparkRTC.value.LastState.DISABLED) {
                  sparkRTC.value.sendAudioStatus(false)
                } else {
                  sparkRTC.value.sendAudioStatus(true)
                }
              }, 2000)
            },
            () => {
              console.log('audience-broadcasting cancelling..')
              //onClose
              updateUser({
                ableToRaiseHand: true,
                isMicrophoneOn: true,
                isCameraOn: true,
              })

              sparkRTC.value.resetAudioVideoState()
              sparkRTC.value.cancelJoinStage(msg.data, true)
              sparkRTC.value.onRaiseHandRejected()
            }
          )
        },
      })

      if (sparkRTC.value) {
        //Init socket and start sparkRTC
        await setupSignalingSocket(host, name, room, isDebugMode.value)
        await start()
      }
    }

    if (meetingStatus.value) {
      setupSparkRTC()
    }
  }, [meetingStatus.value])

  const rejoinMeeting = () => {
    window.location.reload()
  }

  const leaveMeeting = () => {
    window.parent.postMessage('leave', '*')
  }

  return (
    <div class="flex flex-col justify-between min-h-[--doc-height] dark:bg-secondary-1-a bg-white-f-9 text-medium-12 text-gray-800 dark:text-gray-200">
      <TopBar />
      {meetingStatus.value ? (
        <>
          <MeetingBody />
          <BottomBar />
        </>
      ) : (
        <div class="flex flex-col justify-center items-center sm:p-10 rounded-md gap-4 h-full flex-grow">
          <span class="text-bold-18 text-center">You Left The Live Show</span>
          <div class="flex w-full justify-center items-center gap-4 flex-row max-w-[85%] sm:max-w-[400px]">
            <Button onClick={leaveMeeting} variant="outline" class="flex-1 w-full px-0">
              Go To Home Feed
            </Button>
            <Button onClick={rejoinMeeting} variant="primary" class="flex-1 w-full px-0">
              Rejoin
            </Button>
          </div>
        </div>
      )}

      <ToastProvider />
    </div>
  )
}

export default Meeting
