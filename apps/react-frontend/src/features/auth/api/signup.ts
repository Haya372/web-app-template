import { signupResponseSchema } from "@/features/auth/types/auth"
import type { SignupResponse } from "@/features/auth/types/auth"

export async function callSignup(
	name: string,
	email: string,
	password: string,
): Promise<SignupResponse> {
	const baseUrl = import.meta.env.VITE_API_BASE_URL
	if (!baseUrl) {
		throw new Error("VITE_API_BASE_URL is not set")
	}

	const response = await fetch(`${baseUrl}/v1/users/signup`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify({ name, email, password }),
	})

	if (!response.ok) {
		throw new Error(String(response.status))
	}

	const data: unknown = await response.json()
	return signupResponseSchema.parse(data)
}
