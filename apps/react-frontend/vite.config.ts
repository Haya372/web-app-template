import tailwindcss from "@tailwindcss/vite";
import viteReact from "@vitejs/plugin-react";
import { defineConfig } from "vite-plus";
import tsconfigPaths from "vite-tsconfig-paths";

export default defineConfig({
  plugins: [tsconfigPaths({ projects: ["./tsconfig.json"] }), tailwindcss(), viteReact()],
  server: {
    port: parseInt(process.env.VITE_PORT ?? "3000", 10),
  },
  test: {
    environment: "jsdom",
    globals: true,
    environmentOptions: {
      jsdom: {
        resources: "usable",
      },
    },
    setupFiles: ["./src/test-setup.ts"],
  },
  lint: {
    // NOTE: ESLint feature-boundary rules (src/features/<A>/ must not import from src/features/<B>/)
    // were removed as part of the Vite+ migration (#131). Oxlint enforcement is tracked in a follow-up issue.
    ignorePatterns: ["src/routeTree.gen.ts", "src/styles.css", "src/generated/**", "dist/**"],
  },
  fmt: {
    indentStyle: "tab",
    quotes: "double",
  },
  staged: {
    "*": "vp check --fix",
  },
});
