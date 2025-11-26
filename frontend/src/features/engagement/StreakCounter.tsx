import { Flame, Snowflake } from "lucide-react";
import { cn } from "@/lib/utils";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

interface StreakCounterProps {
    streak: number;
    frozen?: boolean;
}

export const StreakCounter = ({ streak, frozen = false }: StreakCounterProps) => {
    return (
        <TooltipProvider>
            <Tooltip>
                <TooltipTrigger asChild>
                    <div className={cn(
                        "flex items-center space-x-1 px-3 py-1.5 rounded-full border transition-all cursor-help",
                        frozen
                            ? "bg-blue-900/20 border-blue-500/30 text-blue-400"
                            : "bg-orange-900/20 border-orange-500/30 text-orange-400"
                    )}>
                        {frozen ? (
                            <Snowflake className="w-4 h-4 animate-pulse" />
                        ) : (
                            <Flame className={cn("w-4 h-4", streak > 0 && "fill-orange-500 text-orange-600 animate-pulse")} />
                        )}
                        <span className="font-bold font-mono text-sm">{streak}</span>
                    </div>
                </TooltipTrigger>
                <TooltipContent className="bg-slate-900 border-slate-800 text-slate-200">
                    <p className="font-bold mb-1">
                        {frozen ? "Streak Frozen ❄️" : `${streak} Day Streak 🔥`}
                    </p>
                    <p className="text-xs text-slate-400">
                        {frozen
                            ? "Your streak is safe for today."
                            : "Fast today to keep it alive!"}
                    </p>
                </TooltipContent>
            </Tooltip>
        </TooltipProvider>
    );
};
