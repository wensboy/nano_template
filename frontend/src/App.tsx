import { useEffect } from "react";
import { Navigate, Route, Routes, useLocation, useNavigate } from "react-router-dom";

import { UNAUTHORIZED_EVENT } from "@/api/interceptor";
import AuthGuard from "@/components/internal/AuthGuard";
import HomePage from "@/pages/HomePage";
import LoginPage from "@/pages/LoginPage";

function UnauthorizedRedirect() {
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    function handleUnauthorized() {
      if (location.pathname !== "/login") {
        navigate("/login", { replace: true });
      }
    }

    window.addEventListener(UNAUTHORIZED_EVENT, handleUnauthorized);
    return () => window.removeEventListener(UNAUTHORIZED_EVENT, handleUnauthorized);
  }, [location.pathname, navigate]);

  return null;
}

export default function App() {
  return (
    <>
      <UnauthorizedRedirect />
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route element={<AuthGuard />}>
          <Route path="/" element={<HomePage />} />
        </Route>
        {/* <Route path="/" element={<HomePage />} /> */}
        <Route path="*" element={<Navigate replace to="/" />} />
      </Routes>
    </>
  );
}
