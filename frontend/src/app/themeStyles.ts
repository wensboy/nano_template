import type { ThemeMode } from "@/app/themeSlice";

export function getThemeButtonClassName(themeMode: ThemeMode) {
  const isDark = themeMode === "dark";

  return `rounded-full transition-all duration-300 hover:scale-110 ${isDark ? "bg-[#3c3836] text-[#ebdbb2]" : "bg-[#ebdbb2] text-[#282828]"}`;
}
