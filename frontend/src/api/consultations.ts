import { get, post, put } from './client'

export interface Consultation {
  id: number
  mother_id: number
  mother_name: string
  specialist_type: string
  consultation_date: string
  notes: string
  status: string
  created_at: string
}

export function getConsultations(motherId: number) {
  return get<Consultation[]>('/api/consultations/mother/' + motherId)
}

export function createConsultation(motherId: number, specialistType: string, date: string, notes: string, status: string) {
  return post<Consultation>('/api/consultations', {
    mother_id: motherId,
    specialist_type: specialistType,
    consultation_date: date,
    notes,
    status
  })
}

export function updateConsultationStatus(id: number, status: string) {
  return put<Consultation>('/api/consultations/' + id, { status })
}
