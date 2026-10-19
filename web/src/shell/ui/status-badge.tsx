import { CheckCircle2, CircleDashed, CircleDot, Loader2, PauseCircle, XCircle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from './cn'

/** Trạng thái luôn có màu, biểu tượng và nhãn chữ (CR-DOC-11, UX7). */
const STATUS = {
  passed: { icon: CheckCircle2, className: 'text-success bg-success/10' },
  failed: { icon: XCircle, className: 'text-danger bg-danger/10' },
  checking: { icon: Loader2, className: 'text-info bg-info/10' },
  recorded: { icon: CircleDot, className: 'text-info bg-info/10' },
  interrupted: { icon: PauseCircle, className: 'text-warning bg-warning/15' },
  pending: { icon: CircleDashed, className: 'text-muted-foreground bg-muted' },
} as const

export type Status = keyof typeof STATUS

export function StatusBadge({ status, label }: { status: Status; label?: string }) {
  const { t } = useTranslation()
  const s = STATUS[status]
  const Icon = s.icon
  return (
    <span className={cn('inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium', s.className)}>
      <Icon className={cn('size-3.5', status === 'checking' && 'animate-spin')} aria-hidden />
      {label ?? t(`status.${status}`)}
    </span>
  )
}
