import { useState } from 'react'
import { useLoad } from '../hooks/useLoad'
import { listUsers, createUser, UserResponse } from '../api/auth'
import { errorText } from '../api/client'

export default function UsersPage() {
  const { data: users, error, reload } = useLoad(listUsers, [])
  const [showForm, setShowForm] = useState(false)
  const [formError, setFormError] = useState('')
  const [formLoading, setFormLoading] = useState(false)

  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [role, setRole] = useState('user')

  const resetForm = () => {
    setLogin('')
    setPassword('')
    setFirstName('')
    setLastName('')
    setRole('user')
    setFormError('')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setFormError('')
    setFormLoading(true)
    try {
      await createUser({ login, password, first_name: firstName, last_name: lastName, role })
      resetForm()
      setShowForm(false)
      await reload()
    } catch (err) {
      setFormError(errorText(err))
    } finally {
      setFormLoading(false)
    }
  }

  if (error) return <p className="card error">{error}</p>

  return (
    <div>
      <div className="page-head">
        <div>
          <h1>Пользователи</h1>
          <p className="muted">Управление пользователями системы</p>
        </div>
        <button className="btn btn-primary" onClick={() => { resetForm(); setShowForm(true) }}>
          Добавить пользователя
        </button>
      </div>

      {showForm && (
        <form className="card" onSubmit={handleSubmit} style={{ marginBottom: 18 }}>
          <h2>Новый пользователь</h2>
          {formError && <div className="card error" style={{ marginBottom: 14, padding: 10 }}>{formError}</div>}

          <div className="form-grid">
            <div className="form-row">
              <label>Имя</label>
              <input type="text" value={firstName} onChange={(e) => setFirstName(e.target.value)} required placeholder="Имя" autoFocus />
            </div>
            <div className="form-row">
              <label>Фамилия</label>
              <input type="text" value={lastName} onChange={(e) => setLastName(e.target.value)} required placeholder="Фамилия" />
            </div>
            <div className="form-row">
              <label>Логин</label>
              <input type="text" value={login} onChange={(e) => setLogin(e.target.value)} required placeholder="Логин для входа" />
            </div>
            <div className="form-row">
              <label>Пароль</label>
              <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required placeholder="Минимум 6 символов" />
            </div>
            <div className="form-row">
              <label>Роль</label>
              <select value={role} onChange={(e) => setRole(e.target.value)}>
                <option value="user">Пользователь</option>
                <option value="admin">Администратор</option>
              </select>
            </div>
          </div>

          <div className="dialog-buttons">
            <button type="button" className="btn" onClick={() => { setShowForm(false); resetForm() }}>Отмена</button>
            <button type="submit" className="btn btn-primary" disabled={formLoading}>
              {formLoading ? 'Создание...' : 'Создать'}
            </button>
          </div>
        </form>
      )}

      <div className="card">
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Логин</th>
                <th>Имя</th>
                <th>Фамилия</th>
                <th>Роль</th>
                <th>Создан</th>
              </tr>
            </thead>
            <tbody>
              {users?.map((u: UserResponse) => (
                <tr key={u.id}>
                  <td>{u.login}</td>
                  <td>{u.first_name}</td>
                  <td>{u.last_name}</td>
                  <td>
                    <span className={`badge ${u.role === 'admin' ? 'badge-ok' : 'badge-muted'}`}>
                      {u.role === 'admin' ? 'Админ' : 'Пользователь'}
                    </span>
                  </td>
                  <td>{new Date(u.created_at).toLocaleDateString('ru-RU')}</td>
                </tr>
              ))}
              {(!users || users.length === 0) && (
                <tr>
                  <td colSpan={5} style={{ textAlign: 'center', padding: 20, color: 'var(--muted)' }}>
                    Пользователей пока нет
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}
