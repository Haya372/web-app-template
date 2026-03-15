import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { callLogin } from "@/features/auth/api/login"

// ---------------------------------------------------------------------------
// Module mocks
// ---------------------------------------------------------------------------

vi.mock("@/generated/sdk.gen", () => ({
	postV1UsersLogin: vi.fn(),
}))

// Import after vi.mock so that the mock is already in place
import { postV1UsersLogin } from "@/generated/sdk.gen"

const mockPostV1UsersLogin = postV1UsersLogin as ReturnType<typeof vi.fn>

// ---------------------------------------------------------------------------
// Fixtures
// ---------------------------------------------------------------------------

const VALID_RESPONSE = {
	token: "test-jwt-token",
	expiresAt: "2026-03-08T00:00:00Z",
	user: {
		id: "user-id-123",
		name: "Test User",
		email: "test@example.com",
	},
}

const API_BASE_URL = "https://api.example.com"

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("callLogin", () => {
	beforeEach(() => {
		vi.stubEnv("VITE_API_BASE_URL", API_BASE_URL)
	})

	afterEach(() => {
		vi.unstubAllEnvs()
		vi.resetAllMocks()
	})

	describe("happy path — 2xx response", () => {
		it("resolves with data when the server returns 200", async () => {
			mockPostV1UsersLogin.mockResolvedValue({
				data: VALID_RESPONSE,
				error: undefined,
				response: { status: 200, ok: true },
			})

			const result = await callLogin("test@example.com", "password123")

			expect(result).toEqual(VALID_RESPONSE)
		})

		it("calls postV1UsersLogin with the correct email and password body", async () => {
			mockPostV1UsersLogin.mockResolvedValue({
				data: VALID_RESPONSE,
				error: undefined,
				response: { status: 200, ok: true },
			})

			await callLogin("user@example.com", "s3cr3t")

			expect(mockPostV1UsersLogin).toHaveBeenCalledOnce()
			expect(mockPostV1UsersLogin).toHaveBeenCalledWith(
				expect.objectContaining({
					body: { email: "user@example.com", password: "s3cr3t" },
				}),
			)
		})
	})

	describe("error path — non-2xx responses", () => {
		const errorCases = [
			{ name: "400 Bad Request", status: 400 },
			{ name: "401 Unauthorized", status: 401 },
			{ name: "403 Forbidden", status: 403 },
			{ name: "404 Not Found", status: 404 },
			{ name: "422 Unprocessable Entity", status: 422 },
			{ name: "500 Internal Server Error", status: 500 },
			{ name: "503 Service Unavailable", status: 503 },
		]

		for (const { name, status } of errorCases) {
			it(`throws an error containing the status code on ${name}`, async () => {
				mockPostV1UsersLogin.mockResolvedValue({
					data: undefined,
					error: { type: "ERROR", title: "error", status },
					response: { status, ok: false },
				})

				await expect(
					callLogin("user@example.com", "wrong-password"),
				).rejects.toThrow(String(status))
			})
		}
	})

	describe("network / runtime errors", () => {
		it("propagates a network-level error when the SDK rejects", async () => {
			const networkError = new Error("Network failure")
			mockPostV1UsersLogin.mockRejectedValue(networkError)

			await expect(
				callLogin("user@example.com", "password123"),
			).rejects.toThrow("Network failure")
		})
	})
})
