import type { Meta, StoryObj } from "@storybook/react";
import { Label } from "./Label";

const meta = {
	title: "Components/Label",
	component: Label,
	args: {
		children: "Label",
	},
} satisfies Meta<typeof Label>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
