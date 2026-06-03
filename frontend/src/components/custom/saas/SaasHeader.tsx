import { useEffect, useRef, useState } from "react";
import { LuArrowLeft, LuChevronRight, LuHouse, LuSearch, LuX } from "react-icons/lu";

import type { ThemeTone } from "@/app/store/themeSlice";
import { cx } from "@/app/themeStyles";
import SaasIconButton from "@/components/custom/saas/SaasIconButton";
import type { SaasHeaderAction, SaasNavItem } from "@/components/custom/saas/types";

type SaasHeaderProps = {
  actions: SaasHeaderAction[];
  activeItem: SaasNavItem;
  canNavigateBack: boolean;
  onNavigateBack: () => void;
  onNavigateHome: () => void;
  onSearchQueryChange: (query: string) => void;
  searchQuery: string;
  tone: ThemeTone;
};

export default function SaasHeader({
  actions,
  activeItem,
  canNavigateBack,
  onNavigateBack,
  onNavigateHome,
  onSearchQueryChange,
  searchQuery,
  tone,
}: SaasHeaderProps) {
  const [searchOpen, setSearchOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const searchRef = useRef<HTMLDivElement | null>(null);
  const hasSearchQuery = searchQuery.trim().length > 0;

  useEffect(() => {
    if (searchOpen) {
      inputRef.current?.focus();
    }
  }, [searchOpen]);

  useEffect(() => {
    if (!searchOpen) {
      return;
    }

    function handlePointerDown(event: PointerEvent) {
      const target = event.target;

      if (!(target instanceof Node) || searchRef.current?.contains(target)) {
        return;
      }

      setSearchOpen(false);
      inputRef.current?.blur();
    }

    document.addEventListener("pointerdown", handlePointerDown);

    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
    };
  }, [searchOpen]);

  const navigationActions: SaasHeaderAction[] = [
    {
      disabled: !canNavigateBack,
      icon: LuArrowLeft,
      id: "back",
      label: "Back",
      onClick: onNavigateBack,
    },
    {
      icon: LuHouse,
      id: "home",
      label: "Home",
      onClick: onNavigateHome,
    },
  ];

  return (
    <header className={cx("col-start-2 row-start-1 flex min-w-0 items-center gap-3 border-b px-5", tone.border, tone.surface)}>
      <div className="flex shrink-0 items-center gap-2">
        {navigationActions.map((action) => (
          <SaasIconButton action={action} key={action.id} tone={tone} />
        ))}
      </div>

      <nav
        aria-label="Breadcrumb"
        className={cx(
          "min-w-0 items-center gap-2 text-sm transition-opacity duration-200 ease-out motion-reduce:transition-none",
          searchOpen ? "hidden lg:flex" : "flex",
        )}
      >
        <span className={cx("truncate font-medium", tone.mutedForeground)}>Workspace</span>
        <LuChevronRight className={cx("shrink-0", tone.iconMuted)} size={16} />
        <span className="min-w-0 truncate font-semibold">{activeItem.label}</span>
      </nav>

      <div className="ml-auto flex shrink-0 items-center justify-end gap-2">
        <div
          className={cx(
            "relative flex h-9 shrink-0 items-center overflow-hidden rounded-md border transition-[width,background-color,border-color] duration-300 ease-out motion-reduce:transition-none",
            searchOpen ? "w-56 sm:w-72" : "w-9",
            tone.border,
            tone.secondary,
          )}
          ref={searchRef}
        >
          <button
            aria-label={`Search ${activeItem.label}`}
            className={cx(
              "absolute left-0 top-0 z-10 flex h-9 w-9 items-center justify-center text-sm transition-colors",
              tone.focusRing,
              tone.mutedForeground,
              tone.secondaryHover,
            )}
            onClick={() => setSearchOpen(true)}
            title={`Search ${activeItem.label}`}
            type="button"
          >
            <LuSearch size={18} />
          </button>
          <input
            aria-label={`Search ${activeItem.label}`}
            className={cx(
              "h-full min-w-0 flex-1 bg-transparent pl-9 pr-9 text-sm outline-none transition-opacity duration-200 ease-out placeholder:text-zinc-400 motion-reduce:transition-none",
              searchOpen ? "opacity-100" : "pointer-events-none opacity-0",
              tone.surfaceForeground,
            )}
            onChange={(event) => onSearchQueryChange(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Escape") {
                setSearchOpen(false);
                inputRef.current?.blur();
              }
            }}
            placeholder={`Search ${activeItem.label}`}
            ref={inputRef}
            tabIndex={searchOpen ? 0 : -1}
            type="search"
            value={searchQuery}
          />
          {searchOpen && hasSearchQuery ? (
            <button
              aria-label="Clear search"
              className={cx(
                "absolute right-0 top-0 flex h-9 w-9 items-center justify-center transition-colors",
                tone.focusRing,
                tone.mutedForeground,
                tone.secondaryHover,
              )}
              onClick={() => {
                onSearchQueryChange("");
                inputRef.current?.focus();
              }}
              type="button"
            >
              <LuX size={16} />
            </button>
          ) : null}
        </div>
        {actions.map((action) => (
          <SaasIconButton action={action} key={action.id} tone={tone} />
        ))}
      </div>
    </header>
  );
}
