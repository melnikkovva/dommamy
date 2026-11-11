import { useState } from 'react'
import { getReport, reportCsvUrl } from '../api/reports'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import Field from '../components/Field'

const reports: [string, string][] = [
  ['aid', 'Оказанная помощь'],
  ['events', 'Участие в мероприятиях'],
  ['mothers', 'Подопечные']
]

export default function ReportsPage() {
  const [name, setName] = useState('aid')
  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')
  const { data, error } = useLoad(() => getReport(name, from, to), [name, from, to])
  return (
    <div>
      <div className="page-head">
        <h1>Отчёты</h1>
        <a className="btn btn-primary" href={reportCsvUrl(name, from, to)}>Скачать в Excel</a>
      </div>
      <div className="report-tabs">
        {reports.map(([key, text]) => (
          <button key={key} className={key === name ? 'btn selected' : 'btn'} onClick={() => setName(key)}>{text}</button>
        ))}
      </div>
      <div className="period">
        <Field label="С даты"><input type="date" value={from} onChange={(e) => setFrom(e.target.value)} /></Field>
        <Field label="По дату"><input type="date" value={to} onChange={(e) => setTo(e.target.value)} /></Field>
        <button className="btn" onClick={() => { setFrom(''); setTo('') }}>Весь период</button>
      </div>
      {error && <div className="card error">Ошибка: {error}</div>}
      {!data && !error && <LoadState error="" />}
      {data && !error && data.name === name && (
        <div>
          <div className="stats">
            {data.totals.map((t) => (
              <div className="stat" key={t.label}>
                <div className="label">{t.label}</div>
                <div className="value" style={{ fontSize: 22 }}>{t.value}</div>
              </div>
            ))}
          </div>
          <div className="card">
            <h2>{data.title} строк: {data.rows.length}</h2>
            <div className="table-wrap">
              <table>
                <thead><tr>{data.headers.map((h) => <th key={h}>{h}</th>)}</tr></thead>
                <tbody>
                  {data.rows.map((row, i) => (
                    <tr key={i}>{row.map((cell, j) => <td key={j}>{cell}</td>)}</tr>
                  ))}
                  {data.rows.length === 0 && <tr><td colSpan={data.headers.length} className="muted">За этот период данных нет</td></tr>}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
