import { describe, expect, it } from "vitest";
import { ServingStatus } from "@/generated/proto/health/v1/health_pb";
import { formatServingStatus } from "./labels";

describe("formatServingStatus", () => {
	it("returns 'Serving' for ServingStatus.SERVING", () => {
		expect(formatServingStatus(ServingStatus.SERVING)).toBe("Serving");
	});

	it("returns 'Not Serving' for ServingStatus.NOT_SERVING", () => {
		expect(formatServingStatus(ServingStatus.NOT_SERVING)).toBe("Not Serving");
	});

	it("returns 'Unknown' for ServingStatus.UNSPECIFIED", () => {
		expect(formatServingStatus(ServingStatus.UNSPECIFIED)).toBe("Unknown");
	});

	it("returns 'Unknown' for an unrecognized numeric value", () => {
		expect(formatServingStatus(99 as ServingStatus)).toBe("Unknown");
	});
});
