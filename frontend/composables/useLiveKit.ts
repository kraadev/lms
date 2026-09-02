import {
  Room,
  RoomEvent,
  VideoPresets,
  createLocalVideoTrack,
  createLocalAudioTrack,
  LocalVideoTrack,
  LocalAudioTrack,
  Participant,
  RemoteParticipant,
  Track,
  ConnectionState
} from 'livekit-client'
import { meetingsService } from '~/services/meetings'
import type { Meeting } from '~/types'

export interface ParticipantTrackInfo {
  participant: Participant
  isLocal: boolean
  isSpeaking: boolean
  isAudioMuted: boolean
  isVideoMuted: boolean
  videoTrack?: Track | null
  audioTrack?: Track | null
}

export function useLiveKit() {
  const room = useState<Room | null>('livekit_room', () => null)
  const connectionState = useState<ConnectionState>('livekit_connection_state', () => ConnectionState.Disconnected)
  const participants = useState<ParticipantTrackInfo[]>('livekit_participants', () => [])
  const activeSpeakers = useState<string[]>('livekit_active_speakers', () => [])
  const isScreenSharing = useState<boolean>('livekit_screen_sharing', () => false)
  const localAudioTrack = useState<LocalAudioTrack | null>('livekit_local_audio', () => null)
  const localVideoTrack = useState<LocalVideoTrack | null>('livekit_local_video', () => null)
  const error = useState<string | null>('livekit_error', () => null)

  function updateParticipantsList(currentRoom: Room) {
    const list: ParticipantTrackInfo[] = []

    // Local participant
    if (currentRoom.localParticipant) {
      const local = currentRoom.localParticipant
      list.push({
        participant: local,
        isLocal: true,
        isSpeaking: local.isSpeaking,
        isAudioMuted: local.isMicrophoneEnabled === false,
        isVideoMuted: local.isCameraEnabled === false,
        videoTrack: local.getTrackPublication(Track.Source.Camera)?.track || null,
        audioTrack: local.getTrackPublication(Track.Source.Microphone)?.track || null
      })
    }

    // Remote participants
    currentRoom.remoteParticipants.forEach((p) => {
      list.push({
        participant: p,
        isLocal: false,
        isSpeaking: p.isSpeaking,
        isAudioMuted: !p.isMicrophoneEnabled,
        isVideoMuted: !p.isCameraEnabled,
        videoTrack: p.getTrackPublication(Track.Source.Camera)?.track || null,
        audioTrack: p.getTrackPublication(Track.Source.Microphone)?.track || null
      })
    })

    participants.value = list
  }

  async function joinMeeting(meetingId: number | string, options: { enableMic: boolean; enableCam: boolean; audioDeviceId?: string; videoDeviceId?: string }) {
    error.value = null

    try {
      // 1. Fetch LiveKit connection info from Go backend
      const response = await meetingsService.join(meetingId)
      const { url, token } = response

      if (!url || !token) {
        throw new Error('URL atau token LiveKit tidak valid dari server.')
      }

      // 2. Initialize LiveKit Room
      const newRoom = new Room({
        adaptiveStream: true,
        dynacast: true,
        videoCaptureDefaults: {
          resolution: VideoPresets.h720.resolution
        }
      })

      // Attach Room Events
      newRoom
        .on(RoomEvent.ConnectionStateChanged, (state) => {
          connectionState.value = state
        })
        .on(RoomEvent.Connected, () => {
          connectionState.value = ConnectionState.Connected
          updateParticipantsList(newRoom)
        })
        .on(RoomEvent.Disconnected, () => {
          connectionState.value = ConnectionState.Disconnected
          participants.value = []
        })
        .on(RoomEvent.ParticipantConnected, () => updateParticipantsList(newRoom))
        .on(RoomEvent.ParticipantDisconnected, () => updateParticipantsList(newRoom))
        .on(RoomEvent.TrackSubscribed, () => updateParticipantsList(newRoom))
        .on(RoomEvent.TrackUnsubscribed, () => updateParticipantsList(newRoom))
        .on(RoomEvent.TrackMuted, () => updateParticipantsList(newRoom))
        .on(RoomEvent.TrackUnmuted, () => updateParticipantsList(newRoom))
        .on(RoomEvent.LocalTrackPublished, () => updateParticipantsList(newRoom))
        .on(RoomEvent.LocalTrackUnpublished, () => updateParticipantsList(newRoom))
        .on(RoomEvent.ActiveSpeakersChanged, (speakers) => {
          activeSpeakers.value = speakers.map(s => s.identity)
          updateParticipantsList(newRoom)
        })

      // 3. Connect to LiveKit server
      await newRoom.connect(url, token)
      room.value = newRoom

      // 4. Publish local tracks based on pre-join choices
      if (options.enableMic) {
        await newRoom.localParticipant.setMicrophoneEnabled(true, {
          deviceId: options.audioDeviceId
        })
      } else {
        await newRoom.localParticipant.setMicrophoneEnabled(false)
      }

      if (options.enableCam && response.meeting?.type !== 'audio') {
        await newRoom.localParticipant.setCameraEnabled(true, {
          deviceId: options.videoDeviceId
        })
      } else {
        await newRoom.localParticipant.setCameraEnabled(false)
      }

      updateParticipantsList(newRoom)
      return { meeting: response.meeting, room: newRoom }
    } catch (err: any) {
      if (err?.status === 403) {
        error.value = 'Anda tidak memiliki akses untuk bergabung ke kelas meeting ini.'
      } else {
        error.value = err?.message || 'Gagal terhubung ke ruang meeting.'
      }
      throw err
    }
  }

  async function toggleMicrophone(): Promise<boolean> {
    if (!room.value?.localParticipant) return false
    const currentState = room.value.localParticipant.isMicrophoneEnabled
    await room.value.localParticipant.setMicrophoneEnabled(!currentState)
    updateParticipantsList(room.value)
    return !currentState
  }

  async function toggleCamera(): Promise<boolean> {
    if (!room.value?.localParticipant) return false
    const currentState = room.value.localParticipant.isCameraEnabled
    await room.value.localParticipant.setCameraEnabled(!currentState)
    updateParticipantsList(room.value)
    return !currentState
  }

  async function toggleScreenShare(): Promise<boolean> {
    if (!room.value?.localParticipant) return false
    const currentState = isScreenSharing.value
    await room.value.localParticipant.setScreenShareEnabled(!currentState)
    isScreenSharing.value = !currentState
    updateParticipantsList(room.value)
    return !currentState
  }

  async function leaveMeeting() {
    if (room.value) {
      room.value.disconnect()
      room.value = null
    }
    participants.value = []
    connectionState.value = ConnectionState.Disconnected
    isScreenSharing.value = false
  }

  return {
    room,
    connectionState,
    participants,
    activeSpeakers,
    isScreenSharing,
    error,
    joinMeeting,
    toggleMicrophone,
    toggleCamera,
    toggleScreenShare,
    leaveMeeting
  }
}
