import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Droplets, Plus, Minus, Bell, BellOff, GlassWater } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { toast } from '@/hooks/use-toast';
import { useTelemetry } from '@/hooks/use-telemetry';
import { format } from 'date-fns';

interface WaterTrackerProps {
  dailyGoal?: number; // in ml
  glassSize?: number; // in ml
}

const WaterTracker = ({ dailyGoal = 2500, glassSize = 250 }: WaterTrackerProps) => {
  const { history, logMetric, isLoading } = useTelemetry('water');
  const [remindersEnabled, setRemindersEnabled] = useState(false);
  const [lastReminderTime, setLastReminderTime] = useState<Date | null>(null);

  // Calculate today's intake from history
  const today = new Date();
  const todayStr = format(today, 'yyyy-MM-dd');
  const intake = history?.find(h => h.date.startsWith(todayStr))?.value || 0;

  const glasses = Math.floor(intake / glassSize);
  const totalGlassesNeeded = Math.ceil(dailyGoal / glassSize);
  const progress = Math.min((intake / dailyGoal) * 100, 100);
  const isGoalReached = intake >= dailyGoal;

  // Reminder effect
  useEffect(() => {
    if (!remindersEnabled) return;

    const reminderInterval = setInterval(() => {
      if (!isGoalReached) {
        const now = new Date();
        const hours = now.getHours();

        // Only remind during waking hours (8am - 10pm)
        if (hours >= 8 && hours <= 22) {
          toast({
            title: "💧 Hydration Reminder",
            description: `You've had ${glasses} glasses today. Stay hydrated!`,
          });
          setLastReminderTime(now);
        }
      }
    }, 60 * 60 * 1000); // Every hour

    return () => clearInterval(reminderInterval);
  }, [remindersEnabled, isGoalReached, glasses]);

  const addWater = (amount: number) => {
    logMetric({ value: amount, unit: 'ml' });
  };

  const toggleReminders = () => {
    setRemindersEnabled(prev => {
      if (!prev) {
        toast({
          title: "Reminders Enabled",
          description: "You'll receive hourly hydration reminders.",
        });
      }
      return !prev;
    });
  };

  // Water level visualization
  const waterLevelHeight = Math.min(progress, 100);

  return (
    <motion.div
      className="bg-card border border-border rounded-2xl p-5"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
    >
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <div className="p-2 rounded-lg bg-blue-500/20">
            <Droplets className="w-5 h-5 text-blue-400" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground">Water Intake</h3>
            <p className="text-xs text-muted-foreground">Goal: {(dailyGoal / 1000).toFixed(1)}L</p>
          </div>
        </div>
        <Button
          size="sm"
          variant="ghost"
          onClick={toggleReminders}
          className={cn(
            'gap-1',
            remindersEnabled && 'text-blue-400'
          )}
        >
          {remindersEnabled ? (
            <Bell className="w-4 h-4" />
          ) : (
            <BellOff className="w-4 h-4" />
          )}
        </Button>
      </div>

      {/* Water Glass Visualization */}
      <div className="flex items-center gap-4 mb-4">
        <div className="relative w-20 h-28 flex-shrink-0">
          {/* Glass outline */}
          <div className="absolute inset-0 border-2 border-blue-400/30 rounded-b-2xl rounded-t-lg bg-muted/20" />

          {/* Water level */}
          <motion.div
            className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-blue-500 to-blue-400/80 rounded-b-2xl"
            initial={{ height: 0 }}
            animate={{ height: `${waterLevelHeight}%` }}
            transition={{ duration: 0.5, ease: 'easeOut' }}
          >
            {/* Bubbles animation */}
            {intake > 0 && (
              <>
                <motion.div
                  className="absolute w-2 h-2 bg-blue-300/40 rounded-full"
                  style={{ left: '20%', bottom: '20%' }}
                  animate={{ y: [-5, -15, -5], opacity: [0.4, 0.8, 0.4] }}
                  transition={{ duration: 2, repeat: Infinity }}
                />
                <motion.div
                  className="absolute w-1.5 h-1.5 bg-blue-300/40 rounded-full"
                  style={{ left: '60%', bottom: '40%' }}
                  animate={{ y: [-3, -12, -3], opacity: [0.3, 0.7, 0.3] }}
                  transition={{ duration: 2.5, repeat: Infinity, delay: 0.5 }}
                />
              </>
            )}
          </motion.div>

          {/* Percentage overlay */}
          <div className="absolute inset-0 flex items-center justify-center">
            <span className="text-lg font-bold text-foreground drop-shadow-lg">
              {Math.round(progress)}%
            </span>
          </div>
        </div>

        <div className="flex-1 space-y-3">
          {/* Current intake */}
          <div>
            <AnimatePresence mode="wait">
              <motion.p
                key={intake}
                className="text-2xl font-bold text-foreground"
                initial={{ opacity: 0, y: 5 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -5 }}
              >
                {intake >= 1000 ? `${(intake / 1000).toFixed(1)}L` : `${intake}ml`}
              </motion.p>
            </AnimatePresence>
            <p className="text-xs text-muted-foreground">
              {dailyGoal - intake > 0
                ? `${((dailyGoal - intake) / 1000).toFixed(1)}L remaining`
                : 'Goal reached! 🎉'}
            </p>
          </div>

          {/* Glasses indicator */}
          <div className="flex flex-wrap gap-1">
            {Array.from({ length: totalGlassesNeeded }).map((_, i) => (
              <motion.div
                key={i}
                initial={{ scale: 0.8, opacity: 0 }}
                animate={{
                  scale: i < glasses ? 1 : 0.8,
                  opacity: 1
                }}
                transition={{ delay: i * 0.03 }}
              >
                <GlassWater
                  className={cn(
                    'w-4 h-4 transition-colors',
                    i < glasses ? 'text-blue-400' : 'text-muted-foreground/30'
                  )}
                />
              </motion.div>
            ))}
          </div>
          <p className="text-xs text-muted-foreground">
            {glasses} / {totalGlassesNeeded} glasses ({glassSize}ml each)
          </p>
        </div>
      </div>

      {/* Quick Add Buttons */}
      <div className="grid grid-cols-4 gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => addWater(-glassSize)}
          disabled={intake < glassSize}
          className="h-12 flex-col p-1"
        >
          <Minus className="w-4 h-4" />
          <span className="text-xs">-1</span>
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => addWater(glassSize)}
          className="h-12 flex-col p-1 border-blue-500/30 hover:bg-blue-500/10"
        >
          <GlassWater className="w-4 h-4 text-blue-400" />
          <span className="text-xs">+1</span>
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => addWater(500)}
          className="h-12 flex-col p-1 border-blue-500/30 hover:bg-blue-500/10"
        >
          <Droplets className="w-4 h-4 text-blue-400" />
          <span className="text-xs">500ml</span>
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={() => addWater(1000)}
          className="h-12 flex-col p-1 border-blue-500/30 hover:bg-blue-500/10"
        >
          <Plus className="w-4 h-4 text-blue-400" />
          <span className="text-xs">1L</span>
        </Button>
      </div>

      {/* Goal reached celebration */}
      <AnimatePresence>
        {isGoalReached && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="mt-4 bg-blue-500/20 border border-blue-500/30 rounded-lg p-3 text-center"
          >
            <p className="text-sm font-medium text-blue-400">
              🎉 Daily hydration goal reached!
            </p>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  );
};

export default WaterTracker;