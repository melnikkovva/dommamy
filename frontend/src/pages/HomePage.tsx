import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { getMothers, getMotherStatuses, getProfiles } from '../api/mothers'
import { getHelp, getRequests } from '../api/help'
import { getEvents } from '../api/events'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import { ProfileBadge, RequestBadge } from '../components/Badges'
import { fmtDate, fmtDateTime, namesById } from '../utils/format'
import { kindNames, label, sumText } from '../utils/labels'

export default function HomePage() {
  const navigate = useNavigate()
  const [search, setSearch] = useState('')
  const { data, error } = useLoad(async () => {
    const [mothers, statuses, profiles, requests, events, help] = await Promise.all([
      getMothers(),
      getMotherStatuses(),
      getProfiles('sent'),
      getRequests('new'),
      getEvents(),
      getHelp()
    ])
    return { mothers, statuses, profiles, requests, events, help }
  }, [])
  if (!data) return <LoadState error={error} />
  const names = namesById(data.mothers)
  const activeCount = data.statuses.filter((s) => s.status === 'active').length
  const now = new Date()
  const upcoming = data.events
    .filter((e) => new Date(e.starts_at) >= now)
    .sort((a, b) => (a.starts_at < b.starts_at ? -1 : 1))
  const lastHelp = [...data.help].sort((a, b) => b.id - a.id).slice(0, 5)
  return (
    <div>
      <input
        className="search"
        placeholder="Найти маму по имени или фамилии"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === 'Enter') navigate('/mothers?q=' + encodeURIComponent(search.trim()))
        }}
      />
      <div className="actions">
        <Link className="btn btn-primary" to="/mothers?new=1">+ Новая подопечная</Link>
        <Link className="btn" to="/help">Выдать помощь</Link>
        <Link className="btn" to="/events?new=1">Новое мероприятие</Link>
      </div>
      <div className="stats">
        <div className="stat"><div className="label">Семьи на сопровождении</div><div className="value">{activeCount}</div></div>
        <div className="stat"><div className="label">Анкеты ждут проверки</div><div className="value">{data.profiles.length}</div></div>
        <div className="stat"><div className="label">Новые запросы на помощь</div><div className="value">{data.requests.length}</div></div>
        <div className="stat"><div className="label">Мероприятий впереди</div><div className="value">{upcoming.length}</div></div>
      </div>
      <div className="grid-2">
        <div>
          <div className="card">
            <div className="row-title"><h2>Анкеты на проверку</h2><Link to="/mothers">Все подопечные →</Link></div>
            <table>
              <thead><tr><th>Мама</th><th>Обновлена</th><th>Статус</th></tr></thead>
              <tbody>
                {data.profiles.slice(0, 5).map((p) => (
                  <tr key={p.mother_id}>
                    <td><Link to={'/mothers/' + p.mother_id}><strong>{names[p.mother_id] || 'Мама № ' + p.mother_id}</strong></Link></td>
                    <td>{fmtDate(p.updated_at)}</td>
                    <td><ProfileBadge status={p.status} /></td>
                  </tr>
                ))}
                {data.profiles.length === 0 && <tr><td colSpan={3} className="muted">Новых анкет нет</td></tr>}
              </tbody>
            </table>
          </div>
          <div className="card">
            <div className="row-title"><h2>Новые запросы</h2><Link to="/requests">Все запросы →</Link></div>
            <table>
              <thead><tr><th>Мама и запрос</th><th>Дата</th><th>Статус</th></tr></thead>
              <tbody>
                {data.requests.slice(0, 5).map((r) => (
                  <tr key={r.id}>
                    <td><strong>{names[r.mother_id] || 'Мама № ' + r.mother_id}</strong>  {label(kindNames, r.kind)}</td>
                    <td>{fmtDate(r.created_at)}</td>
                    <td><RequestBadge status={r.status} /></td>
                  </tr>
                ))}
                {data.requests.length === 0 && <tr><td colSpan={3} className="muted">Новых запросов нет</td></tr>}
              </tbody>
            </table>
          </div>
          <div className="card">
            <div className="row-title"><h2>Последняя помощь</h2><Link to="/help">Вся помощь →</Link></div>
            <table>
              <thead><tr><th>Мама</th><th>Вид</th><th>Объём</th><th>Дата</th></tr></thead>
              <tbody>
                {lastHelp.map((h) => (
                  <tr key={h.id}>
                    <td><strong>{names[h.mother_id] || 'Мама № ' + h.mother_id}</strong></td>
                    <td>{label(kindNames, h.kind)}</td>
                    <td>{sumText(h.kind, h.sum)}</td>
                    <td>{fmtDate(h.created_at)}</td>
                  </tr>
                ))}
                {lastHelp.length === 0 && <tr><td colSpan={4} className="muted">Помощь ещё не выдавали</td></tr>}
              </tbody>
            </table>
          </div>
        </div>
        <div className="card">
          <div className="row-title"><h2>Ближайшие мероприятия</h2><Link to="/events">Все →</Link></div>
          {upcoming.slice(0, 3).map((e) => (
            <div className="event-item" key={e.id}>
              <div className="event-date">{fmtDateTime(e.starts_at)}</div>
              <div className="event-name"><Link to={'/events/' + e.id}>{e.title}</Link></div>
              <div className="muted">{e.place || 'Место не указано'}</div>
            </div>
          ))}
          {upcoming.length === 0 && <p className="muted">Ближайших мероприятий нет</p>}
        </div>
      </div>
    </div>
  )
}
