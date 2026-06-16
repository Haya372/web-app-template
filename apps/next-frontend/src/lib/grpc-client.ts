import { type Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { HealthService } from "@/generated/proto/health/v1/health_pb";

const DEFAULT_BACKEND_GRPC_URL = "http://localhost:8081";

let cachedHealthClient: Client<typeof HealthService> | undefined;

export function getHealthClient(): Client<typeof HealthService> {
	if (cachedHealthClient !== undefined) {
		return cachedHealthClient;
	}
	// biome-ignore lint/style/noProcessEnv: server-side env var for Next.js BFF — no build-time alternative
	const baseUrl = process.env.BACKEND_GRPC_URL ?? DEFAULT_BACKEND_GRPC_URL;
	// NOTE: backend runs HTTP/1.1 only (go-backend/internal/infrastructure/connectrpc/server.go).
	// Switch httpVersion to "2" when h2c is enabled on the backend.
	const transport = createConnectTransport({
		baseUrl,
		httpVersion: "1.1",
	});
	cachedHealthClient = createClient(HealthService, transport);
	return cachedHealthClient;
}

// Test only: resets module-level client cache. Do not call from production code.
export function __resetHealthClientForTest(): void {
	cachedHealthClient = undefined;
}
