import { Loader2, Plus } from "lucide-react";

import { QueryClient } from "@tanstack/react-query";
import { ContentLayout } from "@/components/layouts/content";
import { useParams } from "react-router";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { DraggableCastTable } from "@/features/sketches/components/castTable";
import { ScreenshotGrid } from "@/features/sketches/components/screenshotGrid";

import {
  useSketch,
  sketchQueryOptions,
} from "@/features/sketches/api/getSketch";
import { SketchForm } from "@/features/sketches/forms/sketchForm";

import { useCast } from "@/features/sketches/api/getCast";
import { useCastFormStore } from "@/features/sketches/stores/castFormStore";
import { CastForm } from "@/features/sketches/forms/castForm";

export const clientLoader =
  (queryClient: QueryClient) =>
  async ({ params }: any) => {
    const { id } = params;
    await queryClient.ensureQueryData(sketchQueryOptions({ id: Number(id) }));
    return null;
  };

const UpdateSketchRoute = () => {
  const { id } = useParams<{ id: string }>();

  const { data: sketchData, isLoading: sketchLoading } = useSketch({
    id: Number(id),
  });
  const sketch = sketchData.sketch;

  const { data: castData, isLoading: castLoading } = useCast({
    id: Number(id),
  });
  const { openForm } = useCastFormStore();

  return (
    <ContentLayout title={`Sketch ID ${sketch.id}`}>
      <Tabs defaultValue="cast">
        <TabsList variant="line" className="w-[400px] mb-4 border-orange">
          <TabsTrigger value="metadata">Metadata</TabsTrigger>
          <TabsTrigger value="cast">Cast</TabsTrigger>
          <TabsTrigger value="quotes">Quotes</TabsTrigger>
        </TabsList>
        <TabsContent value="metadata">
          {sketchLoading ? (
            <div className="flex h-screen items-center justify-center">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : (
            <SketchForm mode="update" existingData={sketch} />
          )}
        </TabsContent>
        <TabsContent value="cast">
          <div>
            <div className="mb-4">
              <Button onClick={() => openForm()}>
                <Plus className="mr-2 h-4 w-4" />
                Add Cast Member
              </Button>
            </div>
            {castLoading || !castData ? (
              <div className="flex h-screen items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
              </div>
            ) : (
              <div className="flex flex-col lg:flex-row gap-4">
                <DraggableCastTable
                  sketchId={sketch.id}
                  cast={castData.cast}
                ></DraggableCastTable>
                <ScreenshotGrid
                  screenshots={castData.screenshots}
                ></ScreenshotGrid>
              </div>
            )}
            <CastForm sketchId={sketch.id} />
          </div>
        </TabsContent>
        <TabsContent value="quotes">
          <h1>Quotes Form</h1>
        </TabsContent>
      </Tabs>
    </ContentLayout>
  );
};

export default UpdateSketchRoute;
