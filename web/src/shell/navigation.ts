import type { ParseKeys } from 'i18next'
import type { LucideIcon } from 'lucide-react'
import type { Role } from './auth/roles'
import type { ModuleNamespace } from './i18n/i18n'

/** Khóa dịch có tiền tố không gian của module, ví dụ `content:nav.labs`. */
export type ModuleLabelKey = Extract<ParseKeys<ModuleNamespace[]>, `${ModuleNamespace}:${string}`>

export type NavItem = { to: string; labelKey: ModuleLabelKey; icon: LucideIcon; roles: Role[] }
