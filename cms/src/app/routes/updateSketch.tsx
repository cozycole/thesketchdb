import { QueryClient } from "@tanstack/react-query";
import { ContentLayout } from "@/components/layouts/content";
import { useParams } from "react-router";

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

  const { data } = useSketch({ id: Number(id) });
  const sketch = data.sketch;

  return (
    <ContentLayout title={`Update Sketch ID ${sketch.id}`}>
      <SketchForm mode="update" existingData={data.sketch} />
    </ContentLayout>
  );
};

export default UpdateSketchRoute;
