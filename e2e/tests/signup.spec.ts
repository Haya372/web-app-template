import { test, expect } from "@playwright/test";

test.describe("Signup", () => {
	test("happy path: fills form and redirects to login on success", async ({
		page,
	}) => {
		const unique = Date.now();

		await page.goto("/signup");

		await page.getByLabel("名前").fill(`E2E User ${unique}`);
		await page.getByLabel("メールアドレス").fill(`e2e-signup-${unique}@example.com`);
		await page.getByLabel("パスワード").fill("password123");

		await page.getByRole("button", { name: "登録する" }).click();

		await expect(page).toHaveURL(/\/login/);
	});

	test("shows error when email is already registered", async ({ page }) => {
		const unique = Date.now();
		const email = `e2e-dup-${unique}@example.com`;
		const password = "password123";

		// First signup
		await page.goto("/signup");
		await page.getByLabel("名前").fill(`E2E Dup User ${unique}`);
		await page.getByLabel("メールアドレス").fill(email);
		await page.getByLabel("パスワード").fill(password);
		await page.getByRole("button", { name: "登録する" }).click();
		await page.waitForURL(/\/login/);

		// Second signup with same email
		await page.goto("/signup");
		await page.getByLabel("名前").fill("別のユーザー");
		await page.getByLabel("メールアドレス").fill(email);
		await page.getByLabel("パスワード").fill(password);
		await page.getByRole("button", { name: "登録する" }).click();

		await expect(
			page.getByText("このメールアドレスはすでに登録されています"),
		).toBeVisible();
	});
});
