import { useState } from 'react'
import type { FormEvent } from 'react'
import { createDocument, documentFileUrl, getDocuments, signDocument, uploadDocument } from '../../api/documents'
import { getHelp } from '../../api/help'
import { errorText } from '../../api/client'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import Modal from '../../components/Modal'
import Field from '../../components/Field'
import { DocBadge, docTypeText } from '../../components/Badges'
import { useToast } from '../../components/Toast'
import { fmtDate } from '../../utils/format'
import { docTypeNames, kindNames, label, sumText } from '../../utils/labels'
import type { HelpItem } from '../../types'

function NewDocumentModal({ motherId, help, onClose, onSaved }: { motherId: number; help: HelpItem[]; onClose: () => void; onSaved: () => void }) {
  const toast = useToast()
  const [type, setType] = useState('consent')
  const [helpId, setHelpId] = useState(help.length > 0 ? String(help[0].id) : '')
  async function submit(e: FormEvent) {
    e.preventDefault()
    try {
      await createDocument(motherId, type, type === 'responsibility' ? Number(helpId) : 0)
      toast('Документ создан')
      onSaved()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <Modal title="Новый документ" onClose={onClose}>
      <form onSubmit={submit}>
        <Field label="Тип документа">
          <select value={type} onChange={(e) => setType(e.target.value)}>
            {Object.keys(docTypeNames).map((key) => (
              <option key={key} value={key}>{docTypeNames[key]}</option>
            ))}
          </select>
        </Field>
        {type === 'responsibility' && (
          <Field label="К какой выдаче помощи">
            <select value={helpId} onChange={(e) => setHelpId(e.target.value)} required>
              {help.map((h) => (
                <option key={h.id} value={h.id}>№ {h.id} · {label(kindNames, h.kind)} · {sumText(h.kind, h.sum)} · {fmtDate(h.created_at)}</option>
              ))}
            </select>
          </Field>
        )}
        <div className="dialog-buttons">
          <button type="button" className="btn" onClick={onClose}>Отмена</button>
          <button type="submit" className="btn btn-primary">Создать</button>
        </div>
      </form>
    </Modal>
  )
}

export default function DocumentsTab({ motherId }: { motherId: number }) {
  const toast = useToast()
  const [showNew, setShowNew] = useState(false)
  const { data, error, reload } = useLoad(async () => {
    const [docs, help] = await Promise.all([getDocuments(motherId), getHelp()])
    return { docs, help: help.filter((h) => h.mother_id === motherId) }
  }, [motherId])
  if (!data) return <LoadState error={error} />
  async function upload(id: number, file: File) {
    try {
      await uploadDocument(id, file)
      toast('Скан загружен')
      reload()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  async function sign(id: number) {
    try {
      await signDocument(id)
      toast('Документ подписан')
      reload()
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <div className="card">
      <div className="row-title">
        <h2>Документы · {data.docs.length}</h2>
        <button className="btn btn-primary btn-small" onClick={() => setShowNew(true)}>+ Новый документ</button>
      </div>
      <div className="table-wrap">
        <table>
          <thead><tr><th>Документ</th><th>Статус</th><th>Скан</th><th>Действия</th></tr></thead>
          <tbody>
            {data.docs.map((d) => (
              <tr key={d.id}>
                <td>
                  <strong>{docTypeText(d.type)}</strong>
                  {d.help_id > 0 && <div className="muted">Выдача № {d.help_id}</div>}
                </td>
                <td>
                  <DocBadge status={d.status} />
                  {d.signed_at && <div className="muted">{fmtDate(d.signed_at)}</div>}
                </td>
                <td>
                  {d.has_file ? <a href={documentFileUrl(d.id)} target="_blank" rel="noreferrer">{d.original_name}</a> : '—'}
                </td>
                <td>
                  {d.status !== 'signed' && (
                    <div className="inline-actions">
                      <label className="btn btn-small file-btn">
                        {d.has_file ? 'Заменить скан' : 'Загрузить скан'}
                        <input
                          type="file"
                          accept=".pdf,.jpg,.jpeg,.png"
                          onChange={(e) => {
                            const file = e.target.files ? e.target.files[0] : null
                            if (file) upload(d.id, file)
                            e.target.value = ''
                          }}
                        />
                      </label>
                      {d.has_file && <button className="btn btn-small btn-primary" onClick={() => sign(d.id)}>Подписан</button>}
                    </div>
                  )}
                </td>
              </tr>
            ))}
            {data.docs.length === 0 && <tr><td colSpan={4} className="muted">Документов пока нет</td></tr>}
          </tbody>
        </table>
      </div>
      {showNew && (
        <NewDocumentModal motherId={motherId} help={data.help} onClose={() => setShowNew(false)} onSaved={() => { setShowNew(false); reload() }} />
      )}
    </div>
  )
}
