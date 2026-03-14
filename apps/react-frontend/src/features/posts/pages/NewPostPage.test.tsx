/**
 * Tests for NewPostPage component.
 *
 * Uses ReactDOM + manual DOM querying since @testing-library/react is not installed.
 *
 * Strategy: mock useCreatePostForm entirely so tests verify only that NewPostPage
 * correctly wires up the hook — not the hook's own logic (which is tested
 * separately in useCreatePostForm.test.ts).
 *
 * Mocks:
 *  - @/features/posts/hooks/useCreatePostForm  useCreatePostForm vi.mock
 */

import type React from "react"
import { act } from "react"
import { createRoot } from "react-dom/client"
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"

// ---------------------------------------------------------------------------
// Hoisted mock factories
// ---------------------------------------------------------------------------

const { mockReset, makeFormObject } = vi.hoisted(() => {
	const mockReset = vi.fn()

	function makeFormObject(isSubmitting: boolean) {
		return {
			// handleSubmit wraps the caller-supplied onSubmit so that form submission
			// is synchronous and testable without needing real RHF internals.
			handleSubmit:
				(fn: (values: { content: string }) => void) =>
				(e: React.FormEvent) => {
					e.preventDefault()
					fn({ content: "" })
				},
			control: {},
			formState: { isSubmitting },
			register: () => ({}),
			reset: mockReset,
		}
	}

	return { mockReset, makeFormObject }
})

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

// Default: not submitting, charCount starts at 0.
vi.mock("@/features/posts/hooks/useCreatePostForm", () => ({
	useCreatePostForm: vi.fn(() => ({
		form: makeFormObject(false),
		onSubmit: vi.fn(),
		charCount: 0,
	})),
}))

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { useCreatePostForm } from "@/features/posts/hooks/useCreatePostForm"
import { NewPostPage } from "./NewPostPage"

// ---------------------------------------------------------------------------
// Typed mock reference
// ---------------------------------------------------------------------------

const mockUseCreatePostForm = useCreatePostForm as ReturnType<typeof vi.fn>

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function mountNewPostPage(): HTMLDivElement {
	const container = document.createElement("div")
	document.body.appendChild(container)
	act(() => {
		createRoot(container).render(<NewPostPage />)
	})
	return container
}

async function submitForm(form: HTMLFormElement): Promise<void> {
	await act(async () => {
		form.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true }))
		await Promise.resolve()
		await Promise.resolve()
	})
}

function clearBody(): void {
	while (document.body.firstChild) {
		document.body.removeChild(document.body.firstChild)
	}
}

// ---------------------------------------------------------------------------
// Setup / teardown
// ---------------------------------------------------------------------------

beforeEach(() => {
	mockReset.mockReset()
	// Restore the default mock return value before each test.
	mockUseCreatePostForm.mockReturnValue({
		form: makeFormObject(false),
		onSubmit: vi.fn(),
		charCount: 0,
	})
})

afterEach(() => {
	clearBody()
})

// ---------------------------------------------------------------------------
// Tests — rendering
// ---------------------------------------------------------------------------

describe("NewPostPage — rendering", () => {
	it("renders a heading containing 'New Post'", () => {
		const container = mountNewPostPage()
		expect(container.textContent).toContain("New Post")
	})

	it("renders a form element", () => {
		const container = mountNewPostPage()
		const form = container.querySelector<HTMLFormElement>("form")
		expect(form).not.toBeNull()
	})

	it("renders a textarea for the content field", () => {
		const container = mountNewPostPage()
		const textarea =
			container.querySelector<HTMLTextAreaElement>("textarea") ??
			container.querySelector<HTMLInputElement>(
				'input[name="content"], input[placeholder*="content" i]',
			)
		expect(textarea).not.toBeNull()
	})

	it("renders a submit button", () => {
		const container = mountNewPostPage()
		const submitBtn =
			container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
			container.querySelector<HTMLButtonElement>("button")
		expect(submitBtn).not.toBeNull()
	})

	it("shows character count '0 / 280' initially", () => {
		const container = mountNewPostPage()
		expect(container.textContent).toContain("0 / 280")
	})
})

// ---------------------------------------------------------------------------
// Tests — submit button disabled state
// ---------------------------------------------------------------------------

describe("NewPostPage — submit button disabled state", () => {
	it("disables the submit button when isSubmitting is true", () => {
		mockUseCreatePostForm.mockReturnValue({
			form: makeFormObject(true),
			onSubmit: vi.fn(),
			charCount: 0,
		})

		const container = mountNewPostPage()
		const submitBtn =
			container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
			container.querySelector<HTMLButtonElement>("button")
		expect(submitBtn).not.toBeNull()
		expect(submitBtn?.disabled).toBe(true)
	})

	it("enables the submit button when isSubmitting is false", () => {
		mockUseCreatePostForm.mockReturnValue({
			form: makeFormObject(false),
			onSubmit: vi.fn(),
			charCount: 0,
		})

		const container = mountNewPostPage()
		const submitBtn =
			container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
			container.querySelector<HTMLButtonElement>("button")
		expect(submitBtn).not.toBeNull()
		expect(submitBtn?.disabled).toBe(false)
	})
})

// ---------------------------------------------------------------------------
// Tests — form submission wiring
// ---------------------------------------------------------------------------

describe("NewPostPage — form submission wiring", () => {
	it("calls form.handleSubmit(onSubmit) when the form is submitted", async () => {
		// Provide a spy for handleSubmit so we can assert it was invoked.
		const handleSubmitSpy = vi.fn(
			(fn: (values: { content: string }) => void) =>
				(e: React.FormEvent) => {
					e.preventDefault()
					fn({ content: "" })
				},
		)
		const onSubmitSpy = vi.fn()

		mockUseCreatePostForm.mockReturnValue({
			form: {
				...makeFormObject(false),
				handleSubmit: handleSubmitSpy,
			},
			onSubmit: onSubmitSpy,
			charCount: 0,
		})

		const container = mountNewPostPage()
		const form = container.querySelector<HTMLFormElement>("form")
		expect(form).not.toBeNull()

		await submitForm(form as HTMLFormElement)

		// handleSubmit must have been called with the page's onSubmit handler.
		expect(handleSubmitSpy).toHaveBeenCalledWith(onSubmitSpy)
	})
})
