import { useCallback, useEffect, useState } from 'react'

/** Chế độ màu (W7): lựa chọn lưu trên trình duyệt; `system` theo cài đặt của hệ điều hành. */
export const THEMES = ['light', 'dark', 'system'] as const
export type ThemePreference = (typeof THEMES)[number]

// Cùng khóa với script trong index.html, script đó đặt lớp `dark` trước khi React chạy để không nháy màu.
export const THEME_STORAGE_KEY = 'crucible.theme'
const DARK_QUERY = '(prefers-color-scheme: dark)'

export function readTheme(): ThemePreference {
  try {
    const value = window.localStorage.getItem(THEME_STORAGE_KEY)
    return (THEMES as readonly string[]).includes(value ?? '') ? (value as ThemePreference) : 'system'
  } catch {
    return 'system'
  }
}

export function applyTheme(preference: ThemePreference) {
  const dark = preference === 'dark' || (preference === 'system' && window.matchMedia(DARK_QUERY).matches)
  document.documentElement.classList.toggle('dark', dark)
  document.documentElement.style.colorScheme = dark ? 'dark' : 'light'
}

export function useTheme() {
  const [theme, setThemeState] = useState<ThemePreference>(readTheme)

  useEffect(() => {
    applyTheme(theme)
    if (theme !== 'system') return
    const media = window.matchMedia(DARK_QUERY)
    const onChange = () => applyTheme('system')
    media.addEventListener('change', onChange)
    return () => media.removeEventListener('change', onChange)
  }, [theme])

  const setTheme = useCallback((preference: ThemePreference) => {
    try {
      window.localStorage.setItem(THEME_STORAGE_KEY, preference)
    } catch {
      // Không lưu được thì chỉ áp dụng cho lần mở trang này.
    }
    setThemeState(preference)
  }, [])

  return [theme, setTheme] as const
}
