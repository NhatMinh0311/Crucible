import createClient, { type Middleware } from 'openapi-fetch'
import { User } from 'oidc-client-ts'
import { OIDC_AUTHORITY, OIDC_CLIENT_ID } from '../auth/oidc'
import type { paths as ContentPaths } from './generated/content'
import type { paths as RuntimePaths } from './generated/runtime'

/** Gắn access token của người dùng vào mọi yêu cầu tới dịch vụ. */
const auth: Middleware = {
  onRequest({ request }) {
    const raw = window.sessionStorage.getItem(`oidc.user:${OIDC_AUTHORITY}:${OIDC_CLIENT_ID}`)
    const user = raw ? User.fromStorageString(raw) : null
    if (user?.access_token) request.headers.set('Authorization', `Bearer ${user.access_token}`)
    return request
  },
}

function client<Paths extends object>(baseUrl: string) {
  const c = createClient<Paths>({ baseUrl })
  c.use(auth)
  return c
}

/** Client sinh từ contracts/ (npm run gen:api). Mỗi module tự thêm client cho API nội bộ của mình. */
export const runtimeApi = client<RuntimePaths>('/api/runtime')
export const contentApi = client<ContentPaths>('/api/content')
