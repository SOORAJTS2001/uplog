import { LogLevel } from '@/types/log';
import { Button } from './ui/button';
import { cn } from '@/lib/utils';

interface LogFilterProps {
  selectedLevels: LogLevel[];
  onToggleLevel: (level: LogLevel) => void;
}

const levels: LogLevel[] = ['ERROR', 'WARN', 'INFO', 'DEBUG'];

const levelColors: Record<LogLevel, string> = {
  ERROR: 'data-[active=true]:bg-log-error/20 data-[active=true]:text-log-error data-[active=true]:border-log-error/50',
  WARN: 'data-[active=true]:bg-log-warn/20 data-[active=true]:text-log-warn data-[active=true]:border-log-warn/50',
  INFO: 'data-[active=true]:bg-log-info/20 data-[active=true]:text-log-info data-[active=true]:border-log-info/50',
  DEBUG: 'data-[active=true]:bg-log-debug/20 data-[active=true]:text-log-debug data-[active=true]:border-log-debug/50',
};

export function LogFilter({ selectedLevels, onToggleLevel }: LogFilterProps) {
  return (
    <div className="flex items-center gap-2">
      <span className="text-sm text-muted-foreground mr-1">Filter:</span>
      {levels.map((level) => {
        const isActive = selectedLevels.includes(level);
        return (
          <Button
            key={level}
            variant="outline"
            size="sm"
            onClick={() => onToggleLevel(level)}
            data-active={isActive}
            className={cn(
              'h-7 px-2.5 text-xs font-mono uppercase border-border/50 transition-all',
              levelColors[level],
              !isActive && 'text-muted-foreground hover:text-foreground'
            )}
          >
            {level}
          </Button>
        );
      })}
    </div>
  );
}
