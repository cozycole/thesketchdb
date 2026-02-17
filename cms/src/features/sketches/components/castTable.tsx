import { useState, useEffect } from "react";
import { Image } from "lucide-react";
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { GripVertical, Edit } from "lucide-react";
import { CastMember } from "@/types/api";
import { useCastFormStore } from "../stores/castFormStore";
import { useUpdateCastOrder } from "../api/updateCastOrder";

import { buildImageUrl } from "@/lib/utils";

interface SortableCastRowProps {
  castMember: CastMember;
  onEdit: (cast: CastMember) => void;
}

function SortableCastRow({ castMember, onEdit }: SortableCastRowProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: castMember.id,
  });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const displayName = castMember.actor
    ? `${castMember.actor.first} ${castMember.actor.last}`
    : "";

  const profileImage = castMember.profileImage
    ? buildImageUrl("cast/profile", "small", castMember.profileImage)
    : castMember.actor?.profileImage
      ? buildImageUrl("person", "small", castMember.actor.profileImage)
      : undefined;

  return (
    <TableRow ref={setNodeRef} style={style}>
      <TableCell className="w-12">
        <button
          className="cursor-grab active:cursor-grabbing"
          {...attributes}
          {...listeners}
        >
          <GripVertical className="h-5 w-5 text-muted-foreground" />
        </button>
      </TableCell>
      <TableCell className="font-medium">
        {profileImage ? (
          <img className="min-w-14" src={profileImage} />
        ) : (
          <div className="flex h-full w-full items-center justify-center">
            <Image className="h-14 w-14 text-muted-foreground" />
          </div>
        )}
      </TableCell>
      <TableCell>{displayName}</TableCell>
      <TableCell>{castMember.characterName}</TableCell>
      <TableCell>
        {castMember.castRole
          ? castMember.castRole.charAt(0).toUpperCase() +
            castMember.castRole.slice(1)
          : ""}
      </TableCell>
      <TableCell>{castMember.minorRole ? "Yes" : "No"}</TableCell>
      <TableCell>
        <Button variant="ghost" size="sm" onClick={() => onEdit(castMember)}>
          <Edit className="h-4 w-4" />
        </Button>
      </TableCell>
    </TableRow>
  );
}

interface DraggableCastTableProps {
  cast: CastMember[];
  sketchId: number;
}

export function DraggableCastTable({
  cast,
  sketchId,
}: DraggableCastTableProps) {
  const [localCast, setLocalCast] = useState<CastMember[]>(cast);
  const { openForm } = useCastFormStore();
  const { mutate: updateOrder } = useUpdateCastOrder({
    mutationConfig: {},
  });

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  // Update local state when prop changes
  useEffect(() => {
    setLocalCast(cast);
  }, [cast]);

  const handleDragEnd = ({ active, over }: DragEndEvent) => {
    if (!over || active.id === over.id) return;

    const oldIndex = localCast.findIndex((i) => i.id === active.id);
    const newIndex = localCast.findIndex((i) => i.id === over.id);
    if (oldIndex < 0 || newIndex < 0) return;

    const newItems = arrayMove(localCast, oldIndex, newIndex);

    setLocalCast(
      newItems.map((item, index) => ({
        ...item,
        position: index + 1,
      })),
    );

    updateOrder({
      sketchId,
      castPositions: newItems.map((i) => i.id),
    });
  };

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <SortableContext
        items={localCast.map((c) => c.id)}
        strategy={verticalListSortingStrategy}
      >
        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-12"></TableHead>
                <TableHead className="w-24"></TableHead>
                <TableHead className="w-52">Actor</TableHead>
                <TableHead className="w-52">Character</TableHead>
                <TableHead className="w-20">Role</TableHead>
                <TableHead className="w-16">Minor</TableHead>
                <TableHead className="w-20">Edit</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {localCast.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={7}
                    className="text-center text-muted-foreground"
                  >
                    No cast members yet. Add one to get started.
                  </TableCell>
                </TableRow>
              ) : (
                localCast.map((castMember) => (
                  <SortableCastRow
                    key={castMember.id}
                    castMember={castMember}
                    onEdit={openForm}
                  />
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </SortableContext>
    </DndContext>
  );
}
