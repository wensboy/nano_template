import type { ChatAttachment } from "@/components/internal/AiChatToolboxPanel/types";

export function getMessageHistoryTitle(content: string) {
  const normalizedContent = content.trim() || "空消息";

  if (normalizedContent.length <= 16) {
    return normalizedContent;
  }

  return `${normalizedContent.slice(0, 16)}....`;
}

export function formatAttachmentSize(bytes: number) {
  if (bytes < 1024) {
    return `${bytes} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function isFileMatchedKind(file: File, kind: ChatAttachment["kind"]) {
  if (kind === "image") {
    return file.type.startsWith("image/");
  }
  if (kind === "video") {
    return file.type.startsWith("video/");
  }
  if (kind === "audio") {
    return file.type.startsWith("audio/");
  }
  return true;
}

export function createChatAttachment(file: File, kind: ChatAttachment["kind"]): ChatAttachment {
  return {
    file,
    id: crypto.randomUUID(),
    kind,
    mimeType: file.type || "application/octet-stream",
    name: file.name,
    sizeBytes: file.size,
    sizeLabel: formatAttachmentSize(file.size),
  };
}
