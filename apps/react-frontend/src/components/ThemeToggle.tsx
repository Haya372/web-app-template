import { useEffect, useState } from "react";

type ThemeMode = "light" | "dark" | "system";

function getStoredMode(): ThemeMode {
	if (typeof window === "undefined") return "system";
	const stored = window.localStorage.getItem("theme");
	if (stored === "light" || stored === "dark") return stored;
	return "system";
}

function applyTheme(mode: ThemeMode) {
	const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
	const isDark = mode === "dark" || (mode === "system" && prefersDark);
	document.documentElement.classList.toggle("dark", isDark);
	document.documentElement.style.colorScheme = isDark ? "dark" : "light";
}

export default function ThemeToggle() {
	// Lazy initializer reads stored preference once on mount — avoids calling
	// setState synchronously inside an effect (react-hooks/set-state-in-effect).
	const [mode, setMode] = useState<ThemeMode>(() => getStoredMode());

	useEffect(() => {
		applyTheme(mode);
	}, [mode]);

	useEffect(() => {
		if (mode !== "system") return;
		const media = window.matchMedia("(prefers-color-scheme: dark)");
		const onChange = () => applyTheme("system");
		media.addEventListener("change", onChange);
		return () => media.removeEventListener("change", onChange);
	}, [mode]);

	function cycle() {
		const next: ThemeMode =
			mode === "light" ? "dark" : mode === "dark" ? "system" : "light";
		setMode(next);
		applyTheme(next);
		window.localStorage.setItem("theme", next);
	}

	const label = `Theme: ${mode}. Click to switch.`;

	return (
		<button
			type="button"
			onClick={cycle}
			aria-label={label}
			title={label}
			className="rounded-md border border-border bg-background px-3 py-1.5 text-sm font-medium text-foreground transition-colors hover:bg-accent hover:text-accent-foreground"
		>
			{mode === "system" ? "System" : mode === "dark" ? "Dark" : "Light"}
		</button>
	);
}
