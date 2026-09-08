/**
 * LMS Shared Utility Helper Functions
 * General helper functions for date formatting, text truncation, and currency display.
 */

export const formatDate = (dateString, locale = 'id-ID') => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return new Intl.DateTimeFormat(locale, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

export const truncateText = (text, maxLength = 100) => {
  if (!text || text.length <= maxLength) return text
  return text.slice(0, maxLength).trim() + '...'
}

export const formatCurrency = (amount, currency = 'IDR') => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0
  }).format(amount || 0)
}

export const slugify = (text) => {
  return text
    .toString()
    .toLowerCase()
    .trim()
    .replace(/\s+/g, '-')
    .replace(/[^\w-]+/g, '')
    .replace(/--+/g, '-')
}
