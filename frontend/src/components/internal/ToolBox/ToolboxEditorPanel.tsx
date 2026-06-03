import { useEffect, useRef, useState } from "react";
import { LuCheck, LuClipboard, LuCloud, LuHardDrive, LuRefreshCcw, LuSave } from "react-icons/lu";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeButtonClassName, getThemeTone } from "@/app/themeStyles";
import RichTextEditor, { type DomternalEditorRef } from "@/components/internal/RichTextEditor";

const initialEditorContent = `
<h2>项目草稿</h2>
<p>今天先整理登录态、文件上传和工具箱入口的联动，保留可直接迁移到业务表单里的内容结构。</p>
<ul>
  <li><p>确认页面级 action 的展示顺序。</p></li>
  <li><p>记录富文本输出的 HTML 片段。</p></li>
</ul>
<blockquote><p>后续可以把当前草稿保存到后端配置或用户笔记。</p></blockquote>
`;

function getTextLengthFromHtml(html: string) {
  if (typeof DOMParser === "undefined") {
    return html.replace(/<[^>]*>/g, "").length;
  }

  return new DOMParser().parseFromString(html, "text/html").body.textContent?.length ?? 0;
}

export default function ToolboxEditorPanel() {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const editorRef = useRef<DomternalEditorRef>(null);
  const saveMenuRef = useRef<HTMLDivElement>(null);
  const [content, setContent] = useState(initialEditorContent);
  const [copied, setCopied] = useState(false);
  const [saveMenuOpen, setSaveMenuOpen] = useState(false);
  const textLength = getTextLengthFromHtml(content);

  useEffect(() => {
    if (!saveMenuOpen) {
      return;
    }

    function handlePointerDown(event: PointerEvent) {
      if (event.target instanceof Node && saveMenuRef.current?.contains(event.target)) {
        return;
      }

      setSaveMenuOpen(false);
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setSaveMenuOpen(false);
      }
    }

    document.addEventListener("pointerdown", handlePointerDown);
    window.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [saveMenuOpen]);

  async function handleCopy() {
    await navigator.clipboard?.writeText(editorRef.current?.htmlContent ?? content);
    setCopied(true);
    window.setTimeout(() => setCopied(false), 1200);
  }

  return (
    <div className="flex h-full min-h-0 flex-col">
      <RichTextEditor
        ref={editorRef}
        className="internal-rich-text-editor--contained min-h-0 flex-1"
        outputFormat="html"
        value={content}
        onChange={(value) => setContent(String(value))}
      />

      <div
        className={cx(
          "flex shrink-0 items-center justify-between border border-t-0 px-3 py-2 text-xs",
          tone.border,
          tone.subtle,
          tone.subtleForeground,
        )}
      >
        <span>字符总数: {textLength}</span>

        <div className="relative flex flex-row-reverse items-center gap-2" ref={saveMenuRef}>
          <button
            aria-expanded={saveMenuOpen}
            aria-haspopup="menu"
            aria-label="保存"
            className={cx("flex h-8 w-8 items-center justify-center", getThemeButtonClassName(theme))}
            onClick={() => setSaveMenuOpen((open) => !open)}
            title="保存"
            type="button"
          >
            <LuSave size={16} />
          </button>
          <button
            aria-label="复制 HTML"
            className={cx("flex h-8 w-8 items-center justify-center", getThemeButtonClassName(theme))}
            onClick={handleCopy}
            title="复制 HTML"
            type="button"
          >
            {copied ? <LuCheck size={16} /> : <LuClipboard size={16} />}
          </button>
          <button
            aria-label="重置"
            className={cx("flex h-8 w-8 items-center justify-center", getThemeButtonClassName(theme))}
            onClick={() => setContent(initialEditorContent)}
            title="重置"
            type="button"
          >
            <LuRefreshCcw size={16} />
          </button>

          {saveMenuOpen ? (
            <div
              className={cx(
                "absolute bottom-full right-0 z-50 mb-2 min-w-36 rounded-md border p-1 shadow-lg",
                tone.border,
                tone.surface,
                tone.surfaceForeground,
              )}
              role="menu"
            >
              <button
                className={cx(
                  "flex w-full items-center gap-2 rounded-sm px-2 py-2 text-left text-sm transition-colors",
                  tone.secondaryHover,
                )}
                onClick={() => setSaveMenuOpen(false)}
                role="menuitem"
                type="button"
              >
                <LuHardDrive size={16} />
                本地
              </button>
              <button
                className={cx(
                  "flex w-full items-center gap-2 rounded-sm px-2 py-2 text-left text-sm transition-colors",
                  tone.secondaryHover,
                )}
                onClick={() => setSaveMenuOpen(false)}
                role="menuitem"
                type="button"
              >
                <LuCloud size={16} />
                云
              </button>
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );
}
