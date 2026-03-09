import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { callLogin } from "@/features/auth/api/login"

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
		// Expose VITE_API_BASE_URL so the module under test can read it
		vi.stubEnv("VITE_API_BASE_URL", API_BASE_URL)
	})

	afterEach(() => {
		vi.unstubAllGlobals()
		vi.unstubAllEnvs()
	})

	describe("happy path — 2xx response", () => {
		it("resolves with parsed JSON when the server returns 200", async () => {
			vi.stubGlobal(
				"fetch",
				vi.fn().mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 200)),
			)

			const result = await callLogin("test@example.com", "password123")

			expect(result).toEqual(VALID_RESPONSE)
		})

		it("sends a POST request to the correct endpoint", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 200))
			vi.stubGlobal("fetch", mockFetch)

			await callLogin("user@example.com", "s3cr3t")

			expect(mockFetch).toHaveBeenCalledOnce()
			const [url, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			expect(url).toBe(`${API_BASE_URL}/v1/users/login`)
			expect(init.method).toBe("POST")
		})

		it("sends Content-Type: application/json header", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 200))
			vi.stubGlobal("fetch", mockFetch)

			await callLogin("user@example.com", "s3cr3t")

			const [, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			const headers = new Headers(init.headers as HeadersInit)
			expect(headers.get("Content-Type")).toBe("application/json")
		})

		it("sends email and password serialised as JSON in the request body", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 200))
			vi.stubGlobal("fetch", mockFetch)

			const email = "user@example.com"
			const password = "s3cr3t"
			await callLogin(email, password)

			const [, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			expect(JSON.parse(init.body as string)).toEqual({ email, password })
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
				vi.stubGlobal(
					"fetch",
					vi
						.fn()
						.mockResolvedValue(makeFetchResponse({ message: "error" }, status)),
				)

				await expect(
					callLogin("user@example.com", "wrong-password"),
				).rejects.toThrow(String(status))
			})
		}
	})

	describe("network / runtime errors", () => {
		it("propagates a network-level error when fetch itself rejects", async () => {
			const networkError = new Error("Network failure")
			vi.stubGlobal("fetch", vi.fn().mockRejectedValue(networkError))

			await expect(
				callLogin("user@example.com", "password123"),
			).rejects.toThrow("Network failure")
		})
	})
})
