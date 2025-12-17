import { useState, useEffect } from 'react';
import { AlertTriangle, Flame, Clock, Zap } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { motion, AnimatePresence } from 'framer-motion';
import { api } from '@/api/client';
import { useFasting } from '@/hooks/use-fasting';

interface StreakRiskData {
    is_at_risk: boolean;
    current_streak: number;
    days_since_last_fast: number;
    hours_until_loss: number;
    urgency_level: 'none' | 'warning' | 'critical';
    ai_message: string;
    suggested_action: string;
    motivational_fact: string;
}

export const StreakAlert = () => {
    const [streakRisk, setStreakRisk] = useState<StreakRiskData | null>(null);
    const [isDismissed, setIsDismissed] = useState(false);
    const { startFast } = useFasting();

    useEffect(() => {
        const checkStreak = async () => {
            try {
                const response = await api.get<StreakRiskData>('/fasting/streak-risk');
                setStreakRisk(response.data);
            } catch (error) {
                console.error('Failed to check streak:', error);
            }
        };

        checkStreak();
        // Check every hour
        const interval = setInterval(checkStreak, 60 * 60 * 1000);
        return () => clearInterval(interval);
    }, []);

    if (!streakRisk || !streakRisk.is_at_risk || isDismissed) {
        return null;
    }

    const isCritical = streakRisk.urgency_level === 'critical';
    const bgColor = isCritical ? 'bg-destructive/10' : 'bg-yellow-500/10';
    const borderColor = isCritical ? 'border-destructive' : 'border-yellow-500';
    const textColor = isCritical ? 'text-destructive' : 'text-yellow-600';
    const icon = isCritical ? AlertTriangle : Flame;

    const handleEmergencyFast = () => {
        startFast({
            plan_type: '16:8',
            goal_hours: 16,
            start_time: new Date().toISOString()
        });
        setIsDismissed(true);
    };

    return (
        <AnimatePresence>
            <motion.div
                initial={{ opacity: 0, y: -20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                className="fixed top-20 left-1/2 transform -translate-x-1/2 z-50 w-full max-w-md px-4"
            >
                <Card className={`${bgColor} border-2 ${borderColor} p-4 shadow-lg`}>
                    <div className="flex items-start gap-3">
                        {React.createElement(icon, {
                            className: `w-6 h-6 ${textColor} flex-shrink-0 mt-0.5 ${isCritical ? 'animate-pulse' : ''}`
                        })}
                        <div className="flex-1">
                            <h4 className={`font-bold ${textColor} mb-1`}>
                                {isCritical ? '🚨 CRITICAL: Streak Ending Soon!' : '⚠️ Streak At Risk'}
                            </h4>
                            <div className="space-y-2 text-sm">
                                <div className="flex items-center gap-2 text-muted-foreground">
                                    <Flame className="w-4 h-4" />
                                    <span>{streakRisk.current_streak}-day streak</span>
                                    <Clock className="w-4 h-4 ml-2" />
                                    <span className="font-semibold">{streakRisk.hours_until_loss.toFixed(1)}h left</span>
                                </div>
                                <p className="font-medium">{streakRisk.ai_message}</p>
                                <p className="text-xs italic text-muted-foreground">
                                    {streakRisk.motivational_fact}
                                </p>
                            </div>
                            <div className="flex gap-2 mt-3">
                                <Button
                                    onClick={handleEmergencyFast}
                                    size="sm"
                                    className={isCritical ? 'bg-destructive hover:bg-destructive/90' : ''}
                                >
                                    <Zap className="w-4 h-4 mr-1" />
                                    {streakRisk.suggested_action}
                                </Button>
                                <Button
                                    onClick={() => setIsDismissed(true)}
                                    size="sm"
                                    variant="ghost"
                                >
                                    Dismiss
                                </Button>
                            </div>
                        </div>
                    </div>
                </Card>
            </motion.div>
        </AnimatePresence>
    );
};
