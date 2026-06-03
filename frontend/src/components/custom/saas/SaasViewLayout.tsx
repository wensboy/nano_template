import type { ReactNode } from "react";

import type { ThemeTone } from "@/app/store/themeSlice";
import { cx } from "@/app/themeStyles";

export type SaasViewLayoutNavItem = {
  description: string;
  id: string;
  label: string;
  meta: string;
};

type SaasViewLayoutProps = {
  activeItemId: string;
  description: string;
  emptyMessage: string;
  items: SaasViewLayoutNavItem[];
  label: string;
  onActiveItemChange: (itemId: string) => void;
  renderContent: (item: SaasViewLayoutNavItem) => ReactNode;
  tone: ThemeTone;
};

export default function SaasViewLayout({
  activeItemId,
  description,
  emptyMessage,
  items,
  label,
  onActiveItemChange,
  renderContent,
  tone,
}: SaasViewLayoutProps) {
  const activeItem = items.find((item) => item.id === activeItemId) ?? items[0] ?? null;

  return (
    <section className={cx("flex h-full min-h-0 flex-col overflow-hidden rounded-md border", tone.border, tone.subtle)}>
      <header className={cx("shrink-0 border-b px-5 py-4", tone.border)}>
        <h1 className="truncate text-2xl font-semibold leading-8">{label}</h1>
        <p className={cx("mt-1 max-w-3xl text-sm leading-6", tone.mutedForeground)}>{description}</p>
      </header>

      <div className={cx("grid min-h-0 flex-1 grid-cols-[14rem_minmax(0,1fr)]", tone.surface)}>
        <nav className={cx("min-h-0 overflow-y-auto border-r", tone.border)} aria-label={`${label} view sections`}>
          <div className="flex min-h-full flex-col">
            {items.map((item) => {
              const isActive = item.id === activeItem?.id;

              return (
                <button
                  aria-current={isActive ? "page" : undefined}
                  className={cx(
                    "min-h-9 min-w-0 border-b px-3 py-2 text-left text-sm transition-colors",
                    tone.border,
                    tone.focusRing,
                    isActive
                      ? cx("sticky top-0 bottom-0 z-10 border-l-2 border-l-cyan-500 font-semibold", tone.surface)
                      : cx("font-medium", tone.mutedForeground, tone.secondaryHover),
                  )}
                  key={item.id}
                  onClick={() => onActiveItemChange(item.id)}
                  type="button"
                >
                  <span className="block truncate">{item.label}</span>
                </button>
              );
            })}
          </div>
        </nav>

        <div className="min-h-0 overflow-y-auto">
          {activeItem ? (
            renderContent(activeItem)
          ) : (
            <div className={cx("flex h-full items-center justify-center px-5 text-sm", tone.mutedForeground)}>
              {emptyMessage}
            </div>
          )}
        </div>
      </div>

      <footer className={cx("h-10 shrink-0 border-t", tone.border)} aria-hidden="true" />
    </section>
  );
}
