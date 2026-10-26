import { createBrowserRouter, type RouteObject } from 'react-router'
import { contentRoutes } from '@/modules/content/routes'
import { learningRoutes } from '@/modules/learning/routes'
import { runtimeRoutes } from '@/modules/runtime/routes'
import { RoleGate } from './auth/RoleGate'
import { AppShell } from './layout/AppShell'
import { ErrorPage } from './pages/ErrorPage'
import { Home } from './pages/Home'

export const routes: RouteObject[] = [
  {
    path: '/',
    element: <AppShell />,
    errorElement: <ErrorPage />,
    children: [
      {
        // Route không có đường dẫn: lỗi và chặn theo vai trò hiển thị bên trong khung, menu vẫn dùng được.
        element: <RoleGate />,
        errorElement: <ErrorPage />,
        children: [
          { index: true, element: <Home /> },
          ...learningRoutes,
          ...contentRoutes,
          ...runtimeRoutes,
          { path: '*', element: <ErrorPage kind="notFound" /> },
        ],
      },
    ],
  },
]

export const router = createBrowserRouter(routes)
