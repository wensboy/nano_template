import { FormEvent, useEffect, useRef, useState } from "react";
import { LuChevronDown, LuCircleHelp, LuLock, LuSquareUserRound, LuUserRound } from "react-icons/lu";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import { request } from "@/api/client";
import { UserLogin, UserRegister } from "@/api/user_sys/auth";
import { InputFieldProps, LoginFormState, RegisterFormState } from "@/models/user_sys/user";

import loginBackground from "@/assets/login_bg.png";

const tokenStorageKey = "token";
const roleOptions = [
  { label: "默认", value: "default" },
  { label: "管理", value: "admin" },
  { label: "编辑", value: "editor" },
  { label: "访客", value: "visitor" },
];

type AuthMode = "login" | "register";

function isLocationStateWithFrom(
  state: unknown,
): state is { from: { pathname: string } } {
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

function InputField({ icon, name, onChange, placeholder, type = "text", value }: InputFieldProps) {
  return (
    <label className="group relative block w-full">
      <span className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-[rgba(255,255,255,0.7)]">
        {icon}
      </span>
      <input
        className="h-11 w-full rounded-full border border-[rgba(255,255,255,0.2)] bg-[rgba(255,255,255,0.15)] px-4 pr-4 pl-10 text-sm font-normal text-white outline-none backdrop-blur-[10px] transition duration-200 placeholder:text-[rgba(255,255,255,0.6)] hover:bg-[rgba(255,255,255,0.2)] focus:border-white focus:bg-[rgba(255,255,255,0.25)]"
        name={name}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
        type={type}
        value={value}
      />
    </label>
  );
}

export default function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const token = sessionStorage.getItem(tokenStorageKey);
  const [mode, setMode] = useState<AuthMode>("login");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [feedbackMessage, setFeedbackMessage] = useState("");
  const [selectedRole, setSelectedRole] = useState("");
  const [isRoleMenuOpen, setIsRoleMenuOpen] = useState(false);
  const [loginForm, setLoginForm] = useState<LoginFormState>({
    username: "",
    password: "",
    rememberPassword: false,
  });
  const [registerForm, setRegisterForm] = useState<RegisterFormState>({
    username: "",
    password: "",
    confirmPassword: "",
  });
  const roleMenuRef = useRef<HTMLDivElement | null>(null);
  const redirectPath = isLocationStateWithFrom(location.state) ? location.state.from.pathname : "/";
  const activeRole = roleOptions.find((role) => role.value === (selectedRole || "default")) ?? roleOptions[0];

  useEffect(() => {
    function handlePointerDown(event: MouseEvent) {
      if (roleMenuRef.current && !roleMenuRef.current.contains(event.target as Node)) {
        setIsRoleMenuOpen(false);
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setIsRoleMenuOpen(false);
      }
    }

    document.addEventListener("mousedown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("mousedown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, []);

  function handleRoleChange(role: string) {
    const normalizedRole = role || "default";

    setSelectedRole(role);
    setIsRoleMenuOpen(false);
    setLoginForm((current) => ({
      ...current,
      role: normalizedRole,
    }));
    setRegisterForm((current) => ({
      ...current,
      role: normalizedRole,
    }));
  }

  async function handleLoginSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      setIsSubmitting(true);
      setFeedbackMessage("");

      const response = await request(() => UserLogin(loginForm.username, loginForm.password));

      if (response.code !== 0 || !response.data?.token) {
        setFeedbackMessage(response.message || "登录失败，请稍后重试");
        return;
      }

      if (response.data.token){
        sessionStorage.setItem(tokenStorageKey, response.data.token);
      }
      

      if (!loginForm.rememberPassword) {
        setLoginForm((current) => ({
          ...current,
          password: "",
        }));
      }

      navigate(redirectPath, { replace: true });
    } catch (error) {
      setFeedbackMessage(error instanceof Error ? error.message : "登录失败，请稍后重试");
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleRegisterSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      setIsSubmitting(true);
      setFeedbackMessage("");

      const response = await request(() => UserRegister(registerForm));

      if (response.code !== 0 || !response.data?.user_id) {
        setFeedbackMessage(response.message || "注册失败，请稍后重试");
        return;
      }

      setLoginForm((current) => ({
        ...current,
        username: registerForm.username,
        password: registerForm.password,
      }));
      setRegisterForm({
        username: "",
        password: "",
        confirmPassword: "",
      });
      setMode("login");
      setFeedbackMessage("注册成功，请使用新账号登录");
    } catch (error) {
      setFeedbackMessage(error instanceof Error ? error.message : "注册失败，请稍后重试");
    } finally {
      setIsSubmitting(false);
    }
  }

  if (token) {
    return <Navigate replace to={redirectPath} />;
  }

  return (
    <main
      className="relative flex min-h-screen items-center justify-center overflow-hidden px-6 py-10"
      style={{
        fontFamily:
          '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
        backgroundImage: `linear-gradient(rgba(30,58,138,0.6), rgba(30,58,138,0.6)), url(${loginBackground})`,
        backgroundPosition: "center",
        backgroundSize: "cover",
      }}
    >
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_top,rgba(148,163,184,0.2),transparent_42%)]" />
      <div className="absolute right-6 bottom-6 z-10 w-[220px] max-w-[calc(100vw-48px)]">
        <label className="block rounded-[24px] border border-[rgba(255,255,255,0.22)] bg-[rgba(255,255,255,0.1)] px-4 py-3 text-white shadow-[0_16px_40px_rgba(15,23,42,0.18)] backdrop-blur-xl">
          <span className="mb-2 inline-flex items-center gap-2 text-xs tracking-[0.18em] text-[rgba(255,255,255,0.72)] uppercase">
            <LuUserRound size={14} />
            角色
          </span>
          <div className="relative" ref={roleMenuRef}>
            <button
              aria-expanded={isRoleMenuOpen}
              aria-haspopup="listbox"
              className="flex h-11 w-full items-center justify-between rounded-[16px] border border-[rgba(255,255,255,0.2)] bg-[rgba(255,255,255,0.15)] px-4 text-sm text-white outline-none backdrop-blur-[10px] transition duration-200 hover:bg-[rgba(255,255,255,0.2)] focus:border-white focus:bg-[rgba(255,255,255,0.25)]"
              onClick={() => setIsRoleMenuOpen((current) => !current)}
              type="button"
            >
              <span>{activeRole?.label}</span>
              <span className="pointer-events-none text-[rgba(255,255,255,0.7)]">
                <LuChevronDown
                  className={`transition-transform duration-200 ${isRoleMenuOpen ? "rotate-180" : "rotate-0"}`}
                  size={16}
                />
              </span>
            </button>
            {isRoleMenuOpen ? (
              <div
                className="absolute right-0 bottom-[calc(100%+10px)] left-0 overflow-hidden rounded-[24px] border border-[rgba(255,255,255,0.22)] bg-[rgba(255,255,255,0.12)] p-2 shadow-[0_18px_45px_rgba(15,23,42,0.22)] backdrop-blur-xl"
                role="listbox"
              >
                <div className="space-y-1">
                  {roleOptions.map((role) => {
                    const isActive = role.value === activeRole?.value;

                    return (
                      <button
                        className={`flex w-full items-center rounded-2xl border px-4 py-3 text-left text-sm text-white backdrop-blur-[14px] transition duration-200 ${
                          isActive
                            ? "border-[rgba(255,255,255,0.35)] bg-[rgba(255,255,255,0.22)] shadow-[inset_0_1px_0_rgba(255,255,255,0.18)]"
                            : "border-[rgba(255,255,255,0.14)] bg-[rgba(255,255,255,0.1)] hover:border-[rgba(255,255,255,0.24)] hover:bg-[rgba(255,255,255,0.18)]"
                        }`}
                        key={role.value || "default"}
                        onClick={() => handleRoleChange(role.value)}
                        role="option"
                        type="button"
                      >
                        {role.label}
                      </button>
                    );
                  })}
                </div>
              </div>
            ) : null}
          </div>
        </label>
      </div>
      <section className="relative w-full max-w-[416px] text-white">
        <div className="rounded-[32px] border border-[rgba(255,255,255,0.22)] bg-[rgba(255,255,255,0.08)] px-6 py-8 shadow-[0_24px_64px_rgba(15,23,42,0.28)] backdrop-blur-2xl sm:px-8">
          <div className="mx-auto w-full max-w-[320px]">
          <div className="mb-8">
            <p className="text-2xl font-normal tracking-[0.12em] text-white">LOGO</p>
            <p className="mt-3 text-xs text-[rgba(255,255,255,0.68)]">
              {mode === "login" ? "欢迎回来，请登录继续访问内部页面" : "创建账号后即可进入内部页面"}
            </p>
          </div>

          {feedbackMessage ? (
            <p className="mb-4 rounded-2xl border border-[rgba(255,255,255,0.18)] bg-[rgba(15,23,42,0.18)] px-4 py-3 text-sm text-[rgba(255,255,255,0.88)]">
              {feedbackMessage}
            </p>
          ) : null}

          {mode === "login" ? (
            <form onSubmit={handleLoginSubmit}>
              <div className="space-y-3">
                <InputField
                  icon={<LuUserRound size={16} />}
                  name="username"
                  onChange={(value) => setLoginForm((current) => ({ ...current, username: value }))}
                  placeholder="用户名"
                  value={loginForm.username}
                />
                <InputField
                  icon={<LuLock size={16} />}
                  name="password"
                  onChange={(value) => setLoginForm((current) => ({ ...current, password: value }))}
                  placeholder="密码"
                  type="password"
                  value={loginForm.password}
                />
              </div>

              <button
                className="mt-4 h-11 w-full rounded-full bg-[#2563EB] text-sm font-medium text-white transition duration-200 hover:bg-[#1D4ED8] active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-70"
                disabled={isSubmitting}
                type="submit"
              >
                {isSubmitting ? "登录中..." : "登录"}
              </button>

              <div className="mt-3 flex w-full items-center justify-between text-[12px] text-[rgba(255,255,255,0.6)]">
                <label className="flex cursor-pointer items-center gap-2">
                  <span className="relative flex h-4 w-4 items-center justify-center">
                    <input
                      checked={loginForm.rememberPassword}
                      className="peer absolute inset-0 m-0 h-full w-full cursor-pointer appearance-none rounded-full border border-[rgba(255,255,255,0.7)] bg-transparent"
                      onChange={(event) =>
                        setLoginForm((current) => ({
                          ...current,
                          rememberPassword: event.target.checked,
                        }))
                      }
                      type="checkbox"
                    />
                    <span className="h-2 w-2 rounded-full bg-[#2563EB] opacity-0 transition peer-checked:opacity-100" />
                  </span>
                  <span>记住密码</span>
                </label>
                <button
                  className="transition hover:text-white hover:underline"
                  onClick={() => undefined}
                  type="button"
                >
                  忘记密码
                </button>
              </div>

              <div className="mt-5 flex w-full items-center justify-between text-[12px] text-[rgba(255,255,255,0.6)]">
                <button
                  className="transition hover:text-white hover:underline"
                  onClick={() => {
                    setFeedbackMessage("");
                    setMode("register");
                  }}
                  type="button"
                >
                  创建账号
                </button>
                <button
                  className="inline-flex items-center gap-2 transition hover:text-white hover:underline"
                  onClick={() => undefined}
                  type="button"
                >
                  <LuCircleHelp size={14} />
                  其他帮助
                </button>
              </div>
            </form>
          ) : (
            <form onSubmit={handleRegisterSubmit}>
              <div className="space-y-3">
                <InputField
                  icon={<LuSquareUserRound size={16} />}
                  name="register-username"
                  onChange={(value) => setRegisterForm((current) => ({ ...current, username: value }))}
                  placeholder="用户名"
                  value={registerForm.username}
                />
                <InputField
                  icon={<LuLock size={16} />}
                  name="register-password"
                  onChange={(value) => setRegisterForm((current) => ({ ...current, password: value }))}
                  placeholder="密码"
                  type="password"
                  value={registerForm.password}
                />
                <InputField
                  icon={<LuLock size={16} />}
                  name="register-confirm-password"
                  onChange={(value) =>
                    setRegisterForm((current) => ({ ...current, confirmPassword: value }))
                  }
                  placeholder="确认密码"
                  type="password"
                  value={registerForm.confirmPassword}
                />
              </div>

              <button
                className="mt-4 h-11 w-full rounded-full bg-[#2563EB] text-sm font-medium text-white transition duration-200 hover:bg-[#1D4ED8] active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-70"
                disabled={isSubmitting}
                type="submit"
              >
                {isSubmitting ? "注册中..." : "注册"}
              </button>

              <div className="mt-3 flex w-full items-center justify-between text-[12px] text-[rgba(255,255,255,0.6)]">
                <span>已有账号？</span>
                <button
                  className="transition hover:text-white hover:underline"
                  onClick={() => {
                    setFeedbackMessage("");
                    setMode("login");
                  }}
                  type="button"
                >
                  返回登录
                </button>
              </div>

              <div className="mt-5 flex w-full items-center justify-between text-[12px] text-[rgba(255,255,255,0.6)]">
                <span>注册后可继续完善个人资料</span>
                <button
                  className="inline-flex items-center gap-2 transition hover:text-white hover:underline"
                  onClick={() => undefined}
                  type="button"
                >
                  <LuCircleHelp size={14} />
                  其他帮助
                </button>
              </div>
            </form>
          )}
          </div>
        </div>
      </section>
    </main>
  );
}
