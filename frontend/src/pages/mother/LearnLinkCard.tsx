import { useState, useEffect } from 'react'
import { getLearners, linkLearnUser, unlinkLearnUser } from '../../api/mothers'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { useToast } from '../../components/Toast'

export default function LearnLinkCard({ motherId, learnUserId }: { motherId: number; learnUserId?: string }) {
  const toast = useToast()
  const [selectedUserId, setSelectedUserId] = useState(learnUserId || '')
  const { data: learners, error, reload } = useLoad(() => getLearners(), [])
  
  useEffect(() => {
    setSelectedUserId(learnUserId || '')
  }, [learnUserId])

  const handleLink = async () => {
    if (!selectedUserId) {
      toast('Выберите пользователя из Лерна', true)
      return
    }
    try {
      await linkLearnUser(motherId, selectedUserId)
      toast('Связано с Лерном')
      reload()
    } catch (err: any) {
      toast(err.message || 'Ошибка', true)
    }
  }

  const handleUnlink = async () => {
    try {
      await unlinkLearnUser(motherId)
      setSelectedUserId('')
      toast('Связь с Лерном удалена')
      reload()
    } catch (err: any) {
      toast(err.message || 'Ошибка', true)
    }
  }

  if (!learners) return <LoadState error={error} />

  const linkedLearner = learners.find((l: { user_id: string; first_name: string; last_name: string; email: string }) => l.user_id === learnUserId)

  return (
    <div className="card">
      <h3>Связь с iSpring Learn</h3>
      
      {linkedLearner && (
        <div style={{ marginBottom: '12px', padding: '8px', background: '#f0fdf4', borderRadius: '4px' }}>
          <strong>Связано:</strong> {linkedLearner.first_name} {linkedLearner.last_name} ({linkedLearner.email})
          <button className="btn btn-small" onClick={handleUnlink} style={{ marginLeft: '12px' }}>
            Отвязать
          </button>
        </div>
      )}

      <div style={{ display: 'flex', gap: '8px', alignItems: 'center' }}>
        <select
          value={selectedUserId}
          onChange={(e) => setSelectedUserId(e.target.value)}
          style={{ flex: 1, padding: '6px 8px' }}
        >
          <option value="">— Выберите пользователя —</option>
          {learners.map((l: { user_id: string; first_name: string; last_name: string; email: string }) => (
            <option key={l.user_id} value={l.user_id}>
              {l.first_name} {l.last_name} ({l.email})
            </option>
          ))}
        </select>
        <button className="btn btn-small btn-primary" onClick={handleLink}>
          {linkedLearner ? 'Изменить' : 'Связать'}
        </button>
      </div>
    </div>
  )
}
