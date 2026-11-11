import { useState } from 'react'
import type { FormEvent } from 'react'
import { addChild } from '../../api/mothers'
import { errorText } from '../../api/client'
import Modal from '../../components/Modal'
import Field from '../../components/Field'
import { useToast } from '../../components/Toast'
import { ageYears, fmtDate, yearsText } from '../../utils/format'
import type { Child } from '../../types'

export default function ChildrenCard({ motherId, items, onChanged }: { motherId: number; items: Child[]; onChanged: () => void }) {
  const toast = useToast()
  const [adding, setAdding] = useState(false)
  const [firstName, setFirstName] = useState('')
  const [birthDate, setBirthDate] = useState('')
  async function submit(e: FormEvent) {
    e.preventDefault()
    try {
      await addChild(motherId, firstName, birthDate)
      toast('Ребёнок добавлен')
      setAdding(false)
      setFirstName('')
      setBirthDate('')
      onChanged()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <div className="card">
      <div className="row-title">
        <h2>Дети</h2>
        <a href="#" onClick={(e) => { e.preventDefault(); setAdding(true) }}>Добавить ребёнка</a>
      </div>
      <table>
        <thead><tr><th>Ребёнок</th><th>Дата рождения</th></tr></thead>
        <tbody>
          {items.map((c) => (
            <tr key={c.id}>
              <td><strong>{c.first_name}</strong></td>
              <td>{fmtDate(c.birth_date)} · {yearsText(ageYears(c.birth_date))}</td>
            </tr>
          ))}
          {items.length === 0 && <tr><td colSpan={2} className="muted">Детей пока нет</td></tr>}
        </tbody>
      </table>
      {adding && (
        <Modal title="Добавить ребёнка" onClose={() => setAdding(false)}>
          <form onSubmit={submit}>
            <Field label="Имя *"><input value={firstName} onChange={(e) => setFirstName(e.target.value)} required /></Field>
            <Field label="Дата рождения *"><input type="date" value={birthDate} onChange={(e) => setBirthDate(e.target.value)} required /></Field>
            <div className="dialog-buttons">
              <button type="button" className="btn" onClick={() => setAdding(false)}>Отмена</button>
              <button type="submit" className="btn btn-primary">Добавить</button>
            </div>
          </form>
        </Modal>
      )}
    </div>
  )
}
