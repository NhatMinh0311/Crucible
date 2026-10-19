import { useTranslation } from 'react-i18next'
import { Card, CardTitle } from '../ui/card'
import { StatusBadge } from '../ui/status-badge'
import { useRoles } from '../auth/roles'

/** Trang chủ tạm thời; M3 thay bằng SCR-10 (sinh viên) và SCR-20 (giảng viên). */
export function Home() {
  const roles = useRoles()
  const { t } = useTranslation()
  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <h1 className="text-2xl font-bold tracking-tight">{t('home.greeting')}</h1>
      <Card>
        <CardTitle>{t('home.shellReady')}</CardTitle>
        <p className="mt-2 text-sm text-muted-foreground">{t('home.roles', { roles: roles.map((r) => t(`roles.${r}`)).join(', ') || t('home.noRole') })}</p>
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
