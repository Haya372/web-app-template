import { Link } from "@tanstack/react-router";
import ThemeToggle from "./ThemeToggle";

export default function Header() {
  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/80 px-4 backdrop-blur">
      <nav className="page-wrap flex items-center gap-4 py-3">
        <Link to="/" className="text-sm font-semibold text-foreground no-underline">
          App
        </Link>

        <div className="flex items-center gap-4 text-sm font-medium">
          <Link
            to="/"
            className="text-muted-foreground transition-colors hover:text-foreground [&.active]:text-foreground"
            activeProps={{ className: "active" }}
          >
            Home
          </Link>
          <Link
            to="/login"
            className="text-muted-foreground transition-colors hover:text-foreground [&.active]:text-foreground"
            activeProps={{ className: "active" }}
          >
            Login
          </Link>
        </div>

        <div className="ml-auto">
          <ThemeToggle />
        </div>
      </nav>
    </header>
  );
}
