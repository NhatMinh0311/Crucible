import { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router'
import { useTranslation } from 'react-i18next'
import { ArrowLeft, Upload, UserPlus, Users, GraduationCap, ShieldCheck } from 'lucide-react'
import { Button } from '@/shell/ui/button'
import { Card, CardTitle } from '@/shell/ui/card'
import { fetchClassDetail, fetchClassMembers, addInstructorOrAssistant, type ClassItem, type ClassMember } from '../api/classApi'
import { CsvImportModal } from '../components/CsvImportModal'

export function ClassDetailPage() {
  const { classId } = useParams<{ classId: string }>()
  const { t } = useTranslation()
  const [clazz, setClazz] = useState<ClassItem | null>(null)
  const [members, setMembers] = useState<ClassMember[]>([])
  const [activeTab, setActiveTab] = useState<'members' | 'roadmap' | 'assignments' | 'exams'>('members')
  const [showCsvModal, setShowCsvModal] = useState(false)
  const [showAddInstructorModal, setShowAddInstructorModal] = useState(false)
  const [instEmail, setInstEmail] = useState('')
  const [instName, setInstName] = useState('')
  const [instRole, setInstRole] = useState<'INSTRUCTOR' | 'TEACHING_ASSISTANT'>('TEACHING_ASSISTANT')
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
      setMembers(membersData)
    } catch {
      // In dev handle gracefully
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [classId])

  const handleAddInstructor = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!classId) return
    try {
      await addInstructorOrAssistant(classId, {
        email: instEmail,
        fullName: instName,
        role: instRole,
      })
      setShowAddInstructorModal(false)
      setInstEmail('')
      setInstName('')
      await loadData()
    } catch (err: unknown) {
      console.error(err)
    }
  }

  if (loading && !clazz) {
    return <div className="py-12 text-center font-mono text-sm text-muted-foreground">Đang tải...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Link to="/classes" className="text-muted-foreground hover:text-foreground">
          <ArrowLeft className="size-5" />
        </Link>
        <div>
          <div className="flex items-center gap-2">
            <h1 className="font-display text-2xl font-bold tracking-tight">{clazz?.name ?? 'Lớp học'}</h1>
            <span className="label-mono px-2 py-0.5 border border-rule text-xs text-primary font-bold">
              {clazz?.code}
            </span>
          </div>
          <p className="text-xs font-mono text-muted-foreground mt-0.5">
            {clazz?.semester} • SCR-12
          </p>
        </div>
      </div>

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
          Lộ trình
        </button>
      </div>

      {activeTab === 'members' && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <h2 className="font-display text-lg font-semibold tracking-tight">Danh sách thành viên lớp</h2>
            <div className="flex gap-2">
              <Button variant="outline" size="sm" onClick={() => setShowAddInstructorModal(true)}>
                <UserPlus className="size-4" />
                {t('learning:classes.addInstructor')}
              </Button>
              <Button variant="primary" size="sm" onClick={() => setShowCsvModal(true)}>
                <Upload className="size-4" />
                {t('learning:classes.importCsv')}
              </Button>
            </div>
          </div>

          <Card className="p-0 overflow-hidden">
            <table className="w-full text-sm text-left">
              <thead className="bg-muted/40 border-b border-rule font-mono text-xs text-muted-foreground">
                <tr>
                  <th className="py-3 px-4">Mã số</th>
                  <th className="py-3 px-4">Họ và tên</th>
                  <th className="py-3 px-4">Email</th>
                  <th className="py-3 px-4">Vai trò</th>
                  <th className="py-3 px-4">Trạng thái</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-rule font-mono text-xs">
                {members.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-8 text-center text-muted-foreground">
                      Chưa có thành viên nào trong lớp.
                    </td>
                  </tr>
                ) : (
                  members.map((m) => (
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
                          {m.role === 'INSTRUCTOR' ? 'Giảng viên' : m.role === 'TEACHING_ASSISTANT' ? 'Trợ giảng' : 'Sinh viên'}
                        </span>
                      </td>
                      <td className="py-3 px-4">
                        {m.userId ? (
                          <span className="text-success font-semibold">Đã liên kết</span>
                        ) : (
                          <span className="text-muted-foreground">Chờ đăng nhập</span>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </Card>
        </div>
      )}

      {activeTab === 'roadmap' && (
        <Card className="text-center py-12">
          <p className="text-sm font-mono text-muted-foreground">
            Lộ trình luyện tập của lớp học sẽ được triển khai trong Issue #41.
          </p>
        </Card>
      )}

      {classId && (
        <CsvImportModal
          classId={classId}
          isOpen={showCsvModal}
          onClose={() => setShowCsvModal(false)}
          onSuccess={loadData}
        />
      )}

      {showAddInstructorModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
          <Card className="w-full max-w-md bg-surface border-rule">
            <CardTitle className="border-b border-rule pb-3">{t('learning:classes.addInstructor')}</CardTitle>
            <form onSubmit={handleAddInstructor} className="mt-4 space-y-4">
              <div>
                <label className="block text-xs font-mono mb-1">Email trường</label>
                <input
                  type="email"
                  required
                  value={instEmail}
                  onChange={(e) => setInstEmail(e.target.value)}
                  placeholder="name@sis.hust.edu.vn"
                  className="w-full h-10 px-3 border border-rule bg-background text-sm font-mono focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-xs font-mono mb-1">Họ và tên</label>
                <input
                  required
                  value={instName}
                  onChange={(e) => setInstName(e.target.value)}
                  placeholder="Thầy Nguyễn Văn A"
                  className="w-full h-10 px-3 border border-rule bg-background text-sm focus:border-primary focus:outline-none"
                />
              </div>
              <div>
                <label className="block text-xs font-mono mb-1">Vai trò</label>
                <select
                  value={instRole}
                  onChange={(e) => setInstRole(e.target.value as 'INSTRUCTOR' | 'TEACHING_ASSISTANT')}
                  className="w-full h-10 px-3 border border-rule bg-background text-sm font-mono focus:border-primary focus:outline-none"
                >
                  <option value="TEACHING_ASSISTANT">Trợ giảng (TEACHING_ASSISTANT)</option>
                  <option value="INSTRUCTOR">Giảng viên (INSTRUCTOR)</option>
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
    </div>
  )
}
