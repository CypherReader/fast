import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Zap, Droplets, Brain, Sparkles } from "lucide-react";

interface BioPhaseCardProps {
  hours: number;
}

const getBioPhase = (hours: number) => {
  if (hours < 4) {
    return {
      phase: "Early Fasting",
      description: "Glucose depletion beginning",
      icon: Droplets,
      color: "text-chart-5",
    };
  } else if (hours < 12) {
    return {
      phase: "Fat Burning",
      description: "Switching to ketone production",
      icon: Zap,
      color: "text-chart-1",
    };
  } else if (hours < 18) {
    return {
      phase: "Deep Autophagy",
      description: "Cellular cleanup activated",
      icon: Sparkles,
      color: "text-secondary",
    };
  } else {
    return {
      phase: "Peak Regeneration",
      description: "Maximum cellular renewal",
      icon: Brain,
      color: "text-primary",
    };
  }
};

const BioPhaseCard = ({ hours }: BioPhaseCardProps) => {
  const phase = getBioPhase(hours);
  const Icon = phase.icon;

  return (
    <Card className="relative overflow-hidden bg-gradient-to-br from-card to-muted border-primary/20 group hover:border-primary/40 transition-all duration-300">
      {/* Background glow effect */}
      <div className="absolute inset-0 bg-gradient-to-r from-primary/5 to-secondary/5 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
      
      <CardHeader className="pb-3 relative">
        <CardTitle className="text-sm font-medium flex items-center gap-2">
          <Icon className={`h-5 w-5 ${phase.color} group-hover:animate-pulse`} />
          <span className="group-hover:text-primary transition-colors">Hour {hours}: {phase.phase}</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="relative">
        <p className="text-sm text-muted-foreground">{phase.description}</p>
      </CardContent>
    </Card>
  );
};

export default BioPhaseCard;
