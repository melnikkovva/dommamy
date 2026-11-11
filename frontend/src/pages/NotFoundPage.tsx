import { Link } from 'react-router-dom'

export default function NotFoundPage() {
  return (
    <div className="card">
      Такой страницы нет. <Link to="/">На главную</Link>
    </div>
  )
}
