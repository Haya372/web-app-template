import { test } from "@playwright/test";

// TODO: Implement auth guard E2E tests once route-level authentication is added.
//
// Scenario: unauthenticated user directly accessing a protected route (e.g. /dashboard)
// should be redirected to the login page.
test.describe("Auth guard", () => {
  test.fixme(
    "unauthenticated access to protected route redirects to /login",
    async () => {
      // 1. Ensure no JWT token in localStorage.
      // 2. Navigate directly to the protected route (e.g. /dashboard).
      // 3. Verify the browser is redirected to /login.
    },
  );
});
