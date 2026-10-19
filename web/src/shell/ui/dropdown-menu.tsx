import * as DropdownMenuPrimitive from '@radix-ui/react-dropdown-menu'
import { Check } from 'lucide-react'
import type { ComponentProps } from 'react'
import { cn } from './cn'

// Thành phần theo shadcn/ui (CR-DOC-11 mục 4); Radix lo điều hướng bàn phím và thuộc tính ARIA.
export const DropdownMenu = DropdownMenuPrimitive.Root
export const DropdownMenuTrigger = DropdownMenuPrimitive.Trigger
export const DropdownMenuGroup = DropdownMenuPrimitive.Group
export const DropdownMenuRadioGroup = DropdownMenuPrimitive.RadioGroup

export function DropdownMenuContent({ className, sideOffset = 6, ...props }: ComponentProps<typeof DropdownMenuPrimitive.Content>) {
  return (
    <DropdownMenuPrimitive.Portal>
      <DropdownMenuPrimitive.Content
        sideOffset={sideOffset}
        className={cn('z-50 min-w-60 rounded-lg border border-border bg-surface p-1 text-foreground shadow-lg', className)}
        {...props}
      />
    </DropdownMenuPrimitive.Portal>
  )
}

const itemClass =
  'relative flex cursor-default select-none items-center gap-2 rounded-md px-2 py-1.5 text-sm outline-none focus:bg-muted data-[disabled]:pointer-events-none data-[disabled]:opacity-50'

export function DropdownMenuItem({ className, ...props }: ComponentProps<typeof DropdownMenuPrimitive.Item>) {
  return <DropdownMenuPrimitive.Item className={cn(itemClass, className)} {...props} />
}

export function DropdownMenuRadioItem({ className, children, ...props }: ComponentProps<typeof DropdownMenuPrimitive.RadioItem>) {
  return (
    <DropdownMenuPrimitive.RadioItem className={cn(itemClass, 'pl-8', className)} {...props}>
      <span className="absolute left-2 flex size-4 items-center justify-center">
        <DropdownMenuPrimitive.ItemIndicator>
          <Check className="size-4" aria-hidden />
        </DropdownMenuPrimitive.ItemIndicator>
      </span>
      {children}
    </DropdownMenuPrimitive.RadioItem>
  )
}

export function DropdownMenuLabel({ className, ...props }: ComponentProps<typeof DropdownMenuPrimitive.Label>) {
  return <DropdownMenuPrimitive.Label className={cn('px-2 py-1.5 text-xs font-semibold text-muted-foreground', className)} {...props} />
}

export function DropdownMenuSeparator({ className, ...props }: ComponentProps<typeof DropdownMenuPrimitive.Separator>) {
  return <DropdownMenuPrimitive.Separator className={cn('-mx-1 my-1 h-px bg-border', className)} {...props} />
}
