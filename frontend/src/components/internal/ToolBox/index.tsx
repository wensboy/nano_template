import { useEffect, useRef, useState, type CSSProperties, type PointerEvent as ReactPointerEvent } from "react";
import { createPortal } from "react-dom";
import { LuX } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeTone } from "@/app/themeStyles";
import { useToolbox, type ToolboxItem, type ToolboxPlacement } from "@/app/toolbox";
import ToolboxButton from "@/components/internal/ToolBox/ToolboxButton";

type ToolBoxProps = {
  placement: ToolboxPlacement;
};

type ButtonPosition = {
  left: number;
  top: number;
};

type DragState = {
  moved: boolean;
  originLeft: number;
  originTop: number;
  startX: number;
  startY: number;
};

const buttonSize = 48;
const placementOffset = 24;
const viewportPadding = 8;

function getPlacementStyle(placement: ToolboxPlacement): CSSProperties {
  switch (placement) {
    case "bottom-left":
      return { bottom: placementOffset, left: placementOffset };
    case "top-left":
      return { left: placementOffset, top: placementOffset };
    case "top-right":
      return { right: placementOffset, top: placementOffset };
    case "bottom-right":
    default:
      return { bottom: placementOffset, right: placementOffset };
  }
}

function clampPosition(position: ButtonPosition) {
  const maxLeft = Math.max(viewportPadding, window.innerWidth - buttonSize - viewportPadding);
  const maxTop = Math.max(viewportPadding, window.innerHeight - buttonSize - viewportPadding);

  return {
    left: Math.min(Math.max(position.left, viewportPadding), maxLeft),
    top: Math.min(Math.max(position.top, viewportPadding), maxTop),
  };
}

export default function ToolBox({ placement }: ToolBoxProps) {
  const { activeItem, activeItemId, closeToolbox, isOpen, items, openToolbox, setActiveItem } = useToolbox();
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const profileItem = items.find((item) => item.id === "profile");
  const mainItems = profileItem ? items.filter((item) => item.id !== profileItem.id) : items;
  const dragStateRef = useRef<DragState | null>(null);
  const suppressClickRef = useRef(false);
  const [buttonPosition, setButtonPosition] = useState<ButtonPosition | null>(null);
  const [dragging, setDragging] = useState(false);

  useEffect(() => {
    setButtonPosition(null);
  }, [placement]);

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

  useEffect(() => {
    if (!buttonPosition) {
      return;
    }

    function handleResize() {
      setButtonPosition((position) => (position ? clampPosition(position) : position));
    }

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, [buttonPosition]);

  if (items.length === 0) {
    return null;
  }

  function handleToolboxButtonClick() {
    if (suppressClickRef.current) {
      return;
    }

    openToolbox();
  }

  function handleToolboxButtonPointerDown(event: ReactPointerEvent<HTMLButtonElement>) {
    if (event.button !== 0) {
      return;
    }

    const rect = event.currentTarget.getBoundingClientRect();

    dragStateRef.current = {
      moved: false,
      originLeft: rect.left,
      originTop: rect.top,
      startX: event.clientX,
      startY: event.clientY,
    };
    setButtonPosition({ left: rect.left, top: rect.top });
    event.currentTarget.setPointerCapture(event.pointerId);
  }

  function handleToolboxButtonPointerMove(event: ReactPointerEvent<HTMLButtonElement>) {
    const dragState = dragStateRef.current;

    if (!dragState) {
      return;
    }

    const deltaX = event.clientX - dragState.startX;
    const deltaY = event.clientY - dragState.startY;

    if (!dragState.moved && Math.hypot(deltaX, deltaY) < 4) {
      return;
    }

    dragState.moved = true;
    setDragging(true);
    setButtonPosition(
      clampPosition({
        left: dragState.originLeft + deltaX,
        top: dragState.originTop + deltaY,
      }),
    );
  }

  function finishToolboxButtonDrag(event: ReactPointerEvent<HTMLButtonElement>) {
    const dragState = dragStateRef.current;

    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId);
    }

    if (dragState?.moved) {
      suppressClickRef.current = true;
      window.setTimeout(() => {
        suppressClickRef.current = false;
      }, 0);
    }

    dragStateRef.current = null;
    setDragging(false);
  }

  function renderNavigationButton(item: ToolboxItem) {
    const isActive = item.id === activeItemId;

    return (
      <button
        aria-label={item.label}
        className={cx(
          "flex h-11 w-11 items-center justify-center rounded-md border text-lg shadow-sm transition-colors",
          tone.focusRing,
          isActive
            ? cx("border-transparent", tone.primary, tone.primaryForeground, tone.primaryHover)
            : cx(tone.border, tone.secondary, tone.mutedForeground, tone.secondaryHover),
        )}
        key={item.id}
        onClick={() => setActiveItem(item.id)}
        title={item.label}
        type="button"
      >
        {item.icon}
      </button>
    );
  }

  const modal = isOpen ? (
    <div
      className={cx("fixed inset-0 z-50 flex items-center justify-center px-4 py-8 backdrop-blur-sm", tone.overlay)}
      onClick={closeToolbox}
    >
      <div
        className={cx(
          "relative flex h-[min(68vh,520px)] w-full max-w-3xl overflow-hidden rounded-lg border",
          tone.border,
          tone.surface,
          tone.surfaceForeground,
          tone.shadow,
        )}
        onClick={(event) => event.stopPropagation()}
      >
        <aside className={cx("flex w-20 shrink-0 flex-col items-center border-r px-3 py-6", tone.border, tone.muted)}>
          <div className="flex flex-1 flex-col items-center gap-3">
            {mainItems.map((item) => renderNavigationButton(item))}
          </div>

          {profileItem ? (
            <div className={cx("mt-3 border-t pt-3", tone.border)}>
              {renderNavigationButton(profileItem)}
            </div>
          ) : null}
        </aside>

        <section className={cx("flex min-w-0 flex-1 flex-col", tone.surface)}>
          <header className={cx("flex items-center justify-between border-b px-6 py-4", tone.border)}>
            <div>
              <p className={cx("text-xs font-semibold uppercase", tone.subtleForeground)}>Toolbox</p>
              <h2 className={cx("mt-1 text-lg font-semibold", tone.foreground)}>{activeItem?.label ?? "工具箱"}</h2>
            </div>
            <button
              aria-label="Close toolbox"
              className={cx(
                "flex h-10 w-10 items-center justify-center rounded-md border transition-colors",
                tone.border,
                tone.focusRing,
                tone.secondary,
                tone.mutedForeground,
                tone.secondaryHover,
              )}
              onClick={closeToolbox}
              type="button"
            >
              <LuX size={18} />
            </button>
          </header>

          <div className="min-h-0 flex-1 overflow-y-auto px-6 py-5">
            {activeItem?.content ?? (
              <div className={cx("flex h-full items-center justify-center text-sm", tone.subtleForeground)}>
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
      <ToolboxButton
        dragging={dragging}
        onClick={handleToolboxButtonClick}
        onPointerCancel={finishToolboxButtonDrag}
        onPointerDown={handleToolboxButtonPointerDown}
        onPointerMove={handleToolboxButtonPointerMove}
        onPointerUp={finishToolboxButtonDrag}
        style={buttonPosition ?? getPlacementStyle(placement)}
        theme={theme}
      />
      {typeof document !== "undefined" ? createPortal(modal, document.body) : null}
    </>
  );
}

export { default as ToolboxActionsPanel } from "@/components/internal/ToolBox/ToolboxActionsPanel";
export { default as ToolboxButton } from "@/components/internal/ToolBox/ToolboxButton";
export { default as ToolboxEditorPanel } from "@/components/internal/ToolBox/ToolboxEditorPanel";
export { default as ToolboxFilesPanel } from "@/components/internal/ToolBox/ToolboxFilesPanel";
export { default as ToolboxNotesPanel } from "@/components/internal/ToolBox/ToolboxNotesPanel";
export { default as ToolboxOverviewPanel } from "@/components/internal/ToolBox/ToolboxOverviewPanel";
export { default as ToolboxProfilePanel } from "@/components/internal/ToolBox/ToolboxProfilePanel";
