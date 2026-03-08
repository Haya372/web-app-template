import * as SeparatorPrimitive from "@radix-ui/react-separator";
import type { ComponentProps } from "react";

type SeparatorProps = Omit<
	ComponentProps<typeof SeparatorPrimitive.Root>,
	"className"
>;

function Separator({
	orientation = "horizontal",
	decorative = true,
	...props
}: SeparatorProps) {
	return (
		<SeparatorPrimitive.Root
			decorative={decorative}
			orientation={orientation}
			className={
				orientation === "horizontal"
					? "shrink-0 bg-border h-[1px] w-full"
					: "shrink-0 bg-border w-[1px]"
			}
			{...props}
		/>
	);
}

export { Separator };
