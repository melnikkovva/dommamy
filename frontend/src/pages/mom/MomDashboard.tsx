import { useState } from 'react'
import { useAuth } from '../../context/AuthContext'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { getMothers, getChildren, getProfile, getMotherStatus } from '../../api/mothers'
import { getHelp } from '../../api/help'
import { getMotherRequests } from '../../api/help'
import { getEvents, getMotherEvents } from '../../api/events'
import { getDocuments } from '../../api/documents'
import { getCourses } from '../../api/mothers'
import { kindNames, label, requestStatusNames, sumText } from '../../utils/labels'
import { fmtDate, fmtDateTime, ageYears, yearsText } from '../../utils/format'
import { useToast } from '../../components/Toast'
import { createRequest } from '../../api/help'

type Tab = 'home' | 'help' | 'events' | 'profile'

const navItems: { key: Tab; label: string }[] = [
  { key: 'home', label: 'Главная' },
  { key: 'help', label: 'Помощь' },
  { key: 'events', label: 'Афиша' },
  { key: 'profile', label: 'Профиль' },
]

interface DashboardData {
  mother: any
  children: any[]
  profile: any
  status: any
  helpItems: any[]
  requests: any[]
  events: any[]
  allEvents: any[]
  documents: any[]
  courses: any[]
}

export default function MomDashboard() {
  const { user } = useAuth()
  const motherId = user?.mother_id
  const [tab, setTab] = useState<Tab>('home')
  const [showRequestModal, setShowRequestModal] = useState(false)
  const [requestKind, setRequestKind] = useState('food')
  const [requestDesc, setRequestDesc] = useState('')
  const [selectedRequest, setSelectedRequest] = useState<any>(null)
  const toast = useToast()

  const { data, error, reload } = useLoad(async (): Promise<DashboardData | null> => {
    if (!motherId) return null
    const [mothers, children, profile, status, help, requests, events, allEvents, docs, courses] = await Promise.all([
      getMothers(),
      getChildren(motherId),
      getProfile(motherId),
      getMotherStatus(motherId),
      getHelp(),
      getMotherRequests(motherId),
      getMotherEvents(motherId),
      getEvents(),
      getDocuments(motherId),
      getCourses(motherId),
    ])
    const mother = mothers.find((m: any) => m.id === motherId)
    const myHelp = help.filter((h: any) => h.mother_id === motherId).sort((a: any, b: any) => b.id - a.id)
    return {
      mother: mother || null,
      children,
      profile,
      status,
      helpItems: myHelp,
      requests,
      events,
      allEvents,
      documents: docs,
      courses,
    }
  }, [motherId])

  const handleCreateRequest = async () => {
    if (!motherId || !requestDesc.trim()) return
    try {
      await createRequest(motherId, requestKind, requestDesc.trim())
      setShowRequestModal(false)
      setRequestDesc('')
      toast('Заявка отправлена')
      reload()
    } catch (e: any) {
      toast(e.message || 'Ошибка', true)
    }
  }

  if (!motherId) {
    return (
      <div className="mom-dashboard">
        <div className="mom-content" style={{ textAlign: 'center', padding: '40px 16px' }}>
          <p>Не удалось определить ваш аккаунт. Обратитесь к администратору.</p>
        </div>
      </div>
    )
  }

  if (!data) return <LoadState error={error} />

  const { mother, children, profile, status, helpItems, requests, events, allEvents, documents, courses } = data
  const firstName = mother?.first_name || user?.first_name || ''
  const now = new Date()
  const days = ['воскресенье', 'понедельник', 'вторник', 'среда', 'четверг', 'пятница', 'суббота']
  const greeting = now.getHours() < 12 ? 'Добрый день' : now.getHours() < 18 ? 'Добрый день' : 'Добрый вечер'
  const dateStr = now.getDate() + ' ' + ['января', 'февраля', 'марта', 'апреля', 'мая', 'июня', 'июля', 'августа', 'сентября', 'октября', 'ноября', 'декабря'][now.getMonth()] + ', ' + days[now.getDay()]

  const activeRequests = requests.filter((r: any) => r.status === 'new' || r.status === 'in_work')
  const firstActive = activeRequests[0]

  const upcomingEvents = allEvents
    .filter((e: any) => new Date(e.starts_at) > now)
    .sort((a: any, b: any) => new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime())
    .slice(0, 1)

  const recentHelp = helpItems.slice(0, 1)

  const childrenText = children.map((c: any) => c.first_name + ', ' + yearsText(ageYears(c.birth_date))).join(' и ')

  const requestKinds = [
    { key: 'food', label: 'Продукты' },
    { key: 'clothes', label: 'Вещи' },
    { key: 'furniture', label: 'Мебель' },
    { key: 'stationery', label: 'Канцелярия' },
    { key: 'gifts', label: 'Подарки' },
    { key: 'consultation', label: 'Консультация' },
    { key: 'other', label: 'Другое' },
  ]

  const handleLogout = () => {
    localStorage.removeItem('auth_token')
    localStorage.removeItem('auth_user')
    window.location.href = '/login'
  }

  const navIcons: Record<Tab, React.ReactNode> = {
    home: (
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
        <polyline points="9 22 9 12 15 12 15 22" />
      </svg>
    ),
    help: (
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M17 8l-5-5-5 5" />
        <path d="M12 3v18" />
      </svg>
    ),
    events: (
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
        <line x1="16" y1="2" x2="16" y2="6" />
        <line x1="8" y1="2" x2="8" y2="6" />
        <line x1="3" y1="10" x2="21" y2="10" />
      </svg>
    ),
    profile: (
      <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
        <circle cx="12" cy="7" r="4" />
      </svg>
    ),
  }

  return (
    <div className="app">
      <div className="mom-mobile-header">
        <div className="mom-mobile-header-left">
          <div className="mom-mobile-header-icon">
            <img src="/favicon.png" alt="" />
          </div>
          <div>
            <div className="mom-mobile-header-title">Дом мамы</div>
            <div className="mom-mobile-header-sub">Кабинет мамы</div>
          </div>
        </div>
        <button className="mom-mobile-header-btn" onClick={handleLogout} title="Выйти">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
            <polyline points="16 17 21 12 16 7" />
            <line x1="21" y1="12" x2="9" y2="12" />
          </svg>
        </button>
      </div>

      <aside className="sidebar">
        <div className="logo">
          <div className="logo-icon"><img src="/favicon.png" alt="" width="28" height="28" /></div>
          <div>
            <div className="logo-title">Дом мамы</div>
          </div>
        </div>
        <nav className="nav">
          {navItems.map((item) => (
            <button
              key={item.key}
              className={tab === item.key ? 'active' : ''}
              onClick={() => setTab(item.key)}
              style={{ display: 'block', width: '100%', textAlign: 'left', border: 'none', background: 'none', padding: '11px 14px', marginBottom: '4px', borderRadius: '10px', cursor: 'pointer', fontFamily: 'inherit', fontSize: '15px', color: tab === item.key ? 'var(--green)' : 'var(--text)' }}
            >
              {item.label}
            </button>
          ))}
        </nav>
        <div className="sidebar-footer">
          <div className="user-info">
            <div className="user-name">{mother?.first_name || ''} {mother?.last_name || ''}</div>
            <div className="user-role">Кабинет мамы</div>
          </div>
          <button className="btn btn-small" onClick={handleLogout} style={{ width: '100%', marginTop: 8 }}>
            Выйти
          </button>
        </div>
      </aside>
      <div className="main">
        <main className="content">

          {tab === 'home' && (
            <div>
              <div className="page-head">
                <div>
                  <h1>{greeting}, {firstName}</h1>
                  <p className="muted">{dateStr}. Мы рядом, если нужна поддержка.</p>
                </div>
              </div>

              <div className="card" style={{ background: 'var(--green-light)', border: 'none' }}>
                <p style={{ fontSize: '17px', fontWeight: 600, margin: '0 0 6px', color: 'var(--green)' }}>Что сейчас нужно вашей семье?</p>
                <p className="muted" style={{ margin: '0 0 14px' }}>Расскажите нам. Вместе посмотрим, чем центр может помочь.</p>
                <button className="btn btn-primary" onClick={() => setShowRequestModal(true)}>Запросить помощь</button>
              </div>

              {firstActive && (
                <div className="card">
                  <span className={'badge ' + (firstActive.status === 'in_work' ? 'badge-wait' : 'badge-muted')}>
                    {label(requestStatusNames, firstActive.status)}
                  </span>
                  <p style={{ fontSize: '17px', fontWeight: 600, margin: '8px 0 4px' }}>
                    {firstActive.description || label(kindNames, firstActive.kind)}
                  </p>
                  <p className="muted">Запрос № {String(firstActive.id).padStart(4, '0')}  {fmtDate(firstActive.created_at)}</p>
                  <button className="btn btn-small" onClick={() => setSelectedRequest(firstActive)}>Посмотреть статус</button>
                </div>
              )}

              {upcomingEvents.length > 0 && upcomingEvents.map((ev: any) => (
                <div key={ev.id} className="card">
                  <p className="event-date">Ближайшие мероприятия</p>
                  <p className="event-name">{ev.title}</p>
                  <p className="muted">📍 {ev.place || 'Центр'}</p>
                  <p className="muted">Свободно {ev.max_people > 0 ? ev.max_people : '—'} мест</p>
                </div>
              ))}

              {recentHelp.length > 0 && (
                <div>
                  <div className="row-title">
                    <h2>Последние выдачи</h2>
                    <button className="btn btn-small" onClick={() => setTab('help')}>Вся история</button>
                  </div>
                  {recentHelp.map((h: any) => (
                    <div key={h.id} className="card">
                      <p className="muted">{fmtDate(h.created_at)}</p>
                      <p style={{ fontSize: '17px', fontWeight: 600, margin: '4px 0' }}>
                        {label(kindNames, h.kind)}  {sumText(h.kind, h.sum)}
                      </p>
                      <span className="badge badge-ok">Выдано</span>
                      <p className="muted">Для{childrenText ? ' ' + childrenText : ' вас'}</p>
                    </div>
                  ))}
                </div>
              )}

              {courses.length > 0 && (
                <div className="card">
                  <p style={{ fontSize: '14px', fontWeight: 600, color: 'var(--green)', margin: '0 0 4px' }}>Мои курсы</p>
                  <p className="muted">{courses[0].name_course}  {Math.round(courses[0].progress)}%</p>
                </div>
              )}

              {documents.length > 0 && (
                <div className="card">
                  <p style={{ fontSize: '14px', fontWeight: 600, color: 'var(--green)', margin: '0 0 4px' }}>Мои документы</p>
                  <p className="muted">Анкета, согласие и заявления</p>
                </div>
              )}
            </div>
          )}

          {tab === 'help' && (
            <div>
              <div className="page-head">
                <div>
                  <h1>Моя помощь</h1>
                  <p className="muted">Помощь вашей семье и ваши заявки.</p>
                </div>
                <button className="btn btn-primary" onClick={() => setShowRequestModal(true)}>+ Запросить помощь</button>
              </div>

              <div className="card">
                <h2>Выдачи</h2>
                <p className="muted" style={{ marginBottom: 14 }}>{helpItems.length} {helpItems.length === 1 ? 'выдача' : 'выдач'}</p>
                {helpItems.map((h: any) => (
                  <div key={h.id} className="event-item">
                    <p className="muted">{fmtDate(h.created_at)}</p>
                    <p style={{ fontSize: '16px', fontWeight: 600, margin: '4px 0' }}>
                      {label(kindNames, h.kind)}  {sumText(h.kind, h.sum)}
                    </p>
                    <span className="badge badge-ok">Выдано</span>
                    <p className="muted">Для{childrenText ? ' ' + childrenText : ' семьи'}</p>
                  </div>
                ))}
                {helpItems.length === 0 && <p className="muted">Помощь ещё не выдавали</p>}
              </div>

              {requests.length > 0 && (
                <div className="card">
                  <h2>Заявки</h2>
                  {requests.map((r: any) => (
                    <div key={r.id} className="event-item">
                      <p className="muted">{fmtDate(r.created_at)} </p>
                      <p style={{ fontSize: '16px', fontWeight: 600, margin: '4px 0' }}>
                        {r.description || label(kindNames, r.kind)}
                      </p>
                      <span className={'badge ' + (r.status === 'issued' ? 'badge-ok' : r.status === 'rejected' ? 'badge-error' : r.status === 'in_work' ? 'badge-wait' : 'badge-muted')}>
                        {label(requestStatusNames, r.status)}
                      </span>
                      {r.status === 'rejected' && r.comment && (
                        <div style={{ marginTop: 8, padding: '10px 12px', background: 'var(--red-light)', borderRadius: '8px', fontSize: '13px' }}>
                          <p className="label" style={{ marginBottom: 2 }}>Причина:</p>
                          <p style={{ margin: 0 }}>{r.comment}</p>
                        </div>
                      )}
                      {r.status === 'issued' && r.help_id && (
                        <p className="muted">Помощь выдана (№ {r.help_id})</p>
                      )}
                    </div>
                  ))}
                </div>
              )}

              <div className="hint">
                <strong>Что-то не совпадает?</strong>
                Свяжитесь с центром — поможем уточнить вашу историю помощи.
              </div>
            </div>
          )}

          {tab === 'events' && (
            <div>
              <div className="page-head">
                <div>
                  <h1>Афиша</h1>
                  <p className="muted">Мероприятия, встречи и мастер-классы.</p>
                </div>
              </div>

              {events.length > 0 && (
                <div className="card">
                  <h2>Мои мероприятия</h2>
                  {events.map((ev: any) => (
                    <div key={ev.id} className="event-item">
                      <p className="event-name">{ev.title}</p>
                      {ev.type_name && <p className="muted" style={{fontSize: 12}}>🏷 {ev.type_name}</p>}
                      <p className="muted">{fmtDateTime(ev.starts_at)}</p>
                      <p className="muted">📍 {ev.place || 'Центр'}</p>
                      <span className="badge badge-wait">Записана</span>
                    </div>
                  ))}
                </div>
              )}

              <div className="card">
                <h2>Ближайшие мероприятия</h2>
                {allEvents
                  .filter((e: any) => new Date(e.starts_at) > now)
                  .sort((a: any, b: any) => new Date(a.starts_at).getTime() - new Date(b.starts_at).getTime())
                  .map((ev: any) => (
                    <div key={ev.id} className="event-item">
                      <p className="event-name">{ev.title}</p>
                      {ev.type_name && <p className="muted" style={{fontSize: 12}}>🏷 {ev.type_name}</p>}
                      <p className="muted">{fmtDateTime(ev.starts_at)}</p>
                      <p className="muted">📍 {ev.place || 'Центр'}</p>
                      <p className="muted">Свободно {ev.max_people > 0 ? ev.max_people : '—'} мест из {ev.max_people || '—'}</p>
                    </div>
                  ))
                }
                {allEvents.filter((e: any) => new Date(e.starts_at) > now).length === 0 && (
                  <p className="muted">Ближайших мероприятий нет</p>
                )}
              </div>
            </div>
          )}

          {tab === 'profile' && (
            <div>
              <div className="page-head">
                <div>
                  <h1>Профиль</h1>
                </div>
              </div>

              <div className="card">
                <h2>Личная информация</h2>
                <div style={{ borderBottom: '1px solid var(--border)', padding: '10px 0' }}>
                  <span className="muted" style={{ marginRight: 'auto' }}>Имя</span>
                  <span style={{ float: 'right', fontWeight: 500 }}>{mother?.first_name || '—'}</span>
                </div>
                <div style={{ borderBottom: '1px solid var(--border)', padding: '10px 0' }}>
                  <span className="muted">Фамилия</span>
                  <span style={{ float: 'right', fontWeight: 500 }}>{mother?.last_name || '—'}</span>
                </div>
                <div style={{ borderBottom: '1px solid var(--border)', padding: '10px 0' }}>
                  <span className="muted">Телефон</span>
                  <span style={{ float: 'right', fontWeight: 500 }}>{mother?.phone || '—'}</span>
                </div>
                <div style={{ padding: '10px 0' }}>
                  <span className="muted">Логин</span>
                  <span style={{ float: 'right', fontWeight: 500 }}>{user?.login || '—'}</span>
                </div>
              </div>

              {children.length > 0 && (
                <div className="card">
                  <h2>Дети</h2>
                  {children.map((c: any, i: number) => (
                    <div key={c.id} style={{ borderBottom: i < children.length - 1 ? '1px solid var(--border)' : 'none', padding: '10px 0' }}>
                      <span style={{ fontWeight: 500 }}>{c.first_name}</span>
                      <span className="muted" style={{ marginLeft: 8 }}>{yearsText(ageYears(c.birth_date))}</span>
                    </div>
                  ))}
                </div>
              )}

              {profile && (
                <div className="card">
                  <h2>Анкета</h2>
                  {profile.birth_date && (
                    <div style={{ borderBottom: '1px solid var(--border)', padding: '10px 0' }}>
                      <span className="muted">Дата рождения</span>
                      <span style={{ float: 'right', fontWeight: 500 }}>{fmtDate(profile.birth_date)}</span>
                    </div>
                  )}
                  {profile.education && (
                    <div style={{ padding: '10px 0' }}>
                      <span className="muted">Образование</span>
                      <span style={{ float: 'right', fontWeight: 500 }}>{profile.education}</span>
                    </div>
                  )}
                </div>
              )}

              {status && (
                <div className="card">
                  <h2>Статус</h2>
                  <div style={{ padding: '10px 0' }}>
                    <span className="muted">Сопровождение</span>
                    <span style={{ float: 'right', fontWeight: 500 }}>{label(requestStatusNames, status.status) || status.status}</span>
                  </div>
                </div>
              )}

              {courses.length > 0 && (
                <div className="card">
                  <h2>Мои курсы</h2>
                  {courses.map((c: any) => (
                    <div key={c.course_id} style={{ display: 'flex', justifyContent: 'space-between', padding: '10px 0', borderBottom: '1px solid var(--border)' }}>
                      <span style={{ fontWeight: 500 }}>{c.name_course}</span>
                      <span style={{ fontWeight: 700, color: 'var(--green)' }}>{Math.round(c.progress)}%</span>
                    </div>
                  ))}
                </div>
              )}

              {documents.length > 0 && (
                <div className="card">
                  <h2>Документы</h2>
                  {documents.map((d: any, i: number) => (
                    <div key={d.id} style={{ display: 'flex', justifyContent: 'space-between', padding: '10px 0', borderBottom: i < documents.length - 1 ? '1px solid var(--border)' : 'none' }}>
                      <span>{d.type}</span>
                      <span style={{ color: 'var(--green)' }}>📄</span>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

        </main>
      </div>

      {showRequestModal && (
        <div className="modal-backdrop" onClick={() => setShowRequestModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2 style={{ marginBottom: 18 }}>Запросить помощь</h2>
            <div className="form-row">
              <label>Вид помощи</label>
              <select value={requestKind} onChange={(e) => setRequestKind(e.target.value)}>
                {requestKinds.map((k) => (
                  <option key={k.key} value={k.key}>{k.label}</option>
                ))}
              </select>
            </div>
            <div className="form-row">
              <label>Описание</label>
              <textarea
                placeholder="Опишите, какая помощь нужна вашей семье..."
                value={requestDesc}
                onChange={(e) => setRequestDesc(e.target.value)}
                style={{ minHeight: 100, resize: 'vertical' }}
              />
            </div>
            <div className="dialog-buttons">
              <button className="btn" onClick={() => setShowRequestModal(false)}>Отмена</button>
              <button className="btn btn-primary" onClick={handleCreateRequest}>Отправить заявку</button>
            </div>
          </div>
        </div>
      )}

      {selectedRequest && (
        <div className="modal-backdrop" onClick={() => setSelectedRequest(null)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <h2 style={{ marginBottom: 18 }}>Запрос № {String(selectedRequest.id).padStart(4, '0')}</h2>
            <div style={{ padding: '12px', background: 'var(--bg)', borderRadius: '10px', fontSize: '14px' }}>
              <p className="field-label">Статус</p>
              <p><span className={'badge ' + (selectedRequest.status === 'issued' ? 'badge-ok' : selectedRequest.status === 'rejected' ? 'badge-error' : selectedRequest.status === 'in_work' ? 'badge-wait' : 'badge-muted')}>
                {label(requestStatusNames, selectedRequest.status)}
              </span></p>
              <p className="field-label">Описание</p>
              <p>{selectedRequest.description || label(kindNames, selectedRequest.kind)}</p>
              <p className="field-label">Дата</p>
              <p>{fmtDate(selectedRequest.created_at)}</p>
              {selectedRequest.comment && (
                <>
                  <p className="field-label">Комментарий</p>
                  <p>{selectedRequest.comment}</p>
                </>
              )}
            </div>
            <div className="dialog-buttons">
              <button className="btn" onClick={() => setSelectedRequest(null)}>Закрыть</button>
            </div>
          </div>
        </div>
      )}

      <nav className="mom-bottom-nav">
        {navItems.map((item) => (
          <button
            key={item.key}
            className={'mom-nav-item' + (tab === item.key ? ' active' : '')}
            onClick={() => setTab(item.key)}
          >
            <span className="mom-nav-icon">{navIcons[item.key]}</span>
            <span className="mom-nav-label">{item.label}</span>
          </button>
        ))}
      </nav>
    </div>
  )
}
