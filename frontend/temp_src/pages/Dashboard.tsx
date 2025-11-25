import { useState } from "react";
import FastingTimer from "@/components/FastingTimer";
import BioPhaseCard from "@/components/BioPhaseCard";
import AffiliateCard from "@/components/AffiliateCard";
import { Button } from "@/components/ui/button";
import { Droplets, Share2 } from "lucide-react";

const Dashboard = () => {
  const [isFasting, setIsFasting] = useState(false);
  const [fastingHours] = useState(16);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="text-center mb-8 animate-fade-in">
        <h1 className="text-3xl font-bold text-shimmer animate-float">
          Autophagy Arc
        </h1>
        <p className="text-sm text-muted-foreground mt-2 animate-fade-in" style={{ animationDelay: '0.2s' }}>
          Transform through fasting
        </p>
      </div>

      {/* Main Timer */}
      <FastingTimer onStatusChange={setIsFasting} />

      {/* Bio Phase Indicator */}
      {isFasting && (
        <div className="animate-fade-in-up" style={{ animationDelay: '0.3s' }}>
          <BioPhaseCard hours={fastingHours} />
        </div>
      )}

      {/* Affiliate Card - Only show during fast */}
      {isFasting && (
        <div className="animate-fade-in-up" style={{ animationDelay: '0.5s' }}>
          <AffiliateCard />
        </div>
      )}

      {/* Quick Actions */}
      <div className="grid grid-cols-2 gap-4 animate-fade-in-up" style={{ animationDelay: '0.7s' }}>
        <Button
          variant="outline"
          className="h-20 flex-col gap-2 border-primary/30 hover:bg-primary/10 transition-all duration-300 hover:scale-105 hover:glow-primary group"
        >
          <Droplets className="h-6 w-6 text-primary group-hover:animate-bounce" />
          <span className="text-sm">Log Water</span>
        </Button>
        <Button
          variant="outline"
          className="h-20 flex-col gap-2 border-secondary/30 hover:bg-secondary/10 transition-all duration-300 hover:scale-105 hover:glow-secondary group"
        >
          <Share2 className="h-6 w-6 text-secondary group-hover:animate-pulse" />
          <span className="text-sm">Share Fast</span>
        </Button>
      </div>
    </div>
  );
};

export default Dashboard;
