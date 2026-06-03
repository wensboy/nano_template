import { FormEvent, useEffect, useState } from "react";
import { LuChevronDown, LuLock, LuShieldCheck, LuUserRound } from "react-icons/lu";

import { useAppDispatch } from "@/app/hooks";
import { setUser, type UserState } from "@/app/store/userSlice";
import { request } from "@/api/client";
import { ListRoles, type Role } from "@/api/role/role";
import { UserLogin, UserRegister } from "@/api/user/auth";
import AlertPop, { type AlertPopVariant } from "@/components/internal/AlertPop";

type AuthMode = "login" | "register";

type AuthDialogProps = {
  defaultMode?: AuthMode;
  onAuthenticated?: () => void;
  open: boolean;
};

type JwtPayload = {
  user_id?: number;
  username?: string;
  user_role?: number;
  role_level?: number;
};

function decodeUserState(token: string): UserState | null {
  const [, payload] = token.split(".");
  if (!payload) {
    return null;
  }

  try {
    const normalizedPayload = payload.replace(/-/g, "+").replace(/_/g, "/");
    const paddedPayload = normalizedPayload.padEnd(
      normalizedPayload.length + ((4 - (normalizedPayload.length % 4)) % 4),
      "=",
    );
    const decoded = JSON.parse(window.atob(paddedPayload)) as JwtPayload;
    return {
      user_id: Number(decoded.user_id ?? 0),
      username: decoded.username ?? "",
      role_id: Number(decoded.user_role ?? 0),
      role_level: Number(decoded.role_level ?? 0),
    };
  } catch {
    return null;
  }
}

function numberFromInput(value: string) {
  const normalized = value.trim();
  if (normalized === "") {
    return 0;
  }
  const parsed = Number(normalized);
  if (!Number.isFinite(parsed)) {
    throw new Error("角色 ID 必须是数字");
  }
  return parsed;
}

export default function AuthDialog({ defaultMode = "login", onAuthenticated, open }: AuthDialogProps) {
  const dispatch = useAppDispatch();
  const [mode, setMode] = useState<AuthMode>(defaultMode);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [roleId, setRoleId] = useState("0");
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoadingRoles, setIsLoadingRoles] = useState(false);
  const [feedbackMessage, setFeedbackMessage] = useState("");
  const [feedbackVariant, setFeedbackVariant] = useState<AlertPopVariant>("destructive");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!open) {
      return;
    }

    let isMounted = true;

    async function loadRoles() {
      try {
        setIsLoadingRoles(true);
        const response = await request(() => ListRoles({ page_size: 100, state: 1 }));
        if (!isMounted) {
          return;
        }

        const nextRoles = response.code === 0 ? response.data?.items ?? [] : [];
        setRoles(nextRoles);
        setRoleId((currentRoleId) => {
          if (nextRoles.length === 0) {
            return "0";
          }
          return nextRoles.some((role) => String(role.id) === currentRoleId)
            ? currentRoleId
            : String(nextRoles[0].id);
        });
      } catch {
        if (isMounted) {
          setRoles([]);
        }
      } finally {
        if (isMounted) {
          setIsLoadingRoles(false);
        }
      }
    }

    void loadRoles();

    return () => {
      isMounted = false;
    };
  }, [open]);

  if (!open) {
    return null;
  }

  async function handleLogin(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      setIsSubmitting(true);
      setFeedbackMessage("");

      const response = await request(() => UserLogin(username, password, numberFromInput(roleId)));
      if (response.code !== 0 || !response.data?.token) {
        setFeedbackVariant("destructive");
        setFeedbackMessage(response.message || "登录失败，请稍后重试");
        return;
      }

      const userState = decodeUserState(response.data.token);
      if (userState) {
        dispatch(setUser(userState));
      }
      onAuthenticated?.();
    } catch (error) {
      setFeedbackVariant("destructive");
      setFeedbackMessage(error instanceof Error ? error.message : "登录失败，请稍后重试");
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleRegister(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      setIsSubmitting(true);
      setFeedbackMessage("");

      const response = await request(() =>
        UserRegister(
          {
            username,
            password,
            confirmPassword,
            role: roleId,
          },
          numberFromInput(roleId),
        ),
      );
      if (response.code !== 0) {
        setFeedbackVariant("destructive");
        setFeedbackMessage(response.message || "注册失败，请稍后重试");
        return;
      }

      setMode("login");
      setConfirmPassword("");
      setFeedbackVariant("success");
      setFeedbackMessage("注册成功，请登录");
    } catch (error) {
      setFeedbackVariant("destructive");
      setFeedbackMessage(error instanceof Error ? error.message : "注册失败，请稍后重试");
    } finally {
      setIsSubmitting(false);
    }
  }

  const isLogin = mode === "login";

  return (
    <div className="fixed inset-0 z-50 flex min-h-screen items-center justify-center px-4 py-6">
      <AlertPop message={feedbackMessage} variant={feedbackVariant} />
      <div aria-hidden="true" className="absolute inset-0 bg-zinc-950/45 backdrop-blur-md" />
      <section
        aria-modal="true"
        className="relative w-full max-w-[400px] rounded-lg border border-zinc-200 bg-white p-6 text-zinc-950 shadow-2xl"
        role="dialog"
      >
        <div className="mb-6 flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-md border border-zinc-200 bg-zinc-50 text-zinc-700">
            <LuShieldCheck size={18} />
          </div>
          <div>
            <h1 className="text-lg font-semibold">{isLogin ? "登录" : "注册"}</h1>
            <p className="text-sm text-zinc-500">{isLogin ? "使用账号进入系统" : "创建一个新账号"}</p>
          </div>
        </div>

        <div className="mb-5 flex items-center gap-3">
          <div className="grid flex-1 grid-cols-2 rounded-lg bg-zinc-100 p-1 text-sm">
            <button
              className={`h-9 rounded-md transition ${isLogin ? "bg-white text-zinc-950 shadow-sm" : "text-zinc-500 hover:text-zinc-900"}`}
              onClick={() => {
                setFeedbackMessage("");
                setMode("login");
              }}
              type="button"
            >
              登录
            </button>
            <button
              className={`h-9 rounded-md transition ${!isLogin ? "bg-white text-zinc-950 shadow-sm" : "text-zinc-500 hover:text-zinc-900"}`}
              onClick={() => {
                setFeedbackMessage("");
                setMode("register");
              }}
              type="button"
            >
              注册
            </button>
          </div>

          <span className="relative block w-32 shrink-0">
            <select
              className="h-10 w-full appearance-none rounded-md border border-zinc-200 bg-white px-3 pr-9 text-sm outline-none transition focus:border-zinc-400 focus:ring-2 focus:ring-zinc-200 disabled:cursor-not-allowed disabled:opacity-60"
              disabled={isLoadingRoles}
              onChange={(event) => setRoleId(event.target.value)}
              value={roleId}
            >
              {roles.length > 0 ? (
                roles.map((role) => (
                  <option key={role.id} value={role.id}>
                    {role.name}
                  </option>
                ))
              ) : (
                <option value="0">{isLoadingRoles ? "加载角色中..." : "默认角色"}</option>
              )}
            </select>
            <LuChevronDown
              aria-hidden="true"
              className="pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 text-zinc-400"
              size={16}
            />
          </span>
        </div>

        <form className="space-y-4" onSubmit={isLogin ? handleLogin : handleRegister}>
          <label className="block space-y-2 text-sm font-medium text-zinc-800">
            <span>用户名</span>
            <span className="relative block">
              <LuUserRound className="absolute top-1/2 left-3 -translate-y-1/2 text-zinc-400" size={16} />
              <input
                className="h-10 w-full rounded-md border border-zinc-200 bg-white px-3 pl-9 text-sm outline-none transition focus:border-zinc-400 focus:ring-2 focus:ring-zinc-200"
                onChange={(event) => setUsername(event.target.value)}
                value={username}
              />
            </span>
          </label>

          <div className="block space-y-2 text-sm font-medium text-zinc-800">
            <span className="flex items-center justify-between">
              <label htmlFor="auth-dialog-password">密码</label>
              {isLogin ? (
                <button
                  className="inline-flex w-fit flex-none items-center p-0 text-xs font-medium text-zinc-500"
                  onClick={() => undefined}
                  type="button"
                >
                  <span className="underline-offset-4 transition hover:text-zinc-900 hover:underline">忘记密码</span>
                </button>
              ) : null}
            </span>
            <span className="relative block">
              <LuLock className="absolute top-1/2 left-3 -translate-y-1/2 text-zinc-400" size={16} />
              <input
                className="h-10 w-full rounded-md border border-zinc-200 bg-white px-3 pl-9 text-sm outline-none transition focus:border-zinc-400 focus:ring-2 focus:ring-zinc-200"
                id="auth-dialog-password"
                onChange={(event) => setPassword(event.target.value)}
                type="password"
                value={password}
              />
            </span>
          </div>

          {!isLogin ? (
            <label className="block space-y-2 text-sm font-medium text-zinc-800">
              <span>确认密码</span>
              <span className="relative block">
                <LuLock className="absolute top-1/2 left-3 -translate-y-1/2 text-zinc-400" size={16} />
                <input
                  className="h-10 w-full rounded-md border border-zinc-200 bg-white px-3 pl-9 text-sm outline-none transition focus:border-zinc-400 focus:ring-2 focus:ring-zinc-200"
                  onChange={(event) => setConfirmPassword(event.target.value)}
                  type="password"
                  value={confirmPassword}
                />
              </span>
            </label>
          ) : null}

          <button
            className="h-10 w-full rounded-md bg-zinc-950 text-sm font-medium text-white transition hover:bg-zinc-800 disabled:cursor-not-allowed disabled:opacity-60"
            disabled={isSubmitting}
            type="submit"
          >
            {isSubmitting ? "处理中..." : isLogin ? "登录" : "注册"}
          </button>
        </form>
      </section>
    </div>
  );
}
