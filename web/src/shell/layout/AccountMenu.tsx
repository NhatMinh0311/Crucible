import { ChevronDown, LogOut, Monitor, Moon, Sun } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useAuth } from 'react-oidc-context'
import { useRoles } from '../auth/roles'
import { isLanguage, LANGUAGES } from '../i18n/i18n'
import { useLanguage } from '../i18n/useLanguage'
import { THEMES, type ThemePreference, useTheme } from '../theme/theme'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../ui/dropdown-menu'

const THEME_ICON = { light: Sun, dark: Moon, system: Monitor } as const

function initials(name: string | undefined) {
  const parts = (name ?? '').trim().split(/\s+/).filter(Boolean)
  return ((parts[0]?.[0] ?? '') + (parts.length > 1 ? (parts.at(-1)?.[0] ?? '') : '')).toUpperCase() || '?'
}

/** Menu tài khoản (CR-DOC-11 mục 4.1): thông tin tài khoản, ngôn ngữ, chế độ màu, đăng xuất. */
export function AccountMenu({ onNotice }: { onNotice: (message: string) => void }) {
  const auth = useAuth()
  const roles = useRoles()
  const { t } = useTranslation()
  const { language, setLanguage } = useLanguage()
  const [theme, setTheme] = useTheme()
  const profile = auth.user?.profile
  const name = profile?.name ?? profile?.preferred_username

  const changeLanguage = async (value: string) => {
    if (!isLanguage(value) || value === language) return
    // `t` của lần vẽ trước vẫn mang ngôn ngữ cũ; thông báo phải theo ngôn ngữ vừa chọn.
    if (!(await setLanguage(value))) onNotice(t('account.languageNotSaved', { lng: value }))
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        className="flex items-center gap-2 rounded-lg px-2 py-1 text-sm hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        aria-label={t('account.menu')}
      >
        <span className="flex size-8 items-center justify-center rounded-full bg-primary text-xs font-semibold text-primary-foreground" aria-hidden>
          {initials(name)}
        </span>
        <span className="hidden max-w-40 truncate sm:inline">{name}</span>
        <ChevronDown className="size-4 text-muted-foreground" aria-hidden />
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <div className="px-2 py-1.5">
          <p className="truncate text-sm font-semibold">{name}</p>
          {profile?.email && <p className="truncate text-xs text-muted-foreground">{profile.email}</p>}
          {roles.length > 0 && <p className="text-xs text-muted-foreground">{roles.map((r) => t(`roles.${r}`)).join(', ')}</p>}
        </div>
        <DropdownMenuSeparator />
        <DropdownMenuLabel>{t('account.language')}</DropdownMenuLabel>
        <DropdownMenuRadioGroup value={language} onValueChange={(v) => void changeLanguage(v)}>
          {LANGUAGES.map((l) => (
            <DropdownMenuRadioItem key={l} value={l} lang={l}>
              {t(`language.${l}`)}
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
        <DropdownMenuSeparator />
        <DropdownMenuLabel>{t('account.theme')}</DropdownMenuLabel>
        <DropdownMenuRadioGroup value={theme} onValueChange={(v) => setTheme(v as ThemePreference)}>
          {THEMES.map((th) => {
            const Icon = THEME_ICON[th]
            return (
              <DropdownMenuRadioItem key={th} value={th}>
                <Icon className="size-4 text-muted-foreground" aria-hidden />
                {t(`theme.${th}`)}
              </DropdownMenuRadioItem>
            )
          })}
        </DropdownMenuRadioGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem onSelect={() => void auth.signoutRedirect({ extraQueryParams: { ui_locales: language } })}>
          <LogOut className="size-4 text-muted-foreground" aria-hidden />
          {t('nav.signOut')}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
