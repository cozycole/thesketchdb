import { useState } from "react";
import { ImageOff } from "lucide-react";
import { CastScreenshot } from "@/types/api";
import { useCastFormStore } from "../stores/castFormStore";
import { buildImageUrl } from "@/lib/utils";

interface ScreenshotGridProps {
  screenshots: CastScreenshot[];
}

export function ScreenshotGrid({ screenshots }: ScreenshotGridProps) {
  const { openForm } = useCastFormStore();

  const [failedImages, setFailedImages] = useState<Set<number>>(new Set());
  const [dims, setDims] = useState<Record<number, { w: number; h: number }>>(
    {},
  );

  const handleScreenshotClick = (screenshot: CastScreenshot) => {
    openForm(undefined, screenshot);
  };

  const handleImageError = (
    screenshotId: number,
    e: React.SyntheticEvent<HTMLImageElement>,
  ) => {
    setFailedImages((prev) => new Set(prev).add(screenshotId));
    e.currentTarget.removeAttribute("src");
  };

  const handleImageLoad = (
    screenshotId: number,
    e: React.SyntheticEvent<HTMLImageElement>,
  ) => {
    const img = e.currentTarget;
    setDims((prev) => ({
      ...prev,
      [screenshotId]: { w: img.naturalWidth, h: img.naturalHeight },
    }));
  };

  return (
    <div className="space-y-4 w-full lg:max-w-[500px]">
      <div className="grid grid-cols-3 gap-4">
        {screenshots.length === 0 ? (
          <p className="col-span-full mt-10 text-center text-muted-foreground">
            No screenshots available. Run pipeline to generate screenshots
          </p>
        ) : (
          screenshots.map((screenshot) => (
            <div
              key={screenshot.id}
              className="relative aspect-square cursor-pointer overflow-hidden rounded-md bg-muted transition-all hover:opacity-80 hover:ring-2 hover:ring-primary"
              onClick={() => handleScreenshotClick(screenshot)}
            >
              {failedImages.has(screenshot.id) ? (
                <div className="flex h-full w-full items-center justify-center">
                  <ImageOff className="h-8 w-8 text-muted-foreground" />
                </div>
              ) : (
                <>
                  <img
                    src={buildImageUrl(
                      "cast_auto_screenshots/profile",
                      "",
                      screenshot.profileImage,
                    )}
                    alt={`Screenshot ${screenshot.id}`}
                    className="h-full w-full object-cover"
                    onError={(e) => handleImageError(screenshot.id, e)}
                    onLoad={(e) => handleImageLoad(screenshot.id, e)}
                  />

                  {dims[screenshot.id] && (
                    <div className="absolute bottom-1 left-1 rounded bg-black/70 px-1.5 py-0.5 text-[11px] text-white">
                      {dims[screenshot.id].w}×{dims[screenshot.id].h}
                    </div>
                  )}
                </>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
}
