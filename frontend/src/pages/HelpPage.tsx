import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { getMothers } from '../api/mothers'
import { getHelp, getRequests, giveHelp, issueRequest } from '../api/help'
import { errorText } from '../api/client'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import Field from '../components/Field'
import { useToast } from '../components/Toast'
import { fmtDate, fullName, namesById } from '../utils/format'
import { kindNames, label, sumText } from '../utils/labels'

export default function HelpPage() {
  const [params] = useSearchParams()
  const toast = useToast()
  const [motherId, setMotherId] = useState(params.get('mother') || '')
  const [requestId, setRequestId] = useState(params.get('request') || '')
  const [kind, setKind] = useState('food')
  const [sum, setSum] = useState('')
  const [checkUrl, setCheckUrl] = useState('')
  const { data, error, reload } = useLoad(async () => {
    const [mothers, help, requests] = await Promise.all([getMothers(), getHelp(), getRequests('')])
    return { mothers, help, requests }
  }, [])
  if (!data) return <LoadState error={error} />
  const names = namesById(data.mothers)
  const currentMother = motherId || (data.mothers.length > 0 ? String(data.mothers[0].id) : '')
  const openRequests = data.requests.filter(
    (r) => String(r.mother_id) === currentMother && (r.status === 'new' || r.status === 'in_work')
  )
  const chosenRequest = openRequests.find((r) => String(r.id) === requestId)
  const history = [...data.help].sort((a, b) => b.id - a.id)
  const isFood = kind === 'food'
  async function submit(e: FormEvent) {
    e.preventDefault()
    let item
    try {
      item = await giveHelp(Number(currentMother), kind, Number(sum), checkUrl.trim())
    } catch (err) {
      toast(errorText(err), true)
      return
    }
    if (chosenRequest) {
      try {
        await issueRequest(chosenRequest.id, item.id)
        toast('Помощь записана, заявка закрыта')
      } catch (err) {
        toast('Помощь записана, но заявку закрыть не удалось: ' + errorText(err), true)
      }
    } else {
      toast('Помощь записана')
    }
    setSum('')
    setCheckUrl('')
    setRequestId('')
    reload()
  }
  return (
    <div>
      <div className="page-head">
        <div>
          <h1>Выдать помощь</h1>
          <p className="muted">Запишите, что передали маме. Для продуктов чек обязателен.</p>
        </div>
      </div>
      <div className="card">
        <h2>Что передали</h2>
        <form onSubmit={submit}>
          <div className="form-grid">
            <Field label="Мама *">
              <select value={currentMother} onChange={(e) => { setMotherId(e.target.value); setRequestId('') }} required>
                {data.mothers.map((m) => (
                  <option key={m.id} value={m.id}>{fullName(m)} · № {m.id}</option>
                ))}
              </select>
            </Field>
            <Field label="Вид помощи *">
              <select value={kind} onChange={(e) => setKind(e.target.value)}>
                <option value="food">Продукты</option>
                <option value="clothes">Вещи</option>
                <option value="furniture">Мебель</option>
              </select>
            </Field>
            <Field label={isFood ? 'Сумма по чеку, ₽ *' : 'Количество, шт. *'}>
              <input type="number" min="1" value={sum} onChange={(e) => setSum(e.target.value)} required />
            </Field>
            <Field label="Чек (имя файла или ссылка)">
              <input value={checkUrl} onChange={(e) => setCheckUrl(e.target.value)} />
            </Field>
            <Field label="Закрыть заявку этой выдачей">
              <select value={chosenRequest ? requestId : ''} onChange={(e) => setRequestId(e.target.value)}>
                <option value="">Без заявки</option>
                {openRequests.map((r) => (
                  <option key={r.id} value={r.id}>№ {r.id} · {label(kindNames, r.kind)} · {r.description.slice(0, 40)}</option>
                ))}
              </select>
            </Field>
          </div>
          <button type="submit" className="btn btn-primary">Записать</button>
        </form>
      </div>
      <div className="card">
        <h2>История выдач</h2>
        <div className="table-wrap">
          <table>
            <thead><tr><th>№</th><th>Дата</th><th>Мама</th><th>Вид</th><th>Объём</th><th>Чек</th></tr></thead>
            <tbody>
              {history.map((h) => (
                <tr key={h.id}>
                  <td>{h.id}</td>
                  <td>{fmtDate(h.created_at)}</td>
                  <td><Link to={'/mothers/' + h.mother_id}><strong>{names[h.mother_id] || 'Мама № ' + h.mother_id}</strong></Link></td>
                  <td>{label(kindNames, h.kind)}</td>
                  <td>{sumText(h.kind, h.sum)}</td>
                  <td>{h.check_url || '—'}</td>
                </tr>
              ))}
              {history.length === 0 && <tr><td colSpan={6} className="muted">Помощь ещё не выдавали</td></tr>}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
