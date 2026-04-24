import type { Meta, StoryObj } from "@storybook/react-vite";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./Tabs";

const meta = {
	title: "Components/Tabs",
	component: Tabs,
} satisfies Meta<typeof Tabs>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Tabs defaultValue="tab1">
			<TabsList>
				<TabsTrigger value="tab1">Tab 1</TabsTrigger>
				<TabsTrigger value="tab2">Tab 2</TabsTrigger>
				<TabsTrigger value="tab3">Tab 3</TabsTrigger>
			</TabsList>
			<TabsContent value="tab1">Content for Tab 1</TabsContent>
			<TabsContent value="tab2">Content for Tab 2</TabsContent>
			<TabsContent value="tab3">Content for Tab 3</TabsContent>
		</Tabs>
	),
};

export const WithDisabled: Story = {
	render: () => (
		<Tabs defaultValue="tab1">
			<TabsList>
				<TabsTrigger value="tab1">Active</TabsTrigger>
				<TabsTrigger value="tab2" disabled>
					Disabled
				</TabsTrigger>
			</TabsList>
			<TabsContent value="tab1">Active tab content</TabsContent>
			<TabsContent value="tab2">Disabled tab content</TabsContent>
		</Tabs>
	),
};
