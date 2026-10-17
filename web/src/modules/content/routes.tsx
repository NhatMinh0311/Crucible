// Giao diện của M2 Bài lab và chấm điểm. Chủ module sở hữu thư mục này (CR-DOC-07 mục 4).
import { BookOpen } from 'lucide-react'
import type { RouteObject } from 'react-router'
import type { NavItem } from '@/shell/navigation'

export const contentNav: NavItem[] = [{ to: '/labs', label: 'Thư viện bài lab', icon: BookOpen, roles: ['INSTRUCTOR'] }]

export const contentRoutes: RouteObject[] = [
  { path: 'labs', element: <p className="text-muted-foreground">Thư viện bài lab (SCR-23) — sẽ được xây dựng.</p> },
]
