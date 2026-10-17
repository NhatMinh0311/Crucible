import { describe, expect, it } from 'vitest'
import { rolesFromToken } from './roles'

const token = (payload: object) => `h.${btoa(JSON.stringify(payload))}.s`

describe('rolesFromToken', () => {
  it('chỉ lấy vai trò của Crucible từ realm_access', () => {
    expect(rolesFromToken(token({ realm_access: { roles: ['STUDENT', 'offline_access'] } }))).toEqual(['STUDENT'])
  })
  it('trả về rỗng khi token sai định dạng', () => {
    expect(rolesFromToken('khong-phai-jwt')).toEqual([])
    expect(rolesFromToken(undefined)).toEqual([])
  })
})
