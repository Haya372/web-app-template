import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "@repo/ui";
import { useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { useAuth } from "@/features/auth/hooks/useAuth";
import { postV1UsersLogin } from "@/generated/sdk.gen";

function useLoginFormSchema() {
	const { t } = useTranslation();
	return z.object({
		email: z.email(t("login.validation.emailInvalid")),
		password: z.string().min(1, t("login.validation.passwordRequired")),
	});
}

type LoginFormValues = z.infer<ReturnType<typeof useLoginFormSchema>>;

export function useLoginForm() {
	const { t } = useTranslation();
	const navigate = useNavigate();
	const schema = useLoginFormSchema();
	const { login } = useAuth();

	const form = useForm<LoginFormValues>({
		resolver: zodResolver(schema),
		defaultValues: { email: "", password: "" },
	});

	async function onSubmit(values: LoginFormValues) {
		try {
			const { data, error, response } = await postV1UsersLogin({
				body: { email: values.email, password: values.password },
				baseUrl: import.meta.env.VITE_API_BASE_URL,
			});
			if (error || !data) throw new Error(String(response.status));
			login(data.token);
			navigate({ to: "/" });
		} catch (error) {
			const message = error instanceof Error ? error.message : "";
			toast.error(
				message === "401"
					? t("login.error.invalidCredentials")
					: t("login.error.generic"),
			);
		}
	}

	return { form, onSubmit };
}
