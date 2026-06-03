import { useEffect, useState } from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";

import { request } from "@/api/client";
import { AuthUser } from "@/api/user/auth";
import { useAppDispatch } from "@/app/hooks";
import { setUser } from "@/app/store/userSlice";

type AuthStatus = "checking" | "authorized" | "unauthorized";

export default function AuthGuard() {
  const dispatch = useAppDispatch();
  const location = useLocation();
  const [authStatus, setAuthStatus] = useState<AuthStatus>("checking");

  useEffect(() => {
    let isMounted = true;

    void request(() => AuthUser())
      .then((response) => {
        if (!isMounted) {
          return;
        }
        if (response.code === 0 && response.data) {
          dispatch(setUser({
            user_id: response.data.user_id,
            username: response.data.username,
            role_id: response.data.role_id,
            role_level: response.data.level,
          }));
          setAuthStatus("authorized");
          return;
        }

        setAuthStatus("unauthorized");
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
  }, [dispatch]);

  if (authStatus === "authorized") {
    return <Outlet />;
  }

  if (authStatus === "unauthorized") {
    return <Navigate replace state={{ from: location }} to="/login" />;
  }

  return null;
}
