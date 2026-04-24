import type { Meta, StoryObj } from "@storybook/react-vite";
import { Typography } from "./Typography";

const meta = {
	title: "Components/Typography",
	component: Typography,
	args: {
		children: "The quick brown fox jumps over the lazy dog.",
	},
} satisfies Meta<typeof Typography>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Paragraph: Story = {
	args: { variant: "p" },
};

export const Lead: Story = {
	args: { variant: "lead" },
};

export const Muted: Story = {
	args: { variant: "muted" },
};

export const Small: Story = {
	args: { variant: "small" },
};

export const Blockquote: Story = {
	args: {
		variant: "blockquote",
		children: "Design is not just what it looks like. Design is how it works.",
	},
};

export const AllVariants: Story = {
	render: () => (
		<div className="flex flex-col gap-4">
			<Typography variant="p">
				Paragraph: The quick brown fox jumps over the lazy dog.
			</Typography>
			<Typography variant="lead">
				Lead: The quick brown fox jumps over the lazy dog.
			</Typography>
			<Typography variant="muted">
				Muted: The quick brown fox jumps over the lazy dog.
			</Typography>
			<Typography variant="small">
				Small: The quick brown fox jumps over the lazy dog.
			</Typography>
			<Typography variant="blockquote">
				Blockquote: Design is not just what it looks like.
			</Typography>
		</div>
	),
};
