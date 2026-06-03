import type { ThemeTone } from "@/app/store/themeSlice";
import { cx } from "@/app/themeStyles";
import type { SaasHeaderAction } from "@/components/custom/saas/types";

type SaasIconButtonProps = {
  action: SaasHeaderAction;
  tone: ThemeTone;
};

export default function SaasIconButton({ action, tone }: SaasIconButtonProps) {
  const Icon = action.icon;

  return (
    <button
      aria-label={action.label}
      className={cx(
        "flex h-9 w-9 shrink-0 items-center justify-center rounded-md border text-sm transition-colors",
        "disabled:cursor-not-allowed disabled:opacity-50",
        tone.border,
        tone.focusRing,
        tone.secondary,
        tone.mutedForeground,
        tone.secondaryHover,
      )}
      disabled={action.disabled}
      onClick={action.onClick}
      title={action.label}
      type="button"
    >
      <Icon className={action.iconClassName} size={18} />
    </button>
  );
}
