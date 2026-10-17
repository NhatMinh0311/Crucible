import type { ReactNode } from 'react'
import { type Role, useRoles } from './roles'

export function RequireRole({ roles, children }: { roles: Role[]; children: ReactNode }) {
  const mine = useRoles()
  if (!roles.some((r) => mine.includes(r))) {
    return <p className="text-muted-foreground">Tài khoản không có quyền truy cập trang này.</p>
  }
  return <>{children}</>
}
