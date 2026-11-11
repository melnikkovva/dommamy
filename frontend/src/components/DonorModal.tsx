import { useState } from 'react'
import type { FormEvent } from 'react'
import Modal from './Modal'
import Field from './Field'
import { donorKindNames } from '../utils/labels'
import type { DonorBody } from '../types'

type Props = {
  title: string
  initial: DonorBody
  onClose: () => void
  onSubmit: (body: DonorBody) => Promise<void>
}

export default function DonorModal({ title, initial, onClose, onSubmit }: Props) {
  const [form, setForm] = useState<DonorBody>(initial)
  function change(field: keyof DonorBody, value: string) {
    setForm({ ...form, [field]: value })
  }
  async function submit(e: FormEvent) {
    e.preventDefault()
    await onSubmit(form)
  }
  return (
    <Modal title={title} onClose={onClose}>
      <form onSubmit={submit}>
        <Field label="Имя или название *"><input value={form.name} onChange={(e) => change('name', e.target.value)} required /></Field>
        <Field label="Кто это">
          <select value={form.kind} onChange={(e) => change('kind', e.target.value)}>
            {Object.keys(donorKindNames).map((key) => (
              <option key={key} value={key}>{donorKindNames[key]}</option>
            ))}
          </select>
        </Field>
        <Field label="Телефон"><input value={form.phone} onChange={(e) => change('phone', e.target.value)} /></Field>
        <Field label="Email"><input value={form.email} onChange={(e) => change('email', e.target.value)} /></Field>
        <Field label="Комментарий"><textarea rows={2} value={form.comment} onChange={(e) => change('comment', e.target.value)} /></Field>
        <div className="dialog-buttons">
          <button type="button" className="btn" onClick={onClose}>Отмена</button>
          <button type="submit" className="btn btn-primary">Сохранить</button>
        </div>
      </form>
    </Modal>
  )
}
