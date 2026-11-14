import { Route, Routes, Navigate, NavigateOptions } from 'react-router-dom'
import Layout from './components/Layout'
import { ToastProvider } from './components/Toast'
import { AuthProvider, useAuth } from './context/AuthContext'
import HomePage from './pages/HomePage'
import MothersPage from './pages/MothersPage'
import MotherPage from './pages/MotherPage'
import HelpPage from './pages/HelpPage'
import RequestsPage from './pages/RequestsPage'
import EventsPage from './pages/EventsPage'
import EventPage from './pages/EventPage'
import DonorsPage from './pages/DonorsPage'
import DonorPage from './pages/DonorPage'
import ReportsPage from './pages/ReportsPage'
import UsersPage from './pages/UsersPage'
import LoginPage from './pages/LoginPage'
import NotFoundPage from './pages/NotFoundPage'

function ProtectedRoute({ children }: { children: JSX.Element }) {
  const { user } = useAuth()
  if (!user) {
    const to: NavigateOptions = { replace: true }
    return <Navigate to="/login" {...to} />
  }
  return children
}

function AdminRoute({ children }: { children: JSX.Element }) {
  const { isAdmin } = useAuth()
  if (!isAdmin) {
    const to: NavigateOptions = { replace: true }
    return <Navigate to="/" {...to} />
  }
  return children
}

export default function App() {
  return (
    <AuthProvider>
      <ToastProvider>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route element={<ProtectedRoute><Layout /></ProtectedRoute>}>
            <Route path="/" element={<HomePage />} />
            <Route path="/mothers" element={<MothersPage />} />
            <Route path="/mothers/:id" element={<MotherPage />} />
            <Route path="/help" element={<HelpPage />} />
            <Route path="/requests" element={<RequestsPage />} />
            <Route path="/events" element={<EventsPage />} />
            <Route path="/events/:id" element={<EventPage />} />
            <Route path="/donors" element={<DonorsPage />} />
            <Route path="/donors/:id" element={<DonorPage />} />
            <Route path="/reports" element={<ReportsPage />} />
            <Route path="/users" element={<AdminRoute><UsersPage /></AdminRoute>} />
            <Route path="*" element={<NotFoundPage />} />
          </Route>
        </Routes>
      </ToastProvider>
    </AuthProvider>
  )
}
