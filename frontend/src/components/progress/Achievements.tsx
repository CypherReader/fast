import { motion } from 'framer-motion';
import { Flame, Droplets, Footprints, Trophy, Star, Zap } from 'lucide-react';
import { FastingSession } from '@/hooks/use-fasting';

interface AchievementsProps {
    fastingHistory?: FastingSession[];
    waterStreak?: number;
    stepsToday?: number;
    caloriesLogged?: number;
}

interface Achievement {
    id: string;
    title: string;
    description: string;
    icon: React.ReactNode;
    color: string;
    bgColor: string;
    unlocked: boolean;
}

export const Achievements = ({
    fastingHistory = [],
    waterStreak = 0,
    stepsToday = 0,
    caloriesLogged = 0,
}: AchievementsProps) => {
    // Calculate longest fast this week
    const now = new Date();
    const oneWeekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);

    const thisWeekFasts = fastingHistory.filter(fast => {
        if (!fast.end_time || fast.status !== 'completed') return false;
        return new Date(fast.end_time) >= oneWeekAgo;
    });

    const longestFastHours = thisWeekFasts.reduce((max, fast) => {
        const hours = (new Date(fast.end_time!).getTime() - new Date(fast.start_time).getTime()) / (1000 * 60 * 60);
        return Math.max(max, hours);
    }, 0);

    const achievements: Achievement[] = [
        {
            id: 'longest_fast',
            title: `Longest Fast: ${Math.floor(longestFastHours)}h`,
            description: 'This week',
            icon: <Flame className="w-5 h-5" />,
            color: 'text-orange-400',
            bgColor: 'bg-orange-500/20',
            unlocked: longestFastHours >= 12,
        },
        {
            id: 'water_streak',
            title: `Water Streak: ${waterStreak} days`,
            description: 'Keep it up!',
            icon: <Droplets className="w-5 h-5" />,
            color: 'text-cyan-400',
            bgColor: 'bg-cyan-500/20',
            unlocked: waterStreak >= 1,
        },
        {
            id: 'step_goal',
            title: stepsToday >= 10000 ? 'Step Goal Met!' : `${Math.floor(stepsToday / 100) / 10}k Steps`,
            description: stepsToday >= 10000 ? '10,000+ steps' : 'Keep moving!',
            icon: <Footprints className="w-5 h-5" />,
            color: 'text-green-400',
            bgColor: 'bg-green-500/20',
            unlocked: stepsToday >= 5000,
        },
        {
            id: 'calorie_tracker',
            title: caloriesLogged > 0 ? 'Calories Logged' : 'Log Your Food',
            description: caloriesLogged > 0 ? `${caloriesLogged} cal today` : 'Track your meals',
            icon: <Zap className="w-5 h-5" />,
            color: 'text-amber-400',
            bgColor: 'bg-amber-500/20',
            unlocked: caloriesLogged > 0,
        },
    ];

    const unlockedCount = achievements.filter(a => a.unlocked).length;

    return (
        <motion.div
            className="space-y-4"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
        >
            <div className="flex items-center justify-between">
                <h3 className="font-semibold text-foreground flex items-center gap-2">
                    <Trophy className="w-5 h-5 text-amber-400" />
                    Today's Achievements
                </h3>
                <span className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded-full">
                    {unlockedCount}/{achievements.length} unlocked
                </span>
            </div>

            <div className="grid grid-cols-2 gap-3">
                {achievements.map((achievement, index) => (
                    <motion.div
                        key={achievement.id}
                        className={`p-4 rounded-xl border ${achievement.unlocked
                                ? 'bg-card border-border'
                                : 'bg-muted/30 border-transparent opacity-50'
                            }`}
                        initial={{ opacity: 0, scale: 0.9 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ delay: index * 0.1 }}
                    >
                        <div className="flex items-start gap-3">
                            <div className={`p-2 rounded-lg ${achievement.bgColor} ${achievement.color}`}>
                                {achievement.icon}
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className={`text-sm font-medium truncate ${achievement.unlocked ? 'text-foreground' : 'text-muted-foreground'}`}>
                                    {achievement.title}
                                </p>
                                <p className="text-xs text-muted-foreground">{achievement.description}</p>
                            </div>
                            {achievement.unlocked && (
                                <Star className="w-4 h-4 text-amber-400 fill-amber-400 shrink-0" />
                            )}
                        </div>
                    </motion.div>
                ))}
            </div>
        </motion.div>
    );
};

export default Achievements;
