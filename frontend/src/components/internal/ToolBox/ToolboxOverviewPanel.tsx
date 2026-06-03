import { useAppSelector } from "@/app/hooks";
import { cx, getThemeTone } from "@/app/themeStyles";

export default function ToolboxOverviewPanel() {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);
  const isDark = theme.mode === "dark";

  return (
    <div className="space-y-4">
      <p className={cx("text-sm leading-7", tone.mutedForeground)}>
        这个 toolbox 通过 provider
        注入到页面顶层，右下角按钮会始终悬浮显示，适合挂载页面级的工具、
        说明和快捷入口。
      </p>
      <div className="grid gap-3 sm:grid-cols-2">
        <div className={cx("rounded-lg border p-4", tone.border, tone.subtle)}>
          <p className={cx("text-xs font-medium uppercase", tone.subtleForeground)}>
            Theme
          </p>
          <p className={cx("mt-2 text-lg font-semibold", tone.foreground)}>
            {isDark ? "Dark" : "Light"}
          </p>
        </div>
        <div className={cx("rounded-lg border p-4", tone.border, tone.subtle)}>
          <p className={cx("text-xs font-medium uppercase", tone.subtleForeground)}>
            Auth
          </p>
          <p className={cx("mt-2 text-lg font-semibold", tone.foreground)}>
            Protected Route Ready
          </p>
        </div>
      </div>
    </div>
  );
}
