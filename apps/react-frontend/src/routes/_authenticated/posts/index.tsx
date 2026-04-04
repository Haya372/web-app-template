import { createFileRoute } from "@tanstack/react-router";
import { PostsListPage } from "@/features/posts/pages/PostsListPage";

export const Route = createFileRoute("/_authenticated/posts/")({
	component: PostsListPage,
});
