import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import { callSignup } from "@/features/auth/api/signup"

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
		// Expose VITE_API_BASE_URL so the module under test can read it
		vi.stubEnv("VITE_API_BASE_URL", API_BASE_URL)
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

			const result = await callSignup("Alice", "user@example.com", "P@ssw0rd!")

			expect(result).toEqual(VALID_RESPONSE)
		})

		it("sends a POST request to the correct endpoint", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 201))
			vi.stubGlobal("fetch", mockFetch)

			await callSignup("Alice", "user@example.com", "P@ssw0rd!")

			expect(mockFetch).toHaveBeenCalledOnce()
			const [url, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			expect(url).toBe(`${API_BASE_URL}/v1/users/signup`)
			expect(init.method).toBe("POST")
		})

		it("sends Content-Type: application/json header", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 201))
			vi.stubGlobal("fetch", mockFetch)

			await callSignup("Alice", "user@example.com", "P@ssw0rd!")

			const [, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			const headers = new Headers(init.headers as HeadersInit)
			expect(headers.get("Content-Type")).toBe("application/json")
		})

		it("sends name, email, and password serialised as JSON in the request body", async () => {
			const mockFetch = vi
				.fn()
				.mockResolvedValue(makeFetchResponse(VALID_RESPONSE, 201))
			vi.stubGlobal("fetch", mockFetch)

			const name = "Alice"
			const email = "user@example.com"
			const password = "P@ssw0rd!"
			await callSignup(name, email, password)

			const [, init] = mockFetch.mock.calls[0] as [string, RequestInit]
			expect(JSON.parse(init.body as string)).toEqual({ name, email, password })
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
				vi.stubGlobal(
					"fetch",
					vi
						.fn()
						.mockResolvedValue(makeFetchResponse({ message: "error" }, status)),
				)

				await expect(
					callSignup("Alice", "user@example.com", "P@ssw0rd!"),
				).rejects.toThrow(String(status))
			})
		}
	})

	describe("network / runtime errors", () => {
		it("propagates a network-level error when fetch itself rejects", async () => {
			const networkError = new Error("Network failure")
			vi.stubGlobal("fetch", vi.fn().mockRejectedValue(networkError))

			await expect(
				callSignup("Alice", "user@example.com", "P@ssw0rd!"),
			).rejects.toThrow("Network failure")
		})
	})

	describe("configuration errors", () => {
		it("throws when VITE_API_BASE_URL is not set", async () => {
			vi.unstubAllEnvs()
			// Ensure the env var is absent
			vi.stubEnv("VITE_API_BASE_URL", "")

			await expect(
				callSignup("Alice", "user@example.com", "P@ssw0rd!"),
			).rejects.toThrow("VITE_API_BASE_URL")
		})
	})
})
