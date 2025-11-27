import { useEffect } from 'react';
import { motion } from 'framer-motion';
import { Check, Timer, Users, Compass, Lightbulb, ChevronDown } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import confetti from 'canvas-confetti';
import OnboardingLayout from '@/components/onboarding/OnboardingLayout';
import { Button } from '@/components/ui/button';
import { useOnboarding } from '@/contexts/OnboardingContext';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';
import { useState } from 'react';

const planNames: Record<string, string> = {
  '16:8': '16:8',
  '18:6': '18:6',
  'omad': '23:1 OMAD',
};

const goalNames: Record<string, string> = {
  'weight_loss': 'Weight Loss',
  'metabolic': 'Metabolic Health',
  'discipline': 'Discipline Building',
  'longevity': 'Longevity',
};

const tips = [
  'Start your fast after dinner tonight',
  'Stay hydrated during your fasting window',
  'Track your weight and meals for bonus recovery',
  'Join a tribe for accountability',
  'Enable notifications so you don\'t miss earning',
];

const OnboardingSuccess = () => {
  const navigate = useNavigate();
  const { state, completeOnboarding } = useOnboarding();
  const [tipsOpen, setTipsOpen] = useState(false);

  // Confetti effect on mount
  useEffect(() => {
    completeOnboarding();
    
    const duration = 2000;
    const end = Date.now() + duration;

    const colors = ['#f59e0b', '#10b981', '#a855f7'];

    const frame = () => {
      confetti({
        particleCount: 3,
        angle: 60,
        spread: 55,
        origin: { x: 0, y: 0.7 },
        colors,
      });
      confetti({
        particleCount: 3,
        angle: 120,
        spread: 55,
        origin: { x: 1, y: 0.7 },
        colors,
      });

      if (Date.now() < end) {
        requestAnimationFrame(frame);
      }
    };

    frame();

    // Center burst
    setTimeout(() => {
      confetti({
        particleCount: 100,
        spread: 70,
        origin: { y: 0.6 },
        colors,
      });
    }, 300);
  }, [completeOnboarding]);

  const handleStartFasting = () => {
    // In a real app, this would start a fast and navigate to dashboard
    navigate('/');
  };

  return (
    <OnboardingLayout step={6}>
      <div className="max-w-xl w-full text-center">
        {/* Success Animation */}
        <motion.div
          className="mb-6 inline-flex items-center justify-center w-28 h-28 rounded-full bg-secondary/20"
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: 'spring', stiffness: 200, damping: 15, delay: 0.2 }}
        >
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.4 }}
          >
            <Check className="w-14 h-14 text-secondary" strokeWidth={3} />
          </motion.div>
        </motion.div>

        {/* Headlines */}
        <motion.h1 
          className="text-4xl md:text-5xl font-extrabold text-foreground mb-3"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          Welcome to the Vault!
        </motion.h1>
        <motion.p 
          className="text-lg text-muted-foreground mb-8"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
        >
          Your $20 is secured. Let&apos;s start earning it back.
        </motion.p>

        {/* Summary Card */}
        <motion.div
          className="bg-gradient-to-br from-primary/10 via-card to-secondary/10 border border-secondary/30 rounded-2xl p-6 mb-8 text-left"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
        >
          <h3 className="font-semibold text-foreground mb-4">Your Vault Setup</h3>
          <div className="space-y-2 text-sm">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Goal</span>
              <span className="text-foreground font-medium">
                {state.goal ? goalNames[state.goal] : 'Weight Loss'}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Schedule</span>
              <span className="text-foreground font-medium">
                {state.fastingPlan ? planNames[state.fastingPlan] : '16:8'}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Monthly deposit</span>
              <span className="text-foreground font-medium">$20.00</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Next refund</span>
              <span className="text-foreground font-medium">1st of next month</span>
            </div>
          </div>
          <div className="border-t border-border mt-4 pt-4">
            <p className="text-sm text-secondary font-medium">
              Complete 10 fasts this month → Recover full $20
            </p>
          </div>
        </motion.div>

        {/* Primary CTA - Start First Fast */}
        <motion.div
          className="bg-card border border-border rounded-2xl p-6 mb-6"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6 }}
        >
          <Timer className="w-12 h-12 text-secondary mx-auto mb-3" />
          <h3 className="text-xl font-semibold text-foreground mb-2">Start Your First Fast</h3>
          <p className="text-sm text-muted-foreground mb-4">
            Begin your journey right now and start recovering your deposit
          </p>
          <Button
            size="xl"
            className="w-full bg-gradient-gold hover:scale-105 transition-transform shadow-gold-glow animate-pulse-gold"
            onClick={handleStartFasting}
          >
            Start Fasting
          </Button>
          <p className="text-sm text-primary mt-3 font-medium">
            You&apos;ll recover $2.00 when you complete this fast
          </p>
        </motion.div>

        {/* Secondary Actions */}
        <motion.div
          className="grid grid-cols-2 gap-4 mb-6"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.7 }}
        >
          <button 
            className="bg-card border border-border rounded-xl p-4 hover:border-secondary/50 transition-colors text-left"
            onClick={() => navigate('/')}
          >
            <Users className="w-8 h-8 text-secondary mb-2" />
            <h4 className="font-medium text-foreground text-sm mb-1">Join a Tribe</h4>
            <p className="text-xs text-muted-foreground">2x more likely to succeed</p>
          </button>
          <button 
            className="bg-card border border-border rounded-xl p-4 hover:border-secondary/50 transition-colors text-left"
            onClick={() => navigate('/')}
          >
            <Compass className="w-8 h-8 text-secondary mb-2" />
            <h4 className="font-medium text-foreground text-sm mb-1">Take a Tour</h4>
            <p className="text-xs text-muted-foreground">Learn all the features</p>
          </button>
        </motion.div>

        {/* Tips Section */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.8 }}
        >
          <Collapsible open={tipsOpen} onOpenChange={setTipsOpen}>
            <CollapsibleTrigger className="flex items-center justify-center gap-2 w-full text-muted-foreground hover:text-foreground transition-colors py-2">
              <Lightbulb className="w-4 h-4" />
              <span className="text-sm font-medium">Tips for Success</span>
              <ChevronDown className={`w-4 h-4 transition-transform ${tipsOpen ? 'rotate-180' : ''}`} />
            </CollapsibleTrigger>
            <CollapsibleContent>
              <div className="bg-card border border-border rounded-xl p-4 mt-2 text-left">
                <ul className="space-y-2">
                  {tips.map((tip, index) => (
                    <li key={index} className="flex items-start gap-2 text-sm text-muted-foreground">
                      <span className="text-secondary">💡</span>
                      <span>{tip}</span>
                    </li>
                  ))}
                </ul>
              </div>
            </CollapsibleContent>
          </Collapsible>
        </motion.div>

        {/* Skip Link */}
        <motion.p 
          className="mt-6"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.9 }}
        >
          <button 
            className="text-sm text-muted-foreground hover:text-foreground hover:underline transition-colors"
            onClick={() => navigate('/')}
          >
            I&apos;ll start fasting later
          </button>
        </motion.p>
      </div>
    </OnboardingLayout>
  );
};

export default OnboardingSuccess;
