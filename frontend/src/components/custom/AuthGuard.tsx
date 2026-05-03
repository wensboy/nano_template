import { useEffect, useState } from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";

const tokenStorageKey = "token";
const refreshAttemptStorageKey = "auth:refresh-attempted";

type AuthStatus = "checking" | "authorized" | "unauthorized";

export default function AuthGuard() {
  const location = useLocation();
  const [authStatus, setAuthStatus] = useState<AuthStatus>("checking");

  useEffect(() => {
    const token = sessionStorage.getItem(tokenStorageKey);

    if (token) {
      sessionStorage.removeItem(refreshAttemptStorageKey);
      setAuthStatus("authorized");
      return;
    }

    const hasAttemptedRefresh = sessionStorage.getItem(refreshAttemptStorageKey) === "true";

    if (!hasAttemptedRefresh) {
      sessionStorage.setItem(refreshAttemptStorageKey, "true");
      window.location.reload();
      return;
    }

    sessionStorage.removeItem(refreshAttemptStorageKey);
    setAuthStatus("unauthorized");
  }, []);

  if (authStatus === "authorized") {
    return <Outlet />;
  }

  if (authStatus === "unauthorized") {
    return <Navigate replace state={{ from: location }} to="/login" />;
  }

  return null;
}

