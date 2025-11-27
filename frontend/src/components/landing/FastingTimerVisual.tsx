import { useEffect, useState } from "react";

export const FastingTimerVisual = () => {
  const targetHours = 18;
  const totalSeconds = targetHours * 60 * 60;
  
  const [elapsed, setElapsed] = useState(16 * 60 * 60 + 42 * 60 + 15);
  
  useEffect(() => {
    const interval = setInterval(() => {
      setElapsed((prev) => (prev + 1) % totalSeconds);
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  const hours = Math.floor(elapsed / 3600);
  const minutes = Math.floor((elapsed % 3600) / 60);
  
  const progress = (elapsed / totalSeconds) * 100;
  const circumference = 2 * Math.PI * 180;
  const strokeDashoffset = circumference - (progress / 100) * circumference;

  return (
    <div className="absolute inset-0 pointer-events-none overflow-hidden">
      {/* Animated gradient orbs - pulsing light effect */}
      <div className="absolute top-1/4 left-1/4 w-[600px] h-[600px] rounded-full bg-primary/20 blur-[150px] animate-pulse-slow" />
      <div className="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] rounded-full bg-secondary/15 blur-[120px] animate-pulse-slow" style={{ animationDelay: "1s" }} />
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] rounded-full bg-purple/10 blur-[180px] animate-pulse-slow" style={{ animationDelay: "2s" }} />
      
      {/* Rotating light rays */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[1000px] h-[1000px] animate-spin-slow opacity-20">
        {[...Array(12)].map((_, i) => (
          <div
            key={i}
            className="absolute top-1/2 left-1/2 w-1 h-[400px] bg-gradient-to-t from-transparent via-primary/30 to-transparent origin-bottom"
            style={{ transform: `rotate(${i * 30}deg) translateX(-50%)` }}
          />
        ))}
      </div>

      {/* Concentric pulse rings */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
        {[...Array(4)].map((_, i) => (
          <div
            key={i}
            className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full border border-primary/20 animate-ripple"
            style={{
              width: `${300 + i * 150}px`,
              height: `${300 + i * 150}px`,
              animationDelay: `${i * 0.8}s`,
            }}
          />
        ))}
      </div>

      {/* Main timer circle */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
        <svg width="500" height="500" className="transform -rotate-90 opacity-30">
          <circle
            cx="250"
            cy="250"
            r="180"
            fill="none"
            stroke="hsl(var(--muted))"
            strokeWidth="2"
            opacity="0.3"
          />
          <circle
            cx="250"
            cy="250"
            r="180"
            fill="none"
            stroke="url(#timerGradient)"
            strokeWidth="6"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={strokeDashoffset}
            className="transition-all duration-1000 ease-linear"
            style={{ filter: "drop-shadow(0 0 20px hsl(var(--primary)))" }}
          />
          <defs>
            <linearGradient id="timerGradient" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="hsl(var(--primary))" />
              <stop offset="100%" stopColor="hsl(var(--secondary))" />
            </linearGradient>
          </defs>
        </svg>
        
        {/* Timer display */}
        <div className="absolute inset-0 flex flex-col items-center justify-center opacity-25">
          <div className="font-display text-8xl font-bold text-foreground tabular-nums tracking-tight">
            {String(hours).padStart(2, "0")}:{String(minutes).padStart(2, "0")}
          </div>
          <div className="flex items-center gap-2 mt-4">
            <div className="w-3 h-3 rounded-full bg-secondary animate-pulse" />
            <span className="text-secondary text-lg font-semibold tracking-widest uppercase">
              Fasting
            </span>
          </div>
        </div>
      </div>

      {/* Floating particles */}
      <div className="absolute inset-0">
        {[...Array(20)].map((_, i) => (
          <div
            key={i}
            className="absolute w-1 h-1 rounded-full bg-primary/40 animate-float-particle"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 5}s`,
              animationDuration: `${4 + Math.random() * 4}s`,
            }}
          />
        ))}
      </div>

      {/* Corner glow accents */}
      <div className="absolute top-0 left-0 w-96 h-96 bg-gradient-radial from-primary/10 to-transparent blur-3xl" />
      <div className="absolute bottom-0 right-0 w-96 h-96 bg-gradient-radial from-secondary/10 to-transparent blur-3xl" />
    </div>
  );
};