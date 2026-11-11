import { Link } from 'react-router-dom'
import { getMotherRequests } from '../../api/help'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { RequestBadge } from '../../components/Badges'
import { fmtDate } from '../../utils/format'
import { kindNames, label } from '../../utils/labels'

export default function RequestsTab({ motherId }: { motherId: number }) {
  const { data, error } = useLoad(() => getMotherRequests(motherId), [motherId])
  if (!data) return <LoadState error={error} />
  return (
    <div className="card">
      <div className="row-title">
        <h2>Запросы на помощь · {data.length}</h2>
        <Link to="/requests">Все запросы →</Link>
      </div>
      <table>
        <thead><tr><th>Дата</th><th>Запрос</th><th>Статус</th><th>Комментарий</th></tr></thead>
        <tbody>
          {data.map((r) => (
            <tr key={r.id}>
              <td>{fmtDate(r.created_at)}</td>
              <td><strong>{label(kindNames, r.kind)}</strong><br /><span className="muted">{r.description}</span></td>
              <td><RequestBadge status={r.status} /></td>
              <td>{r.comment || '—'}</td>
            </tr>
          ))}
          {data.length === 0 && <tr><td colSpan={4} className="muted">Запросов пока нет</td></tr>}
        </tbody>
      </table>
    </div>
  )
}
