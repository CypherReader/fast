import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Scale, Droplets, Flame, Footprints, Lock } from "lucide-react";

const ProgressPage = () => {
  const [waterCount, setWaterCount] = useState(5);
  const [showPremiumModal, setShowPremiumModal] = useState(false);
  const waterGoal = 8;
  const steps = 5420;
  const stepGoal = 8000;

  const weightData = [
    { date: "Mon", weight: 180 },
    { date: "Tue", weight: 179.5 },
    { date: "Wed", weight: 179 },
    { date: "Thu", weight: 178.8 },
    { date: "Fri", weight: 178.5 },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="animate-fade-in">
        <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
          Progress Tracking
        </h1>
        <p className="text-sm text-muted-foreground">Monitor your health metrics</p>
      </div>

      {/* Weight Tracker */}
      <Card className="border-primary/20 animate-fade-in-up hover:border-primary/40 transition-all duration-300 hover:shadow-lg hover:shadow-primary/10" style={{ animationDelay: '0.1s' }}>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Scale className="h-5 w-5 text-primary" />
            Weight Trend
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {weightData.map((data, index) => (
              <div key={index} className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">{data.date}</span>
                <span className="text-sm font-semibold">{data.weight} lbs</span>
              </div>
            ))}
          </div>
          <div className="mt-4 pt-4 border-t border-border">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Current</span>
              <span className="font-bold text-primary">178.5 lbs</span>
            </div>
            <div className="flex justify-between text-sm mt-1">
              <span className="text-muted-foreground">This week</span>
              <span className="font-bold text-chart-4">-1.5 lbs</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Hydration Tracker */}
      <Card className="border-primary/20 animate-fade-in-up hover:border-primary/40 transition-all duration-300 hover:shadow-lg hover:shadow-primary/10" style={{ animationDelay: '0.2s' }}>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Droplets className="h-5 w-5 text-primary" />
            Daily Hydration
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2 mb-4">
            {Array.from({ length: waterGoal }).map((_, i) => (
              <button
                key={i}
                onClick={() => setWaterCount(Math.min(i + 1, waterGoal))}
                className={`flex-1 h-12 rounded-lg transition-all duration-300 hover:scale-110 ${
                  i < waterCount
                    ? "bg-primary/20 border-2 border-primary animate-scale-in"
                    : "bg-muted border-2 border-border hover:border-primary/50"
                }`}
                style={{ animationDelay: `${i * 0.05}s` }}
              >
                <Droplets
                  className={`mx-auto h-6 w-6 ${
                    i < waterCount ? "text-primary" : "text-muted-foreground"
                  }`}
                />
              </button>
            ))}
          </div>
          <div className="text-center text-sm text-muted-foreground">
            {waterCount} / {waterGoal} glasses today
          </div>
        </CardContent>
      </Card>

      {/* Keto Gauge - Premium Feature */}
      <Card className="border-secondary/20 relative overflow-hidden animate-fade-in-up hover:border-secondary/40 transition-all duration-300 hover:shadow-lg hover:shadow-secondary/10" style={{ animationDelay: '0.3s' }}>
        <div className="absolute top-2 right-2">
          <Lock className="h-4 w-4 text-secondary" />
        </div>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Flame className="h-5 w-5 text-secondary" />
            Ketosis Level
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="text-center">
              <div className="text-4xl font-bold text-secondary">73</div>
              <div className="text-sm text-muted-foreground">Estimated Score</div>
            </div>
            <Progress value={73} className="h-2" />
            <Button
              onClick={() => setShowPremiumModal(true)}
              variant="outline"
              className="w-full border-secondary/50 hover:bg-secondary/10"
            >
              <Lock className="mr-2 h-4 w-4" />
              Log Precision Blood Ketones
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Step Monitor */}
      <Card className="border-primary/20 animate-fade-in-up hover:border-primary/40 transition-all duration-300 hover:shadow-lg hover:shadow-primary/10" style={{ animationDelay: '0.4s' }}>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Footprints className="h-5 w-5 text-primary" />
            Daily Steps
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-2xl font-bold">{steps.toLocaleString()}</span>
              <span className="text-sm text-muted-foreground">/ {stepGoal.toLocaleString()}</span>
            </div>
            <Progress value={(steps / stepGoal) * 100} className="h-2" />
            <p className="text-sm text-muted-foreground">
              {stepGoal - steps} steps to reach your goal
            </p>
          </div>
        </CardContent>
      </Card>

      {/* Premium Upsell Modal */}
      <Dialog open={showPremiumModal} onOpenChange={setShowPremiumModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Lock className="h-5 w-5 text-secondary" />
              Unlock Premium
            </DialogTitle>
            <DialogDescription className="space-y-4 pt-4">
              <p>
                Get access to precision tracking features including blood ketone logging,
                advanced analytics, and personalized insights.
              </p>
              <div className="bg-gradient-to-r from-primary/10 to-secondary/10 p-4 rounded-lg">
                <h4 className="font-semibold mb-2">Premium Features:</h4>
                <ul className="space-y-1 text-sm">
                  <li>✓ Blood ketone tracking</li>
                  <li>✓ Advanced biomarker analysis</li>
                  <li>✓ Personalized recommendations</li>
                  <li>✓ Extended fasting protocols</li>
                </ul>
              </div>
              <Button className="w-full bg-gradient-to-r from-primary to-secondary">
                Upgrade to Premium
              </Button>
            </DialogDescription>
          </DialogHeader>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default ProgressPage;
