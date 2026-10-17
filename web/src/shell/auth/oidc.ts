import { WebStorageStateStore } from 'oidc-client-ts'
import type { AuthProviderProps } from 'react-oidc-context'

export const OIDC_AUTHORITY: string = import.meta.env.VITE_OIDC_AUTHORITY ?? 'http://localhost:8180/realms/crucible'
export const OIDC_CLIENT_ID: string = import.meta.env.VITE_OIDC_CLIENT_ID ?? 'crucible-web'

/** Đăng nhập bằng Authorization Code + PKCE với Keycloak (ADR-0004). */
export const oidcConfig: AuthProviderProps = {
  authority: OIDC_AUTHORITY,
  client_id: OIDC_CLIENT_ID,
  redirect_uri: `${window.location.origin}/`,
  post_logout_redirect_uri: `${window.location.origin}/`,
  scope: 'openid profile email',
  userStore: new WebStorageStateStore({ store: window.sessionStorage }),
  automaticSilentRenew: true,
  onSigninCallback: () => {
    window.history.replaceState({}, document.title, window.location.pathname)
  },
}
