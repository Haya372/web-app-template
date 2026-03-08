import type { Meta, StoryObj } from "@storybook/react";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectLabel,
	SelectGroup,
	SelectTrigger,
	SelectValue,
} from "./Select";

const meta = {
	title: "Components/Select",
	component: Select,
} satisfies Meta<typeof Select>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Select>
			<SelectTrigger>
				<SelectValue placeholder="Select an option" />
			</SelectTrigger>
			<SelectContent>
				<SelectItem value="apple">Apple</SelectItem>
				<SelectItem value="banana">Banana</SelectItem>
				<SelectItem value="cherry">Cherry</SelectItem>
			</SelectContent>
		</Select>
	),
};

export const WithGroups: Story = {
	render: () => (
		<Select>
			<SelectTrigger>
				<SelectValue placeholder="Select a fruit" />
			</SelectTrigger>
			<SelectContent>
				<SelectGroup>
					<SelectLabel>Fruits</SelectLabel>
					<SelectItem value="apple">Apple</SelectItem>
					<SelectItem value="banana">Banana</SelectItem>
				</SelectGroup>
				<SelectGroup>
					<SelectLabel>Vegetables</SelectLabel>
					<SelectItem value="carrot">Carrot</SelectItem>
					<SelectItem value="potato">Potato</SelectItem>
				</SelectGroup>
			</SelectContent>
		</Select>
	),
};

export const Disabled: Story = {
	render: () => (
		<Select disabled>
			<SelectTrigger>
				<SelectValue placeholder="Disabled" />
			</SelectTrigger>
			<SelectContent>
				<SelectItem value="apple">Apple</SelectItem>
			</SelectContent>
		</Select>
	),
};
