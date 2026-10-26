import { type ReactNode, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { hasAuthParams, useAuth } from 'react-oidc-context'
import { isLanguage } from '../i18n/i18n'
import { storeLanguage } from '../i18n/language'

/**
 * Chuyển tới trang đăng nhập Keycloak khi chưa đăng nhập; trang đăng nhập hiển thị theo ngôn ngữ đang chọn (W1).
 * Sau khi đăng nhập, ngôn ngữ theo claim `locale` của tài khoản nếu có (AC-P2.2).
 */
export function AuthGate({ children }: { children: ReactNode }) {
  const auth = useAuth()
  const { t, i18n } = useTranslation()
  useEffect(() => {
    if (!hasAuthParams() && !auth.isAuthenticated && !auth.activeNavigator && !auth.isLoading) {
      void auth.signinRedirect({ ui_locales: i18n.language })
    }
  }, [auth, i18n])
  const claim = auth.user?.profile.locale
  useEffect(() => {
    if (isLanguage(claim) && claim !== i18n.language) {
      storeLanguage(claim)
      void i18n.changeLanguage(claim)
    }
  }, [claim, i18n])
  if (auth.error) return <p className="p-8 text-danger">{t('auth.error', { message: auth.error.message })}</p>
  if (!auth.isAuthenticated) return <p className="p-8 text-muted-foreground">{t('auth.redirecting')}</p>
  return <>{children}</>
}
