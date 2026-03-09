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
import { useLoginForm } from '@/features/auth/hooks/useLoginForm'

export function LoginPage() {
  const { t } = useTranslation()
  const { form, onSubmit } = useLoginForm()

  return (
    <main className="page-wrap flex min-h-[60vh] items-center justify-center px-4 py-12">
      <section className="island-shell w-full max-w-sm rounded-2xl p-6 sm:p-8">
        <Heading level={1}>{t('login.title')}</Heading>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} noValidate>
            <div className="mb-4">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('login.email.label')}</FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder={t('login.email.placeholder')}
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
                    <FormLabel>{t('login.password.label')}</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder={t('login.password.placeholder')}
                        autoComplete="current-password"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="w-full">
              <Button type="submit" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting
                  ? t('login.submitting')
                  : t('login.submit')}
              </Button>
            </div>
          </form>
        </Form>
      </section>
    </main>
  )
}
