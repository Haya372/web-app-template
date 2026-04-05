import * as TabsPrimitive from "@radix-ui/react-tabs";
import type { ComponentProps } from "react";

const Tabs = TabsPrimitive.Root;

type TabsListProps = Omit<ComponentProps<typeof TabsPrimitive.List>, "className">;

function TabsList({ ...props }: TabsListProps) {
  return (
    <TabsPrimitive.List
      className="inline-flex h-9 items-center justify-center rounded-lg bg-muted p-1 text-muted-foreground"
      {...props}
    />
  );
}

type TabsTriggerProps = Omit<ComponentProps<typeof TabsPrimitive.Trigger>, "className">;

function TabsTrigger({ ...props }: TabsTriggerProps) {
  return (
    <TabsPrimitive.Trigger
      className="inline-flex items-center justify-center whitespace-nowrap rounded-md px-3 py-1 text-sm font-medium ring-offset-background transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow"
      {...props}
    />
  );
}

type TabsContentProps = Omit<ComponentProps<typeof TabsPrimitive.Content>, "className">;

function TabsContent({ ...props }: TabsContentProps) {
  return (
    <TabsPrimitive.Content
      className="mt-2 ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
      {...props}
    />
  );
}

export { Tabs, TabsList, TabsTrigger, TabsContent };
