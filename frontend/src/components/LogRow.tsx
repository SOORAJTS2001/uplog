import { LogEntry } from '@/types/log';
import { LogLevelBadge } from './LogLevelBadge';
import { format } from 'date-fns';

interface LogRowProps {
  log: LogEntry;
  isNew?: boolean;
}

export function LogRow({ log, isNew }: LogRowProps) {
  const date = new Date(log.timestamp);

  return (
    <div
      className={[
        'grid grid-cols-[90px_64px_1fr] gap-3 px-4 py-1.5',
        'font-mono text-sm items-start',
        isNew ? 'animate-fade-in bg-primary/5' : ''
      ].join(' ')}
    >
      {/* Timestamp (compact) */}
      <div className="text-muted-foreground leading-tight">
        <div className="text-foreground/70 text-xs">
          {format(date, 'MMM dd')}
        </div>
        <div className="text-xs">
          {format(date, 'HH:mm:ss')}
        </div>
      </div>

      {/* Level (tight) */}
      <div className="flex items-start">
        <LogLevelBadge level={log.level} />
      </div>

      {/* Message (gets space) */}
      <div className="text-foreground/90 whitespace-pre-wrap break-words leading-relaxed">
        {log.message}
      </div>
    </div>
  );
}
