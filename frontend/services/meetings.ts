import { api } from './api'
import type { Meeting, MeetingType, LiveKitJoinResponse } from '~/types'

export const meetingsService = {
  getByClass(classId: number | string): Promise<Meeting[]> {
    return api.get<Meeting[]>(`/classes/${classId}/meetings`)
  },

  getById(id: number | string): Promise<Meeting> {
    return api.get<Meeting>(`/meetings/${id}`)
  },

  create(classId: number | string, payload: { title: string; type: MeetingType }): Promise<Meeting> {
    return api.post<Meeting>(`/classes/${classId}/meetings`, payload)
  },

  // Authenticate and fetch LiveKit url + token
  join(meetingId: number | string): Promise<LiveKitJoinResponse> {
    return api.post<LiveKitJoinResponse>(`/meetings/${meetingId}/join`)
  },

  // Teacher end meeting for everyone
  end(meetingId: number | string): Promise<{ message: string }> {
    return api.post<{ message: string }>(`/meetings/${meetingId}/end`)
  }
}
