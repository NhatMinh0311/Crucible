import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { createMemoryRouter, RouterProvider } from 'react-router'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { RequireDesktop } from '../auth/RequireDesktop'
import { initI18n } from '../i18n/i18n'
import { ErrorPage } from '../pages/ErrorPage'
import { AppShell } from './AppShell'

let roles: string[] = []
vi.mock('react-oidc-context', () => ({
  useAuth: () => ({
    user: { access_token: `h.${btoa(JSON.stringify({ realm_access: { roles } }))}.s`, profile: { name: 'Người Dùng' } },
    signinSilent: vi.fn(),
    signoutRedirect: vi.fn(),
  }),
}))

const i18n = initI18n('vi')
let desktop = true

function Boom(): never {
  throw new Error('hỏng')
}

function renderAt(path: string) {
  const router = createMemoryRouter(
    [
      {
        path: '/',
        element: <AppShell />,
        children: [
          {
            errorElement: <ErrorPage />,
            children: [
              { index: true, element: <p>home</p> },
              { path: 'boom', element: <Boom /> },
              { path: 'desk', element: <RequireDesktop><p>desk-content</p></RequireDesktop> },
              { path: '*', element: <ErrorPage notFound /> },
            ],
          },
        ],
      },
    ],
    { initialEntries: [path] },
  )
  return render(<RouterProvider router={router} />)
}

beforeEach(async () => {
  await i18n.changeLanguage('vi')
  desktop = true
  vi.stubGlobal('matchMedia', () => ({ matches: desktop, addEventListener() {}, removeEventListener() {} }))
  vi.spyOn(console, 'error').mockImplementation(() => {})
})
afterEach(() => {
  cleanup()
  vi.unstubAllGlobals()
  vi.restoreAllMocks()
})

describe('khung giao diện (W8)', () => {
  it('menu chỉ hiện mục theo vai trò', () => {
    roles = ['STUDENT']
    renderAt('/')
    expect(screen.getByRole('link', { name: 'Lớp học' })).toBeTruthy()
    expect(screen.queryByRole('link', { name: 'Thư viện bài lab' })).toBeNull()
    expect(screen.queryByRole('link', { name: 'Sức chứa và phiên' })).toBeNull()
  })

  it('quản trị thấy mục của M1', () => {
    roles = ['ADMIN']
    renderAt('/')
    expect(screen.getByRole('link', { name: 'Sức chứa và phiên' })).toBeTruthy()
  })

  it('liên kết bỏ qua tới nội dung là phần tử nhận tiêu điểm đầu tiên', () => {
    roles = ['STUDENT']
    const { container } = renderAt('/')
    const first = container.querySelector('a, button, [tabindex]:not([tabindex="-1"])')
    expect(first?.textContent).toBe('Bỏ qua tới nội dung')
    expect(first?.getAttribute('href')).toBe('#main')
    expect(document.getElementById('main')).toBeTruthy()
  })

  it('menu trên điện thoại: mở có lớp phủ; Esc, bấm lớp phủ, chọn mục thì đóng', () => {
    roles = ['STUDENT']
    renderAt('/')
    const toggle = () => screen.getByRole('button', { name: /menu/i })
    fireEvent.click(toggle())
    expect(toggle().getAttribute('aria-expanded')).toBe('true')
    act(() => void window.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' })))
    expect(toggle().getAttribute('aria-expanded')).toBe('false')
    fireEvent.click(toggle())
    fireEvent.click(screen.getByTestId('nav-overlay'))
    expect(screen.queryByTestId('nav-overlay')).toBeNull()
    fireEvent.click(toggle())
    fireEvent.click(screen.getByRole('link', { name: 'Lớp học' }))
    expect(toggle().getAttribute('aria-expanded')).toBe('false')
  })

  it('chuông thông báo dẫn tới trang thông báo', () => {
    roles = ['STUDENT']
    renderAt('/')
    expect(screen.getByRole('link', { name: 'Thông báo' }).getAttribute('href')).toBe('/notifications')
  })

  it('đổi ngôn ngữ thì nhãn menu đổi theo', async () => {
    roles = ['INSTRUCTOR']
    renderAt('/')
    await act(async () => void (await i18n.changeLanguage('en')))
    expect(screen.getByRole('link', { name: 'Lab library' })).toBeTruthy()
  })
})

describe('RequireDesktop (W9)', () => {
  it('màn hình hẹp hơn 1024 px hiển thị yêu cầu dùng máy tính', () => {
    roles = ['STUDENT']
    desktop = false
    renderAt('/desk')
    expect(screen.queryByText('desk-content')).toBeNull()
    expect(screen.getByText('Màn hình này cần máy tính')).toBeTruthy()
  })

  it('máy tính hiển thị nội dung', () => {
    roles = ['STUDENT']
    renderAt('/desk')
    expect(screen.getByText('desk-content')).toBeTruthy()
  })
})

describe('trang lỗi (W10)', () => {
  it('đường dẫn không tồn tại: thông báo đã dịch, nút về trang chủ, menu vẫn còn', () => {
    roles = ['STUDENT']
    renderAt('/khong-co')
    expect(screen.getByText('Không tìm thấy trang.')).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Về trang chủ' }).getAttribute('href')).toBe('/')
    expect(screen.getByRole('link', { name: 'Lớp học' })).toBeTruthy()
  })

  it('màn hình ném lỗi: trang lỗi trong khung', () => {
    roles = ['STUDENT']
    renderAt('/boom')
    expect(screen.getByText('Đã xảy ra lỗi khi hiển thị trang.')).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Lớp học' })).toBeTruthy()
  })
})
