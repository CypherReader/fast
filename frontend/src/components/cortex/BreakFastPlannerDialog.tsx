import { useState } from 'react';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { useBreakFastPlanner } from '@/hooks/use-breakfast-planner';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { ScrollArea } from '@/components/ui/scroll-area';
import {
    UtensilsCrossed,
    Droplets,
    Clock,
    AlertCircle,
    CheckCircle,
    Info,
    Sparkles,
} from 'lucide-react';
import { Skeleton } from '@/components/ui/skeleton';

interface BreakFastPlannerProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    fastDuration: number; // Duration in hours
}

export const BreakFastPlannerDialog = ({
    open,
    onOpenChange,
    fastDuration,
}: BreakFastPlannerProps) => {
    const { guide, isLoading } = useBreakFastPlanner(fastDuration, open);

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-2xl max-h-[90vh] p-0">
                <DialogHeader className="p-6 pb-4">
                    <div className="flex items-center gap-3 mb-2">
                        <div className="p-2 bg-emerald-500/20 rounded-lg">
                            <UtensilsCrossed className="h-6 w-6 text-emerald-400" />
                        </div>
                        <div>
                            <DialogTitle className="text-2xl">Break-Fast Guide</DialogTitle>
                            <DialogDescription>
                                Personalized recommendations for ending your{' '}
                                {fastDuration.toFixed(1)}-hour fast
                            </DialogDescription>
                        </div>
                    </div>
                </DialogHeader>

                <ScrollArea className="max-h-[calc(90vh-120px)]">
                    <div className="px-6 pb-6 space-y-4">
                        {isLoading ? (
                            <LoadingSkeleton />
                        ) : guide ? (
                            <>
                                {/* AI Guidance - Featured at top */}
                                <Card className="bg-gradient-to-br from-purple-500/10 to-pink-500/10 border-purple-500/30">
                                    <CardContent className="p-4">
                                        <div className="flex items-start gap-3">
                                            <Sparkles className="h-5 w-5 text-purple-400 flex-shrink-0 mt-0.5" />
                                            <div>
                                                <h3 className="text-sm font-semibold text-purple-400 mb-1">
                                                    AI Personalized Tip
                                                </h3>
                                                <p className="text-sm leading-relaxed">{guide.ai_guidance}</p>
                                            </div>
                                        </div>
                                    </CardContent>
                                </Card>

                                {/* Meal Type & Portion Size */}
                                <div className="grid grid-cols-2 gap-3">
                                    <Card>
                                        <CardContent className="p-4">
                                            <div className="flex items-center gap-2 mb-1">
                                                <UtensilsCrossed className="h-4 w-4 text-emerald-400" />
                                                <h3 className="text-xs font-semibold text-muted-foreground uppercase">
                                                    Meal Type
                                                </h3>
                                            </div>
                                            <p className="text-sm font-semibold">{guide.meal_type}</p>
                                        </CardContent>
                                    </Card>

                                    <Card>
                                        <CardContent className="p-4">
                                            <div className="flex items-center gap-2 mb-1">
                                                <Info className="h-4 w-4 text-blue-400" />
                                                <h3 className="text-xs font-semibold text-muted-foreground uppercase">
                                                    Portion Size
                                                </h3>
                                            </div>
                                            <p className="text-sm font-semibold">{guide.portion_size}</p>
                                        </CardContent>
                                    </Card>
                                </div>

                                {/* Recommended Foods */}
                                <Card className="border-emerald-500/30">
                                    <CardContent className="p-4">
                                        <div className="flex items-center gap-2 mb-3">
                                            <CheckCircle className="h-5 w-5 text-emerald-400" />
                                            <h3 className="text-sm font-semibold text-emerald-400">
                                                Recommended Foods
                                            </h3>
                                        </div>
                                        <div className="flex flex-wrap gap-2">
                                            {guide.recommended_foods.map((food, idx) => (
                                                <Badge
                                                    key={idx}
                                                    variant="outline"
                                                    className="bg-emerald-500/10 border-emerald-500/30 text-emerald-400"
                                                >
                                                    {food}
                                                </Badge>
                                            ))}
                                        </div>
                                    </CardContent>
                                </Card>

                                {/* Foods to Avoid */}
                                <Card className="border-red-500/30">
                                    <CardContent className="p-4">
                                        <div className="flex items-center gap-2 mb-3">
                                            <AlertCircle className="h-5 w-5 text-red-400" />
                                            <h3 className="text-sm font-semibold text-red-400">
                                                Foods to Avoid
                                            </h3>
                                        </div>
                                        <div className="flex flex-wrap gap-2">
                                            {guide.foods_to_avoid.map((food, idx) => (
                                                <Badge
                                                    key={idx}
                                                    variant="outline"
                                                    className="bg-red-500/10 border-red-500/30 text-red-400"
                                                >
                                                    {food}
                                                </Badge>
                                            ))}
                                        </div>
                                    </CardContent>
                                </Card>

                                <Separator />

                                {/* Detailed Guidelines */}
                                <div className="space-y-4">
                                    {/* Hydration Tip */}
                                    <div className="flex items-start gap-3">
                                        <div className="p-2 bg-blue-500/20 rounded-lg flex-shrink-0">
                                            <Droplets className="h-5 w-5 text-blue-400" />
                                        </div>
                                        <div className="flex-1">
                                            <h3 className="text-sm font-semibold mb-1">Hydration First</h3>
                                            <p className="text-sm text-muted-foreground">
                                                {guide.hydration_tip}
                                            </p>
                                        </div>
                                    </div>

                                    {/* Timing Advice */}
                                    <div className="flex items-start gap-3">
                                        <div className="p-2 bg-orange-500/20 rounded-lg flex-shrink-0">
                                            <Clock className="h-5 w-5 text-orange-400" />
                                        </div>
                                        <div className="flex-1">
                                            <h3 className="text-sm font-semibold mb-1">Timing & Schedule</h3>
                                            <p className="text-sm text-muted-foreground">
                                                {guide.timing_advice}
                                            </p>
                                        </div>
                                    </div>

                                    {/* Reintroduction Plan */}
                                    <div className="flex items-start gap-3">
                                        <div className="p-2 bg-purple-500/20 rounded-lg flex-shrink-0">
                                            <UtensilsCrossed className="h-5 w-5 text-purple-400" />
                                        </div>
                                        <div className="flex-1">
                                            <h3 className="text-sm font-semibold mb-1">
                                                Reintroduction Plan
                                            </h3>
                                            <p className="text-sm text-muted-foreground">
                                                {guide.reintroduction_plan}
                                            </p>
                                        </div>
                                    </div>
                                </div>

                                {/* Footer Note */}
                                <Card className="bg-slate-900/50 border-slate-800">
                                    <CardContent className="p-4">
                                        <p className="text-xs text-muted-foreground text-center">
                                            💡 Listen to your body. If you experience digestive discomfort,
                                            slow down and choose lighter foods. These are general guidelines
                                            tailored to your fast duration.
                                        </p>
                                    </CardContent>
                                </Card>
                            </>
                        ) : (
                            <div className="text-center py-8">
                                <p className="text-muted-foreground">
                                    Unable to load recommendations. Please try again.
                                </p>
                            </div>
                        )}
                    </div>
                </ScrollArea>
            </DialogContent>
        </Dialog>
    );
};

const LoadingSkeleton = () => (
    <div className="space-y-4">
        <Skeleton className="h-24 w-full" />
        <div className="grid grid-cols-2 gap-3">
            <Skeleton className="h-20 w-full" />
            <Skeleton className="h-20 w-full" />
        </div>
        <Skeleton className="h-32 w-full" />
        <Skeleton className="h-32 w-full" />
        <div className="space-y-3">
            <Skeleton className="h-16 w-full" />
            <Skeleton className="h-16 w-full" />
            <Skeleton className="h-16 w-full" />
        </div>
    </div>
);
