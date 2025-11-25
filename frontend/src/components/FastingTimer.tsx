import { useState, useEffect } from "react";
import { Button } from "./ui/button";
import { Card } from "./ui/card";
import { Play, StopCircle } from "lucide-react";
import { cn } from "@/lib/utils";

interface FastingTimerProps {
  onStatusChange?: (isFasting: boolean) => void;
}

const FastingTimer = ({ onStatusChange }: FastingTimerProps) => {
  const [isFasting, setIsFasting] = useState(false);
  const [startTime, setStartTime] = useState<number | null>(null);
  const [elapsedSeconds, setElapsedSeconds] = useState(0);
  const targetHours = 16;

  useEffect(() => {
    if (isFasting && startTime) {
      const interval = setInterval(() => {
        setElapsedSeconds(Math.floor((Date.now() - startTime) / 1000));
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [isFasting, startTime]);

  const handleToggleFast = () => {
    if (!isFasting) {
      setStartTime(Date.now());
      setElapsedSeconds(0);
      setIsFasting(true);
      onStatusChange?.(true);
    } else {
      setIsFasting(false);
      setStartTime(null);
      setElapsedSeconds(0);
      onStatusChange?.(false);
    }
  };

  const hours = Math.floor(elapsedSeconds / 3600);
  const minutes = Math.floor((elapsedSeconds % 3600) / 60);

  const progressPercentage = Math.min((hours / targetHours) * 100, 100);

  return (
    <Card className="relative overflow-hidden bg-gradient-to-br from-card to-muted border-primary/20 animate-scale-in">
      {/* Animated background effects */}
      <div className="absolute inset-0 bg-gradient-to-r from-primary/10 via-secondary/10 to-primary/10 animate-shimmer"
        style={{ backgroundSize: '200% 100%' }} />

      {/* Pulsing glow orbs */}
      <div className="absolute top-0 left-0 w-32 h-32 bg-primary/20 rounded-full blur-3xl animate-glow-pulse" />
      <div className="absolute bottom-0 right-0 w-32 h-32 bg-secondary/20 rounded-full blur-3xl animate-glow-pulse"
        style={{ animationDelay: '1.5s' }} />

      <div className="relative p-8 flex flex-col items-center">
        {/* Circular Progress Ring */}
        <div className="relative w-48 h-48 mb-6 animate-float">
          {/* Outer glow ring */}
          <div className="absolute inset-0 rounded-full glow-primary opacity-50" />

          <svg className="w-full h-full transform -rotate-90">
            {/* Background circle */}
            <circle
              cx="96"
              cy="96"
              r="88"
              stroke="currentColor"
              strokeWidth="8"
              fill="none"
              className="text-muted"
            />
            {/* Progress circle */}
            <circle
              cx="96"
              cy="96"
              r="88"
              stroke="url(#gradient)"
              strokeWidth="8"
              fill="none"
              strokeDasharray={`${2 * Math.PI * 88}`}
              strokeDashoffset={`${2 * Math.PI * 88 * (1 - progressPercentage / 100)}`}
              className="transition-all duration-1000 ease-out"
              strokeLinecap="round"
            />
            <defs>
              <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stopColor="hsl(var(--primary))" />
                <stop offset="100%" stopColor="hsl(var(--secondary))" />
              </linearGradient>
            </defs>
          </svg>

          {/* Center content */}
          <div className="absolute inset-0 flex flex-col items-center justify-center">
            <div className="text-4xl font-bold font-mono text-shimmer transition-all duration-300">
              {String(hours).padStart(2, '0')}:{String(minutes).padStart(2, '0')}
            </div>
            <div className="text-sm text-muted-foreground mt-1 animate-pulse">
              {isFasting ? 'Fasting' : 'Not Fasting'}
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="flex gap-8 mb-6 text-center">
          <div>
            <div className="text-2xl font-bold text-foreground">{hours}h</div>
            <div className="text-xs text-muted-foreground">Elapsed</div>
          </div>
          <div className="w-px bg-border" />
          <div>
            <div className="text-2xl font-bold text-primary">{targetHours}h</div>
            <div className="text-xs text-muted-foreground">Target</div>
          </div>
        </div>

        {/* Action Button */}
        <Button
          onClick={handleToggleFast}
          size="lg"
          className={cn(
            "w-full max-w-xs font-semibold transition-all duration-300 hover:scale-105 active:scale-95",
            isFasting
              ? "bg-destructive hover:bg-destructive/90 hover:shadow-lg hover:shadow-destructive/50"
              : "bg-gradient-to-r from-primary to-secondary hover:opacity-90 glow-primary"
          )}
        >
          {isFasting ? (
            <>
              <StopCircle className="mr-2 h-5 w-5" />
              End Fast
            </>
          ) : (
            <>
              <Play className="mr-2 h-5 w-5" />
              Start Fast
            </>
          )}
        </Button>
      </div>
    </Card>
  );
};

export default FastingTimer;
