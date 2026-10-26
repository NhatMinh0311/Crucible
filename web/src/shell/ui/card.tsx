import type { HTMLAttributes } from 'react'
import { cn } from './cn'

/** Thẻ: viền 1 px màu kẻ chính, không bo góc, không bóng đổ. */
export function Card({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cn('border border-rule bg-surface p-5', className)} {...props} />
}

export function CardTitle({ className, ...props }: HTMLAttributes<HTMLHeadingElement>) {
  return <h2 className={cn('font-display text-lg font-semibold tracking-tight', className)} {...props} />
}
