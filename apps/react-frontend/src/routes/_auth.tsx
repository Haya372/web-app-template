import { Outlet, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_auth")({
  component: AuthLayout,
});

export function AuthLayout() {
  // Minimal pass-through: auth pages (login, signup) need no chrome.
  // A future authentication guard (beforeLoad redirect) will be added here in a separate issue.
  return <Outlet />;
}
