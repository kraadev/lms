import { api } from './api'
import type { Assignment, AssignmentSubmission } from '~/types'

export const assignmentsService = {
  getAll(params?: { class_id?: number | string; status?: string }): Promise<Assignment[]> {
    return api.get<Assignment[]>('/assignments', params)
  },

  getByClass(classId: number | string): Promise<Assignment[]> {
    return api.get<Assignment[]>(`/classes/${classId}/assignments`)
  },

  getById(id: number | string): Promise<Assignment> {
    return api.get<Assignment>(`/assignments/${id}`)
  },

  create(classId: number | string, data: FormData | { title: string; instructions: string; due_date: string; points: number }): Promise<Assignment> {
    return api.post<Assignment>(`/classes/${classId}/assignments`, data)
  },

  update(id: number | string, data: FormData | Partial<Assignment>): Promise<Assignment> {
    return api.put<Assignment>(`/assignments/${id}`, data)
  },

  delete(id: number | string): Promise<{ message: string }> {
    return api.delete<{ message: string }>(`/assignments/${id}`)
  },

  // Student: Submit assignment
  submit(assignmentId: number | string, data: FormData | { text_response?: string }): Promise<AssignmentSubmission> {
    return api.post<AssignmentSubmission>(`/assignments/${assignmentId}/submit`, data)
  },

  // Teacher: Get all submissions for an assignment
  getSubmissions(assignmentId: number | string): Promise<AssignmentSubmission[]> {
    return api.get<AssignmentSubmission[]>(`/assignments/${assignmentId}/submissions`)
  },

  // Teacher: Grade a submission
  gradeSubmission(submissionId: number | string, payload: { score: number; feedback?: string }): Promise<AssignmentSubmission> {
    return api.post<AssignmentSubmission>(`/submissions/${submissionId}/grade`, payload)
  }
}
