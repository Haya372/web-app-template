import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { Outlet, createRootRoute } from "@tanstack/react-router"
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools"
import { useState } from "react"
import { Toaster } from "@repo/ui"

export const Route = createRootRoute({
	component: RootLayout,
})

function RootLayout() {
	const [queryClient] = useState(() => new QueryClient())

	return (
		<QueryClientProvider client={queryClient}>
			<div className="min-h-screen bg-background font-sans antialiased">
				<Outlet />
				<Toaster />
				{import.meta.env.DEV && <TanStackRouterDevtools />}
			</div>
		</QueryClientProvider>
	)
}
