import { api } from './api'
import type { StudentDashboardData, TeacherDashboardData, AdminDashboardData } from '~/types'

export const dashboardService = {
  getStudentDashboard(): Promise<StudentDashboardData> {
    return api.get<StudentDashboardData>('/dashboard/student')
  },

  getTeacherDashboard(): Promise<TeacherDashboardData> {
    return api.get<TeacherDashboardData>('/dashboard/teacher')
  },

  getAdminDashboard(): Promise<AdminDashboardData> {
    return api.get<AdminDashboardData>('/dashboard/admin')
  }
}
