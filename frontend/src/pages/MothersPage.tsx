import { useEffect, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { createMother, deleteMother, getMothers } from '../api/mothers'
import { errorText } from '../api/client'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import Modal from '../components/Modal'
import Field from '../components/Field'
import { MotherStatusBadge, ProfileBadge } from '../components/Badges'
import { useToast } from '../components/Toast'
import { fmtDate, fullName } from '../utils/format'

function ConfirmDeleteModal({ motherId, name, onDeleted, onClose }: { motherId: number; name: string; onDeleted: () => void; onClose: () => void }) {
  const toast = useToast()
  async function handleDelete() {
    try {
      await deleteMother(motherId)
      toast('Подопечная удалена')
      onDeleted()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <Modal title="Удаление подопечной" onClose={onClose}>
      <p style={{ margin: '0 0 16px' }}>
        Вы действительно хотите удалить <strong>{name}</strong>? Данные будут скрыты, но сохранены в системе.
      </p>
      <div className="dialog-buttons">
        <button type="button" className="btn" onClick={onClose}>Отмена</button>
        <button type="button" className="btn btn-danger" onClick={handleDelete}>Удалить</button>
      </div>
    </Modal>
  )
}

function NewMotherModal({ onClose }: { onClose: () => void }) {
  const navigate = useNavigate()
  const toast = useToast()
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [phone, setPhone] = useState('')
  const [email, setEmail] = useState('')
  async function submit(e: FormEvent) {
    e.preventDefault()
    try {
      const mother = await createMother(firstName, lastName, phone, email)
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
        <Field label="Email (для связи с Лерном)"><input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="email@ispring.ru" type="email" /></Field>
        <div className="dialog-buttons">
          <button type="button" className="btn" onClick={onClose}>Отмена</button>
          <button type="submit" className="btn btn-primary">Добавить</button>
        </div>
      </form>
    </Modal>
  )
}

const profileStatuses = [
  { value: '', label: 'Все анкеты' },
  { value: 'draft', label: 'Черновик' },
  { value: 'sent', label: 'Отправлена' },
  { value: 'checked', label: 'Проверена' },
]

const motherStatuses = [
  { value: '', label: 'Все статусы' },
  { value: 'review', label: 'На рассмотрении' },
  { value: 'active', label: 'Активна' },
  { value: 'inactive', label: 'Неактивна' },
]

export default function MothersPage() {
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const [query, setQuery] = useState(params.get('q') || '')
  const [phone, setPhone] = useState(params.get('phone') || '')
  const [profileStatus, setProfileStatus] = useState(params.get('profileStatus') || '')
  const [motherStatus, setMotherStatus] = useState(params.get('motherStatus') || '')
  const [search, setSearch] = useState(query)
  const [searchPhone, setSearchPhone] = useState(phone)
  const [showNew, setShowNew] = useState(params.get('new') === '1')
  const [deleting, setDeleting] = useState<{ id: number; name: string } | null>(null)
  const [showFilters, setShowFilters] = useState(false)
  const toast = useToast()
  const debounceRef = useRef<number | undefined>()

  useEffect(() => {
    debounceRef.current = setTimeout(() => setSearch(query), 300)
    return () => clearTimeout(debounceRef.current)
  }, [query])

  useEffect(() => {
    const timer = setTimeout(() => setSearchPhone(phone), 300)
    return () => clearTimeout(timer)
  }, [phone])

  const hasFilters = query || phone || profileStatus || motherStatus

  const { data, error, reload } = useLoad(async () => {
    const mothers = await getMothers(search, searchPhone, profileStatus, motherStatus)
    return { mothers }
  }, [search, searchPhone, profileStatus, motherStatus])

  function clearFilters() {
    setQuery('')
    setPhone('')
    setProfileStatus('')
    setMotherStatus('')
  }

  return (
    <div>
      <div className="page-head">
        <h1>Подопечные</h1>
        <button className="btn btn-primary" onClick={() => setShowNew(true)}>+ Новая подопечная</button>
      </div>
      <input className="search" placeholder="Найти по имени, фамилии или телефону" value={query} onChange={(e) => setQuery(e.target.value)} />
      <div style={{ marginBottom: 12 }}>
        <button className="btn" style={{ fontSize: 13, padding: '4px 12px' }} onClick={() => setShowFilters(!showFilters)}>
          {showFilters ? 'Скрыть фильтры' : 'Фильтры'}
        </button>
        {hasFilters && (
          <button className="btn" style={{ fontSize: 13, padding: '4px 12px', marginLeft: 8 }} onClick={clearFilters}>
            Сбросить
          </button>
        )}
      </div>
      {showFilters && (
        <div className="card" style={{ marginBottom: 16, display: 'flex', gap: 16, alignItems: 'flex-end', flexWrap: 'wrap' }}>
          <div style={{ flex: '1 1 180px' }}>
            <label style={{ fontSize: 13, fontWeight: 600, display: 'block', marginBottom: 4 }}>Телефон</label>
            <input className="input" style={{ width: '100%' }} value={phone} onChange={(e) => setPhone(e.target.value)} placeholder="+7 900..." />
          </div>
          <div>
            <label style={{ fontSize: 13, fontWeight: 600, display: 'block', marginBottom: 4 }}>Статус анкеты</label>
            <select className="input" style={{ width: 180 }} value={profileStatus} onChange={(e) => setProfileStatus(e.target.value)}>
              {profileStatuses.map((s) => <option key={s.value} value={s.value}>{s.label}</option>)}
            </select>
          </div>
          <div>
            <label style={{ fontSize: 13, fontWeight: 600, display: 'block', marginBottom: 4 }}>Статус сопровождения</label>
            <select className="input" style={{ width: 200 }} value={motherStatus} onChange={(e) => setMotherStatus(e.target.value)}>
              {motherStatuses.map((s) => <option key={s.value} value={s.value}>{s.label}</option>)}
            </select>
          </div>
        </div>
      )}
      {!data ? (
        <LoadState error={error} />
      ) : (
        <div className="card">
          <div className="row-title"><strong>Все подопечные {data.mothers.length}</strong></div>
          <div className="table-wrap">
            <table>
              <thead><tr><th>Номер</th><th>Мама</th><th>Телефон</th><th>Обратилась</th><th>Анкета</th><th>Статус</th><th></th></tr></thead>
              <tbody>
                {data.mothers.map((m) => (
                  <tr key={m.id}>
                    <td>№ {m.id}</td>
                    <td><strong className="clickable" onClick={() => navigate('/mothers/' + m.id)}>{fullName(m)}</strong></td>
                    <td>{m.phone || '—'}</td>
                    <td>{fmtDate(m.created_at)}</td>
                    <td>{m.profile_status ? <ProfileBadge status={m.profile_status} /> : <span className="muted">Нет анкеты</span>}</td>
                    <td><MotherStatusBadge status={m.mother_status || 'review'} /></td>
                    <td>
                      <button className="btn btn-danger" style={{ padding: '2px 8px', fontSize: '13px' }} onClick={() => setDeleting({ id: m.id, name: fullName(m) })}>Удалить</button>
                    </td>
                  </tr>
                ))}
                {data.mothers.length === 0 && <tr><td colSpan={7} className="muted">Никого не нашли</td></tr>}
              </tbody>
            </table>
          </div>
        </div>
      )}
      {showNew && <NewMotherModal onClose={() => setShowNew(false)} />}
      {deleting && (
        <ConfirmDeleteModal
          motherId={deleting.id}
          name={deleting.name}
          onDeleted={() => {
            setDeleting(null)
            reload()
          }}
          onClose={() => setDeleting(null)}
        />
      )}
    </div>
  )
}
