import { Outlet, createRootRoute } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { Toaster } from '@repo/ui'
import Footer from '@/components/Footer'
import Header from '@/components/Header'

export const Route = createRootRoute({
  component: RootLayout,
})

function RootLayout() {
  return (
    <>
      <div className="min-h-screen bg-background font-sans antialiased">
        <Header />
        <Outlet />
        <Footer />
      </div>
      <Toaster />
      <TanStackRouterDevtools />
    </>
  )
}
