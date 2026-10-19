import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'
import { type Role, useRoles } from './roles'

export function RequireRole({ roles, children }: { roles: Role[]; children: ReactNode }) {
  const mine = useRoles()
  const { t } = useTranslation()
  if (!roles.some((r) => mine.includes(r))) {
    return <p className="text-muted-foreground">{t('auth.forbidden')}</p>
  }
  return <>{children}</>
}
