import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { createEvent, getEvents, getEventTypes } from '../api/events'
import { errorText } from '../api/client'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import Modal from '../components/Modal'
import Field from '../components/Field'
import { Badge } from '../components/Badges'
import { useToast } from '../components/Toast'
import { fmtDateTime } from '../utils/format'
import type { EventType } from '../types'

function NewEventModal({ types, onClose }: { types: EventType[]; onClose: () => void }) {
  const navigate = useNavigate()
  const toast = useToast()
  const [typeId, setTypeId] = useState(types.length > 0 ? String(types[0].id) : '')
  const [title, setTitle] = useState('')
  const [startsAt, setStartsAt] = useState('')
  const [place, setPlace] = useState('')
  const [maxPeople, setMaxPeople] = useState('0')
  async function submit(e: FormEvent) {
    e.preventDefault()
    try {
      const event = await createEvent(Number(typeId), title, startsAt, place, Number(maxPeople))
      toast('Мероприятие создано')
      navigate('/events/' + event.id)
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <Modal title="Новое мероприятие" onClose={onClose}>
      <form onSubmit={submit}>
        <Field label="Тип *">
          <select value={typeId} onChange={(e) => setTypeId(e.target.value)} required>
            {types.map((t) => (
              <option key={t.id} value={t.id}>{t.name}</option>
            ))}
          </select>
        </Field>
        <Field label="Название *"><input value={title} onChange={(e) => setTitle(e.target.value)} required /></Field>
        <Field label="Дата и время *"><input type="datetime-local" value={startsAt} onChange={(e) => setStartsAt(e.target.value)} required /></Field>
        <Field label="Место"><input value={place} onChange={(e) => setPlace(e.target.value)} /></Field>
        <Field label="Сколько мест (0 = без лимита)"><input type="number" min="0" value={maxPeople} onChange={(e) => setMaxPeople(e.target.value)} /></Field>
        <div className="dialog-buttons">
          <button type="button" className="btn" onClick={onClose}>Отмена</button>
          <button type="submit" className="btn btn-primary">Создать</button>
        </div>
      </form>
    </Modal>
  )
}

export default function EventsPage() {
  const [params] = useSearchParams()
  const navigate = useNavigate()
  const [showNew, setShowNew] = useState(params.get('new') === '1')
  const { data, error } = useLoad(async () => {
    const [events, types] = await Promise.all([getEvents(), getEventTypes()])
    return { events, types }
  }, [])
  if (!data) return <LoadState error={error} />
  const now = new Date()
  const upcomingCount = data.events.filter((e) => new Date(e.starts_at) >= now).length
  return (
    <div>
      <div className="page-head">
        <h1>Мероприятия</h1>
        <button className="btn btn-primary" onClick={() => setShowNew(true)}>+ Новое мероприятие</button>
      </div>
      <div className="stats">
        <div className="stat"><div className="label">Предстоящие встречи</div><div className="value">{upcomingCount}</div></div>
        <div className="stat"><div className="label">Всего мероприятий</div><div className="value">{data.events.length}</div></div>
      </div>
      <div className="card">
        <h2>Предстоящие и прошедшие</h2>
        <div className="table-wrap">
          <table>
            <thead><tr><th>Мероприятие / тип</th><th>Дата и время</th><th>Место</th><th>Мест</th><th>Статус</th></tr></thead>
            <tbody>
              {data.events.map((e) => {
                const isPast = new Date(e.starts_at) < now
                return (
                  <tr key={e.id} className="clickable" onClick={() => navigate('/events/' + e.id)}>
                    <td><strong>{e.title}</strong><br /><span className="muted">{e.type_name || '—'}</span></td>
                    <td>{fmtDateTime(e.starts_at)}</td>
                    <td>{e.place || '—'}</td>
                    <td>{e.max_people || 'без лимита'}</td>
                    <td><Badge text={isPast ? 'Состоялось' : 'Запись открыта'} kind={isPast ? 'muted' : 'ok'} /></td>
                  </tr>
                )
              })}
              {data.events.length === 0 && <tr><td colSpan={5} className="muted">Мероприятий пока нет</td></tr>}
            </tbody>
          </table>
        </div>
      </div>
      {showNew && <NewEventModal types={data.types} onClose={() => setShowNew(false)} />}
    </div>
  )
}
