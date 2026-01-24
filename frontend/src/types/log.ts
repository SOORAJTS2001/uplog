export type LogLevel = 'ERROR' | 'WARN' | 'INFO' | 'DEBUG';

export interface LogEntry {
  id: string;
  message: string;
  timestamp: string;
  level: LogLevel;
}

export interface LogStream {
  id: string;
  name: string;
  createdAt: string;
}
