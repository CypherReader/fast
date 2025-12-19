import { useState } from 'react';
import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import {
    ArrowLeft, Timer, Droplets, Footprints, Utensils,
    TrendingUp, BarChart3
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useUser } from '@/hooks/use-user';
import { useProgress } from '@/hooks/use-progress';
import { useFasting, useFastingHistory } from '@/hooks/use-fasting';
import { useMeals } from '@/hooks/use-meals';
import UserMenu from '@/components/layout/UserMenu';
import { ProgressRing } from '@/components/progress/ProgressRing';
import { ShareCard } from '@/components/progress/ShareCard';
import { Achievements } from '@/components/progress/Achievements';
import { WeeklyReportCard } from '@/components/progress/WeeklyReportCard';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
} from 'recharts';
import { format, subDays, parseISO } from 'date-fns';

const Progress = () => {
    const navigate = useNavigate();
    const { user, stats } = useUser();
    const { dailyHydration } = useProgress();
    const { currentFast } = useFasting();
    const { history: fastingHistory } = useFastingHistory();
    const { meals } = useMeals();
    const [stepsToday, setStepsToday] = useState(0); // Manual entry for now
    const [showStepsInput, setShowStepsInput] = useState(false);

    // Calculate today's calories
    const todayStr = format(new Date(), 'yyyy-MM-dd');
    const todayCalories = meals?.filter(m => m.logged_at.startsWith(todayStr))
        .reduce((sum, m) => sum + (m.calories || 0), 0) || 0;

    // Calculate current fast progress
    const fastingGoalHours = 16; // Default, could come from user settings
    const currentFastHours = currentFast?.status === 'active' && currentFast.start_time
        ? (new Date().getTime() - new Date(currentFast.start_time).getTime()) / (1000 * 60 * 60)
        : 0;

    // Transform fasting history to chart format
    const fastingChartData = Array.from({ length: 7 }, (_, i) => {
        const day = subDays(new Date(), 6 - i);
        const dayStr = format(day, 'yyyy-MM-dd');

        const dailyFasts = fastingHistory?.filter(fast => {
            if (!fast.end_time || fast.status !== 'completed') return false;
            const fastDate = format(parseISO(fast.end_time), 'yyyy-MM-dd');
            return fastDate === dayStr;
        }) || [];

        const totalHours = dailyFasts.reduce((sum, fast) => {
            const start = new Date(fast.start_time).getTime();
            const end = new Date(fast.end_time!).getTime();
            return sum + (end - start) / (1000 * 60 * 60);
        }, 0);

        return {
            date: format(day, 'EEE'),
            hours: Math.round(totalHours * 10) / 10
        };
    });

    const handleStepsSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setShowStepsInput(false);
    };

    return (
        <div className="min-h-screen bg-background">
            {/* Header */}
            <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-40">
                <div className="container mx-auto px-4 py-4 flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <Button variant="ghost" size="icon" onClick={() => navigate('/dashboard')}>
                            <ArrowLeft className="w-5 h-5" />
                        </Button>
                        <div className="flex items-center gap-2">
                            <img src="/fasthero.png" alt="FastingHero" className="w-6 h-6 rounded-lg" />
                            <h1 className="font-bold text-lg text-foreground">Your Progress</h1>
                        </div>
                    </div>
                    <UserMenu />
                </div>
            </header>

            <main className="container mx-auto px-4 py-6 max-w-4xl space-y-6">
                {/* Share Card Section */}
                <ShareCard
                    userName={user?.name || user?.email || 'Faster'}
                    streak={stats?.current_streak || 0}
                    fastingPlan={currentFast?.plan_type || '16:8'}
                    totalFasts={stats?.fasts_completed || 0}
                    totalHours={Math.round(stats?.total_fasting_hours || 0)}
                />

                {/* Daily Progress Rings */}
                <motion.div
                    className="bg-card border border-border rounded-2xl p-6"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.1 }}
                >
                    <h3 className="font-semibold text-foreground mb-6 flex items-center gap-2">
                        <TrendingUp className="w-5 h-5 text-primary" />
                        Today's Progress
                    </h3>

                    <div className="grid grid-cols-4 gap-4">
                        <ProgressRing
                            value={Math.round(currentFastHours * 10) / 10}
                            max={fastingGoalHours}
                            color="#f59e0b"
                            icon={<Timer className="w-5 h-5 text-amber-400" />}
                            label="Fasting"
                            unit="hours"
                        />
                        <ProgressRing
                            value={dailyHydration?.glasses_count || 0}
                            max={8}
                            color="#06b6d4"
                            icon={<Droplets className="w-5 h-5 text-cyan-400" />}
                            label="Water"
                            unit="glasses"
                        />
                        <div
                            className="cursor-pointer"
                            onClick={() => setShowStepsInput(true)}
                        >
                            <ProgressRing
                                value={stepsToday}
                                max={10000}
                                color="#10b981"
                                icon={<Footprints className="w-5 h-5 text-green-400" />}
                                label="Steps"
                                unit="steps"
                            />
                        </div>
                        <ProgressRing
                            value={todayCalories}
                            max={1800}
                            color="#f97316"
                            icon={<Utensils className="w-5 h-5 text-orange-400" />}
                            label="Calories"
                            unit="cal"
                        />
                    </div>

                    {/* Steps Input Modal */}
                    {showStepsInput && (
                        <motion.div
                            className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50 flex items-center justify-center p-4"
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            onClick={() => setShowStepsInput(false)}
                        >
                            <motion.form
                                className="bg-card border border-border rounded-2xl p-6 w-full max-w-sm"
                                initial={{ scale: 0.9, opacity: 0 }}
                                animate={{ scale: 1, opacity: 1 }}
                                onClick={e => e.stopPropagation()}
                                onSubmit={handleStepsSubmit}
                            >
                                <h4 className="font-semibold text-foreground mb-4">Log Today's Steps</h4>
                                <input
                                    type="number"
                                    placeholder="Enter steps..."
                                    value={stepsToday || ''}
                                    onChange={e => setStepsToday(parseInt(e.target.value) || 0)}
                                    className="w-full px-4 py-3 rounded-lg bg-muted border border-border text-foreground mb-4 focus:outline-none focus:ring-2 focus:ring-primary"
                                    autoFocus
                                />
                                <div className="flex gap-2">
                                    <Button
                                        type="button"
                                        variant="outline"
                                        className="flex-1"
                                        onClick={() => setShowStepsInput(false)}
                                    >
                                        Cancel
                                    </Button>
                                    <Button type="submit" className="flex-1">Save</Button>
                                </div>
                            </motion.form>
                        </motion.div>
                    )}
                </motion.div>

                {/* Achievements Section */}
                <Achievements
                    fastingHistory={fastingHistory}
                    waterStreak={dailyHydration?.glasses_count && dailyHydration.glasses_count >= 6 ? 1 : 0}
                    stepsToday={stepsToday}
                    caloriesLogged={todayCalories}
                />

                {/* Weekly Fasting Chart */}
                <motion.div
                    className="bg-card border border-border rounded-2xl p-6"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.3 }}
                >
                    <div className="flex items-center justify-between mb-6">
                        <h3 className="font-semibold text-foreground flex items-center gap-2">
                            <BarChart3 className="w-5 h-5 text-secondary" />
                            Weekly Fasting Summary
                        </h3>
                    </div>

                    <div className="h-[200px] w-full">
                        <ResponsiveContainer width="100%" height="100%">
                            <BarChart data={fastingChartData}>
                                <defs>
                                    <linearGradient id="barGradient" x1="0" y1="0" x2="0" y2="1">
                                        <stop offset="0%" stopColor="#10b981" />
                                        <stop offset="100%" stopColor="#059669" />
                                    </linearGradient>
                                </defs>
                                <CartesianGrid strokeDasharray="3 3" stroke="#333" vertical={false} />
                                <XAxis
                                    dataKey="date"
                                    stroke="#666"
                                    fontSize={12}
                                    tickLine={false}
                                    axisLine={false}
                                />
                                <YAxis
                                    stroke="#666"
                                    fontSize={12}
                                    tickLine={false}
                                    axisLine={false}
                                    unit="h"
                                />
                                <Tooltip
                                    contentStyle={{
                                        backgroundColor: '#1a1a1a',
                                        border: '1px solid #333',
                                        borderRadius: '8px'
                                    }}
                                    itemStyle={{ color: '#fff' }}
                                    formatter={(value: number) => [`${value} hours`, 'Fasted']}
                                />
                                <Bar
                                    dataKey="hours"
                                    fill="url(#barGradient)"
                                    radius={[4, 4, 0, 0]}
                                />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                </motion.div>

                {/* Weekly AI Report */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.4 }}
                >
                    <WeeklyReportCard />
                </motion.div>
            </main>
        </div>
    );
};

export default Progress;
