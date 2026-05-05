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
});
