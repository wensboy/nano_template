import { useEffect, useState } from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { request } from "@/api/client";
import { GetCurrentUserDetails } from "@/api/user_sys/auth";

const tokenStorageKey = "token";

type AuthStatus = "checking" | "authorized" | "unauthorized";

export default function AuthGuard() {
  const location = useLocation();
  const [authStatus, setAuthStatus] = useState<AuthStatus>("checking");

  useEffect(() => {
    const token = sessionStorage.getItem(tokenStorageKey);
    let isMounted = true;

    if (token) {
      setAuthStatus("authorized");
      return () => {
        isMounted = false;
      };
    }

    void request(() => GetCurrentUserDetails())
      .then((response) => {
        if (!isMounted) {
          return;
        }
        setAuthStatus(response.code === 0 ? "authorized" : "unauthorized");
      })
      .catch(() => {
        if (!isMounted) {
          return;
        }
        setAuthStatus("unauthorized");
      });

    return () => {
      isMounted = false;
    };
  }, []);

  if (authStatus === "authorized") {
    return <Outlet />;
  }

  if (authStatus === "unauthorized") {
    return <Navigate replace state={{ from: location }} to="/login" />;
  }

  return null;
}
