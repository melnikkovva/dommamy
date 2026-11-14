import { post, get, ApiError } from './client'

export interface AuthUser {
  id: number
  login: string
  first_name: string
  last_name: string
  role: string
}

export interface LoginResponse {
  token: string
  id: number
  login: string
  first_name: string
  last_name: string
  role: string
}

export interface CreateUserRequest {
  login: string
  password: string
  first_name: string
  last_name: string
  role: string
}

export interface UserResponse {
  id: number
  login: string
  first_name: string
  last_name: string
  role: string
  created_at: string
}

export async function login(login: string, password: string): Promise<LoginResponse> {
  return post<LoginResponse>('/api/auth/login', { login, password })
}

export async function listUsers(): Promise<UserResponse[]> {
  return get<UserResponse[]>('/api/users')
}

export async function createUser(req: CreateUserRequest): Promise<UserResponse> {
  return post<UserResponse>('/api/users', req)
}

export function getToken(): string | null {
  return localStorage.getItem('auth_token')
}

export function getUser(): AuthUser | null {
  const raw = localStorage.getItem('auth_user')
  if (!raw) return null
  try {
    return JSON.parse(raw) as AuthUser
  } catch {
    return null
  }
}

export function setAuth(token: string, user: AuthUser) {
  localStorage.setItem('auth_token', token)
  localStorage.setItem('auth_user', JSON.stringify(user))
}

export function clearAuth() {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('auth_user')
}

export function isAuthenticated(): boolean {
  return !!getToken()
}

export function isAdmin(): boolean {
  const user = getUser()
  return user?.role === 'admin'
}
