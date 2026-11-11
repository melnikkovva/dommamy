import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { getMothers } from '../api/mothers'
import { createRequest, getRequests, setRequestStatus } from '../api/help'
import { errorText } from '../api/client'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import Modal from '../components/Modal'
import Field from '../components/Field'
import { RequestBadge } from '../components/Badges'
import { useToast } from '../components/Toast'
import { fmtDate, fullName, namesById } from '../utils/format'
import { kindNames, label, requestStatusNames } from '../utils/labels'
import type { HelpRequest, Mother } from '../types'

const steps: Record<string, [string, string][]> = {
  new: [['in_work', 'В работу'], ['rejected', 'Отклонить']],
  in_work: [['rejected', 'Отклонить']],
  rejected: [['new', 'Вернуть в новые']],
  issued: []
}

function NewRequestModal({ mothers, onClose, onSaved }: { mothers: Mother[]; onClose: () => void; onSaved: () => void }) {
  const toast = useToast()
  const [motherId, setMotherId] = useState(mothers.length > 0 ? String(mothers[0].id) : '')
  const [kind, setKind] = useState('food')
  const [description, setDescription] = useState('')
  async function submit(e: FormEvent) {
    e.preventDefault()
    try {
      await createRequest(Number(motherId), kind, description)
      toast('Запрос записан')
      onSaved()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <Modal title="Новый запрос на помощь" onClose={onClose}>
      <form onSubmit={submit}>
        <Field label="Мама *">
          <select value={motherId} onChange={(e) => setMotherId(e.target.value)} required>
            {mothers.map((m) => (
              <option key={m.id} value={m.id}>{fullName(m)}</option>
            ))}
          </select>
        </Field>
        <Field label="Что нужно *">
          <select value={kind} onChange={(e) => setKind(e.target.value)}>
            {Object.keys(kindNames).map((key) => (
              <option key={key} value={key}>{kindNames[key]}</option>
            ))}
          </select>
        </Field>
        <Field label="Описание со слов мамы">
          <textarea rows={3} value={description} onChange={(e) => setDescription(e.target.value)} />
        </Field>
        <div className="dialog-buttons">
          <button type="button" className="btn" onClick={onClose}>Отмена</button>
          <button type="submit" className="btn btn-primary">Записать</button>
        </div>
      </form>
    </Modal>
  )
}

function RejectModal({ request, onClose, onSaved }: { request: HelpRequest; onClose: () => void; onSaved: () => void }) {
  const toast = useToast()
  const [comment, setComment] = useState('')
  async function submit(e: FormEvent) {
    e.preventDefault()
    try {
      await setRequestStatus(request.id, 'rejected', comment)
      toast('Запрос отклонён')
      onSaved()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <Modal title="Отклонить запрос" onClose={onClose}>
      <form onSubmit={submit}>
        <Field label="Причина *">
          <textarea rows={3} value={comment} onChange={(e) => setComment(e.target.value)} required />
        </Field>
        <div className="dialog-buttons">
          <button type="button" className="btn" onClick={onClose}>Отмена</button>
          <button type="submit" className="btn btn-primary">Отклонить</button>
        </div>
      </form>
    </Modal>
  )
}

export default function RequestsPage() {
  const toast = useToast()
  const [status, setStatus] = useState('')
  const [showNew, setShowNew] = useState(false)
  const [rejecting, setRejecting] = useState<HelpRequest | null>(null)
  const { data, error, reload } = useLoad(async () => {
    const [mothers, requests] = await Promise.all([getMothers(), getRequests(status)])
    return { mothers, requests }
  }, [status])
  async function change(request: HelpRequest, to: string) {
    if (to === 'rejected') {
      setRejecting(request)
      return
    }
    try {
      await setRequestStatus(request.id, to, '')
      reload()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <div>
      <div className="page-head">
        <h1>Запросы на помощь</h1>
        <button className="btn btn-primary" onClick={() => setShowNew(true)}>Добавить новый запрос</button>
      </div>
      <div className="actions">
        <select style={{ maxWidth: 260 }} value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">Статус: все</option>
          {Object.keys(requestStatusNames).map((key) => (
            <option key={key} value={key}>{requestStatusNames[key]}</option>
          ))}
        </select>
      </div>
      {!data ? (
        <LoadState error={error} />
      ) : (
        <div className="card">
          <div className="table-wrap">
            <table>
              <thead><tr><th>Дата</th><th>Мама</th><th>Запрос</th><th>Статус</th><th>Комментарий</th><th>Действия</th></tr></thead>
              <tbody>
                {data.requests.map((r) => {
                  const names = namesById(data.mothers)
                  return (
                    <tr key={r.id}>
                      <td>{fmtDate(r.created_at)}</td>
                      <td><Link to={'/mothers/' + r.mother_id}><strong>{names[r.mother_id] || 'Мама № ' + r.mother_id}</strong></Link></td>
                      <td><strong>{label(kindNames, r.kind)}</strong><br /><span className="muted">{r.description}</span></td>
                      <td><RequestBadge status={r.status} /></td>
                      <td>{r.comment || '—'}</td>
                      <td>
                        <div className="inline-actions">
                          {(steps[r.status] || []).map(([to, text]) => (
                            <button key={to} className="btn btn-small" onClick={() => change(r, to)}>{text}</button>
                          ))}
                          {(r.status === 'new' || r.status === 'in_work') && (
                            <Link className="btn btn-small btn-primary" to={'/help?mother=' + r.mother_id + '&request=' + r.id}>Выдать помощь</Link>
                          )}
                          {r.status === 'issued' && <span className="muted">Выдача № {r.help_id}</span>}
                        </div>
                      </td>
                    </tr>
                  )
                })}
                {data.requests.length === 0 && <tr><td colSpan={6} className="muted">Запросов нет</td></tr>}
              </tbody>
            </table>
          </div>
        </div>
      )}
      {showNew && data && (
        <NewRequestModal mothers={data.mothers} onClose={() => setShowNew(false)} onSaved={() => { setShowNew(false); reload() }} />
      )}
      {rejecting && (
        <RejectModal request={rejecting} onClose={() => setRejecting(null)} onSaved={() => { setRejecting(null); reload() }} />
      )}
    </div>
  )
}
