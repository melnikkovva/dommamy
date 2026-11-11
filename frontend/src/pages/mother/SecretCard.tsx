import { useState } from 'react'
import { getSecret, saveSecret } from '../../api/mothers'
import { ApiError, errorText } from '../../api/client'
import Field from '../../components/Field'
import { useToast } from '../../components/Toast'
import type { SecretInfo } from '../../types'

const empty: SecretInfo = { criminal_record: '', nd_pnd: '', drugs: '', chronic_diseases: '' }

export default function SecretCard({ motherId }: { motherId: number }) {
  const toast = useToast()
  const [form, setForm] = useState<SecretInfo | null>(null)
  async function reveal() {
    try {
      setForm(await getSecret(motherId))
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        setForm(empty)
      } else {
        toast(errorText(err), true)
      }
    }
  }
  async function save() {
    if (!form) return
    try {
      await saveSecret(motherId, form)
      toast('Сохранено')
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  function change(field: keyof SecretInfo, value: string) {
    if (form) setForm({ ...form, [field]: value })
  }
  return (
    <div className="card">
      <h2>Защищённые сведения</h2>
      <p className="muted">Показывайте только при необходимости. Просмотр фиксируется в журнале.</p>
      {!form ? (
        <button className="btn" onClick={reveal}>Показать данные</button>
      ) : (
        <div>
          <Field label="Сведения о судимости"><input value={form.criminal_record} onChange={(e) => change('criminal_record', e.target.value)} /></Field>
          <Field label="НД / ПНД"><input value={form.nd_pnd} onChange={(e) => change('nd_pnd', e.target.value)} /></Field>
          <Field label="Употребление наркотиков"><input value={form.drugs} onChange={(e) => change('drugs', e.target.value)} /></Field>
          <Field label="Хронические заболевания"><input value={form.chronic_diseases} onChange={(e) => change('chronic_diseases', e.target.value)} /></Field>
          <div className="actions">
            <button className="btn btn-primary" onClick={save}>Сохранить</button>
            <button className="btn" onClick={() => setForm(null)}>Скрыть</button>
          </div>
        </div>
      )}
    </div>
  )
}
