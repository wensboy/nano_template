import { useEffect, useLayoutEffect, useRef, useState, type ChangeEvent } from "react";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeTone } from "@/app/themeStyles";
import cellsImage from "@/assets/cells.png";
import AiChatAttachmentPreview from "@/components/internal/AiChatToolboxPanel/AiChatAttachmentPreview";
import AiChatComposer from "@/components/internal/AiChatToolboxPanel/AiChatComposer";
import AiChatHeader from "@/components/internal/AiChatToolboxPanel/AiChatHeader";
import {
  AssistantChatMessageBubble,
  UserChatMessageBubble,
} from "@/components/internal/AiChatToolboxPanel/AiChatMessageBubble";
import {
  AiChatMessageHistoryPopover,
  AiChatSessionHistoryPopover,
  AiChatSettingsPopover,
} from "@/components/internal/AiChatToolboxPanel/AiChatPopovers";
import {
  chatModels,
  initialChatMessages,
  initialChatSessions,
  initialSessionId,
  markdownDemoReply,
  popoverGap,
  popoverWidth,
  temperatureMaxOptions,
} from "@/components/internal/AiChatToolboxPanel/constants";
import type {
  ChatAttachment,
  ChatMessage,
  ChatSession,
  OutputRenderMode,
  PopoverPosition,
  RequestMethod,
} from "@/components/internal/AiChatToolboxPanel/types";
import { getMessageHistoryTitle } from "@/components/internal/AiChatToolboxPanel/utils";
import AlertPop from "@/components/internal/AlertPop";

export default function AiChatToolboxPanel() {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const [model, setModel] = useState(chatModels[0]);
  const [chatSessions, setChatSessions] = useState<ChatSession[]>(initialChatSessions);
  const [currentSessionId, setCurrentSessionId] = useState(initialSessionId);
  const [messages, setMessages] = useState<ChatMessage[]>(initialChatMessages);
  const [inputValue, setInputValue] = useState("");
  const [attachments, setAttachments] = useState<ChatAttachment[]>([]);
  const [attachmentAlert, setAttachmentAlert] = useState("");
  const [attachmentAlertKey, setAttachmentAlertKey] = useState(0);
  const [isGenerating, setIsGenerating] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const [isMessageHistoryOpen, setIsMessageHistoryOpen] = useState(false);
  const [isHistoryOpen, setIsHistoryOpen] = useState(false);
  const [messageHistoryQuery, setMessageHistoryQuery] = useState("");
  const [historyQuery, setHistoryQuery] = useState("");
  const [editingSessionId, setEditingSessionId] = useState<string | null>(null);
  const [editingSessionTitle, setEditingSessionTitle] = useState("");
  const [isStreamEnabled, setIsStreamEnabled] = useState(true);
  const [isThinkingEnabled, setIsThinkingEnabled] = useState(false);
  const [baseUrl, setBaseUrl] = useState("https://api.openai.com/v1");
  const [requestMethod, setRequestMethod] = useState<"get" | "post">("post");
  const [apiKey, setApiKey] = useState("");
  const [temperature, setTemperature] = useState(0.7);
  const [temperatureMax, setTemperatureMax] = useState<(typeof temperatureMaxOptions)[number]>(1);
  const [maxTokens, setMaxTokens] = useState(1200);
  const [typewriterSpeedMs, setTypewriterSpeedMs] = useState(18);
  const [isOutputRenderingEnabled, setIsOutputRenderingEnabled] = useState(true);
  const [outputRenderMode, setOutputRenderMode] = useState<OutputRenderMode>("instant");
  const [settingsPopoverPosition, setSettingsPopoverPosition] = useState<PopoverPosition | null>(null);
  const [messageHistoryPopoverPosition, setMessageHistoryPopoverPosition] = useState<PopoverPosition | null>(null);
  const [historyPopoverPosition, setHistoryPopoverPosition] = useState<PopoverPosition | null>(null);
  const [previewAttachment, setPreviewAttachment] = useState<ChatAttachment | null>(null);
  const [isPreviewClosing, setIsPreviewClosing] = useState(false);
  const settingsButtonRef = useRef<HTMLButtonElement>(null);
  const messageHistoryButtonRef = useRef<HTMLButtonElement>(null);
  const historyButtonRef = useRef<HTMLButtonElement>(null);
  const messageItemRefs = useRef<Record<string, HTMLDivElement | null>>({});
  const messageListRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const responseTimerRef = useRef<number | null>(null);

  const sessionName = chatSessions.find((session) => session.id === currentSessionId)?.title ?? "";
  const filteredSessions = chatSessions.filter((session) =>
    session.title.toLowerCase().includes(historyQuery.trim().toLowerCase()),
  );
  const userMessageHistoryItems = messages
    .filter((message) => message.role === "user")
    .map((message) => ({
      id: message.id,
      title: getMessageHistoryTitle(message.content),
    }));
  const filteredMessageHistoryItems = userMessageHistoryItems.filter((item) =>
    item.title.toLowerCase().includes(messageHistoryQuery.trim().toLowerCase()),
  );

  useEffect(() => {
    const list = messageListRef.current;
    if (list) {
      list.scrollTop = list.scrollHeight;
    }
  }, [messages, isGenerating]);

  useEffect(() => {
    const textarea = textareaRef.current;
    if (!textarea) {
      return;
    }

    textarea.style.height = "0px";
    textarea.style.height = `${Math.min(textarea.scrollHeight, 128)}px`;
  }, [inputValue]);

  useEffect(() => {
    return () => {
      if (responseTimerRef.current) {
        window.clearTimeout(responseTimerRef.current);
      }
    };
  }, []);

  function getPopoverPosition(anchor: HTMLButtonElement, align: "end" | "start"): PopoverPosition {
    const rect = anchor.getBoundingClientRect();
    const viewportPadding = 16;
    const preferredLeft = align === "end" ? rect.right - popoverWidth : rect.left;
    const maxLeft = window.innerWidth - popoverWidth - viewportPadding;

    return {
      left: Math.min(Math.max(preferredLeft, viewportPadding), Math.max(maxLeft, viewportPadding)),
      top: rect.bottom + popoverGap,
    };
  }

  function updatePopoverPositions() {
    if (isSettingsOpen && settingsButtonRef.current) {
      setSettingsPopoverPosition(getPopoverPosition(settingsButtonRef.current, "start"));
    }
    if (isMessageHistoryOpen && messageHistoryButtonRef.current) {
      setMessageHistoryPopoverPosition(getPopoverPosition(messageHistoryButtonRef.current, "end"));
    }
    if (isHistoryOpen && historyButtonRef.current) {
      setHistoryPopoverPosition(getPopoverPosition(historyButtonRef.current, "end"));
    }
  }

  useLayoutEffect(() => {
    updatePopoverPositions();
  }, [isHistoryOpen, isMessageHistoryOpen, isSettingsOpen]);

  useEffect(() => {
    if (!isSettingsOpen && !isMessageHistoryOpen && !isHistoryOpen) {
      return;
    }

    window.addEventListener("resize", updatePopoverPositions);
    window.addEventListener("scroll", updatePopoverPositions, true);

    return () => {
      window.removeEventListener("resize", updatePopoverPositions);
      window.removeEventListener("scroll", updatePopoverPositions, true);
    };
  }, [isHistoryOpen, isMessageHistoryOpen, isSettingsOpen]);

  function handleAttachmentsSelected(nextAttachments: ChatAttachment[]) {
    setAttachments((current) => [...current, ...nextAttachments]);
  }

  function showAttachmentAlert(message: string) {
    setAttachmentAlert(message);
    setAttachmentAlertKey((current) => current + 1);
  }

  function handleRemoveAttachment(attachmentId: string) {
    setAttachments((current) => current.filter((attachment) => attachment.id !== attachmentId));
  }

  function handleUploadAttachments() {
    // Reserved for unified attachment upload.
  }

  function handleDownloadAttachment() {
    // Reserved for generated attachment download.
  }

  function handlePreviewAttachment(attachment: ChatAttachment) {
    setPreviewAttachment(attachment);
    setIsPreviewClosing(false);
  }

  function handleClosePreview() {
    setIsPreviewClosing(true);
  }

  function handlePreviewExitComplete() {
    setPreviewAttachment(null);
    setIsPreviewClosing(false);
  }

  function handleClearComposer() {
    setInputValue("");
    setAttachments([]);
  }

  function handleSend() {
    const trimmed = inputValue.trim();
    if ((!trimmed && attachments.length === 0) || isGenerating) {
      return;
    }

    setMessages((current) => [
      ...current,
      {
        id: crypto.randomUUID(),
        role: "user",
        content: trimmed || "已发送附件。",
        attachments,
      },
    ]);
    setInputValue("");
    setAttachments([]);
    setIsGenerating(true);

    responseTimerRef.current = window.setTimeout(() => {
      setMessages((current) => [
        ...current,
        {
          id: crypto.randomUUID(),
          role: "assistant",
          content: markdownDemoReply,
          attachments: [
            {
              id: crypto.randomUUID(),
              kind: "image",
              mimeType: "image/png",
              name: "cells.png",
              previewUrl: cellsImage,
              sizeBytes: 865410,
              sizeLabel: "845.1 KB",
            },
          ],
        },
      ]);
      setIsGenerating(false);
      responseTimerRef.current = null;
    }, 900);
  }

  function handleStop() {
    if (responseTimerRef.current) {
      window.clearTimeout(responseTimerRef.current);
      responseTimerRef.current = null;
    }
    setIsGenerating(false);
  }

  function handleInputChange(event: ChangeEvent<HTMLTextAreaElement>) {
    setInputValue(event.target.value);
  }

  function handleStartEditSession(session: ChatSession) {
    setEditingSessionId(session.id);
    setEditingSessionTitle(session.title);
  }

  function handleCancelEditSession() {
    setEditingSessionId(null);
    setEditingSessionTitle("");
  }

  function handleCommitSessionTitle() {
    const nextTitle = editingSessionTitle.trim();

    if (!editingSessionId) {
      return;
    }

    if (nextTitle) {
      setChatSessions((current) =>
        current.map((session) =>
          session.id === editingSessionId ? { ...session, title: nextTitle } : session,
        ),
      );
    }

    handleCancelEditSession();
  }

  function handleDeleteSession(sessionId: string) {
    if (chatSessions.length <= 1) {
      return;
    }

    const nextSessions = chatSessions.filter((session) => session.id !== sessionId);
    setChatSessions(nextSessions);

    if (sessionId === currentSessionId) {
      setCurrentSessionId(nextSessions[0]?.id ?? initialSessionId);
    }

    if (editingSessionId === sessionId) {
      handleCancelEditSession();
    }
  }

  function handleSelectMessageHistory(messageId: string) {
    messageItemRefs.current[messageId]?.scrollIntoView({ behavior: "smooth", block: "center" });
    setIsMessageHistoryOpen(false);
  }

  return (
    <div className={cx("flex h-full min-h-0 flex-col overflow-hidden rounded-lg border", tone.border, tone.surface)}>
      <AlertPop key={attachmentAlertKey} message={attachmentAlert} variant="warning" />

      <AiChatHeader
        historyButtonRef={historyButtonRef}
        isStreamEnabled={isStreamEnabled}
        isThinkingEnabled={isThinkingEnabled}
        messageHistoryButtonRef={messageHistoryButtonRef}
        model={model}
        onToggleHistory={() => {
          setIsSettingsOpen(false);
          setIsMessageHistoryOpen(false);
          setIsHistoryOpen((value) => !value);
        }}
        onToggleMessageHistory={() => {
          setIsSettingsOpen(false);
          setIsHistoryOpen(false);
          setIsMessageHistoryOpen((value) => !value);
        }}
        onToggleSettings={() => {
          setIsMessageHistoryOpen(false);
          setIsHistoryOpen(false);
          setIsSettingsOpen((value) => !value);
        }}
        sessionName={sessionName}
        settingsButtonRef={settingsButtonRef}
      />

      <AiChatSettingsPopover
        apiKey={apiKey}
        baseUrl={baseUrl}
        isOpen={isSettingsOpen}
        isOutputRenderingEnabled={isOutputRenderingEnabled}
        isStreamEnabled={isStreamEnabled}
        isThinkingEnabled={isThinkingEnabled}
        maxTokens={maxTokens}
        model={model}
        onApiKeyChange={setApiKey}
        onBaseUrlChange={setBaseUrl}
        onMaxTokensChange={setMaxTokens}
        onModelChange={setModel}
        onOutputRenderModeChange={setOutputRenderMode}
        onOutputRenderingEnabledChange={setIsOutputRenderingEnabled}
        onRequestMethodChange={setRequestMethod}
        onStreamEnabledChange={setIsStreamEnabled}
        onTemperatureChange={setTemperature}
        onTemperatureMaxChange={setTemperatureMax}
        onThinkingEnabledChange={setIsThinkingEnabled}
        onTypewriterSpeedChange={setTypewriterSpeedMs}
        outputRenderMode={outputRenderMode}
        position={settingsPopoverPosition}
        requestMethod={requestMethod}
        temperature={temperature}
        temperatureMax={temperatureMax}
        typewriterSpeedMs={typewriterSpeedMs}
      />

      <AiChatMessageHistoryPopover
        isOpen={isMessageHistoryOpen}
        items={filteredMessageHistoryItems}
        onQueryChange={setMessageHistoryQuery}
        onSelectMessage={handleSelectMessageHistory}
        position={messageHistoryPopoverPosition}
        query={messageHistoryQuery}
      />

      <AiChatSessionHistoryPopover
        chatSessionsCount={chatSessions.length}
        currentSessionId={currentSessionId}
        editingSessionId={editingSessionId}
        editingSessionTitle={editingSessionTitle}
        isOpen={isHistoryOpen}
        onCancelEditSession={handleCancelEditSession}
        onCommitSessionTitle={handleCommitSessionTitle}
        onDeleteSession={handleDeleteSession}
        onEditingSessionTitleChange={setEditingSessionTitle}
        onQueryChange={setHistoryQuery}
        onSelectSession={(sessionId) => {
          setCurrentSessionId(sessionId);
          setIsHistoryOpen(false);
        }}
        onStartEditSession={handleStartEditSession}
        position={historyPopoverPosition}
        query={historyQuery}
        sessions={filteredSessions}
      />

      <div className="relative min-h-0 flex-1 overflow-hidden">
        <div
          aria-hidden={previewAttachment ? true : undefined}
          ref={messageListRef}
          className={cx(
            "h-full space-y-6 overflow-y-auto px-3 py-4",
            previewAttachment ? "pointer-events-none" : undefined,
          )}
        >
          {messages.map((message) =>
            message.role === "user" ? (
              <div
                key={message.id}
                ref={(node) => {
                  messageItemRefs.current[message.id] = node;
                }}
              >
                <UserChatMessageBubble message={message} />
              </div>
            ) : (
              <div
                key={message.id}
                ref={(node) => {
                  messageItemRefs.current[message.id] = node;
                }}
              >
                <AssistantChatMessageBubble
                  isOutputRenderingEnabled={isOutputRenderingEnabled}
                  isStreamEnabled={isStreamEnabled}
                  message={message}
                  onDownloadAttachment={handleDownloadAttachment}
                  onPreviewAttachment={handlePreviewAttachment}
                  outputRenderMode={outputRenderMode}
                  typewriterSpeedMs={typewriterSpeedMs}
                />
              </div>
            ),
          )}

          {isGenerating ? (
            <div className="flex justify-start">
              <div className={cx("max-w-[88%] rounded-lg border px-4 py-3 text-sm", tone.border, tone.subtle, tone.mutedForeground)}>
                正在生成回复...
              </div>
            </div>
          ) : null}
        </div>

        {previewAttachment ? (
          <AiChatAttachmentPreview
            attachment={previewAttachment}
            isClosing={isPreviewClosing}
            onBack={handleClosePreview}
            onExitComplete={handlePreviewExitComplete}
          />
        ) : null}
      </div>

      <AiChatComposer
        attachments={attachments}
        inputValue={inputValue}
        isGenerating={isGenerating}
        onAttachmentsSelected={handleAttachmentsSelected}
        onClearComposer={handleClearComposer}
        onInputChange={handleInputChange}
        onRemoveAttachment={handleRemoveAttachment}
        onSend={handleSend}
        onStop={handleStop}
        onUploadAttachments={handleUploadAttachments}
        onValidationError={showAttachmentAlert}
        textareaRef={textareaRef}
      />
    </div>
  );
}

export { default as AiChatAttachmentPreview } from "@/components/internal/AiChatToolboxPanel/AiChatAttachmentPreview";
export { default as AiChatAttachmentToolbar, AttachmentListPopover } from "@/components/internal/AiChatToolboxPanel/AiChatAttachmentToolbar";
export { default as AiChatComposer } from "@/components/internal/AiChatToolboxPanel/AiChatComposer";
export { default as AiChatHeader } from "@/components/internal/AiChatToolboxPanel/AiChatHeader";
export { AssistantChatMessageBubble, UserChatMessageBubble } from "@/components/internal/AiChatToolboxPanel/AiChatMessageBubble";
export {
  AiChatMessageHistoryPopover,
  AiChatSessionHistoryPopover,
  AiChatSettingsPopover,
} from "@/components/internal/AiChatToolboxPanel/AiChatPopovers";
export { default as MarkdownContent } from "@/components/internal/AiChatToolboxPanel/MarkdownContent";
export * from "@/components/internal/AiChatToolboxPanel/constants";
export * from "@/components/internal/AiChatToolboxPanel/types";
export * from "@/components/internal/AiChatToolboxPanel/utils";
