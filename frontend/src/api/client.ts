export class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}
async function send<T>(method: string, url: string, body?: BodyInit, headers?: Record<string, string>): Promise<T> {
  const res = await fetch(url, { method, body, headers })
  const text = await res.text()
  let data: any = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = null
    }
  }
  if (!res.ok) {
    throw new ApiError((data && data.error) || 'Ошибка сервера (код ' + res.status + ')', res.status)
  }
  return data as T
}
export function get<T>(url: string) {
  return send<T>('GET', url)
}
export async function getList<T>(url: string): Promise<T[]> {
  const data = await get<T[] | null>(url)
  return data || []
}
export function post<T>(url: string, body: unknown) {
  return send<T>('POST', url, JSON.stringify(body), { 'Content-Type': 'application/json' })
}
export function put<T>(url: string, body: unknown) {
  return send<T>('PUT', url, JSON.stringify(body), { 'Content-Type': 'application/json' })
}
export function upload<T>(url: string, file: File) {
  const form = new FormData()
  form.append('file', file)
  return send<T>('POST', url, form)
}
export function errorText(e: unknown) {
  return e instanceof Error ? e.message : 'Неизвестная ошибка'
}
