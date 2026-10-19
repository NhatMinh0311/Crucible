import { describe, expect, it } from 'vitest'
import { formatDate, formatDateTime, formatNumber, formatTime } from './format'

// 08:30 UTC là 15:30 giờ Việt Nam; kết quả không phụ thuộc múi giờ của máy chạy kiểm thử.
const instant = '2026-10-19T08:30:00Z'

describe('định dạng ngày giờ, số (AC-P2.3, W5)', () => {
  it('giờ theo Asia/Ho_Chi_Minh ở cả hai ngôn ngữ', () => {
    expect(formatDateTime(instant, 'vi')).toBe('15:30 19 thg 10, 2026')
    expect(formatDateTime(instant, 'en')).toBe('19 Oct 2026, 15:30')
    expect(formatTime(instant, 'en')).toBe('15:30')
  })

  it('qua nửa đêm theo giờ Việt Nam thì sang ngày hôm sau', () => {
    expect(formatDate('2026-10-19T17:30:00Z', 'vi')).toBe('20 thg 10, 2026')
  })

  it('số theo ngôn ngữ', () => {
    expect(formatNumber(1234567.89, 'vi')).toBe('1.234.567,89')
    expect(formatNumber(1234567.89, 'en')).toBe('1,234,567.89')
  })

  it('ngôn ngữ lạ dùng định dạng tiếng Việt', () => {
    expect(formatNumber(1000.5, 'fr')).toBe('1.000,5')
  })
})
