import { useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from 'react-oidc-context'
import { OIDC_AUTHORITY } from '../auth/oidc'
import type { Language } from './i18n'
import { saveAccountLocale, storeLanguage } from './language'

/**
 * Đổi ngôn ngữ (AC-P2.2, W2): áp dụng ngay, lưu trên máy, ghi `locale` của tài khoản rồi làm mới token để dịch vụ
 * nhận claim mới. Trả về `false` khi không lưu được vào tài khoản; giao diện vẫn dùng ngôn ngữ mới trên máy này.
 */
export function useLanguage() {
  const auth = useAuth()
  const { i18n } = useTranslation()

  const setLanguage = useCallback(
    async (language: Language) => {
      storeLanguage(language)
      await i18n.changeLanguage(language)
      const token = auth.user?.access_token
      if (!token) return false
      try {
        await saveAccountLocale(OIDC_AUTHORITY, token, language)
      } catch {
        return false
      }
      await auth.signinSilent().catch(() => undefined)
      return true
    },
    [auth, i18n],
  )

  return { language: i18n.language as Language, setLanguage }
}
