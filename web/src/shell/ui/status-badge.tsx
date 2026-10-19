import { CheckCircle2, CircleDashed, CircleDot, Loader2, PauseCircle, XCircle } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from './cn'

/** Trạng thái luôn có màu, biểu tượng và nhãn chữ (CR-DOC-11, UX7); nhãn mono in hoa, nền nhạt cùng sắc. */
const STATUS = {
  passed: { icon: CheckCircle2, className: 'text-success bg-success-soft' },
  failed: { icon: XCircle, className: 'text-danger bg-danger-soft' },
  checking: { icon: Loader2, className: 'text-info bg-info-soft' },
  recorded: { icon: CircleDot, className: 'text-info bg-info-soft' },
  interrupted: { icon: PauseCircle, className: 'text-signal bg-signal-soft' },
  pending: { icon: CircleDashed, className: 'text-muted-foreground bg-muted' },
} as const

export type Status = keyof typeof STATUS

export function StatusBadge({ status, label }: { status: Status; label?: string }) {
  const { t } = useTranslation()
  const s = STATUS[status]
  const Icon = s.icon
  return (
    <span className={cn('label-mono inline-flex h-6 items-center gap-1.5 px-2 font-semibold tracking-[0.04em] whitespace-nowrap', s.className)}>
      <Icon className={cn('size-3.5', status === 'checking' && 'animate-spin')} aria-hidden />
      {label ?? t(`status.${status}`)}
    </span>
  )
}
