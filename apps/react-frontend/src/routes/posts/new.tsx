import { createFileRoute } from '@tanstack/react-router'
import { NewPostPage } from '@/features/posts/pages/NewPostPage'

export const Route = createFileRoute('/posts/new')({
	component: NewPostPage,
})
