export type Mother = { id: number; first_name: string; last_name: string; phone: string; email: string; photo_url?: string; created_at: string; profile_status: string; mother_status: string; learn_user_id?: string }
export type Child = { id: number; mother_id: number; first_name: string; birth_date: string }
export type SecretInfo = { criminal_record: string; nd_pnd: string; drugs: string; chronic_diseases: string }
export type HelpItem = { id: number; mother_id: number; kind: string; sum: number; check_url: string; created_at: string }
export type ProfileForm = {
  birth_date: string
  birth_place: string
  citizenship: string
  registration: string
  living: string
  email: string
  religion: string
  education: string
  profession: string
  income_source: string
  income_amount: number
  marital_status: string
  father_info: string
  pregnancy_weeks: number
  situation: string
  problem_type: string
  expected_help: string
  source: string
}
export type Profile = ProfileForm & { mother_id: number; status: string; updated_at: string }
export type ProfileSummary = { mother_id: number; status: string; updated_at: string }
export type StatusChange = { status: string; reason: string; created_at: string }
export type MotherStatusInfo = { status: string; history: StatusChange[] }
export type MotherStatusItem = { mother_id: number; status: string }
export type HelpRequest = {
  id: number
  mother_id: number
  kind: string
  description: string
  status: string
  comment: string
  source: string
  help_id: number
  created_at: string
  updated_at: string
}
export type EventType = { id: number; name: string }
export type EventItem = { id: number; type_id: number; type_name: string; title: string; starts_at: string; place: string; max_people: number }
export type Participant = { mother_id: number; status: string }
export type DocumentItem = {
  id: number
  mother_id: number
  type: string
  help_id: number
  status: string
  signed_at: string
  has_file: boolean
  original_name: string
  content_type: string
  size: number
  created_at: string
}
export type DonorBody = { name: string; kind: string; phone: string; email: string; comment: string }
export type Donor = DonorBody & { id: number; created_at: string }
export type DonorSummary = Donor & { total: number; count: number }
export type Donation = { id: number; donor_id: number; amount: number; date: string; comment: string }
export type DonorDetails = { donor: Donor; donations: Donation[] }
export type ReportTotal = { label: string; value: string }
export type Report = { name: string; title: string; headers: string[]; rows: string[][]; totals: ReportTotal[] }
export type CourseItem = {
  course_id: string
  name_course: string
  access_date: string
  due_date?: string
  completion_date?: string
  time_spent: number
  duration: number
  completion_status: string
  progress: number
  awarded_score?: number
  is_overdue: boolean
  views_count: number
}

export type IprTask = {
  task_id: string
  task_type: string
  title: string
  description: string
  status: string
  is_overdue: boolean
  absolute_due_date?: string
  completion_date?: string
}

export type IprStage = {
  stage_id: string
  title: string
  description: string
  start_date?: string
  end_date?: string
  tasks: IprTask[]
}

export type IprPlan = {
  plan_id: string
  title: string
  employee_id: string
  mentor_id: string
  mentor_name: string
  start_date: string
  end_date?: string
  status: string
  progress: number
  completed_tasks: number
  total_tasks: number
  goals: string[]
  stages: IprStage[]
}
