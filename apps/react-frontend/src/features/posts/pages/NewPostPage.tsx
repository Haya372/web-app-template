import {
  Button,
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  Heading,
  Textarea,
} from "@repo/ui";
import { useTranslation } from "react-i18next";
import { useCreatePostForm } from "@/features/posts/hooks/useCreatePostForm";

export function NewPostPage() {
  const { t } = useTranslation();
  const { form, onSubmit, charCount } = useCreatePostForm();

  return (
    <main className="page-wrap flex min-h-[60vh] items-center justify-center px-4 py-12">
      <section className="island-shell w-full max-w-lg rounded-2xl p-6 sm:p-8">
        <Heading level={1}>{t("posts.new.title")}</Heading>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} noValidate>
            <div className="mb-2">
              <FormField
                control={form.control}
                name="content"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel htmlFor="post-content">{t("posts.new.label")}</FormLabel>
                    <FormControl>
                      <Textarea id="post-content" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="mb-4 text-sm text-right">
              <span data-testid="char-count">
                {t("posts.new.charCount", { current: charCount, max: 280 })}
              </span>
            </div>

            <div className="w-full">
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {t("posts.new.submit")}
              </Button>
            </div>
          </form>
        </Form>
      </section>
    </main>
  );
}
