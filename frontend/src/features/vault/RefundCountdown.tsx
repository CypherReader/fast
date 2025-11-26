import { Clock } from "lucide-react";

interface RefundCountdownProps {
    daysRemaining: number;
}

export const RefundCountdown = ({ daysRemaining }: RefundCountdownProps) => {
    if (daysRemaining <= 0) return null;

    return (
        <div className="flex items-center space-x-2 bg-yellow-900/20 border border-yellow-500/30 px-3 py-1.5 rounded-full animate-pulse">
            <Clock className="w-3 h-3 text-yellow-500" />
            <span className="text-xs font-bold text-yellow-400 font-mono">
                {daysRemaining} days until payout
            </span>
        </div>
    );
};
