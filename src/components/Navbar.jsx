import React from 'react'

export const Navbar = ({ user, onLogout }) => {
  return (
    <header className="bg-white dark:bg-surface-900 border-b border-surface-200 dark:border-surface-800 sticky top-0 z-30">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div className="flex items-center space-x-3">
          <span className="text-2xl font-bold bg-gradient-to-r from-brand-600 to-indigo-600 bg-clip-text text-transparent">
            🎓 LMS Portal
          </span>
        </div>
        <nav className="hidden md:flex items-center space-x-6 text-sm font-medium text-surface-600 dark:text-surface-300">
          <a href="/dashboard" className="hover:text-brand-600 dark:hover:text-brand-400 transition">Dashboard</a>
          <a href="/classes" className="hover:text-brand-600 dark:hover:text-brand-400 transition">Kelas Saya</a>
          <a href="/assignments" className="hover:text-brand-600 dark:hover:text-brand-400 transition">Tugas</a>
          <a href="/quizzes" className="hover:text-brand-600 dark:hover:text-brand-400 transition">Kuis</a>
        </nav>
        <div className="flex items-center space-x-4">
          {user ? (
            <div className="flex items-center space-x-3">
              <span className="text-sm font-semibold text-surface-800 dark:text-surface-100">{user.name}</span>
              <button onClick={onLogout} className="text-xs bg-surface-100 dark:bg-surface-800 hover:bg-surface-200 px-3 py-1.5 rounded-md font-medium text-surface-700 dark:text-surface-300 transition">Keluar</button>
            </div>
          ) : (
            <a href="/login" className="text-sm bg-brand-600 hover:bg-brand-700 text-white px-4 py-2 rounded-lg font-medium shadow-sm transition">Masuk</a>
          )}
        </div>
      </div>
    </header>
  )
}

export default Navbar
