import { type ReactNode, useEffect } from 'react'
import { hasAuthParams, useAuth } from 'react-oidc-context'

/** Chuyển tới trang đăng nhập Keycloak khi chưa đăng nhập. */
export function AuthGate({ children }: { children: ReactNode }) {
  const auth = useAuth()
  useEffect(() => {
    if (!hasAuthParams() && !auth.isAuthenticated && !auth.activeNavigator && !auth.isLoading) {
      void auth.signinRedirect()
    }
  }, [auth])
  if (auth.error) return <p className="p-8 text-danger">Không đăng nhập được: {auth.error.message}</p>
  if (!auth.isAuthenticated) return <p className="p-8 text-muted-foreground">Đang chuyển tới trang đăng nhập…</p>
  return <>{children}</>
}
