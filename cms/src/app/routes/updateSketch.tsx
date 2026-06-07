import { useBlocker } from "react-router";
import { QueryClient } from "@tanstack/react-query";
import { FixedContentLayout } from "@/components/layouts/fixedContent";
import { useParams } from "react-router";

import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogTitle,
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogDescription,
} from "@/components/ui/alert-dialog";

import { EditCastPage } from "@/features/sketches/pages/editCastPage";
import { EditQuotesPage } from "@/features/sketches/pages/editQuotesPage";
import { EditSketchPage } from "@/features/sketches/pages/editSketchPage";
import { SketchVideoUpload } from "@/features/sketches/components/videoUpload";

import { sketchQueryOptions } from "@/features/sketches/api/getSketch";

import {
  SketchDirtyProvider,
  useSketchDirty,
} from "@/features/sketches/contexts/sketchDirtyContext";
import { useEffect } from "react";

export const clientLoader =
  (queryClient: QueryClient) =>
  async ({ params }: any) => {
    const { id } = params;
    await queryClient.ensureQueryData(sketchQueryOptions({ id: Number(id) }));
    return null;
  };

const UpdateSketchRoute = () => {
  return (
    <SketchDirtyProvider>
      <UpdateSketchComponent />
    </SketchDirtyProvider>
  );
};

const UpdateSketchComponent = () => {
  const { id } = useParams<{ id: string }>();
  const { hasUnsavedChanges, dirtyMap } = useSketchDirty();
  const {
    state: blockerState,
    reset: blockerReset,
    proceed: blockerProceed,
  } = useBlocker(hasUnsavedChanges);
  const dirtyTabs = Object.entries(dirtyMap)
    .filter(([, isDirty]) => isDirty)
    .map(([tabName]) => tabName);

  useEffect(() => {
    if (!hasUnsavedChanges) return;

    const handler = (event: BeforeUnloadEvent) => {
      event.preventDefault();
      event.returnValue = "";
    };

    window.addEventListener("beforeunload", handler);
    return () => window.removeEventListener("beforeunload", handler);
  }, [hasUnsavedChanges]);

  return (
    <FixedContentLayout title={`Sketch ID ${id}`}>
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
          forceMount
          value="metadata"
          className="w-full flex-1 min-h-0 overflow-hidden data-[state=inactive]:hidden"
        >
          <EditSketchPage id={Number(id)} />
        </TabsContent>
        <TabsContent
          forceMount
          value="cast"
          className="w-full flex-1 min-h-0 overflow-hidden data-[state=inactive]:hidden"
        >
          <EditCastPage sketchId={id} />
        </TabsContent>
        <TabsContent
          forceMount
          className="flex-1 min-h-0 overflow-hidden data-[state=inactive]:hidden"
          value="quotes"
        >
          <EditQuotesPage sketchId={Number(id)} />
        </TabsContent>
        <TabsContent
          forceMount
          value="video"
          className="data-[state=inactive]:hidden"
        >
          <SketchVideoUpload sketchId={Number(id)} />
        </TabsContent>
      </Tabs>
      {blockerState === "blocked" && (
        <AlertDialog open>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Unsaved changes</AlertDialogTitle>
              <AlertDialogDescription>
                You have unsaved changes in:
                <ul className="my-2 list-disc pl-5">
                  {dirtyTabs.map((tab) => (
                    <li key={tab}>{tab}</li>
                  ))}
                </ul>
                Are you sure you want to continue?
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel onClick={() => blockerReset()}>
                Stay
              </AlertDialogCancel>
              <AlertDialogAction onClick={() => blockerProceed()}>
                Leave
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}
    </FixedContentLayout>
  );
};

export default UpdateSketchRoute;
