import { NavLink, Outlet } from 'react-router-dom'

const links: [string, string][] = [
  ['/', 'Главная'],
  ['/mothers', 'Подопечные'],
  ['/help', 'Выдать помощь'],
  ['/requests', 'Запросы на помощь'],
  ['/events', 'Мероприятия'],
  ['/donors', 'Дарители'],
  ['/reports', 'Отчёты']
]

export default function Layout() {
  return (
    <div className="app">
      <aside className="sidebar">
        <div className="logo">
          <div className="logo-icon">Т</div>
          <div>
            <div className="logo-title">Дом мамы</div>
          </div>
        </div>
        <nav className="nav">
          {links.map(([to, text]) => (
            <NavLink key={to} to={to} end={to === '/'} className={({ isActive }) => (isActive ? 'active' : '')}>
              {text}
            </NavLink>
          ))}
        </nav>
      </aside>
      <div className="main">
        <main className="content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
