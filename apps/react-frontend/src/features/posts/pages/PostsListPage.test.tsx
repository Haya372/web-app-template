/**
 * Tests for PostsListPage component.
 *
 * Strategy: mock usePostsList entirely so tests verify only that PostsListPage
 * correctly wires up the hook — not the hook's own logic (which is tested
 * separately in usePostsList.test.ts).
 *
 * Mocks:
 *  - @/features/posts/hooks/usePostsList  usePostsList vi.mock
 */

import { act } from "react"
import { createRoot } from "react-dom/client"
import type { Root } from "react-dom/client"
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"
import type { PostListResponse } from "@/generated/types.gen"

// ---------------------------------------------------------------------------
// Hoisted mock factories
// ---------------------------------------------------------------------------

const { mockUsePostsList } = vi.hoisted(() => {
	const mockUsePostsList = vi.fn()

	return { mockUsePostsList }
})

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

vi.mock("@/features/posts/hooks/usePostsList", () => ({
	usePostsList: mockUsePostsList,
}))

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { PostsListPage } from "./PostsListPage"

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

interface MockHookState {
	posts: PostListResponse["posts"] | undefined
	total: number | undefined
	isLoading: boolean
	isError: boolean
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

let root: Root | undefined
let container: HTMLDivElement | undefined

function mountPostsListPage(): HTMLDivElement {
	container = document.createElement("div")
	document.body.append(container)
	act(() => {
		root = createRoot(container as HTMLDivElement)
		root.render(<PostsListPage />)
	})

	return container
}

function setMockState(state: MockHookState): void {
	mockUsePostsList.mockReturnValue(state)
}

// ---------------------------------------------------------------------------
// Setup / teardown
// ---------------------------------------------------------------------------

beforeEach(() => {
	setMockState({ posts: [], total: 0, isLoading: false, isError: false })
})

afterEach(() => {
	if (root) {
		act(() => { root?.unmount() })
		root = undefined
	}
	container?.remove()
	container = undefined
})

// ---------------------------------------------------------------------------
// Tests — rendering
// ---------------------------------------------------------------------------

describe("PostsListPage — rendering", () => {
	it("renders a heading containing 'Posts'", () => {
		const el = mountPostsListPage()
		expect(el.textContent).toContain("Posts")
	})

	it("shows empty state message when posts is an empty array", () => {
		setMockState({ posts: [], total: 0, isLoading: false, isError: false })
		const el = mountPostsListPage()
		expect(el.textContent).toContain("No posts yet")
	})

	it("renders post content when posts are returned", () => {
		setMockState({
			posts: [
				{
					id: "00000000-0000-0000-0000-000000000001",
					userId: "00000000-0000-0000-0000-000000000002",
					content: "Hello from a post",
					createdAt: "2026-01-01T10:00:00Z",
				},
			],
			total: 1,
			isLoading: false,
			isError: false,
		})
		const el = mountPostsListPage()
		expect(el.textContent).toContain("Hello from a post")
	})

	it("renders multiple posts", () => {
		setMockState({
			posts: [
				{
					id: "00000000-0000-0000-0000-000000000001",
					userId: "00000000-0000-0000-0000-000000000002",
					content: "First post",
					createdAt: "2026-01-02T10:00:00Z",
				},
				{
					id: "00000000-0000-0000-0000-000000000003",
					userId: "00000000-0000-0000-0000-000000000002",
					content: "Second post",
					createdAt: "2026-01-01T10:00:00Z",
				},
			],
			total: 2,
			isLoading: false,
			isError: false,
		})
		const el = mountPostsListPage()
		expect(el.textContent).toContain("First post")
		expect(el.textContent).toContain("Second post")
	})
})

// ---------------------------------------------------------------------------
// Tests — loading state
// ---------------------------------------------------------------------------

describe("PostsListPage — loading state", () => {
	it("renders a loading status element while loading", () => {
		setMockState({ posts: undefined, total: undefined, isLoading: true, isError: false })
		const el = mountPostsListPage()
		const status = el.querySelector("[role='status']")
		expect(status).not.toBeNull()
		expect(el.textContent).toContain("Loading")
	})

	it("does not show empty state message while loading", () => {
		setMockState({ posts: undefined, total: undefined, isLoading: true, isError: false })
		const el = mountPostsListPage()
		expect(el.textContent).not.toContain("No posts yet")
	})
})

// ---------------------------------------------------------------------------
// Tests — error state
// ---------------------------------------------------------------------------

describe("PostsListPage — error state", () => {
	it("renders an alert element when isError is true", () => {
		setMockState({ posts: undefined, total: undefined, isLoading: false, isError: true })
		const el = mountPostsListPage()
		const alert = el.querySelector("[role='alert']")
		expect(alert).not.toBeNull()
	})

	it("renders error message text when isError is true", () => {
		setMockState({ posts: undefined, total: undefined, isLoading: false, isError: true })
		const el = mountPostsListPage()
		expect(el.textContent).toContain("Failed to load posts. Please try again.")
	})

	it("does not show post list or empty state when isError is true", () => {
		setMockState({ posts: undefined, total: undefined, isLoading: false, isError: true })
		const el = mountPostsListPage()
		expect(el.textContent).not.toContain("No posts yet")
		expect(el.querySelector("ul")).toBeNull()
	})
})
