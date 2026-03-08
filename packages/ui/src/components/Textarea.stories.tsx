import type { Meta, StoryObj } from "@storybook/react";
import { Label } from "./Label";
import { Textarea } from "./Textarea";

const meta = {
	title: "Components/Textarea",
	component: Textarea,
} satisfies Meta<typeof Textarea>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const WithPlaceholder: Story = {
	args: { placeholder: "Enter your message..." },
};

export const WithValue: Story = {
	args: { defaultValue: "Some pre-filled content." },
};

export const Disabled: Story = {
	args: { placeholder: "Disabled textarea", disabled: true },
};

export const WithLabel: Story = {
	render: () => (
		<div className="flex flex-col gap-1.5">
			<Label htmlFor="message">Message</Label>
			<Textarea id="message" placeholder="Enter your message..." />
		</div>
	),
};
