/**
 * Tests for HomePage component.
 *
 * Asserts that the page renders the expected placeholder heading text.
 *
 * Mocks:
 *  - @repo/ui   Card family components vi.mock (avoid style resolution issues)
 */

import React, { act } from "react";
import type { Root } from "react-dom/client";
import { createRoot } from "react-dom/client";
import { afterEach, describe, expect, it, vi } from "vitest";

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

vi.mock("@repo/ui", () => ({
	Card: ({ children }: { children: React.ReactNode }) =>
		React.createElement("div", { "data-testid": "card" }, children),
	CardHeader: ({ children }: { children: React.ReactNode }) =>
		React.createElement("div", { "data-testid": "card-header" }, children),
	CardTitle: ({ children }: { children: React.ReactNode }) =>
		React.createElement("h2", { "data-testid": "card-title" }, children),
	CardDescription: ({ children }: { children: React.ReactNode }) =>
		React.createElement("p", { "data-testid": "card-description" }, children),
	CardContent: ({ children }: { children: React.ReactNode }) =>
		React.createElement("div", { "data-testid": "card-content" }, children),
}));

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { HomePage } from "./HomePage";

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

let root: Root | null = null;
let container: HTMLDivElement | null = null;

async function mount(): Promise<void> {
	const div = document.createElement("div");
	document.body.append(div);
	container = div;
	await act(async () => {
		root = createRoot(div);
		root.render(<HomePage />);
	});
}

// ---------------------------------------------------------------------------
// Teardown
// ---------------------------------------------------------------------------

afterEach(async () => {
	if (root) {
		const r = root;
		root = null;
		await act(async () => {
			r.unmount();
		});
	}
	container?.remove();
	container = null;
});

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("HomePage — rendering", () => {
	it("renders the 'TODO: Home page' heading text", async () => {
		await mount();
		expect(
			document.querySelector("[data-testid='card-title']")?.textContent,
		).toContain("TODO: Home page");
	});

	it("renders a card container", async () => {
		await mount();
		expect(document.querySelector("[data-testid='card']")).not.toBeNull();
	});

	it("renders the card description text", async () => {
		await mount();
		expect(
			document.querySelector("[data-testid='card-description']")?.textContent,
		).toContain("Post-login home page content goes here.");
	});
});
