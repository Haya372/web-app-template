import { Heading, List, ListItem, Typography } from "@repo/ui";
import { useTranslation } from "react-i18next";
import { usePostsList } from "@/features/posts/hooks/usePostsList";

export function PostsListPage() {
	const { t } = useTranslation();
	const { posts, isLoading, isError } = usePostsList();

	return (
		<main className="page-wrap px-4 py-12">
			<Heading level={1}>{t("posts.list.title")}</Heading>

			<div aria-live="polite" aria-atomic="true">
				{isLoading && (
					<Typography role="status" variant="muted">
						{t("posts.list.loading")}
					</Typography>
				)}

				{isError && (
					<Typography role="alert" variant="muted">
						{t("posts.list.error.generic")}
					</Typography>
				)}

				{!isLoading &&
					!isError &&
					posts !== undefined &&
					posts.length === 0 && (
						<Typography variant="muted">{t("posts.list.empty")}</Typography>
					)}

				{!isLoading && !isError && posts !== undefined && posts.length > 0 && (
					<List>
						{posts.map((post) => (
							<ListItem key={post.id}>
								<Typography>{post.content}</Typography>
							</ListItem>
						))}
					</List>
				)}
			</div>
		</main>
	);
}
