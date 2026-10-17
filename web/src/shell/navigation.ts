import type { LucideIcon } from 'lucide-react'
import type { Role } from './auth/roles'

export type NavItem = { to: string; label: string; icon: LucideIcon; roles: Role[] }
