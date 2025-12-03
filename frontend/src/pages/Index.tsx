import { Button } from '@/components/ui/button';
import { useNavigate } from 'react-router-dom';
import { NavLink } from '@/components/NavLink';
import { Terminal, Zap, Shield, Link2, ArrowRight, Play, Clock, Filter, BookOpen } from 'lucide-react';

const Index = () => {
  const navigate = useNavigate();

  const handleCreateStream = () => {
    const streamId = `stream_${Date.now().toString(36)}_${Math.random().toString(36).substr(2, 9)}`;
    navigate(`/logs/${streamId}`);
  };

  return (
    <div className="min-h-screen bg-background relative overflow-hidden">
      {/* Dotted grid background */}
      <div className="absolute inset-0 dotted-grid opacity-60" />

      {/* Gradient orbs */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl animate-float" />
      <div className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-primary/5 rounded-full blur-3xl animate-float" style={{ animationDelay: '-3s' }} />

      {/* Header */}
      <header className="relative z-10 border-b border-border/30">
        <div className="container mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-lg bg-primary flex items-center justify-center">
              <Terminal className="w-5 h-5 text-primary-foreground" />
            </div>
            <span className="font-display font-semibold text-xl tracking-tight">Uplog</span>
          </div>
          <nav className="hidden md:flex items-center gap-8">
            <a href="#features" className="text-sm text-muted-foreground hover:text-foreground transition-colors">Features</a>
            <a href="#how-it-works" className="text-sm text-muted-foreground hover:text-foreground transition-colors">How it works</a>
            <a href="https://github.com/SOORAJTS2001/uplog" className="text-sm text-muted-foreground hover:text-foreground transition-colors">Github</a>

            <NavLink to="/docs" className="text-sm text-muted-foreground hover:text-foreground transition-colors">Docs</NavLink>
            <Button variant="ghost" size="sm" onClick={handleCreateStream}>
              Get Started
            </Button>
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <main className="relative z-10">
        <section className="container mx-auto px-6 pt-24 pb-32">
          <div className="max-w-4xl mx-auto text-center">
            {/* Announcement badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-secondary border border-border text-sm mb-10 animate-fade-in">
              <span className="w-2 h-2 rounded-full bg-primary animate-pulse" />
              <span className="text-muted-foreground">In Development</span>
              <ArrowRight className="w-3.5 h-3.5 text-muted-foreground" />
            </div>

            {/* Main heading */}
            <h1 className="font-display text-5xl md:text-7xl font-bold tracking-tight mb-8 animate-slide-up">
              Debug faster with{' '}
              <span className="text-primary neon-text">Live Logs</span>
            </h1>

            {/* Description */}
            <p className="text-lg md:text-xl text-muted-foreground mb-12 max-w-2xl mx-auto animate-slide-up" style={{ animationDelay: '0.1s' }}>
              Anonymous log monitoring that works in milli-seconds.
              No signup, no dependencies, and no code rewrites -
              just plug in the CLI and watch your logs stream live.
            </p>

            {/* CTAs */}
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 animate-slide-up" style={{ animationDelay: '0.2s' }}>
              <Button
                size="lg"
                onClick={handleCreateStream}
                className="h-12 px-8 text-base font-medium neon-glow"
              >
                Try Sample
                <ArrowRight className="w-4 h-4 ml-1" />
              </Button>
              <Button
                variant="outline"
                size="lg"
                className="h-12 px-8 text-base font-medium"
                onClick={() => document.getElementById('how-it-works')?.scrollIntoView({ behavior: 'smooth' })}
              >
                See how it works
              </Button>
            </div>
          </div>
        </section>
        <section className="container mx-auto px-6 pb-20">
          <div className="max-w-5xl mx-auto">
            <div className="relative rounded-xl border border-border bg-card/50 backdrop-blur-sm overflow-hidden animate-slide-up" style={{ animationDelay: '0.3s' }}>
              {/* Window chrome */}
              <div className="flex items-center gap-2 px-4 py-3 border-b border-border bg-secondary/50">
                <div className="flex gap-1.5">
                  <div className="w-3 h-3 rounded-full bg-destructive/60" />
                  <div className="w-3 h-3 rounded-full bg-log-warn/60" />
                  <div className="w-3 h-3 rounded-full bg-primary/60" />
                </div>
                <div className="flex-1 text-center">
                  <span className="text-xs font-mono text-muted-foreground">demo@demo-cli</span>
                </div>
              </div>
              <div className="p-4 font-mono text-sm space-y-2">
                $ uplog run python main.py
                </div>
              </div>
              </div>
        </section>

        {/* Demo preview */}
        <section className="container mx-auto px-6 pb-32">
          <div className="max-w-5xl mx-auto">
            <div className="relative rounded-xl border border-border bg-card/50 backdrop-blur-sm overflow-hidden animate-slide-up" style={{ animationDelay: '0.3s' }}>
              {/* Window chrome */}
              <div className="flex items-center gap-2 px-4 py-3 border-b border-border bg-secondary/50">
                <div className="flex gap-1.5">
                  <div className="w-3 h-3 rounded-full bg-destructive/60" />
                  <div className="w-3 h-3 rounded-full bg-log-warn/60" />
                  <div className="w-3 h-3 rounded-full bg-primary/60" />
                </div>
                <div className="flex-1 text-center">
                  <span className="text-xs font-mono text-muted-foreground">logstream.app/logs/stream_abc123</span>
                </div>
              </div>
              {/* Mock log content */}
              <div className="p-4 font-mono text-sm space-y-2">
                <LogLine time="14:32:01.234" level="INFO" message="Hello world!" />
                <div className="flex items-center gap-2 pt-2 text-muted-foreground">
                  <div className="w-2 h-2 rounded-full bg-primary animate-pulse" />
                  <span className="text-xs">Streaming...</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Features */}
        <section id="features" className="container mx-auto px-6 py-24 border-t border-border/30">
          <div className="max-w-5xl mx-auto">
            <div className="text-center mb-16">
              <h2 className="font-display text-3xl md:text-4xl font-bold mb-4">
                Everything you need to debug
              </h2>
              <p className="text-muted-foreground text-lg max-w-xl mx-auto">
                A minimal, powerful log viewer built for speed and simplicity.
              </p>
            </div>

            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
              <FeatureCard
                icon={<Zap className="w-5 h-5" />}
                title="Instant Streams"
                description="Create a log stream in seconds. No accounts, no configuration, just click and go."
              />
              <FeatureCard
                icon={<Link2 className="w-5 h-5" />}
                title="Shareable Links"
                description="Share your unique stream URL with teammates. Anyone with the link can watch."
              />
              <FeatureCard
                icon={<Play className="w-5 h-5" />}
                title="HTTP Streaming"
                description="Real-time updates via server-sent events. No polling, no WebSocket complexity."
              />
              <FeatureCard
                icon={<Filter className="w-5 h-5" />}
                title="Smart Filtering"
                description="Filter logs by level—ERROR, WARN, INFO, DEBUG. Find issues instantly."
              />
              <FeatureCard
                icon={<Clock className="w-5 h-5" />}
                title="Auto-scroll"
                description="Logs auto-scroll as they arrive. Pause anytime to investigate."
              />
              <FeatureCard
                icon={<Shield className="w-5 h-5" />}
                title="No Sign-up"
                description="Completely anonymous. No accounts, no tracking, no hassle."
              />
            </div>
          </div>
        </section>

        {/* How it works */}
        <section id="how-it-works" className="container mx-auto px-6 py-24 border-t border-border/30">
          <div className="max-w-4xl mx-auto">
            <div className="text-center mb-16">
              <h2 className="font-display text-3xl md:text-4xl font-bold mb-4">
                How it works
              </h2>
              <p className="text-muted-foreground text-lg">
                Three steps to start streaming logs.
              </p>
            </div>

            <div className="grid md:grid-cols-3 gap-8">
              <StepCard
                number="01"
                title="Install our cli"
                description="By just a simple curl command"
              />
              <StepCard
                number="02"
                title="Run you code with uplog"
                description="A ligthweight background process will spawn to monitor and send logs"
              />
              <StepCard
                number="03"
                title="Watch live"
                description="Open your link from uplog and watch logs appear in real-time."
              />
            </div>

            <div className="text-center mt-16">
              <Button
                size="lg"
                onClick={handleCreateStream}
                className="h-12 px-8 text-base font-medium neon-glow"
              >
                Try sample
                <ArrowRight className="w-4 h-4 ml-1" />
              </Button>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="relative z-10 border-t border-border/30 py-8">
        <div className="container mx-auto px-6">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-2">
              <Terminal className="w-4 h-4 text-primary" />
              <span className="text-sm font-medium">Uplog</span>
            </div>
            <p className="text-sm text-muted-foreground">
              Built for developers who need fast, simple log monitoring.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
};

function LogLine({ time, level, message }: { time: string; level: string; message: string }) {
  const levelColors: Record<string, string> = {
    ERROR: 'text-log-error',
    WARN: 'text-log-warn',
    INFO: 'text-log-info',
    DEBUG: 'text-log-debug',
  };

  return (
    <div className="flex items-start gap-4 py-1 hover:bg-muted/30 px-2 -mx-2 rounded transition-colors">
      <span className="text-muted-foreground shrink-0">{time}</span>
      <span className={`${levelColors[level]} font-medium shrink-0 w-12`}>{level}</span>
      <span className="text-foreground/90">{message}</span>
    </div>
  );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="p-6 rounded-xl bg-card/50 border border-border/50 hover:border-border transition-colors">
      <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center text-primary mb-4">
        {icon}
      </div>
      <h3 className="font-display font-semibold text-lg mb-2">{title}</h3>
      <p className="text-sm text-muted-foreground leading-relaxed">{description}</p>
    </div>
  );
}

function StepCard({ number, title, description }: { number: string; title: string; description: string }) {
  return (
    <div className="text-center">
      <div className="inline-flex items-center justify-center w-14 h-14 rounded-full border-2 border-primary/30 text-primary font-display font-bold text-xl mb-4">
        {number}
      </div>
      <h3 className="font-display font-semibold text-lg mb-2">{title}</h3>
      <p className="text-sm text-muted-foreground">{description}</p>
    </div>
  );
}

export default Index;
