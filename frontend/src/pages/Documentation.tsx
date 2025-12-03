import { useState } from "react";
import { NavLink } from "@/components/NavLink";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Copy, Check, Terminal, Send, Zap } from "lucide-react";

const CodeBlock = ({ code, language = "bash" }: { code: string; language?: string }) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative group">
      <pre className="bg-card/50 border border-border rounded-lg p-4 overflow-x-auto text-sm font-mono">
        <code className="text-foreground/90">{code}</code>
      </pre>
      <button
        onClick={handleCopy}
        className="absolute top-3 right-3 p-2 rounded-md bg-muted/50 hover:bg-muted transition-colors opacity-0 group-hover:opacity-100"
      >
        {copied ? (
          <Check className="w-4 h-4 text-accent" />
        ) : (
          <Copy className="w-4 h-4 text-muted-foreground" />
        )}
      </button>
    </div>
  );
};

const Section = ({
  title,
  description,
  children
}: {
  title: string;
  description?: string;
  children: React.ReactNode;
}) => (
  <section className="space-y-4">
    <div>
      <h2 className="text-2xl font-bold font-display text-foreground">{title}</h2>
      {description && <p className="text-muted-foreground mt-2">{description}</p>}
    </div>
    {children}
  </section>
);

export default function Documentation() {
  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border/50 sticky top-0 bg-background/80 backdrop-blur-sm z-50">
        <div className="max-w-4xl mx-auto px-6 py-4 flex items-center justify-between">
          <NavLink to="/" className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors">
            <ArrowLeft className="w-4 h-4" />
            <span>Back to Home</span>
          </NavLink>
          <span className="text-xl font-bold font-display neon-text">Live Logs</span>
        </div>
      </header>

      {/* Content */}
      <main className="max-w-4xl mx-auto px-6 py-12 space-y-16">
        {/* Hero */}
        <div className="space-y-4">
          <h1 className="text-4xl md:text-5xl font-bold font-display">
            <span className="neon-text">Documentation</span>
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl">
            Learn how to integrate Live Logs into your application and start streaming logs in minutes.
          </p>
        </div>

        {/* Quick Start */}
        <Section
          title="Quick Start"
          description="Get up and running with Live Logs in under 5 minutes."
        >
          <div className="space-y-6">
            <div className="flex items-start gap-4">
              <div className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0 mt-1">
                <span className="text-accent font-bold text-sm">1</span>
              </div>
              <div className="space-y-3 flex-1">
                <h3 className="font-semibold text-foreground">Create a Stream</h3>
                <p className="text-muted-foreground text-sm">
                  Click "Start for Free" on the homepage to generate a unique stream URL. No signup required.
                </p>
              </div>
            </div>

            <div className="flex items-start gap-4">
              <div className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0 mt-1">
                <span className="text-accent font-bold text-sm">2</span>
              </div>
              <div className="space-y-3 flex-1">
                <h3 className="font-semibold text-foreground">Send Logs via HTTP</h3>
                <p className="text-muted-foreground text-sm">
                  POST your logs to our API endpoint with your stream ID.
                </p>
                <CodeBlock
                  language="bash"
                  code={`curl -X POST https://api.livelogs.dev/stream/{YOUR_STREAM_ID} \\
  -H "Content-Type: application/json" \\
  -d '{
    "message": "User logged in successfully",
    "timestamp": "2024-01-15T10:30:00Z",
    "log_level": "INFO"
  }'`}
                />
              </div>
            </div>

            <div className="flex items-start gap-4">
              <div className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0 mt-1">
                <span className="text-accent font-bold text-sm">3</span>
              </div>
              <div className="space-y-3 flex-1">
                <h3 className="font-semibold text-foreground">View in Real-Time</h3>
                <p className="text-muted-foreground text-sm">
                  Open your dashboard URL and watch logs appear instantly as they're sent.
                </p>
              </div>
            </div>
          </div>
        </Section>

        {/* Log Format */}
        <Section
          title="Log Format"
          description="Structure your log entries using our simple JSON schema."
        >
          <div className="space-y-4">
            <CodeBlock
              language="json"
              code={`{
  "message": "string",      // Required: Log message content
  "timestamp": "datetime",  // Required: ISO 8601 format
  "log_level": "enum"       // Required: ERROR | WARN | INFO | DEBUG
}`}
            />

            <div className="bg-card/30 border border-border rounded-lg p-4">
              <h4 className="font-semibold text-foreground mb-3">Log Levels</h4>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                <div className="flex items-center gap-2">
                  <span className="w-3 h-3 rounded-full bg-log-error"></span>
                  <span className="text-sm text-muted-foreground">ERROR</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="w-3 h-3 rounded-full bg-log-warn"></span>
                  <span className="text-sm text-muted-foreground">WARN</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="w-3 h-3 rounded-full bg-log-info"></span>
                  <span className="text-sm text-muted-foreground">INFO</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="w-3 h-3 rounded-full bg-log-debug"></span>
                  <span className="text-sm text-muted-foreground">DEBUG</span>
                </div>
              </div>
            </div>
          </div>
        </Section>

        {/* Integration Examples */}
        <Section
          title="Integration Examples"
          description="Copy-paste examples for popular languages and frameworks."
        >
          <div className="space-y-8">
            {/* Node.js */}
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <Terminal className="w-5 h-5 text-accent" />
                <h3 className="font-semibold text-foreground">Node.js</h3>
              </div>
              <CodeBlock
                language="javascript"
                code={`const STREAM_ID = 'your-stream-id';
const API_URL = \`https://api.livelogs.dev/stream/\${STREAM_ID}\`;

async function sendLog(level, message) {
  await fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      message,
      timestamp: new Date().toISOString(),
      log_level: level
    })
  });
}

// Usage
sendLog('INFO', 'Server started on port 3000');
sendLog('ERROR', 'Database connection failed');`}
              />
            </div>

            {/* Python */}
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <Terminal className="w-5 h-5 text-accent" />
                <h3 className="font-semibold text-foreground">Python</h3>
              </div>
              <CodeBlock
                language="python"
                code={`import requests
from datetime import datetime

STREAM_ID = 'your-stream-id'
API_URL = f'https://api.livelogs.dev/stream/{STREAM_ID}'

def send_log(level: str, message: str):
    requests.post(API_URL, json={
        'message': message,
        'timestamp': datetime.utcnow().isoformat() + 'Z',
        'log_level': level
    })

# Usage
send_log('INFO', 'Application initialized')
send_log('WARN', 'High memory usage detected')`}
              />
            </div>

            {/* Go */}
            <div className="space-y-3">
              <div className="flex items-center gap-2">
                <Terminal className="w-5 h-5 text-accent" />
                <h3 className="font-semibold text-foreground">Go</h3>
              </div>
              <CodeBlock
                language="go"
                code={`package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"
)

const streamID = "your-stream-id"
const apiURL = "https://api.livelogs.dev/stream/" + streamID

type LogEntry struct {
    Message   string \`json:"message"\`
    Timestamp string \`json:"timestamp"\`
    LogLevel  string \`json:"log_level"\`
}

func sendLog(level, message string) error {
    entry := LogEntry{
        Message:   message,
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        LogLevel:  level,
    }
    body, _ := json.Marshal(entry)
    _, err := http.Post(apiURL, "application/json", bytes.NewBuffer(body))
    return err
}

// Usage
// sendLog("INFO", "Service started")
// sendLog("DEBUG", "Processing request")`}
              />
            </div>
          </div>
        </Section>

        {/* HTTP Streaming */}
        <Section
          title="Receiving Logs (HTTP Streaming)"
          description="Connect to the stream endpoint to receive logs in real-time."
        >
          <div className="space-y-4">
            <p className="text-muted-foreground text-sm">
              The dashboard automatically connects via Server-Sent Events (SSE) to receive logs.
              You can also connect programmatically:
            </p>
            <CodeBlock
              language="javascript"
              code={`const STREAM_ID = 'your-stream-id';
const eventSource = new EventSource(
  \`https://api.livelogs.dev/stream/\${STREAM_ID}/subscribe\`
);

eventSource.onmessage = (event) => {
  const log = JSON.parse(event.data);
  console.log(\`[\${log.log_level}] \${log.message}\`);
};

eventSource.onerror = (error) => {
  console.error('Stream connection error:', error);
};

// Close connection when done
// eventSource.close();`}
            />
          </div>
        </Section>

        {/* Best Practices */}
        <Section title="Best Practices">
          <div className="grid md:grid-cols-2 gap-4">
            <div className="bg-card/30 border border-border rounded-lg p-4 space-y-2">
              <div className="flex items-center gap-2">
                <Zap className="w-4 h-4 text-accent" />
                <h4 className="font-semibold text-foreground">Batch Logs</h4>
              </div>
              <p className="text-sm text-muted-foreground">
                For high-volume applications, batch multiple log entries into a single request to reduce overhead.
              </p>
            </div>
            <div className="bg-card/30 border border-border rounded-lg p-4 space-y-2">
              <div className="flex items-center gap-2">
                <Send className="w-4 h-4 text-accent" />
                <h4 className="font-semibold text-foreground">Async Logging</h4>
              </div>
              <p className="text-sm text-muted-foreground">
                Send logs asynchronously to avoid blocking your main application thread.
              </p>
            </div>
          </div>
        </Section>

        {/* CTA */}
        <div className="text-center py-8">
          <NavLink to="/">
            <Button variant="hero" size="lg" className="font-semibold">
              Start Streaming Logs
            </Button>
          </NavLink>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-border/50 py-8">
        <div className="max-w-4xl mx-auto px-6 text-center text-muted-foreground text-sm">
          © 2024 Live Logs. Built for developers.
        </div>
      </footer>
    </div>
  );
}
