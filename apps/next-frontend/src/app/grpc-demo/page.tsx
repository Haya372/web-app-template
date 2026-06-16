import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
	Heading,
	Typography,
} from "@repo/ui";
import type { ServingStatus } from "@/generated/proto/health/v1/health_pb";
import { getHealthClient } from "@/lib/grpc-client";
import { formatServingStatus } from "./labels";

type CheckResponse = { status: ServingStatus };
type CheckableClient = {
	check: (req: Record<string, never>) => Promise<CheckResponse>;
};

export default async function GrpcDemoPage() {
	let statusLabel: string;
	try {
		const client = getHealthClient() as unknown as CheckableClient;
		const response = await client.check({});
		statusLabel = formatServingStatus(response.status);
	} catch {
		statusLabel = "Unreachable";
	}

	return (
		<main className="flex min-h-screen flex-col items-center justify-center gap-6 p-8">
			<Heading level={1}>gRPC Demo</Heading>
			<Card>
				<CardHeader>
					<CardTitle>Backend Health</CardTitle>
				</CardHeader>
				<CardContent>
					<Typography>{statusLabel}</Typography>
				</CardContent>
			</Card>
		</main>
	);
}
