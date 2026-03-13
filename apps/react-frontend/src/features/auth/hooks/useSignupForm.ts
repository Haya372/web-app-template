import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'
import { callSignup } from '@/features/auth/api/signup'

function useSignupFormSchema() {
  const { t } = useTranslation()
  return z.object({
    name: z.string().min(1, t('signup.validation.nameRequired')),
    email: z.email(t('signup.validation.emailInvalid')),
    password: z.string().min(1, t('signup.validation.passwordRequired')),
  })
}

type SignupFormValues = z.infer<ReturnType<typeof useSignupFormSchema>>

export function useSignupForm() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const schema = useSignupFormSchema()
  const [errorMessage, setErrorMessage] = useState<string>('')

  const form = useForm<SignupFormValues>({
    resolver: zodResolver(schema),
    defaultValues: { name: '', email: '', password: '' },
  })

  async function onSubmit(values: SignupFormValues) {
    setErrorMessage('')
    try {
      await callSignup(values.name, values.email, values.password)
      navigate({ to: '/login' })
    } catch (err) {
      const message = err instanceof Error ? err.message : ''
      setErrorMessage(
        message === '409'
          ? t('signup.error.emailAlreadyRegistered')
          : t('signup.error.generic'),
      )
    }
  }

  return { form, onSubmit, errorMessage }
}
