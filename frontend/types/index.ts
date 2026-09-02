export type UserRole = 'admin' | 'teacher' | 'student'

export interface User {
  id: number
  name: string
  email: string
  role: UserRole
  avatar?: string
  phone?: string
  status?: 'active' | 'inactive'
  created_at?: string
}

export interface AuthResponse {
  user: User
  message?: string
}

export interface Class {
  id: number
  title: string
  name?: string
  code: string
  description?: string
  academic_year: string
  semester?: string
  teacher_id: number
  teacher?: User
  member_count?: number
  status: 'active' | 'archived'
  created_at?: string
  active_meeting?: Meeting | null
}

export interface ClassMember {
  id: number
  class_id: number
  user_id: number
  user: User
  role: UserRole
  joined_at: string
}

export interface MaterialAttachment {
  name: string
  url: string
  size?: number
}

export interface Material {
  id: number
  class_id: number
  title: string
  description?: string
  file_url?: string
  file_name?: string
  file_size?: number
  external_link?: string
  created_at: string
  updated_at?: string
  created_by?: User
}

export type SubmissionStatus = 'not_submitted' | 'submitted' | 'late' | 'graded'

export interface Assignment {
  id: number
  class_id: number
  class_title?: string
  class?: Class
  title: string
  instructions: string
  due_date: string
  points: number
  attachments?: MaterialAttachment[]
  attachment_url?: string
  attachment_name?: string
  attachment_path?: string
  created_at: string
  submission?: AssignmentSubmission
  my_submission?: AssignmentSubmission
  submissions_count?: number
  graded_count?: number
  total_students?: number
}

export interface AssignmentSubmission {
  id: number
  assignment_id: number
  student_id: number
  student?: User
  text_response?: string
  file_url?: string
  file_name?: string
  file_size?: number
  status: SubmissionStatus
  score?: number | null
  feedback?: string | null
  submitted_at: string
  graded_at?: string | null
}

export type QuizStatus = 'upcoming' | 'active' | 'closed'

export interface QuizOption {
  id: number
  question_id: number
  option_text: string
  order: number
}

export interface QuizQuestion {
  id: number
  quiz_id: number
  question_text: string
  question_type: 'multiple_choice' | 'essay'
  points: number
  order: number
  options?: QuizOption[]
}

export interface Quiz {
  id: number
  class_id: number
  class_title?: string
  title: string
  description?: string
  duration_minutes: number
  attempts_allowed: number
  due_date?: string
  start_time?: string
  end_time?: string
  total_questions: number
  total_points: number
  status: QuizStatus
  my_attempts_count?: number
  my_best_score?: number | null
  my_latest_attempt?: QuizAttempt | null
}

export interface QuizAttemptAnswer {
  question_id: number
  selected_option_id?: number | null
  text_response?: string | null
}

export interface QuizAttempt {
  id: number
  quiz_id: number
  student_id: number
  started_at: string
  submitted_at?: string | null
  time_remaining_seconds?: number
  status: 'in_progress' | 'submitted' | 'timed_out'
  score?: number | null
  total_points?: number
  answers?: QuizAttemptAnswer[]
}

export type MeetingType = 'video' | 'audio'
export type MeetingStatus = 'scheduled' | 'active' | 'ended'

export interface Meeting {
  id: number
  class_id: number
  class_title?: string
  title: string
  type: MeetingType
  status: MeetingStatus
  host_id: number
  host?: User
  room_name?: string
  started_at?: string
  ended_at?: string
  participant_count?: number
  livekit_room_name?: string
  created_at?: string
}

export interface LiveKitJoinResponse {
  url: string
  token: string
  room_name?: string
  meeting: Meeting
}

export interface ChatMessage {
  id: string | number
  class_id: number
  sender_id: number
  sender: User
  message: string
  created_at: string
  status?: 'sending' | 'sent' | 'failed'
}

export interface WebSocketEvent<T = any> {
  type: string
  payload: T
}

export type NotificationType =
  | 'assignment_new'
  | 'assignment_graded'
  | 'quiz_new'
  | 'meeting_started'
  | 'general'

export interface Notification {
  id: number
  user_id: number
  title: string
  message: string
  type: NotificationType
  link?: string
  is_read: boolean
  created_at: string
}

export interface StudentDashboardData {
  greeting: string
  current_classes: Class[]
  upcoming_assignments: Assignment[]
  upcoming_quizzes: Quiz[]
  active_meeting?: Meeting | null
  recent_announcements: { id: number; title: string; class_title: string; content: string; created_at: string }[]
  recent_grades: { id: number; title: string; class_title: string; score: number; max_score: number; graded_at: string }[]
  unread_chat_count: number
}

export interface TeacherDashboardData {
  classes_taught: Class[]
  pending_grading_count: number
  pending_grading: { assignment: Assignment; submission_count: number }[]
  upcoming_assignments: Assignment[]
  recent_submissions: AssignmentSubmission[]
  active_meetings: Meeting[]
  recent_chats: { class_id: number; class_title: string; last_message: string; timestamp: string }[]
  quiz_overview: { total_quizzes: number; active_quizzes: number; total_attempts: number }
}

export interface AdminDashboardData {
  total_students: number
  total_teachers: number
  total_classes: number
  active_classes: number
  system_activity: { id: number; description: string; timestamp: string; type: string }[]
  recent_users: User[]
}

export interface ApiError {
  message: string
  status?: number
  errors?: Record<string, string[]>
}
