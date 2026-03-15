import { postV1Posts } from "@/generated/sdk.gen"
import { getToken } from "@/features/auth/utils/tokenStorage"
import type { CreatePostResponse } from "@/features/posts/types/post"

export async function callCreatePost(content: string): Promise<CreatePostResponse> {
	const baseUrl = import.meta.env.VITE_API_BASE_URL
	if (!baseUrl) {
		throw new Error("VITE_API_BASE_URL is not set")
	}

	const token = getToken()
	if (!token) {
		throw new Error("Unauthenticated: no token available")
	}

	const { data, error, response } = await postV1Posts({
		body: { content },
		baseUrl,
		headers: { Authorization: `Bearer ${token}` },
	})

	if (error || !data) {
		throw new Error(String(response.status))
	}

	return data
}
