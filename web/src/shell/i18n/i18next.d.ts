import 'i18next'
import type { resources } from './i18n'

// Khóa dịch được kiểm tra khi biên dịch theo tệp tiếng Việt; tệp tiếng Anh phải có cùng tập khóa (locales.test.ts).
declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'shell'
    resources: (typeof resources)['vi']
  }
}
