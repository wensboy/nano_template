import { useAppSelector } from "@/app/hooks";
import { cx, getThemeTone } from "@/app/themeStyles";
import FilePicker from "@/components/internal/FilePicker";

export default function ToolboxFilesPanel() {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  return (
    <div className="space-y-4">
      <p className={cx("text-sm leading-7", tone.mutedForeground)}>
        选择需要上传的文件，支持批量添加与单个移除。
      </p>
      <FilePicker
        accept="image/*"
        maxFiles={8}
        onFilesChange={(files) => {
          // eslint-disable-next-line no-console
          console.log(
            "FilePicker selection:",
            files.map((f) => ({
              name: f.name,
              size: f.size,
              type: f.type,
            })),
          );
        }}
      />
      <p className={cx("text-xs", tone.subtleForeground)}>
        当前仅记录文件信息到控制台，后续可接入 OSS presign 上传流程。
      </p>
    </div>
  );
}
