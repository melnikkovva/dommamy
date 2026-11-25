import { get } from './client'
import type { Report } from '../types'

function query(from: string, to: string, problemType: string) {
  let q = '?from=' + encodeURIComponent(from) + '&to=' + encodeURIComponent(to)
  if (problemType) q += '&problemType=' + encodeURIComponent(problemType)
  return q
}
export const getReport = (name: string, from: string, to: string, problemType: string) =>
  get<Report>('/api/reports/' + name + query(from, to, problemType))
export const reportCsvUrl = (name: string, from: string, to: string, problemType: string) =>
  '/api/reports/' + name + '/csv' + query(from, to, problemType)

export async function downloadReportCsv(name: string, from: string, to: string, problemType: string) {
  const token = localStorage.getItem('auth_token')
  const response = await fetch(reportCsvUrl(name, from, to, problemType), {
    headers: token ? { Authorization: 'Bearer ' + token } : undefined,
  })

  if (!response.ok) {
    const text = await response.text()
    let message = 'Ошибка сервера (код ' + response.status + ')'
    if (text) {
      try {
        const data = JSON.parse(text)
        message = data.error || message
      } catch {
        message = text
      }
    }
    throw new Error(message)
  }

  const filename = response.headers.get('Content-Disposition')?.match(/filename="?([^";]+)"?/)?.[1]
    || 'report-' + name + '.csv'
  const url = window.URL.createObjectURL(await response.blob())
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
  window.setTimeout(() => window.URL.revokeObjectURL(url), 1000)
}
