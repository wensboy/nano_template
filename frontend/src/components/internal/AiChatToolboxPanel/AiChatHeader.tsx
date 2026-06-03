import type { RefObject } from "react";
import { LuBrainCircuit, LuHistory, LuMessageSquareText, LuRadio, LuSettings2 } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeButtonClassName, getThemeTone } from "@/app/themeStyles";

type AiChatHeaderProps = {
  historyButtonRef: RefObject<HTMLButtonElement | null>;
  isStreamEnabled: boolean;
  isThinkingEnabled: boolean;
  messageHistoryButtonRef: RefObject<HTMLButtonElement | null>;
  model: string;
  onToggleHistory: () => void;
  onToggleMessageHistory: () => void;
  onToggleSettings: () => void;
  sessionName: string;
  settingsButtonRef: RefObject<HTMLButtonElement | null>;
};

export default function AiChatHeader({
  historyButtonRef,
  isStreamEnabled,
  isThinkingEnabled,
  messageHistoryButtonRef,
  model,
  onToggleHistory,
  onToggleMessageHistory,
  onToggleSettings,
  sessionName,
  settingsButtonRef,
}: AiChatHeaderProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return (
    <header className={cx("grid grid-cols-3 items-center gap-3 border-b px-3 py-2", tone.border)}>
      <div className="relative flex min-w-0 items-center gap-2">
        <button
          aria-label="Chat settings"
          className={cx("flex h-9 w-9 shrink-0 items-center justify-center", getThemeButtonClassName(theme))}
          onClick={onToggleSettings}
          ref={settingsButtonRef}
          type="button"
        >
          <LuSettings2 size={16} />
        </button>
        <span className={cx("truncate text-sm font-medium", tone.foreground)}>{model}</span>
        <span className="flex shrink-0 items-center gap-1 text-green-600">
          {isStreamEnabled ? <LuRadio aria-label="stream enabled" size={14} /> : null}
          {isThinkingEnabled ? <LuBrainCircuit aria-label="thinking enabled" size={14} /> : null}
        </span>
      </div>

      <div className="min-w-0 text-center">
        <p className={cx("truncate text-sm font-semibold", tone.foreground)} title={sessionName}>
          {sessionName}
        </p>
      </div>

      <div className="relative flex justify-end gap-2">
        <button
          aria-label="Message history"
          className={cx("flex h-9 w-9 items-center justify-center", getThemeButtonClassName(theme))}
          onClick={onToggleMessageHistory}
          ref={messageHistoryButtonRef}
          type="button"
        >
          <LuMessageSquareText size={16} />
        </button>
        <button
          aria-label="Chat history"
          className={cx("flex h-9 w-9 items-center justify-center", getThemeButtonClassName(theme))}
          onClick={onToggleHistory}
          ref={historyButtonRef}
          type="button"
        >
          <LuHistory size={16} />
        </button>
      </div>
    </header>
  );
}
