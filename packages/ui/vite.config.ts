import { defineConfig } from "vite-plus";

export default defineConfig({
  lint: {
    ignorePatterns: ["dist/**"],
  },
  fmt: {
    indentStyle: "tab",
    quotes: "double",
  },
  staged: {
    "*": "vp check --fix",
  },
});
