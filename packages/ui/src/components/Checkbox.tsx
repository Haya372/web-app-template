import * as CheckboxPrimitive from "@radix-ui/react-checkbox";
import { Check } from "lucide-react";
import type { ComponentProps } from "react";

type CheckboxProps = Omit<
	ComponentProps<typeof CheckboxPrimitive.Root>,
	"className"
>;

function Checkbox({ ...props }: CheckboxProps) {
	return (
		<CheckboxPrimitive.Root
			className="peer h-4 w-4 shrink-0 rounded-sm border border-primary shadow focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground"
			{...props}
		>
			<CheckboxPrimitive.Indicator className="flex items-center justify-center text-current">
				<Check className="h-4 w-4" />
			</CheckboxPrimitive.Indicator>
		</CheckboxPrimitive.Root>
	);
}

export { Checkbox };
