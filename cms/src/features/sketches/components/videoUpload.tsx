import { useState, useCallback, useRef, useEffect } from "react";

import { useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api-client";
import { useSketchVideos } from "../api/getSketchVideos";
import { SketchVideo, PipelineJob } from "@/types/api";

import { UploadIcon, VideoIcon, XIcon } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";

import { cn, formatBytes } from "@/lib/utils";

interface VideosResponse {
  videos: SketchVideo[];
}

interface UploadUrlResponse {
  uploadUrl: string;
  s3Key: string;
}

const POLL_INTERVAL = 4000;

const STATUS_CONFIG = {
  pending: {
    label: "Pipeline Job Pending",
    color: "bg-yellow-500/15 text-yellow-800 border-yellow-500/30",
    dot: "bg-yellow-400",
    pulse: true,
  },
  processing: {
    label: "Processing",
    color: "bg-orange-500/15 text-orange-800 border-orange-500/30",
    dot: "bg-orange-400",
    pulse: true,
  },
  completed: {
    label: "Pipeline Complete",
    color: "bg-emerald-500/15 text-emerald-800 border-emerald-500/30",
    dot: "bg-emerald-400",
    pulse: false,
  },
  error: {
    label: "Pipeline Error",
    color: "bg-red-500/15 text-red-800 border-red-500/30",
    dot: "bg-red-400",
    pulse: false,
  },
};

async function getUploadUrl(
  fileName: string,
  contentType: string,
  fileSize: number,
): Promise<UploadUrlResponse> {
  const res = await api.post<UploadUrlResponse>("/admin/sketch/upload-url", {
    fileName,
    contentType,
    fileSize,
  });
  return res;
}

async function notifyUploaded(
  sketchId: number,
  s3Key: string,
): Promise<VideosResponse> {
  const res = await api.post<VideosResponse>(
    `/admin/sketch/${sketchId}/video-uploaded`,
    { s3Key },
  );
  return res;
}

// ─── Sub-components ───────────────────────────────────────────────────────────

function StatusBadge({ status }: { status: PipelineJob["status"] }) {
  const cfg = STATUS_CONFIG[status] ?? STATUS_CONFIG.error;
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border",
        cfg.color,
      )}
    >
      <span
        className={cn("w-1.5 h-1.5 rounded-full shrink-0", cfg.dot, {
          "animate-pulse": cfg.pulse,
        })}
      />
      {cfg.label}
    </span>
  );
}

function VideoCard({ video }: { video: SketchVideo }) {
  const latestJob = video.jobs[video.jobs.length - 1];

  return (
    <div className="rounded-xl border border-gray-200 bg-gray-100 text-sm overflow-hidden">
      <div className="flex items-center justify-between px-4 py-3 border-b border-gray-200 bg-gray-100">
        <div className="flex items-start gap-2.5">
          <div className="w-7 h-7 rounded-lg bg-gray-200 flex items-center justify-center">
            <VideoIcon className="w-3.5 h-3.5 text-gray-800" />
          </div>
          <div className="flex flex-col gap-2">
            {video.coldS3Key ? (
              <div>Cold Storage: {video.coldS3Key}</div>
            ) : video.hotS3Key ? (
              <div>Hot Storage: {video.hotS3Key}</div>
            ) : (
              ""
            )}
            {video.archivedAt && (
              <div>
                Archived At: {new Date(video.archivedAt).toLocaleString()}
              </div>
            )}
            {latestJob?.error && (
              <div className="mt-2 p-2 rounded-lg border-red-400 border bg-red-300 text-red-600">
                {latestJob.error}
              </div>
            )}
          </div>
        </div>
        <div className="self-start">
          {latestJob && <StatusBadge status={latestJob.status} />}
        </div>
      </div>
    </div>
  );
}

function UploadZone({
  onSubmit,
  uploading,
  progress,
  error,
}: {
  onSubmit: (file: File) => void;
  uploading: boolean;
  progress: number | null;
  error: string | null;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [dragging, setDragging] = useState(false);
  const onSubmitRef = useRef(onSubmit);
  const [stagedFile, setStagedFile] = useState<File | null>(null);

  useEffect(() => {
    onSubmitRef.current = onSubmit;
  });

  const stageFile = (file: File) => {
    if (file.type.startsWith("video/")) setStagedFile(file);
  };

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setDragging(false);
    const file = e.dataTransfer.files[0];
    if (file) stageFile(file);
  }, []);

  const handleClear = () => {
    setStagedFile(null);
    if (inputRef.current) inputRef.current.value = "";
  };
  // Staged file preview
  if (stagedFile && !uploading) {
    return (
      <div className="space-y-3">
        <div className="rounded-xl border border-gray-200 bg-gray-100 overflow-hidden">
          {/* File info */}
          <div className="flex items-center gap-3 px-4 py-3 border-b border-gray-200">
            <div className="w-9 h-9 rounded-lg bg-gray-100 flex items-center justify-center shrink-0">
              <VideoIcon className="w-4 h-4 text-gray-800" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-sm font-medium text-gray-800 truncate">
                {stagedFile.name}
              </p>
              <p className="text-xs text-gray-800 mt-0.5">
                {formatBytes(stagedFile.size)}
                <span className="mx-1.5 text-gray-800">·</span>
                <span className="font-mono">{stagedFile.type || "video"}</span>
              </p>
            </div>
            <button
              type="button"
              onClick={handleClear}
              className="w-6 h-6 rounded-md flex items-center justify-center text-gray-800 hover:text-gray-700 hover:bg-gray-200 transition-colors shrink-0"
              title="Remove file"
            >
              <XIcon className="w-3.5 h-3.5" />
            </button>
          </div>

          <div className="flex items-center justify-between px-4 py-3 gap-3">
            <button
              type="button"
              onClick={handleClear}
              className="text-xs text-gray-800 hover:text-gray-700 transition-colors"
            >
              Choose a different file
            </button>
            <button
              type="button"
              onClick={() => onSubmitRef.current(stagedFile)}
              className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-orange-500 hover:bg-orange-400 text-white text-xs font-medium transition-colors"
            >
              <UploadIcon className="w-3.5 h-3.5" />
              Upload Video
            </button>
          </div>
        </div>
      </div>
    );
  }

  // Uploading state
  if (uploading) {
    return (
      <div className="space-y-3">
        <div className="rounded-xl border border-gray-200 bg-gray-100 px-4 py-5 flex flex-col items-center gap-3">
          <Spinner className="w-5 h-5 text-orange-400" />
          <div className="text-center">
            <p className="text-sm font-medium text-gray-800">Uploading…</p>
            <p className="text-xs text-gray-600 mt-0.5">
              {progress !== null ? `${progress}%` : "Preparing…"}
            </p>
          </div>
          {progress !== null && (
            <div className="w-full bg-gray-100 rounded-full h-1 overflow-hidden">
              <div
                className="bg-orange-500 h-full rounded-full transition-all duration-300"
                style={{ width: `${progress}%` }}
              />
            </div>
          )}
        </div>
        {error && (
          <div className="text-red-400 text-xs border-red-500/30 bg-red-500/10">
            {error}
          </div>
        )}
      </div>
    );
  }
  return (
    <div className="space-y-3">
      <button
        type="button"
        onClick={() => inputRef.current?.click()}
        onDragOver={(e) => {
          e.preventDefault();
          setDragging(true);
        }}
        onDragLeave={() => setDragging(false)}
        onDrop={handleDrop}
        disabled={uploading}
        className={cn(
          "w-full rounded-xl border-2 border-dashed transition-all duration-200 p-10 flex flex-col items-center gap-3 cursor-pointer group",
          dragging
            ? "border-orange-500 bg-orange-500/5"
            : "border-gray-400 bg-gray-50",
          uploading && "opacity-50 cursor-not-allowed",
        )}
      >
        <div
          className={cn(
            "w-12 h-12 rounded-xl flex items-center justify-center transition-colors",
            dragging
              ? "bg-orange-500/20"
              : "bg-gray-200 group-hover:bg-gray-100",
          )}
        >
          {uploading ? (
            <Spinner className="w-5 h-5 text-orange-400" />
          ) : (
            <UploadIcon
              className={cn(
                "w-5 h-5 transition-colors",
                dragging
                  ? "text-orange-400"
                  : "text-gray-400 group-hover:text-gray-300",
              )}
            />
          )}
        </div>
        <div className="text-center">
          <p className="text-sm font-medium text-gray-400">
            {uploading ? "Uploading…" : "Upload Video"}
          </p>
          <p className="text-xs text-gray-600 mt-0.5">
            {uploading
              ? progress !== null
                ? `${progress}%`
                : "Preparing…"
              : "Drag & drop or click to browse"}
          </p>
        </div>
        <input
          ref={inputRef}
          type="file"
          accept="video/*"
          className="hidden"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) stageFile(file);
          }}
        />
      </button>

      {uploading && progress !== null && (
        <div className="w-full bg-gray-200 rounded-full h-1 overflow-hidden">
          <div
            className="bg-orange-500 h-full rounded-full transition-all duration-300"
            style={{ width: `${progress}%` }}
          />
        </div>
      )}

      {error && <div className="border-red-500/30 bg-red-500/10">{error}</div>}
    </div>
  );
}

type SketchVideoUploadProps = {
  sketchId: number;
};

export function SketchVideoUpload({ sketchId }: SketchVideoUploadProps) {
  const queryClient = useQueryClient();
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState<number | null>(null);

  // Determine if we should poll (any active jobs)
  const shouldPoll = useCallback((data: VideosResponse | undefined) => {
    if (!data) return false;
    return data.videos.some((v) =>
      v.jobs.some((j) => j.status === "pending" || j.status === "processing"),
    );
  }, []);

  const { data, isLoading, isError } = useSketchVideos({ id: sketchId });
  const handleFile = useCallback(
    async (file: File) => {
      setUploadError(null);
      setUploading(true);
      setProgress(0);

      try {
        // 1. Get signed URL
        setProgress(10);
        const { uploadUrl, s3Key } = await getUploadUrl(
          file.name,
          file.type,
          file.size,
        );

        // 2. Upload to S3 with XHR for progress
        await new Promise<void>((resolve, reject) => {
          const xhr = new XMLHttpRequest();
          xhr.open("PUT", uploadUrl);
          xhr.setRequestHeader("Content-Type", file.type);
          xhr.setRequestHeader("Content-Length", String(file.size));
          xhr.upload.onprogress = (e) => {
            if (e.lengthComputable) {
              setProgress(10 + Math.round((e.loaded / e.total) * 80));
            }
          };
          xhr.onload = () =>
            xhr.status >= 200 && xhr.status < 300
              ? resolve()
              : reject(new Error(`S3 upload failed: ${xhr.status}`));
          xhr.onerror = () => reject(new Error("Network error during upload"));
          xhr.send(file);
        });

        // 3. Notify backend
        setProgress(95);
        const result = await notifyUploaded(sketchId, s3Key);
        setProgress(100);

        // 4. Seed query cache with fresh data
        queryClient.setQueryData(["sketch-videos", sketchId], result);
      } catch (err) {
        setUploadError(err instanceof Error ? err.message : "Upload failed");
      } finally {
        setUploading(false);
        setTimeout(() => setProgress(null), 800);
      }
    },
    [queryClient, sketchId],
  );

  const videos = data?.videos ?? [];
  const isEmpty = !isLoading && !isError && videos.length === 0;

  return (
    <div className="pt-4 space-y-3">
      {isError && (
        <div className="border-red-500/30 bg-red-500/10">
          Failed to load videos.
        </div>
      )}

      {isLoading && (
        <div className="flex items-center justify-center py-10 gap-2 text-gray-600">
          <Spinner className="w-4 h-4" />
          <span className="text-xs">Loading…</span>
        </div>
      )}

      {isEmpty && (
        <UploadZone
          onSubmit={handleFile}
          uploading={uploading}
          progress={progress}
          error={uploadError}
        />
      )}

      {!isLoading && videos.length > 0 && (
        <>
          <div className="space-y-3">
            {videos.map((video) => (
              <VideoCard key={video.id} video={video} />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
