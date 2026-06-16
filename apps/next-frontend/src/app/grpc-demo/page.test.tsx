import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import GrpcDemoPage from "./page";

vi.mock("@/lib/grpc-client", () => ({
	getHealthClient: vi.fn(),
}));

import { getHealthClient } from "@/lib/grpc-client";

const mockGetHealthClient = vi.mocked(getHealthClient);

describe("GrpcDemoPage", () => {
	beforeEach(() => {
		mockGetHealthClient.mockReset();
	});

	it("renders the gRPC Demo heading", async () => {
		mockGetHealthClient.mockReturnValue({
			check: async () => ({ status: 1 }),
		} as never);
		render(await GrpcDemoPage());
		expect(
			screen.getByRole("heading", { level: 1, name: /grpc demo/i }),
		).toBeInTheDocument();
	});

	it("displays 'Serving' when backend returns SERVING status", async () => {
		mockGetHealthClient.mockReturnValue({
			check: async () => ({ status: 1 }),
		} as never);
		render(await GrpcDemoPage());
		expect(screen.getByText("Serving")).toBeInTheDocument();
	});

	it("displays 'Unreachable' when backend throws", async () => {
		mockGetHealthClient.mockReturnValue({
			check: async () => {
				throw new Error("connection refused");
			},
		} as never);
		render(await GrpcDemoPage());
		expect(screen.getByText("Unreachable")).toBeInTheDocument();
	});
});
