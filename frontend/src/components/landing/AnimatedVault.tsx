import { useState, useEffect } from "react";
import { motion } from "framer-motion";

export const AnimatedVault = () => {
  const [progress, setProgress] = useState(60);

  useEffect(() => {
    const interval = setInterval(() => {
      setProgress((prev) => (prev >= 100 ? 0 : prev + 1));
    }, 100);
    return () => clearInterval(interval);
  }, []);

  const circumference = 2 * Math.PI * 28;
  const strokeDashoffset = circumference - (progress / 100) * circumference;

  return (
    <div className="relative w-20 h-20">
      {/* Pulsing glow */}
      <motion.div
        className="absolute inset-0 rounded-full"
        animate={{
          boxShadow: [
            "0 0 15px hsl(var(--primary) / 0.3)",
            "0 0 25px hsl(var(--primary) / 0.5)",
            "0 0 15px hsl(var(--primary) / 0.3)",
          ],
        }}
        transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
      />

      {/* Timer circle */}
      <svg width="80" height="80" className="transform -rotate-90">
        {/* Background ring */}
        <circle
          cx="40"
          cy="40"
          r="28"
          fill="none"
          stroke="hsl(var(--muted))"
          strokeWidth="6"
          opacity="0.3"
        />
        {/* Progress ring */}
        <motion.circle
          cx="40"
          cy="40"
          r="28"
          fill="none"
          stroke="url(#miniTimerGradient)"
          strokeWidth="6"
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={strokeDashoffset}
          transition={{ duration: 0.1 }}
        />
        <defs>
          <linearGradient id="miniTimerGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="hsl(var(--secondary))" />
            <stop offset="100%" stopColor="hsl(var(--primary))" />
          </linearGradient>
        </defs>
      </svg>

      {/* Center icon */}
      <div className="absolute inset-0 flex items-center justify-center">
        <motion.span
          className="text-2xl"
          animate={{ scale: [1, 1.1, 1] }}
          transition={{ duration: 2, repeat: Infinity }}
        >
          ⏱️
        </motion.span>
      </div>
    </div>
  );
};
