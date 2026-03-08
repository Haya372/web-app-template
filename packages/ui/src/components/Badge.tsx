import { type VariantProps, cva } from "class-variance-authority";
import type { ComponentProps } from "react";

const badgeVariants = cva(
	"inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold",
	{
		variants: {
			variant: {
				default:
					"border-transparent bg-primary text-primary-foreground shadow",
				secondary:
					"border-transparent bg-secondary text-secondary-foreground",
				destructive:
					"border-transparent bg-destructive text-destructive-foreground shadow",
				outline: "text-foreground",
			},
		},
		defaultVariants: {
			variant: "default",
		},
	},
);

type BadgeProps = Omit<ComponentProps<"div">, "className"> &
	VariantProps<typeof badgeVariants>;

function Badge({ variant, ...props }: BadgeProps) {
	return <div className={badgeVariants({ variant })} {...props} />;
}

export { Badge, badgeVariants };
