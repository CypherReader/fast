import { useEffect, useState } from "react";

const activities = [
  { name: "Sarah M.", action: "recovered $2.00", time: "2m ago" },
  { name: "Mike R.", action: "completed 18h fast", time: "5m ago" },
  { name: "Anna K.", action: "hit 30-day streak!", time: "7m ago" },
  { name: "James D.", action: "got $2.00 refund", time: "9m ago" },
  { name: "Lisa B.", action: "completed 16h fast", time: "11m ago" },
  { name: "Tom W.", action: "reached $0 net cost", time: "14m ago" },
  { name: "Emma S.", action: "recovered full $20", time: "16m ago" },
  { name: "Chris L.", action: "completed 24h fast", time: "19m ago" },
];

export const LiveTickerSection = () => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setIsVisible(true), 500);
    return () => clearTimeout(timer);
  }, []);

  return (
    <div className={`bg-background border-y border-border py-4 overflow-hidden transition-opacity duration-500 ${isVisible ? "opacity-100" : "opacity-0"}`}>
      <div className="flex items-center gap-4">
        {/* Live indicator */}
        <div className="flex-shrink-0 flex items-center gap-2 px-4">
          <div className="w-2.5 h-2.5 rounded-full bg-secondary animate-pulse" />
          <span className="text-sm font-medium text-muted-foreground">Live</span>
        </div>

        {/* Ticker */}
        <div className="flex-1 overflow-hidden">
          <div className="flex animate-ticker">
            {[...activities, ...activities].map((activity, i) => (
              <div
                key={i}
                className="flex-shrink-0 flex items-center gap-4 px-8 whitespace-nowrap"
              >
                <span className="font-semibold text-foreground">{activity.name}</span>
                <span className="text-primary">{activity.action}</span>
                <span className="text-muted-foreground text-sm">• {activity.time}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
