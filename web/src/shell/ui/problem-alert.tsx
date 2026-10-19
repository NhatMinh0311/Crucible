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
    <div role="alert" className={cn('grid grid-cols-[2.25rem_minmax(0,1fr)] border border-signal bg-surface text-sm', className)}>
      <div className="flex items-center justify-center bg-signal text-signal-foreground">
        <AlertTriangle className="size-4" aria-hidden />
      </div>
      <div className="space-y-1 px-3 py-2.5">
        <p className="leading-relaxed">{problemMessage(i18n, problem, ns)}</p>
        {ref && <p className="label-mono text-[10.5px] text-muted-foreground">{t('errors.reference', { ref })}</p>}
      </div>
    </div>
  )
}
