import i18next from 'i18next'
import { initReactI18next } from 'react-i18next'
import enContent from '@/locales/en/content.json'
import enLearning from '@/locales/en/learning.json'
import enRuntime from '@/locales/en/runtime.json'
import enShell from '@/locales/en/shell.json'
import viContent from '@/locales/vi/content.json'
import viLearning from '@/locales/vi/learning.json'
import viRuntime from '@/locales/vi/runtime.json'
import viShell from '@/locales/vi/shell.json'

/** Ngôn ngữ giao diện (CR-DOC-11 mục 4.1); tiếng Việt là mặc định. */
export const LANGUAGES = ['vi', 'en'] as const
export type Language = (typeof LANGUAGES)[number]

/** Không gian khóa: `shell` cho khung chung, mỗi module một không gian và chỉ sửa không gian của mình. */
export const NAMESPACES = ['shell', 'runtime', 'content', 'learning'] as const
export type Namespace = (typeof NAMESPACES)[number]
export type ModuleNamespace = Exclude<Namespace, 'shell'>

export const resources = {
  vi: { shell: viShell, runtime: viRuntime, content: viContent, learning: viLearning },
  en: { shell: enShell, runtime: enRuntime, content: enContent, learning: enLearning },
} as const

export function isLanguage(value: unknown): value is Language {
  return typeof value === 'string' && (LANGUAGES as readonly string[]).includes(value)
}

/** Khởi tạo đồng bộ: tệp dịch nằm sẵn trong bundle nên lần hiển thị đầu đã đúng ngôn ngữ. */
export function initI18n(lng: Language) {
  if (!i18next.isInitialized) {
    i18next.on('languageChanged', (l) => {
      document.documentElement.lang = l
    })
    void i18next.use(initReactI18next).init({
      resources,
      lng,
      fallbackLng: 'vi',
      ns: NAMESPACES,
      defaultNS: 'shell',
      interpolation: { escapeValue: false },
      initAsync: false,
    })
  }
  document.documentElement.lang = i18next.language
  return i18next
}
