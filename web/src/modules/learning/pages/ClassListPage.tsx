import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Plus, Users, BookOpen } from 'lucide-react'
import { useRoles } from '@/shell/auth/roles'
import type { ApiProblem } from '@/shell/i18n/problem'
import { Button } from '@/shell/ui/button'
import { Card, CardTitle } from '@/shell/ui/card'
import { ProblemAlert } from '@/shell/ui/problem-alert'
import { fetchMyClasses, createClass, type ClassItem } from '../api/classApi'

export function ClassListPage() {
  const { t } = useTranslation(['learning', 'shell'])
  const roles = useRoles()
  const canCreateClass = roles.includes('INSTRUCTOR') || roles.includes('ADMIN')
  const [classes, setClasses] = useState<ClassItem[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [newCode, setNewCode] = useState('')
  const [newName, setNewName] = useState('')
  const [newSemester, setNewSemester] = useState('')
  const [weeklyMinutes, setWeeklyMinutes] = useState(1200)
  const [maxConcurrentSessions, setMaxConcurrentSessions] = useState(30)
  const [errorProblem, setErrorProblem] = useState<ApiProblem | null>(null)

  const loadClasses = async () => {
    try {
      setLoading(true)
      const data = await fetchMyClasses()
      setClasses(data)
    } catch {
      // Môi trường dev khi dịch vụ chưa bật
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadClasses()
  }, [])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setErrorProblem(null)
    try {
      await createClass({
        code: newCode,
        name: newName,
        semester: newSemester,
        weeklyMinutes,
        maxConcurrentSessions,
      })
      setShowCreateModal(false)
      setNewCode('')
      setNewName('')
      setNewSemester('')
      setWeeklyMinutes(1200)
      setMaxConcurrentSessions(30)
      await loadClasses()
    } catch (err: unknown) {
      setErrorProblem(typeof err === 'object' && err ? (err as ApiProblem) : {})
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between border-b border-rule pb-4">
        <div>
          <h1 className="font-display text-2xl font-bold tracking-tight">{t('learning:classes.title')}</h1>
        </div>
        {canCreateClass && (
          <Button variant="primary" onClick={() => setShowCreateModal(true)}>
            <Plus className="size-4" />
            {t('learning:classes.create')}
          </Button>
        )}
      </div>

      {loading ? (
        <div className="text-center py-12 text-sm font-mono text-muted-foreground">{t('learning:common.loading')}</div>
      ) : classes.length === 0 ? (
        <Card className="text-center py-12">
          <BookOpen className="size-10 mx-auto text-muted-foreground mb-3" />
          <p className="text-sm text-muted-foreground">{t('learning:classes.empty')}</p>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {classes.map((c) => (
            <Link key={c.id} to={`/classes/${c.id}`} className="group">
              <Card className="h-full hover:border-primary transition-colors flex flex-col justify-between">
                <div>
                  <div className="flex items-center justify-between text-xs font-mono text-muted-foreground mb-2">
                    <span className="font-bold text-primary">{c.code}</span>
                    <span>{c.semester}</span>
                  </div>
                  <CardTitle className="group-hover:text-primary transition-colors">{c.name}</CardTitle>
                </div>
                <div className="mt-6 flex items-center justify-between border-t border-rule/60 pt-3 text-xs font-mono text-muted-foreground">
                  <span className="flex items-center gap-1.5">
                    <Users className="size-3.5" />
                    {c.memberCount} {t('learning:classes.members')}
                  </span>
                  <span className={c.status === 'ACTIVE' ? 'text-success font-semibold' : 'text-muted-foreground'}>
                    {c.status === 'ACTIVE' ? t('learning:classes.status.active') : t('learning:classes.status.archived')}
                  </span>
                </div>
              </Card>
            </Link>
          ))}
        </div>
      )}

      {showCreateModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
          <Card className="w-full max-w-md bg-surface border-rule">
            <CardTitle className="border-b border-rule pb-3">{t('learning:classes.createTitle')}</CardTitle>
            <form onSubmit={handleCreate} className="mt-4 space-y-4">
              {errorProblem && <ProblemAlert problem={errorProblem} ns="learning" />}
              <div>
                <label className="block text-xs font-mono mb-1">{t('learning:classes.code')}</label>
                <input
                  required
                  value={newCode}
                  onChange={(e) => setNewCode(e.target.value)}
                  placeholder={t('learning:classes.placeholders.code')}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm font-mono focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-xs font-mono mb-1">{t('learning:classes.name')}</label>
                <input
                  required
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  placeholder={t('learning:classes.placeholders.name')}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-xs font-mono mb-1">{t('learning:classes.semester')}</label>
                <input
                  required
                  value={newSemester}
                  onChange={(e) => setNewSemester(e.target.value)}
                  placeholder={t('learning:classes.placeholders.semester')}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm font-mono focus:border-primary focus:outline-none"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
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
              </div>
              <div className="flex justify-end gap-2 pt-2 border-t border-rule">
                <Button type="button" variant="outline" onClick={() => setShowCreateModal(false)}>
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
