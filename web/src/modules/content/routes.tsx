// Giao diện của M2 Bài lab và chấm điểm. Chủ module sở hữu thư mục này (CR-DOC-07 mục 4).
import { BookOpen } from 'lucide-react'
import type { RouteObject } from 'react-router'
import type { NavItem } from '@/shell/navigation'
import { ComingSoon } from '@/shell/pages/ComingSoon'

export const contentNav: NavItem[] = [{ to: '/labs', labelKey: 'content:nav.labs', icon: BookOpen, roles: ['INSTRUCTOR'] }]

export const contentRoutes: RouteObject[] = [
  { path: 'labs', element: <ComingSoon labelKey="content:nav.labs" screen="SCR-23" /> },
]
