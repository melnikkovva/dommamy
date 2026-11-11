import { getHelp } from '../../api/help'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { fmtDate } from '../../utils/format'
import { kindNames, label, sumText } from '../../utils/labels'

export default function HelpTab({ motherId }: { motherId: number }) {
  const { data, error } = useLoad(async () => {
    const all = await getHelp()
    return all.filter((h) => h.mother_id === motherId).sort((a, b) => b.id - a.id)
  }, [motherId])
  if (!data) return <LoadState error={error} />
  return (
    <div className="card">
      <h2>История помощи · {data.length}</h2>
      <table>
        <thead><tr><th>№</th><th>Дата</th><th>Вид</th><th>Объём</th><th>Чек</th></tr></thead>
        <tbody>
          {data.map((h) => (
            <tr key={h.id}>
              <td>{h.id}</td>
              <td>{fmtDate(h.created_at)}</td>
              <td><strong>{label(kindNames, h.kind)}</strong></td>
              <td>{sumText(h.kind, h.sum)}</td>
              <td>{h.check_url || '—'}</td>
            </tr>
          ))}
          {data.length === 0 && <tr><td colSpan={5} className="muted">Помощь ещё не выдавали</td></tr>}
        </tbody>
      </table>
    </div>
  )
}
