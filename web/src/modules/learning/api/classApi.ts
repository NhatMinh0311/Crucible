import { learningApi } from '@/shell/api/client'
import type { components } from '@/shell/api/generated/learning'

export type ClassItem = components['schemas']['ClassItem']
export type ClassMember = components['schemas']['ClassMember']
export type ClassRole = components['schemas']['ClassRole']
export type ClassStatus = components['schemas']['ClassStatus']
export type CreateClassRequest = components['schemas']['CreateClassRequest']
export type UpdateClassRequest = components['schemas']['UpdateClassRequest']
export type AddClassStaffRequest = components['schemas']['AddClassStaffRequest']
export type CsvRowError = components['schemas']['CsvRowError']
export type CsvPreviewResponse = components['schemas']['CsvImportPreviewResponse']
export type CsvCommitResponse = components['schemas']['CsvImportCommitResponse']

export async function fetchMyClasses(): Promise<ClassItem[]> {
  const { data, error } = await learningApi.GET('/v1/classes')
  if (error) throw error
  return data ?? []
}

export async function createClass(body: CreateClassRequest): Promise<ClassItem> {
  const { data, error } = await learningApi.POST('/v1/classes', { body })
  if (error || !data) throw error
  return data
}

export async function fetchClassDetail(classId: string): Promise<ClassItem> {
  const { data, error } = await learningApi.GET('/v1/classes/{classId}', {
    params: { path: { classId } },
  })
  if (error || !data) throw error
  return data
}

export async function updateClass(classId: string, body: UpdateClassRequest): Promise<ClassItem> {
  const { data, error } = await learningApi.PATCH('/v1/classes/{classId}', {
    params: { path: { classId } },
    body,
  })
  if (error || !data) throw error
  return data
}

export async function fetchClassMembers(classId: string, role?: ClassRole): Promise<ClassMember[]> {
  const { data, error } = await learningApi.GET('/v1/classes/{classId}/members', {
    params: { path: { classId }, query: role ? { role } : undefined },
  })
  if (error) throw error
  return data ?? []
}

export async function addClassStaff(classId: string, body: AddClassStaffRequest): Promise<ClassMember> {
  const { data, error } = await learningApi.POST('/v1/classes/{classId}/staff', {
    params: { path: { classId } },
    body,
  })
  if (error || !data) throw error
  return data
}

export async function removeClassMember(classId: string, memberId: string): Promise<void> {
  const { error } = await learningApi.DELETE('/v1/classes/{classId}/members/{memberId}', {
    params: { path: { classId, memberId } },
  })
  if (error) throw error
}

export async function previewCsvMembers(classId: string, file: File): Promise<CsvPreviewResponse> {
  const { data, error } = await learningApi.POST('/v1/classes/{classId}/members/csv-preview', {
    params: { path: { classId } },
    body: { file: file.name },
    bodySerializer() {
      const formData = new FormData()
      formData.append('file', file)
      return formData
    },
  })
  if (error || !data) throw error
  return data
}

export async function importCsvMembers(classId: string, file: File): Promise<CsvCommitResponse> {
  const { data, error } = await learningApi.POST('/v1/classes/{classId}/members/csv-import', {
    params: { path: { classId } },
    body: { file: file.name },
    bodySerializer() {
      const formData = new FormData()
      formData.append('file', file)
      return formData
    },
  })
  if (error || !data) throw error
  return data
}
