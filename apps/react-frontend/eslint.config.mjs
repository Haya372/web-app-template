import { readdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { fixupPluginRules } from "@eslint/compat";
import importX from "eslint-plugin-import-x";
import jsxA11y from "eslint-plugin-jsx-a11y";
import reactPlugin from "eslint-plugin-react";
import reactHooks from "eslint-plugin-react-hooks";
import unicornPlugin from "eslint-plugin-unicorn";
import tseslint from "typescript-eslint";

const __dirname = dirname(fileURLToPath(import.meta.url));
const featuresDir = join(__dirname, "src/features");
const features = readdirSync(featuresDir);

/**
 * For each feature directory, restrict imports from every other feature.
 * This enforces the architecture rule: features must not depend on each other.
 * Shared code belongs in src/utils/, src/hooks/, or src/components/.
 */
const featureBoundaryConfigs = features.map((feature) => ({
	files: [`src/features/${feature}/**/*.{ts,tsx}`],
	rules: {
		"no-restricted-imports": [
			"error",
			{
				patterns: features
					.filter((f) => f !== feature)
					.map((otherFeature) => ({
						group: [
							`@/features/${otherFeature}`,
							`@/features/${otherFeature}/**`,
						],
						message: `Cross-feature import forbidden. Move shared code to src/utils/, src/hooks/, or src/components/.`,
					})),
			},
		],
	},
}));

export default tseslint.config(
	{ ignores: ["src/routeTree.gen.ts", "src/generated/**", "dist/**"] },

	// ── Base language options ──────────────────────────────────────────────────
	{
		files: ["src/**/*.{ts,tsx}"],
		languageOptions: {
			parser: tseslint.parser,
			parserOptions: {
				projectService: true,
				tsconfigRootDir: __dirname,
			},
		},
	},

	// ── React (flat/recommended) ───────────────────────────────────────────────
	// Provides rules for JSX correctness and React best practices.
	// eslint-plugin-react v7 uses legacy context APIs removed in ESLint v10;
	// fixupPluginRules() shims those APIs so the rules work unchanged.
	{
		plugins: {
			react: fixupPluginRules(reactPlugin),
		},
		rules: {
			...reactPlugin.configs.flat.recommended.rules,
			...reactPlugin.configs.flat["jsx-runtime"].rules, // React 17+ JSX transform — no `import React` needed
		},
		settings: {
			react: {
				version: "detect",
			},
		},
	},

	// ── React Hooks ───────────────────────────────────────────────────────────
	// Enforces the Rules of Hooks and exhaustive-deps.
	reactHooks.configs.flat["recommended-latest"],

	// ── JSX Accessibility (a11y) ──────────────────────────────────────────────
	// Enforces accessibility best practices for JSX elements.
	jsxA11y.flatConfigs.recommended,

	// ── Import-X ──────────────────────────────────────────────────────────────
	// Validates import paths and catches problematic import patterns.
	// Note: Resolution-based rules (no-unresolved, namespace, default) are disabled
	// because TypeScript (tsc) already validates all import paths and types.
	// The eslint-import-resolver-typescript peer dep is intentionally not installed
	// to keep the toolchain lean — tsc is the authoritative resolver.
	importX.configs["flat/recommended"],
	{
		rules: {
			// TypeScript already catches unresolved imports — disable to avoid duplicates.
			"import-x/no-unresolved": "off",
			"import-x/namespace": "off",
			"import-x/default": "off",
			"import-x/no-named-as-default": "off",
			"import-x/no-named-as-default-member": "off",

			// Covered by Biome's import sorting (assist.organizeImports).
			"import-x/order": "off",
			// Allow missing extensions for TS/TSX — TypeScript resolves them.
			"import-x/extensions": "off",
			// Biome enforces no barrel re-exports; ESLint rule too noisy here.
			"import-x/no-default-export": "off",
			// TanStack Router generates named exports — disable the named rule.
			"import-x/prefer-default-export": "off",
		},
	},

	// ── Unicorn (modern JS idioms) ────────────────────────────────────────────
	// Promotes modern JavaScript patterns and catches common pitfalls.
	unicornPlugin.configs["flat/recommended"],
	{
		rules: {
			// The project uses common abbreviations (props, fn, ref, etc.).
			// Requiring full names would hurt readability more than it helps.
			"unicorn/prevent-abbreviations": "off",

			// PascalCase component filenames are a React convention.
			// Enforce kebab-case only for non-component files via the pattern below.
			"unicorn/filename-case": [
				"error",
				{
					cases: {
						// Allow PascalCase (components, pages) and camelCase (hooks, utils, routes)
						pascalCase: true,
						camelCase: true,
						// Allow kebab-case for config/test files that may use it
						kebabCase: true,
					},
				},
			],

			// `null` is used in React (e.g., conditional rendering, API types).
			// Replacing all nulls with `undefined` would conflict with existing types.
			"unicorn/no-null": "off",

			// Array.forEach is idiomatic in this codebase and readable.
			"unicorn/no-array-for-each": "off",

			// `Array.from` vs spread is a style preference; Biome handles this.
			"unicorn/prefer-spread": "off",

			// TanStack Router / React patterns use `export default` for routes/components.
			"unicorn/no-anonymous-default-export": "off",

			// useReducer and similar hooks rely on switch statements; this rule is too strict.
			"unicorn/no-negated-condition": "off",

			// Ternary nesting is common in JSX; Biome handles formatting complexity.
			"unicorn/no-nested-ternary": "off",

			// Conflicts with TypeScript's `import type` which Biome already handles.
			"unicorn/prefer-module": "off",

			// Number constants are intentionally not extracted in all cases.
			"unicorn/numeric-separators-style": "off",

			// Covered by TypeScript itself — redundant.
			"unicorn/prefer-number-properties": "off",

			// This is a browser React app — `window` is the conventional way to access
			// browser globals (localStorage, matchMedia, etc.). `globalThis` is clearer
			// in isomorphic/Node contexts but adds noise in pure browser code.
			"unicorn/prefer-global-this": "off",

			// TypeScript declaration files use `export {}` to make a file a module.
			// This is an established TS pattern that unicorn incorrectly flags.
			"unicorn/require-module-specifiers": "off",
		},
	},

	// ── Formatting rules OFF (Biome is the formatter) ────────────────────────
	// Biome owns all formatting concerns: indentation, quotes, semicolons, etc.
	// Disable any ESLint rules that would duplicate or conflict with Biome.
	{
		rules: {
			indent: "off",
			quotes: "off",
			semi: "off",
			"comma-dangle": "off",
			"max-len": "off",
			"object-curly-spacing": "off",
			"arrow-parens": "off",
			"space-before-function-paren": "off",
		},
	},

	// ── Feature boundary imports (architecture enforcement) ───────────────────
	...featureBoundaryConfigs,
);
