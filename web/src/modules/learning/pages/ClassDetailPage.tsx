import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { ArrowLeft, Upload, UserPlus, Users, GraduationCap, ShieldCheck, Settings, Trash2, CheckCircle } from 'lucide-react'
import { useRoles } from '@/shell/auth/roles'
import type { ApiProblem } from '@/shell/i18n/problem'
import { Button } from '@/shell/ui/button'
import { Card, CardTitle } from '@/shell/ui/card'
import { ProblemAlert } from '@/shell/ui/problem-alert'
import {
  fetchClassDetail,
  fetchClassMembers,
  addClassStaff,
  updateClass,
  removeClassMember,
  type ClassItem,
  type ClassMember,
} from '../api/classApi'
import { CsvImportModal } from '../components/CsvImportModal'

export function ClassDetailPage() {
  const { classId } = useParams<{ classId: string }>()
  const { t } = useTranslation(['learning', 'shell'])
  const roles = useRoles()
  const canManage = roles.includes('INSTRUCTOR') || roles.includes('ADMIN')
  const [clazz, setClazz] = useState<ClassItem | null>(null)
  const [members, setMembers] = useState<ClassMember[]>([])
  const [activeTab, setActiveTab] = useState<'members' | 'roadmap'>('members')
  const [showCsvModal, setShowCsvModal] = useState(false)
  const [showAddInstructorModal, setShowAddInstructorModal] = useState(false)
  const [showEditQuotaModal, setShowEditQuotaModal] = useState(false)
  const [instEmail, setInstEmail] = useState('')
  const [instName, setInstName] = useState('')
  const [instRole, setInstRole] = useState<'INSTRUCTOR' | 'TEACHING_ASSISTANT'>('TEACHING_ASSISTANT')
  const [weeklyMinutes, setWeeklyMinutes] = useState(1200)
  const [maxConcurrentSessions, setMaxConcurrentSessions] = useState(30)
  const [errorProblem, setErrorProblem] = useState<ApiProblem | null>(null)
  const [importedBannerCount, setImportedBannerCount] = useState<number | null>(null)
  const [loading, setLoading] = useState(true)

  const loadData = async () => {
    if (!classId) return
    try {
      setLoading(true)
      const [classData, membersData] = await Promise.all([
        fetchClassDetail(classId),
        fetchClassMembers(classId),
      ])
      setClazz(classData)
      setWeeklyMinutes(classData.weeklyMinutes)
      setMaxConcurrentSessions(classData.maxConcurrentSessions)
      setMembers(membersData)
    } catch {
      // Môi trường dev khi dịch vụ chưa bật
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadData()
  }, [classId])

  const handleAddInstructor = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!classId) return
    setErrorProblem(null)
    try {
      await addClassStaff(classId, {
        email: instEmail,
        fullName: instName,
        role: instRole,
      })
      setShowAddInstructorModal(false)
      setInstEmail('')
      setInstName('')
      await loadData()
    } catch (err: unknown) {
      setErrorProblem(typeof err === 'object' && err ? (err as ApiProblem) : {})
    }
  }

  const handleUpdateQuota = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!classId) return
    setErrorProblem(null)
    try {
      await updateClass(classId, {
        weeklyMinutes,
        maxConcurrentSessions,
      })
      setShowEditQuotaModal(false)
      await loadData()
    } catch (err: unknown) {
      setErrorProblem(typeof err === 'object' && err ? (err as ApiProblem) : {})
    }
  }

  const handleSetLeadInstructor = async (userId: string) => {
    if (!classId) return
    setErrorProblem(null)
    try {
      await updateClass(classId, { leadInstructorId: userId })
      await loadData()
    } catch (err: unknown) {
      setErrorProblem(typeof err === 'object' && err ? (err as ApiProblem) : {})
    }
  }

  const handleRemoveMember = async (memberId: string) => {
    if (!classId) return
    setErrorProblem(null)
    try {
      await removeClassMember(classId, memberId)
      await loadData()
    } catch (err: unknown) {
      setErrorProblem(typeof err === 'object' && err ? (err as ApiProblem) : {})
    }
  }

  if (loading && !clazz) {
    return <div className="py-12 text-center font-mono text-sm text-muted-foreground">{t('learning:common.loading')}</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <Link to="/classes" className="text-muted-foreground hover:text-foreground">
            <ArrowLeft className="size-5" />
          </Link>
          <div>
            <div className="flex items-center gap-2">
              <h1 className="font-display text-2xl font-bold tracking-tight">{clazz?.name ?? t('learning:classes.title')}</h1>
              <span className="label-mono px-2 py-0.5 border border-rule text-xs text-primary font-bold">
                {clazz?.code}
              </span>
            </div>
            <p className="text-xs font-mono text-muted-foreground mt-0.5">
              {clazz?.semester} · {t('learning:classes.quotaRef')}: {clazz?.quotaAccountRef} · {t('learning:classes.weeklyMinutes')}: {clazz?.weeklyMinutes} · {t('learning:classes.maxConcurrentSessions')}: {clazz?.maxConcurrentSessions}
            </p>
          </div>
        </div>
        {canManage && (
          <Button variant="outline" size="sm" onClick={() => { setErrorProblem(null); setShowEditQuotaModal(true) }}>
            <Settings className="size-4" />
            {t('learning:classes.editQuota')}
          </Button>
        )}
      </div>

      {importedBannerCount !== null && (
        <div className="border border-success bg-surface px-4 py-3 text-sm font-mono text-success flex items-center gap-2">
          <CheckCircle className="size-4 shrink-0" />
          <span>{t('learning:csv.success', { count: importedBannerCount })}</span>
        </div>
      )}

      {errorProblem && !showAddInstructorModal && !showEditQuotaModal && (
        <ProblemAlert problem={errorProblem} ns="learning" />
      )}

      <div className="flex border-b border-rule gap-6 text-sm font-mono">
        <button
          onClick={() => setActiveTab('members')}
          className={`pb-2.5 font-bold border-b-2 -mb-px flex items-center gap-2 transition-colors ${
            activeTab === 'members' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Users className="size-4" />
          {t('learning:classes.members')} ({members.length})
        </button>
        <button
          onClick={() => setActiveTab('roadmap')}
          className={`pb-2.5 font-bold border-b-2 -mb-px flex items-center gap-2 transition-colors ${
            activeTab === 'roadmap' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <GraduationCap className="size-4" />
          {t('learning:classes.roadmap')}
        </button>
      </div>

      {activeTab === 'members' && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="font-display text-lg font-semibold tracking-tight">{t('learning:classes.membersList')}</h2>
            {canManage && (
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={() => { setErrorProblem(null); setShowAddInstructorModal(true) }}>
                  <UserPlus className="size-4" />
                  {t('learning:classes.addInstructor')}
                </Button>
                <Button variant="primary" size="sm" onClick={() => setShowCsvModal(true)}>
                  <Upload className="size-4" />
                  {t('learning:classes.importCsv')}
                </Button>
              </div>
            )}
          </div>

          <Card className="p-0 overflow-hidden">
            <table className="w-full text-sm text-left">
              <thead className="bg-muted/40 border-b border-rule font-mono text-xs text-muted-foreground">
                <tr>
                  <th className="py-3 px-4">{t('learning:common.studentId')}</th>
                  <th className="py-3 px-4">{t('learning:common.fullName')}</th>
                  <th className="py-3 px-4">{t('learning:common.email')}</th>
                  <th className="py-3 px-4">{t('learning:common.role')}</th>
                  <th className="py-3 px-4">{t('learning:common.status')}</th>
                  <th className="py-3 px-4 text-right">{t('learning:common.actions')}</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-rule font-mono text-xs">
                {members.length === 0 ? (
                  <tr>
                    <td colSpan={6} className="py-8 text-center text-muted-foreground">
                      {t('learning:classes.noMembers')}
                    </td>
                  </tr>
                ) : (
                  members.map((m) => {
                    const isLead = Boolean(m.userId && clazz?.leadInstructorId === m.userId)
                    return (
                      <tr key={m.id} className="hover:bg-accent/30">
                        <td className="py-3 px-4 font-bold text-foreground">{m.studentIdNumber || '—'}</td>
                        <td className="py-3 px-4 font-sans text-sm">{m.fullName}</td>
                        <td className="py-3 px-4 text-muted-foreground">{m.email}</td>
                        <td className="py-3 px-4">
                          <span className={`inline-flex items-center gap-1 px-2 py-0.5 border text-[11px] font-bold ${
                            m.role === 'INSTRUCTOR'
                              ? 'border-primary bg-primary/10 text-primary'
                              : m.role === 'TEACHING_ASSISTANT'
                              ? 'border-signal bg-signal-soft/30 text-signal'
                              : 'border-rule text-foreground'
                          }`}>
                            {m.role === 'INSTRUCTOR' && <ShieldCheck className="size-3" />}
                            {m.role === 'INSTRUCTOR'
                              ? t('learning:classes.instructor')
                              : m.role === 'TEACHING_ASSISTANT'
                              ? t('learning:classes.assistant')
                              : t('learning:classes.student')}
                            {isLead && ` (${t('learning:classes.leadInstructor')})`}
                          </span>
                        </td>
                        <td className="py-3 px-4">
                          {m.userId ? (
                            <span className="text-success font-semibold">{t('learning:common.linked')}</span>
                          ) : (
                            <span className="text-muted-foreground">{t('learning:common.pendingLogin')}</span>
                          )}
                        </td>
                        <td className="py-3 px-4 text-right space-x-2">
                          {canManage && m.role === 'INSTRUCTOR' && m.userId && !isLead && (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => void handleSetLeadInstructor(m.userId!)}
                            >
                              {t('learning:classes.setLeadInstructor')}
                            </Button>
                          )}
                          {canManage && !isLead && (
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => void handleRemoveMember(m.id)}
                            >
                              <Trash2 className="size-3.5" />
                              {t('learning:classes.removeMember')}
                            </Button>
                          )}
                        </td>
                      </tr>
                    )
                  })
                )}
              </tbody>
            </table>
          </Card>
        </div>
      )}

      {activeTab === 'roadmap' && (
        <Card className="text-center py-12">
          <p className="text-sm font-mono text-muted-foreground">
            {t('learning:classes.roadmapPlaceholder')}
          </p>
        </Card>
      )}

      {classId && (
        <CsvImportModal
          classId={classId}
          isOpen={showCsvModal}
          onClose={() => setShowCsvModal(false)}
          onSuccess={(importedCount) => {
            setImportedBannerCount(importedCount)
            void loadData()
          }}
        />
      )}

      {showAddInstructorModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
          <Card className="w-full max-w-md bg-surface border-rule">
            <CardTitle className="border-b border-rule pb-3">{t('learning:classes.addInstructor')}</CardTitle>
            <form onSubmit={handleAddInstructor} className="mt-4 space-y-4">
              {errorProblem && <ProblemAlert problem={errorProblem} ns="learning" />}
              <div>
                <label className="block text-xs font-mono mb-1">{t('learning:common.schoolEmail')}</label>
                <input
                  type="email"
                  required
                  value={instEmail}
                  onChange={(e) => setInstEmail(e.target.value)}
                  placeholder={t('learning:classes.placeholders.schoolEmail')}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm font-mono focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-xs font-mono mb-1">{t('learning:common.fullName')}</label>
                <input
                  required
                  value={instName}
                  onChange={(e) => setInstName(e.target.value)}
                  placeholder={t('learning:classes.placeholders.instructorName')}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-xs font-mono mb-1">{t('learning:common.role')}</label>
                <select
                  value={instRole}
                  onChange={(e) => setInstRole(e.target.value as 'INSTRUCTOR' | 'TEACHING_ASSISTANT')}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm font-mono focus:border-primary focus:outline-none"
                >
                  <option value="TEACHING_ASSISTANT">{t('learning:classes.roleAssistant')}</option>
                  <option value="INSTRUCTOR">{t('learning:classes.roleInstructor')}</option>
                </select>
              </div>
              <div className="flex justify-end gap-2 pt-2 border-t border-rule">
                <Button type="button" variant="outline" onClick={() => setShowAddInstructorModal(false)}>
                  {t('learning:classes.cancel')}
                </Button>
                <Button type="submit" variant="primary">
                  {t('learning:classes.save')}
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}

      {showEditQuotaModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
          <Card className="w-full max-w-md bg-surface border-rule">
            <CardTitle className="border-b border-rule pb-3">{t('learning:classes.editQuota')}</CardTitle>
            <form onSubmit={handleUpdateQuota} className="mt-4 space-y-4">
              {errorProblem && <ProblemAlert problem={errorProblem} ns="learning" />}
              <div>
                <label className="block text-xs font-mono mb-1">{t('learning:classes.weeklyMinutes')}</label>
                <input
                  type="number"
                  min={0}
                  max={100000}
                  required
                  value={weeklyMinutes}
                  onChange={(e) => setWeeklyMinutes(Number(e.target.value))}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm font-mono focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-xs font-mono mb-1">{t('learning:classes.maxConcurrentSessions')}</label>
                <input
                  type="number"
                  min={0}
                  max={500}
                  required
                  value={maxConcurrentSessions}
                  onChange={(e) => setMaxConcurrentSessions(Number(e.target.value))}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm font-mono focus:border-primary focus:outline-none"
                />
              </div>
              <div className="flex justify-end gap-2 pt-2 border-t border-rule">
                <Button type="button" variant="outline" onClick={() => setShowEditQuotaModal(false)}>
                  {t('learning:classes.cancel')}
                </Button>
                <Button type="submit" variant="primary">
                  {t('learning:classes.save')}
                </Button>
              </div>
            </form>
          </Card>
        </div>
      )}
    </div>
  )
}
