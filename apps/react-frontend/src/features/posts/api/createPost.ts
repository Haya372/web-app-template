import { getToken } from "@/features/auth/utils/tokenStorage"
import { createPostResponseSchema } from "@/features/posts/types/post"
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

	const response = await fetch(`${baseUrl}/v1/posts`, {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			Authorization: `Bearer ${token}`,
		},
		body: JSON.stringify({ content }),
	})

	if (!response.ok) {
		throw new Error(String(response.status))
	}

	const data: unknown = await response.json()
	return createPostResponseSchema.parse(data)
}
