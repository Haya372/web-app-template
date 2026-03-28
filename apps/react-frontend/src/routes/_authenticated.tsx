import { Outlet, createFileRoute, redirect } from "@tanstack/react-router"
import Footer from "@/components/Footer"
import Header from "@/components/Header"
import { getToken } from "@/utils/tokenStorage"

export const Route = createFileRoute("/_authenticated")({
	beforeLoad: () => {
		if (!getToken()) {
			throw redirect({ to: "/login" })
		}
	},
	component: AuthenticatedLayout,
})

export function AuthenticatedLayout() {
	return (
		<>
			<Header />
			<Outlet />
			<Footer />
		</>
	)
}
