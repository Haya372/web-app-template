import type { Meta, StoryObj } from "@storybook/react";
import { List, ListItem } from "./List";

const meta = {
	title: "Components/List",
	component: List,
} satisfies Meta<typeof List>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<List>
			<ListItem>First item</ListItem>
			<ListItem>Second item</ListItem>
			<ListItem>Third item</ListItem>
		</List>
	),
};

export const Empty: Story = {
	render: () => <List />,
};

export const SingleItem: Story = {
	render: () => (
		<List>
			<ListItem>Only item</ListItem>
		</List>
	),
};

export const AllVariants: Story = {
	render: () => (
		<div className="space-y-8">
			<List>
				<ListItem>Item one</ListItem>
				<ListItem>Item two</ListItem>
			</List>
		</div>
	),
};
