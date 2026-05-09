import { fileURLToPath } from "node:url";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vitest/config";

function stub(path: string): string {
	return fileURLToPath(new URL(path, import.meta.url));
}

export default defineConfig({
	plugins: [react()],
	test: {
		environment: "jsdom",
		globals: true,
		setupFiles: ["./src/test-setup.ts"],
	},
	resolve: {
		alias: {
			"@": fileURLToPath(new URL("./src", import.meta.url)),
			// Stubs for packages added in Issue #106 (not yet installed via pnpm install).
			// vi.mock() factories in test files override these stubs at runtime.
			// Remove these aliases after running `pnpm install`.
			"@connectrpc/connect": stub("./src/__stubs__/connectrpc-connect.ts"),
			"@connectrpc/connect-node": stub(
				"./src/__stubs__/connectrpc-connect-node.ts",
			),
			// @repo/ui is a workspace package whose deps live in packages/ui/node_modules.
			// Remove after running `pnpm install`.
			"@repo/ui": fileURLToPath(
				new URL("../../packages/ui/src/index.ts", import.meta.url),
			),
		},
	},
});
