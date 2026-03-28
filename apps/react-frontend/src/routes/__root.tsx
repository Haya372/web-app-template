import { Outlet, createRootRoute } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"
import { Toaster } from "@repo/ui"
import { AuthProvider } from "@/features/auth/contexts/AuthContext"

export const Route = createRootRoute({
	component: RootLayout,
})

function RootLayout() {
	return (
		<AuthProvider>
			<div className="min-h-screen bg-background font-sans antialiased">
				<Outlet />
				<Toaster />
				{import.meta.env.DEV && <TanStackRouterDevtools />}
			</div>
		</AuthProvider>
	)
}
