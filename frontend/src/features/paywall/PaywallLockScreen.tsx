import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Lock } from "lucide-react";

interface PaywallLockScreenProps {
    onJoinVault: () => void;
    potentialEarnings?: number;
    completedFasts?: number;
    activeUsers?: number;
}

export const PaywallLockScreen = ({
    onJoinVault,
    potentialEarnings = 14.50,
    completedFasts = 7,
    activeUsers = 12453
}: PaywallLockScreenProps) => {
    return (
        <div className="min-h-screen flex items-center justify-center p-4 bg-slate-950">
            <Card className="max-w-lg w-full border-2 border-yellow-500/50 shadow-2xl shadow-yellow-500/20 bg-slate-900">
                <CardContent className="p-12 text-center">
                    {/* Lock Icon */}
                    <div className="mb-6">
                        <div className="mx-auto w-20 h-20 bg-yellow-500/20 rounded-full flex items-center justify-center mb-4 animate-pulse">
                            <Lock className="h-10 w-10 text-yellow-500" />
                        </div>
                        <h2 className="text-3xl font-bold mb-2 text-white">Free Trial Complete</h2>
                        <p className="text-slate-400">
                            You've proven you can fast. Now let's prove you can get paid for it.
                        </p>
                    </div>

                    {/* What They Would Have Earned */}
                    <Card className="bg-gradient-to-r from-emerald-500/20 to-purple-500/20 border-emerald-500/30 mb-6">
                        <CardContent className="p-6">
                            <p className="text-sm text-slate-400 mb-2">
                                You would have earned
                            </p>
                            <p className="text-6xl font-bold text-yellow-400 mb-2">
                                ${potentialEarnings.toFixed(2)}
                            </p>
                            <p className="text-xs text-slate-500">
                                in the last 7 days if you had the Vault
                            </p>

                            {/* Breakdown */}
                            <div className="mt-4 pt-4 border-t border-slate-700 space-y-2 text-sm">
                                <div className="flex justify-between">
                                    <span className="text-slate-400">Fasts completed:</span>
                                    <span className="font-semibold text-white">{completedFasts}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-slate-400">Avg earning per fast:</span>
                                    <span className="font-semibold text-emerald-400">
                                        ${(potentialEarnings / completedFasts).toFixed(2)}
                                    </span>
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Social Proof */}
                    <div className="mb-6 flex items-center justify-center gap-3">
                        <div className="flex -space-x-2">
                            {[...Array(5)].map((_, i) => (
                                <Avatar key={i} className="border-2 border-slate-900 h-8 w-8">
                                    <AvatarFallback className="bg-slate-700 text-xs text-slate-300">U{i}</AvatarFallback>
                                </Avatar>
                            ))}
                        </div>
                        <p className="text-sm text-slate-400">
                            <strong className="text-white">{activeUsers.toLocaleString()}</strong> people earning
                        </p>
                    </div>

                    {/* CTA */}
                    <Button
                        size="lg"
                        className="w-full h-14 text-lg bg-gradient-to-r from-emerald-500 to-purple-500 hover:from-emerald-600 hover:to-purple-600 text-white font-bold border-0"
                        onClick={onJoinVault}
                    >
                        Join Vault - Start Earning Today
                    </Button>

                    {/* Fine Print */}
                    <p className="text-xs text-slate-500 mt-4">
                        $20/month • Cancel anytime • Full refund if you maintain discipline
                    </p>
                </CardContent>
            </Card>
        </div>
    );
};
