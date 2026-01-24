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
        </div>
      </header>

      {/* Content */}
      <main className="max-w-4xl mx-auto px-6 py-12 space-y-16">
        {/* Hero */}
        <div className="space-y-4">
          <h1 className="text-4xl md:text-5xl font-bold font-display">
            <span className="text">Documentation</span>
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl">
            Learn how to integrate Live Logs into your application and start streaming logs in minutes.
          </p>
          <p className="text-xl text-muted-foreground max-w-2xl">
            Please note that uplog is currently under active development and this page may change
          </p>
        </div>

        {/* Quick Start */}
        <Section
          title="Quick Start"
          description="Get up and running with Live Logs in under 2 minutes."
        >
          <div className="space-y-6">
            <div className="flex items-start gap-4">
              <div className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0 mt-1">
                <span className="text-accent font-bold text-sm">1</span>
              </div>
              <div className="space-y-3 flex-1">
                <h3 className="font-semibold text-foreground">Installation (Both x86 and ARM are supported)</h3>
                <p className="text-muted-foreground text-sm">
                  Linux / Macos
                </p>
                <CodeBlock
                  language="bash"
                  code={`curl -fsSL https://uplog.live/install.sh | sh`}
                />
                <p className="text-muted-foreground text-sm">
                  Download for Windows from{" "}
                  <a
                    href="https://github.com/SOORAJTS2001/uplog/releases/latest"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary underline hover:text-primary/80"
                  >
                    here
                  </a>
                </p>

              </div>
            </div>

            <div className="flex items-start gap-4">
              <div className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0 mt-1">
                <span className="text-accent font-bold text-sm">2</span>
              </div>
              <div className="space-y-3 flex-1">
                <h3 className="font-semibold text-foreground">
                  Usage
                </h3>
                <CodeBlock
                  language="bash"
                  code={`$ uplog <executable>`}
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
                  Open your dashboard URL from uplog and watch logs appear instantly as they're sent.
                </p>
              </div>
            </div>
          </div>
        </Section>
        <Section title="Advanced Usage">
  <p className="text-sm text-muted-foreground">
    Set batch upload size
  </p>
  <CodeBlock
    language="bash"
    code={`uplog --poll <batch_size> <executable>`}
  />

  <p className="text-sm text-muted-foreground mt-4">
    Tag a session
  </p>
  <CodeBlock
    language="bash"
    code={`uplog --tag <tag> <executable>`}
  />

  <p className="text-sm text-muted-foreground mt-4">
    You can tag and batch at the same time
  </p>
  <CodeBlock
    language="bash"
    code={`uplog --poll <batch_size> --tag <tag> <executable>`}
  />

  <p className="text-sm text-muted-foreground mt-4">
    List all recorded sessions
  </p>
  <CodeBlock
    language="bash"
    code={`uplog list`}
  />

  <p className="text-sm text-muted-foreground mt-4">
    Delete all recorded sessions
  </p>
  <CodeBlock
    language="bash"
    code={`uplog purge`}
  />

  <p className="text-sm text-muted-foreground mt-4">
    Delete a single session
  </p>
  <CodeBlock
    language="bash"
    code={`uplog delete <session_id>`}
  />
</Section>

        {/* Log Format */}
        <Section
          title="Log Levels"
        >
            <div className="bg-card/30 border border-border rounded-lg p-4">
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
        </Section>


      </main>

      {/* Footer */}
      <footer className="border-t border-border/50 py-8">
        <div className="max-w-4xl mx-auto px-6 text-center text-muted-foreground text-sm">
          Built with ❤️ from Kerala, India
        </div>
      </footer>
    </div>
  );
}
