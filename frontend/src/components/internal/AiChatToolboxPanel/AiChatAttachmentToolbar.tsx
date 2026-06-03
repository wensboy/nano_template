import { useEffect, useLayoutEffect, useRef, useState, type ChangeEvent, type ReactNode } from "react";
import { createPortal } from "react-dom";
import { LuChevronDown, LuChevronUp, LuFile, LuImage, LuMic, LuPlus, LuSearch, LuUpload, LuVideo, LuX } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeButtonClassName, getThemeTone } from "@/app/themeStyles";
import type { ChatAttachment } from "@/components/internal/AiChatToolboxPanel/types";
import { createChatAttachment, isFileMatchedKind } from "@/components/internal/AiChatToolboxPanel/utils";

const popoverGap = 8;

const attachmentToolItems = [
  { accept: "", icon: LuFile, kind: "file", label: "文件", selectable: true },
  { accept: "image/*", icon: LuImage, kind: "image", label: "图像", selectable: true },
  { accept: "video/*", icon: LuVideo, kind: "video", label: "视频", selectable: true },
  { accept: "audio/*", icon: LuMic, kind: "audio", label: "语音", selectable: false },
] satisfies Array<{
  accept: string;
  icon: typeof LuFile;
  kind: ChatAttachment["kind"];
  label: string;
  selectable: boolean;
}>;

function getAttachmentIcon(kind: ChatAttachment["kind"]) {
  return attachmentToolItems.find((item) => item.kind === kind)?.icon ?? LuFile;
}

type AttachmentInfoItemProps = {
  attachment: ChatAttachment;
  onRemove?: (attachmentId: string) => void;
  renderActions?: (attachment: ChatAttachment, closePopover: () => void) => ReactNode;
  onClosePopover: () => void;
};

function AttachmentInfoItem({ attachment, onClosePopover, onRemove, renderActions }: AttachmentInfoItemProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const Icon = getAttachmentIcon(attachment.kind);
  const actions = renderActions?.(attachment, onClosePopover);

  return (
    <div className={cx("group relative flex items-center gap-3 rounded-md border p-3", tone.border, tone.subtle)}>
      <div className="flex h-10 w-10 shrink-0 items-center justify-center">
        <Icon className={tone.iconMuted} size={22} />
      </div>
      <div className={cx("min-w-0 flex-1", onRemove ? "pr-5" : undefined)}>
        <p className={cx("truncate text-sm font-medium", tone.foreground)}>{attachment.name}</p>
        <p className={cx("mt-0.5 text-xs", tone.subtleForeground)}>{attachment.sizeLabel}</p>
      </div>
      {actions ? <div className="flex shrink-0 items-center gap-1">{actions}</div> : null}
      {onRemove ? (
        <button
          aria-label={`Remove ${attachment.name}`}
          className={cx(
            "absolute right-1.5 top-1.5 flex h-6 w-6 items-center justify-center rounded-md opacity-0 transition-colors group-hover:opacity-100 focus-visible:opacity-100",
            tone.destructive,
            tone.destructiveForeground,
            tone.destructiveHover,
            tone.focusRing,
          )}
          onClick={() => onRemove(attachment.id)}
          type="button"
        >
          <LuX size={14} />
        </button>
      ) : null}
    </div>
  );
}

type AttachmentListPopoverProps = {
  attachments: ChatAttachment[];
  closeOnOutsideClick?: boolean;
  id: string;
  onRemoveAttachment?: (attachmentId: string) => void;
  onUploadAttachments?: () => void;
  renderAttachmentActions?: (attachment: ChatAttachment, closePopover: () => void) => ReactNode;
  size?: "compact" | "default";
};

export function AttachmentListPopover({
  attachments,
  closeOnOutsideClick = false,
  id,
  onRemoveAttachment,
  onUploadAttachments,
  renderAttachmentActions,
  size = "default",
}: AttachmentListPopoverProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const [isOpen, setIsOpen] = useState(false);
  const [popoverPosition, setPopoverPosition] = useState<{ bottom: number; left: number } | null>(null);
  const [query, setQuery] = useState("");
  const popoverRef = useRef<HTMLDivElement>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);
  const filteredAttachments = attachments.filter((attachment) =>
    attachment.name.toLowerCase().includes(query.trim().toLowerCase()),
  );
  const isCompact = size === "compact";

  function updateAttachmentPopoverPosition() {
    if (!triggerRef.current) {
      return;
    }

    const rect = triggerRef.current.getBoundingClientRect();
    const viewportPadding = 16;
    const popoverWidth = 320;
    const maxLeft = window.innerWidth - popoverWidth - viewportPadding;

    setPopoverPosition({
      bottom: window.innerHeight - rect.top + popoverGap,
      left: Math.min(Math.max(rect.left, viewportPadding), Math.max(maxLeft, viewportPadding)),
    });
  }

  useLayoutEffect(() => {
    if (isOpen) {
      updateAttachmentPopoverPosition();
    }
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen) {
      return;
    }

    window.addEventListener("resize", updateAttachmentPopoverPosition);
    window.addEventListener("scroll", updateAttachmentPopoverPosition, true);

    return () => {
      window.removeEventListener("resize", updateAttachmentPopoverPosition);
      window.removeEventListener("scroll", updateAttachmentPopoverPosition, true);
    };
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen || !closeOnOutsideClick) {
      return;
    }

    function handlePointerDown(event: PointerEvent) {
      const target = event.target;

      if (!(target instanceof Node)) {
        return;
      }

      if (triggerRef.current?.contains(target) || popoverRef.current?.contains(target)) {
        return;
      }

      setIsOpen(false);
    }

    document.addEventListener("pointerdown", handlePointerDown);

    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
    };
  }, [closeOnOutsideClick, isOpen]);

  return (
    <div className="relative flex min-w-0 items-center gap-2">
      <button
        aria-controls={id}
        aria-expanded={isOpen}
        aria-label="Toggle attachment list"
        className={cx(
          "flex shrink-0 items-center gap-1.5",
          isCompact ? "h-7 px-2 text-xs" : "h-9 px-3 text-sm",
          getThemeButtonClassName(theme),
        )}
        onClick={() => setIsOpen((value) => !value)}
        ref={triggerRef}
        type="button"
      >
        <span className="tabular-nums">{attachments.length}</span>
        {isOpen ? <LuChevronDown size={isCompact ? 12 : 14} /> : <LuChevronUp size={isCompact ? 12 : 14} />}
      </button>

      {onUploadAttachments && attachments.length > 0 ? (
        <button
          aria-label="Upload chat attachments"
          className={cx("flex h-9 shrink-0 items-center justify-center px-3 text-sm", getThemeButtonClassName(theme))}
          onClick={onUploadAttachments}
          type="button"
        >
          <LuUpload size={14} />
        </button>
      ) : null}

      {isOpen && popoverPosition && typeof document !== "undefined"
        ? createPortal(
            <div
              className={cx("fixed z-[100] w-80 rounded-lg border p-3 shadow-xl", tone.border, tone.surface, tone.surfaceForeground)}
              id={id}
              ref={popoverRef}
              style={{ bottom: popoverPosition.bottom, left: popoverPosition.left }}
            >
              <div className="relative">
                <LuSearch className={cx("absolute left-3 top-1/2 -translate-y-1/2", tone.iconMuted)} size={14} />
                <input
                  className={cx("h-9 w-full rounded-md border px-3 pl-8 text-sm", tone.border, tone.surface, tone.foreground, tone.focusRing)}
                  onChange={(event) => setQuery(event.target.value)}
                  placeholder="搜索附件"
                  value={query}
                />
              </div>

              <div className="mt-3 max-h-56 space-y-2 overflow-y-auto pr-1">
                {filteredAttachments.length > 0 ? (
                  filteredAttachments.map((attachment) => (
                    <AttachmentInfoItem
                      attachment={attachment}
                      key={attachment.id}
                      onClosePopover={() => setIsOpen(false)}
                      {...(onRemoveAttachment ? { onRemove: onRemoveAttachment } : {})}
                      {...(renderAttachmentActions ? { renderActions: renderAttachmentActions } : {})}
                    />
                  ))
                ) : (
                  <div className={cx("rounded-md border px-3 py-6 text-center text-sm", tone.border, tone.subtleForeground)}>
                    没有匹配的附件
                  </div>
                )}
              </div>
            </div>,
            document.body,
          )
        : null}
    </div>
  );
}

type AiChatAttachmentToolbarProps = {
  attachments: ChatAttachment[];
  onAttachmentsSelected: (attachments: ChatAttachment[]) => void;
  onRemoveAttachment: (attachmentId: string) => void;
  onUploadAttachments: () => void;
  onValidationError: (message: string) => void;
};

export default function AiChatAttachmentToolbar({
  attachments,
  onAttachmentsSelected,
  onRemoveAttachment,
  onUploadAttachments,
  onValidationError,
}: AiChatAttachmentToolbarProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const [isToolMenuOpen, setIsToolMenuOpen] = useState(false);
  const [fileInputAccept, setFileInputAccept] = useState("");
  const selectedKindRef = useRef<ChatAttachment["kind"]>("file");
  const fileInputRef = useRef<HTMLInputElement>(null);

  function handleToolSelect(kind: ChatAttachment["kind"]) {
    const toolItem = attachmentToolItems.find((item) => item.kind === kind);
    if (!toolItem?.selectable) {
      return;
    }

    selectedKindRef.current = kind;
    setFileInputAccept(toolItem.accept);
    window.setTimeout(() => fileInputRef.current?.click(), 0);
  }

  function handleFileChange(event: ChangeEvent<HTMLInputElement>) {
    const selectedFiles = Array.from(event.target.files ?? []);
    const kind = selectedKindRef.current;
    const validAttachments: ChatAttachment[] = [];
    const invalidFileNames: string[] = [];

    for (const file of selectedFiles) {
      if (isFileMatchedKind(file, kind)) {
        validAttachments.push(createChatAttachment(file, kind));
      } else {
        invalidFileNames.push(file.name);
      }
    }

    if (validAttachments.length > 0) {
      onAttachmentsSelected(validAttachments);
    }

    if (invalidFileNames.length > 0) {
      const label = attachmentToolItems.find((item) => item.kind === kind)?.label ?? "附件";
      onValidationError(`${invalidFileNames.join("、")} 不符合${label}类型要求`);
    }

    event.target.value = "";
  }

  return (
    <div className="relative flex min-w-0 items-center gap-2">
      <input
        ref={fileInputRef}
        accept={fileInputAccept}
        className="hidden"
        multiple
        onChange={handleFileChange}
        type="file"
      />

      <button
        aria-label="Add chat attachment"
        className={cx("flex h-9 w-9 shrink-0 items-center justify-center", getThemeButtonClassName(theme))}
        onClick={() => {
          setIsToolMenuOpen((value) => !value);
        }}
        type="button"
      >
        <LuPlus size={16} />
      </button>

      {isToolMenuOpen ? (
        <div className={cx("absolute bottom-11 left-0 z-20 w-44 rounded-lg border p-1 shadow-xl", tone.border, tone.surface, tone.surfaceForeground)}>
          {attachmentToolItems.map((item) => {
            const Icon = item.icon;

            return (
              <button
                className={cx("flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm transition-colors", tone.secondaryHover)}
                key={item.kind}
                onClick={() => {
                  handleToolSelect(item.kind);
                  setIsToolMenuOpen(false);
                }}
                type="button"
              >
                <Icon className={tone.iconMuted} size={16} />
                <span>{item.label}</span>
              </button>
            );
          })}
        </div>
      ) : null}

      <AttachmentListPopover
        attachments={attachments}
        id="ai-chat-attachment-list"
        onRemoveAttachment={onRemoveAttachment}
        onUploadAttachments={onUploadAttachments}
      />
    </div>
  );
}
