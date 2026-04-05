/**
 * Tests for LoginPage component.
 *
 * Uses ReactDOM + manual DOM querying since @testing-library/react is not installed.
 *
 * Mocks:
 *  - @/generated/sdk.gen               postV1UsersLogin vi.mock
 *  - @/utils/tokenStorage              saveToken vi.mock
 *  - @tanstack/react-router            useNavigate vi.mock
 *  - @repo/ui                          toast vi.mock
 */

import React, { act } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vite-plus/test";

// ---------------------------------------------------------------------------
// Hoisted mock functions (must be defined before vi.mock calls)
// ---------------------------------------------------------------------------

const { mockNavigate, mockToastError, mockLogin } = vi.hoisted(() => ({
  mockNavigate: vi.fn(),
  mockToastError: vi.fn(),
  mockLogin: vi.fn(),
}));

// ---------------------------------------------------------------------------
// Module-level mocks
// ---------------------------------------------------------------------------

vi.mock("@/generated/sdk.gen", () => ({
  postV1UsersLogin: vi.fn(),
}));

vi.mock("@/utils/tokenStorage", () => ({
  saveToken: vi.fn(),
  getToken: vi.fn().mockReturnValue(null),
  removeToken: vi.fn(),
}));

vi.mock("@/features/auth/hooks/useAuth", () => ({
  useAuth: () => ({
    token: null,
    isAuthenticated: false,
    login: mockLogin,
    logout: vi.fn(),
  }),
}));

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => mockNavigate,
  Link: ({ children, to }: { children: React.ReactNode; to: string }) =>
    React.createElement("a", { href: to }, children),
}));

vi.mock("@repo/ui", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@repo/ui")>();
  return {
    ...actual,
    toast: {
      ...actual.toast,
      error: mockToastError,
    },
  };
});

// ---------------------------------------------------------------------------
// Imports after mocks
// ---------------------------------------------------------------------------

import { postV1UsersLogin } from "@/generated/sdk.gen";
import { LoginPage } from "./LoginPage";

// ---------------------------------------------------------------------------
// Typed mock references
// ---------------------------------------------------------------------------

const mockPostV1UsersLogin = postV1UsersLogin as ReturnType<typeof vi.fn>;

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function mountLoginPage(): HTMLDivElement {
  const container = document.createElement("div");
  document.body.append(container);
  act(() => {
    createRoot(container).render(<LoginPage />);
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
  mockPostV1UsersLogin.mockReset();
  mockLogin.mockReset();
  mockNavigate.mockReset();
  mockToastError.mockReset();
});

afterEach(() => {
  clearBody();
  vi.unstubAllEnvs();
});

// ---------------------------------------------------------------------------
// Tests — rendering
// ---------------------------------------------------------------------------

describe("LoginPage — rendering", () => {
  it("renders an email input", () => {
    const container = mountLoginPage();
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    expect(emailInput).not.toBeNull();
  });

  it("renders a password input with type='password'", () => {
    const container = mountLoginPage();
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    expect(passwordInput).not.toBeNull();
  });

  it("renders a submit button", () => {
    const container = mountLoginPage();
    const submitBtn =
      container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
      container.querySelector<HTMLButtonElement>("button");
    expect(submitBtn).not.toBeNull();
  });
});

// ---------------------------------------------------------------------------
// Tests — client-side validation
// ---------------------------------------------------------------------------

describe("LoginPage — client-side validation", () => {
  it("shows 'Invalid email address' when email is empty on submit", async () => {
    const container = mountLoginPage();
    const form = container.querySelector<HTMLFormElement>("form");
    expect(form).not.toBeNull();

    await submitForm(form as HTMLFormElement);

    expect(container.textContent).toContain("Invalid email address");
  });

  it("shows 'Password is required' when password is empty on submit", async () => {
    const container = mountLoginPage();
    const form = container.querySelector<HTMLFormElement>("form");
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    expect(form).not.toBeNull();
    expect(emailInput).not.toBeNull();

    fillInput(emailInput as HTMLInputElement, "user@example.com");
    await submitForm(form as HTMLFormElement);

    expect(container.textContent).toContain("Password is required");
  });

  it("does not call postV1UsersLogin when validation fails", async () => {
    const container = mountLoginPage();
    const form = container.querySelector<HTMLFormElement>("form");
    await submitForm(form as HTMLFormElement);

    expect(mockPostV1UsersLogin).not.toHaveBeenCalled();
  });
});

// ---------------------------------------------------------------------------
// Tests — loading state
// ---------------------------------------------------------------------------

describe("LoginPage — loading state", () => {
  it("disables the submit button while the API call is in-flight", async () => {
    mockPostV1UsersLogin.mockImplementation(() => new Promise<never>(() => {}));

    const container = mountLoginPage();
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(emailInput as HTMLInputElement, "user@example.com");
    fillInput(passwordInput as HTMLInputElement, "password123");
    await submitForm(form as HTMLFormElement);

    const submitBtn =
      container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
      container.querySelector<HTMLButtonElement>("button");
    expect(submitBtn?.disabled).toBe(true);
  });
});

// ---------------------------------------------------------------------------
// Tests — successful login
// ---------------------------------------------------------------------------

describe("LoginPage — successful login", () => {
  const VALID_RESPONSE = {
    token: "jwt-abc-123",
    expiresAt: "2026-12-31T00:00:00Z",
    user: { id: "u1", name: "Alice", email: "alice@example.com" },
  };

  it("calls login with the token returned by the API", async () => {
    mockPostV1UsersLogin.mockResolvedValue({
      data: VALID_RESPONSE,
      error: undefined,
      response: { status: 200 },
    });

    const container = mountLoginPage();
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(emailInput as HTMLInputElement, "alice@example.com");
    fillInput(passwordInput as HTMLInputElement, "secret");

    await submitForm(form as HTMLFormElement);

    expect(mockLogin).toHaveBeenCalledWith("jwt-abc-123");
  });

  it("navigates to '/' after a successful login", async () => {
    mockPostV1UsersLogin.mockResolvedValue({
      data: VALID_RESPONSE,
      error: undefined,
      response: { status: 200 },
    });

    const container = mountLoginPage();
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(emailInput as HTMLInputElement, "alice@example.com");
    fillInput(passwordInput as HTMLInputElement, "secret");

    await submitForm(form as HTMLFormElement);

    expect(mockNavigate).toHaveBeenCalledWith({ to: "/" });
  });
});

// ---------------------------------------------------------------------------
// Tests — failed login
// ---------------------------------------------------------------------------

describe("LoginPage — failed login", () => {
  const errorCases = [
    {
      name: "shows 'Invalid email or password' toast on a 401 error",
      mock: { data: undefined, error: {}, response: { status: 401 } },
      expectedMessage: "Invalid email or password",
    },
    {
      name: "shows generic error toast on a 500 server error",
      mock: { data: undefined, error: {}, response: { status: 500 } },
      expectedMessage: "Login failed. Please try again.",
    },
    {
      name: "shows generic error toast on a 403 error",
      mock: { data: undefined, error: {}, response: { status: 403 } },
      expectedMessage: "Login failed. Please try again.",
    },
  ];

  for (const { name, mock, expectedMessage } of errorCases) {
    it(name, async () => {
      mockPostV1UsersLogin.mockResolvedValue(mock);

      const container = mountLoginPage();
      const emailInput = container.querySelector<HTMLInputElement>(
        'input[type="email"], input[name="email"]',
      );
      const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
      const form = container.querySelector<HTMLFormElement>("form");

      fillInput(emailInput as HTMLInputElement, "user@example.com");
      fillInput(passwordInput as HTMLInputElement, "password123");

      await submitForm(form as HTMLFormElement);

      expect(mockToastError).toHaveBeenCalledWith(expect.stringContaining(expectedMessage));
    });
  }

  it("shows generic error toast on a network-level failure", async () => {
    mockPostV1UsersLogin.mockRejectedValue(new Error("Network failure"));

    const container = mountLoginPage();
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(emailInput as HTMLInputElement, "user@example.com");
    fillInput(passwordInput as HTMLInputElement, "password123");

    await submitForm(form as HTMLFormElement);

    expect(mockToastError).toHaveBeenCalledWith(
      expect.stringContaining("Login failed. Please try again."),
    );
  });

  it("re-enables the submit button after a failed login", async () => {
    mockPostV1UsersLogin.mockResolvedValue({
      data: undefined,
      error: {},
      response: { status: 500 },
    });

    const container = mountLoginPage();
    const emailInput = container.querySelector<HTMLInputElement>(
      'input[type="email"], input[name="email"]',
    );
    const passwordInput = container.querySelector<HTMLInputElement>('input[type="password"]');
    const form = container.querySelector<HTMLFormElement>("form");

    fillInput(emailInput as HTMLInputElement, "user@example.com");
    fillInput(passwordInput as HTMLInputElement, "password123");

    await submitForm(form as HTMLFormElement);

    const submitBtn =
      container.querySelector<HTMLButtonElement>('button[type="submit"]') ??
      container.querySelector<HTMLButtonElement>("button");
    expect(submitBtn?.disabled).toBe(false);
  });
});
