// Giao diện của M3 Lớp học, lộ trình và thi. Chủ module sở hữu thư mục này (CR-DOC-07 mục 4).
import { GraduationCap } from 'lucide-react'
import type { RouteObject } from 'react-router'
import type { NavItem } from '@/shell/navigation'

export const learningNav: NavItem[] = [{ to: '/classes', label: 'Lớp học', icon: GraduationCap, roles: ['STUDENT', 'INSTRUCTOR'] }]

export const learningRoutes: RouteObject[] = [
  { path: 'classes', element: <p className="text-muted-foreground">Lớp học (SCR-11, SCR-21) — sẽ được xây dựng.</p> },
]
