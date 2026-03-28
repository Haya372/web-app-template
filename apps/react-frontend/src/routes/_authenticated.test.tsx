/**
 * Tests for _authenticated layout route.
 *
 * Asserts that:
 *  1. The authenticated shell renders Header, Footer, and the router Outlet.
 *  2. The beforeLoad guard does NOT throw when a token is present.
 *  3. The beforeLoad guard throws a redirect to /login when no token is present.
 *
 * Mocks:
 *  - @tanstack/react-router   Outlet / redirect / createFileRoute   vi.mock
 *  - @/components/Header      default export                        vi.mock
 *  - @/components/Footer      default export                        vi.mock
 *  - @/utils/tokenStorage     getToken                              vi.mock
 */

import React, { act } from "react"
import { createRoot } from "react-dom/client"
import type { Root } from "react-dom/client"
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"

// ---------------------------------------------------------------------------
// Hoisted mock functions (must be defined before vi.mock calls)
// ---------------------------------------------------------------------------

const { mockGetToken, mockRedirectResult } = vi.hoisted(() => ({
	mockGetToken: vi.fn<() => string | null>(),
	// redirect() returns a sentinel object; the beforeLoad guard throws it
	mockRedirectResult: { __isRedirect: true, to: "/login" },
}))

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

function passThrough(options: Record<string, unknown>) {
	return options
}

// Pass the full route config through so both `component` and `beforeLoad` are
// accessible on the exported Route object.
vi.mock("@tanstack/react-router", () => ({
	Outlet: () => React.createElement("div", { "data-testid": "outlet" }),
	createFileRoute: () => passThrough,
	redirect: vi.fn(() => mockRedirectResult),
}))

vi.mock("@/components/Header", () => ({
	default: () => React.createElement("header", { "data-testid": "header" }),
}))

vi.mock("@/components/Footer", () => ({
	default: () => React.createElement("footer", { "data-testid": "footer" }),
}))

vi.mock("@/utils/tokenStorage", () => ({
	getToken: mockGetToken,
}))

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { redirect } from "@tanstack/react-router"
import { AuthenticatedLayout, Route } from "./_authenticated"

// ---------------------------------------------------------------------------
// Typed mock references
// ---------------------------------------------------------------------------

const mockRedirect = redirect as ReturnType<typeof vi.fn>

// ---------------------------------------------------------------------------
// Extract beforeLoad from the Route config
// ---------------------------------------------------------------------------

// With the mock in place, Route is the raw config object passed to
// createFileRoute(...)({...}), so beforeLoad is directly accessible.
const { beforeLoad } = Route as unknown as {
	beforeLoad: () => void
	component: React.ComponentType
}

// ---------------------------------------------------------------------------
// Helpers — component rendering
// ---------------------------------------------------------------------------

let root: Root | null = null
let container: HTMLDivElement | null = null

async function mount(): Promise<void> {
	const div = document.createElement("div")
	document.body.append(div)
	container = div
	await act(async () => {
		root = createRoot(div)
		root.render(<AuthenticatedLayout />)
	})
}

// ---------------------------------------------------------------------------
// Teardown
// ---------------------------------------------------------------------------

afterEach(async () => {
	if (root) {
		const r = root
		root = null
		await act(async () => {
			r.unmount()
		})
	}
	container?.remove()
	container = null
})

// ---------------------------------------------------------------------------
// Setup
// ---------------------------------------------------------------------------

beforeEach(() => {
	mockGetToken.mockReset()
	mockRedirect.mockReset()
	mockRedirect.mockReturnValue(mockRedirectResult)
})

// ---------------------------------------------------------------------------
// Tests — component rendering
// ---------------------------------------------------------------------------

describe("AuthenticatedLayout — rendering", () => {
	it("renders the Header component", async () => {
		await mount()
		expect(document.querySelector("[data-testid='header']")).not.toBeNull()
	})

	it("renders the Footer component", async () => {
		await mount()
		expect(document.querySelector("[data-testid='footer']")).not.toBeNull()
	})

	it("renders the router Outlet placeholder", async () => {
		await mount()
		expect(document.querySelector("[data-testid='outlet']")).not.toBeNull()
	})
})

// ---------------------------------------------------------------------------
// Tests — beforeLoad guard
// ---------------------------------------------------------------------------

describe("_authenticated — beforeLoad guard", () => {
	it("does not throw when getToken returns a non-null token", () => {
		mockGetToken.mockReturnValue("valid-jwt-token")

		expect(() => beforeLoad()).not.toThrow()
	})

	it("throws a redirect to /login when getToken returns null", () => {
		mockGetToken.mockReturnValue(null)

		let thrown: unknown
		try {
			beforeLoad()
		} catch (error) {
			thrown = error
		}

		// The guard must throw — and what it throws must be the redirect sentinel
		expect(thrown).toBeDefined()
		expect(thrown).toBe(mockRedirectResult)
		expect(mockRedirect).toHaveBeenCalledWith({ to: "/login" })
	})
})
