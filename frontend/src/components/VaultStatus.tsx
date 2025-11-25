import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Lock, Unlock, DollarSign, TrendingUp } from "lucide-react";

interface VaultStatusProps {
    deposit: number;
    earned: number;
    potentialRefund: number;
}

const VaultStatus: React.FC<VaultStatusProps> = ({ deposit, earned, potentialRefund }) => {
    const progress = Math.min((earned / deposit) * 100, 100);
    const isFullyUnlocked = earned >= deposit;

    return (
        <Card className="w-full max-w-md border-emerald-500/30 bg-gradient-to-br from-slate-900 to-slate-950 shadow-xl overflow-hidden relative">
            {/* Background Glow */}
            <div className="absolute top-0 right-0 w-32 h-32 bg-emerald-500/10 rounded-full blur-3xl -mr-10 -mt-10" />

            <CardHeader className="pb-2">
                <div className="flex justify-between items-center">
                    <div className="flex items-center gap-2">
                        <div className={`p-2 rounded-lg ${isFullyUnlocked ? 'bg-emerald-500/20 text-emerald-400' : 'bg-slate-800 text-slate-400'}`}>
                            {isFullyUnlocked ? <Unlock className="w-5 h-5" /> : <Lock className="w-5 h-5" />}
                        </div>
                        <div>
                            <CardTitle className="text-lg font-bold text-white">Commitment Vault</CardTitle>
                            <p className="text-xs text-slate-400">Monthly Deposit: ${deposit.toFixed(2)}</p>
                        </div>
                    </div>
                    <div className="text-right">
                        <div className="text-2xl font-bold font-mono text-emerald-400">
                            ${potentialRefund.toFixed(2)}
                        </div>
                        <div className="text-xs text-emerald-500/80 font-medium">Unlocked Refund</div>
                    </div>
                </div>
            </CardHeader>

            <CardContent>
                <div className="space-y-4">
                    {/* Progress Bar */}
                    <div className="space-y-2">
                        <div className="flex justify-between text-xs">
                            <span className="text-slate-400">Progress to Max Refund</span>
                            <span className="text-white font-mono">{progress.toFixed(0)}%</span>
                        </div>
                        <Progress value={progress} className="h-2 bg-slate-800" indicatorClassName="bg-gradient-to-r from-emerald-500 to-teal-400" />
                    </div>

                    {/* Stats Grid */}
                    <div className="grid grid-cols-2 gap-3 pt-2">
                        <div className="bg-slate-900/50 p-3 rounded-lg border border-slate-800">
                            <div className="flex items-center gap-2 mb-1 text-slate-400 text-xs">
                                <DollarSign className="w-3 h-3" />
                                <span>Net Cost</span>
                            </div>
                            <div className="font-mono font-bold text-white">
                                ${(30 - potentialRefund).toFixed(2)}
                            </div>
                        </div>
                        <div className="bg-slate-900/50 p-3 rounded-lg border border-slate-800">
                            <div className="flex items-center gap-2 mb-1 text-slate-400 text-xs">
                                <TrendingUp className="w-3 h-3" />
                                <span>Next Goal</span>
                            </div>
                            <div className="text-xs text-white">
                                +$0.50 (Log Dinner)
                            </div>
                        </div>
                    </div>

                    <div className="text-xs text-center text-slate-500 mt-2">
                        {isFullyUnlocked
                            ? "🎉 Maximum refund unlocked! Great discipline."
                            : "Keep consistent to unlock your full deposit."}
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};

export default VaultStatus;
