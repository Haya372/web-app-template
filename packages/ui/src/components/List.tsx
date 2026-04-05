import type { ComponentProps } from "react";

function List({ ...props }: Omit<ComponentProps<"ul">, "className">) {
  return <ul className="space-y-4" {...props} />;
}

function ListItem({ ...props }: Omit<ComponentProps<"li">, "className">) {
  return <li className="rounded-xl border bg-card text-card-foreground shadow p-4" {...props} />;
}

export { List, ListItem };
