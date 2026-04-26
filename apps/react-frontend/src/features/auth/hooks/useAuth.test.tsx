/**
 * Tests for useAuth hook.
 *
 * useAuth reads from AuthContext and throws when called outside AuthProvider.
 * The hook is exercised via minimal components rendered with createRoot + act.
 *
 * Mocks:
 *  - @/utils/tokenStorage   getToken / saveToken / removeToken   vi.mock
 *    (required because AuthProvider imports tokenStorage at module load time)
 */

import React, { act } from "react";
import type { Root } from "react-dom/client";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

// ---------------------------------------------------------------------------
// Hoisted mock functions (must be defined before vi.mock calls)
// ---------------------------------------------------------------------------

const { mockGetToken } = vi.hoisted(() => ({
	mockGetToken: vi.fn<() => string | null>(),
}));

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

vi.mock("@/utils/tokenStorage", () => ({
	getToken: mockGetToken,
	saveToken: vi.fn(),
	removeToken: vi.fn(),
}));

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { AuthProvider } from "@/features/auth/contexts/AuthContext";
import { useAuth } from "./useAuth";

// ---------------------------------------------------------------------------
// Test components
//
// Defined at module scope to satisfy unicorn/consistent-function-scoping.
// Render auth state into data attributes so tests can assert via DOM queries
// without mutating outer-scope variables.
// ---------------------------------------------------------------------------

function AuthConsumer(): React.ReactElement {
	const { token, isAuthenticated, login, logout } = useAuth();
	return React.createElement(
		"span",
		{
			"data-token": token ?? "null",
			"data-is-authenticated": String(isAuthenticated),
			"data-has-login": String(typeof login === "function"),
			"data-has-logout": String(typeof logout === "function"),
		},
		"ok",
	);
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

let root: Root | null = null;
let container: HTMLDivElement | null = null;

function clearBody(): void {
	while (document.body.firstChild) {
		document.body.firstChild.remove();
	}
}

async function mountWithProvider(): Promise<HTMLDivElement> {
	const div = document.createElement("div");
	document.body.append(div);
	container = div;
	await act(async () => {
		root = createRoot(div);
		root.render(
			React.createElement(
				AuthProvider,
				null,
				React.createElement(AuthConsumer),
			),
		);
	});
	return div;
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
	clearBody();
});

// ---------------------------------------------------------------------------
// Setup
// ---------------------------------------------------------------------------

beforeEach(() => {
	mockGetToken.mockReset();
});

// ---------------------------------------------------------------------------
// Tests — happy path (inside AuthProvider)
// ---------------------------------------------------------------------------

describe("useAuth — inside AuthProvider", () => {
	it("returns the AuthContext value when called inside AuthProvider", async () => {
		mockGetToken.mockReturnValue("test-token");

		const div = await mountWithProvider();
		const span = div.querySelector("span");

		expect(span).not.toBeNull();
		expect(span?.dataset.token).toBe("test-token");
		expect(span?.dataset.isAuthenticated).toBe("true");
		expect(span?.dataset.hasLogin).toBe("true");
		expect(span?.dataset.hasLogout).toBe("true");
	});

	it("returns isAuthenticated false when token is null in AuthProvider", async () => {
		mockGetToken.mockReturnValue(null);

		const div = await mountWithProvider();
		const span = div.querySelector("span");

		expect(span?.dataset.isAuthenticated).toBe("false");
		expect(span?.dataset.token).toBe("null");
	});
});

// ---------------------------------------------------------------------------
// Tests — error path (outside AuthProvider)
// ---------------------------------------------------------------------------

describe("useAuth — outside AuthProvider", () => {
	it("throws 'useAuth must be used within AuthProvider' when called without a provider", async () => {
		// `onError` is a spy so BadComponent can report the caught error without
		// mutating a closed-over variable (react-hooks/immutability).
		const onError = vi.fn<(error: Error) => void>();

		function BadComponent(): React.ReactElement {
			try {
				// biome-ignore lint/correctness/useHookAtTopLevel: intentionally testing hook called without a provider
				useAuth();
			} catch (error) {
				if (error instanceof Error) {
					onError(error);
				}
			}
			return React.createElement("span", null, "bad");
		}

		// Suppress React's error-boundary console.error noise during this test
		// biome-ignore lint/suspicious/noConsole: silencing React render errors in tests
		const originalConsoleError = console.error;
		console.error = vi.fn();

		const div = document.createElement("div");
		document.body.append(div);
		container = div;

		await act(async () => {
			root = createRoot(div);
			root.render(React.createElement(BadComponent));
		});

		console.error = originalConsoleError;

		expect(onError).toHaveBeenCalledOnce();
		expect(onError.mock.calls[0][0]).toBeInstanceOf(Error);
		expect(onError.mock.calls[0][0].message).toBe(
			"useAuth must be used within AuthProvider",
		);
	});
});
