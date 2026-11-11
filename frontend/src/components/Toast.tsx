import { createContext, useCallback, useContext, useState } from 'react'
import type { ReactNode } from 'react'

type ToastItem = { id: number; text: string; isError: boolean }
type ShowToast = (text: string, isError?: boolean) => void

const ToastContext = createContext<ShowToast>(() => {})

export function useToast() {
  return useContext(ToastContext)
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<ToastItem[]>([])
  const show = useCallback((text: string, isError = false) => {
    const id = Date.now() + Math.random()
    setItems((list) => [...list, { id, text, isError }])
    setTimeout(() => setItems((list) => list.filter((item) => item.id !== id)), 3500)
  }, [])
  return (
    <ToastContext.Provider value={show}>
      {children}
      {items.map((item, index) => (
        <div key={item.id} className={item.isError ? 'toast error' : 'toast'} style={{ bottom: 24 + index * 64 }}>
          {item.text}
        </div>
      ))}
    </ToastContext.Provider>
  )
}
