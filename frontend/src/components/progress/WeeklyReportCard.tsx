import { motion } from 'framer-motion';
import {
    Calendar,
    Clock,
    Flame,
    Trophy,
    TrendingUp,
    Sparkles,
    Target,
    ChevronRight,
    Lightbulb
} from 'lucide-react';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useWeeklyReport } from '@/hooks/use-weekly-report';
import { cn } from '@/lib/utils';

// Stat card component
const StatCard = ({
    icon: Icon,
    label,
    value,
    sublabel,
    color = 'text-primary',
    bgColor = 'bg-primary/10'
}: {
    icon: React.ElementType;
    label: string;
    value: string | number;
    sublabel?: string;
    color?: string;
    bgColor?: string;
}) => (
    <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex items-center gap-3 p-4 rounded-xl bg-card/50 border border-border"
    >
        <div className={cn('p-2 rounded-lg', bgColor)}>
            <Icon className={cn('w-5 h-5', color)} />
        </div>
        <div className="flex-1">
            <p className="text-2xl font-bold text-foreground">{value}</p>
            <p className="text-sm text-muted-foreground">{label}</p>
            {sublabel && <p className="text-xs text-muted-foreground">{sublabel}</p>}
        </div>
    </motion.div>
);

// Progress bar component
const ProgressBar = ({ value, max, label, color }: {
    value: number;
    max: number;
    label: string;
    color: string;
}) => {
    const percentage = Math.min((value / max) * 100, 100);
    return (
        <div className="space-y-1">
            <div className="flex justify-between text-xs text-muted-foreground">
                <span>{label}</span>
                <span>{value.toFixed(1)}h / {max}h</span>
            </div>
            <div className="h-2 bg-muted/30 rounded-full overflow-hidden">
                <motion.div
                    className={cn('h-full rounded-full', color)}
                    initial={{ width: 0 }}
                    animate={{ width: `${percentage}%` }}
                    transition={{ duration: 1, ease: 'easeOut' }}
                />
            </div>
        </div>
    );
};

export function WeeklyReportCard() {
    const { data: report, isLoading, isError } = useWeeklyReport();

    if (isLoading) {
        return (
            <Card className="p-6 animate-pulse">
                <div className="h-6 bg-muted rounded w-1/3 mb-4"></div>
                <div className="grid grid-cols-2 gap-4">
                    <div className="h-24 bg-muted rounded"></div>
                    <div className="h-24 bg-muted rounded"></div>
                    <div className="h-24 bg-muted rounded"></div>
                    <div className="h-24 bg-muted rounded"></div>
                </div>
            </Card>
        );
    }

    if (isError || !report) {
        return (
            <Card className="p-6">
                <div className="text-center py-8">
                    <Calendar className="w-12 h-12 text-muted-foreground mx-auto mb-3" />
                    <p className="text-muted-foreground">Complete your first fast to see your weekly report!</p>
                </div>
            </Card>
        );
    }

    const weekStart = new Date(report.week_start).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
    const weekEnd = new Date(report.week_end).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });

    return (
        <Card className="p-6 overflow-hidden relative">
            {/* Background decoration */}
            <div className="absolute top-0 right-0 w-64 h-64 bg-gradient-to-bl from-primary/5 to-transparent rounded-full -translate-y-1/2 translate-x-1/2" />

            {/* Header */}
            <div className="flex items-center justify-between mb-6 relative">
                <div>
                    <h2 className="text-xl font-bold text-foreground flex items-center gap-2">
                        <Calendar className="w-5 h-5 text-primary" />
                        Weekly Progress Report
                    </h2>
                    <p className="text-sm text-muted-foreground mt-1">
                        {weekStart} - {weekEnd}
                    </p>
                </div>
                <div className="flex items-center gap-2 text-xs bg-secondary/20 text-secondary px-3 py-1.5 rounded-full">
                    <Sparkles className="w-3 h-3" />
                    AI Powered
                </div>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-2 gap-3 mb-6">
                <StatCard
                    icon={Flame}
                    label="Fasts Completed"
                    value={report.fasts_completed}
                    color="text-orange-400"
                    bgColor="bg-orange-400/10"
                />
                <StatCard
                    icon={Clock}
                    label="Avg Duration"
                    value={`${report.average_duration.toFixed(1)}h`}
                    color="text-blue-400"
                    bgColor="bg-blue-400/10"
                />
                <StatCard
                    icon={Trophy}
                    label="Longest Fast"
                    value={`${report.longest_fast.toFixed(1)}h`}
                    sublabel="Personal best this week!"
                    color="text-amber-400"
                    bgColor="bg-amber-400/10"
                />
                <StatCard
                    icon={TrendingUp}
                    label="Total Fasting"
                    value={`${report.total_fasting_hours.toFixed(1)}h`}
                    color="text-purple-400"
                    bgColor="bg-purple-400/10"
                />
            </div>

            {/* Progress visualization */}
            <div className="space-y-4 mb-6">
                <h3 className="text-sm font-medium text-foreground">Weekly Progress</h3>
                <ProgressBar
                    value={report.total_fasting_hours}
                    max={112} // 16h * 7 days
                    label="Total fasting hours"
                    color="bg-gradient-to-r from-primary to-secondary"
                />
                <div className="grid grid-cols-2 gap-4">
                    <div className="p-3 rounded-lg bg-green-500/10 border border-green-500/20">
                        <p className="text-xs text-green-400 mb-1">Best Day</p>
                        <p className="font-semibold text-foreground">{report.best_day}</p>
                    </div>
                    <div className="p-3 rounded-lg bg-orange-500/10 border border-orange-500/20">
                        <p className="text-xs text-orange-400 mb-1">Challenge Day</p>
                        <p className="font-semibold text-foreground">{report.challenge_day}</p>
                    </div>
                </div>
            </div>

            {/* AI Insights */}
            <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ delay: 0.3 }}
                className="p-4 rounded-xl bg-gradient-to-br from-primary/5 via-secondary/5 to-primary/5 border border-primary/20 mb-6"
            >
                <div className="flex items-start gap-3">
                    <div className="p-2 rounded-lg bg-primary/10">
                        <Sparkles className="w-4 h-4 text-primary" />
                    </div>
                    <div className="flex-1">
                        <h4 className="font-medium text-foreground mb-2">AI Insights</h4>
                        <p className="text-sm text-muted-foreground leading-relaxed">
                            {report.ai_insights}
                        </p>
                    </div>
                </div>
            </motion.div>

            {/* Predictions */}
            <div className="grid grid-cols-3 gap-3 mb-6">
                <div className="text-center p-3 rounded-lg bg-muted/30">
                    <p className="text-2xl font-bold text-foreground">{report.predictions.next_week_fasts_estimate}</p>
                    <p className="text-xs text-muted-foreground">Predicted fasts</p>
                </div>
                <div className="text-center p-3 rounded-lg bg-muted/30">
                    <p className="text-2xl font-bold text-secondary">{report.predictions.success_probability}%</p>
                    <p className="text-xs text-muted-foreground">Success rate</p>
                </div>
                <div className="text-center p-3 rounded-lg bg-muted/30">
                    <p className="text-lg font-bold text-foreground capitalize">{report.predictions.discipline_trend}</p>
                    <p className="text-xs text-muted-foreground">Trend</p>
                </div>
            </div>

            {/* Recommendations */}
            {report.recommendations.length > 0 && (
                <div className="space-y-3">
                    <h3 className="text-sm font-medium text-foreground flex items-center gap-2">
                        <Lightbulb className="w-4 h-4 text-amber-400" />
                        Recommendations
                    </h3>
                    <div className="space-y-2">
                        {report.recommendations.slice(0, 3).map((rec, idx) => (
                            <motion.div
                                key={idx}
                                initial={{ opacity: 0, x: -20 }}
                                animate={{ opacity: 1, x: 0 }}
                                transition={{ delay: 0.4 + idx * 0.1 }}
                                className="flex items-center gap-3 p-3 rounded-lg bg-card border border-border hover:border-primary/50 transition-colors"
                            >
                                <ChevronRight className="w-4 h-4 text-primary flex-shrink-0" />
                                <span className="text-sm text-foreground">{rec}</span>
                            </motion.div>
                        ))}
                    </div>
                </div>
            )}

            {/* Goal Achievement */}
            {report.goal_achievement_date && (
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.6 }}
                    className="mt-6 p-4 rounded-xl bg-gradient-to-r from-amber-500/10 to-orange-500/10 border border-amber-500/20"
                >
                    <div className="flex items-center gap-3">
                        <Target className="w-5 h-5 text-amber-400" />
                        <div>
                            <p className="text-sm text-foreground font-medium">Projected Goal Date</p>
                            <p className="text-lg font-bold text-amber-400">{report.goal_achievement_date}</p>
                        </div>
                    </div>
                </motion.div>
            )}
        </Card>
    );
}
