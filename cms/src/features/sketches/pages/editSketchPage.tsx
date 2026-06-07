import { Loader2 } from "lucide-react";

import { useSketch } from "@/features/sketches/api/getSketch";

import { useSketchDirty } from "../contexts/sketchDirtyContext";
import { SketchForm } from "@/features/sketches/forms/sketchForm";

type EditSketchPageProps = {
  id: number;
};
export function EditSketchPage({ id }: EditSketchPageProps) {
  const { data: sketchData, isLoading: sketchLoading } = useSketch({
    id: Number(id),
  });
  const { setDirty } = useSketchDirty();
  const sketch = sketchData.sketch;
  return sketchLoading ? (
    <div className="flex h-screen gap-2 mt-20 justify-center">
      Loading
      <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
    </div>
  ) : (
    <SketchForm mode="update" existingData={sketch} onDirtyChange={setDirty} />
  );
}
