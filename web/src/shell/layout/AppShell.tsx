import { Bell, FlaskConical, Menu, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { NavLink, Outlet } from 'react-router'
import { learningNav } from '@/modules/learning/routes'
import { contentNav } from '@/modules/content/routes'
import { runtimeNav } from '@/modules/runtime/routes'
import { useRoles } from '../auth/roles'
import { Button } from '../ui/button'
import { cn } from '../ui/cn'
import { AccountMenu } from './AccountMenu'

const NAV = [...learningNav, ...contentNav, ...runtimeNav]

/** Khung chung (P02): menu theo vai trò, chuông thông báo, menu tài khoản (W8). */
export function AppShell() {
  const roles = useRoles()
  const { t } = useTranslation(['shell', 'runtime', 'content', 'learning'])
  const [open, setOpen] = useState(false)
  const [notice, setNotice] = useState<string | null>(null)
  const items = NAV.filter((i) => i.roles.some((r) => roles.includes(r)))

  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && setOpen(false)
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open])

  return (
    <div className="flex min-h-dvh">
      <a
        href="#main"
        className="sr-only z-50 rounded-lg bg-primary px-4 py-2 text-primary-foreground focus:not-sr-only focus:fixed focus:top-2 focus:left-2"
      >
        {t('nav.skipToContent')}
      </a>
      {open && <div className="fixed inset-0 z-20 bg-black/40 md:hidden" onClick={() => setOpen(false)} aria-hidden data-testid="nav-overlay" />}
      <aside
        id="main-nav"
        className={cn(
          'fixed inset-y-0 left-0 z-30 w-64 -translate-x-full bg-sidebar text-sidebar-foreground transition-transform md:static md:translate-x-0',
          open && 'translate-x-0',
        )}
      >
        <div className="flex h-16 items-center gap-2 px-5 text-lg font-bold tracking-tight">
          <FlaskConical className="size-6 text-primary" aria-hidden />
          Crucible
        </div>
        <nav className="flex flex-col gap-1 px-3" aria-label={t('nav.main')}>
          {items.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              onClick={() => setOpen(false)}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 rounded-lg px-3 py-2 text-sm hover:bg-sidebar-active focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
                  isActive && 'bg-sidebar-active font-semibold text-white',
                )
              }
            >
              <item.icon className="size-4" aria-hidden />
              {t(item.labelKey)}
            </NavLink>
          ))}
        </nav>
      </aside>

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-16 items-center gap-2 border-b border-border bg-surface px-4">
          <Button
            variant="ghost"
            size="sm"
            className="md:hidden"
            onClick={() => setOpen(!open)}
            aria-label={t(open ? 'nav.closeMenu' : 'nav.openMenu')}
            aria-expanded={open}
            aria-controls="main-nav"
          >
            {open ? <X className="size-5" /> : <Menu className="size-5" />}
          </Button>
          <div className="flex-1" />
          <NavLink
            to="/notifications"
            className="inline-flex h-8 items-center rounded-lg px-3 hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            aria-label={t('nav.notifications')}
          >
            <Bell className="size-5" />
          </NavLink>
          <AccountMenu onNotice={setNotice} />
        </header>
        <main id="main" tabIndex={-1} className="flex-1 p-4 outline-none md:p-8">
          {notice && (
            <div role="status" className="mb-4 flex items-start gap-3 rounded-lg border border-warning/40 bg-warning/15 p-3 text-sm">
              <p className="flex-1">{notice}</p>
              <Button variant="ghost" size="sm" onClick={() => setNotice(null)} aria-label={t('account.dismiss')}>
                <X className="size-4" />
              </Button>
            </div>
          )}
          <Outlet />
        </main>
      </div>
    </div>
  )
}
