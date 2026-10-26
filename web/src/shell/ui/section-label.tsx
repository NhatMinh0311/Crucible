import type { HTMLAttributes } from 'react'
import { cn } from './cn'

/**
 * Nhãn mục kiểu bản vẽ: "B — HƯỚNG DẪN". `index` là ký hiệu mục (A, B, 01…), không qua khóa dịch vì không phụ thuộc
 * ngôn ngữ; `children` là văn bản đã dịch.
 */
export function SectionLabel({ index, className, children, ...props }: HTMLAttributes<HTMLDivElement> & { index?: string }) {
  return (
    <div className={cn('label-mono text-primary', className)} {...props}>
      {index && <span aria-hidden>{index} — </span>}
      {children}
    </div>
  )
}
