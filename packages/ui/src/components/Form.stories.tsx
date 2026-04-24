import type { Meta, StoryObj } from "@storybook/react-vite";
import { useForm } from "react-hook-form";
import { Button } from "./Button";
import {
	Form,
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "./Form";
import { Input } from "./Input";

const meta = {
	title: "Components/Form",
	component: Form,
} satisfies Meta<typeof Form>;

export default meta;
type Story = StoryObj<typeof meta>;

function BasicFormExample() {
	const form = useForm({
		defaultValues: { email: "", password: "" },
	});
	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(() => {})} className="space-y-4 w-80">
				<FormField
					control={form.control}
					name="email"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Email</FormLabel>
							<FormControl>
								<Input type="email" placeholder="you@example.com" {...field} />
							</FormControl>
							<FormDescription>Enter your email address.</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>
				<FormField
					control={form.control}
					name="password"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Password</FormLabel>
							<FormControl>
								<Input type="password" placeholder="••••••••" {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
				<Button type="submit">Submit</Button>
			</form>
		</Form>
	);
}

export const Default: Story = {
	render: () => <BasicFormExample />,
};

function ValidationErrorExample() {
	const form = useForm({
		defaultValues: { email: "" },
		errors: {
			email: { type: "required", message: "Email is required." },
		},
	});
	return (
		<Form {...form}>
			<form className="space-y-4 w-80">
				<FormField
					control={form.control}
					name="email"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Email</FormLabel>
							<FormControl>
								<Input type="email" placeholder="you@example.com" {...field} />
							</FormControl>
							<FormMessage />
						</FormItem>
					)}
				/>
			</form>
		</Form>
	);
}

export const WithValidationError: Story = {
	render: () => <ValidationErrorExample />,
};

function WithDescriptionExample() {
	const form = useForm({ defaultValues: { username: "" } });
	return (
		<Form {...form}>
			<form className="space-y-4 w-80">
				<FormField
					control={form.control}
					name="username"
					render={({ field }) => (
						<FormItem>
							<FormLabel>Username</FormLabel>
							<FormControl>
								<Input placeholder="johndoe" {...field} />
							</FormControl>
							<FormDescription>
								This is your public display name.
							</FormDescription>
							<FormMessage />
						</FormItem>
					)}
				/>
			</form>
		</Form>
	);
}

export const WithDescription: Story = {
	render: () => <WithDescriptionExample />,
};
