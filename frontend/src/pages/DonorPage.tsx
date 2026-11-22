import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useParams } from 'react-router-dom'
import { addDonation, getDonor, updateDonor } from '../api/donors'
import { errorText } from '../api/client'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import Modal from '../components/Modal'
import Field from '../components/Field'
import DonorModal from '../components/DonorModal'
import { useToast } from '../components/Toast'
import { fmtDate, money } from '../utils/format'
import { donorKindNames, label } from '../utils/labels'
import type { DonorBody } from '../types'

export default function DonorPage() {
  const id = Number(useParams().id)
  const toast = useToast()
  const [editing, setEditing] = useState(false)
  const [adding, setAdding] = useState(false)
  const [amount, setAmount] = useState('')
  const [date, setDate] = useState('')
  const [comment, setComment] = useState('')
  const { data, error, reload } = useLoad(() => getDonor(id), [id])
  if (!data) return <LoadState error={error} />
  const donor = data.donor
  const total = data.donations.reduce((sum, n) => sum + n.amount, 0)
  async function save(body: DonorBody) {
    try {
      await updateDonor(id, body)
      toast('Сохранено')
      setEditing(false)
      reload()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  async function submitDonation(e: FormEvent) {
    e.preventDefault()
    try {
      await addDonation(id, Number(amount), date, comment)
      toast('Пожертвование записано')
      setAdding(false)
      setAmount('')
      setDate('')
      setComment('')
      reload()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <div>
      <div className="page-head">
        <div>
          <h1>{donor.name}</h1>
          <p className="muted"><Link to="/donors">Все дарители</Link> {label(donorKindNames, donor.kind)}</p>
        </div>
        <div className="actions">
          <button className="btn" onClick={() => setEditing(true)}>Изменить</button>
          <button className="btn btn-primary" onClick={() => setAdding(true)}>+ Пожертвование</button>
        </div>
      </div>
      <div className="grid-2">
        <div className="card">
          <h2>Пожертвования: {data.donations.length}</h2>
          <table>
            <thead><tr><th>Дата</th><th>Сумма</th><th>Комментарий</th></tr></thead>
            <tbody>
              {data.donations.map((n) => (
                <tr key={n.id}>
                  <td>{fmtDate(n.date)}</td>
                  <td>{money(n.amount)}</td>
                  <td>{n.comment || '—'}</td>
                </tr>
              ))}
              {data.donations.length === 0 && <tr><td colSpan={3} className="muted">Пожертвований пока нет</td></tr>}
            </tbody>
          </table>
        </div>
        <div className="card">
          <h2>Контакты</h2>
          <div className="field-label">Всего пожертвовал</div><div>{money(total)}</div>
          <div className="field-label">Телефон</div><div>{donor.phone || '—'}</div>
          <div className="field-label">Email</div><div>{donor.email || '—'}</div>
          <div className="field-label">Комментарий</div><div>{donor.comment || '—'}</div>
        </div>
      </div>
      {editing && (
        <DonorModal
          title="Изменить дарителя"
          initial={{ name: donor.name, kind: donor.kind, phone: donor.phone, email: donor.email, comment: donor.comment }}
          onClose={() => setEditing(false)}
          onSubmit={save}
        />
      )}
      {adding && (
        <Modal title="Новое пожертвование" onClose={() => setAdding(false)}>
          <form onSubmit={submitDonation}>
            <Field label="Сумма, ₽ *"><input type="number" min="1" value={amount} onChange={(e) => setAmount(e.target.value)} required /></Field>
            <Field label="Дата (пусто = сегодня)"><input type="date" value={date} onChange={(e) => setDate(e.target.value)} /></Field>
            <Field label="Комментарий"><input value={comment} onChange={(e) => setComment(e.target.value)} /></Field>
            <div className="dialog-buttons">
              <button type="button" className="btn" onClick={() => setAdding(false)}>Отмена</button>
              <button type="submit" className="btn btn-primary">Записать</button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  )
}
