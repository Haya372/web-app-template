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

export const signupResponseSchema = z.object({
	id: z.string(),
	name: z.string(),
	email: z.string(),
	status: z.string(),
	createdAt: z.string(),
})

export type SignupResponse = z.infer<typeof signupResponseSchema>
