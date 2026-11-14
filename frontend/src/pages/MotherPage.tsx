import { useState, useRef } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getChildren, getMothers, getMotherStatus, getProfile, uploadPhoto } from '../api/mothers'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import { MotherStatusBadge } from '../components/Badges'
import { useToast } from '../components/Toast'
import InfoTab from './mother/InfoTab'
import DocumentsTab from './mother/DocumentsTab'
import HelpTab from './mother/HelpTab'
import EventsTab from './mother/EventsTab'
import RequestsTab from './mother/RequestsTab'
import CoursesTab from './mother/CoursesTab'
import { ageYears, fmtDate, fullName, yearsText } from '../utils/format'

const tabs: [string, string][] = [
  ['info', 'Анкета и дети'],
  ['documents', 'Документы'],
  ['help', 'История помощи'],
  ['events', 'Мероприятия'],
  ['courses', 'Курсы'],
  ['requests', 'Заметки']
]

export default function MotherPage() {
  const id = Number(useParams().id)
  const [tab, setTab] = useState('info')
  const [viewPhoto, setViewPhoto] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const toast = useToast()
  const { data, error, reload } = useLoad(async () => {
    const [mothers, children, profile, status] = await Promise.all([getMothers(), getChildren(id), getProfile(id), getMotherStatus(id)])
    return { mother: mothers.find((m) => m.id === id) || null, children, profile, status }
  }, [id])
  if (!data) return <LoadState error={error} />
  const mother = data.mother
  if (!mother) {
    return <div className="card">Такой мамы нет. <Link to="/mothers">К списку</Link></div>
  }
  const initials = ((mother.last_name[0] || '') + (mother.first_name[0] || '')).toUpperCase()
  const childrenText = data.children.map((c) => c.first_name + ', ' + yearsText(ageYears(c.birth_date))).join(' · ')
  const hasPhoto = !!mother.photo_url

  async function handlePhotoUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    if (file.size > 5 * 1024 * 1024) {
      toast('Фото не должно превышать 5 МБ', true)
      return
    }
    try {
      await uploadPhoto(id, file)
      toast('Фото загружено')
      reload()
    } catch (err: any) {
      toast(err.message || 'Ошибка загрузки', true)
    }
    if (fileInputRef.current) fileInputRef.current.value = ''
  }

  return (
    <div>
      <div className="page-head">
        <div>
          <h1>{fullName(mother)}</h1>
          <p className="muted"><Link to="/mothers">← Все подопечные</Link> </p>
        </div>
        <Link className="btn btn-primary" to={'/help?mother=' + mother.id}>+ Выдать помощь</Link>
      </div>
      <div className="card summary">
        <div className="avatar-wrap">
          {hasPhoto ? (
            <img className="avatar" src={mother.photo_url} alt="" onClick={() => setViewPhoto(mother.photo_url || null)} />
          ) : (
            <div className="avatar">{initials}</div>
          )}
          <input ref={fileInputRef} type="file" accept="image/jpeg,image/png,image/webp" onChange={handlePhotoUpload} style={{ display: 'none' }} />
          <button className="btn-upload-photo" onClick={() => fileInputRef.current?.click()} title="Загрузить фото">+</button>
        </div>
        <div><div className="label">Контакт</div>{mother.phone || '—'}</div>
        <div><div className="label">Дети</div>{childrenText || '—'}</div>
        <div><div className="label">Обратилась</div>{fmtDate(mother.created_at)}</div>
        <div><div className="label">Статус</div><MotherStatusBadge status={data.status.status} /></div>
      </div>
      <div className="tabs">
        {tabs.map(([key, text]) => (
          <button key={key} className={key === tab ? 'active' : ''} onClick={() => setTab(key)}>{text}</button>
        ))}
      </div>
      {tab === 'info' && <InfoTab motherId={id} profile={data.profile} status={data.status} childList={data.children} onChanged={reload} />}
      {tab === 'documents' && <DocumentsTab motherId={id} />}
      {tab === 'help' && <HelpTab motherId={id} />}
      {tab === 'events' && <EventsTab motherId={id} />}
      {tab === 'courses' && <CoursesTab motherId={id} />}
      {tab === 'requests' && <RequestsTab motherId={id} />}
      {viewPhoto && (
        <div className="modal-backdrop" onClick={() => setViewPhoto(null)}>
          <img className="photo-view" src={viewPhoto} alt="Фото" onClick={(e) => e.stopPropagation()} />
        </div>
      )}
    </div>
  )
}
