import { test, expect } from "@playwright/test";

const API_BASE_URL = process.env.E2E_API_BASE_URL ?? "http://localhost:8080";

async function registerUser(
	name: string,
	email: string,
	password: string,
): Promise<void> {
	const res = await fetch(`${API_BASE_URL}/v1/users/signup`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ name, email, password }),
	});
	if (!res.ok) {
		throw new Error(`Signup API failed: ${res.status}`);
	}
}

test.describe("Login", () => {
	test("happy path: valid credentials redirect to home", async ({ page }) => {
		const unique = Date.now();
		const email = `e2e-login-${unique}@example.com`;
		const password = "password123";
		await registerUser(`E2E Login User ${unique}`, email, password);

		await page.goto("/login");

		await page.getByLabel("メールアドレス").fill(email);
		await page.getByLabel("パスワード").fill(password);
		await page.getByRole("button", { name: "ログイン" }).click();

		await expect(page).toHaveURL("/");
	});

	test("error case: wrong password shows error toast", async ({ page }) => {
		const unique = Date.now();
		const email = `e2e-login-err-${unique}@example.com`;
		const password = "password123";
		await registerUser(`E2E Login Err ${unique}`, email, password);

		await page.goto("/login");

		await page.getByLabel("メールアドレス").fill(email);
		await page.getByLabel("パスワード").fill("wrongpassword");
		await page.getByRole("button", { name: "ログイン" }).click();

		await expect(
			page.getByText("メールアドレスまたはパスワードが正しくありません"),
		).toBeVisible();
		await expect(page).toHaveURL(/\/login/);
	});
});
