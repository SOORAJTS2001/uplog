import { Radio } from 'lucide-react';

interface StreamingIndicatorProps {
  isStreaming: boolean;
  logCount: number;
}

export function StreamingIndicator({ isStreaming, logCount }: StreamingIndicatorProps) {
  return (
    <div className="flex items-center gap-3">
      <div className="flex items-center gap-2">
        {isStreaming ? (
          <>
            <div className="streaming-dot" />
            <span className="text-sm text-primary font-medium">Live</span>
          </>
        ) : (
          <>
            <div className="w-2 h-2 rounded-full bg-muted-foreground/50" />
            <span className="text-sm text-muted-foreground">Paused</span>
          </>
        )}
      </div>
      <div className="h-4 w-px bg-border" />
      <div className="flex items-center gap-1.5 text-muted-foreground">
        <Radio className="w-3.5 h-3.5" />
        <span className="text-sm font-mono">{logCount.toLocaleString()} events</span>
      </div>
    </div>
  );
}
