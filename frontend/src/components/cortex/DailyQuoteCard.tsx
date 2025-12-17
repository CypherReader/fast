import { Card, CardContent } from '@/components/ui/card';
import { useDailyQuote } from '@/hooks/use-daily-quote';
import { Sparkles, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useQueryClient } from '@tantml:query/react-query';
import { Skeleton } from '@/components/ui/skeleton';

export const DailyQuoteCard = () => {
    const { quote, isLoading, error } = useDailyQuote();
    const queryClient = useQueryClient();

    const handleRefresh = () => {
        queryClient.invalidateQueries({ queryKey: ['daily-quote'] });
    };

    if (isLoading) {
        return (
            <Card className="bg-gradient-to-br from-purple-500/10 to-pink-500/10 border-purple-500/20">
                <CardContent className="p-6">
                    <div className="flex items-start gap-3">
                        <div className="p-2 bg-purple-500/20 rounded-lg flex-shrink-0">
                            <Sparkles className="h-5 w-5 text-purple-400" />
                        </div>
                        <div className="flex-1 space-y-2">
                            <Skeleton className="h-4 w-3/4" />
                            <Skeleton className="h-4 w-full" />
                        </div>
                    </div>
                </CardContent>
            </Card>
        );
    }

    if (error) {
        return (
            <Card className="bg-gradient-to-br from-purple-500/10 to-pink-500/10 border-purple-500/20">
                <CardContent className="p-6">
                    <div className="flex items-start gap-3">
                        <div className="p-2 bg-purple-500/20 rounded-lg flex-shrink-0">
                            <Sparkles className="h-5 w-5 text-purple-400" />
                        </div>
                        <div className="flex-1">
                            <p className="text-sm text-muted-foreground italic">
                                "Every hour of discipline builds the person you're becoming."
                            </p>
                        </div>
                    </div>
                </CardContent>
            </Card>
        );
    }

    return (
        <Card className="bg-gradient-to-br from-purple-500/10 to-pink-500/10 border-purple-500/20 hover:border-purple-500/30 transition-all">
            <CardContent className="p-6">
                <div className="flex items-start gap-3">
                    {/* Icon */}
                    <div className="p-2 bg-purple-500/20 rounded-lg flex-shrink-0 animate-pulse-slow">
                        <Sparkles className="h-5 w-5 text-purple-400" />
                    </div>

                    {/* Quote Content */}
                    <div className="flex-1">
                        <div className="flex items-start justify-between gap-2 mb-2">
                            <h3 className="text-xs font-semibold text-purple-400 uppercase tracking-wide">
                                Daily Motivation
                            </h3>
                            <Button
                                variant="ghost"
                                size="icon"
                                className="h-6 w-6 opacity-60 hover:opacity-100"
                                onClick={handleRefresh}
                                title="Get new quote"
                            >
                                <RefreshCw className="h-3 w-3" />
                            </Button>
                        </div>
                        <p className="text-sm leading-relaxed text-foreground/90 italic">
                            "{quote}"
                        </p>
                    </div>
                </div>

                {/* AI Attribution */}
                <div className="mt-3 pt-3 border-t border-purple-500/10">
                    <p className="text-xs text-muted-foreground text-center">
                        Personalized by Cortex AI
                    </p>
                </div>
            </CardContent>
        </Card>
    );
};
