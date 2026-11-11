import { get } from './client'
import type { Report } from '../types'

function query(from: string, to: string) {
  return '?from=' + encodeURIComponent(from) + '&to=' + encodeURIComponent(to)
}
export const getReport = (name: string, from: string, to: string) =>
  get<Report>('/api/reports/' + name + query(from, to))
export const reportCsvUrl = (name: string, from: string, to: string) =>
  '/api/reports/' + name + '/csv' + query(from, to)
