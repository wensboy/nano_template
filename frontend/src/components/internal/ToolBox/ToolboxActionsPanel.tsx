import { useAppDispatch, useAppSelector } from "@/app/hooks";
import { toggleThemeMode } from "@/app/store/themeSlice";
import {
  cx,
  getThemeButtonClassName,
  getThemePrimaryButtonClassName,
  getThemeTone,
} from "@/app/themeStyles";

type ToolboxActionsPanelProps = {
  onLogout: () => void;
};

export default function ToolboxActionsPanel({ onLogout }: ToolboxActionsPanelProps) {
  const dispatch = useAppDispatch();
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return (
    <div className="space-y-4">
      <p className={cx("text-sm leading-7", tone.mutedForeground)}>
        这里可以放当前页面常用的快捷操作，避免打断主视图布局。
      </p>
      <div className="flex flex-wrap gap-3">
        <button
          className={cx("h-9 px-4 text-sm font-medium", getThemePrimaryButtonClassName(theme))}
          onClick={() => dispatch(toggleThemeMode())}
          type="button"
        >
          切换主题
        </button>
        <button
          className={cx("h-9 px-4 text-sm font-medium", getThemeButtonClassName(theme))}
          onClick={onLogout}
          type="button"
        >
          退出登录
        </button>
      </div>
    </div>
  );
}
