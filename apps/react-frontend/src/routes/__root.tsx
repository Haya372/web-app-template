import { Toaster } from "@repo/ui";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { useState } from "react";
import { AuthProvider } from "@/features/auth/contexts/AuthContext";

export const Route = createRootRoute({
	component: RootLayout,
});

function RootLayout() {
	const [queryClient] = useState(() => new QueryClient());

	return (
		<QueryClientProvider client={queryClient}>
			<AuthProvider>
				<div className="min-h-screen bg-background font-sans antialiased">
					<Outlet />
					<Toaster />
					{import.meta.env.DEV && <TanStackRouterDevtools />}
				</div>
			</AuthProvider>
		</QueryClientProvider>
	);
}
