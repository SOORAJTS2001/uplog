import { useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { LogTable } from '@/components/LogTable';
import { LogFilter } from '@/components/LogFilter';
import { StreamingIndicator } from '@/components/StreamingIndicator';
import { useLogStream } from '@/hooks/useLogStream';
import { LogLevel } from '@/types/log';
import {
  Terminal,
  Pause,
  Play,
  Trash2,
  Copy,
  Check,
  ArrowDown,
  ArrowLeft
} from 'lucide-react';
import { toast } from 'sonner';

export default function Dashboard() {
  const { streamId } = useParams<{ streamId: string }>();
  const navigate = useNavigate();
  const { logs, isStreaming, newLogIds, toggleStream, clearLogs } = useLogStream(streamId);

  const [selectedLevels, setSelectedLevels] = useState<LogLevel[]>([]);
  const [autoScroll, setAutoScroll] = useState(true);
  const [copied, setCopied] = useState(false);

  const handleToggleLevel = useCallback((level: LogLevel) => {
    setSelectedLevels((prev) =>
      prev.includes(level)
        ? prev.filter((l) => l !== level)
        : [...prev, level]
    );
  }, []);

  const handleCopyUrl = useCallback(async () => {
    const url = window.location.href;
    await navigator.clipboard.writeText(url);
    setCopied(true);
    toast.success('Stream URL copied to clipboard');
    setTimeout(() => setCopied(false), 2000);
  }, []);

  const handleClear = useCallback(() => {
    clearLogs();
    toast.success('Logs cleared');
  }, [clearLogs]);

  return (
    <div className="min-h-screen bg-background flex flex-col">
      {/* Header */}
      <header className="border-b border-border/50 bg-card/50 backdrop-blur-sm sticky top-0 z-10">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Button
                variant="ghost"
                size="icon"
                onClick={() => navigate('/')}
                className="text-muted-foreground hover:text-foreground"
              >
                <ArrowLeft className="w-4 h-4" />
              </Button>
              <div className="flex items-center gap-2">
                <div className="w-8 h-8 rounded-lg bg-primary/20 flex items-center justify-center">
                  <Terminal className="w-4 h-4 text-primary" />
                </div>
                <div>
                  <span className="font-semibold">LogStream</span>
                  <p className="text-xs text-muted-foreground font-mono">{streamId}</p>
                </div>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={handleCopyUrl}
                className="gap-2"
              >
                {copied ? (
                  <Check className="w-3.5 h-3.5" />
                ) : (
                  <Copy className="w-3.5 h-3.5" />
                )}
                {copied ? 'Copied!' : 'Copy URL'}
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Toolbar */}
      <div className="border-b border-border/50 bg-card/30">
        <div className="container mx-auto px-6 py-3">
          <div className="flex items-center justify-between flex-wrap gap-4">
            <div className="flex items-center gap-4">
              <StreamingIndicator isStreaming={isStreaming} logCount={logs.length} />
              <div className="h-4 w-px bg-border hidden sm:block" />
              <LogFilter selectedLevels={selectedLevels} onToggleLevel={handleToggleLevel} />
            </div>

            <div className="flex items-center gap-2">
              <Button
                variant={autoScroll ? 'secondary' : 'outline'}
                size="sm"
                onClick={() => setAutoScroll(!autoScroll)}
                className="gap-2"
              >
                <ArrowDown className="w-3.5 h-3.5" />
                Auto-scroll
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={handleClear}
                className="gap-2"
              >
                <Trash2 className="w-3.5 h-3.5" />
                Clear
              </Button>
              <Button
                variant={isStreaming ? 'secondary' : 'default'}
                size="sm"
                onClick={toggleStream}
                className="gap-2"
              >
                {isStreaming ? (
                  <>
                    <Pause className="w-3.5 h-3.5" />
                    Pause
                  </>
                ) : (
                  <>
                    <Play className="w-3.5 h-3.5" />
                    Resume
                  </>
                )}
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Log Table */}
      <div className="flex-1 container mx-auto px-6 py-6 flex flex-col min-h-0">
        <LogTable
          logs={logs}
          selectedLevels={selectedLevels}
          autoScroll={autoScroll}
          newLogIds={newLogIds}
        />
      </div>

      {/* Connection Info Footer */}
      <footer className="border-t border-border/50 bg-card/30 py-3">
        <div className="container mx-auto px-6">
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <div className="flex items-center gap-4">
              <span>Endpoint: <code className="text-foreground/70 bg-muted px-1.5 py-0.5 rounded">POST /api/logs/{streamId}</code></span>
            </div>
            <div>
              <span>Format: <code className="text-foreground/70 bg-muted px-1.5 py-0.5 rounded">{'{ message, timestamp, log_level }'}</code></span>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
