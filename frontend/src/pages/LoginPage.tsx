import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [loginValue, setLoginValue] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const err = await login(loginValue, password)
      if (err) {
        setError(err)
      } else {
        navigate('/')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-split">
        <div className="login-image">
          <img src="/favicon.png" alt="Дом мамы" />
        </div>
        <div className="login-form-side">
          <div className="login-container">
            <div className="login-header">
              <div className="logo-icon"><img src="/favicon.png" alt="" width="28" height="28" /></div>
              <h1>Дом мамы</h1>
              <p>Вход в систему</p>
            </div>

            <form className="card login-form" onSubmit={handleSubmit}>
              {error && <div className="card error" style={{ marginBottom: 14 }}>{error}</div>}

              <div className="form-row">
                <label>Логин</label>
                <input
                  type="text"
                  value={loginValue}
                  onChange={(e) => setLoginValue(e.target.value)}
                  placeholder="Введите логин"
                  required
                  autoFocus
                />
              </div>

              <div className="form-row">
                <label>Пароль</label>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="Введите пароль"
                  required
                />
              </div>

              <button className="btn btn-primary" type="submit" disabled={loading} style={{ width: '100%', marginTop: 8 }}>
                {loading ? 'Вход...' : 'Войти'}
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  )
}
