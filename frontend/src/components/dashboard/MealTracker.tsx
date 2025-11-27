import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Utensils, Plus, Trash2, Flame, Clock, Apple, Coffee, Sandwich, Pizza } from 'lucide-react';
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useMeals, Meal } from '@/hooks/use-meals';
import { format, parseISO } from 'date-fns';

interface MealTrackerProps {
  dailyCalorieGoal?: number;
}

const mealTypeConfig: Record<string, { icon: any; label: string; color: string; bg: string }> = {
  breakfast: { icon: Coffee, label: 'Breakfast', color: 'text-orange-400', bg: 'bg-orange-500/20' },
  lunch: { icon: Sandwich, label: 'Lunch', color: 'text-green-400', bg: 'bg-green-500/20' },
  dinner: { icon: Pizza, label: 'Dinner', color: 'text-blue-400', bg: 'bg-blue-500/20' },
  snack: { icon: Apple, label: 'Snack', color: 'text-pink-400', bg: 'bg-pink-500/20' },
};

const quickMeals = [
  { name: 'Eggs & Bacon', calories: 350, type: 'breakfast' as const },
  { name: 'Greek Yogurt', calories: 150, type: 'breakfast' as const },
  { name: 'Grilled Chicken Salad', calories: 400, type: 'lunch' as const },
  { name: 'Salmon & Veggies', calories: 500, type: 'dinner' as const },
  { name: 'Protein Shake', calories: 200, type: 'snack' as const },
  { name: 'Almonds (1oz)', calories: 160, type: 'snack' as const },
];

const MealTracker = ({ dailyCalorieGoal = 1800 }: MealTrackerProps) => {
  const { meals, logMeal, isLoading } = useMeals();
  const [showAddModal, setShowAddModal] = useState(false);
  const [newMealName, setNewMealName] = useState('');
  const [newMealCalories, setNewMealCalories] = useState('');
  const [newMealType, setNewMealType] = useState<'breakfast' | 'lunch' | 'dinner' | 'snack'>('lunch');

  // Filter meals for today
  const todayStr = format(new Date(), 'yyyy-MM-dd');
  const todaysMeals = meals.filter(meal => {
    // Handle potential timezone issues or just check date string match if backend returns ISO
    return meal.logged_at.startsWith(todayStr);
  });

  const totalCalories = todaysMeals.reduce((sum, meal) => sum + (meal.calories || 0), 0);
  const progress = Math.min((totalCalories / dailyCalorieGoal) * 100, 100);
  const remainingCalories = dailyCalorieGoal - totalCalories;
  const isOverGoal = totalCalories > dailyCalorieGoal;

  const handleAddMeal = (name: string, calories: number, mealType: 'breakfast' | 'lunch' | 'dinner' | 'snack') => {
    logMeal({
      name,
      calories,
      meal_type: mealType,
    });
  };

  const handleAddCustomMeal = () => {
    const calories = parseInt(newMealCalories);
    if (newMealName && !isNaN(calories) && calories > 0) {
      handleAddMeal(newMealName, calories, newMealType);
      setNewMealName('');
      setNewMealCalories('');
      setShowAddModal(false);
    }
  };

  // Note: Backend doesn't support deleting meals yet via API in this scope, 
  // so removeMeal is disabled or needs to be implemented in backend.
  // For now, I'll remove the delete button or make it a TODO.
  const removeMeal = (id: string) => {
    console.log('Delete meal not implemented yet', id);
  };

  const getProgressColor = () => {
    if (isOverGoal) return 'from-red-500 to-orange-500';
    if (progress >= 80) return 'from-yellow-500 to-orange-500';
    return 'from-secondary to-primary';
  };

  // Group meals by type
  const mealsByType = todaysMeals.reduce((acc, meal) => {
    const type = meal.meal_type || 'snack'; // Default to snack if undefined
    if (!acc[type]) acc[type] = [];
    acc[type].push(meal);
    return acc;
  }, {} as Record<string, Meal[]>);

  return (
    <motion.div
      className="bg-card border border-border rounded-2xl overflow-hidden"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
    >
      {/* Header */}
      <div className="p-5 pb-4">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <div className="p-2 rounded-lg bg-orange-500/20">
              <Utensils className="w-5 h-5 text-orange-400" />
            </div>
            <div>
              <h3 className="font-semibold text-foreground">Meal Log</h3>
              <p className="text-xs text-muted-foreground">Daily Goal: {dailyCalorieGoal.toLocaleString()} cal</p>
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
                <DialogTitle>Log a Meal</DialogTitle>
              </DialogHeader>
              <div className="space-y-4 py-4">
                {/* Quick Add */}
                <div>
                  <p className="text-sm font-medium text-foreground mb-2">Quick Add</p>
                  <div className="grid grid-cols-2 gap-2">
                    {quickMeals.map((meal) => {
                      const config = mealTypeConfig[meal.type];
                      return (
                        <button
                          key={meal.name}
                          onClick={() => {
                            handleAddMeal(meal.name, meal.calories, meal.type);
                            setShowAddModal(false);
                          }}
                          className={cn(
                            'p-3 rounded-lg border border-border text-left hover:bg-muted/50 transition-colors',
                          )}
                        >
                          <div className="flex items-center gap-2 mb-1">
                            <config.icon className={cn('w-4 h-4', config.color)} />
                            <span className="text-sm font-medium text-foreground truncate">{meal.name}</span>
                          </div>
                          <span className="text-xs text-muted-foreground">{meal.calories} cal</span>
                        </button>
                      );
                    })}
                  </div>
                </div>

                <div className="relative">
                  <div className="absolute inset-0 flex items-center">
                    <span className="w-full border-t border-border" />
                  </div>
                  <div className="relative flex justify-center text-xs uppercase">
                    <span className="bg-card px-2 text-muted-foreground">Or add custom</span>
                  </div>
                </div>

                {/* Custom Meal Form */}
                <div className="space-y-3">
                  <Input
                    placeholder="Meal name..."
                    value={newMealName}
                    onChange={(e) => setNewMealName(e.target.value)}
                  />
                  <div className="flex gap-2">
                    <Input
                      type="number"
                      placeholder="Calories"
                      value={newMealCalories}
                      onChange={(e) => setNewMealCalories(e.target.value)}
                      className="flex-1"
                    />
                    <Select value={newMealType} onValueChange={(v) => setNewMealType(v as 'breakfast' | 'lunch' | 'dinner' | 'snack')}>
                      <SelectTrigger className="w-32">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {Object.entries(mealTypeConfig).map(([key, config]) => (
                          <SelectItem key={key} value={key}>
                            <div className="flex items-center gap-2">
                              <config.icon className={cn('w-4 h-4', config.color)} />
                              {config.label}
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <Button
                    onClick={handleAddCustomMeal}
                    className="w-full"
                    disabled={!newMealName || !newMealCalories}
                  >
                    Add Meal
                  </Button>
                </div>
              </div>
            </DialogContent>
          </Dialog>
        </div>

        {/* Calorie Summary */}
        <div className="flex items-center gap-4 mb-4">
          <div className="flex-1">
            <div className="flex items-baseline gap-2 mb-1">
              <AnimatePresence mode="wait">
                <motion.span
                  key={totalCalories}
                  className={cn(
                    'text-3xl font-bold',
                    isOverGoal ? 'text-red-400' : 'text-foreground'
                  )}
                  initial={{ opacity: 0, y: 5 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -5 }}
                >
                  {totalCalories.toLocaleString()}
                </motion.span>
              </AnimatePresence>
              <span className="text-muted-foreground">/ {dailyCalorieGoal.toLocaleString()} cal</span>
            </div>

            {/* Progress Bar */}
            <div className="h-2 bg-muted/30 rounded-full overflow-hidden">
              <motion.div
                className={cn('h-full rounded-full bg-gradient-to-r', getProgressColor())}
                initial={{ width: 0 }}
                animate={{ width: `${Math.min(progress, 100)}%` }}
                transition={{ duration: 0.5 }}
              />
            </div>

            <p className={cn(
              'text-xs mt-1',
              isOverGoal ? 'text-red-400' : 'text-muted-foreground'
            )}>
              {isOverGoal
                ? `${Math.abs(remainingCalories).toLocaleString()} cal over goal`
                : `${remainingCalories.toLocaleString()} cal remaining`
              }
            </p>
          </div>

          {/* Flame indicator */}
          <div className={cn(
            'p-3 rounded-xl',
            isOverGoal ? 'bg-red-500/20' : 'bg-secondary/20'
          )}>
            <Flame className={cn(
              'w-8 h-8',
              isOverGoal ? 'text-red-400' : 'text-secondary'
            )} />
          </div>
        </div>
      </div>

      {/* Meal List */}
      <div className="border-t border-border max-h-48 overflow-y-auto">
        {todaysMeals.length === 0 ? (
          <div className="p-6 text-center">
            <Utensils className="w-8 h-8 text-muted-foreground/50 mx-auto mb-2" />
            <p className="text-sm text-muted-foreground">No meals logged yet</p>
            <p className="text-xs text-muted-foreground">Tap "Log" to add your first meal</p>
          </div>
        ) : (
          <div className="divide-y divide-border">
            {todaysMeals.map((meal) => {
              const config = mealTypeConfig[meal.meal_type || 'snack'];
              return (
                <motion.div
                  key={meal.id}
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: 10 }}
                  className="flex items-center gap-3 p-3 hover:bg-muted/30 transition-colors group"
                >
                  <div className={cn('p-2 rounded-lg', config.bg)}>
                    <config.icon className={cn('w-4 h-4', config.color)} />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-foreground truncate">{meal.name}</p>
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <span>{config.label}</span>
                      <span>•</span>
                      <Clock className="w-3 h-3" />
                      <span>{format(parseISO(meal.logged_at), 'h:mm a')}</span>
                    </div>
                  </div>
                  <span className="text-sm font-medium text-foreground">{meal.calories} cal</span>
                  {/* Delete button removed as backend support is pending */}
                </motion.div>
              );
            })}
          </div>
        )}
      </div>

      {/* Macro Summary Footer */}
      {todaysMeals.length > 0 && (
        <div className="border-t border-border p-4 grid grid-cols-4 gap-2 text-center">
          <div>
            <p className="text-lg font-bold text-foreground">{todaysMeals.length}</p>
            <p className="text-xs text-muted-foreground">meals</p>
          </div>
          <div>
            <p className="text-lg font-bold text-orange-400">
              {mealsByType.breakfast?.reduce((s, m) => s + (m.calories || 0), 0) || 0}
            </p>
            <p className="text-xs text-muted-foreground">breakfast</p>
          </div>
          <div>
            <p className="text-lg font-bold text-green-400">
              {mealsByType.lunch?.reduce((s, m) => s + (m.calories || 0), 0) || 0}
            </p>
            <p className="text-xs text-muted-foreground">lunch</p>
          </div>
          <div>
            <p className="text-lg font-bold text-blue-400">
              {mealsByType.dinner?.reduce((s, m) => s + (m.calories || 0), 0) || 0}
            </p>
            <p className="text-xs text-muted-foreground">dinner</p>
          </div>
        </div>
      )}
    </motion.div>
  );
};

export default MealTracker;