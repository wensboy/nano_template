import { useEffect, useMemo, useRef, useState } from "react";
import { LuBell, LuBellOff, LuSearch, LuTrash2 } from "react-icons/lu";

import type { ThemeTone } from "@/app/store/themeSlice";
import { cx } from "@/app/themeStyles";
import SaasIconButton from "@/components/custom/saas/SaasIconButton";
import SaasNotificationItem from "@/components/custom/saas/SaasNotificationItem";
import type { SaasNotification } from "@/components/custom/saas/types";

type SaasNotificationDrawerProps = {
  muted: boolean;
  notifications: SaasNotification[];
  onClose: () => void;
  onMutedChange: (muted: boolean) => void;
  onNotificationsChange: (notifications: SaasNotification[]) => void;
  open: boolean;
  tone: ThemeTone;
};

const DRAWER_ANIMATION_MS = 300;

export default function SaasNotificationDrawer({
  muted,
  notifications,
  onClose,
  onMutedChange,
  onNotificationsChange,
  open,
  tone,
}: SaasNotificationDrawerProps) {
  const [entered, setEntered] = useState(false);
  const [expandedIds, setExpandedIds] = useState<string[]>([]);
  const [mounted, setMounted] = useState(open);
  const [searchQuery, setSearchQuery] = useState("");
  const searchInputRef = useRef<HTMLInputElement | null>(null);
  const normalizedSearchQuery = searchQuery.trim().toLowerCase();
  const visibleNotifications = useMemo(
    () =>
      normalizedSearchQuery
        ? notifications.filter((notification) => notification.title.toLowerCase().includes(normalizedSearchQuery))
        : notifications,
    [normalizedSearchQuery, notifications],
  );

  useEffect(() => {
    if (open) {
      setMounted(true);

      const frame = window.requestAnimationFrame(() => {
        setEntered(true);
      });

      return () => {
        window.cancelAnimationFrame(frame);
      };
    }

    setEntered(false);

    const timeout = window.setTimeout(() => {
      setMounted(false);
    }, DRAWER_ANIMATION_MS);

    return () => {
      window.clearTimeout(timeout);
    };
  }, [open]);

  useEffect(() => {
    if (!entered) {
      return;
    }

    searchInputRef.current?.focus();
  }, [entered]);

  useEffect(() => {
    if (!mounted) {
      return;
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        onClose();
      }
    }

    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [mounted, onClose]);

  function deleteNotification(notificationId: string) {
    onNotificationsChange(notifications.filter((notification) => notification.id !== notificationId));
    setExpandedIds((ids) => ids.filter((id) => id !== notificationId));
  }

  function deleteVisibleNotifications() {
    const visibleIds = new Set(visibleNotifications.map((notification) => notification.id));

    onNotificationsChange(notifications.filter((notification) => !visibleIds.has(notification.id)));
    setExpandedIds((ids) => ids.filter((id) => !visibleIds.has(id)));
  }

  function toggleExpanded(notificationId: string) {
    setExpandedIds((ids) =>
      ids.includes(notificationId) ? ids.filter((id) => id !== notificationId) : [...ids, notificationId],
    );
  }

  if (!mounted) {
    return null;
  }

  return (
    <div
      aria-hidden={!open}
      className={cx(
        "relative z-30 col-start-2 row-start-1 row-span-2 h-full min-h-0 overflow-hidden",
        mounted ? "pointer-events-auto" : "pointer-events-none",
      )}
    >
      <button
        aria-label="Close notifications"
        className={cx(
          "absolute inset-0 cursor-default bg-zinc-950/55 transition-opacity duration-300 ease-out motion-reduce:transition-none",
          entered ? "opacity-100" : "opacity-0",
        )}
        onClick={onClose}
        type="button"
      />

      <aside
        aria-label="Notifications"
        className={cx(
          "absolute right-0 top-0 flex h-full w-full max-w-[26rem] flex-col border-l shadow-2xl transition-transform duration-300 ease-out motion-reduce:transition-none sm:w-[26rem]",
          entered ? "translate-x-0" : "translate-x-full",
          tone.border,
          tone.surface,
        )}
      >
        <div className={cx("border-b p-3", tone.border)}>
          <label className="relative block">
            <LuSearch className={cx("pointer-events-none absolute left-3 top-1/2 -translate-y-1/2", tone.iconMuted)} size={17} />
            <input
              aria-label="Search notification titles"
              className={cx(
                "h-9 w-full rounded-md border bg-transparent pl-9 pr-3 text-sm outline-none transition-colors placeholder:text-zinc-400",
                tone.border,
                tone.focusRing,
                tone.surfaceForeground,
              )}
              onChange={(event) => setSearchQuery(event.target.value)}
              placeholder="Search messages"
              ref={searchInputRef}
              type="search"
              value={searchQuery}
            />
          </label>
        </div>

        <div className="min-h-0 flex-1 overflow-y-auto">
          {visibleNotifications.length > 0 ? (
            visibleNotifications.map((notification) => (
              <SaasNotificationItem
                expanded={expandedIds.includes(notification.id)}
                key={notification.id}
                notification={notification}
                onDelete={deleteNotification}
                onToggleExpanded={toggleExpanded}
                tone={tone}
              />
            ))
          ) : (
            <div className={cx("border-b px-4 py-3 text-sm", tone.border, tone.mutedForeground)}>
              No matching messages.
            </div>
          )}
        </div>

        <footer className={cx("flex h-14 shrink-0 items-center gap-3 border-t px-4", tone.border)}>
          <span className="min-w-0 text-sm font-semibold tabular-nums">{visibleNotifications.length}</span>
          <div className="ml-auto flex shrink-0 items-center gap-2">
            <SaasIconButton
              action={{
                disabled: visibleNotifications.length === 0,
                icon: LuTrash2,
                id: "delete-visible-notifications",
                label: "Delete visible messages",
                onClick: deleteVisibleNotifications,
              }}
              tone={tone}
            />
            <SaasIconButton
              action={{
                icon: muted ? LuBellOff : LuBell,
                id: "toggle-do-not-disturb",
                label: muted ? "No disturb" : "Disturb",
                onClick: () => onMutedChange(!muted),
              }}
              tone={tone}
            />
          </div>
        </footer>
      </aside>
    </div>
  );
}
