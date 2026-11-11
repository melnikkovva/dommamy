import { useState } from 'react'
import type { FormEvent } from 'react'
import { setMotherStatus } from '../../api/mothers'
import { errorText } from '../../api/client'
import Modal from '../../components/Modal'
import Field from '../../components/Field'
import { MotherStatusBadge } from '../../components/Badges'
import { useToast } from '../../components/Toast'
import { fmtDate } from '../../utils/format'
import { label, motherStatusNames } from '../../utils/labels'
import type { MotherStatusInfo } from '../../types'

const steps: Record<string, [string, string][]> = {
  review: [['active', 'Принять на сопровождение'], ['inactive', 'Помощь не оказываем']],
  active: [['inactive', 'Завершить сопровождение']],
  inactive: [['review', 'Пересмотреть решение']]
}

export default function StatusCard({ motherId, status, onChanged }: { motherId: number; status: MotherStatusInfo; onChanged: () => void }) {
  const toast = useToast()
  const [askReason, setAskReason] = useState(false)
  const [reason, setReason] = useState('')
  async function change(to: string, text: string) {
    try {
      await setMotherStatus(motherId, to, text)
      toast('Статус изменён')
      setAskReason(false)
      setReason('')
      onChanged()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  function submitReason(e: FormEvent) {
    e.preventDefault()
    change('inactive', reason)
  }
  return (
    <div className="card">
      <h2>Статус <MotherStatusBadge status={status.status} /></h2>
      <div className="actions">
        {(steps[status.status] || []).map(([to, text]) => (
          <button key={to} className="btn btn-small" onClick={() => (to === 'inactive' ? setAskReason(true) : change(to, ''))}>{text}</button>
        ))}
      </div>
      <div className="muted">История</div>
      {status.history.map((h, i) => (
        <div key={i}>
          <div className="field-label">{fmtDate(h.created_at)} · {label(motherStatusNames, h.status)}</div>
          <div>{h.reason || '—'}</div>
        </div>
      ))}
      {status.history.length === 0 && <p className="muted">Смен статуса ещё не было</p>}
      {askReason && (
        <Modal title="Помощь не оказываем" onClose={() => setAskReason(false)}>
          <form onSubmit={submitReason}>
            <Field label="Причина *"><textarea rows={4} value={reason} onChange={(e) => setReason(e.target.value)} required /></Field>
            <div className="dialog-buttons">
              <button type="button" className="btn" onClick={() => setAskReason(false)}>Отмена</button>
              <button type="submit" className="btn btn-primary">Сохранить</button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  )
}
