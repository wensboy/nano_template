import { useMemo, useState } from "react";
import {
  LuBell,
  LuChartColumnIncreasing,
  LuLayoutDashboard,
  LuBellOff,
  LuLogs,
  LuLogOut,
  LuMail,
  LuMessageSquareText,
  LuMoon,
  LuSettings,
  LuSun,
  LuUserRound,
} from "react-icons/lu";

import { useAppDispatch, useAppSelector } from "@/app/hooks";
import { toggleThemeMode } from "@/app/store/themeSlice";
import { cx, getThemeTone } from "@/app/themeStyles";
import SaasBody from "@/components/custom/saas/SaasBody";
import SaasHeader from "@/components/custom/saas/SaasHeader";
import SaasNotificationDrawer from "@/components/custom/saas/SaasNotificationDrawer";
import SaasSidebar from "@/components/custom/saas/SaasSidebar";
import type { SaasHeaderAction, SaasNavItem, SaasNotification } from "@/components/custom/saas/types";

type SaasLayoutProps = {
  onLogout?: () => void;
};

const NAV_ITEMS: SaasNavItem[] = [
  {
    description: "Track operating signals, pending work, and workspace health in one focused surface.",
    icon: LuLayoutDashboard,
    id: "dashboard",
    label: "Dashboard",
  },
  {
    description: "Inspect usage, conversion, and business metrics without leaving the app shell.",
    icon: LuChartColumnIncreasing,
    id: "analytics",
    label: "Analytics",
  },
  {
    description: "Review workspace activity, system events, and operational history in order.",
    icon: LuLogs,
    id: "log",
    label: "Log",
  },
];

const UTILITY_ITEMS: SaasNavItem[] = [
  {
    description: "Tune workspace preferences, security defaults, and product-level controls.",
    icon: LuSettings,
    id: "settings",
    label: "Settings",
  },
  {
    description: "View profile details, personal preferences, and account-level activity.",
    icon: LuUserRound,
    id: "profile",
    label: "Profile",
  },
];

const INITIAL_NOTIFICATIONS: SaasNotification[] = [
  {
    body: "The delivery lane for the workspace review has three customer notes waiting on owner confirmation before the next checkpoint can be marked ready. This sample intentionally runs past three compact lines so the collapsed state exercises truncation and the expand action has meaningful space to reveal.",
    icon: LuMessageSquareText,
    id: "workspace-review",
    time: "2026-06-01 09:42",
    title: "Workspace review needs owner confirmation before the next checkpoint can be marked ready",
  },
  {
    body: "A new billing document was generated and attached to the account record. Review the document metadata before sending it to the customer contact.",
    icon: LuMail,
    id: "billing-document",
    time: "2026-06-01 10:15",
    title: "Billing document is ready for review",
  },
  {
    body: "Analytics ingestion finished for the current reporting window.",
    id: "analytics-ingestion",
    time: "2026-06-01 11:08",
    title: "Analytics ingestion completed",
  },
  {
    body: "Organization access changes were queued by the admin team. Validate role updates before the next scheduled sync runs.",
    icon: LuBell,
    id: "organization-access",
    time: "2026-06-01 12:31",
    title: "Organization access changes queued",
  },
  {
    body: "The project intake queue received six new customer requests. Assign reviewers before the daily triage window closes.",
    icon: LuMessageSquareText,
    id: "project-intake",
    time: "2026-06-01 13:04",
    title: "Project intake queue received new requests",
  },
  {
    body: "A profile security update was detected from a trusted device. No action is required unless the activity looks unfamiliar.",
    id: "profile-security",
    time: "2026-06-01 13:28",
    title: "Profile security update detected",
  },
  {
    body: "Document export finished with warnings for two archived records. The active customer documents were exported successfully.",
    icon: LuMail,
    id: "document-export",
    time: "2026-06-01 14:02",
    title: "Document export finished with warnings",
  },
  {
    body: "Dashboard health signals moved below the configured threshold for workspace response time. Review the operational summary before acknowledging.",
    icon: LuBell,
    id: "dashboard-health",
    time: "2026-06-01 14:37",
    title: "Dashboard health signal needs review",
  },
  {
    body: "Settings changes for notification preferences were saved. The new defaults apply to future workspace members only.",
    id: "settings-notifications",
    time: "2026-06-01 15:11",
    title: "Notification preference settings saved",
  },
  {
    body: "A customer workspace moved from implementation to active support. Confirm ownership and update the handoff notes when the team is ready.",
    icon: LuMessageSquareText,
    id: "support-handoff",
    time: "2026-06-01 15:44",
    title: "Customer workspace moved to active support",
  },
  {
    body: "Analytics anomaly detection found an unusual spike in conversion events. Compare this with the campaign calendar before creating an incident.",
    icon: LuBell,
    id: "analytics-anomaly",
    time: "2026-06-01 16:20",
    title: "Analytics anomaly detected in conversion events",
  },
  {
    body: "The organization invite batch completed. Three invitations were delivered, one address was skipped because it already belongs to an active member, and one invite remains pending retry.",
    icon: LuMail,
    id: "invite-batch",
    time: "2026-06-01 16:58",
    title: "Organization invite batch completed",
  },
  {
    body: "The monthly workspace report is available for review.",
    id: "monthly-report",
    time: "2026-06-01 17:23",
    title: "Monthly workspace report available",
  },
  {
    body: "A scheduled sync was delayed by upstream rate limits. The next retry is queued and should resume without manual intervention.",
    icon: LuBell,
    id: "scheduled-sync-delay",
    time: "2026-06-01 18:05",
    title: "Scheduled sync delayed by upstream limits",
  },
  {
    body: "The account document request includes a very long title and medium-length body so the drawer can be checked for horizontal stability, title truncation, and scrolling behavior while multiple messages remain visible in a compact list.",
    icon: LuMessageSquareText,
    id: "long-title-document-request",
    time: "2026-06-01 18:41",
    title: "Account document request waiting on customer procurement contact and internal approval routing before final delivery",
  },
  {
    body: "End-of-day digest generated for workspace activity, project updates, analytics changes, and organization access events.",
    icon: LuMail,
    id: "daily-digest",
    time: "2026-06-01 19:00",
    title: "End-of-day digest generated",
  },
];

export default function SaasLayout({ onLogout }: SaasLayoutProps) {
  const [activeItemId, setActiveItemId] = useState(NAV_ITEMS[0]?.id ?? "");
  const [navigationHistory, setNavigationHistory] = useState<string[]>([]);
  const [collapsed, setCollapsed] = useState(false);
  const [notifications, setNotifications] = useState(INITIAL_NOTIFICATIONS);
  const [notificationDrawerOpen, setNotificationDrawerOpen] = useState(false);
  const [notificationsMuted, setNotificationsMuted] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const dispatch = useAppDispatch();
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const isDark = theme.mode === "dark";
  const allItems = [...NAV_ITEMS, ...UTILITY_ITEMS];
  const activeItem = allItems.find((item) => item.id === activeItemId) ?? allItems[0];
  const homeItemId = NAV_ITEMS[0]?.id ?? activeItem.id;
  const canNavigateBack = navigationHistory.length > 0;

  function navigateToItem(itemId: string) {
    if (itemId === activeItemId) {
      return;
    }

    setNavigationHistory((history) => [...history, activeItemId]);
    setActiveItemId(itemId);
  }

  function navigateBack() {
    setNavigationHistory((history) => {
      const previousItemId = history.at(-1);

      if (!previousItemId) {
        return history;
      }

      setActiveItemId(previousItemId);
      return history.slice(0, -1);
    });
  }

  const headerActions = useMemo(() => {
    const actions: SaasHeaderAction[] = [
      {
        icon: notificationsMuted ? LuBellOff : LuBell,
        id: "notifications",
        label: notificationsMuted ? "Notifications muted" : "Notifications",
        onClick: () => setNotificationDrawerOpen(true),
      },
      {
        icon: isDark ? LuSun : LuMoon,
        id: "theme",
        label: isDark ? "Switch to light mode" : "Switch to dark mode",
        onClick: () => dispatch(toggleThemeMode()),
      },
    ];

    if (onLogout) {
      actions.push({ icon: LuLogOut, id: "logout", label: "Log out", onClick: onLogout });
    }

    return actions;
  }, [dispatch, isDark, notificationsMuted, onLogout]);

  return (
    <div
      className={cx(
        "grid h-screen overflow-hidden transition-[background-color,color,grid-template-columns] duration-300 ease-out motion-reduce:transition-none",
        tone.page,
        tone.foreground,
      )}
      style={{
        gridTemplateColumns: `${collapsed ? "4.75rem" : "17rem"} minmax(0, 1fr)`,
        gridTemplateRows: "4.5rem minmax(0, 1fr)",
      }}
    >
      <SaasSidebar
        activeItemId={activeItem.id}
        collapsed={collapsed}
        items={NAV_ITEMS}
        onSelectItem={navigateToItem}
        onToggleCollapsed={() => setCollapsed((value) => !value)}
        tone={tone}
        utilityItems={UTILITY_ITEMS}
      />
      <SaasHeader
        actions={headerActions}
        activeItem={activeItem}
        canNavigateBack={canNavigateBack}
        onNavigateBack={navigateBack}
        onNavigateHome={() => navigateToItem(homeItemId)}
        onSearchQueryChange={setSearchQuery}
        searchQuery={searchQuery}
        tone={tone}
      />
      <SaasBody activeItem={activeItem} searchQuery={searchQuery} tone={tone} />
      <SaasNotificationDrawer
        muted={notificationsMuted}
        notifications={notifications}
        onClose={() => setNotificationDrawerOpen(false)}
        onMutedChange={setNotificationsMuted}
        onNotificationsChange={setNotifications}
        open={notificationDrawerOpen}
        tone={tone}
      />
    </div>
  );
}
