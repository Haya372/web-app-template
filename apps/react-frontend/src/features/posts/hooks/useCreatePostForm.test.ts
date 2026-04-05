/**
 * Tests for useCreatePostForm hook.
 *
 * Since @testing-library/react is not installed, the hook is exercised via a
 * minimal TestComponent rendered with ReactDOM.createRoot + act.  The
 * TestComponent renders a <form> that wires up the hook's `form` and
 * `onSubmit`, and a <span> that displays `charCount`, so every observable
 * return value is reachable through the DOM.
 *
 * Mocks:
 *  - @/generated/sdk.gen               postV1Posts    vi.mock
 *  - @/utils/tokenStorage              getToken       vi.mock
 *  - @tanstack/react-router            useNavigate    vi.mock
 *  - @repo/ui                          toast          vi.mock
 */

import { createElement, act } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vite-plus/test";

// ---------------------------------------------------------------------------
// Hoisted mock functions (must be defined before vi.mock calls)
// ---------------------------------------------------------------------------

const { mockNavigate, mockToastSuccess, mockToastError, mockGetToken } = vi.hoisted(() => ({
  mockNavigate: vi.fn(),
  mockToastSuccess: vi.fn(),
  mockToastError: vi.fn(),
  mockGetToken: vi.fn<() => string | null>(),
}));

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

vi.mock("@/generated/sdk.gen", () => ({
  postV1Posts: vi.fn(),
}));

vi.mock("@/utils/tokenStorage", () => ({
  getToken: mockGetToken,
}));

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => mockNavigate,
}));

vi.mock("@repo/ui", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@repo/ui")>();
  return {
    ...actual,
    toast: {
      ...actual.toast,
      success: mockToastSuccess,
      error: mockToastError,
    },
  };
});

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { postV1Posts } from "@/generated/sdk.gen";
import { useCreatePostForm } from "@/features/posts/hooks/useCreatePostForm";

// ---------------------------------------------------------------------------
// Typed mock references
// ---------------------------------------------------------------------------

const mockPostV1Posts = postV1Posts as ReturnType<typeof vi.fn>;

// ---------------------------------------------------------------------------
// TestComponent
//
// Renders the hook's return values into the DOM so tests can interact with
// them without @testing-library/react.
//
// DOM structure:
//   <form>
//     <textarea name="content" />
//     <span data-testid="char-count">{charCount}</span>
//     <button type="submit" disabled={isSubmitting}>Submit</button>
//     {errors.content && <p data-testid="content-error">{message}</p>}
//   </form>
// ---------------------------------------------------------------------------

function TestComponent(): React.ReactElement {
  const { form, onSubmit, charCount } = useCreatePostForm();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = form;

  return createElement(
    "form",
    { onSubmit: handleSubmit(onSubmit) },
    createElement("textarea", {
      ...register("content"),
    }),
    createElement("span", { "data-testid": "char-count" }, String(charCount)),
    errors.content &&
      createElement("p", { "data-testid": "content-error" }, errors.content.message),
    createElement("button", { type: "submit", disabled: isSubmitting }, "Submit"),
  );
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function mountTestComponent(): HTMLDivElement {
  const container = document.createElement("div");
  document.body.append(container);
  act(() => {
    createRoot(container).render(createElement(TestComponent));
  });
  return container;
}

async function submitForm(form: HTMLFormElement): Promise<void> {
  await act(async () => {
    form.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true }));
    // Allow react-hook-form's async validation + submit handler to resolve
    await Promise.resolve();
    await Promise.resolve();
  });
}

function fillTextarea(textarea: HTMLTextAreaElement, value: string): void {
  act(() => {
    const nativeValueSetter = Object.getOwnPropertyDescriptor(
      HTMLTextAreaElement.prototype,
      "value",
    )?.set;
    nativeValueSetter?.call(textarea, value);
    textarea.dispatchEvent(new Event("input", { bubbles: true }));
    textarea.dispatchEvent(new Event("change", { bubbles: true }));
  });
}

function clearBody(): void {
  while (document.body.firstChild) {
    document.body.firstChild.remove();
  }
}

// ---------------------------------------------------------------------------
// Setup / teardown
// ---------------------------------------------------------------------------

beforeEach(() => {
  vi.stubEnv("VITE_API_BASE_URL", "http://localhost:8080");
  mockPostV1Posts.mockReset();
  mockNavigate.mockReset();
  mockToastSuccess.mockReset();
  mockToastError.mockReset();
  mockGetToken.mockReset();
});

afterEach(() => {
  clearBody();
  vi.unstubAllEnvs();
});

// ---------------------------------------------------------------------------
// Tests — rendering / auth redirect
// ---------------------------------------------------------------------------

describe("useCreatePostForm — auth redirect on mount", () => {
  it("redirects to /login when no token is present on mount", () => {
    mockGetToken.mockReturnValue(null);

    mountTestComponent();

    expect(mockNavigate).toHaveBeenCalledWith({ to: "/login" });
  });

  it("does not redirect when a token is present on mount", () => {
    mockGetToken.mockReturnValue("valid-jwt-token");

    mountTestComponent();

    expect(mockNavigate).not.toHaveBeenCalled();
  });
});

// ---------------------------------------------------------------------------
// Tests — client-side validation
// ---------------------------------------------------------------------------

describe("useCreatePostForm — client-side validation", () => {
  it("shows an error when content is empty on submit", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");

    const container = mountTestComponent();
    const form = container.querySelector<HTMLFormElement>("form");
    expect(form).not.toBeNull();

    await submitForm(form as HTMLFormElement);

    // Validation must surface an error — postV1Posts must not be called
    expect(mockPostV1Posts).not.toHaveBeenCalled();
    // The error message should be visible in the DOM
    const errorEl = container.querySelector<HTMLElement>("[data-testid='content-error']");
    expect(errorEl).not.toBeNull();
  });

  it("shows 'Content must be 280 characters or less' when content exceeds 280 characters", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");
    expect(textarea).not.toBeNull();
    expect(form).not.toBeNull();

    // 281 characters — one over the limit
    fillTextarea(textarea as HTMLTextAreaElement, "a".repeat(281));
    await submitForm(form as HTMLFormElement);

    expect(mockPostV1Posts).not.toHaveBeenCalled();
    expect(container.textContent).toContain("Content must be 280 characters or less");
  });

  it("does not show a validation error when content is exactly 280 characters", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    mockPostV1Posts.mockResolvedValue({
      data: { id: "post-1", content: "a".repeat(280), createdAt: "2026-03-14T00:00:00Z" },
      error: undefined,
      response: { status: 201 },
    });

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");
    expect(textarea).not.toBeNull();
    expect(form).not.toBeNull();

    fillTextarea(textarea as HTMLTextAreaElement, "a".repeat(280));
    await submitForm(form as HTMLFormElement);

    expect(container.textContent).not.toContain("Content must be 280 characters or less");
  });
});

// ---------------------------------------------------------------------------
// Tests — charCount
// ---------------------------------------------------------------------------

describe("useCreatePostForm — charCount", () => {
  it("reflects the current character count of the content field", () => {
    mockGetToken.mockReturnValue("valid-jwt-token");

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    expect(textarea).not.toBeNull();

    fillTextarea(textarea as HTMLTextAreaElement, "Hello!");

    const charCountEl = container.querySelector<HTMLElement>("[data-testid='char-count']");
    expect(charCountEl?.textContent).toBe("6");
  });

  it("starts at 0 before any input", () => {
    mockGetToken.mockReturnValue("valid-jwt-token");

    const container = mountTestComponent();

    const charCountEl = container.querySelector<HTMLElement>("[data-testid='char-count']");
    expect(charCountEl?.textContent).toBe("0");
  });
});

// ---------------------------------------------------------------------------
// Tests — happy path submit
// ---------------------------------------------------------------------------

describe("useCreatePostForm — successful submit", () => {
  const VALID_CONTENT = "Hello, world!";
  const MOCK_RESPONSE = {
    id: "post-abc",
    content: VALID_CONTENT,
    createdAt: "2026-03-14T00:00:00Z",
  };

  it("calls postV1Posts with the entered content", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    mockPostV1Posts.mockResolvedValue({
      data: MOCK_RESPONSE,
      error: undefined,
      response: { status: 201 },
    });

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");

    fillTextarea(textarea as HTMLTextAreaElement, VALID_CONTENT);
    await submitForm(form as HTMLFormElement);

    expect(mockPostV1Posts).toHaveBeenCalledWith(
      expect.objectContaining({ body: { content: VALID_CONTENT } }),
    );
  });

  it("shows a success toast after a successful submit", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    mockPostV1Posts.mockResolvedValue({
      data: MOCK_RESPONSE,
      error: undefined,
      response: { status: 201 },
    });

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");

    fillTextarea(textarea as HTMLTextAreaElement, VALID_CONTENT);
    await submitForm(form as HTMLFormElement);

    expect(mockToastSuccess).toHaveBeenCalledTimes(1);
  });

  it("resets the form after a successful submit", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    mockPostV1Posts.mockResolvedValue({
      data: MOCK_RESPONSE,
      error: undefined,
      response: { status: 201 },
    });

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");

    fillTextarea(textarea as HTMLTextAreaElement, VALID_CONTENT);
    await submitForm(form as HTMLFormElement);

    // After reset, the textarea value should be empty and charCount back to 0
    const charCountEl = container.querySelector<HTMLElement>("[data-testid='char-count']");
    expect(charCountEl?.textContent).toBe("0");
  });
});

// ---------------------------------------------------------------------------
// Tests — submit loading state
// ---------------------------------------------------------------------------

describe("useCreatePostForm — loading state", () => {
  it("disables the submit button while the API call is in-flight", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    // Never resolves — keeps isSubmitting true
    mockPostV1Posts.mockImplementation(() => new Promise<never>(() => {}));

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");

    fillTextarea(textarea as HTMLTextAreaElement, "Some post content");
    await submitForm(form as HTMLFormElement);

    const submitBtn =
      container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
      container.querySelector<HTMLButtonElement>("button");
    expect(submitBtn?.disabled).toBe(true);
  });
});

// ---------------------------------------------------------------------------
// Tests — error path submit
// ---------------------------------------------------------------------------

describe("useCreatePostForm — failed submit", () => {
  it("navigates to /login on a 401 error from postV1Posts", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    mockPostV1Posts.mockResolvedValue({
      data: undefined,
      error: { type: "UNAUTHORIZED", title: "Unauthorized", status: 401 },
      response: { status: 401 },
    });

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");

    fillTextarea(textarea as HTMLTextAreaElement, "Some post content");
    await submitForm(form as HTMLFormElement);

    expect(mockNavigate).toHaveBeenCalledWith({ to: "/login" });
  });

  it("shows an error toast on a 500 server error", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    mockPostV1Posts.mockResolvedValue({
      data: undefined,
      error: { type: "INTERNAL_ERROR", title: "Internal Server Error", status: 500 },
      response: { status: 500 },
    });

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");

    fillTextarea(textarea as HTMLTextAreaElement, "Some post content");
    await submitForm(form as HTMLFormElement);

    expect(mockToastError).toHaveBeenCalledTimes(1);
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it("shows an error toast on a network-level failure", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    mockPostV1Posts.mockRejectedValue(new Error("Network failure"));

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");

    fillTextarea(textarea as HTMLTextAreaElement, "Some post content");
    await submitForm(form as HTMLFormElement);

    expect(mockToastError).toHaveBeenCalledTimes(1);
    expect(mockNavigate).not.toHaveBeenCalled();
  });

  it("does not navigate to /login on a non-401 error", async () => {
    mockGetToken.mockReturnValue("valid-jwt-token");
    mockPostV1Posts.mockResolvedValue({
      data: undefined,
      error: { type: "FORBIDDEN", title: "Forbidden", status: 403 },
      response: { status: 403 },
    });

    const container = mountTestComponent();
    const textarea = container.querySelector<HTMLTextAreaElement>("textarea");
    const form = container.querySelector<HTMLFormElement>("form");

    fillTextarea(textarea as HTMLTextAreaElement, "Some post content");
    await submitForm(form as HTMLFormElement);

    expect(mockNavigate).not.toHaveBeenCalled();
  });
});
