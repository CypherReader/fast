import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import { Lock, ArrowLeft, TrendingUp, Calendar, Target, Scale } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useUser } from '@/hooks/use-user';
import UserMenu from '@/components/layout/UserMenu';
import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    AreaChart,
    Area
} from 'recharts';
import { useProgress } from '@/hooks/use-progress';
import { format, parseISO } from 'date-fns';

const Progress = () => {
    const navigate = useNavigate();
    const { user, stats } = useUser();
    const { weightHistory } = useProgress();

    // Transform weight history to chart format
    const weightData = weightHistory && weightHistory.length > 0
        ? [...weightHistory].reverse().map(entry => ({
            date: format(parseISO(entry.logged_at), 'EEE'),
            weight: entry.weight_lbs
        }))
        : [];

    const fastingData = [
        { date: 'Mon', hours: 16 },
        { date: 'Tue', hours: 16.5 },
        { date: 'Wed', hours: 15.5 },
        { date: 'Thu', hours: 17 },
        { date: 'Fri', hours: 16 },
        { date: 'Sat', hours: 18 },
        { date: 'Sun', hours: 16 },
    ];

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
                {/* Stats Overview */}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    {[
                        { label: 'Current Streak', value: stats?.current_streak || 0, icon: TrendingUp, unit: 'days' },
                        { label: 'Total Fasts', value: stats?.fasts_completed || 0, icon: Calendar, unit: 'fasts' },
                        { label: 'Total Hours', value: stats?.total_fasting_hours || 0, icon: Target, unit: 'hours' },
                        { label: 'Vault Balance', value: `$${stats?.vault_balance || 0}`, icon: Lock, unit: 'saved' },
                    ].map((stat, index) => (
                        <motion.div
                            key={stat.label}
                            className="bg-card border border-border rounded-xl p-4"
                            initial={{ opacity: 0, y: 20 }}
                            animate={{ opacity: 1, y: 0 }}
                            transition={{ delay: index * 0.1 }}
                        >
                            <stat.icon className="w-5 h-5 text-muted-foreground mb-2" />
                            <div className="text-2xl font-bold text-foreground">
                                {stat.value}
                                <span className="text-xs font-normal text-muted-foreground ml-1">{stat.unit}</span>
                            </div>
                            <div className="text-xs text-muted-foreground">{stat.label}</div>
                        </motion.div>
                    ))}
                </div>

                {/* Weight Chart */}
                <motion.div
                    className="bg-card border border-border rounded-2xl p-6"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.2 }}
                >
                    <div className="flex items-center justify-between mb-6">
                        <div className="flex items-center gap-2">
                            <Scale className="w-5 h-5 text-primary" />
                            <h3 className="font-semibold text-foreground">Weight Trend</h3>
                        </div>
                        <select className="bg-muted/50 border-none text-sm rounded-lg px-3 py-1">
                            <option>Last 7 Days</option>
                            <option>Last 30 Days</option>
                            <option>Last 3 Months</option>
                        </select>
                    </div>

                    <div className="h-[300px] w-full">
                        {weightData.length === 0 ? (
                            <div className="h-full flex flex-col items-center justify-center text-center p-6">
                                <Scale className="w-12 h-12 text-muted-foreground/30 mb-3" />
                                <p className="text-foreground font-medium mb-1">No weight data yet</p>
                                <p className="text-sm text-muted-foreground">Log your weight to see your progress chart!</p>
                            </div>
                        ) : (
                            <ResponsiveContainer width="100%" height="100%">
                                <LineChart data={weightData}>
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
                                        domain={['dataMin - 1', 'dataMax + 1']}
                                    />
                                    <Tooltip
                                        contentStyle={{ backgroundColor: '#1a1a1a', border: '1px solid #333', borderRadius: '8px' }}
                                        itemStyle={{ color: '#fff' }}
                                    />
                                    <Line
                                        type="monotone"
                                        dataKey="weight"
                                        stroke="#f59e0b"
                                        strokeWidth={3}
                                        dot={{ fill: '#f59e0b', strokeWidth: 2 }}
                                        activeDot={{ r: 6 }}
                                    />
                                </LineChart>
                            </ResponsiveContainer>
                        )}
                    </div>
                </motion.div>

                {/* Fasting Consistency Chart */}
                <motion.div
                    className="bg-card border border-border rounded-2xl p-6"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.3 }}
                >
                    <div className="flex items-center justify-between mb-6">
                        <div className="flex items-center gap-2">
                            <Target className="w-5 h-5 text-secondary" />
                            <h3 className="font-semibold text-foreground">Fasting Consistency</h3>
                        </div>
                    </div>

                    <div className="h-[200px] w-full">
                        <ResponsiveContainer width="100%" height="100%">
                            <AreaChart data={fastingData}>
                                <defs>
                                    <linearGradient id="colorHours" x1="0" y1="0" x2="0" y2="1">
                                        <stop offset="5%" stopColor="#10b981" stopOpacity={0.3} />
                                        <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
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
                                />
                                <Tooltip
                                    contentStyle={{ backgroundColor: '#1a1a1a', border: '1px solid #333', borderRadius: '8px' }}
                                    itemStyle={{ color: '#fff' }}
                                />
                                <Area
                                    type="monotone"
                                    dataKey="hours"
                                    stroke="#10b981"
                                    fillOpacity={1}
                                    fill="url(#colorHours)"
                                />
                            </AreaChart>
                        </ResponsiveContainer>
                    </div>
                </motion.div>
            </main>
        </div>
    );
};

export default Progress;
