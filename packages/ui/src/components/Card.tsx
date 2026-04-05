import type { ComponentProps } from "react";

function Card({ ...props }: Omit<ComponentProps<"div">, "className">) {
  return <div className="rounded-xl border bg-card text-card-foreground shadow" {...props} />;
}

function CardHeader({ ...props }: Omit<ComponentProps<"div">, "className">) {
  return <div className="flex flex-col space-y-1.5 p-6" {...props} />;
}

function CardTitle({ ...props }: Omit<ComponentProps<"div">, "className">) {
  return <div className="font-semibold leading-none tracking-tight" {...props} />;
}

function CardDescription({ ...props }: Omit<ComponentProps<"div">, "className">) {
  return <div className="text-sm text-muted-foreground" {...props} />;
}

function CardContent({ ...props }: Omit<ComponentProps<"div">, "className">) {
  return <div className="p-6 pt-0" {...props} />;
}

function CardFooter({ ...props }: Omit<ComponentProps<"div">, "className">) {
  return <div className="flex items-center p-6 pt-0" {...props} />;
}

export { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter };
