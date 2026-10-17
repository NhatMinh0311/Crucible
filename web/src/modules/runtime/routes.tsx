// Giao diện của M1 Phiên thực hành và tài nguyên. Chủ module sở hữu thư mục này (CR-DOC-07 mục 4).
import { Server } from 'lucide-react'
import type { RouteObject } from 'react-router'
import type { NavItem } from '@/shell/navigation'

export const runtimeNav: NavItem[] = [{ to: '/admin/cluster', label: 'Sức chứa và phiên', icon: Server, roles: ['ADMIN'] }]

export const runtimeRoutes: RouteObject[] = [
  { path: 'admin/cluster', element: <p className="text-muted-foreground">Sức chứa và phiên (SCR-40) — sẽ được xây dựng.</p> },
]
