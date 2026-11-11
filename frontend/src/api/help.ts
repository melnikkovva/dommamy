import { getList, post, put } from './client'
import type { HelpItem, HelpRequest } from '../types'

export const getHelp = () => getList<HelpItem>('/api/help')
export const giveHelp = (mother_id: number, kind: string, sum: number, check_url: string) =>
  post<HelpItem>('/api/help', { mother_id, kind, sum, check_url })
export const getRequests = (status = '') => getList<HelpRequest>('/api/requests?status=' + encodeURIComponent(status))
export const getMotherRequests = (motherId: number) => getList<HelpRequest>('/api/mother/' + motherId + '/requests')
export const createRequest = (mother_id: number, kind: string, description: string) =>
  post<HelpRequest>('/api/requests', { mother_id, kind, description })
export const setRequestStatus = (id: number, status: string, comment: string) =>
  put<HelpRequest>('/api/requests/' + id + '/status', { status, comment })
export const issueRequest = (id: number, help_id: number) =>
  post<HelpRequest>('/api/requests/' + id + '/issue', { help_id })
