import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { createMother, getMothers, getMotherStatuses, getProfiles } from '../api/mothers'
import { errorText } from '../api/client'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import Modal from '../components/Modal'
import Field from '../components/Field'
import { MotherStatusBadge, ProfileBadge } from '../components/Badges'
import { useToast } from '../components/Toast'
import { fmtDate, fullName } from '../utils/format'

function NewMotherModal({ onClose }: { onClose: () => void }) {
  const navigate = useNavigate()
  const toast = useToast()
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [phone, setPhone] = useState('')
  async function submit(e: FormEvent) {
    e.preventDefault()
    try {
      const mother = await createMother(firstName, lastName, phone)
      toast('Подопечная добавлена')
      navigate('/mothers/' + mother.id)
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <Modal title="Новая подопечная" onClose={onClose}>
      <form onSubmit={submit}>
        <Field label="Имя *"><input value={firstName} onChange={(e) => setFirstName(e.target.value)} required /></Field>
        <Field label="Фамилия *"><input value={lastName} onChange={(e) => setLastName(e.target.value)} required /></Field>
        <Field label="Телефон"><input value={phone} onChange={(e) => setPhone(e.target.value)} placeholder="+7 900 000-00-00" /></Field>
        <div className="dialog-buttons">
          <button type="button" className="btn" onClick={onClose}>Отмена</button>
          <button type="submit" className="btn btn-primary">Добавить</button>
        </div>
      </form>
    </Modal>
  )
}

export default function MothersPage() {
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const [query, setQuery] = useState(params.get('q') || '')
  const [search, setSearch] = useState(query)
  const [showNew, setShowNew] = useState(params.get('new') === '1')
  useEffect(() => {
    const timer = setTimeout(() => setSearch(query), 300)
    return () => clearTimeout(timer)
  }, [query])
  const { data, error } = useLoad(async () => {
    const [mothers, profiles, statuses] = await Promise.all([getMothers(search), getProfiles(''), getMotherStatuses()])
    const profileById: Record<number, string> = {}
    for (const p of profiles) profileById[p.mother_id] = p.status
    const statusById: Record<number, string> = {}
    for (const s of statuses) statusById[s.mother_id] = s.status
    return { mothers, profileById, statusById }
  }, [search])
  return (
    <div>
      <div className="page-head">
        <h1>Подопечные</h1>
        <button className="btn btn-primary" onClick={() => setShowNew(true)}>+ Новая подопечная</button>
      </div>
      <input className="search" placeholder="Найти по имени или фамилии" value={query} onChange={(e) => setQuery(e.target.value)} />
      {!data ? (
        <LoadState error={error} />
      ) : (
        <div className="card">
          <div className="row-title"><strong>Все подопечные {data.mothers.length}</strong></div>
          <div className="table-wrap">
            <table>
              <thead><tr><th>Номер</th><th>Мама</th><th>Телефон</th><th>Обратилась</th><th>Анкета</th><th>Статус</th></tr></thead>
              <tbody>
                {data.mothers.map((m) => (
                  <tr key={m.id} className="clickable" onClick={() => navigate('/mothers/' + m.id)}>
                    <td>№ {m.id}</td>
                    <td><strong>{fullName(m)}</strong></td>
                    <td>{m.phone || '—'}</td>
                    <td>{fmtDate(m.created_at)}</td>
                    <td>{data.profileById[m.id] ? <ProfileBadge status={data.profileById[m.id]} /> : <span className="muted">Нет анкеты</span>}</td>
                    <td><MotherStatusBadge status={data.statusById[m.id] || 'review'} /></td>
                  </tr>
                ))}
                {data.mothers.length === 0 && <tr><td colSpan={6} className="muted">Никого не нашли</td></tr>}
              </tbody>
            </table>
          </div>
        </div>
      )}
      {showNew && <NewMotherModal onClose={() => setShowNew(false)} />}
    </div>
  )
}
