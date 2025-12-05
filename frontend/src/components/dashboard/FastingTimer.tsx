import { useEffect, useRef, useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Flame, Zap, Sparkles, Utensils, Pause, Square, Trophy } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import confetti from 'canvas-confetti';

interface FastingPhase {
  id: string;
  name: string;
  icon: React.ElementType;
  startHour: number;
  endHour: number;
  color: string;
  bgColor: string;
  glowColor: string;
  confettiColors: string[];
  description: string;
  celebration: string;
}

const fastingPhases: FastingPhase[] = [
  {
    id: 'fed',
    name: 'Fed',
    icon: Utensils,
    startHour: 0,
    endHour: 4,
    color: 'text-blue-400',
    bgColor: 'bg-blue-500/20',
    glowColor: 'shadow-blue-500/50',
    confettiColors: ['#60A5FA', '#3B82F6'],
    description: 'Digesting food, insulin elevated',
    celebration: 'Fast started! 🚀',
  },
  {
    id: 'fat-burning',
    name: 'Fat Burning',
    icon: Flame,
    startHour: 4,
    endHour: 12,
    color: 'text-orange-400',
    bgColor: 'bg-orange-500/20',
    glowColor: 'shadow-orange-500/50',
    confettiColors: ['#FB923C', '#F97316', '#EA580C'],
    description: 'Glycogen depleting, fat burning begins',
    celebration: 'Fat burning activated! 🔥',
  },
  {
    id: 'ketosis',
    name: 'Ketosis',
    icon: Zap,
    startHour: 12,
    endHour: 18,
    color: 'text-secondary',
    bgColor: 'bg-secondary/20',
    glowColor: 'shadow-emerald-500/50',
    confettiColors: ['#34D399', '#10B981', '#059669'],
    description: 'Full ketone production, mental clarity',
    celebration: 'Ketosis achieved! ⚡',
  },
  {
    id: 'autophagy',
    name: 'Autophagy',
    icon: Sparkles,
    startHour: 18,
    endHour: 24,
    color: 'text-purple-400',
    bgColor: 'bg-purple-500/20',
    glowColor: 'shadow-purple-500/50',
    confettiColors: ['#C084FC', '#A855F7', '#9333EA'],
    description: 'Cellular cleanup and regeneration',
    celebration: 'Autophagy unlocked! ✨',
  },
];

interface FastingTimerProps {
  isActive: boolean;
  elapsedSeconds: number;
  targetHours: number;
  onStart: () => void;
  onPause: () => void;
  onStop: () => void;
  insight?: string | null;
}

const FastingTimer = ({
  isActive,
  elapsedSeconds,
  targetHours,
  onStart,
  onPause,
  onStop,
  insight,
}: FastingTimerProps) => {
  const elapsedHours = elapsedSeconds / 3600;
  const progress = Math.min((elapsedHours / targetHours) * 100, 100);
  const previousPhaseRef = useRef<string | null>(null);
  const [showCelebration, setShowCelebration] = useState(false);
  const [celebrationMessage, setCelebrationMessage] = useState('');

  // Find current phase
  const currentPhase = fastingPhases.find(
    (phase) => elapsedHours >= phase.startHour && elapsedHours < phase.endHour
  ) || fastingPhases[fastingPhases.length - 1];

  // Detect phase transitions and trigger celebration
  useEffect(() => {
    if (previousPhaseRef.current && previousPhaseRef.current !== currentPhase.id && isActive) {
      // Phase changed! Trigger celebration
      triggerCelebration(currentPhase);
    }
    previousPhaseRef.current = currentPhase.id;
  }, [currentPhase.id, currentPhase, isActive]);

  const triggerCelebration = (phase: FastingPhase) => {
    setCelebrationMessage(phase.celebration);
    setShowCelebration(true);

    // Fire confetti
    const duration = 3000;
    const animationEnd = Date.now() + duration;
    const defaults = { startVelocity: 30, spread: 360, ticks: 60, zIndex: 100 };

    const randomInRange = (min: number, max: number) => Math.random() * (max - min) + min;

    const interval = setInterval(() => {
      const timeLeft = animationEnd - Date.now();
      if (timeLeft <= 0) {
        clearInterval(interval);
        return;
      }

      const particleCount = 50 * (timeLeft / duration);

      confetti({
        ...defaults,
        particleCount,
        origin: { x: randomInRange(0.1, 0.3), y: Math.random() - 0.2 },
        colors: phase.confettiColors,
      });
      confetti({
        ...defaults,
        particleCount,
        origin: { x: randomInRange(0.7, 0.9), y: Math.random() - 0.2 },
        colors: phase.confettiColors,
      });
    }, 250);

    // Hide celebration message after delay
    setTimeout(() => setShowCelebration(false), 4000);
  };

  // Calculate time to next phase
  const nextPhase = fastingPhases.find((phase) => phase.startHour > elapsedHours);
  const hoursToNextPhase = nextPhase ? nextPhase.startHour - elapsedHours : 0;

  const formatTime = (seconds: number) => {
    const hrs = Math.floor(seconds / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    return `${hrs.toString().padStart(2, '0')}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const formatHoursMinutes = (hours: number) => {
    const h = Math.floor(hours);
    const m = Math.floor((hours - h) * 60);
    if (h === 0) return `${m}m`;
    if (m === 0) return `${h}h`;
    return `${h}h ${m}m`;
  };

  // SVG circle calculations
  const size = 200;
  const strokeWidth = 8;
  const radius = (size - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const strokeDashoffset = circumference - (progress / 100) * circumference;

  return (
    <div className="space-y-6 relative">
      {/* Celebration Overlay */}
      <AnimatePresence>
        {showCelebration && (
          <motion.div
            initial={{ opacity: 0, scale: 0.8, y: -20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.8, y: -20 }}
            className="absolute -top-4 left-1/2 -translate-x-1/2 z-50"
          >
            <div className={cn(
              'px-6 py-3 rounded-2xl border-2 flex items-center gap-3',
              'bg-card/95 backdrop-blur-sm shadow-2xl',
              currentPhase.bgColor,
              `border-${currentPhase.color.replace('text-', '')}`
            )}>
              <motion.div
                animate={{ rotate: [0, 15, -15, 0], scale: [1, 1.2, 1] }}
                transition={{ duration: 0.5, repeat: 3 }}
              >
                <Trophy className={cn('w-6 h-6', currentPhase.color)} />
              </motion.div>
              <span className={cn('font-bold text-lg', currentPhase.color)}>
                {celebrationMessage}
              </span>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Circular Timer */}
      <div className="flex flex-col items-center">
        <div className="relative">
          {/* Glow effect on phase change */}
          <motion.div
            className={cn(
              'absolute inset-0 rounded-full blur-xl opacity-0',
              currentPhase.bgColor.replace('/20', '/40')
            )}
            animate={{
              opacity: showCelebration ? [0, 0.8, 0] : 0,
              scale: showCelebration ? [1, 1.3, 1] : 1,
            }}
            transition={{ duration: 1.5, repeat: showCelebration ? 2 : 0 }}
          />

          <svg width={size} height={size} className="transform -rotate-90 relative z-10">
            {/* Background circle */}
            <circle
              cx={size / 2}
              cy={size / 2}
              r={radius}
              fill="none"
              stroke="currentColor"
              strokeWidth={strokeWidth}
              className="text-muted/30"
            />
            {/* Progress circle */}
            <motion.circle
              cx={size / 2}
              cy={size / 2}
              r={radius}
              fill="none"
              stroke="url(#timerGradient)"
              strokeWidth={strokeWidth}
              strokeLinecap="round"
              strokeDasharray={circumference}
              initial={{ strokeDashoffset: circumference }}
              animate={{ strokeDashoffset }}
              transition={{ duration: 0.5, ease: 'easeOut' }}
            />
            <defs>
              <linearGradient id="timerGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" stopColor="hsl(var(--primary))" />
                <stop offset="100%" stopColor="hsl(var(--secondary))" />
              </linearGradient>
            </defs>
          </svg>

          {/* Center content */}
          <div className="absolute inset-0 flex flex-col items-center justify-center z-10">
            <motion.div
              key={currentPhase.id}
              className={cn('p-2 rounded-full mb-2', currentPhase.bgColor)}
              initial={{ scale: 0.5, opacity: 0 }}
              animate={{
                scale: isActive ? [1, 1.1, 1] : 1,
                opacity: 1,
              }}
              transition={{
                scale: { repeat: isActive ? Infinity : 0, duration: 2 },
                opacity: { duration: 0.3 }
              }}
            >
              <currentPhase.icon className={cn('w-6 h-6', currentPhase.color)} />
            </motion.div>
            <motion.div
              className="text-3xl font-bold font-mono text-foreground"
              key={formatTime(elapsedSeconds)}
            >
              {formatTime(elapsedSeconds)}
            </motion.div>
            <div className="text-xs text-muted-foreground mt-1">
              {isActive ? 'Fasting' : 'Paused'} • Goal: {targetHours}h
            </div>
          </div>
        </div>

        {/* Current Phase Badge with animation */}
        <AnimatePresence mode="wait">
          <motion.div
            key={currentPhase.id}
            className={cn(
              'mt-4 px-4 py-2 rounded-full border flex items-center gap-2',
              currentPhase.bgColor,
              showCelebration && 'shadow-lg',
              showCelebration && currentPhase.glowColor
            )}
            initial={{ opacity: 0, y: 20, scale: 0.8 }}
            animate={{
              opacity: 1,
              y: 0,
              scale: 1,
              boxShadow: showCelebration
                ? ['0 0 0px rgba(0,0,0,0)', '0 0 30px rgba(0,0,0,0.3)', '0 0 0px rgba(0,0,0,0)']
                : '0 0 0px rgba(0,0,0,0)'
            }}
            exit={{ opacity: 0, y: -20, scale: 0.8 }}
            transition={{ duration: 0.4 }}
          >
            <motion.div
              animate={showCelebration ? { rotate: [0, 360] } : {}}
              transition={{ duration: 0.5 }}
            >
              <currentPhase.icon className={cn('w-4 h-4', currentPhase.color)} />
            </motion.div>
            <span className={cn('font-medium text-sm', currentPhase.color)}>
              {currentPhase.name}
            </span>
          </motion.div>
        </AnimatePresence>

        <motion.p
          key={insight || currentPhase.description}
          className="text-sm text-muted-foreground mt-2 text-center max-w-xs"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.2 }}
        >
          {insight ? (
            <span className="text-primary font-medium flex items-center gap-1 justify-center">
              <Sparkles className="w-3 h-3" />
              {insight}
            </span>
          ) : (
            currentPhase.description
          )}
        </motion.p>

        {nextPhase && isActive && (
          <motion.p
            className="text-xs text-muted-foreground mt-2"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
          >
            {formatHoursMinutes(hoursToNextPhase)} until {nextPhase.name}
          </motion.p>
        )}
      </div>

      {/* Phase Progress Indicator */}
      <div className="space-y-3">
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>Fasting Phases</span>
          <span>{Math.round(progress)}% complete</span>
        </div>

        <div className="flex gap-1 h-2 rounded-full overflow-hidden bg-muted/30">
          {fastingPhases.map((phase) => {
            const phaseWidth = ((phase.endHour - phase.startHour) / targetHours) * 100;
            const isComplete = elapsedHours >= phase.endHour;
            const isCurrent = elapsedHours >= phase.startHour && elapsedHours < phase.endHour;
            const phaseProgress = isCurrent
              ? ((elapsedHours - phase.startHour) / (phase.endHour - phase.startHour)) * 100
              : isComplete
                ? 100
                : 0;

            return (
              <div
                key={phase.id}
                className="relative h-full overflow-hidden"
                style={{ width: `${phaseWidth}%` }}
              >
                <motion.div
                  className={cn(
                    'absolute inset-0',
                    phase.bgColor.replace('/20', '/60')
                  )}
                  initial={{ width: 0 }}
                  animate={{ width: `${phaseProgress}%` }}
                  transition={{ duration: 0.5, ease: 'easeOut' }}
                />
                {/* Shimmer effect when phase is active */}
                {isCurrent && isActive && (
                  <motion.div
                    className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent"
                    animate={{ x: ['-100%', '200%'] }}
                    transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
                  />
                )}
              </div>
            );
          })}
        </div>

        {/* Phase Labels */}
        <div className="flex justify-between">
          {fastingPhases.map((phase, index) => {
            const isActivePhase = elapsedHours >= phase.startHour && elapsedHours < phase.endHour;
            const isComplete = elapsedHours >= phase.endHour;

            return (
              <motion.div
                key={phase.id}
                className={cn(
                  'flex flex-col items-center gap-1 text-center transition-all duration-300',
                  isActivePhase ? phase.color : isComplete ? 'text-muted-foreground' : 'text-muted-foreground/50'
                )}
                animate={isActivePhase ? { scale: [1, 1.1, 1] } : {}}
                transition={{ duration: 2, repeat: isActivePhase ? Infinity : 0 }}
              >
                <motion.div
                  animate={isComplete ? { scale: 1.2 } : {}}
                  transition={{ type: 'spring', stiffness: 500 }}
                >
                  <phase.icon className="w-4 h-4" />
                </motion.div>
                <span className="text-xs font-medium hidden sm:block">{phase.name}</span>
                <span className="text-xs">{phase.startHour}h</span>
              </motion.div>
            );
          })}
        </div>
      </div>

      {/* Control Buttons */}
      <div className="flex gap-3 justify-center">
        {!isActive && elapsedSeconds === 0 ? (
          <Button
            className="bg-gradient-gold hover:scale-105 transition-transform shadow-gold-glow px-8"
            onClick={onStart}
          >
            <Flame className="w-4 h-4 mr-2" />
            Start Fast
          </Button>
        ) : (
          <>
            <Button
              variant="outline"
              onClick={isActive ? onPause : onStart}
              className="flex-1 max-w-32"
            >
              {isActive ? (
                <>
                  <Pause className="w-4 h-4 mr-2" />
                  Pause
                </>
              ) : (
                <>
                  <Flame className="w-4 h-4 mr-2" />
                  Resume
                </>
              )}
            </Button>
            <Button
              variant="destructive"
              onClick={onStop}
              className="flex-1 max-w-32"
            >
              <Square className="w-4 h-4 mr-2" />
              End Fast
            </Button>
          </>
        )}
      </div>
    </div>
  );
};

export default FastingTimer;