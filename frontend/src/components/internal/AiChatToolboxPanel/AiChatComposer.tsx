import type { ChangeEvent, RefObject } from "react";
import { LuSendHorizontal, LuSquare, LuTrash2 } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import {
  cx,
  getThemeButtonClassName,
  getThemePrimaryButtonClassName,
  getThemeTone,
} from "@/app/themeStyles";
import AiChatAttachmentToolbar from "@/components/internal/AiChatToolboxPanel/AiChatAttachmentToolbar";
import type { ChatAttachment } from "@/components/internal/AiChatToolboxPanel/types";

type AiChatComposerProps = {
  attachments: ChatAttachment[];
  inputValue: string;
  isGenerating: boolean;
  onAttachmentsSelected: (attachments: ChatAttachment[]) => void;
  onClearComposer: () => void;
  onInputChange: (event: ChangeEvent<HTMLTextAreaElement>) => void;
  onRemoveAttachment: (attachmentId: string) => void;
  onSend: () => void;
  onStop: () => void;
  onUploadAttachments: () => void;
  onValidationError: (message: string) => void;
  textareaRef: RefObject<HTMLTextAreaElement | null>;
};

export default function AiChatComposer({
  attachments,
  inputValue,
  isGenerating,
  onAttachmentsSelected,
  onClearComposer,
  onInputChange,
  onRemoveAttachment,
  onSend,
  onStop,
  onUploadAttachments,
  onValidationError,
  textareaRef,
}: AiChatComposerProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return (
    <footer className={cx("border-t p-3", tone.border)}>
      <textarea
        ref={textareaRef}
        className={cx(
          "max-h-32 min-h-12 w-full resize-none rounded-md border px-3 py-2 text-sm leading-6",
          tone.border,
          tone.surface,
          tone.foreground,
          tone.focusRing,
        )}
        onChange={onInputChange}
        placeholder="输入消息..."
        rows={1}
        value={inputValue}
      />

      <div className="mt-3 flex items-center justify-between gap-3">
        <AiChatAttachmentToolbar
          attachments={attachments}
          onAttachmentsSelected={onAttachmentsSelected}
          onRemoveAttachment={onRemoveAttachment}
          onUploadAttachments={onUploadAttachments}
          onValidationError={onValidationError}
        />

        <div className="flex shrink-0 items-center gap-2">
          <button
            aria-label="Clear current chat input"
            className={cx("flex h-9 w-9 items-center justify-center", getThemeButtonClassName(theme))}
            disabled={!inputValue && attachments.length === 0}
            onClick={onClearComposer}
            type="button"
          >
            <LuTrash2 size={14} />
          </button>

          {isGenerating ? (
            <button
              className={cx("flex h-9 shrink-0 items-center gap-2 px-4 text-sm font-medium", tone.destructive, tone.destructiveForeground, tone.destructiveHover, tone.focusRing, "rounded-md")}
              onClick={onStop}
              type="button"
            >
              <LuSquare size={14} />
              终止
            </button>
          ) : (
            <button
              className={cx("flex h-9 shrink-0 items-center gap-2 px-4 text-sm font-medium", getThemePrimaryButtonClassName(theme))}
              disabled={!inputValue.trim() && attachments.length === 0}
              onClick={onSend}
              type="button"
            >
              <LuSendHorizontal size={14} />
              发送
            </button>
          )}
        </div>
      </div>
    </footer>
  );
}
