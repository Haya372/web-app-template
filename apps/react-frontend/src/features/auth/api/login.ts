import { postV1UsersLogin } from "@/generated/sdk.gen"
import type { LoginResponse } from "@/features/auth/types/auth"

export async function callLogin(
	email: string,
	password: string,
): Promise<LoginResponse> {
	const baseUrl = import.meta.env.VITE_API_BASE_URL
	if (!baseUrl) {
		throw new Error("VITE_API_BASE_URL is not set")
	}

	const { data, error, response } = await postV1UsersLogin({
		body: { email, password },
		baseUrl,
	})

	if (error || !data) {
		throw new Error(String(response.status))
	}

	return data
}
