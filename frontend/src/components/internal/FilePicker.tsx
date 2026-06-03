import { useCallback, useId, useRef, useState } from "react";
import {
  LuCheck,
  LuChevronDown,
  LuFile,
  LuPlus,
  LuUpload,
  LuX,
} from "react-icons/lu";

import { request } from "@/api/client";
import { PresignOssObject } from "@/api/aliyun/oss";
import { useAppSelector } from "@/app/hooks";
import { cx, getThemeTone } from "@/app/themeStyles";
import AlertPop from "@/components/internal/AlertPop";

// ── Types ────────────────────────────────────────────────────────────────────

export type FilePickerProps = {
  onFilesChange: (files: File[]) => void;
  accept?: string;
  maxFiles?: number;
  maxSize?: number;
};

type UploadStatus = {
  progress: number; // 0–100
  done: boolean;
};

// ── Helpers ──────────────────────────────────────────────────────────────────

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

/** Tiny SVG ring that fills clockwise as progress goes from 0 → 100. */
function ProgressRing({
  progress,
  progressClassName,
  size = 64,
  trackClassName,
}: {
  progress: number;
  progressClassName: string;
  size?: number;
  trackClassName: string;
}) {
  const strokeW = 3;
  const r = (size - strokeW) / 2;
  const circ = 2 * Math.PI * r;
  const offset = circ * (1 - Math.min(progress, 100) / 100);

  return (
    <svg
      aria-hidden
      className="absolute inset-0 pointer-events-none"
      viewBox={`0 0 ${size} ${size}`}
    >
      {/* track */}
      <circle
        cx={size / 2}
        cy={size / 2}
        fill="none"
        r={r}
        stroke="currentColor"
        strokeWidth={strokeW}
        className={trackClassName}
      />
      {/* fill arc – rotates from 12 o'clock clockwise */}
      <circle
        cx={size / 2}
        cy={size / 2}
        fill="none"
        r={r}
        stroke="currentColor"
        strokeLinecap="round"
        strokeWidth={strokeW}
        strokeDasharray={circ}
        strokeDashoffset={offset}
        transform={`rotate(-90 ${size / 2} ${size / 2})`}
        className={cx("transition-[stroke-dashoffset] duration-300 ease-linear", progressClassName)}
      />
    </svg>
  );
}

async function uploadSingleFile(
  file: File,
  onProgress: (pct: number) => void,
): Promise<void> {
  const res = await request(() =>
    PresignOssObject({
      object_key: file.name,
      mime: file.type || "application/octet-stream",
      size: file.size,
    }),
  );

  if (res.code !== 0 || !res.data) {
    throw new Error(res.message || "获取预签名地址失败");
  }

  const { signed_url, method, signed_headers } = res.data;

  await new Promise<void>((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open(method || "PUT", signed_url);

    for (const h of signed_headers ?? []) {
      xhr.setRequestHeader(h.key, h.value);
    }

    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) {
        onProgress(Math.round((e.loaded / e.total) * 100));
      }
    };

    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        resolve();
      } else {
        reject(new Error(`上传失败 HTTP ${xhr.status}`));
      }
    };
    xhr.onerror = () => reject(new Error("上传网络异常"));
    xhr.send(file);
  });
}

export default function FilePicker({
  onFilesChange,
  accept,
  maxFiles = 10,
  maxSize,
}: FilePickerProps) {
  const theme = useAppSelector((state) => state.theme);
  const tone = getThemeTone(theme);

  const [files, setFiles] = useState<File[]>([]);
  const [open, setOpen] = useState(false);
  const [uploadMessage, setUploadMessage] = useState("");
  const [uploading, setUploading] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const popoverId = useId();

  const uploadMapRef = useRef<Map<File, UploadStatus>>(new Map());
  const [, setTick] = useState(0);

  const handleAdd = useCallback(() => {
    inputRef.current?.click();
  }, []);

  const handleFileChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const incoming = Array.from(e.target.files ?? []);
      if (incoming.length === 0) return;

      setUploadMessage("");

      let valid = incoming;
      if (maxSize != null) {
        valid = valid.filter((f) => f.size <= maxSize);
      }

      const merged = [...files, ...valid].slice(0, maxFiles);
      setFiles(merged);
      onFilesChange(merged);

      if (inputRef.current) inputRef.current.value = "";
    },
    [files, maxFiles, maxSize, onFilesChange],
  );

  const handleRemove = useCallback(
    (index: number) => {
      setFiles((prev) => {
        const next = prev.filter((_, i) => i !== index);
        onFilesChange(next);
        return next;
      });
    },
    [onFilesChange],
  );

  const handleUpload = useCallback(async () => {
    if (uploading) return;
    const snapshot = [...files];
    if (snapshot.length === 0) return;

    // setOpen(false);
    setUploadMessage("");
    setUploading(true);
    const map = uploadMapRef.current;

    for (const file of snapshot) {
      if (!files.includes(file)) continue;

      map.set(file, { progress: 0, done: false });
      setTick((t) => t + 1);

      try {
        await uploadSingleFile(file, (pct) => {
          map.set(file, { progress: pct, done: false });
          setTick((t) => t + 1);
        });

        map.set(file, { progress: 100, done: true });
        setTick((t) => t + 1);
        await new Promise((r) => setTimeout(r, 700));

        setFiles((prev) => {
          const next = prev.filter((f) => f !== file);
          onFilesChange(next);
          return next;
        });
      } catch (err) {
        const reason = err instanceof Error ? err.message : "上传失败，请稍后重试";
        setUploadMessage(`${file.name}: ${reason}`);
      } finally {
        map.delete(file);
        setTick((t) => t + 1);
      }
    }

    setUploading(false);
  }, [files, uploading, onFilesChange]);

  return (
    <div className="relative inline-flex items-center">
      <AlertPop message={uploadMessage} variant="destructive" />

      {/* Hidden native file input */}
      <input
        ref={inputRef}
        accept={accept}
        className="hidden"
        multiple
        onChange={handleFileChange}
        type="file"
      />

      {/* ── collapsed bar ──────────────────────────────────────────────── */}
      <div
        className={cx(
          "flex items-center gap-1 rounded-md border px-1.5 py-1 text-sm font-medium shadow-sm transition-colors",
          tone.border,
          tone.surface,
          tone.surfaceForeground,
        )}
      >
        {/* add button */}
        <button
          aria-label="Add files"
          className={cx(
            "flex h-7 w-7 items-center justify-center rounded-md transition-colors",
            tone.focusRing,
            tone.secondaryHover,
          )}
          disabled={uploading}
          onClick={handleAdd}
          type="button"
        >
          <LuPlus size={16} />
        </button>

        {/* count + expand toggle */}
        <button
          aria-controls={popoverId}
          aria-expanded={open}
          aria-label="Toggle file list"
          className={cx(
            "flex h-7 items-center gap-1 rounded-md px-2 transition-colors",
            tone.focusRing,
            tone.secondaryHover,
          )}
          onClick={() => setOpen((v) => !v)}
          type="button"
        >
          <span className="tabular-nums">{files.length}</span>
          <LuChevronDown
            className={`transition-transform duration-200 ${open ? "rotate-180" : ""}`}
            size={14}
          />
        </button>

        {/* upload button – only visible when there are files */}
        {files.length > 0 && (
          <button
            aria-label="Upload all files"
            className={cx(
              "flex h-7 w-7 items-center justify-center rounded-md transition-colors",
              tone.focusRing,
              tone.secondaryHover,
              uploading && "pointer-events-none opacity-50",
            )}
            disabled={uploading}
            onClick={handleUpload}
            type="button"
          >
            <LuUpload size={14} />
          </button>
        )}
      </div>

      {/* ── popover file list ──────────────────────────────────────────── */}
      {open && files.length > 0 && (
        <div
          className={cx(
            "absolute left-0 top-full z-50 mt-2 w-64 rounded-lg border p-3 shadow-xl",
            tone.border,
            tone.surface,
            tone.surfaceForeground,
          )}
          id={popoverId}
        >
          <ul className="grid grid-cols-3 gap-2 justify-items-center">
            {files.map((file, index) => {
              const status = uploadMapRef.current.get(file);
              const isDone = status?.done;

              return (
                <li
                  className={cx(
                    "group relative flex h-16 w-16 flex-col items-center justify-center rounded-md border text-center transition-colors",
                    isDone ? tone.fileDone : cx(tone.muted, tone.border),
                  )}
                  key={`${file.name}-${file.size}-${index}`}
                  title={`${file.name}\n${formatSize(file.size)}`}
                >
                  {/* progress ring overlay */}
                  {status && !isDone && (
                    <ProgressRing
                      progress={status.progress}
                      progressClassName={tone.progress}
                      size={64}
                      trackClassName={tone.progressTrack}
                    />
                  )}

                  {/* done checkmark */}
                  {isDone ? (
                    <LuCheck className="shrink-0" size={26} />
                  ) : (
                    <>
                      <LuFile className={cx("shrink-0 text-lg", tone.iconMuted)} />
                      <span className={cx("mt-0.5 block w-full truncate px-1 text-[10px] leading-tight", tone.subtleForeground)}>
                        {file.name}
                      </span>
                    </>
                  )}

                  {/* remove button — hidden during upload for this file */}
                  {!status && (
                    <button
                      aria-label={`Remove ${file.name}`}
                      className={cx(
                        "absolute -right-1.5 -top-1.5 flex h-4 w-4 items-center justify-center rounded-full opacity-0 shadow-sm transition-colors group-hover:opacity-100",
                        tone.destructive,
                        tone.destructiveForeground,
                        tone.destructiveHover,
                      )}
                      onClick={(e) => {
                        e.stopPropagation();
                        handleRemove(index);
                      }}
                      type="button"
                    >
                      <LuX size={10} />
                    </button>
                  )}
                </li>
              );
            })}
          </ul>
        </div>
      )}
    </div>
  );
}
