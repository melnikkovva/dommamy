import { getList, post, put, upload } from './client'
import type { DocumentItem } from '../types'

export const getDocuments = (motherId: number) => getList<DocumentItem>('/api/mother/' + motherId + '/documents')
export const createDocument = (motherId: number, type: string, help_id: number) =>
  post<DocumentItem>('/api/mother/' + motherId + '/documents', { type, help_id })
export const uploadDocument = (id: number, file: File) => upload<DocumentItem>('/api/documents/' + id + '/file', file)
export const signDocument = (id: number) => put<DocumentItem>('/api/documents/' + id + '/sign', {})
export const documentFileUrl = (id: number) => '/api/documents/' + id + '/file'
