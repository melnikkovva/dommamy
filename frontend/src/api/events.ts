import { getList, post, put } from './client'
import type { EventItem, EventType, Participant } from '../types'

export const getEventTypes = () => getList<EventType>('/api/event-types')
export const getEvents = () => getList<EventItem>('/api/events')
export const createEvent = (type_id: number, title: string, starts_at: string, place: string, max_people: number) =>
  post<EventItem>('/api/events', { type_id, title, starts_at, place, max_people })
export const getMotherEvents = (motherId: number) => getList<EventItem>('/api/mother/' + motherId + '/events')
export const getParticipants = (eventId: number) => getList<Participant>('/api/events/' + eventId + '/participants')
export const registerMother = (eventId: number, mother_id: number) =>
  post<{ status: string }>('/api/events/' + eventId + '/participants', { mother_id })
export const markParticipant = (eventId: number, motherId: number, status: string) =>
  put<{ status: string }>('/api/events/' + eventId + '/participants/' + motherId, { status })
