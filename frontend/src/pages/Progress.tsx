import { useState, useEffect, useRef } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Scale, Droplets, Flame, Camera, Lock } from "lucide-react";
import { api, mealApi } from "@/api/client";

const ProgressPage = () => {
  const [waterCount, setWaterCount] = useState(5);
  const [showPremiumModal, setShowPremiumModal] = useState(false);
  const [mealsLogged, setMealsLogged] = useState(0);
  const [earnedBack, setEarnedBack] = useState(0.00);
  const waterGoal = 8;
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [recentMeals, setRecentMeals] = useState<any[]>([]);

  const weightData = [
    { date: "Mon", weight: 180 },
    { date: "Tue", weight: 179.5 },
    { date: "Wed", weight: 179 },
    { date: "Thu", weight: 178.8 },
    { date: "Fri", weight: 178.5 },
  ];

  const fetchMeals = async () => {
    try {
      const res = await mealApi.list();
      const meals = res.data || [];
      setRecentMeals(meals);
      // Count meals logged today
      const today = new Date().toISOString().split('T')[0];
      const todayMeals = meals.filter((m: any) => m.logged_at.startsWith(today));
      setMealsLogged(todayMeals.length);
      setEarnedBack(todayMeals.length * 0.50);
    } catch (e) {
      console.error("Failed to fetch meals", e);
    }
  };

  useEffect(() => {
    fetchMeals();
  }, []);

  const handleCameraClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = async () => {
        const base64String = reader.result as string;
        // Log meal
        try {
          await mealApi.log(base64String, "Meal logged via camera");
          alert("Meal logged successfully! +$0.50 earned back.");
          fetchMeals();
        } catch (e) {
          alert("Failed to log meal");
        }
      };
      reader.readAsDataURL(file);
    }
  };

  const [showWeightModal, setShowWeightModal] = useState(false);
  const [manualWeight, setManualWeight] = useState("");
  const [loading, setLoading] = useState(false);
  const [weightUnit, setWeightUnit] = useState<'kg' | 'lbs'>('lbs'); // Unit preference

  const [currentWeight, setCurrentWeight] = useState<number | null>(null);

  const fetchWeight = async () => {
    try {
      console.log("Fetching weight...");
      const res = await api.get('/telemetry/metric', { params: { type: 'weight' } });
      console.log("Weight fetched:", res.data);
      setCurrentWeight(res.data.value);
    } catch (e) {
      console.log("No weight data found", e);
    }
  };

  useEffect(() => {
    fetchWeight();
  }, []);

  const handleLogWeight = async () => {
    if (!manualWeight) return;
    setLoading(true);
    try {
      console.log("Logging weight:", manualWeight, weightUnit);
      // Convert to lbs for storage if in kg
      const weightInLbs = weightUnit === 'kg' ? parseFloat(manualWeight) * 2.20462 : parseFloat(manualWeight);
      await api.post('/telemetry/manual', {
        type: 'weight',
        value: weightInLbs,
        unit: 'lbs' // Always store in lbs
      });
      console.log("Weight logged successfully");
      setManualWeight("");
      setShowWeightModal(false);
      alert("Weight logged successfully");
      await fetchWeight(); // Refresh data immediately
    } catch (error) {
      console.error("Failed to log weight:", error);
      alert("Failed to log weight");
    } finally {
      setLoading(false);
    }
  };

  // Convert weight for display
  const displayWeight = (weightInLbs: number | null) => {
    if (weightInLbs === null) return "No data";
    if (weightUnit === 'kg') {
      return `${(weightInLbs / 2.20462).toFixed(1)} kg`;
    }
    return `${weightInLbs} lbs`;
  };

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
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-base">
            <Scale className="h-5 w-5 text-primary" />
            Weight Trend
          </CardTitle>
          <Button
            size="sm"
            variant="outline"
            className="h-8"
            onClick={() => setShowWeightModal(true)}
          >
            Log Weight
          </Button>
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
              <span className="font-bold text-primary">{displayWeight(currentWeight)}</span>
            </div>
            <div className="flex justify-between text-sm mt-1">
              <span className="text-muted-foreground">This week</span>
              <span className="font-bold text-chart-4">-{weightUnit === 'kg' ? '0.7 kg' : '1.5 lbs'}</span>
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
                className={`flex-1 h-12 rounded-lg transition-all duration-300 hover:scale-110 ${i < waterCount
                  ? "bg-primary/20 border-2 border-primary animate-scale-in"
                  : "bg-muted border-2 border-border hover:border-primary/50"
                  }`}
                style={{ animationDelay: `${i * 0.05}s` }}
              >
                <Droplets
                  className={`mx-auto h-6 w-6 ${i < waterCount ? "text-primary" : "text-muted-foreground"
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

      {/* Food Logging - Commitment Vault */}
      <Card className="border-emerald-500/20 animate-fade-in-up hover:border-emerald-500/40 transition-all duration-300 hover:shadow-lg hover:shadow-emerald-500/10" style={{ animationDelay: '0.3s' }}>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Camera className="h-5 w-5 text-emerald-500" />
            Food Logging
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <div>
                <div className="text-2xl font-bold">{mealsLogged} / 3</div>
                <div className="text-sm text-muted-foreground">Meals Logged Today</div>
              </div>
              <div className="text-right">
                <div className="text-xl font-bold text-emerald-400">+${earnedBack.toFixed(2)}</div>
                <div className="text-xs text-emerald-500/80">Earned Back</div>
              </div>
            </div>
            <Progress value={(mealsLogged / 3) * 100} className="h-2 bg-slate-800" indicatorClassName="bg-emerald-500" />

            <input
              type="file"
              accept="image/*"
              capture="environment"
              ref={fileInputRef}
              onChange={handleFileChange}
              className="hidden"
            />

            <Button
              onClick={handleCameraClick}
              disabled={mealsLogged >= 3}
              className="w-full bg-emerald-600 hover:bg-emerald-700 text-white group disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <Camera className="mr-2 h-4 w-4 group-hover:scale-110 transition-transform" />
              {mealsLogged >= 3 ? "Daily Limit Reached" : "Log Meal Photo (+$0.50)"}
            </Button>

            <p className="text-xs text-center text-muted-foreground">
              Photos are analyzed for nutritional content.
            </p>

            {/* Recent Meals */}
            <div className="mt-4">
              <h4 className="text-sm font-semibold mb-2">Recent Meals</h4>
              <div className="space-y-3">
                {recentMeals.slice(0, 3).map((meal, i) => (
                  <div key={i} className="flex gap-3 bg-muted/50 p-2 rounded-lg">
                    <div className="relative w-16 h-16 rounded-md overflow-hidden flex-shrink-0">
                      <img src={meal.image} alt="Meal" className="object-cover w-full h-full" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium truncate">{meal.description}</p>
                      <p className="text-xs text-muted-foreground line-clamp-2">{meal.analysis}</p>
                      <div className="flex gap-2 mt-1">
                        {meal.is_keto && <span className="text-[10px] bg-green-500/20 text-green-500 px-1.5 py-0.5 rounded">Keto</span>}
                        {!meal.is_authentic && <span className="text-[10px] bg-red-500/20 text-red-500 px-1.5 py-0.5 rounded">Fake?</span>}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Keto Gauge - Premium Feature */}
      <Card className="border-secondary/20 relative overflow-hidden animate-fade-in-up hover:border-secondary/40 transition-all duration-300 hover:shadow-lg hover:shadow-secondary/10" style={{ animationDelay: '0.4s' }}>
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

      {/* Manual Weight Modal */}
      <Dialog open={showWeightModal} onOpenChange={setShowWeightModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Log Weight</DialogTitle>
            <DialogDescription>
              Enter your current weight.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 pt-4">
            {/* Unit Selector */}
            <div className="flex gap-2 mb-4">
              <button
                onClick={() => setWeightUnit('lbs')}
                className={`flex-1 px-4 py-2 rounded-md transition-all ${weightUnit === 'lbs'
                  ? 'bg-primary text-white'
                  : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
                  }`}
              >
                lbs
              </button>
              <button
                onClick={() => setWeightUnit('kg')}
                className={`flex-1 px-4 py-2 rounded-md transition-all ${weightUnit === 'kg'
                  ? 'bg-primary text-white'
                  : 'bg-slate-800 text-slate-400 hover:bg-slate-700'
                  }`}
              >
                kg
              </button>
            </div>
            <div className="flex gap-2">
              <input
                type="number"
                placeholder={weightUnit === 'kg' ? 'e.g. 79.5' : 'e.g. 175.5'}
                value={manualWeight}
                onChange={(e) => setManualWeight(e.target.value)}
                className="flex-1 bg-slate-900 border border-slate-700 rounded-md px-3 py-2 text-white focus:outline-none focus:border-primary"
                step="0.1"
              />
              <Button
                onClick={handleLogWeight}
                disabled={loading || !manualWeight}
              >
                Save
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default ProgressPage;
