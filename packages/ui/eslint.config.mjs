import { dirname } from "node:path";
import { fileURLToPath } from "node:url";
import { fixupPluginRules } from "@eslint/compat";
import importX from "eslint-plugin-import-x";
import jsxA11y from "eslint-plugin-jsx-a11y";
import reactPlugin from "eslint-plugin-react";
import reactHooks from "eslint-plugin-react-hooks";
import unicornPlugin from "eslint-plugin-unicorn";
import tseslint from "typescript-eslint";
import { formattingRulesOff } from "../../eslint.config.base.mjs";

const __dirname = dirname(fileURLToPath(import.meta.url));

export default tseslint.config(
	{ ignores: ["dist/**", "storybook-static/**"] },

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
	// Note: Resolution-based rules are disabled because TypeScript already
	// validates all import paths and types.
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
			// Story files use default exports (Storybook meta) — disable the named rule.
			"import-x/prefer-default-export": "off",
		},
	},

	// ── Unicorn (modern JS idioms) ────────────────────────────────────────────
	// Promotes modern JavaScript patterns and catches common pitfalls.
	unicornPlugin.configs["flat/recommended"],
	{
		rules: {
			// The project uses common abbreviations (props, fn, ref, etc.).
			"unicorn/prevent-abbreviations": "off",

			// PascalCase component filenames are a React convention.
			"unicorn/filename-case": [
				"error",
				{
					cases: {
						pascalCase: true,
						camelCase: true,
						kebabCase: true,
					},
				},
			],

			// `null` is used in React (e.g., conditional rendering, Radix UI).
			"unicorn/no-null": "off",

			// Array.forEach is idiomatic in this codebase and readable.
			"unicorn/no-array-for-each": "off",

			// `Array.from` vs spread is a style preference; Biome handles this.
			"unicorn/prefer-spread": "off",

			// Story files use default exports (Storybook meta pattern).
			"unicorn/no-anonymous-default-export": "off",

			// CVA and Radix patterns use switch statements; this rule is too strict.
			"unicorn/no-negated-condition": "off",

			// Ternary nesting is common in JSX; Biome handles formatting complexity.
			"unicorn/no-nested-ternary": "off",

			// Conflicts with TypeScript's `import type` which Biome already handles.
			"unicorn/prefer-module": "off",

			// Number constants are intentionally not extracted in all cases.
			"unicorn/numeric-separators-style": "off",

			// Covered by TypeScript itself — redundant.
			"unicorn/prefer-number-properties": "off",

			// This is a browser React component library — `window` is idiomatic.
			"unicorn/prefer-global-this": "off",

			// TypeScript declaration files use `export {}` to make a file a module.
			"unicorn/require-module-specifiers": "off",
		},
	},

	// ── Formatting rules OFF (Biome is the formatter) ────────────────────────
	// Biome owns all formatting concerns: indentation, quotes, semicolons, etc.
	formattingRulesOff,
);
