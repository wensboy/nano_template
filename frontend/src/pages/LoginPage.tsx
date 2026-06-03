import { useLocation, useNavigate } from "react-router-dom";

import AuthDialog from "@/components/internal/AuthDialog";

function isLocationStateWithFrom(state: unknown): state is { from: { pathname: string } } {
  return (
    typeof state === "object" &&
    state !== null &&
    "from" in state &&
    typeof state.from === "object" &&
    state.from !== null &&
    "pathname" in state.from &&
    typeof state.from.pathname === "string"
  );
}

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const redirectPath = isLocationStateWithFrom(location.state) ? location.state.from.pathname : "/";

  return (
    <main className="min-h-screen bg-zinc-100">
      <AuthDialog open onAuthenticated={() => navigate(redirectPath, { replace: true })} />
    </main>
  );
}
