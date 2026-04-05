/**
 * Tests for usePostsList hook.
 *
 * Strategy: mock getV1Posts and getToken so tests verify that usePostsList
 * correctly delegates to TanStack Query with the right fetch function,
 * token header, and options.
 *
 * Hook state is observed by rendering a test component that exposes values
 * as data attributes on a DOM node.
 */

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, createElement } from "react";
import { createRoot } from "react-dom/client";
import type { Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vite-plus/test";
import type { PostListResponse } from "@/generated/types.gen";

// ---------------------------------------------------------------------------
// Hoisted mock factories
// ---------------------------------------------------------------------------

const { mockGetV1Posts, mockGetToken, mockNavigate } = vi.hoisted(() => {
  const mockGetV1Posts = vi.fn();
  const mockGetToken = vi.fn();
  const mockNavigate = vi.fn();

  return { mockGetV1Posts, mockGetToken, mockNavigate };
});

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

vi.mock("@/generated/sdk.gen", () => ({
  getV1Posts: mockGetV1Posts,
}));

vi.mock("@/utils/tokenStorage", () => ({
  getToken: mockGetToken,
}));

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => mockNavigate,
}));

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { usePostsList } from "./usePostsList";

// ---------------------------------------------------------------------------
// Helpers — test component renders hook state into the DOM
// ---------------------------------------------------------------------------

function TestComponent() {
  const { posts, total, isLoading, isError } = usePostsList();

  return createElement(
    "div",
    { "data-testid": "hook-output" },
    createElement("span", { "data-loading": String(isLoading) }),
    createElement("span", { "data-error": String(isError) }),
    createElement("span", { "data-total": total !== undefined ? String(total) : "" }),
    createElement(
      "ul",
      null,
      ...(posts ?? []).map((p) =>
        createElement("li", { key: p.id, "data-content": p.content }, p.content),
      ),
    ),
  );
}

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
}

function mountHook(queryClient: QueryClient): { container: HTMLDivElement; root: Root } {
  const container = document.createElement("div");
  document.body.append(container);

  let root!: Root;
  act(() => {
    root = createRoot(container);
    root.render(
      createElement(QueryClientProvider, { client: queryClient }, createElement(TestComponent)),
    );
  });

  return { container, root };
}

function getOutput(container: HTMLDivElement) {
  const el = container.querySelector("[data-testid='hook-output']");
  const isLoading =
    (el?.querySelector("[data-loading]") as HTMLElement | null)?.dataset.loading === "true";
  const isError =
    (el?.querySelector("[data-error]") as HTMLElement | null)?.dataset.error === "true";
  const totalAttr = (el?.querySelector("[data-total]") as HTMLElement | null)?.dataset.total;
  const total = totalAttr ? Number(totalAttr) : undefined;
  const posts = [...(el?.querySelectorAll("[data-content]") ?? [])].map(
    (node) => (node as HTMLElement).dataset.content ?? "",
  );

  return { isLoading, isError, total, posts };
}

// ---------------------------------------------------------------------------
// Setup / teardown
// ---------------------------------------------------------------------------

let cleanup: (() => void) | undefined;

beforeEach(() => {
  mockGetV1Posts.mockReset();
  mockGetToken.mockReset();
  mockNavigate.mockReset();
  cleanup = undefined;
});

afterEach(() => {
  cleanup?.();
  while (document.body.firstChild) {
    document.body.firstChild.remove();
  }
});

// ---------------------------------------------------------------------------
// Tests — happy path
// ---------------------------------------------------------------------------

describe("usePostsList — happy path", () => {
  it("returns posts and total when fetch succeeds", async () => {
    mockGetToken.mockReturnValue("test-token");
    const mockData: PostListResponse = {
      posts: [
        {
          id: "00000000-0000-0000-0000-000000000001",
          userId: "00000000-0000-0000-0000-000000000002",
          content: "Hello world",
          createdAt: "2026-01-01T10:00:00Z",
        },
      ],
      total: 1,
      limit: 20,
      offset: 0,
    };
    mockGetV1Posts.mockResolvedValue({
      data: mockData,
      error: undefined,
      response: { status: 200 },
    });

    const queryClient = makeQueryClient();
    const { container, root } = mountHook(queryClient);
    cleanup = () => {
      act(() => {
        root.unmount();
      });
      container.remove();
    };

    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: ["posts", "list"] });
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    const output = getOutput(container);
    expect(output.posts).toHaveLength(1);
    expect(output.posts[0]).toBe("Hello world");
    expect(output.total).toBe(1);
    expect(output.isLoading).toBe(false);
    expect(output.isError).toBe(false);
  });

  it("returns empty posts array when response is empty", async () => {
    mockGetToken.mockReturnValue("test-token");
    const mockData: PostListResponse = { posts: [], total: 0, limit: 20, offset: 0 };
    mockGetV1Posts.mockResolvedValue({
      data: mockData,
      error: undefined,
      response: { status: 200 },
    });

    const queryClient = makeQueryClient();
    const { container, root } = mountHook(queryClient);
    cleanup = () => {
      act(() => {
        root.unmount();
      });
      container.remove();
    };

    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: ["posts", "list"] });
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    const output = getOutput(container);
    expect(output.posts).toHaveLength(0);
    expect(output.total).toBe(0);
  });
});

// ---------------------------------------------------------------------------
// Tests — error path
// ---------------------------------------------------------------------------

describe("usePostsList — error path", () => {
  it("sets isError to true when API returns an error", async () => {
    mockGetToken.mockReturnValue("test-token");
    mockGetV1Posts.mockResolvedValue({
      data: undefined,
      error: { type: "INTERNAL_SERVER_ERROR" },
      response: { status: 500 },
    });

    const queryClient = makeQueryClient();
    const { container, root } = mountHook(queryClient);
    cleanup = () => {
      act(() => {
        root.unmount();
      });
      container.remove();
    };

    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: ["posts", "list"] });
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    const output = getOutput(container);
    expect(output.isError).toBe(true);
    expect(output.posts).toHaveLength(0);
  });

  it("sets isError to true when fetch throws a network error", async () => {
    mockGetToken.mockReturnValue("test-token");
    mockGetV1Posts.mockRejectedValue(new Error("Network error"));

    const queryClient = makeQueryClient();
    const { container, root } = mountHook(queryClient);
    cleanup = () => {
      act(() => {
        root.unmount();
      });
      container.remove();
    };

    await act(async () => {
      await queryClient.invalidateQueries({ queryKey: ["posts", "list"] });
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    const output = getOutput(container);
    expect(output.isError).toBe(true);
  });
});

// ---------------------------------------------------------------------------
// Tests — API integration (token header)
// ---------------------------------------------------------------------------

describe("usePostsList — API integration", () => {
  it("calls getV1Posts with a Bearer token header when token is present", async () => {
    mockGetToken.mockReturnValue("my-secret-token");
    const mockData: PostListResponse = { posts: [], total: 0, limit: 20, offset: 0 };
    mockGetV1Posts.mockResolvedValue({
      data: mockData,
      error: undefined,
      response: { status: 200 },
    });

    const queryClient = makeQueryClient();
    const { container, root } = mountHook(queryClient);
    cleanup = () => {
      act(() => {
        root.unmount();
      });
      container.remove();
    };

    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    expect(mockGetV1Posts).toHaveBeenCalledOnce();
    expect(mockGetV1Posts).toHaveBeenCalledWith(
      expect.objectContaining({
        headers: { Authorization: "Bearer my-secret-token" },
      }),
    );
  });

  it("does not call getV1Posts when token is null (query disabled)", async () => {
    mockGetToken.mockReturnValue(null);

    const queryClient = makeQueryClient();
    const { container, root } = mountHook(queryClient);
    cleanup = () => {
      act(() => {
        root.unmount();
      });
      container.remove();
    };

    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    expect(mockGetV1Posts).not.toHaveBeenCalled();
  });

  it("redirects to login and does not call getV1Posts when token is null", async () => {
    mockGetToken.mockReturnValue(null);

    const queryClient = makeQueryClient();
    const { container, root } = mountHook(queryClient);
    cleanup = () => {
      act(() => {
        root.unmount();
      });
      container.remove();
    };

    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    expect(mockNavigate).toHaveBeenCalledWith({ to: "/login" });
    expect(mockGetV1Posts).not.toHaveBeenCalled();
  });
});
