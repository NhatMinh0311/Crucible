import { Bell, Menu, X } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { NavLink, Outlet } from 'react-router'
import { learningNav } from '@/modules/learning/routes'
import { contentNav } from '@/modules/content/routes'
import { runtimeNav } from '@/modules/runtime/routes'
import { useRoles } from '../auth/roles'
import { cn } from '../ui/cn'
import { SectionLabel } from '../ui/section-label'
import { AccountMenu } from './AccountMenu'

const NAV = [...learningNav, ...contentNav, ...runtimeNav]
const NAV_SECTION = 'A'

/** Dấu hiệu của Crucible: khung bản vẽ có đường tâm và vòng tròn. */
function Mark() {
  return (
    <svg className="size-[22px] text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" aria-hidden>
      <rect x="3" y="3" width="18" height="18" />
      <path d="M3 12h18M12 3v18" />
      <circle cx="12" cy="12" r="4" />
    </svg>
  )
}

const cell = 'flex items-center justify-center border-l border-rule hover:bg-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring'

/** Khung chung (P02) theo hướng "Bản vẽ kỹ thuật": thanh bên đánh số, header dạng ô, nền lưới (W8, W12). */
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
        className="sr-only z-50 bg-primary px-4 py-2 font-display font-bold text-primary-foreground focus:not-sr-only focus:fixed focus:top-2 focus:left-2"
      >
        {t('nav.skipToContent')}
      </a>
      {open && <div className="fixed inset-0 z-20 bg-foreground/40 md:hidden" onClick={() => setOpen(false)} aria-hidden data-testid="nav-overlay" />}
      <aside
        id="main-nav"
        className={cn(
          'fixed inset-y-0 left-0 z-30 flex w-56 -translate-x-full flex-col border-r border-rule bg-surface transition-transform md:static md:translate-x-0',
          open && 'translate-x-0',
        )}
      >
        <div className="flex h-16 shrink-0 items-center gap-2.5 border-b border-rule px-5">
          <Mark />
          <span className="font-display text-[19px] font-bold tracking-tight uppercase">Crucible</span>
        </div>
        <nav className="flex flex-col py-4" aria-label={t('nav.main')}>
          <SectionLabel index={NAV_SECTION} className="px-5 pb-2.5 text-muted-foreground">
            {t('nav.section')}
          </SectionLabel>
          {items.map((item, i) => (
            <NavLink
              key={item.to}
              to={item.to}
              onClick={() => setOpen(false)}
              className={({ isActive }) =>
                cn(
                  'flex h-10 items-center gap-3 border-l-[3px] border-transparent px-5 text-sm text-muted-foreground hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring',
                  isActive && 'border-primary bg-accent font-semibold text-foreground',
                )
              }
            >
              <span className="font-mono text-[11px] text-muted-foreground" aria-hidden>
                {String(i + 1).padStart(2, '0')}
              </span>
              {t(item.labelKey)}
            </NavLink>
          ))}
        </nav>
      </aside>

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-16 shrink-0 items-stretch border-b border-rule bg-surface">
          <button
            type="button"
            className={cn(cell, 'w-16 border-l-0 border-r md:hidden')}
            onClick={() => setOpen(!open)}
            aria-label={t(open ? 'nav.closeMenu' : 'nav.openMenu')}
            aria-expanded={open}
            aria-controls="main-nav"
          >
            {open ? <X className="size-5" /> : <Menu className="size-5" />}
          </button>
          <div className="flex-1" />
          <NavLink to="/notifications" className={cn(cell, 'w-16')} aria-label={t('nav.notifications')}>
            <Bell className="size-[18px]" />
          </NavLink>
          <AccountMenu onNotice={setNotice} />
        </header>
        <main id="main" tabIndex={-1} className="bg-blueprint flex-1 p-4 outline-none md:p-7">
          {notice && (
            <div role="status" className="mb-5 flex items-start gap-3 border border-signal bg-signal-soft py-2 pr-2 pl-3 text-sm">
              <p className="flex-1 py-1">{notice}</p>
              <button
                type="button"
                className="flex size-7 items-center justify-center hover:bg-surface focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                onClick={() => setNotice(null)}
                aria-label={t('account.dismiss')}
              >
                <X className="size-4" />
              </button>
            </div>
          )}
          <Outlet />
        </main>
      </div>
    </div>
  )
}
