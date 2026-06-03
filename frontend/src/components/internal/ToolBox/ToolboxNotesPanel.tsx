import { useAppSelector } from "@/app/hooks";
import { cx, getThemeTone } from "@/app/themeStyles";

export default function ToolboxNotesPanel() {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return (
    <div className="space-y-3">
      <p className={cx("text-sm font-medium", tone.foreground)}>接入说明</p>
      <ul className={cx("space-y-2 text-sm leading-7", tone.mutedForeground)}>
        <li>1. 用 `ToolboxProvider` 包住页面。</li>
        <li>2. 传入带 `icon / label / content` 的工具项数组。</li>
        <li>3. 页面右下角会自动获得 toolbox 入口。</li>
      </ul>
    </div>
  );
}
