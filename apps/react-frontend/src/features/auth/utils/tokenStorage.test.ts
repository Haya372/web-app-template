import { afterEach, describe, expect, it } from "vitest"
import {
	getToken,
	removeToken,
	saveToken,
} from "@/features/auth/utils/tokenStorage"

// ---------------------------------------------------------------------------
// The localStorage key the module is expected to use.
// Keeping it here makes the intent explicit and guards against typos.
// ---------------------------------------------------------------------------
const AUTH_TOKEN_KEY = "auth_token"

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("tokenStorage", () => {
	// Always start each test with a clean slate so tests do not bleed into
	// each other regardless of execution order.
	afterEach(() => {
		localStorage.clear()
	})

	// -------------------------------------------------------------------------
	// saveToken
	// -------------------------------------------------------------------------
	describe("saveToken", () => {
		it("stores the token in localStorage under the key 'auth_token'", () => {
			saveToken("my-jwt")

			expect(localStorage.getItem(AUTH_TOKEN_KEY)).toBe("my-jwt")
		})

		const overwriteCases = [
			{ name: "short token", first: "first-token", second: "second-token" },
			{ name: "empty string token", first: "some-token", second: "" },
		]

		for (const { name, first, second } of overwriteCases) {
			it(`overwrites a previously saved token — ${name}`, () => {
				saveToken(first)
				saveToken(second)

				expect(localStorage.getItem(AUTH_TOKEN_KEY)).toBe(second)
			})
		}
	})

	// -------------------------------------------------------------------------
	// getToken
	// -------------------------------------------------------------------------
	describe("getToken", () => {
		it("returns the token that was previously saved", () => {
			localStorage.setItem(AUTH_TOKEN_KEY, "stored-jwt")

			expect(getToken()).toBe("stored-jwt")
		})

		it("returns null when no token has been saved", () => {
			// localStorage starts empty thanks to afterEach
			expect(getToken()).toBeNull()
		})

		it("returns null after the token has been removed", () => {
			localStorage.setItem(AUTH_TOKEN_KEY, "to-be-removed")
			localStorage.removeItem(AUTH_TOKEN_KEY)

			expect(getToken()).toBeNull()
		})
	})

	// -------------------------------------------------------------------------
	// removeToken
	// -------------------------------------------------------------------------
	describe("removeToken", () => {
		it("removes an existing token so getToken returns null afterwards", () => {
			saveToken("jwt-to-delete")

			removeToken()

			expect(getToken()).toBeNull()
		})

		it("does not throw when called with no token present", () => {
			// Expect removeToken to be a safe no-op when nothing is stored
			expect(() => removeToken()).not.toThrow()
		})

		it("does not affect other keys in localStorage", () => {
			localStorage.setItem("other_key", "other_value")
			saveToken("jwt-to-delete")

			removeToken()

			expect(localStorage.getItem("other_key")).toBe("other_value")
		})
	})

	// -------------------------------------------------------------------------
	// Round-trip: save → get → remove → get
	// -------------------------------------------------------------------------
	describe("full lifecycle round-trip", () => {
		it("save then get returns the saved token, remove then get returns null", () => {
			const token = "round-trip-token"

			saveToken(token)
			expect(getToken()).toBe(token)

			removeToken()
			expect(getToken()).toBeNull()
		})
	})
})
