import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { callSignup } from "@/features/auth/api/signup"

// ---------------------------------------------------------------------------
// Module mocks
// ---------------------------------------------------------------------------

vi.mock("@/generated/sdk.gen", () => ({
	postV1UsersSignup: vi.fn(),
}))

// Import after vi.mock so that the mock is already in place
import { postV1UsersSignup } from "@/generated/sdk.gen"

const mockPostV1UsersSignup = postV1UsersSignup as ReturnType<typeof vi.fn>

// ---------------------------------------------------------------------------
// Fixtures
// ---------------------------------------------------------------------------

const VALID_RESPONSE = {
	id: "550e8400-e29b-41d4-a716-446655440000",
	name: "Alice",
	email: "user@example.com",
	status: "active",
	createdAt: "2025-01-01T00:00:00Z",
}

const API_BASE_URL = "https://api.example.com"

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("callSignup", () => {
	beforeEach(() => {
		vi.stubEnv("VITE_API_BASE_URL", API_BASE_URL)
	})

	afterEach(() => {
		vi.unstubAllEnvs()
		vi.resetAllMocks()
	})

	describe("happy path — 2xx response", () => {
		it("resolves with data when the server returns 201", async () => {
			mockPostV1UsersSignup.mockResolvedValue({
				data: VALID_RESPONSE,
				error: undefined,
				response: { status: 201, ok: true },
			})

			const result = await callSignup("Alice", "user@example.com", "P@ssw0rd!")

			expect(result).toEqual(VALID_RESPONSE)
		})

		it("calls postV1UsersSignup with the correct name, email, and password body", async () => {
			mockPostV1UsersSignup.mockResolvedValue({
				data: VALID_RESPONSE,
				error: undefined,
				response: { status: 201, ok: true },
			})

			await callSignup("Alice", "user@example.com", "P@ssw0rd!")

			expect(mockPostV1UsersSignup).toHaveBeenCalledOnce()
			expect(mockPostV1UsersSignup).toHaveBeenCalledWith(
				expect.objectContaining({
					body: { name: "Alice", email: "user@example.com", password: "P@ssw0rd!" },
				}),
			)
		})
	})

	describe("error path — non-2xx responses", () => {
		const errorCases = [
			{ name: "400 Bad Request (validation error)", status: 400 },
			{ name: "409 Conflict (duplicate email)", status: 409 },
			{ name: "500 Internal Server Error", status: 500 },
		]

		for (const { name, status } of errorCases) {
			it(`throws an error containing the status code on ${name}`, async () => {
				mockPostV1UsersSignup.mockResolvedValue({
					data: undefined,
					error: { type: "ERROR", title: "error", status },
					response: { status, ok: false },
				})

				await expect(
					callSignup("Alice", "user@example.com", "P@ssw0rd!"),
				).rejects.toThrow(String(status))
			})
		}
	})

	describe("network / runtime errors", () => {
		it("propagates a network-level error when the SDK rejects", async () => {
			const networkError = new Error("Network failure")
			mockPostV1UsersSignup.mockRejectedValue(networkError)

			await expect(
				callSignup("Alice", "user@example.com", "P@ssw0rd!"),
			).rejects.toThrow("Network failure")
		})
	})

	describe("configuration errors", () => {
		it("throws when VITE_API_BASE_URL is not set", async () => {
			vi.unstubAllEnvs()
			vi.stubEnv("VITE_API_BASE_URL", "")

			await expect(
				callSignup("Alice", "user@example.com", "P@ssw0rd!"),
			).rejects.toThrow("VITE_API_BASE_URL")
		})
	})
})
