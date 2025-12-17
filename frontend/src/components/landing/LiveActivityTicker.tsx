import { motion } from "framer-motion";
import { TrendingUp, Users, Zap } from "lucide-react";
import { useEffect, useState } from "react";

const activities = [
    { icon: Users, name: "Sarah from NYC", action: "just started her transformation", color: "text-primary" },
    { icon: Zap, name: "Mike R.", action: "hit his 30-day streak", color: "text-secondary" },
    { icon: TrendingUp, name: "Emma L.", action: "lost 15 lbs this month", color: "text-accent" },
    { icon: Users, name: "David from LA", action: "joined the Fast tribe", color: "text-primary" },
    { icon: Zap, name: "Lisa M.", action: "completed her first 16:8 fast", color: "text-secondary" },
    { icon: TrendingUp, name: "James W.", action: "reached his goal weight", color: "text-accent" },
    { icon: Users, name: "Anna from Chicago", action: "started using Cortex AI", color: "text-primary" },
    { icon: Zap, name: "Tom H.", action: "unlocked 60-day badge", color: "text-secondary" },
];

export const LiveActivityTicker = () => {
    const [currentIndex, setCurrentIndex] = useState(0);

    useEffect(() => {
        const interval = setInterval(() => {
            setCurrentIndex((prev) => (prev + 1) % activities.length);
        }, 3000); // Change every 3 seconds

        return () => clearInterval(interval);
    }, []);

    const current = activities[currentIndex];

    return (
        <section className="py-8 bg-card/50 border-y border-border/50 overflow-hidden">
            <div className="container px-4">
                <div className="flex items-center justify-center gap-3">
                    {/* Live indicator */}
                    <div className="flex items-center gap-2">
                        <motion.div
                            animate={{ scale: [1, 1.2, 1], opacity: [1, 0.5, 1] }}
                            transition={{ duration: 2, repeat: Infinity }}
                            className="w-2 h-2 bg-red-500 rounded-full"
                        />
                        <span className="text-xs text-muted-foreground font-medium">LIVE</span>
                    </div>

                    {/* Activity message */}
                    <motion.div
                        key={currentIndex}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -10 }}
                        transition={{ duration: 0.5 }}
                        className="flex items-center gap-2 text-sm"
                    >
                        <current.icon className={`w-4 h-4 ${current.color}`} />
                        <span className="font-semibold text-foreground">{current.name}</span>
                        <span className="text-muted-foreground">{current.action}</span>
                    </motion.div>

                    {/* Counter */}
                    <div className="hidden md:flex items-center gap-1 ml-4 px-3 py-1 bg-primary/10 rounded-full">
                        <Users className="w-3 h-3 text-primary" />
                        <span className="text-xs font-semibold text-primary">127 joined today</span>
                    </div>
                </div>
            </div>
        </section>
    );
};
