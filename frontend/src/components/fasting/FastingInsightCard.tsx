import { Card } from '../ui/card';
import { Badge } from '../ui/badge';
import { Sparkles, TrendingUp } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

interface FastingInsightCardProps {
    hours: number;
    milestone: string;
    insight: string;
    benefits: string[];
    motivation: string;
}

export function FastingInsightCard({
    hours,
    milestone,
    insight,
    benefits,
    motivation
}: FastingInsightCardProps) {
    return (
        <AnimatePresence>
            <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.4 }}
            >
                <Card className="bg-gradient-to-br from-blue-50 to-purple-50 dark:from-blue-950 dark:to-purple-950 border-2 border-purple-200 dark:border-purple-800">
                    <div className="p-5 space-y-4">
                        {/* Header with Milestone Badge */}
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <Sparkles className="w-5 h-5 text-purple-600" />
                                <h3 className="font-semibold text-lg">What's Happening Now</h3>
                            </div>
                            <Badge className="bg-purple-600 text-white px-3 py-1">
                                {milestone}
                            </Badge>
                        </div>

                        {/* AI Insight */}
                        <p className="text-sm text-gray-700 dark:text-gray-300 leading-relaxed">
                            {insight}
                        </p>

                        {/* Benefits Pills */}
                        {benefits && benefits.length > 0 && (
                            <div className="flex flex-wrap gap-2">
                                {benefits.map((benefit, index) => (
                                    <motion.div
                                        key={benefit}
                                        initial={{ scale: 0 }}
                                        animate={{ scale: 1 }}
                                        transition={{ delay: index * 0.1 }}
                                    >
                                        <Badge
                                            variant="secondary"
                                            className="bg-white/80 dark:bg-gray-800/80 text-purple-700 dark:text-purple-300 border border-purple-200 dark:border-purple-700"
                                        >
                                            <TrendingUp className="w-3 h-3 mr-1" />
                                            {benefit}
                                        </Badge>
                                    </motion.div>
                                ))}
                            </div>
                        )}

                        {/* Motivation Quote */}
                        {motivation && (
                            <div className="pt-3 border-t border-purple-200 dark:border-purple-800">
                                <p className="text-xs italic text-purple-700 dark:text-purple-300 text-center">
                                    "{motivation}"
                                </p>
                            </div>
                        )}

                        {/* Hours Display */}
                        <div className="text-xs text-gray-500 dark:text-gray-400 text-center">
                            {hours.toFixed(1)} hours into your fast
                        </div>
                    </div>
                </Card>
            </motion.div>
        </AnimatePresence>
    );
}
