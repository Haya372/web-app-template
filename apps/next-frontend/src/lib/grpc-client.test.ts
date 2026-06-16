import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const { mockCreateClient, mockCreateConnectTransport } = vi.hoisted(() => {
	const mockCreateConnectTransport = vi.fn(() => ({ _fake: "transport" }));
	const mockCreateClient = vi.fn(() => ({ _fake: "client" }));
	return { mockCreateClient, mockCreateConnectTransport };
});

vi.mock("@connectrpc/connect", () => ({
	createClient: mockCreateClient,
}));

vi.mock("@connectrpc/connect-node", () => ({
	createConnectTransport: mockCreateConnectTransport,
}));

vi.mock("@/generated/proto/health/v1/health_pb", () => ({
	HealthService: { typeName: "health.v1.HealthService" },
}));

import { __resetHealthClientForTest, getHealthClient } from "./grpc-client";

describe("getHealthClient", () => {
	beforeEach(() => {
		__resetHealthClientForTest();
		mockCreateClient.mockClear();
		mockCreateConnectTransport.mockClear();
		mockCreateClient.mockReturnValue({ _fake: "client" } as never);
		mockCreateConnectTransport.mockReturnValue({ _fake: "transport" } as never);
	});

	afterEach(() => {
		vi.unstubAllEnvs();
	});

	it("returns the same client instance on repeated calls", () => {
		const first = getHealthClient();
		const second = getHealthClient();
		expect(first).toBe(second);
		expect(mockCreateClient).toHaveBeenCalledTimes(1);
	});

	it("uses BACKEND_GRPC_URL env var as baseUrl", () => {
		vi.stubEnv("BACKEND_GRPC_URL", "http://custom-backend:9090");
		getHealthClient();
		expect(mockCreateConnectTransport).toHaveBeenCalledWith({
			baseUrl: "http://custom-backend:9090",
			httpVersion: "1.1",
		});
	});

	it("falls back to http://localhost:8081 when BACKEND_GRPC_URL is not set", () => {
		getHealthClient();
		expect(mockCreateConnectTransport).toHaveBeenCalledWith({
			baseUrl: "http://localhost:8081",
			httpVersion: "1.1",
		});
	});
});
