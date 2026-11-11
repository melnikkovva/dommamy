import type { ReactNode } from 'react'

type Props = { label: string; children: ReactNode }

export default function Field({ label, children }: Props) {
  return (
    <div className="form-row">
      <label>{label}</label>
      {children}
    </div>
  )
}
