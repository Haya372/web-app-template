import { test, expect } from "@playwright/test";

const API_BASE_URL = process.env.E2E_API_BASE_URL ?? "http://localhost:8080";

async function registerAndLogin(
	name: string,
	email: string,
	password: string,
): Promise<string> {
	const signupRes = await fetch(`${API_BASE_URL}/v1/users/signup`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ name, email, password }),
	});
	if (!signupRes.ok) {
		throw new Error(`Signup API failed: ${signupRes.status}`);
	}

	const loginRes = await fetch(`${API_BASE_URL}/v1/users/login`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ email, password }),
	});
	if (!loginRes.ok) {
		throw new Error(`Login API failed: ${loginRes.status}`);
	}

	const { token } = (await loginRes.json()) as { token: string };
	return token;
}

test.describe("Post creation", () => {
	test("happy path: logged-in user can create a post", async ({ page }) => {
		const unique = Date.now();
		const email = `e2e-post-${unique}@example.com`;
		const password = "password123";
		const token = await registerAndLogin(
			`E2E Post User ${unique}`,
			email,
			password,
		);

		await page.goto("/");
		await page.evaluate((t) => {
			localStorage.setItem("auth_token", t);
		}, token);
		await page.goto("/posts/new");

		await page.getByLabel("内容").fill("E2E test post content");
		await page.getByRole("button", { name: "投稿する" }).click();

		await expect(page.getByText("投稿しました")).toBeVisible();
	});

	test("validation error: empty content shows error message", async ({
		page,
	}) => {
		const unique = Date.now();
		const email = `e2e-post-empty-${unique}@example.com`;
		const password = "password123";
		const token = await registerAndLogin(
			`E2E Post Empty ${unique}`,
			email,
			password,
		);

		await page.goto("/");
		await page.evaluate((t) => {
			localStorage.setItem("auth_token", t);
		}, token);
		await page.goto("/posts/new");

		await page.getByRole("button", { name: "投稿する" }).click();

		await expect(page.getByText("内容を入力してください")).toBeVisible();
	});

	test("validation error: content exceeding 280 characters shows error", async ({
		page,
	}) => {
		const unique = Date.now();
		const email = `e2e-post-long-${unique}@example.com`;
		const password = "password123";
		const token = await registerAndLogin(
			`E2E Post Long ${unique}`,
			email,
			password,
		);

		await page.goto("/");
		await page.evaluate((t) => {
			localStorage.setItem("auth_token", t);
		}, token);
		await page.goto("/posts/new");

		await page.getByLabel("内容").fill("a".repeat(281));
		await page.getByRole("button", { name: "投稿する" }).click();

		await expect(
			page.getByText("内容は280文字以内で入力してください"),
		).toBeVisible();
	});
});
