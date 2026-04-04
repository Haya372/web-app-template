import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { postV1UsersSignup } from "@/generated/sdk.gen";

function useSignupFormSchema() {
	const { t } = useTranslation();
	return z.object({
		name: z.string().min(1, t("signup.validation.nameRequired")),
		email: z.email(t("signup.validation.emailInvalid")),
		password: z.string().min(1, t("signup.validation.passwordRequired")),
	});
}

type SignupFormValues = z.infer<ReturnType<typeof useSignupFormSchema>>;

export function useSignupForm() {
	const { t } = useTranslation();
	const navigate = useNavigate();
	const schema = useSignupFormSchema();
	const [errorMessage, setErrorMessage] = useState<string>("");

	const form = useForm<SignupFormValues>({
		resolver: zodResolver(schema),
		defaultValues: { name: "", email: "", password: "" },
	});

	async function onSubmit(values: SignupFormValues) {
		setErrorMessage("");
		try {
			const { data, error, response } = await postV1UsersSignup({
				body: {
					name: values.name,
					email: values.email,
					password: values.password,
				},
				baseUrl: import.meta.env.VITE_API_BASE_URL,
			});
			if (error || !data) throw new Error(String(response.status));
			navigate({ to: "/login" });
		} catch (error) {
			const message = error instanceof Error ? error.message : "";
			setErrorMessage(
				message === "409"
					? t("signup.error.emailAlreadyRegistered")
					: t("signup.error.generic"),
			);
		}
	}

	return { form, onSubmit, errorMessage };
}
