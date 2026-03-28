import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@repo/ui"

export function HomePage() {
	return (
		<main className="page-wrap px-4 py-12">
			<Card>
				<CardHeader>
					<CardTitle>TODO: Home page</CardTitle>
					<CardDescription>Post-login home page content goes here.</CardDescription>
				</CardHeader>
				<CardContent>
					<p className="text-sm text-muted-foreground">
						Implement authenticated content for this page.
					</p>
				</CardContent>
			</Card>
		</main>
	)
}
