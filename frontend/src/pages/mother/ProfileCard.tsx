import { useState } from 'react'
import type { FormEvent } from 'react'
import { saveProfile, setProfileStatus } from '../../api/mothers'
import { errorText } from '../../api/client'
import Modal from '../../components/Modal'
import Field from '../../components/Field'
import { ProfileBadge } from '../../components/Badges'
import { useToast } from '../../components/Toast'
import { fmtDate } from '../../utils/format'
import type { Profile, ProfileForm } from '../../types'

type FieldInfo = { key: keyof ProfileForm; label: string; type: 'text' | 'date' | 'number' }

const fields: FieldInfo[] = [
  { key: 'birth_date', label: 'Дата рождения', type: 'date' },
  { key: 'birth_place', label: 'Место рождения', type: 'text' },
  { key: 'citizenship', label: 'Гражданство', type: 'text' },
  { key: 'email', label: 'Email', type: 'text' },
  { key: 'registration', label: 'Место регистрации', type: 'text' },
  { key: 'living', label: 'Место проживания', type: 'text' },
  { key: 'religion', label: 'Вероисповедание', type: 'text' },
  { key: 'education', label: 'Образование', type: 'text' },
  { key: 'profession', label: 'Профессия и опыт работы', type: 'text' },
  { key: 'income_source', label: 'Источник дохода', type: 'text' },
  { key: 'income_amount', label: 'Доход, ₽ в месяц', type: 'number' },
  { key: 'marital_status', label: 'Семейное положение', type: 'text' },
  { key: 'father_info', label: 'Сведения об отце ребёнка', type: 'text' },
  { key: 'pregnancy_weeks', label: 'Срок беременности, недель', type: 'number' },
  { key: 'source', label: 'Откуда узнали о центре', type: 'text' }
]

const nextSteps: Record<string, [string, string][]> = {
  draft: [['sent', 'Отправить на проверку']],
  sent: [['checked', 'Отметить проверенной'], ['draft', 'Вернуть в черновик']],
  checked: [['sent', 'Вернуть на проверку']]
}

function toForm(profile: Profile) {
  const form: Record<string, string> = { situation: profile.situation, expected_help: profile.expected_help }
  for (const f of fields) {
    const value = profile[f.key]
    form[f.key] = value === 0 ? '' : String(value)
  }
  return form
}

function ProfileModal({ motherId, profile, onClose, onSaved }: { motherId: number; profile: Profile; onClose: () => void; onSaved: () => void }) {
  const toast = useToast()
  const [form, setForm] = useState(toForm(profile))
  function change(key: string, value: string) {
    setForm({ ...form, [key]: value })
  }
  async function submit(e: FormEvent) {
    e.preventDefault()
    const body: Record<string, string | number> = { situation: form.situation, expected_help: form.expected_help }
    for (const f of fields) {
      body[f.key] = f.type === 'number' ? Number(form[f.key]) : form[f.key]
    }
    try {
      await saveProfile(motherId, body)
      toast('Анкета сохранена')
      onSaved()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <Modal title="Анкета" onClose={onClose}>
      <form onSubmit={submit}>
        <div className="form-grid">
          {fields.map((f) => (
            <Field key={f.key} label={f.label}>
              <input type={f.type} value={form[f.key]} onChange={(e) => change(f.key, e.target.value)} />
            </Field>
          ))}
        </div>
        <Field label="Опишите ситуацию"><textarea rows={3} value={form.situation} onChange={(e) => change('situation', e.target.value)} /></Field>
        <Field label="Какую помощь ждёте от центра"><textarea rows={3} value={form.expected_help} onChange={(e) => change('expected_help', e.target.value)} /></Field>
        <div className="dialog-buttons">
          <button type="button" className="btn" onClick={onClose}>Отмена</button>
          <button type="submit" className="btn btn-primary">Сохранить</button>
        </div>
      </form>
    </Modal>
  )
}

export default function ProfileCard({ motherId, profile, onChanged }: { motherId: number; profile: Profile; onChanged: () => void }) {
  const toast = useToast()
  const [editing, setEditing] = useState(false)
  async function changeStatus(status: string) {
    try {
      await setProfileStatus(motherId, status)
      onChanged()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <div className="card">
      <div className="row-title">
        <h2>Анкета <ProfileBadge status={profile.status} /></h2>
        <a href="#" onClick={(e) => { e.preventDefault(); setEditing(true) }}>Изменить</a>
      </div>
      {fields.map((f) => {
        const value = profile[f.key]
        if (value === '' || value === 0) return null
        return (
          <div key={f.key}>
            <div className="field-label">{f.label}</div>
            <div>{f.type === 'date' ? fmtDate(String(value)) : String(value)}</div>
          </div>
        )
      })}
      <div className="field-label">Что сейчас важно</div>
      <div>{profile.situation || '—'}</div>
      <div className="field-label">Какую помощь ждёт от центра</div>
      <div>{profile.expected_help || '—'}</div>
      <div className="actions" style={{ marginTop: 16 }}>
        {(nextSteps[profile.status] || []).map(([status, text]) => (
          <button key={status} className="btn btn-small" onClick={() => changeStatus(status)}>{text}</button>
        ))}
      </div>
      {editing && (
        <ProfileModal motherId={motherId} profile={profile} onClose={() => setEditing(false)} onSaved={() => { setEditing(false); onChanged() }} />
      )}
    </div>
  )
}
