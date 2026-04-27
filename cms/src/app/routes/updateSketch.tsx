import { Loader2 } from "lucide-react";

import { QueryClient } from "@tanstack/react-query";
import { ContentLayout } from "@/components/layouts/content";
import { useParams } from "react-router";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { EditCastPage } from "@/features/sketches/editCastPage";
import { EditQuotesPage } from "@/features/sketches/components/editQuotesPage";
import { SketchVideoUpload } from "@/features/sketches/components/videoUpload";

import {
  useSketch,
  sketchQueryOptions,
} from "@/features/sketches/api/getSketch";
import { SketchForm } from "@/features/sketches/forms/sketchForm";

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

  return (
    <ContentLayout title={`Sketch ID ${sketch.id}`}>
      <Tabs defaultValue="metadata" className="flex h-full min-h-0 flex-col">
        <TabsList
          variant="line"
          className="shrink-0 w-[400px] my-2 border-orange"
        >
          <TabsTrigger value="metadata">Metadata</TabsTrigger>
          <TabsTrigger value="cast">Cast</TabsTrigger>
          <TabsTrigger value="quotes">Quotes</TabsTrigger>
          <TabsTrigger value="video">Video</TabsTrigger>
        </TabsList>
        <TabsContent
          value="metadata"
          className="w-full flex-1 min-h-0 overflow-hidden"
        >
          {sketchLoading ? (
            <div className="flex h-screen gap-2 mt-20 justify-center">
              Loading
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : (
            <SketchForm mode="update" existingData={sketch} />
          )}
        </TabsContent>
        <TabsContent
          value="cast"
          className="w-full flex-1 min-h-0 overflow-hidden"
        >
          <EditCastPage sketchId={sketch.id} />
        </TabsContent>
        <TabsContent className="flex-1 min-h-0 overflow-hidden" value="quotes">
          <EditQuotesPage sketchId={sketch.id} />
        </TabsContent>
        <TabsContent value="video">
          <SketchVideoUpload sketchId={sketch.id} />
        </TabsContent>
      </Tabs>
    </ContentLayout>
  );
};

export default UpdateSketchRoute;
