import { User } from 'oidc-client-ts'
import { OIDC_AUTHORITY, OIDC_CLIENT_ID } from '@/shell/auth/oidc'

function getAuthHeader(): Record<string, string> {
  const raw = window.sessionStorage.getItem(`oidc.user:${OIDC_AUTHORITY}:${OIDC_CLIENT_ID}`)
  const user = raw ? User.fromStorageString(raw) : null
  if (user?.access_token) {
    return { Authorization: `Bearer ${user.access_token}` }
  }
  return {}
}

export interface ClassItem {
  id: string
  code: string
  name: string
  semester: string
  instructorQuotaAccountRef?: string
  sessionLimitPerStudent: number
  status: 'ACTIVE' | 'ARCHIVED'
  memberCount: number
  createdAt: string
  updatedAt: string
}

export interface ClassMember {
  id: string
  classId: string
  userId?: string
  role: 'INSTRUCTOR' | 'TEACHING_ASSISTANT' | 'STUDENT'
  studentIdNumber?: string
  fullName: string
  email: string
  createdAt: string
}

export interface CsvRowError {
  line: number
  field: string
  message: string
}

export interface CsvPreviewResponse {
  totalRows: number
  validRows: number
  errorRows: number
  validStudents: Array<{
    line: number
    studentId: string
    fullName: string
    email: string
  }>
  errors: CsvRowError[]
}

const BASE_URL = '/api/learning/api/v1/classes'

export async function fetchMyClasses(): Promise<ClassItem[]> {
  const res = await fetch(BASE_URL, {
    headers: { ...getAuthHeader(), 'Content-Type': 'application/json' },
  })
  if (!res.ok) throw await res.json()
  return res.json()
}

export async function createClass(data: {
  code: string
  name: string
  semester: string
  instructorQuotaAccountRef?: string
  sessionLimitPerStudent?: number
}): Promise<ClassItem> {
  const res = await fetch(BASE_URL, {
    method: 'POST',
    headers: { ...getAuthHeader(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw await res.json()
  return res.json()
}

export async function fetchClassDetail(classId: string): Promise<ClassItem> {
  const res = await fetch(`${BASE_URL}/${classId}`, {
    headers: { ...getAuthHeader(), 'Content-Type': 'application/json' },
  })
  if (!res.ok) throw await res.json()
  return res.json()
}

export async function fetchClassMembers(classId: string, role?: string): Promise<ClassMember[]> {
  const url = role ? `${BASE_URL}/${classId}/members?role=${role}` : `${BASE_URL}/${classId}/members`
  const res = await fetch(url, {
    headers: { ...getAuthHeader(), 'Content-Type': 'application/json' },
  })
  if (!res.ok) throw await res.json()
  return res.json()
}

export async function addInstructorOrAssistant(classId: string, data: {
  email: string
  fullName: string
  role: 'INSTRUCTOR' | 'TEACHING_ASSISTANT'
}): Promise<ClassMember> {
  const res = await fetch(`${BASE_URL}/${classId}/instructors`, {
    method: 'POST',
    headers: { ...getAuthHeader(), 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!res.ok) throw await res.json()
  return res.json()
}

export async function importCsv(classId: string, file: File, commit = false): Promise<CsvPreviewResponse | { success: boolean; importedCount: number; message: string }> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await fetch(`${BASE_URL}/${classId}/members/import-csv?commit=${commit}`, {
    method: 'POST',
    headers: { ...getAuthHeader() },
    body: formData,
  })
  if (!res.ok) throw await res.json()
  return res.json()
}
