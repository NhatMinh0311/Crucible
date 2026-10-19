import { isLanguage, type Language } from './i18n'

const STORAGE_KEY = 'crucible.locale'

/** AC-P2.2: tiếng Anh nếu ngôn ngữ ưu tiên đầu tiên của trình duyệt là tiếng Anh, còn lại tiếng Việt. */
export function browserLanguage(languages: readonly string[] = navigator.languages): Language {
  return languages[0]?.toLowerCase().startsWith('en') ? 'en' : 'vi'
}

export function storedLanguage(): Language | undefined {
  try {
    const value = window.localStorage.getItem(STORAGE_KEY)
    return isLanguage(value) ? value : undefined
  } catch {
    return undefined
  }
}

export function storeLanguage(language: Language) {
  try {
    window.localStorage.setItem(STORAGE_KEY, language)
  } catch {
    // Trình duyệt chặn lưu trữ: giao diện vẫn chạy, chỉ không nhớ lựa chọn.
  }
}

/** Thứ tự: claim `locale` của tài khoản, lựa chọn lưu trên máy, ngôn ngữ của trình duyệt. */
export function initialLanguage(claim?: unknown, languages?: readonly string[]): Language {
  if (isLanguage(claim)) return claim
  return storedLanguage() ?? browserLanguage(languages)
}

type AccountRepresentation = {
  username: string
  email?: string
  firstName?: string
  lastName?: string
  attributes?: Record<string, string[]>
}

/** Ghi thuộc tính `locale` của người dùng qua Account API của Keycloak. */
export async function saveAccountLocale(authority: string, accessToken: string, language: Language) {
  const url = `${authority}/account`
  const headers = { Authorization: `Bearer ${accessToken}`, Accept: 'application/json' }
  const current = await fetch(url, { headers })
  if (!current.ok) throw new Error(`Account API ${current.status}`)
  const { username, email, firstName, lastName, attributes } = (await current.json()) as AccountRepresentation
  const updated = await fetch(url, {
    method: 'POST',
    headers: { ...headers, 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, firstName, lastName, attributes: { ...attributes, locale: [language] } }),
  })
  if (!updated.ok) throw new Error(`Account API ${updated.status}`)
}
