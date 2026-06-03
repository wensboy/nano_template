import { forwardRef, type ReactNode, useMemo } from "react";
import {
  BlockColor,
  Blockquote,
  Bold,
  BulletList,
  CharacterCount,
  ClearFormatting,
  Code,
  CodeBlock,
  Dropcursor,
  FontFamily,
  FontSize,
  Gapcursor,
  HardBreak,
  Heading,
  Highlight,
  HorizontalRule,
  InvisibleChars,
  Italic,
  LineHeight,
  Link,
  LinkPopover,
  ListIndent,
  ListItem,
  ListKeymap,
  NotionColorPicker,
  OrderedList,
  Placeholder,
  Selection,
  SelectionDecoration,
  Strike,
  Subscript,
  Superscript,
  TaskItem,
  TaskList,
  TextAlign,
  TextColor,
  TextStyle,
  Typography,
  Underline,
  UniqueID,
  type AnyExtension,
  type ToolbarLayoutEntry,
} from "@domternal/core";
import {
  DomternalBubbleMenu,
  DomternalEditor,
  type DomternalEditorProps,
  type DomternalEditorRef,
  DomternalNotionColorPicker,
  DomternalToolbar,
} from "@domternal/react";
import "@domternal/theme/css";

import { useAppSelector } from "@/app/hooks";
import { cx } from "@/app/themeStyles";

export const richTextEditorToolbarLayout: ToolbarLayoutEntry[] = [
  "undo",
  "redo",
  "|",
  "heading",
  "fontFamily",
  "fontSize",
  "lineHeight",
  "|",
  "bold",
  "italic",
  "underline",
  "strike",
  "code",
  "subscript",
  "superscript",
  "|",
  "textColor",
  "highlight",
  "textAlign",
  "clearFormatting",
  "|",
  "bulletList",
  "orderedList",
  "taskList",
  "blockquote",
  "codeBlock",
  "horizontalRule",
  "|",
  "link",
  "invisibleChars",
];

export function createRichTextEditorExtensions(placeholder: string): AnyExtension[] {
  return [
    Heading.configure({ levels: [1, 2, 3, 4] }),
    Blockquote,
    CodeBlock,
    BulletList,
    OrderedList,
    ListItem,
    HorizontalRule,
    HardBreak,
    TaskList,
    TaskItem.configure({ nested: true }),
    Bold,
    Italic,
    Underline,
    Strike,
    Code,
    Link.configure({
      HTMLAttributes: {
        class: "underline underline-offset-2",
        rel: "noopener noreferrer",
        target: "_blank",
      },
      openOnClick: "whenNotEditable",
    }),
    Subscript,
    Superscript,
    Dropcursor,
    Gapcursor,
    ListKeymap,
    ListIndent,
    LinkPopover,
    SelectionDecoration,
    ...(placeholder ? [Placeholder.configure({ placeholder, showOnlyCurrent: false })] : []),
    CharacterCount,
    Typography,
    TextStyle,
    TextColor,
    Highlight,
    FontFamily.configure({
      fontFamilies: ["Inter", "Arial", "Verdana", "Georgia", "Times New Roman", "Courier New"],
    }),
    FontSize.configure({
      fontSizes: ["12px", "14px", "16px", "18px", "24px", "32px"],
      showReset: true,
    }),
    TextAlign.configure({
      types: ["heading", "paragraph"],
    }),
    LineHeight.configure({
      lineHeights: ["1", "1.25", "1.5", "1.75", "2"],
      types: ["heading", "paragraph"],
    }),
    UniqueID,
    BlockColor,
    Selection,
    InvisibleChars,
    NotionColorPicker,
    ClearFormatting,
  ];
}

export type RichTextEditorProps = Omit<DomternalEditorProps, "children" | "className" | "extensions"> & {
  className?: string;
  editorClassName?: string;
  extensions?: AnyExtension[];
  placeholder?: string;
  showBubbleMenu?: boolean;
  showNotionColorPicker?: boolean;
  showToolbar?: boolean;
  toolbarActions?: ReactNode;
  toolbarLayout?: ToolbarLayoutEntry[];
};

const RichTextEditor = forwardRef<DomternalEditorRef, RichTextEditorProps>(function RichTextEditor(
  {
    className,
    editorClassName,
    extensions: extensionOverrides,
    placeholder = "",
    showBubbleMenu = false,
    showNotionColorPicker = true,
    showToolbar = true,
    toolbarActions,
    toolbarLayout = richTextEditorToolbarLayout,
    ...editorProps
  },
  ref,
) {
  const themeMode = useAppSelector((state) => state.theme.mode);
  const extensions = useMemo(() => {
    const builtInExtensions = createRichTextEditorExtensions(placeholder);

    return extensionOverrides?.length ? [...builtInExtensions, ...extensionOverrides] : builtInExtensions;
  }, [extensionOverrides, placeholder]);

  return (
    <div className={cx("internal-rich-text-editor", `dm-theme-${themeMode}`, className)}>
      <DomternalEditor
        {...editorProps}
        ref={ref}
        className={cx("internal-rich-text-editor__surface", editorClassName)}
        extensions={extensions}
      >
        {showToolbar || toolbarActions ? (
          <div className="internal-rich-text-editor__toolbar-row" data-dm-editor-ui="">
            {showToolbar ? <DomternalToolbar layout={toolbarLayout} /> : null}
            {toolbarActions ? <div className="internal-rich-text-editor__toolbar-actions">{toolbarActions}</div> : null}
          </div>
        ) : null}
        {showBubbleMenu ? <DomternalBubbleMenu placement="top" /> : null}
        {showNotionColorPicker ? <DomternalNotionColorPicker /> : null}
      </DomternalEditor>
    </div>
  );
});

export default RichTextEditor;
export type { DomternalEditorRef, ToolbarLayoutEntry };
