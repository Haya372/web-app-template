import { Outlet, createFileRoute } from "@tanstack/react-router"
import Footer from "@/components/Footer"
import Header from "@/components/Header"

export const Route = createFileRoute("/_authenticated")({
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
