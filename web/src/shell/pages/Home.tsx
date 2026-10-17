import { Card, CardTitle } from '../ui/card'
import { StatusBadge } from '../ui/status-badge'
import { useRoles } from '../auth/roles'

/** Trang chủ tạm thời; M3 thay bằng SCR-10 (sinh viên) và SCR-20 (giảng viên). */
export function Home() {
  const roles = useRoles()
  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <h1 className="text-2xl font-bold tracking-tight">Xin chào</h1>
      <Card>
        <CardTitle>Khung giao diện đang hoạt động</CardTitle>
        <p className="mt-2 text-sm text-muted-foreground">Vai trò: {roles.join(', ') || 'chưa có'}</p>
        <div className="mt-4 flex flex-wrap gap-2">
          <StatusBadge status="passed" />
          <StatusBadge status="failed" />
          <StatusBadge status="checking" />
          <StatusBadge status="recorded" />
          <StatusBadge status="interrupted" />
          <StatusBadge status="pending" />
        </div>
      </Card>
    </div>
  )
}
