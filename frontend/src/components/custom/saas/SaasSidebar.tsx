import { LuPanelLeftClose, LuPanelLeftOpen } from "react-icons/lu";

import type { ThemeTone } from "@/app/store/themeSlice";
import { cx } from "@/app/themeStyles";
import type { SaasNavItem } from "@/components/custom/saas/types";

type SaasSidebarProps = {
  activeItemId: string;
  collapsed: boolean;
  items: SaasNavItem[];
  onSelectItem: (itemId: string) => void;
  onToggleCollapsed: () => void;
  tone: ThemeTone;
  utilityItems: SaasNavItem[];
};

function NavButton({
  activeItemId,
  collapsed,
  item,
  onSelectItem,
  tone,
}: {
  activeItemId: string;
  collapsed: boolean;
  item: SaasNavItem;
  onSelectItem: (itemId: string) => void;
  tone: ThemeTone;
}) {
  const Icon = item.icon;
  const isActive = item.id === activeItemId;

  return (
    <button
      aria-label={item.label}
      className={cx(
        "flex h-10 items-center overflow-hidden rounded-md text-left text-sm font-medium transition-[background-color,color,padding,gap,width] duration-200 ease-out motion-reduce:transition-none",
        tone.focusRing,
        collapsed ? "mx-auto w-10 justify-center gap-0 px-0" : "w-full gap-3 px-3",
        isActive
          ? cx(tone.primary, tone.primaryForeground, tone.primaryHover)
          : cx(tone.mutedForeground, tone.secondaryHover),
      )}
      onClick={() => onSelectItem(item.id)}
      title={collapsed ? item.label : undefined}
      type="button"
    >
      <Icon className="shrink-0" size={18} />
      <span
        aria-hidden={collapsed}
        className={cx(
          "min-w-0 truncate transition-[max-width,opacity,transform] duration-200 ease-out motion-reduce:transition-none",
          collapsed ? "max-w-0 translate-x-1 opacity-0" : "max-w-36 translate-x-0 opacity-100",
        )}
      >
        {item.label}
      </span>
    </button>
  );
}

export default function SaasSidebar({
  activeItemId,
  collapsed,
  items,
  onSelectItem,
  onToggleCollapsed,
  tone,
  utilityItems,
}: SaasSidebarProps) {
  const CollapseIcon = collapsed ? LuPanelLeftOpen : LuPanelLeftClose;

  return (
    <>
      <div
        className={cx(
          "col-start-1 row-start-1 flex min-w-0 items-center gap-3 border-r border-b px-4",
          tone.border,
          tone.surface,
        )}
      >
        <div
          className={cx(
            "flex h-10 w-10 shrink-0 items-center justify-center rounded-md border text-sm font-semibold",
            tone.border,
            tone.surfaceForeground,
          )}
        >
          NT
        </div>
        <div
          aria-hidden={collapsed}
          className={cx(
            "min-w-0 overflow-hidden transition-[max-width,opacity,transform] duration-200 ease-out motion-reduce:transition-none",
            collapsed ? "max-w-0 translate-x-1 opacity-0" : "max-w-44 translate-x-0 opacity-100",
          )}
        >
          <p className="truncate text-sm font-semibold leading-5">Nano Template</p>
          <p className={cx("truncate text-xs leading-4", tone.subtleForeground)}>SaaS Workspace</p>
        </div>
      </div>

      <aside className={cx("col-start-1 row-start-2 flex min-h-0 flex-col border-r", tone.border, tone.surface)}>
        <nav className="flex flex-1 flex-col gap-1 overflow-y-auto p-3" aria-label="Primary navigation">
          {items.map((item) => (
            <NavButton
              activeItemId={activeItemId}
              collapsed={collapsed}
              item={item}
              key={item.id}
              onSelectItem={onSelectItem}
              tone={tone}
            />
          ))}
        </nav>

        <div className={cx("flex flex-col gap-1 border-t p-3", tone.border)}>
          {utilityItems.map((item) => (
            <NavButton
              activeItemId={activeItemId}
              collapsed={collapsed}
              item={item}
              key={item.id}
              onSelectItem={onSelectItem}
              tone={tone}
            />
          ))}

          <button
            aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
            className={cx(
              "flex h-10 items-center overflow-hidden rounded-md text-sm font-medium transition-[background-color,color,padding,gap,width] duration-200 ease-out motion-reduce:transition-none",
              tone.focusRing,
              tone.mutedForeground,
              tone.secondaryHover,
              collapsed ? "mx-auto w-10 justify-center gap-0 px-0" : "w-full gap-3 px-3",
            )}
            onClick={onToggleCollapsed}
            title={collapsed ? "Expand sidebar" : undefined}
            type="button"
          >
            <CollapseIcon className="shrink-0" size={18} />
            <span
              aria-hidden={collapsed}
              className={cx(
                "min-w-0 truncate transition-[max-width,opacity,transform] duration-200 ease-out motion-reduce:transition-none",
                collapsed ? "max-w-0 translate-x-1 opacity-0" : "max-w-36 translate-x-0 opacity-100",
              )}
            >
              Collapse
            </span>
          </button>
        </div>
      </aside>
    </>
  );
}
