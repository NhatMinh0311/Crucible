import { act, renderHook } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { readTheme, THEME_STORAGE_KEY, useTheme } from './theme'

/** Giả lập prefers-color-scheme; `set` đổi cài đặt của hệ điều hành và báo cho bên đăng ký. */
function mockSystemDark(initial: boolean) {
  let dark = initial
  const listeners = new Set<() => void>()
  vi.stubGlobal('matchMedia', (query: string) => ({
    get matches() {
      return query.includes('dark') && dark
    },
    addEventListener: (_: string, l: () => void) => listeners.add(l),
    removeEventListener: (_: string, l: () => void) => listeners.delete(l),
  }))
  return {
    set(value: boolean) {
      dark = value
      listeners.forEach((l) => l())
    },
  }
}

const isDark = () => document.documentElement.classList.contains('dark')

describe('chế độ màu (W7)', () => {
  beforeEach(() => {
    window.localStorage.clear()
    document.documentElement.classList.remove('dark')
  })
  afterEach(() => vi.unstubAllGlobals())

  it('lần đầu theo hệ điều hành', () => {
    mockSystemDark(true)
    expect(readTheme()).toBe('system')
    renderHook(() => useTheme())
    expect(isDark()).toBe(true)
  })

  it('lựa chọn được lưu trên trình duyệt và áp dụng ngay', () => {
    mockSystemDark(true)
    const { result } = renderHook(() => useTheme())
    act(() => result.current[1]('light'))
    expect(isDark()).toBe(false)
    expect(window.localStorage.getItem(THEME_STORAGE_KEY)).toBe('light')
    expect(readTheme()).toBe('light')
  })

  it('Theo hệ thống: đổi theo khi hệ điều hành đổi; chế độ cố định thì không', () => {
    const system = mockSystemDark(false)
    const { result } = renderHook(() => useTheme())
    expect(isDark()).toBe(false)
    act(() => system.set(true))
    expect(isDark()).toBe(true)
    act(() => result.current[1]('light'))
    act(() => system.set(true))
    expect(isDark()).toBe(false)
  })

  it('giá trị lưu không hợp lệ được coi là Theo hệ thống', () => {
    window.localStorage.setItem(THEME_STORAGE_KEY, 'tím')
    expect(readTheme()).toBe('system')
  })
})
