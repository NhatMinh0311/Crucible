import './shell/theme/theme.css'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { AuthProvider } from 'react-oidc-context'
import { RouterProvider } from 'react-router'
import { AuthGate } from './shell/auth/AuthGate'
import { oidcConfig } from './shell/auth/oidc'
import { router } from './shell/router'

const queryClient = new QueryClient()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider {...oidcConfig}>
      <AuthGate>
        <QueryClientProvider client={queryClient}>
          <RouterProvider router={router} />
        </QueryClientProvider>
      </AuthGate>
    </AuthProvider>
  </StrictMode>,
)
