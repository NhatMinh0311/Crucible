import { Bell, FlaskConical, LogOut, Menu, Moon, Sun } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from 'react-oidc-context'
import { NavLink, Outlet } from 'react-router'
import { learningNav } from '@/modules/learning/routes'
import { contentNav } from '@/modules/content/routes'
import { runtimeNav } from '@/modules/runtime/routes'
import { useRoles } from '../auth/roles'
import { Button } from '../ui/button'
import { cn } from '../ui/cn'

const NAV = [...learningNav, ...contentNav, ...runtimeNav]

export function AppShell() {
  const auth = useAuth()
  const { t } = useTranslation(['shell', 'runtime', 'content', 'learning'])
  const roles = useRoles()
  const [open, setOpen] = useState(false)
  const [dark, setDark] = useState(() => document.documentElement.classList.contains('dark'))
  const items = NAV.filter((i) => i.roles.some((r) => roles.includes(r)))

  const toggleDark = () => {
    document.documentElement.classList.toggle('dark')
    setDark(!dark)
  }

  return (
    <div className="flex min-h-dvh">
      <aside
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
                  'flex items-center gap-3 rounded-lg px-3 py-2 text-sm hover:bg-sidebar-active',
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
          <Button variant="ghost" size="sm" className="md:hidden" onClick={() => setOpen(!open)} aria-label={t('nav.openMenu')}>
            <Menu className="size-5" />
          </Button>
          <div className="flex-1" />
          <Button variant="ghost" size="sm" aria-label={t('nav.notifications')}>
            <Bell className="size-5" />
          </Button>
          <Button variant="ghost" size="sm" onClick={toggleDark} aria-label={t('nav.toggleTheme')}>
            {dark ? <Sun className="size-5" /> : <Moon className="size-5" />}
          </Button>
          <span className="hidden text-sm text-muted-foreground sm:inline">{auth.user?.profile.name}</span>
          <Button variant="ghost" size="sm" onClick={() => void auth.signoutRedirect()} aria-label={t('nav.signOut')}>
            <LogOut className="size-5" />
          </Button>
        </header>
        <main className="flex-1 p-4 md:p-8">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
