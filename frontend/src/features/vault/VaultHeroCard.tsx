import { Lock, TrendingUp, DollarSign } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";

interface VaultHeroCardProps {
    earned: number;
    deposit: number;
    daysUntilRefund: number;
    fastsRemaining: number;
}

export const VaultHeroCard = ({ earned, deposit, daysUntilRefund, fastsRemaining }: VaultHeroCardProps) => {
    const progress = (earned / deposit) * 100;

    return (
        <Card className="bg-gradient-to-br from-yellow-900/40 to-slate-900 border-yellow-500/30 overflow-hidden relative group">
            <div className="absolute top-0 right-0 p-3 opacity-10 group-hover:opacity-20 transition-opacity">
                <DollarSign className="w-24 h-24 text-yellow-500" />
            </div>

            <CardContent className="p-6">
                <div className="flex justify-between items-start mb-4">
                    <div>
                        <div className="flex items-center space-x-2 mb-1">
                            <span className="bg-yellow-500/20 text-yellow-400 text-xs font-bold px-2 py-1 rounded border border-yellow-500/30 flex items-center">
                                <Lock className="w-3 h-3 mr-1" /> VAULT BALANCE
                            </span>
                        </div>
                        <h3 className="text-3xl font-bold text-white tracking-tight flex items-baseline">
                            ${earned.toFixed(2)}
                            <span className="text-sm text-slate-400 font-normal ml-2">/ ${deposit.toFixed(2)}</span>
                        </h3>
                    </div>
                    <div className="text-right">
                        <div className="text-xs text-yellow-200/70 mb-1">Payout in</div>
                        <div className="font-mono text-xl font-bold text-yellow-400">{daysUntilRefund}d</div>
                    </div>
                </div>

                <div className="space-y-2">
                    <div className="flex justify-between text-xs text-slate-300">
                        <span>Progress to Refund</span>
                        <span className="text-yellow-400 font-bold">{progress.toFixed(0)}%</span>
                    </div>
                    <Progress value={progress} className="h-3 bg-slate-800" indicatorClassName="bg-gradient-to-r from-yellow-600 to-yellow-400" />
                    <p className="text-xs text-slate-400 mt-2 flex items-center">
                        <TrendingUp className="w-3 h-3 mr-1 text-emerald-400" />
                        Complete <span className="text-white font-bold mx-1">{fastsRemaining} more fasts</span> to unlock full payout.
                    </p>
                </div>
            </CardContent>
        </Card>
    );
};
