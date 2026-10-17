import { useAuth } from 'react-oidc-context'

export type Role = 'STUDENT' | 'INSTRUCTOR' | 'ADMIN'

/** Đọc vai trò realm từ access token. Quyết định phân quyền thật nằm ở dịch vụ; giao diện chỉ dùng để hiển thị. */
export function rolesFromToken(accessToken: string | undefined): Role[] {
  if (!accessToken) return []
  const part = accessToken.split('.')[1]
  if (!part) return []
  try {
    const json = JSON.parse(atob(part.replace(/-/g, '+').replace(/_/g, '/'))) as {
      realm_access?: { roles?: string[] }
    }
    const roles = json.realm_access?.roles ?? []
    return roles.filter((r): r is Role => r === 'STUDENT' || r === 'INSTRUCTOR' || r === 'ADMIN')
  } catch {
    return []
  }
}

export function useRoles(): Role[] {
  const auth = useAuth()
  return rolesFromToken(auth.user?.access_token)
}
