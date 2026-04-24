import type { Meta, StoryObj } from "@storybook/react-vite";
import { Heading } from "./Heading";

const meta = {
	title: "Components/Heading",
	component: Heading,
	args: {
		children: "The quick brown fox",
	},
} satisfies Meta<typeof Heading>;

export default meta;
type Story = StoryObj<typeof meta>;

export const H1: Story = {
	args: { level: 1 },
};

export const H2: Story = {
	args: { level: 2 },
};

export const H3: Story = {
	args: { level: 3 },
};

export const H4: Story = {
	args: { level: 4 },
};

export const AllLevels: Story = {
	render: () => (
		<div className="flex flex-col gap-4">
			<Heading level={1}>Heading 1</Heading>
			<Heading level={2}>Heading 2</Heading>
			<Heading level={3}>Heading 3</Heading>
			<Heading level={4}>Heading 4</Heading>
		</div>
	),
};
