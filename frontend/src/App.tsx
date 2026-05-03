import { Navigate, Route, Routes } from "react-router-dom";

import AuthGuard from "@/components/custom/AuthGuard";
import HomePage from "@/pages/HomePage";
import LoginPage from "@/pages/LoginPage";

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route element={<AuthGuard />}>
        <Route path="/" element={<HomePage />} />
      </Route>
      <Route path="*" element={<Navigate replace to="/" />} />
    </Routes>
  );
}
