import { useEffect } from "react";
import { createPortal } from "react-dom";
import { LuX } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import { useToolbox } from "@/app/toolbox";
import ToolboxButton from "@/components/internal/ToolboxButton";

export default function ToolBox() {
  const { activeItem, activeItemId, closeToolbox, isOpen, items, openToolbox, setActiveItem } = useToolbox();
  const themeMode = useAppSelector((state) => state.theme.mode);
  const isDark = themeMode === "dark";

  useEffect(() => {
    if (!isOpen) {
      return;
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        closeToolbox();
      }
    }

    const originalOverflow = document.body.style.overflow;

    document.body.style.overflow = "hidden";
    window.addEventListener("keydown", handleKeyDown);

    return () => {
      document.body.style.overflow = originalOverflow;
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [closeToolbox, isOpen]);

  if (items.length === 0) {
    return null;
  }

  const modal = isOpen ? (
    <div
      className={`fixed inset-0 z-50 flex items-center justify-center px-4 py-8 backdrop-blur-sm ${isDark ? "bg-black/55" : "bg-black/35"}`}
      onClick={closeToolbox}
    >
      <div
        className={`relative flex h-[min(68vh,520px)] w-full max-w-3xl overflow-hidden rounded-[14px] border text-black shadow-[0_28px_100px_rgba(0,0,0,0.35)] ${isDark ? "border-white/15 bg-white" : "border-black/10 bg-[#fbf1c7]"}`}
        onClick={(event) => event.stopPropagation()}
      >
        <aside className={`flex w-20 shrink-0 flex-col items-center gap-3 border-r px-3 py-6 ${isDark ? "border-black/10 bg-[#f4f4f4]" : "border-black/10 bg-[#ebdbb2]"}`}>
          {items.map((item) => {
            const isActive = item.id === activeItemId;

            return (
              <button
                aria-label={item.label}
                className={`flex h-12 w-12 items-center justify-center rounded-2xl border text-lg transition ${isActive ? "border-black bg-black text-white shadow-[0_10px_24px_rgba(0,0,0,0.16)]" : "border-black/10 bg-white text-black/70 hover:border-black/25 hover:text-black"}`}
                key={item.id}
                onClick={() => setActiveItem(item.id)}
                title={item.label}
                type="button"
              >
                {item.icon}
              </button>
            );
          })}
        </aside>

        <section className={`flex min-w-0 flex-1 flex-col ${isDark ? "bg-white" : "bg-[#fbf1c7]"}`}>
          <header className="flex items-center justify-between border-b border-black/10 px-6 py-4">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.28em] text-black/45">Toolbox</p>
              <h2 className="mt-1 text-lg font-semibold text-black">{activeItem?.label ?? "工具箱"}</h2>
            </div>
            <button
              aria-label="Close toolbox"
              className="flex h-10 w-10 items-center justify-center rounded-full border border-black/10 text-black/70 transition hover:border-black/25 hover:text-black"
              onClick={closeToolbox}
              type="button"
            >
              <LuX size={18} />
            </button>
          </header>

          <div className="min-h-0 flex-1 overflow-y-auto px-6 py-5">
            {activeItem?.content ?? (
              <div className="flex h-full items-center justify-center text-sm text-black/50">
                当前没有可显示的工具内容。
              </div>
            )}
          </div>
        </section>
      </div>
    </div>
  ) : null;

  return (
    <>
      <ToolboxButton onClick={openToolbox} themeMode={themeMode} />
      {typeof document !== "undefined" ? createPortal(modal, document.body) : null}
    </>
  );
}
