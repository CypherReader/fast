import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  X, Lock, Timer, CheckCircle, Circle, Scale, Utensils,
  Users, Bell, TrendingUp, Flame, Target, ChevronRight
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { useOnboarding } from '@/contexts/OnboardingContext';
import { cn } from '@/lib/utils';
import FastingTimer from '@/components/dashboard/FastingTimer';
import StepTracker from '@/components/dashboard/StepTracker';
import WeightTracker from '@/components/dashboard/WeightTracker';
import WaterTracker from '@/components/dashboard/WaterTracker';
import MealTracker from '@/components/dashboard/MealTracker';
import { useFasting } from '@/hooks/use-fasting';
import { useUser } from '@/hooks/use-user';
import { useProgress } from '@/hooks/use-progress';
import confetti from 'canvas-confetti';
import UserMenu from '@/components/layout/UserMenu';

const planNames: Record<string, string> = {
  '16:8': '16:8',
  '18:6': '18:6',
  'omad': '23:1 OMAD',
};

const checklistItems = [
  { id: 'account', label: 'Create account', completed: true, icon: CheckCircle },
  { id: 'deposit', label: 'Deposit to Vault', completed: true, icon: Lock },
  { id: 'first_fast', label: 'Start your first fast', completed: false, icon: Timer, action: 'start_fast' },
  { id: 'log_weight', label: 'Log your weight', completed: false, icon: Scale, action: 'log_weight' },
  { id: 'log_steps', label: 'Log your steps', completed: false, icon: TrendingUp, action: 'log_steps' },
  { id: 'join_tribe', label: 'Join a tribe', completed: false, icon: Users, action: 'join_tribe' },
  { id: 'notifications', label: 'Enable notifications', completed: false, icon: Bell, action: 'notifications' },
];

const Dashboard = () => {
  const navigate = useNavigate();
  const { state } = useOnboarding();
  const [showWelcomeBanner, setShowWelcomeBanner] = useState(true);
  const [checklist, setChecklist] = useState(checklistItems);
  const [isFirstVisit] = useState(true); // In real app, check from user data
  const [showSetupComplete, setShowSetupComplete] = useState(false);
  const [setupJustCompleted, setSetupJustCompleted] = useState(false);

  const { currentFast, startFast, stopFast, isLoading, getInsight } = useFasting();
  const { user, stats } = useUser();
  const { weightHistory, dailyHydration } = useProgress();
  const [elapsedSeconds, setElapsedSeconds] = useState(0);
  const [insight, setInsight] = useState<string | null>(null);

  // Dialog control states for Quick Actions
  const [showMealDialog, setShowMealDialog] = useState(false);
  const [showWeightDialog, setShowWeightDialog] = useState(false);

  // Auto-dismiss welcome banner after 10 seconds
  useEffect(() => {
    if (showWelcomeBanner) {
      const timer = setTimeout(() => setShowWelcomeBanner(false), 10000);
      return () => clearTimeout(timer);
    }
  }, [showWelcomeBanner]);

  // Fetch insight when fast is active (once per session load for now)
  useEffect(() => {
    if (currentFast?.status === 'active' && !insight) {
      const hours = (new Date().getTime() - new Date(currentFast.start_time).getTime()) / (1000 * 60 * 60);
      getInsight(hours).then((data: any) => {
        if (data.insight) setInsight(data.insight);
      }).catch(console.error);
    }
  }, [currentFast, getInsight, insight]);

  // Calculate elapsed time
  useEffect(() => {
    console.log('Dashboard: currentFast:', currentFast);
    if (currentFast?.status === 'active' && currentFast.start_time) {
      const startTime = new Date(currentFast.start_time).getTime();

      const updateTimer = () => {
        const now = new Date().getTime();
        setElapsedSeconds(Math.floor((now - startTime) / 1000));
      };

      updateTimer(); // Update immediately
      const interval = setInterval(updateTimer, 1000);
      return () => clearInterval(interval);
    } else {
      setElapsedSeconds(0);
    }
  }, [currentFast]);

  const getTargetHours = () => {
    const plan = state.fastingPlan || '16:8';
    if (plan === '18:6') return 18;
    if (plan === 'omad') return 23;
    return 16;
  };

  const handleStartFast = () => {
    const plan = state.fastingPlan || '16:8';
    const goalHours = getTargetHours();

    startFast({
      plan_type: plan,
      goal_hours: goalHours,
      start_time: new Date().toISOString()
    }, {
      onSuccess: () => {
        setChecklist(prev => prev.map(item =>
          item.id === 'first_fast' ? { ...item, completed: true } : item
        ));
      }
    });
  };

  const handleStopFast = () => {
    stopFast();
  };

  // Update checklist based on actual progress data
  useEffect(() => {
    setChecklist(prev => prev.map(item => {
      // Mark fast complete if there's an active or completed fast
      if (item.id === 'first_fast' && currentFast) {
        return { ...item, completed: true };
      }
      // Mark weight complete if there's weight history
      if (item.id === 'log_weight' && weightHistory && weightHistory.length > 0) {
        return { ...item, completed: true };
      }
      // Mark steps complete if there's hydration data (temp, until steps endpoint exists)
      if (item.id === 'log_steps' && dailyHydration && dailyHydration.glasses_count > 0) {
        return { ...item, completed: true };
      }
      return item;
    }));
  }, [currentFast, weightHistory, dailyHydration]);

  // Check for setup completion and trigger celebration
  useEffect(() => {
    const completedCount = checklist.filter(item => item.completed).length;
    const wasComplete = completedCount === checklist.length;

    if (wasComplete && !setupJustCompleted && !showSetupComplete) {
      setSetupJustCompleted(true);
      setShowSetupComplete(true);

      // Trigger confetti
      const duration = 3000;
      const end = Date.now() + duration;

      const colors = ['#f59e0b', '#10b981', '#a855f7'];

      const frame = () => {
        confetti({
          particleCount: 5,
          angle: 60,
          spread: 55,
          origin: { x: 0, y: 0.7 },
          colors,
        });
        confetti({
          particleCount: 5,
          angle: 120,
          spread: 55,
          origin: { x: 1, y: 0.7 },
          colors,
        });

        if (Date.now() < end) {
          requestAnimationFrame(frame);
        }
      };

      frame();

      // Center burst
      setTimeout(() => {
        confetti({
          particleCount: 100,
          spread: 70,
          origin: { y: 0.6 },
          colors,
        });
      }, 200);
    }
  }, [checklist, setupJustCompleted, showSetupComplete]);

  const handleChecklistAction = (action: string) => {
    if (action === 'start_fast') {
      handleStartFast();
    }
    // Other actions would navigate to respective features
  };

  const completedCount = checklist.filter(item => item.completed).length;
  const progressPercentage = (completedCount / checklist.length) * 100;

  // Calculate vault progress
  const vaultDeposit = stats?.vault_total || 20;
  const recoveredAmount = stats?.vault_balance || 0;
  const fastsCompleted = stats?.fasts_completed || 0;
  const fastsNeeded = 10; // This could also be dynamic based on plan

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-40">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <img src="/fasthero.png" alt="FastingHero" className="w-8 h-8 rounded-lg" />
            <span className="font-bold text-lg text-foreground">FastingHero</span>
          </div>
          <nav className="hidden md:flex items-center gap-6 text-sm">
            <button className="text-foreground font-medium">Dashboard</button>
            <button onClick={() => navigate('/progress')} className="text-muted-foreground hover:text-foreground transition-colors">Progress</button>
            <button onClick={() => navigate('/community')} className="text-muted-foreground hover:text-foreground transition-colors">Community</button>
            <button onClick={() => navigate('/resources')} className="text-muted-foreground hover:text-foreground transition-colors">Resources</button>
          </nav>
          <UserMenu />
        </div>
      </header>

      <main className="container mx-auto px-4 py-6 max-w-6xl">
        {/* Welcome Banner */}
        <AnimatePresence>
          {showWelcomeBanner && isFirstVisit && (
            <motion.div
              initial={{ opacity: 0, y: -20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className="relative bg-gradient-to-r from-primary/20 via-secondary/20 to-primary/20 border border-primary/30 rounded-2xl p-5 mb-6"
            >
              <button
                onClick={() => setShowWelcomeBanner(false)}
                className="absolute top-4 right-4 text-muted-foreground hover:text-foreground transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
              <div className="flex items-center gap-3">
                <span className="text-2xl">👋</span>
                <div>
                  <h2 className="font-semibold text-foreground">Welcome to FastingHero!</h2>
                  <p className="text-sm text-muted-foreground">
                    You&apos;re all set! Start your first fast to begin recovering your $20 deposit.
                  </p>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Setup Complete Celebration */}
        <AnimatePresence>
          {showSetupComplete && (
            <motion.div
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.9 }}
              className="fixed inset-0 z-50 flex items-center justify-center bg-background/80 backdrop-blur-sm"
              onClick={() => setShowSetupComplete(false)}
            >
              <motion.div
                initial={{ y: 50 }}
                animate={{ y: 0 }}
                className="relative bg-gradient-to-br from-primary/20 via-card to-secondary/20 border-2 border-primary/50 rounded-3xl p-8 max-w-md mx-4 text-center shadow-2xl"
                onClick={(e) => e.stopPropagation()}
              >
                <motion.div
                  animate={{ rotate: [0, 10, -10, 10, 0] }}
                  transition={{ duration: 0.5, repeat: 3 }}
                  className="text-6xl mb-4"
                >
                  🎉
                </motion.div>
                <h2 className="text-3xl font-bold text-foreground mb-2">
                  Setup Complete!
                </h2>
                <p className="text-lg text-secondary font-semibold mb-4">
                  You're all set to begin your journey!
                </p>
                <p className="text-muted-foreground mb-6">
                  "The journey of a thousand miles begins with one step. You've taken yours. Now watch as small daily actions transform into extraordinary results."
                </p>
                <div className="flex gap-3">
                  <Button
                    onClick={() => setShowSetupComplete(false)}
                    className="flex-1 bg-gradient-gold hover:scale-105 transition-transform"
                  >
                    Let's Go! 🚀
                  </Button>
                </div>
              </motion.div>
            </motion.div>
          )}
        </AnimatePresence>

        <div className="grid lg:grid-cols-3 gap-6">
          {/* Main Content - Left Side */}
          <div className="lg:col-span-2 space-y-6">
            {/* Vault Hero Card */}
            <motion.div
              className="bg-card border border-border rounded-2xl overflow-hidden"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
            >
              <div className="p-6">
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-2">
                    <Lock className="w-5 h-5 text-primary" />
                    <h3 className="font-semibold text-foreground">Your Vault: Month 1</h3>
                  </div>
                  <span className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded-full">
                    {state.fastingPlan ? planNames[state.fastingPlan] : '16:8'} Plan
                  </span>
                </div>

                {/* Vault Balance */}
                <div className="flex items-baseline gap-2 mb-2">
                  <span className="text-4xl font-bold text-primary">${recoveredAmount.toFixed(2)}</span>
                  <span className="text-muted-foreground">/ ${vaultDeposit.toFixed(2)} recovered</span>
                </div>

                {/* Progress Bar */}
                <Progress value={(recoveredAmount / vaultDeposit) * 100} className="h-3 mb-4" />

                <div className="flex items-center justify-between text-sm mb-6">
                  <span className="text-muted-foreground">
                    {fastsCompleted} / {fastsNeeded} fasts completed
                  </span>
                  <span className="text-secondary font-medium">
                    {fastsNeeded - fastsCompleted} more for full refund
                  </span>
                </div>

                {/* Fasting Timer with Phase Indicators */}
                <FastingTimer
                  isActive={currentFast?.status === 'active'}
                  elapsedSeconds={elapsedSeconds}
                  targetHours={getTargetHours()}
                  onStart={handleStartFast}
                  onPause={() => { }} // Pause not supported by backend yet
                  onStop={handleStopFast}
                  insight={insight}
                />
              </div>
            </motion.div>

            {/* Quick Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[
                { label: 'Discipline', value: '-- / 100', icon: Target, sublabel: 'Complete a fast to see' },
                { label: 'Fasts', value: fastsCompleted.toString(), icon: Timer, sublabel: 'completed' },
                { label: 'Vault Balance', value: `$${recoveredAmount}`, icon: Lock, sublabel: `of $${vaultDeposit}` },
                { label: 'Streak', value: '0', icon: Flame, sublabel: 'days' },
              ].map((stat, index) => (
                <motion.div
                  key={stat.label}
                  className="bg-card border border-border rounded-xl p-4"
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: 0.1 * index }}
                >
                  <stat.icon className="w-5 h-5 text-muted-foreground mb-2" />
                  <div className="text-xl font-bold text-foreground">{stat.value}</div>
                  <div className="text-xs text-muted-foreground">{stat.label}</div>
                  <div className="text-xs text-muted-foreground/70">{stat.sublabel}</div>
                </motion.div>
              ))}
            </div>

            {/* Weight Tracker with Chart */}
            <WeightTracker externalOpen={showWeightDialog} externalOnOpenChange={setShowWeightDialog} />

            {/* Meal Tracker */}
            <MealTracker dailyCalorieGoal={1800} externalOpen={showMealDialog} externalOnOpenChange={setShowMealDialog} />
          </div>

          {/* Sidebar - Right Side */}
          <div className="space-y-6">
            {/* Getting Started Checklist */}
            <motion.div
              className="bg-card border border-border rounded-2xl p-5"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
            >
              <div className="flex items-center justify-between mb-4">
                <h3 className="font-semibold text-foreground">Complete Your Setup</h3>
                <span className="text-xs text-muted-foreground">
                  {completedCount}/{checklist.length}
                </span>
              </div>

              <Progress value={progressPercentage} className="h-2 mb-4" />

              <div className="space-y-2">
                {checklist.map((item) => (
                  <button
                    key={item.id}
                    onClick={() => !item.completed && item.action && handleChecklistAction(item.action)}
                    disabled={item.completed}
                    className={cn(
                      'w-full flex items-center gap-3 p-3 rounded-lg transition-all text-left',
                      item.completed
                        ? 'bg-secondary/10 cursor-default'
                        : 'bg-muted/30 hover:bg-muted/50 cursor-pointer'
                    )}
                  >
                    {item.completed ? (
                      <CheckCircle className="w-5 h-5 text-secondary flex-shrink-0" />
                    ) : (
                      <Circle className="w-5 h-5 text-muted-foreground flex-shrink-0" />
                    )}
                    <span className={cn(
                      'text-sm flex-1',
                      item.completed ? 'text-muted-foreground line-through' : 'text-foreground'
                    )}>
                      {item.label}
                    </span>
                    {!item.completed && (
                      <ChevronRight className="w-4 h-4 text-muted-foreground" />
                    )}
                  </button>
                ))}
              </div>
            </motion.div>

            {/* Quick Actions */}
            <motion.div
              className="bg-card border border-border rounded-2xl p-5"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.1 }}
            >
              <h3 className="font-semibold text-foreground mb-4">Quick Actions</h3>
              <div className="space-y-2">
                <button
                  onClick={() => setShowWeightDialog(true)}
                  className="w-full flex items-center gap-3 p-3 rounded-lg bg-muted/30 hover:bg-muted/50 transition-colors"
                >
                  <Scale className="w-5 h-5 text-secondary" />
                  <span className="text-sm text-foreground">Log Weight</span>
                </button>
                <button
                  onClick={() => setShowMealDialog(true)}
                  className="w-full flex items-center gap-3 p-3 rounded-lg bg-muted/30 hover:bg-muted/50 transition-colors"
                >
                  <Utensils className="w-5 h-5 text-secondary" />
                  <span className="text-sm text-foreground">Log Meal</span>
                </button>
                <button className="w-full flex items-center gap-3 p-3 rounded-lg bg-muted/30 hover:bg-muted/50 transition-colors">
                  <Users className="w-5 h-5 text-secondary" />
                  <span className="text-sm text-foreground">Find a Tribe</span>
                </button>
              </div>
            </motion.div>

            {/* Step Tracker */}
            <StepTracker dailyGoal={10000} />

            {/* Water Tracker */}
            <WaterTracker dailyGoal={2500} />

            {/* Help Card */}
            <motion.div
              className="bg-gradient-to-br from-secondary/10 to-primary/10 border border-secondary/30 rounded-2xl p-5"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: 0.2 }}
            >
              <h4 className="font-medium text-foreground mb-2">Need Help?</h4>
              <p className="text-sm text-muted-foreground mb-3">
                Check out our guides to get the most out of FastingHero.
              </p>
              <Button variant="outline" size="sm" className="w-full" onClick={() => navigate('/resources')}>
                View Resources
              </Button>
            </motion.div>
          </div>
        </div>
      </main>

      {/* Bottom Navigation - Mobile */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-card border-t border-border px-4 py-2 z-40">
        <div className="flex justify-around">
          {[
            { icon: Lock, label: 'Dashboard', active: true },
            { icon: TrendingUp, label: 'Progress' },
            { icon: Timer, label: 'Activity' },
            { icon: Users, label: 'Community' },
          ].map((item) => (
            <button
              key={item.label}
              onClick={() => {
                if (item.label === 'Progress') navigate('/progress');
                if (item.label === 'Community') navigate('/community');
              }}
              className={cn(
                'flex flex-col items-center gap-1 p-2',
                item.active ? 'text-primary' : 'text-muted-foreground'
              )}
            >
              <item.icon className="w-5 h-5" />
              <span className="text-xs">{item.label}</span>
            </button>
          ))}
        </div>
      </nav>
    </div>
  );
};

export default Dashboard;
