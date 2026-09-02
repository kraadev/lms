import { api } from './api'
import type { Quiz, QuizQuestion, QuizAttempt, QuizAttemptAnswer } from '~/types'

export interface CreateQuizPayload {
  title: string
  description?: string
  duration_minutes: number
  attempts_allowed: number
  due_date?: string
  questions: {
    question_text: string
    question_type: 'multiple_choice' | 'essay'
    points: number
    order: number
    options?: { option_text: string; order: number; is_correct?: boolean }[]
  }[]
}

export const quizzesService = {
  getAll(params?: { class_id?: number | string }): Promise<Quiz[]> {
    return api.get<Quiz[]>('/quizzes', params)
  },

  getByClass(classId: number | string): Promise<Quiz[]> {
    return api.get<Quiz[]>(`/classes/${classId}/quizzes`)
  },

  getById(id: number | string): Promise<Quiz & { questions?: QuizQuestion[] }> {
    return api.get<Quiz & { questions?: QuizQuestion[] }>(`/quizzes/${id}`)
  },

  create(classId: number | string, payload: CreateQuizPayload): Promise<Quiz> {
    return api.post<Quiz>(`/classes/${classId}/quizzes`, payload)
  },

  delete(id: number | string): Promise<{ message: string }> {
    return api.delete<{ message: string }>(`/quizzes/${id}`)
  },

  // Student: Start a new attempt
  startAttempt(quizId: number | string): Promise<QuizAttempt & { questions: QuizQuestion[] }> {
    return api.post<QuizAttempt & { questions: QuizQuestion[] }>(`/quizzes/${quizId}/attempts`)
  },

  // Student: Get current attempt
  getAttempt(attemptId: number | string): Promise<QuizAttempt & { questions: QuizQuestion[] }> {
    return api.get<QuizAttempt & { questions: QuizQuestion[] }>(`/attempts/${attemptId}`)
  },

  // Student: Submit attempt
  submitAttempt(attemptId: number | string, answers: QuizAttemptAnswer[]): Promise<QuizAttempt> {
    return api.post<QuizAttempt>(`/attempts/${attemptId}/submit`, { answers })
  }
}
