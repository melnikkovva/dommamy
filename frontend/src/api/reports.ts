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
