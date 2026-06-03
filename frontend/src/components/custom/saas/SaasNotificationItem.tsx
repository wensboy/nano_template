import { LuChevronDown, LuTrash2 } from "react-icons/lu";

import type { ThemeTone } from "@/app/store/themeSlice";
import { cx } from "@/app/themeStyles";
import SaasIconButton from "@/components/custom/saas/SaasIconButton";
import type { SaasNotification } from "@/components/custom/saas/types";

type SaasNotificationItemProps = {
  expanded: boolean;
  notification: SaasNotification;
  onDelete: (notificationId: string) => void;
  onToggleExpanded: (notificationId: string) => void;
  tone: ThemeTone;
};

export default function SaasNotificationItem({
  expanded,
  notification,
  onDelete,
  onToggleExpanded,
  tone,
}: SaasNotificationItemProps) {
  const MessageIcon = notification.icon;

  return (
    <article className={cx("border-b", tone.border)}>
      <div className="flex min-w-0 items-center gap-2 px-4 py-2">
        {MessageIcon ? <MessageIcon className={cx("shrink-0", tone.iconMuted)} size={17} /> : null}
        <h3 className="min-w-0 truncate text-sm font-medium leading-5">{notification.title}</h3>
      </div>

      <div className={cx("border-t border-b px-4 py-2", tone.border)}>
        <div
          className={cx(
            "overflow-hidden transition-[max-height] duration-300 ease-out motion-reduce:transition-none",
            expanded ? "max-h-64" : "max-h-[3.75rem]",
          )}
        >
          <p
            className={cx(
              "text-sm leading-5",
              tone.mutedForeground,
              expanded ? "whitespace-pre-line" : "overflow-hidden [display:-webkit-box] [-webkit-box-orient:vertical] [-webkit-line-clamp:3]",
            )}
          >
            {notification.body}
          </p>
        </div>
      </div>

      <div className="flex min-w-0 items-center gap-3 px-4 py-2">
        <time className={cx("min-w-0 truncate text-xs", tone.subtleForeground)}>{notification.time}</time>
        <div className="ml-auto flex shrink-0 items-center gap-2">
          <SaasIconButton
            action={{
              icon: LuChevronDown,
              iconClassName: cx(
                "transition-transform duration-200 ease-out motion-reduce:transition-none",
                expanded ? "rotate-180" : "rotate-0",
              ),
              id: `${notification.id}-expand`,
              label: expanded ? "Collapse message" : "Expand message",
              onClick: () => onToggleExpanded(notification.id),
            }}
            tone={tone}
          />
          <SaasIconButton
            action={{
              icon: LuTrash2,
              id: `${notification.id}-delete`,
              label: "Delete message",
              onClick: () => onDelete(notification.id),
            }}
            tone={tone}
          />
        </div>
      </div>
    </article>
  );
}
