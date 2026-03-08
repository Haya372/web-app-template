import * as LabelPrimitive from "@radix-ui/react-label";
import type { ComponentProps } from "react";

type LabelProps = Omit<
	ComponentProps<typeof LabelPrimitive.Root>,
	"className"
>;

function Label({ ...props }: LabelProps) {
	return (
		<LabelPrimitive.Root
			className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
			{...props}
		/>
	);
}

export { Label };
