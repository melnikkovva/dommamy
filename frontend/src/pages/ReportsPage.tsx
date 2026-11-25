import { useState } from 'react'
import { downloadReportCsv, getReport } from '../api/reports'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import Field from '../components/Field'
import { errorText } from '../api/client'
import { useToast } from '../components/Toast'

const reports: [string, string][] = [
  ['aid', 'Оказанная помощь'],
  ['events', 'Участие в мероприятиях'],
  ['mothers', 'Подопечные']
]

const problemTypes: [string, string][] = [
  ['', 'Все типы'],
  ['housing', 'Жилищная'],
  ['unemployment', 'Безработица'],
  ['single_parent', 'Одинокая мама'],
  ['low_income', 'Низкий доход'],
  ['violence', 'Насилие'],
  ['other', 'Другое']
]
const problemTypeTotalLabels = new Set(['Не указана', ...problemTypes.slice(1).map(([, text]) => text)])

export default function ReportsPage() {
  const toast = useToast()
  const [name, setName] = useState('aid')
  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')
  const [problemType, setProblemType] = useState('')
  const [downloading, setDownloading] = useState(false)
  const { data, error } = useLoad(() => getReport(name, from, to, problemType), [name, from, to, problemType])

  async function download() {
    setDownloading(true)
    try {
      await downloadReportCsv(name, from, to, problemType)
    } catch (err) {
      toast(errorText(err), true)
    } finally {
      setDownloading(false)
    }
  }

  return (
    <div>
      <div className="page-head">
        <h1>Отчёты</h1>
        <button className="btn btn-primary" onClick={download} disabled={downloading}>
          {downloading ? 'Скачивание...' : 'Скачать в Excel'}
        </button>
      </div>
      <div className="report-tabs">
        {reports.map(([key, text]) => (
          <button key={key} className={key === name ? 'btn selected' : 'btn'} onClick={() => setName(key)}>{text}</button>
        ))}
      </div>
      <div className="period">
        <Field label="С даты"><input type="date" value={from} onChange={(e) => setFrom(e.target.value)} /></Field>
        <Field label="По дату"><input type="date" value={to} onChange={(e) => setTo(e.target.value)} /></Field>
        {name === 'mothers' && (
          <Field label="Тип проблемы">
            <select value={problemType} onChange={(e) => setProblemType(e.target.value)}>
              {problemTypes.map(([key, text]) => (
                <option key={key} value={key}>{text}</option>
              ))}
            </select>
          </Field>
        )}
        <button className="btn" onClick={() => { setFrom(''); setTo(''); setProblemType('') }}>Весь период</button>
      </div>
      {error && <div className="card error">Ошибка: {error}</div>}
      {!data && !error && <LoadState error="" />}
      {data && !error && data.name === name && (
        <div>
          <div className="stats">
            {data.totals.filter((t) =>
              !(name === 'events' && t.label === 'Разных мам') &&
              !(name === 'mothers' && (
                problemTypeTotalLabels.has(t.label) ||
                ['Всего детей', 'Всего консультаций', 'Психолог', 'Юрист', 'Батюшка', 'Центр занятости'].includes(t.label)
              ))
            ).map((t) => (
              <div className="stat" key={t.label}>
                <div className="label">{t.label}</div>
                <div className="value" style={{ fontSize: 22 }}>{t.value}</div>
              </div>
            ))}
          </div>
          <div className="card">
            <h2>{data.title} : {data.rows.length}</h2>
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
