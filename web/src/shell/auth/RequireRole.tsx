import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import { type Role, useRoles } from './roles'

/** Ẩn một phần của màn hình theo vai trò; chặn cả route thì khai báo `handle: { roles }` trên route (RoleGate). */
export function RequireRole({ roles, children }: { roles: Role[]; children: ReactNode }) {
  const mine = useRoles()
  const { t } = useTranslation()
  if (!roles.some((r) => mine.includes(r))) {
    return <p className="border border-rule bg-surface p-6 text-muted-foreground">{t('auth.forbidden')}</p>
  }
  return <>{children}</>
}
