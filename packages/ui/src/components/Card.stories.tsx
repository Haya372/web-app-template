import type { Meta, StoryObj } from "@storybook/react";
import { Button } from "./Button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "./Card";

const meta = {
  title: "Components/Card",
  component: Card,
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  render: () => (
    <Card>
      <CardHeader>
        <CardTitle>Card Title</CardTitle>
        <CardDescription>Card description goes here.</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Card content goes here.</p>
      </CardContent>
      <CardFooter>
        <Button>Action</Button>
      </CardFooter>
    </Card>
  ),
};

export const Simple: Story = {
  render: () => (
    <Card>
      <CardContent>
        <p>A simple card with only content.</p>
      </CardContent>
    </Card>
  ),
};
