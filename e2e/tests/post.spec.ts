import { test } from "@playwright/test";

// TODO: Implement post creation E2E tests once the post creation UI is built.
// Tracked in a separate ticket.
//
// Scenario: logged-in user fills the post content field and submits;
// the new post should appear in the post list.
test.describe("Post creation", () => {
  test.fixme(
    "happy path: logged-in user can create a post and it appears in the list",
    async () => {
      // 1. Log in via the login page (or set JWT in localStorage directly).
      // 2. Navigate to the post creation page.
      // 3. Fill the content field and click Submit.
      // 4. Verify the new post appears in the post list.
    },
  );
});
