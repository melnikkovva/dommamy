import { Link } from 'react-router-dom'
import { getMotherEvents } from '../../api/events'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { fmtDateTime } from '../../utils/format'

export default function EventsTab({ motherId }: { motherId: number }) {
  const { data, error } = useLoad(() => getMotherEvents(motherId), [motherId])
  if (!data) return <LoadState error={error} />
  return (
    <div className="card">
      <h2>Мероприятия · {data.length}</h2>
      <table>
        <thead><tr><th>Дата и время</th><th>Мероприятие</th><th>Место</th></tr></thead>
        <tbody>
          {data.map((e) => (
            <tr key={e.id}>
              <td>{fmtDateTime(e.starts_at)}</td>
              <td><Link to={'/events/' + e.id}><strong>{e.title}</strong></Link></td>
              <td>{e.place || '—'}</td>
            </tr>
          ))}
          {data.length === 0 && <tr><td colSpan={3} className="muted">Мама пока не записана ни на одно мероприятие</td></tr>}
        </tbody>
      </table>
    </div>
  )
}
