import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { callCreatePost } from "@/features/posts/api/createPost"

// ---------------------------------------------------------------------------
// Module mocks
// ---------------------------------------------------------------------------

vi.mock("@/generated/sdk.gen", () => ({
	postV1Posts: vi.fn(),
}))

vi.mock("@/utils/tokenStorage", () => ({
	getToken: vi.fn(),
}))

// Import after vi.mock so that the mocks are already in place
import { postV1Posts } from "@/generated/sdk.gen"
import { getToken } from "@/utils/tokenStorage"

const mockPostV1Posts = postV1Posts as ReturnType<typeof vi.fn>
const mockGetToken = getToken as ReturnType<typeof vi.fn>

// ---------------------------------------------------------------------------
// Fixtures
// ---------------------------------------------------------------------------

const VALID_RESPONSE = {
	id: "550e8400-e29b-41d4-a716-446655440000",
	userId: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	content: "Hello World!",
	createdAt: "2025-01-01T00:00:00Z",
}

const API_BASE_URL = "https://api.example.com"

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("callCreatePost", () => {
	beforeEach(() => {
		vi.stubEnv("VITE_API_BASE_URL", API_BASE_URL)
		mockGetToken.mockReturnValue("test-jwt-token")
	})

	afterEach(() => {
		vi.unstubAllEnvs()
		vi.resetAllMocks()
	})

	describe("happy path — 2xx response", () => {
		it("resolves with data when the server returns 201", async () => {
			mockPostV1Posts.mockResolvedValue({
				data: VALID_RESPONSE,
				error: undefined,
				response: { status: 201, ok: true },
			})

			const result = await callCreatePost("Hello World!")

			expect(result).toEqual(VALID_RESPONSE)
		})

		it("calls postV1Posts with the correct content body", async () => {
			mockPostV1Posts.mockResolvedValue({
				data: VALID_RESPONSE,
				error: undefined,
				response: { status: 201, ok: true },
			})

			await callCreatePost("Hello World!")

			expect(mockPostV1Posts).toHaveBeenCalledOnce()
			expect(mockPostV1Posts).toHaveBeenCalledWith(
				expect.objectContaining({ body: { content: "Hello World!" } }),
			)
		})

		it("passes the Authorization Bearer header from getToken", async () => {
			mockPostV1Posts.mockResolvedValue({
				data: VALID_RESPONSE,
				error: undefined,
				response: { status: 201, ok: true },
			})

			await callCreatePost("Hello World!")

			expect(mockPostV1Posts).toHaveBeenCalledWith(
				expect.objectContaining({
					headers: { Authorization: "Bearer test-jwt-token" },
				}),
			)
		})
	})

	describe("error path — non-2xx responses", () => {
		const errorCases = [
			{ name: "400 Bad Request", status: 400 },
			{ name: "401 Unauthorized", status: 401 },
			{ name: "403 Forbidden", status: 403 },
			{ name: "500 Internal Server Error", status: 500 },
		]

		for (const { name, status } of errorCases) {
			it(`throws an error containing the status code on ${name}`, async () => {
				mockPostV1Posts.mockResolvedValue({
					data: undefined,
					error: { type: "ERROR", title: "error", status },
					response: { status, ok: false },
				})

				await expect(callCreatePost("Hello World!")).rejects.toThrow(
					String(status),
				)
			})
		}
	})

	describe("network / runtime errors", () => {
		it("propagates a network-level error when the SDK rejects", async () => {
			const networkError = new Error("Network failure")
			mockPostV1Posts.mockRejectedValue(networkError)

			await expect(callCreatePost("Hello World!")).rejects.toThrow(
				"Network failure",
			)
		})
	})

	describe("authentication errors", () => {
		it("throws before making a network request when getToken returns null", async () => {
			mockGetToken.mockReturnValue(null)

			await expect(callCreatePost("Hello World!")).rejects.toThrow(
				"Unauthenticated",
			)
			expect(mockPostV1Posts).not.toHaveBeenCalled()
		})
	})
})
