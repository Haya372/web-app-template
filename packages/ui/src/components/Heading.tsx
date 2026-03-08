import { cva, type VariantProps } from "class-variance-authority";
import type { ComponentProps } from "react";

const headingVariants = cva("font-semibold tracking-tight", {
	variants: {
		level: {
			1: "text-4xl",
			2: "text-3xl",
			3: "text-2xl",
			4: "text-xl",
		},
	},
	defaultVariants: {
		level: 1,
	},
});

type HeadingProps = Omit<ComponentProps<"h1">, "className"> &
	VariantProps<typeof headingVariants> & {
		level?: 1 | 2 | 3 | 4;
	};

function Heading({ level = 1, children, ...props }: HeadingProps) {
	const Tag = `h${level}` as "h1" | "h2" | "h3" | "h4";
	return (
		<Tag className={headingVariants({ level })} {...props}>
			{children}
		</Tag>
	);
}

export { Heading };
