import { get, getList, post, put } from './client'
import type { Child, CourseItem, IprPlan, Mother, MotherStatusInfo, MotherStatusItem, Profile, ProfileSummary, SecretInfo } from '../types'

export const getMothers = (q = '') => getList<Mother>('/api/mother?q=' + encodeURIComponent(q))
export const createMother = (first_name: string, last_name: string, phone: string) =>
  post<Mother>('/api/mother', { first_name, last_name, phone })
export const uploadPhoto = (motherId: number, file: File) => {
  const formData = new FormData()
  formData.append('photo', file)
  return fetch('/api/mother/' + motherId + '/photo', { method: 'POST', body: formData }).then((res) => {
    if (!res.ok) return res.json().then((e) => { throw new Error(e.error || 'upload failed') })
    return res.json()
  })
}
export const getChildren = (motherId: number) => getList<Child>('/api/mother/' + motherId + '/children')
export const addChild = (motherId: number, first_name: string, birth_date: string) =>
  post<Child>('/api/mother/' + motherId + '/children', { first_name, birth_date })
export const getSecret = (motherId: number) => get<SecretInfo>('/api/mother/' + motherId + '/sensitive')
export const saveSecret = (motherId: number, data: SecretInfo) =>
  put<{ status: string }>('/api/mother/' + motherId + '/sensitive', data)
export const getProfile = (motherId: number) => get<Profile>('/api/mother/' + motherId + '/profile')
export const saveProfile = (motherId: number, data: Record<string, string | number>) =>
  put<{ status: string }>('/api/mother/' + motherId + '/profile', data)
export const setProfileStatus = (motherId: number, status: string) =>
  put<{ status: string }>('/api/mother/' + motherId + '/profile/status', { status })
export const getProfiles = (status: string) => getList<ProfileSummary>('/api/profiles?status=' + encodeURIComponent(status))
export const getMotherStatus = (motherId: number) => get<MotherStatusInfo>('/api/mother/' + motherId + '/status')
export const setMotherStatus = (motherId: number, status: string, reason: string) =>
  put<{ status: string }>('/api/mother/' + motherId + '/status', { status, reason })
export const getMotherStatuses = () => getList<MotherStatusItem>('/api/mother-statuses')
export const getCourses = (motherId: number) => getList<CourseItem>('/api/mother/' + motherId + '/courses')
export const getIPR = (motherId: number) => getList<IprPlan>('/api/mother/' + motherId + '/ipr')
