import { test, expect } from "@playwright/test";

test.describe("Auth guard", () => {
	test("unauthenticated access to /posts/new redirects to /login", async ({
		page,
	}) => {
		await page.goto("/posts/new");

		await expect(page).toHaveURL(/\/login/);
	});
});
