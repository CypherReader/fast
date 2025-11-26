import { Dialog, DialogContent } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { TrendingUp, Users, Brain, Trophy, Shield, CreditCard } from "lucide-react";

interface VaultPromptDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onJoinVault: () => void;
    onContinueFree: () => void;
    discipline?: number;
    projectedFasts?: number;
}

export const VaultPromptDialog = ({
    open,
    onOpenChange,
    onJoinVault,
    onContinueFree,
    discipline = 72,
    projectedFasts = 12
}: VaultPromptDialogProps) => {
    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="max-w-lg bg-slate-950 border-slate-800 text-slate-200">
                <div className="text-center">
                    {/* Celebration Header */}
                    <div className="mb-6">
                        <div className="text-7xl mb-4 animate-bounce">🎯</div>
                        <h2 className="text-3xl font-bold mb-2 text-white">You're Crushing It!</h2>
                        <p className="text-slate-400">
                            You've completed 3 fasts. Ready to get paid for your discipline?
                        </p>
                    </div>

                    {/* Vault Value Proposition */}
                    <Card className="bg-slate-900 border-slate-800 mb-6">
                        <CardContent className="p-6">
                            <div className="text-left space-y-4">
                                {/* Deposit */}
                                <div className="flex items-center justify-between">
                                    <span className="text-slate-400">Monthly Deposit</span>
                                    <span className="text-2xl font-bold text-red-400">-$20.00</span>
                                </div>

                                {/* Earnings Potential */}
                                <div className="flex items-center justify-between">
                                    <div>
                                        <p className="text-slate-400">Earn Back</p>
                                        <p className="text-xs text-slate-500">
                                            (based on your current pace)
                                        </p>
                                    </div>
                                    <span className="text-2xl font-bold text-emerald-400">+$20.00</span>
                                </div>

                                <Separator className="bg-slate-800" />

                                {/* Net Cost */}
                                <div className="flex items-center justify-between">
                                    <span className="font-bold text-lg text-white">Your Real Cost</span>
                                    <span className="text-4xl font-bold text-emerald-400">$0.00</span>
                                </div>

                                {/* Explanation */}
                                <div className="bg-slate-800/50 rounded-lg p-3">
                                    <p className="text-xs text-slate-400">
                                        At your current discipline level ({discipline}/100), you're projected to complete
                                        <strong className="text-white"> {projectedFasts} fasts/month</strong>,
                                        earning back your full deposit.
                                    </p>
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    {/* Bonus Features List */}
                    <div className="text-left mb-6 space-y-2">
                        <p className="text-sm font-semibold mb-3 text-white">Vault Members Also Get:</p>
                        {[
                            { icon: TrendingUp, text: "Advanced analytics & trends" },
                            { icon: Users, text: "Join tribes for 2x accountability" },
                            { icon: Brain, text: "Unlimited AI coaching (Cortex)" },
                            { icon: Trophy, text: "Compete on global leaderboard" }
                        ].map((feature, idx) => (
                            <div key={idx} className="flex items-center gap-3 text-sm text-slate-300">
                                <feature.icon className="h-4 w-4 text-emerald-400 flex-shrink-0" />
                                <span>{feature.text}</span>
                            </div>
                        ))}
                    </div>

                    {/* CTA Buttons */}
                    <div className="space-y-3">
                        <Button
                            size="lg"
                            className="w-full bg-gradient-to-r from-emerald-500 to-purple-500 hover:from-emerald-600 hover:to-purple-600 text-lg h-14 font-bold text-white border-0"
                            onClick={onJoinVault}
                        >
                            Start Earning - $20/month
                        </Button>

                        <Button
                            variant="ghost"
                            size="sm"
                            className="w-full text-slate-500 hover:text-slate-400 hover:bg-slate-900"
                            onClick={onContinueFree}
                        >
                            Continue with limited features (free)
                        </Button>
                    </div>

                    {/* Trust Badges */}
                    <div className="mt-6 flex items-center justify-center gap-4 text-xs text-slate-500">
                        <div className="flex items-center gap-1">
                            <Shield className="h-4 w-4 text-emerald-400" />
                            <span>Cancel anytime</span>
                        </div>
                        <Separator orientation="vertical" className="h-4 bg-slate-800" />
                        <div className="flex items-center gap-1">
                            <CreditCard className="h-4 w-4 text-emerald-400" />
                            <span>Secure payment</span>
                        </div>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
};
