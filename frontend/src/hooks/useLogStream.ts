import { useState, useCallback, useRef, useEffect } from 'react';
import { LogEntry, LogLevel } from '@/types/log';

// Mock log generator for demo
const generateMockLog = (): LogEntry => {
  const levels: LogLevel[] = ['ERROR', 'WARN', 'INFO', 'DEBUG'];
  const weights = [0.1, 0.15, 0.5, 0.25]; // Probability weights

  const random = Math.random();
  let cumulative = 0;
  let selectedLevel: LogLevel = 'INFO';

  for (let i = 0; i < levels.length; i++) {
    cumulative += weights[i];
    if (random < cumulative) {
      selectedLevel = levels[i];
      break;
    }
  }

  const messages: Record<LogLevel, string[]> = {
    ERROR: [
      'Connection refused: ECONNREFUSED 127.0.0.1:5432',
      'Failed to parse JSON response: Unexpected token',
      'Authentication failed for user: invalid_credentials',
      'Database query timeout after 30000ms',
      'Memory allocation failed: OutOfMemoryError',
    ],
    WARN: [
      'Deprecated API endpoint called: /api/v1/users',
      'Rate limit approaching: 950/1000 requests',
      'Certificate expires in 7 days',
      'Slow query detected: 2500ms execution time',
      'Cache miss ratio exceeds threshold: 45%',
    ],
    INFO: [
      'Server started on port 3000',
      'User session created: usr_8x7k2m',
      'Request completed: GET /api/health 200 OK',
      'Background job completed: email_notifications',
      'Database connection pool initialized: 10 connections',
      'Webhook delivered successfully to endpoint',
      'File uploaded: document_v2.pdf (2.4MB)',
    ],
    DEBUG: [
      'Entering function: processPayment()',
      'Cache key generated: user:profile:12345',
      'SQL query: SELECT * FROM users WHERE id = $1',
      'Response headers: Content-Type: application/json',
      'Memory usage: 245MB / 512MB',
    ],
  };

  return {
    id: `log_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
    message: messages[selectedLevel][Math.floor(Math.random() * messages[selectedLevel].length)],
    timestamp: new Date().toISOString(),
    log_level: selectedLevel,
  };
};

export function useLogStream(streamId: string | undefined) {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [newLogIds, setNewLogIds] = useState<Set<string>>(new Set());
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  const startStream = useCallback(() => {
    if (intervalRef.current) return;

    setIsStreaming(true);

    // Simulate initial batch of logs
    const initialLogs: LogEntry[] = [];
    for (let i = 0; i < 15; i++) {
      const log = generateMockLog();
      log.timestamp = new Date(Date.now() - (15 - i) * 2000).toISOString();
      initialLogs.push(log);
    }
    setLogs(initialLogs);

    // Simulate streaming new logs
    intervalRef.current = setInterval(() => {
      const newLog = generateMockLog();
      setLogs((prev) => [...prev.slice(-500), newLog]); // Keep last 500 logs
      setNewLogIds((prev) => {
        const next = new Set(prev);
        next.add(newLog.id);
        // Remove old IDs after animation
        setTimeout(() => {
          setNewLogIds((current) => {
            const updated = new Set(current);
            updated.delete(newLog.id);
            return updated;
          });
        }, 1000);
        return next;
      });
    }, Math.random() * 2000 + 500); // Random interval between 500ms and 2.5s
  }, []);

  const stopStream = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
    setIsStreaming(false);
  }, []);

  const toggleStream = useCallback(() => {
    if (isStreaming) {
      stopStream();
    } else {
      startStream();
    }
  }, [isStreaming, startStream, stopStream]);

  const clearLogs = useCallback(() => {
    setLogs([]);
    setNewLogIds(new Set());
  }, []);

  useEffect(() => {
    if (streamId) {
      startStream();
    }
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [streamId, startStream]);

  return {
    logs,
    isStreaming,
    newLogIds,
    startStream,
    stopStream,
    toggleStream,
    clearLogs,
  };
}
