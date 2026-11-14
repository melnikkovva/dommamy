import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { login as apiLogin, clearAuth, setAuth, getUser, AuthUser } from '../api/auth'
import { errorText } from '../api/client'

interface AuthContextType {
  user: AuthUser | null
  login: (login: string, password: string) => Promise<string>
  logout: () => void
  isAdmin: boolean
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  login: async () => { throw new Error('not implemented') },
  logout: () => {},
  isAdmin: false,
})

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(getUser())

  useEffect(() => {
    if (user) {
      setAuth('', user)
    }
  }, [user])

  const doLogin = async (loginStr: string, password: string): Promise<string> => {
    try {
      const res = await apiLogin(loginStr, password)
      const authUser: AuthUser = {
        id: res.id,
        login: res.login,
        first_name: res.first_name,
        last_name: res.last_name,
        role: res.role,
      }
      setAuth(res.token, authUser)
      setUser(authUser)
      return ''
    } catch (e) {
      return errorText(e)
    }
  }

  const doLogout = () => {
    clearAuth()
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, login: doLogin, logout: doLogout, isAdmin: user?.role === 'admin' }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
