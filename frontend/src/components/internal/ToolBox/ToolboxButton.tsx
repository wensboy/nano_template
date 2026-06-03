import type { CSSProperties, PointerEventHandler } from "react";
import { LuWrench } from "react-icons/lu";

import type { ThemeState } from "@/app/store/themeSlice";
import { cx, getThemeButtonClassName } from "@/app/themeStyles";

type ToolboxButtonProps = {
  dragging: boolean;
  onClick: () => void;
  onPointerCancel: PointerEventHandler<HTMLButtonElement>;
  onPointerDown: PointerEventHandler<HTMLButtonElement>;
  onPointerMove: PointerEventHandler<HTMLButtonElement>;
  onPointerUp: PointerEventHandler<HTMLButtonElement>;
  style: CSSProperties;
  theme: ThemeState;
};

export default function ToolboxButton({
  dragging,
  onClick,
  onPointerCancel,
  onPointerDown,
  onPointerMove,
  onPointerUp,
  style,
  theme,
}: ToolboxButtonProps) {
  return (
    <button
      aria-label="Open toolbox"
      className={cx(
        "fixed z-40 flex h-12 w-12 touch-none items-center justify-center",
        dragging ? "cursor-grabbing" : "cursor-grab",
        getThemeButtonClassName(theme),
      )}
      onClick={onClick}
      onPointerCancel={onPointerCancel}
      onPointerDown={onPointerDown}
      onPointerMove={onPointerMove}
      onPointerUp={onPointerUp}
      style={style}
      title="Open toolbox"
      type="button"
    >
      <LuWrench size={22} />
    </button>
  );
}
