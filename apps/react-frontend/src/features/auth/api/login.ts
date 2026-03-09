import { loginResponseSchema } from "@/features/auth/types/auth"
import type { LoginResponse } from "@/features/auth/types/auth"

export async function callLogin(
	email: string,
	password: string,
): Promise<LoginResponse> {
	const baseUrl = import.meta.env.VITE_API_BASE_URL
	if (!baseUrl) {
		throw new Error("VITE_API_BASE_URL is not set")
	}

	const response = await fetch(`${baseUrl}/v1/users/login`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({ email, password }),
	})

	if (!response.ok) {
		throw new Error(String(response.status))
	}

	const data: unknown = await response.json()
	return loginResponseSchema.parse(data)
}
