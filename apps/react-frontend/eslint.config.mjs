import { readdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
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
	{ ignores: ["src/routeTree.gen.ts", "dist/**"] },
	{
		files: ["src/**/*.{ts,tsx}"],
		languageOptions: {
			parser: tseslint.parser,
		},
	},
	...featureBoundaryConfigs,
);
