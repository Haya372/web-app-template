/**
 * Tests for SignupPage component.
 *
 * Uses ReactDOM + manual DOM querying since @testing-library/react is not installed.
 *
 * Mocks:
 *  - @/generated/sdk.gen   postV1UsersSignup vi.mock
 *  - @tanstack/react-router       useNavigate vi.mock
 *
 * Note: errors are shown inline in the page (NOT via toast).
 */

import { act } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vite-plus/test";

// ---------------------------------------------------------------------------
// Hoisted mock functions (must be defined before vi.mock calls)
// ---------------------------------------------------------------------------

const { mockNavigate } = vi.hoisted(() => ({
  mockNavigate: vi.fn(),
}));

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

vi.mock("@/generated/sdk.gen", () => ({
  postV1UsersSignup: vi.fn(),
}));

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => mockNavigate,
}));

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { postV1UsersSignup } from "@/generated/sdk.gen";
import { SignupPage } from "./SignupPage";

// ---------------------------------------------------------------------------
// Typed mock references
// ---------------------------------------------------------------------------

const mockPostV1UsersSignup = postV1UsersSignup as ReturnType<typeof vi.fn>;

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function mountSignupPage(): HTMLDivElement {
  const container = document.createElement("div");
  document.body.append(container);
  act(() => {
    createRoot(container).render(<SignupPage />);
  });
  return container;
}

async function submitForm(form: HTMLFormElement): Promise<void> {
  await act(async () => {
    form.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true }));
    // Allow react-hook-form's async validation to resolve
    await Promise.resolve();
    await Promise.resolve();
  });
}

function fillInput(input: HTMLInputElement, value: string): void {
  act(() => {
    const nativeValueSetter = Object.getOwnPropertyDescriptor(
      HTMLInputElement.prototype,
      "value",
    )?.set;
    nativeValueSetter?.call(input, value);
    input.dispatchEvent(new Event("input", { bubbles: true }));
    input.dispatchEvent(new Event("change", { bubbles: true }));
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
  mockPostV1UsersSignup.mockReset();
  mockNavigate.mockReset();
});

afterEach(() => {
  clearBody();
  vi.unstubAllEnvs();
});

// ---------------------------------------------------------------------------
// Tests — rendering
// ---------------------------------------------------------------------------

describe("SignupPage — rendering", () => {
  it("renders a name input", () => {
    const container = mountSignupPage();
    const nameInput = container.querySelector<HTMLInputElement>(
      'input[name="name"], input[placeholder*="name" i]',
    );
    expect(nameInput).not.toBeNull();
  });

  it("renders an email input", () => {
    const container = mountSignupPage();
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    expect(emailInput).not.toBeNull();
  });

  it("renders a password input with type='password'", () => {
    const container = mountSignupPage();
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    expect(passwordInput).not.toBeNull();
  });

  it("renders a submit button", () => {
    const container = mountSignupPage();
    const submitBtn =
      container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
      container.querySelector<HTMLButtonElement>("button");
    expect(submitBtn).not.toBeNull();
  });
});

// ---------------------------------------------------------------------------
// Tests — client-side validation
// ---------------------------------------------------------------------------

describe("SignupPage — client-side validation", () => {
  it("shows 'Name is required' when name is empty on submit", async () => {
    const container = mountSignupPage();
    const form = container.querySelector<HTMLFormElement>("form");
    expect(form).not.toBeNull();

    await submitForm(form as HTMLFormElement);

    expect(container.textContent).toContain("Name is required");
  });

  it("shows 'Invalid email address' when email is empty on submit", async () => {
    const container = mountSignupPage();
    const form = container.querySelector<HTMLFormElement>("form");
    const nameInput = container.querySelector<HTMLInputElement>(
      'input[name="name"], input[placeholder*="name" i]',
    );
    expect(form).not.toBeNull();
    expect(nameInput).not.toBeNull();

    fillInput(nameInput as HTMLInputElement, "Alice");
    await submitForm(form as HTMLFormElement);

    expect(container.textContent).toContain("Invalid email address");
  });

  it("shows 'Password is required' when password is empty on submit", async () => {
    const container = mountSignupPage();
    const form = container.querySelector<HTMLFormElement>("form");
    const nameInput = container.querySelector<HTMLInputElement>(
      'input[name="name"], input[placeholder*="name" i]',
    );
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    expect(form).not.toBeNull();
    expect(nameInput).not.toBeNull();
    expect(emailInput).not.toBeNull();

    fillInput(nameInput as HTMLInputElement, "Alice");
    fillInput(emailInput as HTMLInputElement, "alice@example.com");
    await submitForm(form as HTMLFormElement);

    expect(container.textContent).toContain("Password is required");
  });

  it("does not call postV1UsersSignup when validation fails", async () => {
    const container = mountSignupPage();
    const form = container.querySelector<HTMLFormElement>("form");
    await submitForm(form as HTMLFormElement);

    expect(mockPostV1UsersSignup).not.toHaveBeenCalled();
  });
});

// ---------------------------------------------------------------------------
// Tests — loading state
// ---------------------------------------------------------------------------

describe("SignupPage — loading state", () => {
  it("disables the submit button while the API call is in-flight", async () => {
    mockPostV1UsersSignup.mockImplementation(() => new Promise<never>(() => {}));

    const container = mountSignupPage();
    const nameInput = container.querySelector<HTMLInputElement>(
      'input[name="name"], input[placeholder*="name" i]',
    );
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(nameInput as HTMLInputElement, "Alice");
    fillInput(emailInput as HTMLInputElement, "alice@example.com");
    fillInput(passwordInput as HTMLInputElement, "password123");
    await submitForm(form as HTMLFormElement);

    const submitBtn =
      container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
      container.querySelector<HTMLButtonElement>("button");
    expect(submitBtn?.disabled).toBe(true);
  });
});

// ---------------------------------------------------------------------------
// Tests — successful signup
// ---------------------------------------------------------------------------

describe("SignupPage — successful signup", () => {
  it("navigates to '/login' after a successful signup", async () => {
    mockPostV1UsersSignup.mockResolvedValue({
      data: {
        id: "u1",
        name: "Alice",
        email: "alice@example.com",
        status: "active",
        createdAt: "2026-03-15T00:00:00Z",
      },
      error: undefined,
      response: { status: 201 },
    });

    const container = mountSignupPage();
    const nameInput = container.querySelector<HTMLInputElement>(
      'input[name="name"], input[placeholder*="name" i]',
    );
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(nameInput as HTMLInputElement, "Alice");
    fillInput(emailInput as HTMLInputElement, "alice@example.com");
    fillInput(passwordInput as HTMLInputElement, "password123");

    await submitForm(form as HTMLFormElement);

    expect(mockNavigate).toHaveBeenCalledWith({ to: "/login" });
  });
});

// ---------------------------------------------------------------------------
// Tests — failed signup (inline errors, NOT toast)
// ---------------------------------------------------------------------------

describe("SignupPage — failed signup", () => {
  const errorCases = [
    {
      name: "shows 'Email already registered' inline error on a 409 conflict",
      mock: { data: undefined, error: {}, response: { status: 409 } },
      expectedMessage: "Email already registered",
    },
    {
      name: "shows generic inline error on a 500 server error",
      mock: { data: undefined, error: {}, response: { status: 500 } },
      expectedMessage: "Signup failed",
    },
  ];

  for (const { name, mock, expectedMessage } of errorCases) {
    it(name, async () => {
      mockPostV1UsersSignup.mockResolvedValue(mock);

      const container = mountSignupPage();
      const nameInput = container.querySelector<HTMLInputElement>(
        'input[name="name"], input[placeholder*="name" i]',
      );
      const emailInput = container.querySelector<HTMLInputElement>(
        'input[type="email"], input[name="email"]',
      );
      const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
      const form = container.querySelector<HTMLFormElement>("form");

      fillInput(nameInput as HTMLInputElement, "Alice");
      fillInput(emailInput as HTMLInputElement, "alice@example.com");
      fillInput(passwordInput as HTMLInputElement, "password123");

      await submitForm(form as HTMLFormElement);

      // Error must be visible inline in the page, not in a toast
      expect(container.textContent).toContain(expectedMessage);
    });
  }

  it("shows generic inline error on a network-level failure", async () => {
    mockPostV1UsersSignup.mockRejectedValue(new Error("Network failure"));

    const container = mountSignupPage();
    const nameInput = container.querySelector<HTMLInputElement>(
      'input[name="name"], input[placeholder*="name" i]',
    );
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(nameInput as HTMLInputElement, "Alice");
    fillInput(emailInput as HTMLInputElement, "alice@example.com");
    fillInput(passwordInput as HTMLInputElement, "password123");

    await submitForm(form as HTMLFormElement);

    expect(container.textContent).toContain("Signup failed");
  });

  it("re-enables the submit button after a failed signup", async () => {
    mockPostV1UsersSignup.mockResolvedValue({
      data: undefined,
      error: {},
      response: { status: 500 },
    });

    const container = mountSignupPage();
    const nameInput = container.querySelector<HTMLInputElement>(
      'input[name="name"], input[placeholder*="name" i]',
    );
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(nameInput as HTMLInputElement, "Alice");
    fillInput(emailInput as HTMLInputElement, "alice@example.com");
    fillInput(passwordInput as HTMLInputElement, "password123");

    await submitForm(form as HTMLFormElement);

    const submitBtn =
      container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
      container.querySelector<HTMLButtonElement>("button");
    expect(submitBtn?.disabled).toBe(false);
  });
});
