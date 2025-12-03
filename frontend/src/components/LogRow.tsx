import { LogEntry } from '@/types/log';
import { LogLevelBadge } from './LogLevelBadge';
import { format } from 'date-fns';

interface LogRowProps {
  log: LogEntry;
  isNew?: boolean;
}

export function LogRow({ log, isNew }: LogRowProps) {
  const formattedTime = format(new Date(log.timestamp), 'HH:mm:ss.SSS');
  const formattedDate = format(new Date(log.timestamp), 'MMM dd');

  return (
    <div className={`log-row grid grid-cols-[140px_80px_1fr] gap-4 px-4 py-3 ${isNew ? 'animate-fade-in bg-primary/5' : ''}`}>
      <div className="flex items-center gap-2 text-muted-foreground font-mono text-sm">
        <span className="text-foreground/70">{formattedDate}</span>
        <span>{formattedTime}</span>
      </div>
      <div className="flex items-center">
        <LogLevelBadge level={log.log_level} />
      </div>
      <div className="font-mono text-sm text-foreground/90 break-all">
        {log.message}
      </div>
    </div>
  );
}
