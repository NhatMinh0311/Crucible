import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { initI18n } from '../i18n/i18n'
import { AccountMenu } from './AccountMenu'

const token = `h.${btoa(JSON.stringify({ realm_access: { roles: ['STUDENT'] } }))}.s`
const auth = {
  user: { access_token: token, profile: { name: 'Sinh Viên 1', email: 'sv1@crucible.local', locale: undefined } },
  signinSilent: vi.fn(async () => null),
  signoutRedirect: vi.fn(async () => undefined),
}
vi.mock('react-oidc-context', () => ({ useAuth: () => auth }))

const i18n = initI18n('vi')

/** Radix mở menu bằng phím; jsdom không có các API đo kích thước mà Radix dùng để định vị. */
function openMenu() {
  fireEvent.keyDown(screen.getByRole('button', { name: /Tài khoản|Account/ }), { key: 'Enter' })
}

beforeEach(async () => {
  window.localStorage.clear()
  await i18n.changeLanguage('vi')
  vi.stubGlobal('ResizeObserver', class { observe() {} unobserve() {} disconnect() {} })
  vi.stubGlobal('matchMedia', () => ({ matches: false, addEventListener() {}, removeEventListener() {} }))
})
afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
  vi.clearAllMocks()
})

describe('menu tài khoản (W2, W7)', () => {
  it('hiển thị tên, email, vai trò theo ngôn ngữ', () => {
    render(<AccountMenu onNotice={() => {}} />)
    openMenu()
    expect(screen.getByText('sv1@crucible.local')).toBeTruthy()
    expect(screen.getByText('Sinh viên')).toBeTruthy()
  })

  it('đổi sang tiếng Anh: giao diện đổi ngay, lưu locale vào tài khoản, làm mới token (AC-P2.2)', async () => {
    const fetchMock = vi.fn(async (_url: string, init?: RequestInit) =>
      init?.method === 'POST' ? new Response(null, { status: 204 }) : Response.json({ username: 'sv1' }),
    )
    vi.stubGlobal('fetch', fetchMock)
    const onNotice = vi.fn()
    render(<AccountMenu onNotice={onNotice} />)
    openMenu()
    await act(async () => fireEvent.click(screen.getByRole('menuitemradio', { name: 'English' })))
    expect(i18n.language).toBe('en')
    expect(document.documentElement.lang).toBe('en')
    expect(window.localStorage.getItem('crucible.locale')).toBe('en')
    await waitFor(() => expect(auth.signinSilent).toHaveBeenCalled())
    expect(JSON.parse(fetchMock.mock.calls[1]![1]!.body as string).attributes.locale).toEqual(['en'])
    expect(onNotice).not.toHaveBeenCalled()
  })

  it('lưu vào tài khoản thất bại: vẫn đổi ngôn ngữ trên máy, báo không lưu được', async () => {
    vi.stubGlobal('fetch', vi.fn(async () => new Response(null, { status: 401 })))
    const onNotice = vi.fn()
    render(<AccountMenu onNotice={onNotice} />)
    openMenu()
    await act(async () => fireEvent.click(screen.getByRole('menuitemradio', { name: 'English' })))
    expect(i18n.language).toBe('en')
    await waitFor(() => expect(onNotice).toHaveBeenCalledWith(expect.stringMatching(/^Could not save the language/)))
    expect(auth.signinSilent).not.toHaveBeenCalled()
  })

  it('chọn chế độ Tối', () => {
    render(<AccountMenu onNotice={() => {}} />)
    openMenu()
    fireEvent.click(screen.getByRole('menuitemradio', { name: 'Tối' }))
    expect(document.documentElement.classList.contains('dark')).toBe(true)
    expect(window.localStorage.getItem('crucible.theme')).toBe('dark')
  })

  it('đăng xuất truyền ngôn ngữ cho trang của Keycloak', () => {
    render(<AccountMenu onNotice={() => {}} />)
    openMenu()
    fireEvent.click(screen.getByRole('menuitem', { name: 'Đăng xuất' }))
    expect(auth.signoutRedirect).toHaveBeenCalledWith({ extraQueryParams: { ui_locales: 'vi' } })
  })
})
