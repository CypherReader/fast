import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Scale, Plus, Target, TrendingDown, TrendingUp, Minus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  ReferenceLine,
  TooltipProps,
} from 'recharts';
import { useProgress } from '@/hooks/use-progress';
import { format } from 'date-fns';

interface WeightTrackerProps {
  onWeightChange?: (weight: number) => void;
  externalOpen?: boolean;
  externalOnOpenChange?: (open: boolean) => void;
}

const WeightTracker = ({ onWeightChange, externalOpen, externalOnOpenChange }: WeightTrackerProps) => {
  const { weightHistory, logWeight, isWeightLoading } = useProgress();

  const [currentWeight, setCurrentWeight] = useState<number | null>(null);
  const [goalWeight, setGoalWeight] = useState(170); // TODO: Fetch from user profile
  const [newWeight, setNewWeight] = useState('');
  const [newGoal, setNewGoal] = useState('');
  const [internalShowModal, setInternalShowModal] = useState(false);
  const [showGoalModal, setShowGoalModal] = useState(false);

  // Use external control if provided, otherwise use internal state
  const showAddModal = externalOpen ?? internalShowModal;
  const setShowAddModal = externalOnOpenChange ?? setInternalShowModal;

  useEffect(() => {
    if (weightHistory && weightHistory.length > 0) {
      const latest = weightHistory[0]; // Assuming sorted DESC from backend, but let's check
      // Actually backend sorts DESC, so [0] is latest
      setCurrentWeight(latest.weight_lbs);
      onWeightChange?.(latest.weight_lbs);
    }
  }, [weightHistory, onWeightChange]);

  const handleAddWeight = () => {
    const weight = parseFloat(newWeight);
    if (!isNaN(weight) && weight > 0) {
      logWeight({ weight, unit: 'lbs' }); // Assuming lbs for now
      setNewWeight('');
      setShowAddModal(false);
    }
  };

  const handleSetGoal = () => {
    const goal = parseFloat(newGoal);
    if (!isNaN(goal) && goal > 0) {
      setGoalWeight(goal);
      setNewGoal('');
      setShowGoalModal(false);
      // TODO: Save goal to backend
    }
  };

  // Process history for chart
  // Backend returns DESC, chart needs ASC
  const chartData = weightHistory ? [...weightHistory].reverse().map(entry => ({
    date: entry.logged_at,
    weight: entry.weight_lbs,
    displayDate: format(new Date(entry.logged_at), 'MMM d'),
  })) : [];

  const startWeight = chartData[0]?.weight || currentWeight || 0;
  const totalToLose = startWeight - goalWeight;
  const currentLoss = startWeight - (currentWeight || startWeight);
  const progressPercent = totalToLose > 0 ? Math.min((currentLoss / totalToLose) * 100, 100) : 0;
  const remainingToGoal = (currentWeight || 0) - goalWeight;
  const isGaining = remainingToGoal < 0;

  // Calculate trend
  const recentEntries = chartData.slice(-3);
  const trend = recentEntries.length >= 2
    ? recentEntries[recentEntries.length - 1].weight - recentEntries[0].weight
    : 0;

  // Chart domain
  const allWeights = [...chartData.map(e => e.weight), goalWeight, currentWeight || 0].filter(Boolean);
  const minWeight = allWeights.length > 0 ? Math.min(...allWeights) - 5 : 0;
  const maxWeight = allWeights.length > 0 ? Math.max(...allWeights) + 5 : 200;

  const CustomTooltip = ({ active, payload }: TooltipProps<number, string>) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-card border border-border rounded-lg px-3 py-2 shadow-lg">
          <p className="text-sm font-medium text-foreground">{payload[0].value} lbs</p>
          <p className="text-xs text-muted-foreground">{payload[0].payload.displayDate}</p>
        </div>
      );
    }
    return null;
  };

  return (
    <motion.div
      className="bg-card border border-border rounded-2xl overflow-hidden"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
    >
      {/* Header */}
      <div className="p-5 pb-0">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <div className="p-2 rounded-lg bg-secondary/20">
              <Scale className="w-5 h-5 text-secondary" />
            </div>
            <div>
              <h3 className="font-semibold text-foreground">Weight Tracker</h3>
              <p className="text-xs text-muted-foreground">Track your progress</p>
            </div>
          </div>
          <Dialog open={showAddModal} onOpenChange={setShowAddModal}>
            <DialogTrigger asChild>
              <Button size="sm" variant="outline" className="gap-1">
                <Plus className="w-4 h-4" />
                Log
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>Log Weight</DialogTitle>
              </DialogHeader>
              <div className="space-y-4 py-4">
                <div className="flex items-center gap-2">
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => setNewWeight(prev => {
                      const val = parseFloat(prev) || (currentWeight || 0);
                      return Math.max(0, val - 0.5).toString();
                    })}
                  >
                    <Minus className="w-4 h-4" />
                  </Button>
                  <Input
                    id="weight-input"
                    name="weight"
                    type="number"
                    step="0.1"
                    placeholder="Enter weight..."
                    value={newWeight}
                    onChange={(e) => setNewWeight(e.target.value)}
                    className="flex-1 text-center text-2xl font-bold h-14"
                  />
                  <Button
                    variant="outline"
                    size="icon"
                    onClick={() => setNewWeight(prev => {
                      const val = parseFloat(prev) || (currentWeight || 0);
                      return (val + 0.5).toString();
                    })}
                  >
                    <Plus className="w-4 h-4" />
                  </Button>
                </div>
                <p className="text-sm text-center text-muted-foreground">
                  {currentWeight && `Current: ${currentWeight} lbs`}
                </p>
                <Button onClick={handleAddWeight} className="w-full" disabled={!newWeight}>
                  Save Weight
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>

        {/* Current Weight & Goal */}
        <div className="grid grid-cols-2 gap-4 mb-4">
          <div className="bg-muted/30 rounded-xl p-4">
            <p className="text-xs text-muted-foreground mb-1">Current</p>
            <AnimatePresence mode="wait">
              <motion.p
                key={currentWeight || 'loading'}
                className="text-2xl font-bold text-foreground"
                initial={{ opacity: 0, y: 5 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -5 }}
              >
                {currentWeight ? `${currentWeight} lbs` : '-- lbs'}
              </motion.p>
            </AnimatePresence>
            {trend !== 0 && (
              <div className={cn(
                'flex items-center gap-1 mt-1',
                trend < 0 ? 'text-secondary' : 'text-orange-400'
              )}>
                {trend < 0 ? (
                  <TrendingDown className="w-3 h-3" />
                ) : (
                  <TrendingUp className="w-3 h-3" />
                )}
                <span className="text-xs font-medium">
                  {Math.abs(trend).toFixed(1)} lbs
                </span>
              </div>
            )}
          </div>

          <Dialog open={showGoalModal} onOpenChange={setShowGoalModal}>
            <DialogTrigger asChild>
              <button className="bg-muted/30 rounded-xl p-4 text-left hover:bg-muted/50 transition-colors">
                <div className="flex items-center justify-between">
                  <p className="text-xs text-muted-foreground mb-1">Goal</p>
                  <Target className="w-3 h-3 text-muted-foreground" />
                </div>
                <p className="text-2xl font-bold text-primary">{goalWeight} lbs</p>
                {!isGaining && remainingToGoal > 0 && (
                  <p className="text-xs text-muted-foreground mt-1">
                    {remainingToGoal.toFixed(1)} lbs to go
                  </p>
                )}
                {remainingToGoal <= 0 && (
                  <p className="text-xs text-secondary mt-1">Goal reached! 🎉</p>
                )}
              </button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>Set Goal Weight</DialogTitle>
              </DialogHeader>
              <div className="space-y-4 py-4">
                <Input
                  id="goal-weight-input"
                  name="goalWeight"
                  type="number"
                  step="0.1"
                  placeholder="Enter goal weight..."
                  value={newGoal}
                  onChange={(e) => setNewGoal(e.target.value)}
                  className="text-center text-2xl font-bold h-14"
                />
                <p className="text-sm text-center text-muted-foreground">
                  Current goal: {goalWeight} lbs
                </p>
                <Button onClick={handleSetGoal} className="w-full" disabled={!newGoal}>
                  Set Goal
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>

        {/* Progress Bar */}
        <div className="space-y-2 mb-4">
          <div className="flex items-center justify-between text-xs">
            <span className="text-muted-foreground">Progress to Goal</span>
            <span className="font-medium text-secondary">{Math.round(progressPercent)}%</span>
          </div>
          <div className="h-2 bg-muted/30 rounded-full overflow-hidden">
            <motion.div
              className="h-full bg-gradient-to-r from-primary to-secondary rounded-full"
              initial={{ width: 0 }}
              animate={{ width: `${progressPercent}%` }}
              transition={{ duration: 0.5, ease: 'easeOut' }}
            />
          </div>
          <div className="flex justify-between text-xs text-muted-foreground">
            <span>{startWeight} lbs</span>
            <span>{goalWeight} lbs</span>
          </div>
        </div>
      </div>

      {/* Chart */}
      <div className="h-40 px-2 pb-4">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={chartData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
            <XAxis
              dataKey="displayDate"
              tick={{ fontSize: 10, fill: 'hsl(var(--muted-foreground))' }}
              axisLine={false}
              tickLine={false}
            />
            <YAxis
              domain={[minWeight, maxWeight]}
              tick={{ fontSize: 10, fill: 'hsl(var(--muted-foreground))' }}
              axisLine={false}
              tickLine={false}
              width={35}
            />
            <Tooltip content={<CustomTooltip />} />
            <ReferenceLine
              y={goalWeight}
              stroke="hsl(var(--secondary))"
              strokeDasharray="5 5"
              strokeWidth={2}
            />
            <Line
              type="monotone"
              dataKey="weight"
              stroke="hsl(var(--primary))"
              strokeWidth={2}
              dot={{ fill: 'hsl(var(--primary))', strokeWidth: 0, r: 4 }}
              activeDot={{ fill: 'hsl(var(--primary))', strokeWidth: 0, r: 6 }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Stats Footer */}
      <div className="border-t border-border p-4 grid grid-cols-3 gap-2">
        <div className="text-center">
          <p className="text-lg font-bold text-foreground">{currentLoss.toFixed(1)}</p>
          <p className="text-xs text-muted-foreground">lbs lost</p>
        </div>
        <div className="text-center border-x border-border">
          <p className="text-lg font-bold text-foreground">{chartData.length}</p>
          <p className="text-xs text-muted-foreground">entries</p>
        </div>
        <div className="text-center">
          <p className="text-lg font-bold text-secondary">
            {remainingToGoal > 0 ? remainingToGoal.toFixed(1) : '0'}
          </p>
          <p className="text-xs text-muted-foreground">to goal</p>
        </div>
      </div>
    </motion.div>
  );
};

export default WeightTracker;