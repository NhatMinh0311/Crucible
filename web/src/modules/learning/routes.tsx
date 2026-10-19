// Giao diện của M3 Lớp học, lộ trình và thi. Chủ module sở hữu thư mục này (CR-DOC-07 mục 4).
import { GraduationCap } from 'lucide-react'
import type { RouteObject } from 'react-router'
import type { RouteHandle } from '@/shell/auth/RoleGate'
import type { NavItem } from '@/shell/navigation'
import { ComingSoon } from '@/shell/pages/ComingSoon'

export const learningNav: NavItem[] = [{ to: '/classes', labelKey: 'learning:nav.classes', icon: GraduationCap, roles: ['STUDENT', 'INSTRUCTOR'] }]

export const learningRoutes: RouteObject[] = [
  { path: 'notifications', element: <ComingSoon labelKey="learning:pages.notifications" screen="SCR-02" /> },
  { path: 'classes', handle: { roles: ['STUDENT', 'INSTRUCTOR'] } satisfies RouteHandle, element: <ComingSoon labelKey="learning:nav.classes" screen="SCR-11, SCR-21" /> },
]
