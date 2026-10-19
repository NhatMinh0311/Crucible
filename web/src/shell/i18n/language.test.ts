import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { initI18n } from './i18n'
import { browserLanguage, initialLanguage, saveAccountLocale, storeLanguage } from './language'

describe('chọn ngôn ngữ (AC-P2.2, W1)', () => {
  beforeEach(() => window.localStorage.clear())

  it('tiếng Anh khi trình duyệt ưu tiên tiếng Anh, còn lại tiếng Việt', () => {
    expect(browserLanguage(['en-US', 'vi'])).toBe('en')
    expect(browserLanguage(['EN'])).toBe('en')
    expect(browserLanguage(['vi-VN', 'en'])).toBe('vi')
    expect(browserLanguage(['fr-FR', 'en'])).toBe('vi')
    expect(browserLanguage([])).toBe('vi')
  })

  it('claim locale của tài khoản thắng lựa chọn trên máy và trình duyệt', () => {
    storeLanguage('vi')
    expect(initialLanguage('en', ['vi'])).toBe('en')
  })

  it('bỏ qua claim locale không phải vi, en', () => {
    expect(initialLanguage('fr', ['en-US'])).toBe('en')
    expect(initialLanguage(undefined, ['vi'])).toBe('vi')
  })

  it('lựa chọn lưu trên máy thắng ngôn ngữ trình duyệt', () => {
    storeLanguage('en')
    expect(initialLanguage(undefined, ['vi-VN'])).toBe('en')
  })

  it('vẫn chạy khi trình duyệt chặn localStorage', () => {
    const spy = vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new DOMException('blocked')
    })
    expect(initialLanguage(undefined, ['en'])).toBe('en')
    spy.mockRestore()
  })
})

describe('thuộc tính lang của trang (W6)', () => {
  it('đổi theo ngôn ngữ đang chọn', async () => {
    const i18n = initI18n('vi')
    await i18n.changeLanguage('en')
    expect(document.documentElement.lang).toBe('en')
    await i18n.changeLanguage('vi')
    expect(document.documentElement.lang).toBe('vi')
  })
})

describe('lưu locale qua Account API (AC-P2.2, W2)', () => {
  afterEach(() => vi.unstubAllGlobals())

  it('giữ thông tin tài khoản và thuộc tính khác, chỉ đổi locale', async () => {
    const fetchMock = vi.fn(async (_url: string, init?: RequestInit) =>
      init?.method === 'POST'
        ? new Response(null, { status: 204 })
        : Response.json({
            id: 'x',
            username: 'sv1',
            email: 'sv1@crucible.local',
            firstName: 'Sinh',
            lastName: 'Viên 1',
            attributes: { studentId: ['2021'] },
            userProfileMetadata: {},
          }),
    )
    vi.stubGlobal('fetch', fetchMock)
    await saveAccountLocale('http://kc/realms/crucible', 'token', 'en')
    const [url, init] = fetchMock.mock.calls[1]!
    expect(url).toBe('http://kc/realms/crucible/account')
    expect(JSON.parse(init!.body as string)).toEqual({
      username: 'sv1',
      email: 'sv1@crucible.local',
      firstName: 'Sinh',
      lastName: 'Viên 1',
      attributes: { studentId: ['2021'], locale: ['en'] },
    })
  })

  it('báo lỗi khi Keycloak từ chối', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(null, { status: 401 })))
    await expect(saveAccountLocale('http://kc/realms/crucible', 'token', 'en')).rejects.toThrow('401')
  })
})
