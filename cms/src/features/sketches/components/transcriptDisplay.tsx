import { TranscriptLine } from "@/types/api";

interface TranscriptDisplayProps {
  transcript: TranscriptLine[];
}

export function TranscriptDisplay({ transcript }: TranscriptDisplayProps) {
  return (
    <div className="grid grid-cols-1 gap-4 my-2 mx-4">
      {transcript.length === 0 ? (
        <p className="col-span-full mt-10 text-center text-muted-foreground">
          No transcription available. Run pipeline to generate one.
        </p>
      ) : (
        transcript.map((line) => (
          <div className="p-3 rounded-lg bg-gray-200">{line.text}</div>
        ))
      )}
    </div>
  );
}
