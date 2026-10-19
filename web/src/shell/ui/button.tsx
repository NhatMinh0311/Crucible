import { cva, type VariantProps } from 'class-variance-authority'
import type { ButtonHTMLAttributes } from 'react'
import { cn } from './cn'

// Hướng "Bản vẽ kỹ thuật": góc vuông, viền 1 px, nút chính dùng phông tiêu đề (CR-DOC-11 mục 4).
const buttonVariants = cva(
  'inline-flex items-center justify-center gap-2 border text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: {
        primary: 'border-primary bg-primary font-display font-bold text-primary-foreground hover:bg-primary/90',
        outline: 'border-rule bg-surface text-foreground hover:bg-accent',
        ghost: 'border-transparent hover:bg-accent',
        danger: 'border-danger bg-danger font-display font-bold text-surface hover:bg-danger/90',
      },
      size: { sm: 'h-8 px-3', md: 'h-10 px-4', lg: 'h-11 px-5 text-base' },
    },
    defaultVariants: { variant: 'primary', size: 'md' },
  },
)

export type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & VariantProps<typeof buttonVariants>

export function Button({ className, variant, size, ...props }: ButtonProps) {
  return <button className={cn(buttonVariants({ variant, size }), className)} {...props} />
}
