/**
 * Application Global Constants and Route Definitions
 */

export const APP_CONFIG = {
  NAME: 'LMS - Learning Management System',
  VERSION: '1.0.0',
  DEFAULT_LOCALE: 'id-ID',
  SUPPORT_EMAIL: 'support@lms.local'
}

export const USER_ROLES = {
  ADMIN: 'admin',
  TEACHER: 'teacher',
  STUDENT: 'student'
}

export const NAV_ROUTES = [
  { name: 'Dashboard', path: '/dashboard', icon: 'Home', roles: ['admin', 'teacher', 'student'] },
  { name: 'Classes', path: '/classes', icon: 'BookOpen', roles: ['admin', 'teacher', 'student'] },
  { name: 'Assignments', path: '/assignments', icon: 'FileText', roles: ['admin', 'teacher', 'student'] },
  { name: 'Quizzes', path: '/quizzes', icon: 'Award', roles: ['admin', 'teacher', 'student'] },
  { name: 'Meetings', path: '/meetings', icon: 'Video', roles: ['admin', 'teacher', 'student'] },
  { name: 'Messages', path: '/messages', icon: 'MessageSquare', roles: ['teacher', 'student'] },
  { name: 'Administration', path: '/admin', icon: 'Shield', roles: ['admin'] }
]

export const MEETING_STATUS = {
  SCHEDULED: 'scheduled',
  ACTIVE: 'active',
  ENDED: 'ended'
}
