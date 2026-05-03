import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { LuMoon, LuSun } from "react-icons/lu";

export default function HomePage() {
  const [isDark, setIsDark] = useState(true);
  const navigate = useNavigate();

  function handleLogout() {
    sessionStorage.removeItem("token");
    navigate("/login", { replace: true });
  }

  return (
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
          className={`rounded-full px-4 py-3 text-sm font-semibold transition-all duration-300 hover:scale-105 ${isDark ? "bg-[#cc241d] text-[#fbf1c7]" : "bg-[#9d0006] text-[#fbf1c7]"}`}
          onClick={handleLogout}
          type="button"
        >
          Logout
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
  );
}

