import i18next from 'i18next'
import { useTranslation } from 'react-i18next'
import { isLanguage, type Language } from './i18n'

/** Giờ luôn hiển thị theo giờ Việt Nam, không theo múi giờ của máy (AC-P2.3, AGENTS.md mục 4). */
export const TIME_ZONE = 'Asia/Ho_Chi_Minh'

// en-GB: giờ 24 tiếng, ngày trước tháng như tiếng Việt; giờ thi đọc giống nhau ở hai ngôn ngữ.
const INTL_LOCALE: Record<Language, string> = { vi: 'vi-VN', en: 'en-GB' }

type DateInput = Date | string | number

const intlLocale = (language: string) => INTL_LOCALE[isLanguage(language) ? language : 'vi']
const toDate = (value: DateInput) => (value instanceof Date ? value : new Date(value))

export function formatDateTime(value: DateInput, language: string = i18next.language) {
  return new Intl.DateTimeFormat(intlLocale(language), { dateStyle: 'medium', timeStyle: 'short', timeZone: TIME_ZONE }).format(toDate(value))
}

export function formatDate(value: DateInput, language: string = i18next.language) {
  return new Intl.DateTimeFormat(intlLocale(language), { dateStyle: 'medium', timeZone: TIME_ZONE }).format(toDate(value))
}

export function formatTime(value: DateInput, language: string = i18next.language) {
  return new Intl.DateTimeFormat(intlLocale(language), { timeStyle: 'short', timeZone: TIME_ZONE }).format(toDate(value))
}

export function formatNumber(value: number, language: string = i18next.language, options?: Intl.NumberFormatOptions) {
  return new Intl.NumberFormat(intlLocale(language), options).format(value)
}

/** Hàm định dạng theo ngôn ngữ đang chọn; thành phần dùng hook này được vẽ lại khi đổi ngôn ngữ. */
export function useFormat() {
  const { i18n } = useTranslation()
  const language = i18n.language
  return {
    dateTime: (value: DateInput) => formatDateTime(value, language),
    date: (value: DateInput) => formatDate(value, language),
    time: (value: DateInput) => formatTime(value, language),
    number: (value: number, options?: Intl.NumberFormatOptions) => formatNumber(value, language, options),
  }
}
