// Giao diện của M1 Phiên thực hành và tài nguyên. Chủ module sở hữu thư mục này (CR-DOC-07 mục 4).
import { Server } from 'lucide-react'
import type { RouteObject } from 'react-router'
import type { RouteHandle } from '@/shell/auth/RoleGate'
import type { NavItem } from '@/shell/navigation'
import { ComingSoon } from '@/shell/pages/ComingSoon'

export const runtimeNav: NavItem[] = [{ to: '/admin/cluster', labelKey: 'runtime:nav.cluster', icon: Server, roles: ['ADMIN'] }]

export const runtimeRoutes: RouteObject[] = [
  { path: 'admin/cluster', handle: { roles: ['ADMIN'] } satisfies RouteHandle, element: <ComingSoon labelKey="runtime:nav.cluster" screen="SCR-40" /> },
]
