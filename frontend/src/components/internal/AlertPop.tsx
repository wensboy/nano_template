import { useEffect, useState, type ReactNode } from "react";
import {
  LuCircleAlert,
  LuCircleCheck,
  LuInfo,
  LuTriangleAlert,
} from "react-icons/lu";

type AlertPopVariant = "default" | "destructive" | "success" | "warning";

type AlertPopProps = {
  durationMs?: number;
  message?: ReactNode;
  variant?: AlertPopVariant;
};

const variantStyles: Record<AlertPopVariant, string> = {
  default: "border-zinc-200 bg-white text-zinc-950 shadow-lg",
  destructive: "border-red-200 bg-red-50 text-red-950 shadow-lg",
  success: "border-emerald-200 bg-emerald-50 text-emerald-950 shadow-lg",
  warning: "border-amber-200 bg-amber-50 text-amber-950 shadow-lg",
};

const iconStyles: Record<AlertPopVariant, string> = {
  default: "text-zinc-500",
  destructive: "text-red-600",
  success: "text-emerald-600",
  warning: "text-amber-600",
};

const icons = {
  default: LuInfo,
  destructive: LuCircleAlert,
  success: LuCircleCheck,
  warning: LuTriangleAlert,
} satisfies Record<AlertPopVariant, typeof LuInfo>;

export type { AlertPopVariant };

export default function AlertPop({ durationMs = 3000, message, variant = "default" }: AlertPopProps) {
  const [isVisible, setIsVisible] = useState(Boolean(message));

  useEffect(() => {
    if (!message) {
      setIsVisible(false);
      return;
    }

    setIsVisible(true);
    const timer = window.setTimeout(() => {
      setIsVisible(false);
    }, durationMs);

    return () => window.clearTimeout(timer);
  }, [durationMs, message]);

  if (!message || !isVisible) {
    return null;
  }

  const Icon = icons[variant];

  return (
    <div className="pointer-events-none fixed top-8 left-1/2 z-[100] w-[calc(100%-2rem)] max-w-md -translate-x-1/2 sm:top-10">
      <div
        className={`flex items-start gap-3 rounded-lg border px-4 py-3 text-sm shadow-sm ${variantStyles[variant]}`}
        role="alert"
      >
        <Icon className={`mt-0.5 shrink-0 ${iconStyles[variant]}`} size={16} />
        <div className="min-w-0 leading-5 break-words">{message}</div>
      </div>
    </div>
  );
}
