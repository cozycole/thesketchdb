import { TranscriptLine } from "@/types/api";
import { formatHMS } from "@/lib/utils";
import { NewspaperIcon } from "lucide-react";

interface TranscriptPanelProps {
  transcript: TranscriptLine[];
}

export function TranscriptPanel({ transcript }: TranscriptPanelProps) {
  return (
    <div className="grid grid-cols-1 gap-4 my-2 mx-2">
      {transcript.length === 0 ? (
        <p className="col-span-full mt-10 text-center text-muted-foreground">
          <NewspaperIcon className="mx-auto mb-6" />
          No transcript available. Upload a video to generate one.
        </p>
      ) : (
        transcript.map((line) => (
          <div className="p-3 rounded-lg bg-gray-200" key={line.id}>
            <div>{formatHMS(Math.round(line.startMs / 1000))}</div>
            {line.text}
          </div>
        ))
      )}
    </div>
  );
}
