import { createBrowserRouter } from 'react-router'
import { contentRoutes } from '@/modules/content/routes'
import { learningRoutes } from '@/modules/learning/routes'
import { runtimeRoutes } from '@/modules/runtime/routes'
import { AppShell } from './layout/AppShell'
import { ErrorPage } from './pages/ErrorPage'
import { Home } from './pages/Home'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppShell />,
    errorElement: <ErrorPage />,
    children: [
      {
        // Route không có đường dẫn để lỗi của một màn hình hiển thị bên trong khung, menu vẫn dùng được.
        errorElement: <ErrorPage />,
        children: [
          { index: true, element: <Home /> },
          ...learningRoutes,
          ...contentRoutes,
          ...runtimeRoutes,
          { path: '*', element: <ErrorPage notFound /> },
        ],
      },
    ],
  },
])
