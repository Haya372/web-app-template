import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { callCreatePost } from "@/features/posts/api/createPost"

// ---------------------------------------------------------------------------
// Module mocks
// ---------------------------------------------------------------------------

vi.mock("@/features/auth/utils/tokenStorage", () => ({
	getToken: vi.fn(),
}))

// Import after vi.mock so that the mock is already in place
import { getToken } from "@/features/auth/utils/tokenStorage"

const mockGetToken = getToken as ReturnType<typeof vi.fn>

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/** Builds a minimal Response-like object that globalThis.fetch can return. */
function makeFetchResponse(body: unknown, status: number): Response {
	return {
		ok: status >= 200 && status < 300,
		status,
		json: () => Promise.resolve(body),
	} as Response
}

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
		// Expose VITE_API_BASE_URL so the module under test can read it
		vi.stubEnv("VITE_API_BASE_URL", API_BASE_URL)
		// Provide a token for every test by default
		mockGetToken.mockReturnValue("test-jwt-token")
	})

	afterEach(() => {
		vi.unstubAllGlobals()
		vi.unstubAllEnvs()
	})

	describe("happy path — 2xx response", () => {
		it("resolves with parsed JSON when the server returns 201", async () => {
			vi.stubGlobal(
				"fetch",
				vi.fn().mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 201)),
			)

			const result = await callCreatePost("Hello World!")

			expect(result).toEqual(VALID_RESPONSE)
		})

		it("sends a POST request to the correct endpoint", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 201))
			vi.stubGlobal("fetch", mockFetch)

			await callCreatePost("Hello World!")

			expect(mockFetch).toHaveBeenCalledOnce()
			const [url, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			expect(url).toBe(`${API_BASE_URL}/v1/posts`)
			expect(init.method).toBe("POST")
		})

		it("sends Content-Type: application/json header", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 201))
			vi.stubGlobal("fetch", mockFetch)

			await callCreatePost("Hello World!")

			const [, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			const headers = new Headers(init.headers as HeadersInit)
			expect(headers.get("Content-Type")).toBe("application/json")
		})

		it("sends Authorization Bearer header with the token from getToken", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 201))
			vi.stubGlobal("fetch", mockFetch)

			await callCreatePost("Hello World!")

			const [, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			const headers = new Headers(init.headers as HeadersInit)
			expect(headers.get("Authorization")).toBe("Bearer test-jwt-token")
		})

		it("sends the content serialised as JSON in the request body", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 201))
			vi.stubGlobal("fetch", mockFetch)

			const content = "Hello World!"
			await callCreatePost(content)

			const [, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			expect(JSON.parse(init.body as string)).toEqual({ content })
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
				vi.stubGlobal(
					"fetch",
					vi
						.fn()
						.mockResolvedValue(makeFetchResponse({ message: "error" }, status)),
				)

				await expect(callCreatePost("Hello World!")).rejects.toThrow(
					String(status),
				)
			})
		}
	})

	describe("network / runtime errors", () => {
		it("propagates a network-level error when fetch itself rejects", async () => {
			const networkError = new Error("Network failure")
			vi.stubGlobal("fetch", vi.fn().mockRejectedValue(networkError))

			await expect(callCreatePost("Hello World!")).rejects.toThrow(
				"Network failure",
			)
		})
	})

	describe("authentication errors", () => {
		it("throws before making a network request when getToken returns null", async () => {
			const mockFetch = vi.fn()
			vi.stubGlobal("fetch", mockFetch)
			mockGetToken.mockReturnValue(null)

			await expect(callCreatePost("Hello World!")).rejects.toThrow(
				"Unauthenticated",
			)
			expect(mockFetch).not.toHaveBeenCalled()
		})
	})
})
