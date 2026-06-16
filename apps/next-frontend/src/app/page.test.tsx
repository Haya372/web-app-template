import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import HomePage from "./page";

describe("HomePage", () => {
	it("renders the title heading", () => {
		render(<HomePage />);
		expect(
			screen.getByRole("heading", { level: 1, name: /next\.js frontend/i }),
		).toBeInTheDocument();
	});

	it("renders a link to the grpc-demo page", () => {
		render(<HomePage />);
		expect(screen.getByRole("link", { name: /grpc demo/i })).toHaveAttribute(
			"href",
			"/grpc-demo",
		);
	});
});
