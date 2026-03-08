import type { Meta, StoryObj } from "@storybook/react";
import { Button } from "./Button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "./Dialog";

const meta = {
	title: "Components/Dialog",
	component: Dialog,
} satisfies Meta<typeof Dialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<Dialog>
			<DialogTrigger asChild>
				<Button variant="outline">Open Dialog</Button>
			</DialogTrigger>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>Dialog Title</DialogTitle>
					<DialogDescription>
						This is a dialog description. It provides context for the dialog
						content.
					</DialogDescription>
				</DialogHeader>
				<p className="text-sm">Dialog body content goes here.</p>
				<DialogFooter>
					<Button variant="outline">Cancel</Button>
					<Button>Confirm</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	),
};
