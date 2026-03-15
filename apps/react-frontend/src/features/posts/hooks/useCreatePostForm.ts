import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { toast } from '@repo/ui'
import { useEffect } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { callCreatePost } from '@/features/posts/api/createPost'
import { getToken } from '@/utils/tokenStorage'

function useCreatePostFormSchema() {
  const { t } = useTranslation()
  return z.object({
    content: z
      .string()
      .min(1, t('posts.new.validation.contentRequired'))
      .max(280, t('posts.new.validation.contentTooLong')),
  })
}

type CreatePostFormValues = z.infer<ReturnType<typeof useCreatePostFormSchema>>

export function useCreatePostForm() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const schema = useCreatePostFormSchema()

  const form = useForm<CreatePostFormValues>({
    resolver: zodResolver(schema),
    defaultValues: { content: '' },
  })

  // Redirect to login if no token is present on mount.
  // getToken() reads localStorage synchronously and is not reactive, so it is
  // intentionally omitted from the dependency array.
  useEffect(() => {
    if (getToken() === null) {
      void navigate({ to: '/login' })
    }
  }, [navigate])

  // useWatch scopes re-renders to this consumer only, avoiding a full form
  // re-render on every keystroke that form.watch() would cause.
  // useWatch returns undefined before RHF initialises its internal state on the
  // first render. The nullish fallback ensures charCount is always a number.
  const content = useWatch({ control: form.control, name: 'content' }) ?? ''
  const charCount = content.length

  async function onSubmit(values: CreatePostFormValues) {
    try {
      await callCreatePost(values.content)
      toast.success(t('posts.new.success'))
      form.reset()
    } catch (err) {
      // callCreatePost throws new Error(String(response.status)) for HTTP errors.
      // Only a server-returned 401 redirects to login; all other failures show
      // a generic toast. The pre-flight "Unauthenticated" error from
      // callCreatePost (no token) is already prevented by the mount guard above,
      // so it does not need special handling here.
      if (err instanceof Error && err.message === '401') {
        void navigate({ to: '/login' })
      } else {
        toast.error(t('posts.new.error.generic'))
      }
    }
  }

  return { form, onSubmit, charCount }
}
