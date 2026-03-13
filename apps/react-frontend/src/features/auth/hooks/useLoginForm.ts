import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { toast } from '@repo/ui'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { callLogin } from '@/features/auth/api/login'
import { saveToken } from '@/features/auth/utils/tokenStorage'

function useLoginFormSchema() {
  const { t } = useTranslation()
  return z.object({
    email: z.email(t('login.validation.emailInvalid')),
    password: z.string().min(1, t('login.validation.passwordRequired')),
  })
}

type LoginFormValues = z.infer<ReturnType<typeof useLoginFormSchema>>

export function useLoginForm() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const schema = useLoginFormSchema()

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: '', password: '' },
  })

  async function onSubmit(values: LoginFormValues) {
    try {
      const response = await callLogin(values.email, values.password)
      saveToken(response.token)
      navigate({ to: '/' })
    } catch (err) {
      const message = err instanceof Error ? err.message : ''
      toast.error(
        message === '401'
          ? t('login.error.invalidCredentials')
          : t('login.error.generic'),
      )
    }
  }

  return { form, onSubmit }
}
