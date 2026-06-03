import { useEffect, useMemo, useState } from "react";

import type { ThemeTone } from "@/app/store/themeSlice";
import { cx } from "@/app/themeStyles";
import SaasViewLayout, { type SaasViewLayoutNavItem } from "@/components/custom/saas/SaasViewLayout";
import type { SaasNavItem } from "@/components/custom/saas/types";

type SaasBodyProps = {
  activeItem: SaasNavItem;
  searchQuery: string;
  tone: ThemeTone;
};

const VIEW_SECTION_FIXTURES = [
  "Overview",
  "Queue",
  "Approvals",
  "Long-running operations with an intentionally verbose section label",
  "Segments",
  "Automation",
  "Ownership",
  "Escalations",
  "Exports",
  "Audit Trail",
  "Webhook retries",
  "Regional rollout status",
  "Archived edge cases",
  "SLA exceptions",
  "Customer handoff notes",
  "Forecast",
  "Permissions",
  "Import validation",
  "Empty-state review",
  "Rate limits",
  "Data retention",
  "Incident response",
  "Release readiness",
  "Backlog",
  "Final checkpoint with another long label for truncation coverage",
];

export default function SaasBody({ activeItem, searchQuery, tone }: SaasBodyProps) {
  const [activeViewItemId, setActiveViewItemId] = useState("");
  const normalizedSearchQuery = searchQuery.trim().toLowerCase();

  const viewItems = useMemo<SaasViewLayoutNavItem[]>(
    () =>
      VIEW_SECTION_FIXTURES.map((section, index) => ({
        description:
          index % 5 === 3
            ? `${activeItem.label} ${section} covers high-volume records, long copy, and scroll-heavy review states.`
            : `${activeItem.label} ${section} keeps the current workspace context focused and inspectable.`,
        id: `${activeItem.id}-${index}`,
        label: section,
        meta: index === 0 ? "Default" : `${index + 3} records`,
      })),
    [activeItem.id, activeItem.label],
  );

  const visibleViewItems = useMemo(
    () =>
      normalizedSearchQuery
        ? viewItems.filter((item) =>
            [activeItem.label, activeItem.description, item.label, item.description, item.meta].some((value) =>
              value.toLowerCase().includes(normalizedSearchQuery),
            ),
          )
        : viewItems,
    [activeItem.description, activeItem.label, normalizedSearchQuery, viewItems],
  );

  useEffect(() => {
    const nextActiveItem = visibleViewItems.find((item) => item.id === activeViewItemId) ?? visibleViewItems[0] ?? null;

    setActiveViewItemId(nextActiveItem?.id ?? "");
  }, [activeViewItemId, visibleViewItems]);

  function renderContent(viewItem: SaasViewLayoutNavItem) {
    const rows = [
      { label: `${viewItem.label} primary workflow`, owner: "Core Team", status: "Healthy" },
      {
        label: `${viewItem.label} review item with a deliberately long title for table overflow coverage`,
        owner: "Ops",
        status: "Needs review",
      },
      { label: `${viewItem.label} export`, owner: "Data", status: "Ready" },
      { label: `${viewItem.label} archived checkpoint`, owner: "Support", status: "Blocked" },
      { label: `${viewItem.label} final validation`, owner: "Product", status: "Queued" },
    ];

    return (
      <div className="flex min-h-full flex-col p-5">
        <div className={cx("border-b pb-4", tone.border)}>
          <p className={cx("text-xs font-semibold uppercase", tone.subtleForeground)}>Selected Section</p>
          <h2 className="mt-1 text-xl font-semibold leading-7">{viewItem.label}</h2>
          <p className={cx("mt-1 max-w-3xl text-sm leading-6", tone.mutedForeground)}>{viewItem.description}</p>
        </div>

        <div className="grid gap-0 border-b sm:grid-cols-3">
          {rows.slice(0, 3).map((row, index) => (
            <div
              className={cx("min-w-0 px-4 py-4", index > 0 && "border-t sm:border-t-0 sm:border-l", tone.border)}
              key={row.label}
            >
              <p className="truncate text-sm font-medium">{row.label}</p>
              <p className={cx("mt-1 text-xs", tone.subtleForeground)}>{row.owner}</p>
              <p className="mt-3 text-sm font-semibold">{row.status}</p>
            </div>
          ))}
        </div>

        <div
          className={cx(
            "grid grid-cols-[minmax(0,1fr)_8rem_8rem] border-b px-4 py-3 text-xs font-semibold uppercase",
            tone.border,
            tone.subtleForeground,
          )}
        >
          <span>Item</span>
          <span>Owner</span>
          <span>Status</span>
        </div>
        {rows.map((row) => (
          <div
            className={cx(
              "grid grid-cols-[minmax(0,1fr)_8rem_8rem] items-center border-b px-4 py-3 text-sm last:border-b-0",
              tone.border,
            )}
            key={`${viewItem.id}-${row.label}`}
          >
            <span className="min-w-0 truncate font-medium">{row.label}</span>
            <span className={cx("truncate", tone.mutedForeground)}>{row.owner}</span>
            <span className="truncate">{row.status}</span>
          </div>
        ))}

        <div className={cx("mt-5 rounded-md border px-4 py-3 text-sm leading-6", tone.border, tone.mutedForeground)}>
          This content intentionally has enough structure to test nested scrolling, long labels, and table truncation
          while the left text navigation keeps the active section pinned during scroll.
        </div>
      </div>
    );
  }

  const description = normalizedSearchQuery
    ? `${activeItem.description} Search is scoped to this view and currently matches ${visibleViewItems.length} of ${viewItems.length} sections.`
    : activeItem.description;

  return (
    <main className="col-start-2 row-start-2 min-h-0 overflow-hidden">
      <div className="h-full w-full p-5">
        <SaasViewLayout
          activeItemId={activeViewItemId}
          description={description}
          emptyMessage={`No matching sections in ${activeItem.label}.`}
          items={visibleViewItems}
          label={activeItem.label}
          onActiveItemChange={setActiveViewItemId}
          renderContent={renderContent}
          tone={tone}
        />
      </div>
    </main>
  );
}
