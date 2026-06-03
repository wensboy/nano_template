import { useEffect, useState } from "react";
import { LuCheck, LuCopy, LuDownload, LuEye, LuRefreshCcw } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeButtonClassName, getThemeTone } from "@/app/themeStyles";
import { AttachmentListPopover } from "@/components/internal/AiChatToolboxPanel/AiChatAttachmentToolbar";
import type { ChatAttachment, ChatMessage, OutputRenderMode } from "@/components/internal/AiChatToolboxPanel/types";
import MarkdownContent from "@/components/internal/AiChatToolboxPanel/MarkdownContent";

type ChatMessageFooterProps = {
  attachments?: ChatAttachment[];
  content: string;
  messageId: string;
  onDownloadAttachment?: (attachment: ChatAttachment) => void;
  onPreviewAttachment?: (attachment: ChatAttachment) => void;
  showAttachments?: boolean;
  showRefresh?: boolean;
};

function CopyMessageButton({ content }: { content: string }) {
  const theme = useAppSelector((state) => state.theme);
  const [isCopied, setIsCopied] = useState(false);

  async function handleCopy() {
    await navigator.clipboard?.writeText(content);
    setIsCopied(true);
  }

  useEffect(() => {
    if (!isCopied) {
      return;
    }

    const timeoutId = window.setTimeout(() => setIsCopied(false), 1200);

    return () => window.clearTimeout(timeoutId);
  }, [isCopied]);

  return (
    <button
      aria-label={isCopied ? "Message copied" : "Copy message"}
      className={cx("flex h-7 w-7 items-center justify-center", getThemeButtonClassName(theme))}
      onClick={() => void handleCopy()}
      type="button"
    >
      {isCopied ? <LuCheck size={13} /> : <LuCopy size={13} />}
    </button>
  );
}

function ChatMessageFooter({
  attachments = [],
  content,
  messageId,
  onDownloadAttachment,
  onPreviewAttachment,
  showAttachments = false,
  showRefresh = false,
}: ChatMessageFooterProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const hasAttachments = showAttachments && attachments.length > 0;
  const renderAttachmentActions =
    onPreviewAttachment && onDownloadAttachment
      ? (attachment: ChatAttachment, closePopover: () => void) => (
          <>
            <button
              aria-label={`Preview ${attachment.name}`}
              className={cx("flex h-7 w-7 items-center justify-center", getThemeButtonClassName(theme))}
              onClick={() => {
                closePopover();
                onPreviewAttachment(attachment);
              }}
              type="button"
            >
              <LuEye size={13} />
            </button>
            <button
              aria-label={`Download ${attachment.name}`}
              className={cx("flex h-7 w-7 items-center justify-center", getThemeButtonClassName(theme))}
              onClick={() => onDownloadAttachment(attachment)}
              type="button"
            >
              <LuDownload size={13} />
            </button>
          </>
        )
      : null;

  function handleRefresh() {
    // Reserved for retrying the same chat task.
  }

  return (
    <div
      className={cx(
        "invisible mt-1 border-t pt-1 opacity-0 transition-opacity duration-150 group-hover:visible group-hover:opacity-100 group-focus-within:visible group-focus-within:opacity-100",
        tone.border,
      )}
    >
      <div className="flex min-w-0 items-center justify-between gap-2">
        <div className="min-w-0">
          {hasAttachments ? (
            <AttachmentListPopover
              attachments={attachments}
              closeOnOutsideClick
              id={`ai-chat-message-attachments-${messageId}`}
              {...(renderAttachmentActions ? { renderAttachmentActions } : {})}
              size="compact"
            />
          ) : null}
        </div>
        <div className="flex shrink-0 items-center gap-1">
          <CopyMessageButton content={content} />
          {showRefresh ? (
            <button
              aria-label="Refresh message task"
              className={cx("flex h-7 w-7 items-center justify-center", getThemeButtonClassName(theme))}
              onClick={handleRefresh}
              type="button"
            >
              <LuRefreshCcw size={13} />
            </button>
          ) : null}
        </div>
      </div>
    </div>
  );
}

type ChatMessageBubbleProps = {
  message: ChatMessage;
};

type AssistantChatMessageBubbleProps = ChatMessageBubbleProps & {
  isOutputRenderingEnabled: boolean;
  isStreamEnabled: boolean;
  onDownloadAttachment?: (attachment: ChatAttachment) => void;
  onPreviewAttachment?: (attachment: ChatAttachment) => void;
  outputRenderMode: OutputRenderMode;
  typewriterSpeedMs: number;
};

function useTypewriterText(content: string, isEnabled: boolean, speedMs: number) {
  const [visibleText, setVisibleText] = useState(isEnabled ? "" : content);

  useEffect(() => {
    if (!isEnabled) {
      setVisibleText(content);
      return;
    }

    setVisibleText("");
    let cursor = 0;
    const intervalId = window.setInterval(() => {
      cursor += 1;
      setVisibleText(content.slice(0, cursor));

      if (cursor >= content.length) {
        window.clearInterval(intervalId);
      }
    }, speedMs);

    return () => window.clearInterval(intervalId);
  }, [content, isEnabled, speedMs]);

  return visibleText;
}

export function AssistantChatMessageBubble({
  isOutputRenderingEnabled,
  isStreamEnabled,
  message,
  onDownloadAttachment,
  onPreviewAttachment,
  outputRenderMode,
  typewriterSpeedMs,
}: AssistantChatMessageBubbleProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const shouldTypewrite = isStreamEnabled;
  const visibleContent = useTypewriterText(message.content, shouldTypewrite, typewriterSpeedMs);
  const displayContent = shouldTypewrite ? visibleContent : message.content;
  const shouldRenderMarkdown =
    isOutputRenderingEnabled && (!shouldTypewrite || outputRenderMode === "instant" || displayContent.length === message.content.length);

  return (
    <div className="flex justify-start">
      <div
        className={cx(
          "group max-w-[88%] rounded-lg border px-4 py-3 text-sm leading-6",
          tone.border,
          tone.subtle,
          tone.mutedForeground,
        )}
      >
        {shouldRenderMarkdown ? (
          <MarkdownContent content={displayContent} />
        ) : (
          <p className="whitespace-pre-wrap break-words">{displayContent}</p>
        )}
        <ChatMessageFooter
          attachments={message.attachments ?? []}
          content={message.content}
          messageId={message.id}
          onDownloadAttachment={onDownloadAttachment}
          onPreviewAttachment={onPreviewAttachment}
          showAttachments
        />
      </div>
    </div>
  );
}

export function UserChatMessageBubble({ message }: ChatMessageBubbleProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return (
    <div className="flex justify-end">
      <div
        className={cx(
          "group max-w-[88%] rounded-lg border px-4 py-3 text-sm leading-6",
          tone.border,
          tone.muted,
          tone.foreground,
        )}
      >
        <p className="whitespace-pre-wrap break-words">{message.content}</p>
        <ChatMessageFooter
          attachments={message.attachments ?? []}
          content={message.content}
          messageId={message.id}
          showAttachments
          showRefresh
        />
      </div>
    </div>
  );
}
