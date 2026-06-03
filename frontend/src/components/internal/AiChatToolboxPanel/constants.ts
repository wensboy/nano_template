import cellsImage from "@/assets/cells.png";
import type { ChatMessage, ChatSession } from "@/components/internal/AiChatToolboxPanel/types";

export const chatModels = ["gpt-4.1", "gpt-4.1-mini", "o4-mini", "claued-code-opus-4.7-plus-max"];
export const temperatureMaxOptions = [1, 2] as const;
export const popoverGap = 8;
export const popoverWidth = 288;
export const initialSessionId = "session-current";

export const initialChatSessions: ChatSession[] = [
  { id: initialSessionId, title: "Nano Template AI Chat Session" },
  { id: "session-integration", title: "Nano Template 集成讨论" },
  { id: "session-auth", title: "权限校验和 Cookie 会话" },
  { id: "session-oss", title: "OSS 文件上传体验优化" },
  { id: "session-toolbox", title: "Toolbox 组件设计" },
  { id: "session-theme", title: "Shadcn 主题调色板" },
  { id: "session-role", title: "Role Select 表单改造" },
];

export const markdownDemoReply = `# AI Chat Markdown 演示

这是一段用于验证 **Markdown 渲染组件** 的默认回复。它应该能够在父容器宽度内稳定显示，并处理常见结构。

## 能力清单

- 支持普通段落、**加粗**、链接和列表
- 支持引用块
- 支持行内代码，例如 \`pnpm run build\`
- 支持代码块，并在内容较长时保持横向滚动
- 支持 GFM 表格、任务列表和删除线
- 支持 math 语法解析

> 这是一段引用内容，用于检查边框、缩进和主题色调是否协调。

1. 保持 AI 气泡主体可读
2. 不影响用户气泡的纯文本显示
3. 不引入额外行为副作用

## GFM 展示

| 类型 | 语法 | 状态 |
| --- | --- | --- |
| 表格 | \`remark-gfm\` | 已启用 |
| 删除线 | ~~过期方案~~ | 已启用 |
| 任务列表 | \`- [x]\` | 已启用 |

- [x] 渲染 markdown 段落
- [x] 渲染 GFM 扩展语法
- [ ] 接入真实 AI 回复流

## Math 展示

行内公式: $E = mc^2$

块级公式:

$$
\\int_0^1 x^2 dx = \\frac{1}{3}
$$

\`\`\`tsx
function MarkdownPreview() {
  return <MarkdownContent content={markdown} />;
}
\`\`\`

[React Markdown](https://github.com/remarkjs/react-markdown) 当前用于完成 Markdown 渲染，插件由 \`remark-gfm\` 和 \`remark-math\` 扩展。`;

export const initialChatMessages: ChatMessage[] = [
  {
    id: "m1",
    role: "assistant",
    content: markdownDemoReply,
  },
  {
    id: "m2",
    role: "user",
    content: "帮我看一下这个模板如何接入一个简洁的聊天面板。",
    attachments: [
      {
        id: "a1",
        kind: "image",
        mimeType: "image/png",
        name: "layout-reference.png",
        sizeBytes: 421888,
        sizeLabel: "412 KB",
      },
      {
        id: "a2",
        kind: "file",
        mimeType: "text/markdown",
        name: "requirements.md",
        sizeBytes: 18432,
        sizeLabel: "18 KB",
      },
    ],
  },
  {
    id: "m3",
    role: "assistant",
    content: "建议保持头部、消息列表和输入区三段式结构。头部用于模型与历史，会话主体独立滚动，底部工具栏固定在输入框下方。",
    attachments: [
      {
        id: "a3",
        kind: "image",
        mimeType: "image/png",
        name: "cells.png",
        previewUrl: cellsImage,
        sizeBytes: 865410,
        sizeLabel: "845.1 KB",
      },
      {
        id: "a4",
        kind: "file",
        mimeType: "text/markdown",
        name: "ai-summary.md",
        previewText: "## AI 生成文件\n\n这是用于验证 AI 气泡附件预览的示例文件内容。为了很好地演示滚动效果, 接下来需要接入一些完整的mock内容:\n Nim is a statically typed compiled systems programming language. It combines successful concepts from mature languages like Python, Ada and Modula. ",
        sizeBytes: 96,
        sizeLabel: "96 B",
      },
    ],
  },
];
