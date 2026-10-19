import { beforeEach, describe, expect, it } from 'vitest'
import { initI18n } from './i18n'
import { problemMessage, traceReference } from './problem'

const i18n = initI18n('vi')

describe('thông điệp lỗi API theo code (AC-P2.3, W4)', () => {
  beforeEach(async () => {
    await i18n.changeLanguage('vi')
  })

  it('tra không gian của module trả lỗi trước', () => {
    expect(problemMessage(i18n, { code: 'SESSION_LIMIT', status: 409 }, 'runtime')).toMatch(/^Bạn đang có một phiên/)
  })

  it('mã dùng chung lấy từ shell', () => {
    expect(problemMessage(i18n, { code: 'FORBIDDEN', status: 403 }, 'content')).toBe('Tài khoản không có quyền thực hiện thao tác này.')
  })

  it('đổi ngôn ngữ thì thông điệp đổi theo', async () => {
    await i18n.changeLanguage('en')
    expect(problemMessage(i18n, { code: 'EXEC_BUSY' }, 'runtime')).toBe('Another check is running. Wait for its result and try again.')
  })

  it('mã chưa có bản dịch: thông điệp chung kèm code, không hiển thị detail', () => {
    const text = problemMessage(i18n, { code: 'LAB_NOT_PUBLISHED', detail: 'lab 42 version 3 draft' }, 'runtime')
    expect(text).toContain('LAB_NOT_PUBLISHED')
    expect(text).not.toContain('lab 42')
  })

  it('không có phản hồi: lỗi kết nối', () => {
    expect(problemMessage(i18n, undefined, 'learning')).toMatch(/^Không kết nối được/)
  })

  it('mã tra cứu là 8 ký tự đầu của traceId', () => {
    expect(traceReference({ traceId: '4bf92f3577b34da6a3ce929d0e0e4736' })).toBe('4bf92f35')
    expect(traceReference({ code: 'X' })).toBeUndefined()
  })
})
