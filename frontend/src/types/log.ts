export type LogLevel = 'ERROR' | 'WARN' | 'INFO' | 'DEBUG';

export interface LogEntry {
  id: string;
  message: string;
  timestamp: string;
  log_level: LogLevel;
}

export interface LogStream {
  id: string;
  name: string;
  createdAt: string;
}
