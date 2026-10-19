import { MonitorSmartphone } from 'lucide-react'
import { type ReactNode, useSyncExternalStore } from 'react'
import { useTranslation } from 'react-i18next'

const DESKTOP_QUERY = '(min-width: 1024px)'

function subscribe(onChange: () => void) {
  const media = window.matchMedia(DESKTOP_QUERY)
  media.addEventListener('change', onChange)
  return () => media.removeEventListener('change', onChange)
}

/** Màn hình chỉ dùng trên máy tính (CR-DOC-11 mục 3, W9): trang làm bài, phòng chờ thi, soạn bài lab, chạy thử. */
export function RequireDesktop({ children }: { children: ReactNode }) {
  const { t } = useTranslation()
  const desktop = useSyncExternalStore(subscribe, () => window.matchMedia(DESKTOP_QUERY).matches)
  if (desktop) return <>{children}</>
  return (
    <div className="mx-auto flex max-w-sm flex-col items-center gap-3 py-16 text-center">
      <MonitorSmartphone className="size-10 text-muted-foreground" aria-hidden />
      <p className="font-semibold">{t('pages.desktopOnly')}</p>
      <p className="text-sm text-muted-foreground">{t('pages.desktopOnlyHint')}</p>
    </div>
  )
}
