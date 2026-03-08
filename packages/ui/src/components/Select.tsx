import * as SelectPrimitive from "@radix-ui/react-select";
import { Check, ChevronDown, ChevronUp } from "lucide-react";
import type { ComponentProps } from "react";

const Select = SelectPrimitive.Root;
const SelectGroup = SelectPrimitive.Group;
const SelectValue = SelectPrimitive.Value;

type SelectTriggerProps = Omit<
	ComponentProps<typeof SelectPrimitive.Trigger>,
	"className"
>;

function SelectTrigger({ children, ...props }: SelectTriggerProps) {
	return (
		<SelectPrimitive.Trigger
			className="flex h-9 w-full items-center justify-between whitespace-nowrap rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-ring disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1"
			{...props}
		>
			{children}
			<SelectPrimitive.Icon asChild>
				<ChevronDown className="h-4 w-4 opacity-50" />
			</SelectPrimitive.Icon>
		</SelectPrimitive.Trigger>
	);
}

type SelectScrollUpButtonProps = Omit<
	ComponentProps<typeof SelectPrimitive.ScrollUpButton>,
	"className"
>;

function SelectScrollUpButton({ ...props }: SelectScrollUpButtonProps) {
	return (
		<SelectPrimitive.ScrollUpButton
			className="flex cursor-default items-center justify-center py-1"
			{...props}
		>
			<ChevronUp className="h-4 w-4" />
		</SelectPrimitive.ScrollUpButton>
	);
}

type SelectScrollDownButtonProps = Omit<
	ComponentProps<typeof SelectPrimitive.ScrollDownButton>,
	"className"
>;

function SelectScrollDownButton({ ...props }: SelectScrollDownButtonProps) {
	return (
		<SelectPrimitive.ScrollDownButton
			className="flex cursor-default items-center justify-center py-1"
			{...props}
		>
			<ChevronDown className="h-4 w-4" />
		</SelectPrimitive.ScrollDownButton>
	);
}

type SelectContentProps = Omit<
	ComponentProps<typeof SelectPrimitive.Content>,
	"className"
>;

function SelectContent({
	children,
	position = "popper",
	...props
}: SelectContentProps) {
	return (
		<SelectPrimitive.Portal>
			<SelectPrimitive.Content
				className="relative z-50 max-h-96 min-w-[8rem] overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2"
				position={position}
				{...props}
			>
				<SelectScrollUpButton />
				<SelectPrimitive.Viewport
					className={
						position === "popper"
							? "h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-content-available-width)] p-1"
							: "p-1"
					}
				>
					{children}
				</SelectPrimitive.Viewport>
				<SelectScrollDownButton />
			</SelectPrimitive.Content>
		</SelectPrimitive.Portal>
	);
}

type SelectLabelProps = Omit<
	ComponentProps<typeof SelectPrimitive.Label>,
	"className"
>;

function SelectLabel({ ...props }: SelectLabelProps) {
	return (
		<SelectPrimitive.Label
			className="px-2 py-1.5 text-xs font-semibold"
			{...props}
		/>
	);
}

type SelectItemProps = Omit<
	ComponentProps<typeof SelectPrimitive.Item>,
	"className"
>;

function SelectItem({ children, ...props }: SelectItemProps) {
	return (
		<SelectPrimitive.Item
			className="relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-2 pr-8 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
			{...props}
		>
			<span className="absolute right-2 flex h-3.5 w-3.5 items-center justify-center">
				<SelectPrimitive.ItemIndicator>
					<Check className="h-4 w-4" />
				</SelectPrimitive.ItemIndicator>
			</span>
			<SelectPrimitive.ItemText>{children}</SelectPrimitive.ItemText>
		</SelectPrimitive.Item>
	);
}

type SelectSeparatorProps = Omit<
	ComponentProps<typeof SelectPrimitive.Separator>,
	"className"
>;

function SelectSeparator({ ...props }: SelectSeparatorProps) {
	return (
		<SelectPrimitive.Separator
			className="-mx-1 my-1 h-px bg-muted"
			{...props}
		/>
	);
}

export {
	Select,
	SelectGroup,
	SelectValue,
	SelectTrigger,
	SelectContent,
	SelectLabel,
	SelectItem,
	SelectSeparator,
	SelectScrollUpButton,
	SelectScrollDownButton,
};
