import { DraggableCastTable } from "@/features/sketches/components/castTable";
import { ScreenshotGrid } from "@/features/sketches/components/screenshotGrid";
import { CastForm } from "@/features/sketches/forms/castForm";

import { Loader2, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";

import { useCast } from "@/features/sketches/api/getCast";
import { useCastFormStore } from "@/features/sketches/stores/castFormStore";

export function EditCastPage({ sketchId }) {
  const { data, isLoading: castLoading } = useCast({
    id: sketchId,
  });
  const { openForm } = useCastFormStore();
  return (
    <div className="flex h-full min-h-0 flex-col overflow-hidden">
      <div className="flex min-h-0 flex-1 flex-col overflow-hidden">
        <div className="mb-4 shrink-0">
          <Button className="text-white font-bold" onClick={() => openForm()}>
            <Plus className="mr-2 h-4 w-4" />
            Add Cast Member
          </Button>
        </div>
        {castLoading || !data ? (
          <div className="flex min-h-0 flex-1 items-center justify-center">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : (
          <div className="flex flex-col flex-1 min-h-0 xl:flex-row gap-4 overflow-hidden">
            <div className="min-h-0 min-w-0 flex-1 overflow-hidden">
              <DraggableCastTable sketchId={sketchId} cast={data.cast} />
            </div>

            <div className="min-h-0 min-w-0 flex-1 overflow-hidden xl:max-w-[500px]">
              <ScreenshotGrid screenshots={data.screenshots} />
            </div>
          </div>
        )}
      </div>

      <CastForm sketchId={sketchId} />
    </div>
  );
}
