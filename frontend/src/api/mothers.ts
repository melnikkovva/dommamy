import { get, getList, post, put, remove, send } from './client'
import type { Child, CourseItem, IprPlan, Mother, MotherStatusInfo, MotherStatusItem, Profile, ProfileSummary, SecretInfo } from '../types'

export const getMothers = (q = '', phone = '', profileStatus = '', motherStatus = '') => {
  const params = new URLSearchParams()
  if (q) params.set('q', q)
  if (phone) params.set('phone', phone)
  if (profileStatus) params.set('profileStatus', profileStatus)
  if (motherStatus) params.set('motherStatus', motherStatus)
  return getList<Mother>('/api/mother?' + params.toString())
}
export const createMother = (first_name: string, last_name: string, phone: string, email: string) =>
  post<Mother>('/api/mother', { first_name, last_name, phone, email })
export const getLearners = () => getList<{ user_id: string; first_name: string; last_name: string; email: string }>('/api/learners')
export const linkLearnUser = (motherId: number, learnUserId: string) =>
  put<{ status: string }>('/api/mother/' + motherId + '/learn-link', { learn_user_id: learnUserId })
export const unlinkLearnUser = (motherId: number) =>
  remove('/api/mother/' + motherId + '/learn-link')
export const uploadPhoto = (motherId: number, file: File) => {
  const form = new FormData()
  form.append('photo', file)
  return send('POST', '/api/mother/' + motherId + '/photo', form)
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
export const deleteMother = (motherId: number) => remove(`/api/mother/${motherId}`)
