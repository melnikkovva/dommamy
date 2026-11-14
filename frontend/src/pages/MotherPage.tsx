import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getChildren, getMothers, getMotherStatus, getProfile } from '../api/mothers'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import { MotherStatusBadge } from '../components/Badges'
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
        <div className="avatar">{initials}</div>
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
    </div>
  )
}
