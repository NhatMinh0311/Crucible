import { AlertTriangle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import type { ModuleNamespace } from '../i18n/i18n'
import { type ApiProblem, problemMessage, traceReference } from '../i18n/problem'
import { cn } from './cn'

/** Hiển thị lỗi API cho người dùng: thông điệp đã dịch theo `code` và mã tra cứu (UX5). */
export function ProblemAlert({ problem, ns, className }: { problem: ApiProblem | undefined; ns: ModuleNamespace; className?: string }) {
  const { t, i18n } = useTranslation()
  const ref = traceReference(problem)
  return (
    <div role="alert" className={cn('flex gap-3 rounded-lg border border-danger/30 bg-danger/10 p-3 text-sm', className)}>
      <AlertTriangle className="mt-0.5 size-4 shrink-0 text-danger" aria-hidden />
      <div>
        <p>{problemMessage(i18n, problem, ns)}</p>
        {ref && <p className="mt-1 text-xs text-muted-foreground">{t('errors.reference', { ref })}</p>}
      </div>
    </div>
  )
}
