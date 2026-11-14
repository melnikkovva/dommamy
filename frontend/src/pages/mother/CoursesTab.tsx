import { getCourses } from '../../api/mothers'
import { useLoad } from '../../hooks/useLoad'
import LoadState from '../../components/LoadState'
import { fmtDate } from '../../utils/format'
import type { CourseItem } from '../../types'

function CourseCard({ course }: { course: CourseItem }) {
  const isCompleted = course.completion_status === 'completed'
  const progressPercent = Math.round(course.progress)

  const lessonsCompleted = Math.floor((course.views_count / 100) * course.progress)
  const lessonsTotal = course.views_count > 0 ? Math.ceil(course.views_count / 100 * 100 / (course.progress || 1)) : 0

  return (
    <div className="card course-card">
      <div className="course-header">
        <h3>{course.name_course}</h3>
        <span className={`badge ${isCompleted ? 'badge-ok' : 'badge-wait'}`}>
          {isCompleted ? 'Завершён' : 'В процессе'}
        </span>
      </div>

      <p className="muted">Learn · Старт {fmtDate(course.access_date)}</p>

      <div className="progress-section">
        <span className="progress-label">Прогресс обучения</span>
        <span className="progress-value">{progressPercent}%</span>
        <div className="progress-bar">
          <div className="progress-fill" style={{ width: `${progressPercent}%` }} />
        </div>
      </div>

      <div className="course-results">
        <span className="field-label">Результаты</span>
        {isCompleted
          ? <p>{course.views_count} из {course.views_count} уроков · Итоговый тест: {course.awarded_score ?? '—'}</p>
          : <p>Пройдено {lessonsCompleted} из {lessonsTotal || '?'} уроков · Тесты: {course.awarded_score ?? '—'}</p>
        }
      </div>

      <div className="course-dates">
        <span className="field-label">Дата завершения</span>
        {isCompleted
          ? <p>{fmtDate(course.completion_date || '')}</p>
          : <p>Ещё не завершён · План: {fmtDate(course.due_date || '')}</p>
        }
      </div>

      <div className="course-actions">
        {isCompleted
          ? <button className="btn btn-small">📄 Сертификат.pdf</button>
          : <button className="btn btn-small">Обновить сведения</button>
        }
      </div>
    </div>
  )
}

export default function CoursesTab({ motherId }: { motherId: number }) {
  const { data, error } = useLoad(() => getCourses(motherId), [motherId])
  if (!data) return <LoadState error={error} />

  const courses: CourseItem[] = data
  const inProgress = courses.filter((c) => c.completion_status !== 'completed')
  const completed = courses.filter((c) => c.completion_status === 'completed')

  return (
    <div>
      {inProgress.length > 0 && (
        <div className="courses-grid">
          {inProgress.map((c) => (
            <CourseCard key={c.course_id} course={c} />
          ))}
        </div>
      )}

      {completed.length > 0 && (
        <div className="courses-grid">
          {completed.map((c) => (
            <CourseCard key={c.course_id} course={c} />
          ))}
        </div>
      )}

      {courses.length === 0 && (
        <div className="card">
          <h2>Другие курсы</h2>
          <div className="empty-courses">
            <span className="empty-icon">📖</span>
            <span className="muted">Новых курсов пока нет</span>
            <button className="btn">Добавить курс</button>
          </div>
        </div>
      )}
    </div>
  )
}
