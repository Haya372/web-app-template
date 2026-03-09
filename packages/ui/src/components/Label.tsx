import { Root as LabelRoot } from "@radix-ui/react-label";
import type { ComponentProps } from "react";

type LabelProps = Omit<ComponentProps<typeof LabelRoot>, "className">;

function Label({ ...props }: LabelProps) {
	return (
		<LabelRoot
			className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
			{...props}
		/>
	);
}

export { Label };
