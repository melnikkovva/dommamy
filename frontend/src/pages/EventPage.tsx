import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getMothers } from '../api/mothers'
import { getEvents, getParticipants, markParticipant, registerMother } from '../api/events'
import { errorText } from '../api/client'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import { ParticipantBadge } from '../components/Badges'
import { useToast } from '../components/Toast'
import { fmtDateTime, fullName, namesById } from '../utils/format'

export default function EventPage() {
  const id = Number(useParams().id)
  const toast = useToast()
  const [motherId, setMotherId] = useState('')
  const { data, error, reload } = useLoad(async () => {
    const [events, mothers, participants] = await Promise.all([getEvents(), getMothers(), getParticipants(id)])
    return { event: events.find((e) => e.id === id) || null, mothers, participants }
  }, [id])
  if (!data) return <LoadState error={error} />
  const event = data.event
  if (!event) {
    return <div className="card">Такого мероприятия нет. <Link to="/events">К списку</Link></div>
  }
  const names = namesById(data.mothers)
  const currentMother = motherId || (data.mothers.length > 0 ? String(data.mothers[0].id) : '')
  async function join() {
    if (!currentMother) return
    try {
      await registerMother(id, Number(currentMother))
      toast('Мама записана')
      reload()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  async function mark(mother: number, status: string) {
    try {
      await markParticipant(id, mother, status)
      reload()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <div>
      <div className="page-head">
        <div>
          <h1>{event.title}</h1>
          <p className="muted">
            <Link to="/events">Все мероприятия</Link> · {fmtDateTime(event.starts_at)} · {event.place || 'место не указано'} · мест: {event.max_people || 'без лимита'}
          </p>
        </div>
      </div>
      <div className="card">
        <h2>Записать маму</h2>
        <div className="actions">
          <select style={{ maxWidth: 320 }} value={currentMother} onChange={(e) => setMotherId(e.target.value)}>
            {data.mothers.map((m) => (
              <option key={m.id} value={m.id}>{fullName(m)}</option>
            ))}
          </select>
          <button className="btn btn-primary" onClick={join}>Записать</button>
        </div>
      </div>
      <div className="card">
        <h2>Участницы {data.participants.length}</h2>
        <table>
          <thead><tr><th>Мама</th><th>Статус</th><th>Отметить</th></tr></thead>
          <tbody>
            {data.participants.map((p) => (
              <tr key={p.mother_id}>
                <td><Link to={'/mothers/' + p.mother_id}><strong>{names[p.mother_id] || 'Мама'}</strong></Link></td>
                <td><ParticipantBadge status={p.status} /></td>
                <td>
                  <div className="inline-actions">
                    <button className="btn btn-small" onClick={() => mark(p.mother_id, 'came')}>Пришла</button>
                    <button className="btn btn-small" onClick={() => mark(p.mother_id, 'absent')}>Не пришла</button>
                  </div>
                </td>
              </tr>
            ))}
            {data.participants.length === 0 && <tr><td colSpan={3} className="muted">Пока никто не записан</td></tr>}
          </tbody>
        </table>
      </div>
    </div>
  )
}
