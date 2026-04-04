/**
 * ESLint base configuration shared across all workspaces.
 *
 * This config turns off formatting rules that Biome owns.
 * Workspaces import this and spread it into their own config array.
 *
 * @example
 * import { formattingRulesOff } from "../../eslint.config.base.mjs";
 * export default tseslint.config(
 *   // ... workspace-specific configs ...
 *   formattingRulesOff,
 * );
 */

/** Disable ESLint formatting rules that conflict with Biome. */
export const formattingRulesOff = {
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
};
