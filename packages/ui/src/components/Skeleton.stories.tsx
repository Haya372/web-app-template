import type { Meta, StoryObj } from "@storybook/react";
import { Skeleton } from "./Skeleton";

const meta = {
	title: "Components/Skeleton",
	component: Skeleton,
} satisfies Meta<typeof Skeleton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: { style: { width: 200, height: 20 } },
};

export const CardSkeleton: Story = {
	render: () => (
		<div className="flex flex-col gap-2 w-64">
			<Skeleton style={{ height: 160 }} />
			<Skeleton style={{ height: 16, width: "75%" }} />
			<Skeleton style={{ height: 16, width: "50%" }} />
		</div>
	),
};

export const AvatarSkeleton: Story = {
	render: () => (
		<div className="flex items-center gap-3">
			<Skeleton style={{ width: 40, height: 40, borderRadius: "50%" }} />
			<div className="flex flex-col gap-1.5">
				<Skeleton style={{ width: 120, height: 14 }} />
				<Skeleton style={{ width: 80, height: 14 }} />
			</div>
		</div>
	),
};
