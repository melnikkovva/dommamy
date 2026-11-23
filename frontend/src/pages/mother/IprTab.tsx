import { useState } from 'react'
import { getIPR } from '../../api/mothers'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { fmtDate } from '../../utils/format'
import type { IprPlan, IprTask } from '../../types'

const statusLabels: Record<string, string> = {
  NOT_STARTED: 'Не начат',
  IN_PROGRESS: 'В процессе',
  WAITING_DECISION: 'Ожидает решения',
  COMPLETED: 'Завершён',
}

const statusBadgeClass = (status: string) => {
  if (status === 'COMPLETED') return 'task-badge task-badge-done'
  if (status === 'IN_PROGRESS') return 'task-badge task-badge-progress'
  if (status === 'WAITING_DECISION') return 'task-badge task-badge-wait'
  return 'task-badge task-badge-todo'
}

const statusBadgeLabel = (status: string) => statusLabels[status] || status

function tryParseQuillDescription(raw: string): string | null {
  try {
    const obj = JSON.parse(raw)
    if (obj.state && Array.isArray(obj.state)) {
      const texts = obj.state
        .filter((s: any) => s.T && typeof s.T === 'string')
        .map((s: any) => s.T)
      return texts.join(' ') || null
    }
  } catch {
  }
  return null
}

function TaskItem({ task }: { task: IprTask }) {
  const isCompleted = task.status === 'COMPLETED'
  const descText = task.description ? tryParseQuillDescription(task.description) : null

  return (
    <div className={`ipr-task ${isCompleted ? 'task-completed' : ''} ${task.is_overdue ? 'task-overdue' : ''}`}>
      <div className="ipr-task-left">
        <span className={`task-check ${isCompleted ? 'check-done' : task.status === 'IN_PROGRESS' ? 'check-progress' : ''}`}>
          {isCompleted ? '✓' : task.status === 'IN_PROGRESS' ? '◐' : '○'}
        </span>
        <div className="ipr-task-info">
          <span className="task-title">{task.title || 'Без названия'}</span>
          {descText && <span className="task-desc-text">{descText}</span>}
        </div>
      </div>
      <div className="ipr-task-right">
        <span className={statusBadgeClass(task.status)}>{statusBadgeLabel(task.status)}</span>
        {task.is_overdue && <span className="badge badge-danger">Просрочено</span>}
        {task.absolute_due_date && (
          <span className="task-due">до {fmtDate(task.absolute_due_date)}</span>
        )}
        {task.completion_date && (
          <span className="task-done-date">✓ {fmtDate(task.completion_date)}</span>
        )}
      </div>
    </div>
  )
}

function StageSection({ stage }: { stage: { stage_id: string; title: string; start_date?: string; end_date?: string; tasks: IprTask[] } }) {
  const [isOpen, setIsOpen] = useState(true)
  const completedTasks = stage.tasks.filter((t) => t.status === 'COMPLETED').length
  const totalTasks = stage.tasks.length
  const stageProgress = totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 100) : 0

  return (
    <div className="ipr-stage">
      <button className="ipr-stage-header" onClick={() => setIsOpen(!isOpen)}>
        <span className="stage-arrow">{isOpen ? '▾' : '▸'}</span>
        <div className="stage-info">
          <span className="stage-title">{stage.title}</span>
          {(stage.start_date || stage.end_date) && (
            <span className="stage-dates">
              {stage.start_date && fmtDate(stage.start_date)}
              {stage.start_date && stage.end_date && ' — '}
              {stage.end_date && fmtDate(stage.end_date)}
            </span>
          )}
        </div>
        <div className="stage-progress-mini">
          <div className="stage-progress-track">
            <div className="stage-progress-fill" style={{ width: `${stageProgress}%` }} />
          </div>
          <span className="stage-progress-text">
            {completedTasks}/{totalTasks}
          </span>
        </div>
      </button>
      {isOpen && (
        <div className="ipr-tasks">
          {stage.tasks.length === 0 && <p className="muted" style={{ paddingLeft: '28px' }}>Нет задач</p>}
          {stage.tasks.map((t) => (
            <TaskItem key={t.task_id} task={t} />
          ))}
        </div>
      )}
    </div>
  )
}

function IPRCard({ plan }: { plan: IprPlan }) {
  const progressPercent = Math.round(plan.progress)
  const statusLabel = statusLabels[plan.status] || plan.status
  const isCompleted = plan.status === 'COMPLETED'
  const isInProgress = plan.status === 'IN_PROGRESS'

  return (
    <div className="card ipr-card">
      <div className="ipr-card-header">
        <div className="ipr-card-title-row">
          <h3 className="ipr-plan-title">{plan.title}</h3>
          <span className={`ipr-status-badge ${isCompleted ? 'ipr-status-done' : isInProgress ? 'ipr-status-active' : ''}`}>
            {statusLabel}
          </span>
        </div>
        {plan.mentor_name && (
          <div className="ipr-mentor">
            <span className="ipr-mentor-icon"></span> {plan.mentor_name}
          </div>
        )}
      </div>

      <div className="ipr-meta-row">
        <div className="ipr-meta-item">
          <span>
            {fmtDate(plan.start_date)}
            {plan.end_date && ` — ${fmtDate(plan.end_date)}`}
          </span>
        </div>
        <div className="ipr-meta-item">
          <span>{plan.total_tasks} задач</span>
        </div>
      </div>

      <div className="ipr-progress-section">
        <div className="ipr-progress-header">
          <span className="ipr-progress-label">Общий прогресс</span>
          <span className="ipr-progress-percent">{progressPercent}%</span>
        </div>
        <div className="ipr-progress-track">
          <div
            className={`ipr-progress-fill ${isCompleted ? 'ipr-progress-done' : ''}`}
            style={{ width: `${progressPercent}%` }}
          />
        </div>
        <span className="ipr-progress-detail">
          {plan.completed_tasks} из {plan.total_tasks} задач выполнено
        </span>
      </div>

      {plan.goals && plan.goals.length > 0 && (
        <div className="ipr-goals-section">
          <span className="ipr-section-label"> Цели</span>
          <div className="ipr-goals-list">
            {plan.goals.map((goal, i) => (
              <span key={i} className="ipr-goal-chip">{goal}</span>
            ))}
          </div>
        </div>
      )}

      <div className="ipr-stages-section">
        <span className="ipr-section-label">Этапы</span>
        {plan.stages.map((stage) => (
          <StageSection key={stage.stage_id} stage={stage} />
        ))}
      </div>
    </div>
  )
}

export default function IprTab({ motherId }: { motherId: number }) {
  const { data, error } = useLoad(() => getIPR(motherId), [motherId])
  if (!data) return <LoadState error={error} />

  const plans: IprPlan[] = data
  const active = plans.filter((p) => p.status !== 'COMPLETED')
  const completed = plans.filter((p) => p.status === 'COMPLETED')

  return (
    <div className="ipr-page">
      {active.length > 0 && (
        <section className="ipr-section">
          <h2 className="ipr-section-title">
            <span className="section-icon"></span> Активные планы
          </h2>
          <div className="ipr-card-list">
            {active.map((p) => (
              <IPRCard key={p.plan_id} plan={p} />
            ))}
          </div>
        </section>
      )}

      {completed.length > 0 && (
        <section className="ipr-section">
          <h2 className="ipr-section-title">
            <span className="section-icon"></span> Завершённые планы
          </h2>
          <div className="ipr-card-list">
            {completed.map((p) => (
              <IPRCard key={p.plan_id} plan={p} />
            ))}
          </div>
        </section>
      )}

      {plans.length === 0 && (
        <div className="card ipr-empty">
          <h3>Индивидуальные планы развития</h3>
          <p className="muted">ИПР пока не назначены</p>
        </div>
      )}
    </div>
  )
}
