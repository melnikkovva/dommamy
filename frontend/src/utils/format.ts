import type { Mother } from '../types'

export function fmtDate(iso: string) {
  if (!iso) return '—'
  const d = new Date(iso)
  if (isNaN(d.getTime())) return iso
  return d.toLocaleDateString('ru-RU')
}
export function fmtDateTime(iso: string) {
  if (!iso) return '—'
  const d = new Date(iso)
  if (isNaN(d.getTime())) return iso
  return d.toLocaleDateString('ru-RU') + ' · ' + d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
}
export function ageYears(birthDate: string) {
  const b = new Date(birthDate)
  const now = new Date()
  let years = now.getFullYear() - b.getFullYear()
  const notYet = now.getMonth() < b.getMonth() || (now.getMonth() === b.getMonth() && now.getDate() < b.getDate())
  if (notYet) years--
  return years
}
export function yearsText(n: number) {
  const last = n % 10
  const last2 = n % 100
  if (last2 >= 11 && last2 <= 14) return n + ' лет'
  if (last === 1) return n + ' год'
  if (last >= 2 && last <= 4) return n + ' года'
  return n + ' лет'
}
export function money(n: number) {
  return n.toLocaleString('ru-RU') + ' ₽'
}
export function fullName(m: Mother) {
  return m.last_name + ' ' + m.first_name
}
export function namesById(mothers: Mother[]) {
  const result: Record<number, string> = {}
  for (const m of mothers) {
    result[m.id] = fullName(m)
  }
  return result
}
export function todayString() {
  return new Date().toISOString().slice(0, 10)
}
