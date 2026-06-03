export type ChatAttachment = {
  file?: File;
  id: string;
  kind: "audio" | "image" | "video" | "file";
  mimeType: string;
  name: string;
  previewText?: string;
  previewUrl?: string;
  sizeBytes: number;
  sizeLabel: string;
};

export type ChatMessage = {
  id: string;
  role: "assistant" | "user";
  content: string;
  attachments?: ChatAttachment[];
};

export type ChatSession = {
  id: string;
  title: string;
};

export type PopoverPosition = {
  left: number;
  top: number;
};

export type RequestMethod = "get" | "post";

export type OutputRenderMode = "instant" | "result";
