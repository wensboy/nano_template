import { useEffect, useState } from "react";
import "katex/dist/katex.min.css";
import { LuCheck, LuCopy } from "react-icons/lu";
import ReactMarkdown from "react-markdown";
import rehypeKatex from "rehype-katex";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { oneLight } from "react-syntax-highlighter/dist/esm/styles/prism";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";

import { useAppSelector } from "@/app/hooks";
import { cx, getThemeButtonClassName, getThemeTone } from "@/app/themeStyles";

type MarkdownContentProps = {
  content: string;
};

type MarkdownCodeBlockProps = {
  code: string;
  language: string;
};

function MarkdownCodeBlock({ code, language }: MarkdownCodeBlockProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const [isCopied, setIsCopied] = useState(false);
  const normalizedLanguage = language || "text";
  const shouldHighlight = normalizedLanguage !== "math";

  async function handleCopy() {
    await navigator.clipboard?.writeText(code);
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
    <div className={cx("my-2 max-w-full overflow-hidden rounded-md border", tone.border, tone.muted)}>
      <div className={cx("flex h-9 items-center justify-between border-b px-3", tone.border, tone.secondary)}>
        <span className={cx("truncate text-xs font-medium uppercase tracking-wide", tone.mutedForeground)}>
          {normalizedLanguage}
        </span>
        <div className="flex shrink-0 items-center gap-1">
          <button
            aria-label={isCopied ? "Code copied" : "Copy code"}
            className={cx("flex h-7 w-7 items-center justify-center", getThemeButtonClassName(theme))}
            onClick={() => void handleCopy()}
            type="button"
          >
            {isCopied ? <LuCheck size={13} /> : <LuCopy size={13} />}
          </button>
        </div>
      </div>

      <div className="max-w-full overflow-x-auto p-3 text-xs leading-5">
        {shouldHighlight ? (
          <SyntaxHighlighter
            PreTag="div"
            codeTagProps={{
              style: {
                background: "transparent",
              },
            }}
            customStyle={{
              background: "transparent",
              margin: 0,
              padding: 0,
            }}
            language={normalizedLanguage}
            wrapLongLines
            style={oneLight}
          >
            {code}
          </SyntaxHighlighter>
        ) : (
          <code className={cx("block whitespace-pre-wrap break-words", tone.foreground)}>{code}</code>
        )}
      </div>
    </div>
  );
}

export default function MarkdownContent({ content }: MarkdownContentProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return (
    <ReactMarkdown
      components={{
        a: ({ children, ...props }) => (
          <a
            className="break-words underline underline-offset-4"
            rel="noreferrer"
            target="_blank"
            {...props}
          >
            {children}
          </a>
        ),
        blockquote: ({ children, ...props }) => (
          <blockquote className={cx("border-l-2 pl-3", tone.border)} {...props}>
            {children}
          </blockquote>
        ),
        code: ({ children, className, node: _node, ...props }) => {
          const language = /language-(\w+)/.exec(className ?? "")?.[1];
          const codeContent = String(children).replace(/\n$/, "");
          const isBlockCode = codeContent.includes("\n");

          if (language && isBlockCode) {
            return <MarkdownCodeBlock code={codeContent} language={language} />;
          }

          return (
            <code
              className={cx(
                language === "math" ? "block whitespace-pre-wrap break-words" : "rounded px-1 py-0.5 text-[0.92em]",
                language === "math" ? tone.foreground : cx(tone.muted, tone.foreground),
              )}
              {...props}
            >
              {codeContent}
            </code>
          );
        },
        ol: ({ children, ...props }) => (
          <ol className="list-decimal space-y-1 pl-5" {...props}>
            {children}
          </ol>
        ),
        p: ({ children, ...props }) => (
          <p className="whitespace-pre-wrap break-words" {...props}>
            {children}
          </p>
        ),
        pre: ({ children }) => <>{children}</>,
        table: ({ children, ...props }) => (
          <div className="max-w-full overflow-x-auto">
            <table className={cx("w-full border-collapse text-sm", tone.border)} {...props}>
              {children}
            </table>
          </div>
        ),
        td: ({ children, ...props }) => (
          <td className={cx("border px-3 py-2 align-top", tone.border)} {...props}>
            {children}
          </td>
        ),
        th: ({ children, ...props }) => (
          <th className={cx("border px-3 py-2 text-left font-semibold", tone.border, tone.muted)} {...props}>
            {children}
          </th>
        ),
        ul: ({ children, ...props }) => (
          <ul className="list-disc space-y-1 pl-5" {...props}>
            {children}
          </ul>
        ),
      }}
      rehypePlugins={[rehypeKatex]}
      remarkPlugins={[remarkGfm, remarkMath]}
    >
      {content}
    </ReactMarkdown>
  );
}
