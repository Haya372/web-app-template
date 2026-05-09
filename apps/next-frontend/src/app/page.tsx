import { Button, Heading, Typography } from "@repo/ui";
import Link from "next/link";

export default function HomePage() {
	return (
		<main className="flex min-h-screen flex-col items-center justify-center gap-6 p-8">
			<Heading level={1}>Next.js Frontend</Heading>
			<Typography variant="lead">
				Bootstrapped for the web-app-template monorepo.
			</Typography>
			<Button asChild variant="outline">
				<Link href="/grpc-demo">gRPC Demo</Link>
			</Button>
		</main>
	);
}
