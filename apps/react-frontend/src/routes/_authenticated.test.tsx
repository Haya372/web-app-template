/**
 * Tests for _authenticated layout route.
 *
 * Asserts that the authenticated shell renders Header, Footer,
 * and the router Outlet placeholder.
 *
 * Mocks:
 *  - @tanstack/react-router   Outlet vi.mock
 *  - @/components/Header      default export vi.mock
 *  - @/components/Footer      default export vi.mock
 */

import React, { act } from "react"
import { createRoot } from "react-dom/client"
import type { Root } from "react-dom/client"
import { afterEach, describe, expect, it, vi } from "vitest"

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

function routeOptions(options: { component: unknown }) {
	return options
}

vi.mock("@tanstack/react-router", () => ({
	Outlet: () => React.createElement("div", { "data-testid": "outlet" }),
	createFileRoute: () => routeOptions,
}))

vi.mock("@/components/Header", () => ({
	default: () => React.createElement("header", { "data-testid": "header" }),
}))

vi.mock("@/components/Footer", () => ({
	default: () => React.createElement("footer", { "data-testid": "footer" }),
}))

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { AuthenticatedLayout } from "./_authenticated"

// ---------------------------------------------------------------------------
// Helpers
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
// Tests
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
