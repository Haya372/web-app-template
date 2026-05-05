import { dirname } from "node:path";
import { fileURLToPath } from "node:url";
import { fixupPluginRules } from "@eslint/compat";
import nextPlugin from "@next/eslint-plugin-next";
import importX from "eslint-plugin-import-x";
import jsxA11y from "eslint-plugin-jsx-a11y";
import reactPlugin from "eslint-plugin-react";
import reactHooks from "eslint-plugin-react-hooks";
import unicornPlugin from "eslint-plugin-unicorn";
import tseslint from "typescript-eslint";
import { formattingRulesOff } from "../../eslint.config.base.mjs";

const __dirname = dirname(fileURLToPath(import.meta.url));

const featureBoundaryConfigs = [];

export default tseslint.config(
	{ ignores: [".next/**", "next-env.d.ts", "dist/**"] },

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

	// ── Next.js ───────────────────────────────────────────────────────────────
	// @next/eslint-plugin-next uses legacy plugin context APIs removed in ESLint v10;
	// fixupPluginRules() shims those APIs so the rules work unchanged.
	{
		plugins: {
			"@next/next": fixupPluginRules(nextPlugin),
		},
		rules: {
			...nextPlugin.configs.recommended.rules,
			...nextPlugin.configs["core-web-vitals"].rules,
		},
	},

	// ── React (flat/recommended) ───────────────────────────────────────────────
	{
		plugins: {
			react: fixupPluginRules(reactPlugin),
		},
		rules: {
			...reactPlugin.configs.flat.recommended.rules,
			...reactPlugin.configs.flat["jsx-runtime"].rules,
		},
		settings: {
			react: {
				version: "detect",
			},
		},
	},

	// ── React Hooks ───────────────────────────────────────────────────────────
	reactHooks.configs.flat["recommended-latest"],

	// ── JSX Accessibility (a11y) ──────────────────────────────────────────────
	jsxA11y.flatConfigs.recommended,

	// ── Import-X ──────────────────────────────────────────────────────────────
	importX.configs["flat/recommended"],
	{
		rules: {
			"import-x/no-unresolved": "off",
			"import-x/namespace": "off",
			"import-x/default": "off",
			"import-x/no-named-as-default": "off",
			"import-x/no-named-as-default-member": "off",
			"import-x/order": "off",
			"import-x/extensions": "off",
			"import-x/no-default-export": "off",
			"import-x/prefer-default-export": "off",
		},
	},

	// ── Unicorn (modern JS idioms) ────────────────────────────────────────────
	unicornPlugin.configs["flat/recommended"],
	{
		rules: {
			"unicorn/prevent-abbreviations": "off",
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
			"unicorn/no-null": "off",
			"unicorn/no-array-for-each": "off",
			"unicorn/prefer-spread": "off",
			"unicorn/no-anonymous-default-export": "off",
			"unicorn/no-negated-condition": "off",
			"unicorn/no-nested-ternary": "off",
			"unicorn/prefer-module": "off",
			"unicorn/numeric-separators-style": "off",
			"unicorn/prefer-number-properties": "off",
			"unicorn/prefer-global-this": "off",
			"unicorn/require-module-specifiers": "off",
		},
	},

	// ── Formatting rules OFF (Biome is the formatter) ────────────────────────
	formattingRulesOff,

	// ── Feature boundary imports (architecture enforcement) ───────────────────
	...featureBoundaryConfigs,
);
