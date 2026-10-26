import { cleanup, render, screen } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { initI18n } from '../i18n/i18n'
import { routes } from '../router'

let roles: string[] = []
vi.mock('react-oidc-context', () => ({
  useAuth: () => ({
    user: { access_token: `h.${btoa(JSON.stringify({ realm_access: { roles } }))}.s`, profile: { name: 'Người Dùng' } },
    signinSilent: vi.fn(),
    signoutRedirect: vi.fn(),
  }),
}))

const i18n = initI18n('vi')
const open = (path: string) => render(<RouterProvider router={createMemoryRouter(routes, { initialEntries: [path] })} />)

beforeEach(async () => {
  await i18n.changeLanguage('vi')
  vi.stubGlobal('matchMedia', () => ({ matches: true, addEventListener() {}, removeEventListener() {} }))
})
afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
})

describe('chặn route theo vai trò khai báo trên route (W8)', () => {
  it('sinh viên gõ thẳng /admin/cluster: trang không có quyền trong khung, có nút về trang chủ', () => {
    roles = ['STUDENT']
    open('/admin/cluster')
    expect(screen.getByText('Tài khoản không có quyền truy cập trang này.')).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Về trang chủ' })).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Lớp học' })).toBeTruthy()
  })

  it('quản trị vào được /admin/cluster', () => {
    roles = ['ADMIN']
    open('/admin/cluster')
    expect(screen.queryByText('Tài khoản không có quyền truy cập trang này.')).toBeNull()
    expect(screen.getByText(/SCR-40/)).toBeTruthy()
  })

  it('sinh viên không vào thư viện bài lab; giảng viên vào được', () => {
    roles = ['STUDENT']
    open('/labs')
    expect(screen.getByText('Tài khoản không có quyền truy cập trang này.')).toBeTruthy()
    cleanup()
    roles = ['INSTRUCTOR']
    open('/labs')
    expect(screen.getByText(/SCR-23/)).toBeTruthy()
  })

  it('route không khai báo vai trò thì ai cũng vào được', () => {
    roles = ['ADMIN']
    open('/notifications')
    expect(screen.getByText(/SCR-02/)).toBeTruthy()
  })
})
