import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "@repo/ui";
import { useEffect } from "react";
import { useForm, useWatch } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { postV1Posts } from "@/generated/sdk.gen";
import { getToken } from "@/utils/tokenStorage";

function useCreatePostFormSchema() {
  const { t } = useTranslation();
  return z.object({
    content: z
      .string()
      .min(1, t("posts.new.validation.contentRequired"))
      .max(280, t("posts.new.validation.contentTooLong")),
  });
}

type CreatePostFormValues = z.infer<ReturnType<typeof useCreatePostFormSchema>>;

export function useCreatePostForm() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const schema = useCreatePostFormSchema();

  const form = useForm<CreatePostFormValues>({
    resolver: zodResolver(schema),
    defaultValues: { content: "" },
  });

  // Redirect to login if no token is present on mount.
  // getToken() reads localStorage synchronously and is not reactive, so it is
  // intentionally omitted from the dependency array.
  useEffect(() => {
    if (getToken() === null) {
      void navigate({ to: "/login" });
    }
  }, [navigate]);

  // useWatch scopes re-renders to this consumer only, avoiding a full form
  // re-render on every keystroke that form.watch() would cause.
  // useWatch returns undefined before RHF initialises its internal state on the
  // first render. The nullish fallback ensures charCount is always a number.
  const content = useWatch({ control: form.control, name: "content" }) ?? "";
  const charCount = content.length;

  async function onSubmit(values: CreatePostFormValues) {
    try {
      const token = getToken();
      const { data, error, response } = await postV1Posts({
        body: { content: values.content },
        baseUrl: import.meta.env.VITE_API_BASE_URL,
        headers: token ? { Authorization: `Bearer ${token}` } : undefined,
      });
      if (error || !data) throw new Error(String(response.status));
      toast.success(t("posts.new.success"));
      form.reset();
    } catch (error) {
      // postV1Posts throws on network-level failures; for HTTP errors we throw
      // new Error(String(response.status)) ourselves above.
      // Only a server-returned 401 redirects to login; all other failures show
      // a generic toast. The pre-flight mount guard already redirects if no
      // token is present, so an auth failure here is unexpected but handled.
      if (error instanceof Error && error.message === "401") {
        void navigate({ to: "/login" });
      } else {
        toast.error(t("posts.new.error.generic"));
      }
    }
  }

  return { form, onSubmit, charCount };
}
