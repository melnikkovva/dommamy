import { get, getList, post, put } from './client'
import type { Donation, Donor, DonorBody, DonorDetails, DonorSummary } from '../types'

export const getDonors = (q = '') => getList<DonorSummary>('/api/donors?q=' + encodeURIComponent(q))
export const getDonor = (id: number) => get<DonorDetails>('/api/donors/' + id)
export const createDonor = (body: DonorBody) => post<Donor>('/api/donors', body)
export const updateDonor = (id: number, body: DonorBody) => put<Donor>('/api/donors/' + id, body)
export const addDonation = (id: number, amount: number, date: string, comment: string) =>
  post<Donation>('/api/donors/' + id + '/donations', { amount, date, comment })
