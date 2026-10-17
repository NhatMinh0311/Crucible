import { createBrowserRouter } from 'react-router'
import { contentRoutes } from '@/modules/content/routes'
import { learningRoutes } from '@/modules/learning/routes'
import { runtimeRoutes } from '@/modules/runtime/routes'
import { AppShell } from './layout/AppShell'
import { Home } from './pages/Home'
import { NotFound } from './pages/NotFound'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppShell />,
    children: [
      { index: true, element: <Home /> },
      ...learningRoutes,
      ...contentRoutes,
      ...runtimeRoutes,
      { path: '*', element: <NotFound /> },
    ],
  },
])
