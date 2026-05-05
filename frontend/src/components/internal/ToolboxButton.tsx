import { LuWrench } from "react-icons/lu";

import type { ThemeMode } from "@/app/themeSlice";
import { getThemeButtonClassName } from "@/app/themeStyles";

type ToolboxButtonProps = {
  onClick: () => void;
  themeMode: ThemeMode;
};

export default function ToolboxButton({ onClick, themeMode }: ToolboxButtonProps) {
  return (
    <button
      aria-label="Open toolbox"
      className={`fixed right-6 bottom-6 z-40 flex h-14 w-14 items-center justify-center active:scale-95 ${getThemeButtonClassName(themeMode)}`}
      onClick={onClick}
      type="button"
    >
      <LuWrench size={22} />
    </button>
  );
}
