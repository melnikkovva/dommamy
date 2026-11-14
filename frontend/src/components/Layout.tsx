import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

const links: [string, string][] = [
  ['/', 'Главная'],
  ['/mothers', 'Подопечные'],
  ['/help', 'Выдать помощь'],
  ['/requests', 'Запросы на помощь'],
  ['/events', 'Мероприятия'],
  ['/donors', 'Дарители'],
  ['/reports', 'Отчёты']
]

const adminLinks: [string, string][] = [
  ['/users', 'Пользователи']
]

export default function Layout() {
  const { user, logout, isAdmin } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const allLinks = isAdmin ? [...links, ...adminLinks] : links

  return (
    <div className="app">
      <aside className="sidebar">
        <div className="logo">
          <div className="logo-icon"><img src="/favicon.png" alt="" width="28" height="28" /></div>
          <div>
            <div className="logo-title">Дом мамы</div>
          </div>
        </div>
        <nav className="nav">
          {allLinks.map(([to, text]) => (
            <NavLink key={to} to={to} end={to === '/'} className={({ isActive }) => (isActive ? 'active' : '')}>
              {text}
            </NavLink>
          ))}
        </nav>

        {user && (
          <div className="sidebar-footer">
            <div className="user-info">
              <div className="user-name">{user.first_name} {user.last_name}</div>
              <div className="user-role">{user.role === 'admin' ? 'Администратор' : 'Пользователь'}</div>
            </div>
            <button className="btn btn-small" onClick={handleLogout} style={{ width: '100%', marginTop: 8 }}>
              Выйти
            </button>
          </div>
        )}
      </aside>
      <div className="main">
        <main className="content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
