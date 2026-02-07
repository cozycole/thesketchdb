import { ContentLayout } from "@/components/layouts/content";
import { SketchForm } from "@/features/sketches/forms/sketchForm";

const CreateSketchRoute = () => {
  return (
    <ContentLayout title="Create Sketch">
      <SketchForm mode="create" />
    </ContentLayout>
  );
};

export default CreateSketchRoute;
