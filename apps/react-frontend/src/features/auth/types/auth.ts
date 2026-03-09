import { z } from "zod"

export const loginResponseSchema = z.object({
	token: z.string(),
	expiresAt: z.string(),
	user: z.object({
		id: z.string(),
		name: z.string(),
		email: z.string(),
	}),
})

export type LoginResponse = z.infer<typeof loginResponseSchema>
