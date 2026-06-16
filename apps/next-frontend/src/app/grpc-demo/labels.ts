import { ServingStatus } from "@/generated/proto/health/v1/health_pb";

export function formatServingStatus(status: ServingStatus): string {
	switch (status) {
		case ServingStatus.SERVING: {
			return "Serving";
		}
		case ServingStatus.NOT_SERVING: {
			return "Not Serving";
		}
		default: {
			return "Unknown";
		}
	}
}
