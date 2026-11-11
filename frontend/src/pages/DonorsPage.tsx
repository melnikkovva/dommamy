import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createDonor, getDonors } from '../api/donors'
import { errorText } from '../api/client'
import { useLoad } from '../hooks/useLoad'
import LoadState from '../components/LoadState'
import DonorModal from '../components/DonorModal'
import { useToast } from '../components/Toast'
import { money } from '../utils/format'
import { donorKindNames, label } from '../utils/labels'
import type { DonorBody } from '../types'

const emptyDonor: DonorBody = { name: '', kind: 'person', phone: '', email: '', comment: '' }

export default function DonorsPage() {
  const navigate = useNavigate()
  const toast = useToast()
  const [query, setQuery] = useState('')
  const [search, setSearch] = useState('')
  const [showNew, setShowNew] = useState(false)
  useEffect(() => {
    const timer = setTimeout(() => setSearch(query), 300)
    return () => clearTimeout(timer)
  }, [query])
  const { data, error } = useLoad(() => getDonors(search), [search])
  async function create(body: DonorBody) {
    try {
      const donor = await createDonor(body)
      toast('Даритель добавлен')
      navigate('/donors/' + donor.id)
    } catch (err) {
      toast(errorText(err), true)
    }
  }
  return (
    <div>
      <div className="page-head">
        <h1>Дарители</h1>
        <button className="btn btn-primary" onClick={() => setShowNew(true)}>Новый даритель</button>
      </div>
      <input className="search" placeholder="Найти по имени или телефону" value={query} onChange={(e) => setQuery(e.target.value)} />
      {!data ? (
        <LoadState error={error} />
      ) : (
        <div className="card">
          <div className="table-wrap">
            <table>
              <thead><tr><th>Даритель</th><th>Кто</th><th>Телефон</th><th>Пожертвований</th><th>Всего</th></tr></thead>
              <tbody>
                {data.map((d) => (
                  <tr key={d.id} className="clickable" onClick={() => navigate('/donors/' + d.id)}>
                    <td><strong>{d.name}</strong></td>
                    <td>{label(donorKindNames, d.kind)}</td>
                    <td>{d.phone || '—'}</td>
                    <td>{d.count}</td>
                    <td>{money(d.total)}</td>
                  </tr>
                ))}
                {data.length === 0 && <tr><td colSpan={5} className="muted">Дарителей пока нет</td></tr>}
              </tbody>
            </table>
          </div>
        </div>
      )}
      {showNew && <DonorModal title="Новый даритель" initial={emptyDonor} onClose={() => setShowNew(false)} onSubmit={create} />}
    </div>
  )
}
