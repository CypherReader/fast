import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Lock, Check, AlertTriangle } from "lucide-react";
import VaultStatus from "@/components/VaultStatus";

const Profile = () => {
  const [showMedicalModal, setShowMedicalModal] = useState(false);
  const [showPremiumModal, setShowPremiumModal] = useState(false);


  const fastingPlans = [
    {
      id: "16-8",
      name: "16:8 Beginner",
      description: "16 hours fasting, 8 hours eating",
      locked: false,
      recommended: true,
    },
    {
      id: "18-6",
      name: "18:6 Intermediate",
      description: "18 hours fasting, 6 hours eating",
      locked: false,
    },
    {
      id: "omad",
      name: "OMAD",
      description: "One meal a day",
      locked: false,
    },
    {
      id: "5-day",
      name: "5-Day Water Fast",
      description: "Extended deep cleanse",
      locked: true,
      requiresMedical: true,
    },
    {
      id: "7-day",
      name: "7-Day Deep Cleanse",
      description: "Maximum regeneration protocol",
      locked: true,
      requiresMedical: true,
    },
  ];

  const handlePlanSelect = (plan: typeof fastingPlans[0]) => {
    if (plan.requiresMedical) {

      setShowMedicalModal(true);
    } else if (plan.locked) {
      setShowPremiumModal(true);
    } else {
      // Handle plan selection
      console.log("Selected plan:", plan.id);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="animate-fade-in">
        <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
          Profile & Plans
        </h1>
        <p className="text-sm text-muted-foreground">Customize your fasting journey</p>
      </div>

      {/* Fasting Plans */}
      <div>
        <h2 className="text-lg font-semibold mb-4">Choose Your Protocol</h2>
        <div className="space-y-3">
          {fastingPlans.map((plan, index) => (
            <Card
              key={plan.id}
              className={`border-primary/20 cursor-pointer transition-all duration-300 hover:border-primary/40 hover:scale-[1.02] hover:shadow-lg animate-fade-in-up ${plan.recommended ? "border-primary/40 bg-primary/5 glow-primary" : ""
                }`}
              style={{ animationDelay: `${index * 0.1}s` }}
              onClick={() => handlePlanSelect(plan)}
            >
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <CardTitle className="text-base flex items-center gap-2">
                      {plan.name}
                      {plan.recommended && (
                        <Badge variant="secondary" className="text-xs">
                          Recommended
                        </Badge>
                      )}
                      {plan.locked && (
                        <Lock className="h-4 w-4 text-secondary" />
                      )}
                    </CardTitle>
                    <CardDescription className="text-sm mt-1">
                      {plan.description}
                    </CardDescription>
                  </div>
                  {!plan.locked && !plan.requiresMedical && (
                    <Check className="h-5 w-5 text-primary" />
                  )}
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>

      {/* Commitment Vault */}
      <VaultStatus
        deposit={20.00}
        earned={5.50}
        potentialRefund={5.50}
      />

      {/* Elite Medical Data Section */}
      <Card className="border-secondary/20 relative overflow-hidden">
        <div className="absolute top-2 right-2">
          <Lock className="h-4 w-4 text-secondary" />
        </div>
        <CardHeader>
          <CardTitle className="text-base">Elite Medical Data</CardTitle>
          <CardDescription>Track advanced biomarkers</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            <div className="flex justify-between items-center py-2 border-b border-border">
              <span className="text-sm text-muted-foreground">Fasting Insulin</span>
              <Badge variant="outline" className="text-xs">Premium</Badge>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-border">
              <span className="text-sm text-muted-foreground">CRP Level</span>
              <Badge variant="outline" className="text-xs">Premium</Badge>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-border">
              <span className="text-sm text-muted-foreground">HbA1c</span>
              <Badge variant="outline" className="text-xs">Premium</Badge>
            </div>
            <Button
              onClick={() => setShowPremiumModal(true)}
              variant="outline"
              className="w-full mt-4 border-secondary/50 hover:bg-secondary/10"
            >
              <Lock className="mr-2 h-4 w-4" />
              Unlock Medical Tracking
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Medical Disclaimer Modal */}
      <Dialog open={showMedicalModal} onOpenChange={setShowMedicalModal}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-destructive">
              <AlertTriangle className="h-5 w-5" />
              Medical Consultation Required
            </DialogTitle>
            <DialogDescription className="space-y-4 pt-4">
              <div className="bg-destructive/10 p-4 rounded-lg border border-destructive/20">
                <p className="text-sm text-foreground font-semibold mb-2">
                  ⚠️ Important Warning
                </p>
                <p className="text-sm">
                  Extended fasting requires medical supervision. Please confirm you have
                  consulted with a healthcare professional before proceeding.
                </p>
              </div>
              <ul className="space-y-2 text-sm">
                <li>✓ Risk assessment completed</li>
                <li>✓ Medical history reviewed</li>
                <li>✓ Doctor approval obtained</li>
              </ul>
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="flex-col gap-2 sm:flex-col">
            <Button
              onClick={() => setShowMedicalModal(false)}
              variant="outline"
              className="w-full"
            >
              Cancel
            </Button>
            <Button
              onClick={() => {
                setShowMedicalModal(false);
                // Handle confirmation
              }}
              className="w-full bg-gradient-to-r from-primary to-secondary"
            >
              I Confirm - Start Protocol
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Premium Upsell Modal */}
      <Dialog open={showPremiumModal} onOpenChange={setShowPremiumModal}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Lock className="h-5 w-5 text-secondary" />
              Unlock Premium Features
            </DialogTitle>
            <DialogDescription className="space-y-4 pt-4">
              <div className="bg-gradient-to-r from-primary/10 to-secondary/10 p-4 rounded-lg">
                <h4 className="font-semibold mb-3">Premium Includes:</h4>
                <ul className="space-y-2 text-sm">
                  <li>✓ Extended fasting protocols</li>
                  <li>✓ Advanced biomarker tracking</li>
                  <li>✓ Blood ketone logging</li>
                  <li>✓ Personalized AI insights</li>
                  <li>✓ Priority support</li>
                  <li>✓ Medical data export</li>
                </ul>
              </div>
              <div className="text-center py-2">
                <div className="text-3xl font-bold text-primary">$9.99</div>
                <div className="text-sm text-muted-foreground">per month</div>
              </div>
              <Button className="w-full bg-gradient-to-r from-primary to-secondary text-lg h-12">
                Upgrade Now
              </Button>
            </DialogDescription>
          </DialogHeader>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default Profile;
