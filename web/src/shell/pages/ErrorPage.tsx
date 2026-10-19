import { Home } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { isRouteErrorResponse, Link, useRouteError } from 'react-router'

/** Trang lỗi của router (W10): đường dẫn không tồn tại hoặc lỗi khi hiển thị một màn hình. */
export function ErrorPage({ notFound = false }: { notFound?: boolean }) {
  const { t } = useTranslation()
  const error = useRouteError()
  const missing = notFound || (isRouteErrorResponse(error) && error.status === 404)
  if (error) console.error(error)
  return (
    <div className="mx-auto max-w-md space-y-4 py-16 text-center">
      <h1 className="text-xl font-semibold">{t(missing ? 'pages.notFound' : 'pages.error')}</h1>
      {!missing && <p className="text-sm text-muted-foreground">{t('pages.errorHint')}</p>}
      <Link
        to="/"
        className="inline-flex h-10 items-center gap-2 rounded-lg bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
      >
        <Home className="size-4" aria-hidden />
        {t('pages.backHome')}
      </Link>
    </div>
  )
}
