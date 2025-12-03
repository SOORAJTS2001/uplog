import { LogLevel } from '@/types/log';
import { cn } from '@/lib/utils';

interface LogLevelBadgeProps {
  level: LogLevel;
}

const levelStyles: Record<LogLevel, string> = {
  ERROR: 'log-level-error',
  WARN: 'log-level-warn',
  INFO: 'log-level-info',
  DEBUG: 'log-level-debug',
};

export function LogLevelBadge({ level }: LogLevelBadgeProps) {
  return (
    <span className={cn('log-level-badge', levelStyles[level])}>
      {level}
    </span>
  );
}
