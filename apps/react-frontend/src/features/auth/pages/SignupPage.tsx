import {
  Button,
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  Heading,
  Input,
} from '@repo/ui'
import { useTranslation } from 'react-i18next'
import { useSignupForm } from '@/features/auth/hooks/useSignupForm'

export function SignupPage() {
  const { t } = useTranslation()
  const { form, onSubmit, errorMessage } = useSignupForm()

  return (
    <main className="page-wrap flex min-h-[60vh] items-center justify-center px-4 py-12">
      <section className="island-shell w-full max-w-sm rounded-2xl p-6 sm:p-8">
        <Heading level={1}>{t('signup.title')}</Heading>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} noValidate>
            <div className="mb-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('signup.name.label')}</FormLabel>
                    <FormControl>
                      <Input
                        type="text"
                        placeholder={t('signup.name.placeholder')}
                        autoComplete="name"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="mb-4">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('signup.email.label')}</FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder={t('signup.email.placeholder')}
                        autoComplete="email"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="mb-6">
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('signup.password.label')}</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder={t('signup.password.placeholder')}
                        autoComplete="new-password"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {errorMessage && (
              <div className="mb-4">
                <p className="text-sm text-red-600">{errorMessage}</p>
              </div>
            )}

            <div className="w-full">
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting
                  ? t('signup.submitting')
                  : t('signup.submit')}
              </Button>
            </div>
          </form>
        </Form>
      </section>
    </main>
  )
}
