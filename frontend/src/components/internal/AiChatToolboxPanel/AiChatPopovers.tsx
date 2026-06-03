import type { KeyboardEvent } from "react";
import { createPortal } from "react-dom";
import { LuBrainCircuit, LuPencil, LuRadio, LuSearch, LuTrash2 } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeButtonClassName, getThemeTone } from "@/app/themeStyles";
import { chatModels, temperatureMaxOptions } from "@/components/internal/AiChatToolboxPanel/constants";
import type { ChatSession, OutputRenderMode, PopoverPosition, RequestMethod } from "@/components/internal/AiChatToolboxPanel/types";

type SettingsPopoverProps = {
  apiKey: string;
  baseUrl: string;
  isOpen: boolean;
  isOutputRenderingEnabled: boolean;
  isStreamEnabled: boolean;
  isThinkingEnabled: boolean;
  maxTokens: number;
  model: string;
  onApiKeyChange: (value: string) => void;
  onBaseUrlChange: (value: string) => void;
  onMaxTokensChange: (value: number) => void;
  onModelChange: (value: string) => void;
  onOutputRenderModeChange: (value: OutputRenderMode) => void;
  onOutputRenderingEnabledChange: (value: boolean) => void;
  onRequestMethodChange: (value: RequestMethod) => void;
  onStreamEnabledChange: (value: boolean) => void;
  onTemperatureChange: (value: number) => void;
  onTemperatureMaxChange: (value: (typeof temperatureMaxOptions)[number]) => void;
  onThinkingEnabledChange: (value: boolean) => void;
  onTypewriterSpeedChange: (value: number) => void;
  outputRenderMode: OutputRenderMode;
  position: PopoverPosition | null;
  requestMethod: RequestMethod;
  temperature: number;
  temperatureMax: (typeof temperatureMaxOptions)[number];
  typewriterSpeedMs: number;
};

export function AiChatSettingsPopover({
  apiKey,
  baseUrl,
  isOpen,
  isOutputRenderingEnabled,
  isStreamEnabled,
  isThinkingEnabled,
  maxTokens,
  model,
  onApiKeyChange,
  onBaseUrlChange,
  onMaxTokensChange,
  onModelChange,
  onOutputRenderModeChange,
  onOutputRenderingEnabledChange,
  onRequestMethodChange,
  onStreamEnabledChange,
  onTemperatureChange,
  onTemperatureMaxChange,
  onThinkingEnabledChange,
  onTypewriterSpeedChange,
  outputRenderMode,
  position,
  requestMethod,
  temperature,
  temperatureMax,
  typewriterSpeedMs,
}: SettingsPopoverProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  if (!isOpen || !position || typeof document === "undefined") {
    return null;
  }

  return createPortal(
    <div
      className={cx(
        "fixed z-[100] max-h-52 w-72 overflow-y-auto rounded-lg border p-3 shadow-xl",
        tone.border,
        tone.surface,
        tone.surfaceForeground,
      )}
      style={{ left: position.left, top: position.top }}
    >
      <div className="space-y-4">
        <section className="space-y-3">
          <span className={cx("inline-flex h-6 items-center rounded-md border px-2 text-xs font-semibold uppercase", tone.border, tone.muted, tone.foreground)}>
            模型
          </span>
          <label className="block space-y-1.5 text-sm">
            <span className={tone.subtleForeground}>模型名称</span>
            <select
              className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
              onChange={(event) => onModelChange(event.target.value)}
              value={model}
            >
              {chatModels.map((chatModel) => (
                <option key={chatModel} value={chatModel}>
                  {chatModel}
                </option>
              ))}
            </select>
          </label>
        </section>

        <section className={cx("space-y-3 border-t pt-3", tone.border)}>
          <span className={cx("inline-flex h-6 items-center rounded-md border px-2 text-xs font-semibold uppercase", tone.border, tone.muted, tone.foreground)}>
            连接
          </span>
          <label className="flex items-center justify-between gap-3 text-sm">
            <span className={cx("flex items-center gap-1.5", tone.subtleForeground)}>
              <LuRadio size={14} />
              流式对话
            </span>
            <input
              checked={isStreamEnabled}
              className="h-4 w-4 accent-current"
              onChange={(event) => onStreamEnabledChange(event.target.checked)}
              type="checkbox"
            />
          </label>
          <label className="block space-y-1.5 text-sm">
            <span className={tone.subtleForeground}>调用方法</span>
            <select
              className={cx("h-9 w-full rounded-md border px-3 text-sm uppercase", tone.border, tone.surface, tone.foreground, tone.focusRing)}
              onChange={(event) => onRequestMethodChange(event.target.value as RequestMethod)}
              value={requestMethod}
            >
              <option value="post">POST</option>
              <option value="get">GET</option>
            </select>
          </label>
          <label className="block space-y-1.5 text-sm">
            <span className={tone.subtleForeground}>Base URL</span>
            <input
              className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
              onChange={(event) => onBaseUrlChange(event.target.value)}
              placeholder="https://api.openai.com/v1"
              value={baseUrl}
            />
          </label>
          <label className="block space-y-1.5 text-sm">
            <span className={tone.subtleForeground}>API Key</span>
            <input
              className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
              onChange={(event) => onApiKeyChange(event.target.value)}
              placeholder="sk-..."
              type="password"
              value={apiKey}
            />
          </label>
        </section>

        <section className={cx("space-y-3 border-t pt-3", tone.border)}>
          <span className={cx("inline-flex h-6 items-center rounded-md border px-2 text-xs font-semibold uppercase", tone.border, tone.muted, tone.foreground)}>
            生成参数
          </span>
          <label className="flex items-center justify-between gap-3 text-sm">
            <span className={cx("flex items-center gap-1.5", tone.subtleForeground)}>
              <LuBrainCircuit size={14} />
              思考模式
            </span>
            <input
              checked={isThinkingEnabled}
              className="h-4 w-4 accent-current"
              onChange={(event) => onThinkingEnabledChange(event.target.checked)}
              type="checkbox"
            />
          </label>
          <label className="block space-y-1.5 text-sm">
            <span className="flex items-center justify-between gap-3">
              <span className={tone.subtleForeground}>Temperature: {temperature.toFixed(1)}</span>
              <span className={cx("grid grid-cols-2 rounded-md border p-0.5", tone.border, tone.muted)}>
                {temperatureMaxOptions.map((maxValue) => (
                  <button
                    className={cx(
                      "h-6 rounded-sm px-2 text-xs transition-colors",
                      temperatureMax === maxValue
                        ? cx(tone.primary, tone.primaryForeground)
                        : cx(tone.mutedForeground, tone.secondaryHover),
                    )}
                    key={maxValue}
                    onClick={() => {
                      onTemperatureMaxChange(maxValue);
                      onTemperatureChange(Math.min(temperature, maxValue));
                    }}
                    type="button"
                  >
                    {maxValue.toFixed(1)}
                  </button>
                ))}
              </span>
            </span>
            <input
              className="w-full accent-current"
              max={temperatureMax}
              min="0"
              onChange={(event) => onTemperatureChange(Number(event.target.value))}
              step="0.1"
              type="range"
              value={temperature}
            />
          </label>
          <label className="block space-y-1.5 text-sm">
            <span className={tone.subtleForeground}>Max tokens</span>
            <input
              className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
              inputMode="numeric"
              onChange={(event) => onMaxTokensChange(Number(event.target.value))}
              value={maxTokens}
            />
          </label>
        </section>

        <section className={cx("space-y-3 border-t pt-3", tone.border)}>
          <span className={cx("inline-flex h-6 items-center rounded-md border px-2 text-xs font-semibold uppercase", tone.border, tone.muted, tone.foreground)}>
            输出
          </span>
          <label className="block space-y-1.5 text-sm">
            <span className={tone.subtleForeground}>打字速度: {typewriterSpeedMs}ms</span>
            <input
              className="w-full accent-current"
              max="80"
              min="4"
              onChange={(event) => onTypewriterSpeedChange(Number(event.target.value))}
              step="2"
              type="range"
              value={typewriterSpeedMs}
            />
          </label>
          <label className="flex items-center justify-between gap-3 text-sm">
            <span className={tone.subtleForeground}>输出渲染</span>
            <input
              checked={isOutputRenderingEnabled}
              className="h-4 w-4 accent-current"
              onChange={(event) => onOutputRenderingEnabledChange(event.target.checked)}
              type="checkbox"
            />
          </label>
          <label className="block space-y-1.5 text-sm">
            <span className={tone.subtleForeground}>渲染模式</span>
            <select
              className={cx("h-9 w-full rounded-md border px-3 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
              disabled={!isOutputRenderingEnabled}
              onChange={(event) => onOutputRenderModeChange(event.target.value as OutputRenderMode)}
              value={outputRenderMode}
            >
              <option value="instant">即时渲染</option>
              <option value="result">结果渲染</option>
            </select>
          </label>
        </section>
      </div>
    </div>,
    document.body,
  );
}

type MessageHistoryItem = {
  id: string;
  title: string;
};

type MessageHistoryPopoverProps = {
  items: MessageHistoryItem[];
  isOpen: boolean;
  onQueryChange: (value: string) => void;
  onSelectMessage: (messageId: string) => void;
  position: PopoverPosition | null;
  query: string;
};

export function AiChatMessageHistoryPopover({
  items,
  isOpen,
  onQueryChange,
  onSelectMessage,
  position,
  query,
}: MessageHistoryPopoverProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  if (!isOpen || !position || typeof document === "undefined") {
    return null;
  }

  return createPortal(
    <div
      className={cx("fixed z-[100] w-72 rounded-lg border p-3 shadow-xl", tone.border, tone.surface, tone.surfaceForeground)}
      style={{ left: position.left, top: position.top }}
    >
      <div className="relative">
        <LuSearch className={cx("absolute left-3 top-1/2 -translate-y-1/2", tone.iconMuted)} size={14} />
        <input
          className={cx("h-9 w-full rounded-md border px-3 pl-8 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
          onChange={(event) => onQueryChange(event.target.value)}
          placeholder="搜索消息索引"
          value={query}
        />
      </div>
      <div className="mt-3 max-h-52 space-y-1 overflow-y-auto pr-1">
        {items.length > 0 ? (
          items.map((item) => (
            <button
              className={cx("block w-full truncate rounded-md px-3 py-2 text-left text-sm transition-colors", tone.secondaryHover)}
              key={item.id}
              onClick={() => onSelectMessage(item.id)}
              title={item.title}
              type="button"
            >
              {item.title}
            </button>
          ))
        ) : (
          <div className={cx("rounded-md border px-3 py-6 text-center text-sm", tone.border, tone.subtleForeground)}>
            没有匹配的消息
          </div>
        )}
      </div>
    </div>,
    document.body,
  );
}

type SessionHistoryPopoverProps = {
  chatSessionsCount: number;
  currentSessionId: string;
  editingSessionId: string | null;
  editingSessionTitle: string;
  isOpen: boolean;
  onCancelEditSession: () => void;
  onCommitSessionTitle: () => void;
  onDeleteSession: (sessionId: string) => void;
  onEditingSessionTitleChange: (value: string) => void;
  onQueryChange: (value: string) => void;
  onSelectSession: (sessionId: string) => void;
  onStartEditSession: (session: ChatSession) => void;
  position: PopoverPosition | null;
  query: string;
  sessions: ChatSession[];
};

export function AiChatSessionHistoryPopover({
  chatSessionsCount,
  currentSessionId,
  editingSessionId,
  editingSessionTitle,
  isOpen,
  onCancelEditSession,
  onCommitSessionTitle,
  onDeleteSession,
  onEditingSessionTitleChange,
  onQueryChange,
  onSelectSession,
  onStartEditSession,
  position,
  query,
  sessions,
}: SessionHistoryPopoverProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  if (!isOpen || !position || typeof document === "undefined") {
    return null;
  }

  function handleEditKeyDown(event: KeyboardEvent<HTMLInputElement>) {
    if (event.key === "Enter") {
      onCommitSessionTitle();
    }
    if (event.key === "Escape") {
      onCancelEditSession();
    }
  }

  return createPortal(
    <div
      className={cx("fixed z-[100] w-72 rounded-lg border p-3 shadow-xl", tone.border, tone.surface, tone.surfaceForeground)}
      style={{ left: position.left, top: position.top }}
    >
      <div className="relative">
        <LuSearch className={cx("absolute left-3 top-1/2 -translate-y-1/2", tone.iconMuted)} size={14} />
        <input
          className={cx("h-9 w-full rounded-md border px-3 pl-8 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
          onChange={(event) => onQueryChange(event.target.value)}
          placeholder="搜索历史会话"
          value={query}
        />
      </div>
      <div className="mt-3 max-h-52 space-y-1 overflow-y-auto pr-1">
        {sessions.map((session) => (
          <div
            className={cx("flex min-w-0 items-center gap-1 rounded-md px-2 py-1.5 transition-colors", tone.secondaryHover)}
            key={session.id}
          >
            {editingSessionId === session.id ? (
              <input
                autoFocus
                className={cx(
                  "h-8 min-w-0 flex-1 rounded-md border px-2 text-sm",
                  tone.border,
                  tone.surface,
                  tone.foreground,
                  tone.focusRing,
                )}
                onBlur={onCommitSessionTitle}
                onChange={(event) => onEditingSessionTitleChange(event.target.value)}
                onKeyDown={handleEditKeyDown}
                value={editingSessionTitle}
              />
            ) : (
              <button
                className={cx(
                  "min-w-0 flex-1 truncate rounded-sm px-1 py-1 text-left text-sm",
                  session.id === currentSessionId ? "text-blue-600" : tone.foreground,
                )}
                onClick={() => onSelectSession(session.id)}
                title={session.title}
                type="button"
              >
                {session.title}
              </button>
            )}

            <div className="flex shrink-0 items-center gap-1">
              <button
                aria-label={`Edit ${session.title}`}
                className={cx("flex h-7 w-7 items-center justify-center", getThemeButtonClassName(theme))}
                onClick={() => onStartEditSession(session)}
                type="button"
              >
                <LuPencil size={13} />
              </button>
              <button
                aria-label={`Delete ${session.title}`}
                className={cx(
                  "flex h-7 w-7 items-center justify-center",
                  chatSessionsCount <= 1 ? "cursor-not-allowed opacity-50" : getThemeButtonClassName(theme),
                )}
                disabled={chatSessionsCount <= 1}
                onClick={() => onDeleteSession(session.id)}
                type="button"
              >
                <LuTrash2 className="text-red-600" size={13} />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>,
    document.body,
  );
}
