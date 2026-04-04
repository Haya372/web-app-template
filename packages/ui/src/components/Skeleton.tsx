import type { ComponentProps } from "react";

type SkeletonProps = Omit<ComponentProps<"div">, "className">;

function Skeleton({ ...props }: SkeletonProps) {
	return <div className="animate-pulse rounded-md bg-primary/10" {...props} />;
}

export { Skeleton };
