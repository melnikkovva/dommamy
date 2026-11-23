import { useState } from 'react'
import { getConsultations, createConsultation, updateConsultationStatus } from '../../api/consultations'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { useToast } from '../../components/Toast'
import Field from '../../components/Field'
import { fmtDate } from '../../utils/format'

const specialistTypes: [string, string][] = [
  ['psychologist', 'Психолог'],
  ['lawyer', 'Юрист'],
  ['priest', 'Батюшка'],
  ['employment', 'Центр занятости']
]

const specialistLabels: Record<string, string> = Object.fromEntries(specialistTypes)

const statusOptions: [string, string][] = [
  ['scheduled', 'Запланирована'],
  ['completed', 'Прошла'],
  ['cancelled', 'Отменена']
]

const statusLabels: Record<string, string> = Object.fromEntries(statusOptions)

export default function ConsultationsTab({ motherId }: { motherId: number }) {
  const toast = useToast()
  const { data, error, reload } = useLoad(() => getConsultations(motherId), [motherId])
  const [showForm, setShowForm] = useState(false)
  const [specialistType, setSpecialistType] = useState('psychologist')
  const [date, setDate] = useState('')
  const [notes, setNotes] = useState('')
  const [status, setStatus] = useState('scheduled')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!date) {
      toast('Укажите дату', true)
      return
    }
    try {
      await createConsultation(motherId, specialistType, date, notes, status)
      toast('Консультация добавлена')
      setShowForm(false)
      setDate('')
      setNotes('')
      setStatus('scheduled')
      reload()
    } catch (err: any) {
      toast(err.message || 'Ошибка', true)
    }
  }

  async function changeStatus(id: number, newStatus: string) {
    try {
      await updateConsultationStatus(id, newStatus)
      toast('Статус обновлён')
      reload()
    } catch (err: any) {
      toast(err.message || 'Ошибка', true)
    }
  }

  if (!data) return <LoadState error={error} />

  const byType: Record<string, number> = {}
  for (const c of data) {
    if (c.status === 'completed') {
      byType[c.specialist_type] = (byType[c.specialist_type] || 0) + 1
    }
  }

  return (
    <div>
      <div className="card">
        <div className="card-head">
          <h3>История консультаций</h3>
          <button className="btn btn-primary" onClick={() => setShowForm(!showForm)}>
            {showForm ? 'Отмена' : '+ Записать на консультацию'}
          </button>
        </div>

        {Object.keys(byType).length > 0 && (
          <div className="consultation-summary">
            <div className="field-label">Визиты по специалистам</div>
            {specialistTypes.map(([key, label]) => byType[key] && (
              <div key={key} className="summary-item">
                <span>{label}</span>
                <span className="badge">{byType[key]}</span>
              </div>
            ))}
          </div>
        )}

        {showForm && (
          <form onSubmit={handleSubmit} className="dialog">
            <Field label="Специалист">
              <select value={specialistType} onChange={(e) => setSpecialistType(e.target.value)}>
                {specialistTypes.map(([key, text]) => (
                  <option key={key} value={key}>{text}</option>
                ))}
              </select>
            </Field>
            <Field label="Дата">
              <input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
            </Field>
            <Field label="Статус">
              <select value={status} onChange={(e) => setStatus(e.target.value)}>
                {statusOptions.map(([key, text]) => (
                  <option key={key} value={key}>{text}</option>
                ))}
              </select>
            </Field>
            <Field label="Заметки">
              <textarea rows={3} value={notes} onChange={(e) => setNotes(e.target.value)} />
            </Field>
            <div className="dialog-buttons">
              <button type="button" className="btn" onClick={() => setShowForm(false)}>Отмена</button>
              <button type="submit" className="btn btn-primary">Сохранить</button>
            </div>
          </form>
        )}

        {data.length === 0 ? (
          <p className="muted">Консультаций пока нет</p>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Дата</th>
                <th>Специалист</th>
                <th>Статус</th>
                <th>Заметки</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {data.map((c: any) => (
                <tr key={c.id}>
                  <td>{fmtDate(c.consultation_date)}</td>
                  <td>{specialistLabels[c.specialist_type] || c.specialist_type}</td>
                  <td>
                    <select
                      value={c.status}
                      onChange={(e) => changeStatus(c.id, e.target.value)}
                      className={'status-select ' + c.status}
                    >
                      {statusOptions.map(([key, text]) => (
                        <option key={key} value={key}>{text}</option>
                      ))}
                    </select>
                  </td>
                  <td className="notes">{c.notes || '—'}</td>
                  <td>
                    {c.status === 'scheduled' && (
                      <button className="btn-sm" onClick={() => changeStatus(c.id, 'completed')}>✓</button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
