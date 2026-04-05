import { cva, type VariantProps } from "class-variance-authority";
import type { HTMLAttributes } from "react";

const typographyVariants = cva("", {
  variants: {
    variant: {
      p: "text-base leading-7",
      lead: "text-xl text-muted-foreground",
      muted: "text-sm text-muted-foreground",
      small: "text-sm font-medium leading-none",
      blockquote: "border-l-2 border-border pl-4 italic text-muted-foreground",
    },
  },
  defaultVariants: {
    variant: "p",
  },
});

const tagMap = {
  p: "p",
  lead: "p",
  muted: "p",
  small: "small",
  blockquote: "blockquote",
} as const;

type TypographyVariant = NonNullable<VariantProps<typeof typographyVariants>["variant"]>;

type TypographyProps = Omit<HTMLAttributes<HTMLElement>, "className"> & {
  variant?: TypographyVariant;
};

function Typography({ variant = "p", ...props }: TypographyProps) {
  const Tag = tagMap[variant];
  return <Tag className={typographyVariants({ variant })} {...props} />;
}

export { Typography };
