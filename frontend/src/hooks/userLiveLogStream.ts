import { useState, useCallback, useEffect, useRef } from "react";
import { LogEntry, LogLevel } from "@/types/log";

export async function* connectSSE(
  sessionId: string
): AsyncGenerator<Omit<LogEntry, "id">[], void, unknown> {
  const backendUrl = import.meta.env.VITE_BACKEND_URL;
  console.log(backendUrl,"backendUrl")

  const url = `${backendUrl}session/consume?session_id=${encodeURIComponent(sessionId)}`;
  const source = new EventSource(url);

  let queue: Omit<LogEntry, "id">[][] = [];
  let resolveNext:
    | ((value: IteratorResult<Omit<LogEntry, "id">[]>) => void)
    | null = null;

  source.onmessage = (event) => {
    try {
      const parsed = JSON.parse(event.data);
      console.log(parsed)
      if (!Array.isArray(parsed)) {
        console.error("SSE payload is not an array:", parsed);
        return;
      }

      const batch = parsed.map((item) => ({
        message: String(item.message),
        timestamp: String(item.timestamp),
        log_level: item.log_level as LogLevel,
      }));

      if (resolveNext) {
        resolveNext({ value: batch, done: false });
        resolveNext = null;
      } else {
        queue.push(batch);
      }
    } catch (err) {
      console.error("Invalid SSE payload:", event.data, err);
    }
  };

  source.onerror = (err) => {
    console.error("SSE connection error", err);
    source.close();

    if (resolveNext) {
      resolveNext({ value: undefined as any, done: true });
      resolveNext = null;
    }
  };

  try {
    while (true) {
      if (queue.length > 0) {
        yield queue.shift()!;
      } else {
        const next = await new Promise<
          IteratorResult<Omit<LogEntry, "id">[]>
        >((resolve) => {
          resolveNext = resolve;
        });

        if (next.done) return;
        yield next.value;
      }
    }
  } finally {
    // ALWAYS close the connection
    source.close();
  }
}


function makeLogId() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

export function useLogStream(streamId: string | undefined) {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [newLogIds, setNewLogIds] = useState<Set<string>>(new Set());

  // used to stop the async loop
  const cancelRef = useRef(false);

const startStream = useCallback(() => {
  if (!streamId || isStreaming) return;

  cancelRef.current = false;   // ← REQUIRED
  setIsStreaming(true);

  (async () => {
    try {
      for await (const batch of connectSSE(streamId)) {
        // if (cancelRef.current) break;

        const logsWithIds = batch.map((log) => ({
          ...log,
          id: makeLogId(),
        }));

        setLogs((prev) => [...prev, ...logsWithIds].slice(-500));

        setNewLogIds((prev) => {
          const next = new Set(prev);
          logsWithIds.forEach((l) => next.add(l.id));
          return next;
        });

        setTimeout(() => {
          setNewLogIds((current) => {
            const next = new Set(current);
            logsWithIds.forEach((l) => next.delete(l.id));
            return next;
          });
        }, 1000);
      }
    } finally {
      setIsStreaming(false);
    }
  })();
}, [streamId, isStreaming]);

  const stopStream = useCallback(() => {
    cancelRef.current = true;
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

  // auto-start when streamId changes
  useEffect(() => {
    if (!streamId) return;

    startStream();

    return () => {
      cancelRef.current = true;
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
