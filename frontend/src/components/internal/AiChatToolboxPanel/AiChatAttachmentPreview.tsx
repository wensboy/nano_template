import { useEffect, useState } from "react";
import { LuArrowLeft, LuFile } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeButtonClassName, getThemeTone } from "@/app/themeStyles";
import type { ChatAttachment } from "@/components/internal/AiChatToolboxPanel/types";

type AiChatAttachmentPreviewProps = {
  attachment: ChatAttachment;
  isClosing?: boolean;
  onBack: () => void;
  onExitComplete?: () => void;
};

function isTextPreview(attachment: ChatAttachment) {
  return (
    attachment.mimeType.startsWith("text/") ||
    attachment.mimeType.includes("json") ||
    attachment.mimeType.includes("xml") ||
    attachment.mimeType.includes("csv") ||
    attachment.mimeType.includes("markdown")
  );
}

export default function AiChatAttachmentPreview({
  attachment,
  isClosing = false,
  onBack,
  onExitComplete,
}: AiChatAttachmentPreviewProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const [objectUrl, setObjectUrl] = useState<string | null>(null);
  const [fileText, setFileText] = useState("");
  const sourceUrl = attachment.previewUrl ?? objectUrl;

  useEffect(() => {
    if (!attachment.file || attachment.previewUrl) {
      setObjectUrl(null);
      return;
    }

    const nextUrl = URL.createObjectURL(attachment.file);
    setObjectUrl(nextUrl);

    return () => URL.revokeObjectURL(nextUrl);
  }, [attachment.file, attachment.previewUrl]);

  useEffect(() => {
    if (attachment.previewText || !attachment.file || !isTextPreview(attachment)) {
      setFileText("");
      return;
    }

    let isActive = true;
    void attachment.file.text().then((text) => {
      if (isActive) {
        setFileText(text);
      }
    });

    return () => {
      isActive = false;
    };
  }, [attachment]);

  const textContent = attachment.previewText ?? fileText;

  return (
    <section
      className={cx(
        "absolute inset-0 z-10 flex min-h-0 flex-col",
        isClosing ? "ai-chat-preview-slide-out" : "ai-chat-preview-slide-in",
        tone.surface,
        tone.foreground,
      )}
      onAnimationEnd={(event) => {
        if (isClosing && event.currentTarget === event.target) {
          onExitComplete?.();
        }
      }}
    >
      <header className={cx("grid h-12 shrink-0 grid-cols-[40px_minmax(0,1fr)_40px] items-center border-b px-3", tone.border)}>
        <button
          aria-label="Back to chat"
          className={cx("flex h-8 w-8 items-center justify-center", getThemeButtonClassName(theme))}
          onClick={onBack}
          type="button"
        >
          <LuArrowLeft size={16} />
        </button>
        <h2 className="truncate px-2 text-center text-sm font-medium">{attachment.name}</h2>
        <div aria-hidden="true" />
      </header>

      <div className="min-h-0 flex-1 overflow-x-hidden overflow-y-auto p-4">
        {attachment.kind === "image" && sourceUrl ? (
          <div className="flex h-full min-h-0 items-center justify-center overflow-hidden">
            <img alt={attachment.name} className="max-h-full max-w-full object-contain" src={sourceUrl} />
          </div>
        ) : attachment.kind === "video" && sourceUrl ? (
          <div className="flex h-full min-h-0 items-center justify-center overflow-hidden">
            <video className="max-h-full max-w-full" controls src={sourceUrl} />
          </div>
        ) : textContent ? (
          <pre className={cx("whitespace-pre-wrap break-words rounded-md border p-4 text-sm leading-6", tone.border, tone.subtle)}>
            {textContent}
          </pre>
        ) : (
          <div className={cx("flex h-full min-h-52 flex-col items-center justify-center gap-3 text-center", tone.subtleForeground)}>
            <LuFile size={32} />
            <p className="text-sm">暂无可预览内容</p>
          </div>
        )}
      </div>
    </section>
  );
}
