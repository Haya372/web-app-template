import type { Meta, StoryObj } from "@storybook/react";
import { Input } from "./Input";
import { Label } from "./Label";

const meta = {
	title: "Components/Input",
	component: Input,
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithPlaceholder: Story = {
	args: { placeholder: "Enter text..." },
};

export const WithValue: Story = {
	args: { defaultValue: "Hello, world!" },
};

export const Disabled: Story = {
	args: { placeholder: "Disabled input", disabled: true },
};

export const Password: Story = {
	args: { type: "password", placeholder: "Enter password" },
};

export const WithLabel: Story = {
	render: () => (
		<div className="flex flex-col gap-1.5">
			<Label htmlFor="email">Email</Label>
			<Input id="email" type="email" placeholder="you@example.com" />
		</div>
	),
};
