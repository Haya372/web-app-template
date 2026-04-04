/**
 * Tests for AuthProvider and AuthContext.
 *
 * AuthProvider initialises token from localStorage, exposes login/logout
 * actions, and keeps isAuthenticated in sync with the token state.
 *
 * Mocks:
 *  - @/utils/tokenStorage   getToken / saveToken / removeToken   vi.mock
 */

import React, { act, useContext } from "react";
import type { Root } from "react-dom/client";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

// ---------------------------------------------------------------------------
// Hoisted mock functions (must be defined before vi.mock calls)
// ---------------------------------------------------------------------------

const { mockGetToken, mockSaveToken, mockRemoveToken } = vi.hoisted(() => ({
	mockGetToken: vi.fn<() => string | null>(),
	mockSaveToken: vi.fn<(token: string) => void>(),
	mockRemoveToken: vi.fn<() => void>(),
}));

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

vi.mock("@/utils/tokenStorage", () => ({
	getToken: mockGetToken,
	saveToken: mockSaveToken,
	removeToken: mockRemoveToken,
}));

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { AuthContext, AuthProvider } from "./AuthContext";

// ---------------------------------------------------------------------------
// BareConsumer — renders context presence into a data attribute.
// Defined at module scope to satisfy unicorn/consistent-function-scoping.
// ---------------------------------------------------------------------------

function BareConsumer(): React.ReactElement {
	const ctx = useContext(AuthContext);
	return React.createElement(
		"span",
		{ "data-is-null": String(ctx === null) },
		"bare",
	);
}

// ---------------------------------------------------------------------------
// TestConsumer
//
// Reads the AuthContext value and projects every field into the DOM so tests
// can observe them without @testing-library/react.
//
// DOM structure:
//   <div>
//     <span data-testid="token">{token ?? "null"}</span>
//     <span data-testid="is-authenticated">{String(isAuthenticated)}</span>
//     <button data-testid="login-btn" onClick={() => login("new-token")} />
//     <button data-testid="logout-btn" onClick={() => logout()} />
//   </div>
// ---------------------------------------------------------------------------

function TestConsumer(): React.ReactElement {
	const ctx = useContext(AuthContext);
	if (ctx === null) {
		return React.createElement(
			"div",
			{ "data-testid": "no-context" },
			"no context",
		);
	}
	const { token, isAuthenticated, login, logout } = ctx;
	return React.createElement(
		"div",
		null,
		React.createElement("span", { "data-testid": "token" }, token ?? "null"),
		React.createElement(
			"span",
			{ "data-testid": "is-authenticated" },
			String(isAuthenticated),
		),
		React.createElement(
			"button",
			{
				"data-testid": "login-btn",
				onClick: () => login("new-token"),
				type: "button",
			},
			"Login",
		),
		React.createElement(
			"button",
			{
				"data-testid": "logout-btn",
				onClick: () => logout(),
				type: "button",
			},
			"Logout",
		),
	);
}

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
		root.render(
			React.createElement(
				AuthProvider,
				null,
				React.createElement(TestConsumer),
			),
		);
	});
}

function getSpanText(testId: string): string | null {
	return (
		container?.querySelector<HTMLElement>(`[data-testid="${testId}"]`)
			?.textContent ?? null
	);
}

function clickButton(testId: string): Promise<void> {
	return act(async () => {
		const btn = container?.querySelector<HTMLButtonElement>(
			`[data-testid="${testId}"]`,
		);
		btn?.click();
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
// Setup
// ---------------------------------------------------------------------------

beforeEach(() => {
	mockGetToken.mockReset();
	mockSaveToken.mockReset();
	mockRemoveToken.mockReset();
});

// ---------------------------------------------------------------------------
// Tests — initial state
// ---------------------------------------------------------------------------

describe("AuthProvider — initial state", () => {
	it("initializes token from localStorage when one exists", async () => {
		mockGetToken.mockReturnValue("existing-token");

		await mount();

		expect(getSpanText("token")).toBe("existing-token");
	});

	it("initializes token as null when localStorage is empty", async () => {
		mockGetToken.mockReturnValue(null);

		await mount();

		expect(getSpanText("token")).toBe("null");
	});

	it("sets isAuthenticated to true when a token exists in localStorage", async () => {
		mockGetToken.mockReturnValue("existing-token");

		await mount();

		expect(getSpanText("is-authenticated")).toBe("true");
	});

	it("sets isAuthenticated to false when localStorage is empty", async () => {
		mockGetToken.mockReturnValue(null);

		await mount();

		expect(getSpanText("is-authenticated")).toBe("false");
	});
});

// ---------------------------------------------------------------------------
// Tests — login action
// ---------------------------------------------------------------------------

describe("AuthProvider — login", () => {
	it("updates the token state after calling login", async () => {
		mockGetToken.mockReturnValue(null);

		await mount();
		await clickButton("login-btn");

		expect(getSpanText("token")).toBe("new-token");
	});

	it("calls saveToken with the provided token", async () => {
		mockGetToken.mockReturnValue(null);

		await mount();
		await clickButton("login-btn");

		expect(mockSaveToken).toHaveBeenCalledWith("new-token");
	});

	it("sets isAuthenticated to true after login", async () => {
		mockGetToken.mockReturnValue(null);

		await mount();
		await clickButton("login-btn");

		expect(getSpanText("is-authenticated")).toBe("true");
	});
});

// ---------------------------------------------------------------------------
// Tests — logout action
// ---------------------------------------------------------------------------

describe("AuthProvider — logout", () => {
	it("clears the token state after calling logout", async () => {
		mockGetToken.mockReturnValue("existing-token");

		await mount();
		await clickButton("logout-btn");

		expect(getSpanText("token")).toBe("null");
	});

	it("calls removeToken when logout is invoked", async () => {
		mockGetToken.mockReturnValue("existing-token");

		await mount();
		await clickButton("logout-btn");

		expect(mockRemoveToken).toHaveBeenCalledTimes(1);
	});

	it("sets isAuthenticated to false after logout", async () => {
		mockGetToken.mockReturnValue("existing-token");

		await mount();
		await clickButton("logout-btn");

		expect(getSpanText("is-authenticated")).toBe("false");
	});
});

// ---------------------------------------------------------------------------
// Tests — isAuthenticated derived from token
// ---------------------------------------------------------------------------

describe("AuthProvider — isAuthenticated reflects token", () => {
	it("is false initially, becomes true after login, then false after logout", async () => {
		mockGetToken.mockReturnValue(null);

		await mount();
		expect(getSpanText("is-authenticated")).toBe("false");

		await clickButton("login-btn");
		expect(getSpanText("is-authenticated")).toBe("true");

		await clickButton("logout-btn");
		expect(getSpanText("is-authenticated")).toBe("false");
	});
});

// ---------------------------------------------------------------------------
// Tests — context value shape (type assertion without rendering)
// ---------------------------------------------------------------------------

describe("AuthContext — default value", () => {
	it("has a null default value outside of AuthProvider", async () => {
		// The default value passed to createContext is null; verify via a consumer
		// rendered without wrapping AuthProvider.
		const div = document.createElement("div");
		document.body.append(div);
		await act(async () => {
			createRoot(div).render(React.createElement(BareConsumer));
		});

		expect(div.querySelector("[data-is-null='true']")).not.toBeNull();
		div.remove();
	});
});
