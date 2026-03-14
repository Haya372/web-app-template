import { z } from "zod"

export const createPostResponseSchema = z.object({
	id: z.string(),
	userId: z.string(),
	content: z.string(),
	createdAt: z.string(),
})

export type CreatePostResponse = z.infer<typeof createPostResponseSchema>
