import { Outlet, createRootRoute } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"
import { Toaster } from "@repo/ui"

export const Route = createRootRoute({
	component: RootLayout,
})

function RootLayout() {
	return (
		<div className="min-h-screen bg-background font-sans antialiased">
			<Outlet />
			<Toaster />
			{import.meta.env.DEV && <TanStackRouterDevtools />}
		</div>
	)
}
