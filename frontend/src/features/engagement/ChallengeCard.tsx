import { useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { CheckCircle, Circle, Trophy } from "lucide-react";
import { cn } from "@/lib/utils";

interface Challenge {
    id: string;
    title: string;
    reward: string;
    completed: boolean;
}

export const ChallengeCard = () => {
    const [challenges, setChallenges] = useState<Challenge[]>([
        { id: '1', title: 'Log your first meal', reward: '+5 pts', completed: false },
        { id: '2', title: 'Drink 2L of water', reward: '+5 pts', completed: false },
        { id: '3', title: 'Walk 10 minutes', reward: '+10 pts', completed: false },
    ]);

    const toggleChallenge = (id: string) => {
        setChallenges(prev => prev.map(c => {
            if (c.id === id) {
                // Simple confetti effect could be triggered here
                return { ...c, completed: !c.completed };
            }
            return c;
        }));
    };

    const allCompleted = challenges.every(c => c.completed);

    return (
        <Card className={cn(
            "border-slate-800 transition-all duration-500",
            allCompleted ? "bg-emerald-900/20 border-emerald-500/30" : "bg-slate-900"
        )}>
            <CardContent className="p-4">
                <div className="flex justify-between items-center mb-4">
                    <h3 className="font-bold text-white flex items-center">
                        <Trophy className={cn("w-5 h-5 mr-2", allCompleted ? "text-emerald-400" : "text-yellow-500")} />
                        Daily Challenges
                    </h3>
                    <span className="text-xs font-mono text-slate-400">
                        {challenges.filter(c => c.completed).length}/{challenges.length}
                    </span>
                </div>

                <div className="space-y-3">
                    {challenges.map((challenge) => (
                        <div
                            key={challenge.id}
                            onClick={() => toggleChallenge(challenge.id)}
                            className={cn(
                                "flex items-center justify-between p-3 rounded-lg cursor-pointer transition-all border",
                                challenge.completed
                                    ? "bg-emerald-500/10 border-emerald-500/20"
                                    : "bg-slate-800/50 border-transparent hover:bg-slate-800"
                            )}
                        >
                            <div className="flex items-center space-x-3">
                                {challenge.completed ? (
                                    <CheckCircle className="w-5 h-5 text-emerald-400" />
                                ) : (
                                    <Circle className="w-5 h-5 text-slate-500" />
                                )}
                                <span className={cn(
                                    "text-sm font-medium transition-colors",
                                    challenge.completed ? "text-emerald-200 line-through opacity-70" : "text-slate-200"
                                )}>
                                    {challenge.title}
                                </span>
                            </div>
                            <span className={cn(
                                "text-xs font-bold px-2 py-1 rounded",
                                challenge.completed ? "bg-emerald-500/20 text-emerald-300" : "bg-yellow-500/10 text-yellow-500"
                            )}>
                                {challenge.reward}
                            </span>
                        </div>
                    ))}
                </div>

                {allCompleted && (
                    <div className="mt-4 text-center animate-in zoom-in duration-300">
                        <p className="text-emerald-400 text-sm font-bold">
                            🎉 All challenges completed! Bonus +20 pts earned.
                        </p>
                    </div>
                )}
            </CardContent>
        </Card>
    );
};
