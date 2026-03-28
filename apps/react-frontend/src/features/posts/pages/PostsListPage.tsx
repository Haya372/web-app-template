import { Heading } from "@repo/ui"
import { useTranslation } from "react-i18next"
import { usePostsList } from "@/features/posts/hooks/usePostsList"

export function PostsListPage() {
	const { t } = useTranslation()
	const { posts, isLoading, isError } = usePostsList()

	return (
		<main className="page-wrap px-4 py-12">
			<Heading level={1}>{t("posts.list.title")}</Heading>

			<div aria-live="polite" aria-atomic="true">
				{isLoading && (
					<p role="status" className="mt-4 text-muted-foreground">
						{t("posts.list.loading")}
					</p>
				)}

				{isError && (
					<p role="alert" className="mt-4 text-destructive">
						{t("posts.list.error.generic")}
					</p>
				)}

				{!isLoading && !isError && posts !== undefined && posts.length === 0 && (
					<p className="mt-4 text-muted-foreground">{t("posts.list.empty")}</p>
				)}

				{!isLoading && !isError && posts !== undefined && posts.length > 0 && (
					<ul className="mt-6 space-y-4">
						{posts.map((post) => (
							<li key={post.id} className="island-shell rounded-xl p-4">
								<p>{post.content}</p>
							</li>
						))}
					</ul>
				)}
			</div>
		</main>
	)
}
