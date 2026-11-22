import { useState } from 'react'
import { getIPR } from '../../api/mothers'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { fmtDate } from '../../utils/format'
import type { IprPlan } from '../../types'

const statusLabels: Record<string, string> = {
  NOT_STARTED: 'Не начат',
  IN_PROGRESS: 'В процессе',
  WAITING_DECISION: 'Ожидает решения',
  COMPLETED: 'Завершён',
}

const taskStatusIcon = (status: string) => {
  if (status === 'COMPLETED') return '✓'
  if (status === 'IN_PROGRESS') return '◐'
  return '○'
}

function TaskItem({ task }: { task: { task_id: string; title: string; description: string; status: string; is_overdue: boolean; absolute_due_date?: string; completion_date?: string } }) {
  const icon = taskStatusIcon(task.status)
  const isCompleted = task.status === 'COMPLETED'

  return (
    <div className={`ipr-task ${isCompleted ? 'task-completed' : ''} ${task.is_overdue ? 'task-overdue' : ''}`}>
      <span className="task-icon">{icon}</span>
      <div className="task-content">
        <span className="task-title">{task.title}</span>
        {task.description && <span className="task-description">{task.description}</span>}
      </div>
      <div className="task-meta">
        {task.is_overdue && <span className="badge badge-danger">Просрочено</span>}
        {task.absolute_due_date && <span className="muted">до {fmtDate(task.absolute_due_date)}</span>}
        {task.completion_date && <span className="muted">выполнено {fmtDate(task.completion_date)}</span>}
      </div>
    </div>
  )
}

function StageSection({ stage, planProgress }: { stage: { stage_id: string; title: string; description: string; start_date?: string; end_date?: string; tasks: { task_id: string; title: string; description: string; status: string; is_overdue: boolean; absolute_due_date?: string; completion_date?: string }[] }; planProgress: number }) {
  const [isOpen, setIsOpen] = useState(true)
  const completedTasks = stage.tasks.filter((t) => t.status === 'COMPLETED').length
  const totalTasks = stage.tasks.length
  const stageProgress = totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 100) : 0

  return (
    <div className="ipr-stage">
      <button className="ipr-stage-header" onClick={() => setIsOpen(!isOpen)}>
        <span className="stage-arrow">{isOpen ? '▼' : '▶'}</span>
        <div className="stage-title">
          {stage.title}
          {stage.description && <span className="stage-description">{stage.description}</span>}
        </div>
        <div className="stage-progress">
          <span className="muted">{completedTasks}/{totalTasks}</span>
          <span className="badge {stageProgress === 100 ? 'badge-ok' : 'badge-wait'}">{stageProgress}%</span>
        </div>
      </button>
      {isOpen && (
        <div className="ipr-tasks">
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
      <div className="ipr-header">
        <h3>{plan.title}</h3>
        <span className={`badge ${isCompleted ? 'badge-ok' : isInProgress ? 'badge-wait' : 'badge-info'}`}>
          {statusLabel}
        </span>
      </div>

      {plan.mentor_name && (
        <p className="muted">Наставник: {plan.mentor_name}</p>
      )}

      <div className="ipr-dates">
        <span>Период: {fmtDate(plan.start_date)}</span>
        {plan.end_date && <span>— {fmtDate(plan.end_date)}</span>}
      </div>

      <div className="progress-section">
        <span className="progress-label">Прогресс</span>
        <span className="progress-value">{progressPercent}% ({plan.completed_tasks} из {plan.total_tasks} задач)</span>
        <div className="progress-bar">
          <div className="progress-fill" style={{ width: `${progressPercent}%` }} />
        </div>
      </div>

      {plan.goals && plan.goals.length > 0 && (
        <div className="ipr-goals">
          <span className="field-label">Цели</span>
          <ul>
            {plan.goals.map((goal, i) => (
              <li key={i}>{goal}</li>
            ))}
          </ul>
        </div>
      )}

      {plan.stages.map((stage) => (
        <StageSection key={stage.stage_id} stage={stage} planProgress={plan.progress} />
      ))}
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
    <div>
      {active.length > 0 && (
        <div>
          <h2>Активные планы</h2>
          <div className="ipr-grid">
            {active.map((p) => (
              <IPRCard key={p.plan_id} plan={p} />
            ))}
          </div>
        </div>
      )}

      {completed.length > 0 && (
        <div>
          <h2>Завершённые планы</h2>
          <div className="ipr-grid">
            {completed.map((p) => (
              <IPRCard key={p.plan_id} plan={p} />
            ))}
          </div>
        </div>
      )}

      {plans.length === 0 && (
        <div className="card">
          <h2>Индивидуальные планы развития</h2>
          <div className="empty-courses">
            <span className="empty-icon">📋</span>
            <span className="muted">ИПР пока не назначены</span>
          </div>
        </div>
      )}
    </div>
  )
}
