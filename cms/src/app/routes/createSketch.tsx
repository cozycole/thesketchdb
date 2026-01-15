import { ContentLayout } from "@/components/layouts/content";
import { SketchForm } from "@/features/sketches/forms/sketchForm";

const AddSketchRoute = () => {
  return (
    <ContentLayout title="Add Sketch">
      <SketchForm mode="create" />
    </ContentLayout>
  );
};

export default AddSketchRoute;
