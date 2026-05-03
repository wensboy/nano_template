import { Navigate, useLocation, useNavigate } from "react-router-dom";

const tokenStorageKey = "token";

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const token = sessionStorage.getItem(tokenStorageKey);
  const redirectPath =
    typeof location.state === "object" &&
    location.state !== null &&
    "from" in location.state &&
    typeof location.state.from === "object" &&
    location.state.from !== null &&
    "pathname" in location.state.from &&
    typeof location.state.from.pathname === "string"
      ? location.state.from.pathname
      : "/";

  function handleMockLogin() {
    sessionStorage.setItem(tokenStorageKey, "mock-token");
    navigate(redirectPath, { replace: true });
  }

  if (token) {
    return <Navigate replace to={redirectPath} />;
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-[linear-gradient(135deg,#fdf6e3_0%,#fe8019_45%,#1d2021_100%)] px-6 py-10">
      <section className="w-full max-w-md rounded-[28px] bg-[#f9f5d7]/92 p-8 shadow-[0_24px_80px_rgba(29,32,33,0.28)] backdrop-blur">
        <p className="mb-3 text-sm font-semibold uppercase tracking-[0.35em] text-[#9d0006]">
          Auth Gateway
        </p>
        <h1 className="text-4xl font-black text-[#1d2021]">Login</h1>
        <p className="mt-4 text-sm leading-7 text-[#3c3836]">
          内部页面会先检查 <code>sessionStorage</code> 中的 <code>token</code>。如果刷新后仍不存在，
          就会自动跳转到这里。
        </p>
        <button
          className="mt-8 w-full rounded-2xl bg-[#b57614] px-5 py-4 text-base font-bold text-[#fbf1c7] transition-transform duration-200 hover:scale-[1.02]"
          onClick={handleMockLogin}
          type="button"
        >
          Mock Login
        </button>
      </section>
    </main>
  );
}

