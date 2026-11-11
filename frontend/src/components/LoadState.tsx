export default function LoadState({ error }: { error: string }) {
  if (error) return <div className="card error">Ошибка: {error}</div>
  return <p className="muted">Загрузка</p>
}
