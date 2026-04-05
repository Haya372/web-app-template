/**
 * Tests for _auth layout route.
 *
 * Asserts that the auth shell renders ONLY the router Outlet placeholder —
 * Header and Footer must NOT appear in this layout.
 *
 * Mocks:
 *  - @tanstack/react-router   Outlet vi.mock
 */

import React, { act } from "react";
import { createRoot } from "react-dom/client";
import type { Root } from "react-dom/client";
import { afterEach, describe, expect, it, vi } from "vite-plus/test";

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

function routeOptions(options: { component: unknown }) {
  return options;
}

vi.mock("@tanstack/react-router", () => ({
  Outlet: () => React.createElement("div", { "data-testid": "outlet" }),
  createFileRoute: () => routeOptions,
}));

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { AuthLayout } from "./_auth";

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
    root.render(<AuthLayout />);
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

describe("AuthLayout — rendering", () => {
  it("renders the router Outlet placeholder", async () => {
    await mount();
    expect(document.querySelector("[data-testid='outlet']")).not.toBeNull();
  });

  it("does NOT render the Header component", async () => {
    await mount();
    expect(document.querySelector("[data-testid='header']")).toBeNull();
  });

  it("does NOT render the Footer component", async () => {
    await mount();
    expect(document.querySelector("[data-testid='footer']")).toBeNull();
  });
});
