export interface StudentGradeRecord {
  studentName: string
  studentEmail: string
  assignmentTitle: string
  score: number
  submittedAt: string
}

export function exportGradesToCSV(records: StudentGradeRecord[], filename = 'rekap_nilai.csv') {
  if (!import.meta.client || !records.length) return

  const headers = ['Nama Siswa', 'Email', 'Judul Tugas', 'Nilai', 'Waktu Pengumpulan']
  const rows = records.map((r) => [
    `"${r.studentName.replace(/"/g, '""')}"`,
    `"${r.studentEmail}"`,
    `"${r.assignmentTitle.replace(/"/g, '""')}"`,
    r.score,
    `"${r.submittedAt}"`
  ])

  const csvContent = [headers.join(','), ...rows.map((row) => row.join(','))].join('\n')
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)

  const link = document.createElement('a')
  link.setAttribute('href', url)
  link.setAttribute('download', filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
