import { useRef, useEffect } from 'react';
import { LogEntry, LogLevel } from '@/types/log';
import { LogRow } from './LogRow';
import { ScrollArea } from './ui/scroll-area';

interface LogTableProps {
  logs: LogEntry[];
  selectedLevels: LogLevel[];
  autoScroll: boolean;
  newLogIds: Set<string>;
}

export function LogTable({ logs, selectedLevels, autoScroll, newLogIds }: LogTableProps) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);
  console.log(logs)

  const filteredLogs = logs.filter((log) =>
    selectedLevels.length === 0 || selectedLevels.includes(log.log_level)
  );

  useEffect(() => {
    if (autoScroll && bottomRef.current) {
      bottomRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [logs, autoScroll]);

  return (
    <div className="flex-1 glass-card overflow-hidden flex flex-col">
      {/* Header */}
      <div className="grid grid-cols-[140px_80px_1fr] gap-4 px-4 py-3 border-b border-border/50 bg-muted/30">
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">Timestamp</span>
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">Level</span>
        <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">Message</span>
      </div>

      {/* Log rows */}
      <ScrollArea className="flex-1" ref={scrollRef}>
        <div className="min-h-0">
          {filteredLogs.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
              <div className="w-12 h-12 rounded-full bg-muted/50 flex items-center justify-center mb-4">
                <div className="w-6 h-6 border-2 border-muted-foreground/30 border-t-muted-foreground rounded-full animate-spin" />
              </div>
              <p className="text-sm">Waiting for logs...</p>
              <p className="text-xs mt-1">Events will appear here in real-time</p>
            </div>
          ) : (
            filteredLogs.map((log) => (
              <LogRow
                key={log.id}
                log={log}
                isNew={newLogIds.has(log.id)}
              />
            ))
          )}
          <div ref={bottomRef} />
        </div>
      </ScrollArea>
    </div>
  );
}
