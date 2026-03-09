import { Root as LabelRoot } from "@radix-ui/react-label";
import { Slot } from "@radix-ui/react-slot";
import { createContext, useContext, useId, type ComponentProps } from "react";
import {
	Controller,
	FormProvider,
	useFormContext,
	type ControllerProps,
	type FieldPath,
	type FieldValues,
} from "react-hook-form";
import { cn } from "../lib/utils";

const Form = FormProvider;

interface FormFieldContextValue<
	TFieldValues extends FieldValues = FieldValues,
	TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
> {
	name: TName;
}

const FormFieldContext = createContext<FormFieldContextValue>(
	{} as FormFieldContextValue,
);

function FormField<
	TFieldValues extends FieldValues = FieldValues,
	TName extends FieldPath<TFieldValues> = FieldPath<TFieldValues>,
>({ ...props }: ControllerProps<TFieldValues, TName>) {
	return (
		<FormFieldContext.Provider value={{ name: props.name }}>
			<Controller {...props} />
		</FormFieldContext.Provider>
	);
}

function useFormField() {
	const fieldContext = useContext(FormFieldContext);
	const itemContext = useContext(FormItemContext);
	const { getFieldState, formState } = useFormContext();

	if (!fieldContext) {
		throw new Error("useFormField should be used within <FormField>");
	}

	const fieldState = getFieldState(fieldContext.name, formState);

	const { id } = itemContext;

	return {
		id,
		name: fieldContext.name,
		formItemId: `${id}-form-item`,
		formDescriptionId: `${id}-form-item-description`,
		formMessageId: `${id}-form-item-message`,
		...fieldState,
	};
}

interface FormItemContextValue {
	id: string;
}

const FormItemContext = createContext<FormItemContextValue>(
	{} as FormItemContextValue,
);

type FormItemProps = Omit<ComponentProps<"div">, "className">;

function FormItem({ ...props }: FormItemProps) {
	const id = useId();
	return (
		<FormItemContext.Provider value={{ id }}>
			<div className="space-y-2" {...props} />
		</FormItemContext.Provider>
	);
}

type FormLabelProps = Omit<ComponentProps<typeof LabelRoot>, "className">;

function FormLabel({ ...props }: FormLabelProps) {
	const { error, formItemId } = useFormField();
	return (
		<LabelRoot
			className={cn(
				"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
				error && "text-destructive",
			)}
			htmlFor={formItemId}
			{...props}
		/>
	);
}

type FormControlProps = Omit<ComponentProps<typeof Slot>, "className">;

function FormControl({ ...props }: FormControlProps) {
	const { error, formItemId, formDescriptionId, formMessageId } =
		useFormField();
	return (
		<Slot
			id={formItemId}
			aria-describedby={
				!error
					? `${formDescriptionId}`
					: `${formDescriptionId} ${formMessageId}`
			}
			aria-invalid={!!error}
			{...props}
		/>
	);
}

type FormDescriptionProps = Omit<ComponentProps<"p">, "className">;

function FormDescription({ ...props }: FormDescriptionProps) {
	const { formDescriptionId } = useFormField();
	return (
		<p
			id={formDescriptionId}
			className="text-sm text-muted-foreground"
			{...props}
		/>
	);
}

type FormMessageProps = Omit<ComponentProps<"p">, "className">;

function FormMessage({ children, ...props }: FormMessageProps) {
	const { error, formMessageId } = useFormField();
	const body = error ? String(error?.message ?? "") : children;

	if (!body) return null;

	return (
		<p
			id={formMessageId}
			className="text-sm font-medium text-destructive"
			{...props}
		>
			{body}
		</p>
	);
}

export {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
	useFormField,
};
