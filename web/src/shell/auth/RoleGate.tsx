import { Outlet, useMatches } from 'react-router'
import { ErrorPage } from '../pages/ErrorPage'
import { type Role, useRoles } from './roles'

/** Thông tin route mà module khai báo: `handle: { roles: ['ADMIN'] }` chỉ cho vai trò trong danh sách vào route. */
export type RouteHandle = { roles?: Role[] }

/**
 * Chặn route theo vai trò khai báo trên route, để module không phải tự bọc từng màn hình. Đây chỉ là trải nghiệm:
 * quyền thật do dịch vụ kiểm tra (CR-DOC-03 mục 3.2).
 */
export function RoleGate() {
  const mine = useRoles()
  const allowed = useMatches().every((m) => {
    const roles = (m.handle as RouteHandle | undefined)?.roles
    return !roles || roles.some((r) => mine.includes(r))
  })
  return allowed ? <Outlet /> : <ErrorPage kind="forbidden" />
}
