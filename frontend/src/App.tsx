import { Route, Routes } from 'react-router-dom'
import Layout from './components/Layout'
import { ToastProvider } from './components/Toast'
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
import NotFoundPage from './pages/NotFoundPage'

export default function App() {
  return (
    <ToastProvider>
      <Routes>
        <Route element={<Layout />}>
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
          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Routes>
    </ToastProvider>
  )
}
