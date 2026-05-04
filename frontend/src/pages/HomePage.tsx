import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { LuCompass, LuLayoutGrid, LuLogOut, LuMoon, LuSparkles, LuSun } from "react-icons/lu";

import { ToolboxProvider, type ToolboxItem } from "@/app/toolbox";
import { request } from "@/api/client";
import { UserLogout } from "@/api/user_sys/auth";

export default function HomePage() {
  const [isDark, setIsDark] = useState(true);
  const navigate = useNavigate();

  async function handleLogout() {
    try {
      await request(() => UserLogout());
    } catch {
      // Keep frontend logout resilient even if the backend cookie is already invalid.
    } finally {
      sessionStorage.removeItem("token");
      navigate("/login", { replace: true });
    }
  }

  const toolboxItems: ToolboxItem[] = [
    {
      id: "overview",
      label: "Overview",
      icon: <LuLayoutGrid size={18} />,
      content: (
        <div className="space-y-4">
          <p className="text-sm leading-7 text-black/70">
            这个 toolbox 通过 provider 注入到页面顶层，右下角按钮会始终悬浮显示，适合挂载页面级的工具、
            说明和快捷入口。
          </p>
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="rounded-2xl border border-black/10 bg-[#f7f7f7] p-4">
              <p className="text-xs uppercase tracking-[0.2em] text-black/45">Theme</p>
              <p className="mt-2 text-lg font-semibold text-black">{isDark ? "Dark Atmosphere" : "Warm Light"}</p>
            </div>
            <div className="rounded-2xl border border-black/10 bg-[#f7f7f7] p-4">
              <p className="text-xs uppercase tracking-[0.2em] text-black/45">Auth</p>
              <p className="mt-2 text-lg font-semibold text-black">Protected Route Ready</p>
            </div>
          </div>
        </div>
      ),
    },
    {
      id: "actions",
      label: "Actions",
      icon: <LuSparkles size={18} />,
      content: (
        <div className="space-y-4">
          <p className="text-sm leading-7 text-black/70">这里可以放当前页面常用的快捷操作，避免打断主视图布局。</p>
          <div className="flex flex-wrap gap-3">
            <button
              className="rounded-full bg-black px-4 py-2 text-sm font-medium text-white transition hover:bg-black/85"
              onClick={() => setIsDark((current) => !current)}
              type="button"
            >
              切换主题
            </button>
            <button
              className="rounded-full border border-black/15 px-4 py-2 text-sm font-medium text-black transition hover:border-black/30 hover:bg-black/5"
              onClick={handleLogout}
              type="button"
            >
              退出登录
            </button>
          </div>
        </div>
      ),
    },
    {
      id: "notes",
      label: "Notes",
      icon: <LuCompass size={18} />,
      content: (
        <div className="space-y-3">
          <p className="text-sm font-medium text-black">接入说明</p>
          <ul className="space-y-2 text-sm leading-7 text-black/70">
            <li>1. 用 `ToolboxProvider` 包住页面。</li>
            <li>2. 传入带 `icon / label / content` 的工具项数组。</li>
            <li>3. 页面右下角会自动获得 toolbox 入口。</li>
          </ul>
        </div>
      ),
    },
  ];

  return (
    <ToolboxProvider items={toolboxItems}>
      <div
        className={`relative flex h-screen items-center justify-center overflow-hidden transition-colors duration-500 ${isDark ? "bg-[#282828]" : "bg-[#fbf1c7]"}`}
      >
        <div className="absolute right-6 top-6 z-10 flex items-center gap-3">
          <button
            className={`rounded-full p-3 transition-all duration-300 hover:scale-110 ${isDark ? "bg-[#3c3836] text-[#ebdbb2]" : "bg-[#ebdbb2] text-[#282828]"}`}
            onClick={() => setIsDark(!isDark)}
          >
            {isDark ? <LuSun size={24} /> : <LuMoon size={24} />}
          </button>
          <button
            className={`rounded-full p-3 transition-all duration-300 hover:scale-110 ${isDark ? "bg-[#3c3836] text-[#ebdbb2]" : "bg-[#ebdbb2] text-[#282828]"}`}
            onClick={handleLogout}
            type="button"
          >
            <LuLogOut size={24} />
          </button>
        </div>

        <div className="absolute inset-0 overflow-hidden">
          <div
            className={`absolute left-1/4 top-1/4 h-64 w-64 rounded-full opacity-30 mix-blend-multiply blur-3xl filter animate-float1 ${isDark ? "bg-[#d3869b]" : "bg-[#b8bb26]"}`}
          />
          <div
            className={`absolute right-1/4 top-1/3 h-72 w-72 rounded-full opacity-30 mix-blend-multiply blur-3xl filter animate-float2 ${isDark ? "bg-[#b8bb26]" : "bg-[#d3869b]"}`}
          />
          <div
            className={`absolute bottom-1/4 left-1/3 h-80 w-80 rounded-full opacity-25 mix-blend-multiply blur-3xl filter animate-float3 ${isDark ? "bg-[#83a598]" : "bg-[#fabd2f]"}`}
          />
        </div>

        <h1
          className={`relative bg-clip-text text-6xl font-bold text-transparent animate-pulse ${isDark ? "bg-gradient-to-r from-[#fabd2f] via-[#d3869b] to-[#83a598]" : "bg-gradient-to-r from-[#d3869b] via-[#fabd2f] to-[#b8bb26]"}`}
        >
          Nano Template
        </h1>

        <style>{`
          @keyframes float1 {
            0%, 100% { transform: translate(0, 0) scale(1); }
            50% { transform: translate(30px, -30px) scale(1.1); }
          }
          @keyframes float2 {
            0%, 100% { transform: translate(0, 0) scale(1); }
            50% { transform: translate(-40px, 20px) scale(1.05); }
          }
          @keyframes float3 {
            0%, 100% { transform: translate(0, 0) scale(1); }
            50% { transform: translate(20px, 40px) scale(1.1); }
          }
          .animate-float1 { animation: float1 8s ease-in-out infinite; }
          .animate-float2 { animation: float2 10s ease-in-out infinite; }
          .animate-float3 { animation: float3 12s ease-in-out infinite; }
        `}</style>
      </div>
    </ToolboxProvider>
  );
}
