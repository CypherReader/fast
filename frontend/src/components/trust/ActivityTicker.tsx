import { Card, CardContent } from "@/components/ui/card";
import { useState, useEffect } from "react";

interface Activity {
    userName: string;
    action: string;
    amount: string;
    timestamp: string;
}

export const ActivityTicker = () => {
    const [recentActivity, setRecentActivity] = useState<Activity>({
        userName: "Sarah K.",
        action: "just earned",
        amount: "2.00",
        timestamp: "2m ago"
    });

    // Mock activity updates
    useEffect(() => {
        const activities = [
            { userName: "Mike R.", action: "joined the vault", amount: "20.00", timestamp: "just now" },
            { userName: "Alex T.", action: "hit a 7-day streak", amount: "5.00", timestamp: "1m ago" },
            { userName: "Sarah K.", action: "just earned", amount: "2.00", timestamp: "2m ago" },
            { userName: "David L.", action: "got refunded", amount: "20.00", timestamp: "5m ago" }
        ];

        let index = 0;
        const interval = setInterval(() => {
            index = (index + 1) % activities.length;
            setRecentActivity(activities[index]);
        }, 5000);

        return () => clearInterval(interval);
    }, []);

    return (
        <Card className="bg-slate-900/50 border-slate-800 mb-4 overflow-hidden animate-fade-in">
            <CardContent className="p-0">
                <div className="flex items-center gap-3 px-4 py-3">
                    {/* Pulsing Indicator */}
                    <div className="relative">
                        <div className="h-2 w-2 bg-emerald-400 rounded-full" />
                        <div className="absolute inset-0 h-2 w-2 bg-emerald-400 rounded-full animate-ping" />
                    </div>

                    {/* Activity Text */}
                    <p className="text-sm animate-fade-in key={recentActivity.timestamp}">
                        <span className="font-semibold text-white">{recentActivity.userName}</span>
                        {' '}{recentActivity.action}
                        {' '}<span className="text-emerald-400">${recentActivity.amount}</span>
                    </p>

                    {/* Timestamp */}
                    <span className="text-xs text-slate-500 ml-auto">
                        {recentActivity.timestamp}
                    </span>
                </div>
            </CardContent>
        </Card>
    );
};
