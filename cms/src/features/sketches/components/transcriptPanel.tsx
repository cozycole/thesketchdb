import { TranscriptLine } from "@/types/api";
import { formatHMS, cn } from "@/lib/utils";
import { NewspaperIcon } from "lucide-react";

import { useDraggable } from "@dnd-kit/core";

function DraggableTranscriptLine({ line }: { line: TranscriptLine }) {
  const { attributes, listeners, setNodeRef, transform } = useDraggable({
    id: `transcript-${line.id}`,
    data: { type: "transcriptLine", lineId: line.id },
  });

  const style: React.CSSProperties | undefined = transform
    ? { transform: `translate3d(${transform.x}px, ${transform.y}px, 0)` }
    : undefined;

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      className={cn(
        `
        group flex gap-3 p-3 rounded-lg bg-gray-100 border border-gray-200
        hover:bg-gray-50 transition-colors cursor-grab active:cursor-grabbing
        `,
      )}
    >
      <div
        className="
          text-xs font-mono text-orange-600 bg-orange-50 px-2 py-1
          rounded shrink-0 self-start
        "
      >
        {formatHMS(Math.round(line.startMs / 1000))}
      </div>

      <div className="text-sm text-gray-800 leading-relaxed">{line.text}</div>
    </div>
  );
}

interface TranscriptPanelProps {
  transcript: TranscriptLine[];
  onSelect: (id: number, e: MouseEvent) => void;
  selectedIds: Set<number>;
}

export function TranscriptPanel({
  transcript,
  selectedIds,
  onSelect,
}: TranscriptPanelProps) {
  return (
    <div className="grid grid-cols-1 gap-4 my-2 mx-2">
      {transcript.length === 0 ? (
        <p className="col-span-full mt-10 text-center text-muted-foreground">
          <NewspaperIcon className="mx-auto mb-6" />
          No transcript available. Upload a video to generate one.
        </p>
      ) : (
        transcript.map((line) => (
          <div
            key={line.id}
            onClick={(e) => onSelect(line.id, e)}
            className={
              selectedIds.has(line.id)
                ? "rounded-lg ring-2 ring-orange-400"
                : ""
            }
          >
            <DraggableTranscriptLine line={line} />
          </div>
        ))
      )}
    </div>
  );
}
