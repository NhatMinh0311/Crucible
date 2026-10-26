import type { i18n as I18n } from 'i18next'
import type { components } from '../api/generated/runtime'
import type { ModuleNamespace } from './i18n'

/** Lỗi RFC 9457 của mọi module (CR-DOC-03 mục 3.4); các hợp đồng dùng chung lược đồ này. */
export type ApiProblem = Partial<components['schemas']['Problem']>

/** Mã tra cứu người dùng gửi cho giảng viên, quản trị để tìm log: 8 ký tự đầu của `traceId`. */
export function traceReference(problem: ApiProblem | undefined) {
  return problem?.traceId ? problem.traceId.slice(0, 8) : undefined
}

/**
 * Thông điệp theo `code` (AC-P2.3, W4): tra `errors.<code>` ở không gian của module trả lỗi rồi tới `shell`.
 * `problem` rỗng nghĩa là không nhận được phản hồi (mất kết nối). `detail` không bao giờ được hiển thị.
 */
export function problemMessage(i18n: I18n, problem: ApiProblem | undefined, ns: ModuleNamespace) {
  const t = i18n.getFixedT(i18n.language, 'shell')
  if (!problem) return t('errors.network')
  const code = problem.code
  if (code) {
    for (const n of [ns, 'shell']) {
      const text = i18n.getResource(i18n.language, n, `errors.${code}`) as unknown
      if (typeof text === 'string') return text
    }
  }
  return t('errors.unknown', { code: code ?? String(problem.status ?? '') })
}
