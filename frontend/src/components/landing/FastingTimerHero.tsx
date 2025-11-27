import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

const phases = [
  { name: "Digestion", icon: "🍽️", hours: 4, color: "text-muted-foreground" },
  { name: "Glycogen Depletion", icon: "⚡", hours: 12, color: "text-secondary" },
  { name: "Ketosis", icon: "🔥", hours: 16, color: "text-primary" },
  { name: "Autophagy", icon: "✨", hours: 24, color: "text-accent" },
];

const getPhase = (hours: number) => {
  if (hours < 4) return phases[0];
  if (hours < 12) return phases[1];
  if (hours < 16) return phases[2];
  return phases[3];
};

const getRecoveryAmount = (hours: number) => {
  // $0.125 per hour, max $2.00 at 16 hours
  return Math.min(hours * 0.125, 2.0);
};

export const FastingTimerHero = () => {
  const [elapsed, setElapsed] = useState(14 * 3600 + 23 * 60 + 45); // Start at 14:23:45
  const [showMilestone, setShowMilestone] = useState(false);
  
  useEffect(() => {
    const interval = setInterval(() => {
      setElapsed((prev) => {
        const newVal = prev + 1;
        // Check for hour milestones
        if (newVal % 3600 === 0) {
          setShowMilestone(true);
          setTimeout(() => setShowMilestone(false), 1500);
        }
        return newVal;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  const hours = Math.floor(elapsed / 3600);
  const minutes = Math.floor((elapsed % 3600) / 60);
  const seconds = elapsed % 60;
  
  const progress = Math.min((elapsed / (24 * 3600)) * 100, 100);
  const circumference = 2 * Math.PI * 130;
  const strokeDashoffset = circumference - (progress / 100) * circumference;
  
  const currentPhase = getPhase(hours);
  const recoveryAmount = getRecoveryAmount(hours);
  const phaseProgress = Math.min(((hours % currentPhase.hours) / currentPhase.hours) * 100, 100);

  return (
    <div className="relative">
      {/* Money recovery counter */}
      <motion.div
        className="text-center mb-4"
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
      >
        <span className="text-muted-foreground text-sm">Recovering</span>
        <motion.div
          className="font-display text-3xl font-bold text-primary"
          key={recoveryAmount.toFixed(2)}
          initial={{ y: -10, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
        >
          ${recoveryAmount.toFixed(2)}
        </motion.div>
      </motion.div>

      {/* Main timer circle */}
      <div className="relative w-[300px] h-[300px] mx-auto">
        {/* Pulsing glow effect */}
        <motion.div
          className="absolute inset-0 rounded-full"
          animate={{
            boxShadow: [
              "0 0 30px hsl(var(--primary) / 0.3)",
              "0 0 60px hsl(var(--primary) / 0.5)",
              "0 0 30px hsl(var(--primary) / 0.3)",
            ],
          }}
          transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
        />

        <svg width="300" height="300" className="transform -rotate-90">
          {/* Background ring */}
          <circle
            cx="150"
            cy="150"
            r="130"
            fill="none"
            stroke="hsl(var(--muted))"
            strokeWidth="12"
            opacity="0.3"
          />
          
          {/* Progress ring with gradient */}
          <motion.circle
            cx="150"
            cy="150"
            r="130"
            fill="none"
            stroke="url(#timerGradientHero)"
            strokeWidth="12"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={strokeDashoffset}
            initial={{ strokeDashoffset: circumference }}
            animate={{ strokeDashoffset }}
            transition={{ duration: 0.5, ease: "easeOut" }}
          />
          
          <defs>
            <linearGradient id="timerGradientHero" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="hsl(var(--secondary))" />
              <stop offset="100%" stopColor="hsl(var(--primary))" />
            </linearGradient>
          </defs>
        </svg>

        {/* Center content */}
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          {/* Time display */}
          <motion.div
            className="font-mono text-5xl font-bold text-foreground tabular-nums tracking-tight"
            animate={showMilestone ? { scale: [1, 1.1, 1] } : {}}
            transition={{ duration: 0.3 }}
          >
            {String(hours).padStart(2, "0")}:
            {String(minutes).padStart(2, "0")}:
            {String(seconds).padStart(2, "0")}
          </motion.div>

          {/* Phase indicator */}
          <AnimatePresence mode="wait">
            <motion.div
              key={currentPhase.name}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              className={`flex items-center gap-2 mt-3 ${currentPhase.color}`}
            >
              <span className="text-xl">{currentPhase.icon}</span>
              <span className="text-sm font-semibold">{currentPhase.name}</span>
            </motion.div>
          </AnimatePresence>

          {/* Phase progress bar */}
          <div className="w-32 h-1 bg-muted rounded-full mt-2 overflow-hidden">
            <motion.div
              className="h-full bg-gradient-gold rounded-full"
              initial={{ width: 0 }}
              animate={{ width: `${phaseProgress}%` }}
              transition={{ duration: 0.5 }}
            />
          </div>
        </div>

        {/* Milestone celebration effect */}
        <AnimatePresence>
          {showMilestone && (
            <motion.div
              initial={{ scale: 0.5, opacity: 0 }}
              animate={{ scale: 1.5, opacity: 1 }}
              exit={{ scale: 2, opacity: 0 }}
              className="absolute inset-0 flex items-center justify-center"
            >
              <div className="text-6xl">🎉</div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>

      {/* Floating +$0.XX notifications */}
      <FloatingMoneyNotifications elapsed={elapsed} />
    </div>
  );
};

const FloatingMoneyNotifications = ({ elapsed }: { elapsed: number }) => {
  const [notifications, setNotifications] = useState<{ id: number; amount: string }[]>([]);

  useEffect(() => {
    // Add notification every 30 seconds
    if (elapsed % 30 === 0 && elapsed > 0) {
      const id = Date.now();
      setNotifications((prev) => [...prev, { id, amount: "+$0.06" }]);
      setTimeout(() => {
        setNotifications((prev) => prev.filter((n) => n.id !== id));
      }, 2000);
    }
  }, [elapsed]);

  return (
    <div className="absolute top-0 right-0 pointer-events-none">
      <AnimatePresence>
        {notifications.map((notification) => (
          <motion.div
            key={notification.id}
            initial={{ opacity: 0, y: 20, x: 0 }}
            animate={{ opacity: 1, y: -30, x: 10 }}
            exit={{ opacity: 0, y: -60 }}
            className="text-primary font-bold text-lg"
          >
            {notification.amount}
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  );
};
