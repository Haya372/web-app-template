import { postV1UsersSignup } from "@/generated/sdk.gen"
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

	const { data, error, response } = await postV1UsersSignup({
		body: { name, email, password },
		baseUrl,
	})

	if (error || !data) {
		throw new Error(String(response.status))
	}

	return data
}
