import { motion } from 'framer-motion';
import { Check, Star } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import OnboardingLayout from '@/components/onboarding/OnboardingLayout';
import { Button } from '@/components/ui/button';
import { useOnboarding, FastingPlan } from '@/contexts/OnboardingContext';
import { cn } from '@/lib/utils';

const plans = [
  {
    id: '16:8' as FastingPlan,
    title: '16:8 - Time-Restricted Eating',
    fastHours: 16,
    eatHours: 8,
    description: 'Fast for 16 hours, eat within 8-hour window',
    example: 'Example: Eat 12pm-8pm, fast 8pm-12pm',
    difficulty: 1,
    badge: 'Recommended for beginners',
    badgeColor: 'bg-secondary text-secondary-foreground',
  },
  {
    id: '18:6' as FastingPlan,
    title: '18:6 - Extended Fasting',
    fastHours: 18,
    eatHours: 6,
    description: 'Fast for 18 hours, eat within 6-hour window',
    example: 'Example: Eat 2pm-8pm, fast 8pm-2pm',
    difficulty: 2,
    badge: 'Most Popular',
    badgeColor: 'bg-primary text-primary-foreground',
  },
  {
    id: 'omad' as FastingPlan,
    title: '23:1 - One Meal A Day',
    fastHours: 23,
    eatHours: 1,
    description: 'Fast 23 hours, eat one meal per day',
    example: 'Example: Eat 6pm-7pm, fast rest of day',
    difficulty: 3,
    badge: null,
    warning: 'Not recommended for beginners',
  },
];

const TimelineBar = ({ fastHours, eatHours }: { fastHours: number; eatHours: number }) => {
  const eatPercentage = (eatHours / 24) * 100;
  const fastPercentage = (fastHours / 24) * 100;

  return (
    <div className="mt-4">
      <div className="flex h-3 rounded-full overflow-hidden bg-muted">
        <div
          className="bg-muted-foreground/30"
          style={{ width: `${fastPercentage}%` }}
        />
        <div
          className="bg-secondary"
          style={{ width: `${eatPercentage}%` }}
        />
      </div>
      <div className="flex justify-between text-xs text-muted-foreground mt-1">
        <span>12am</span>
        <span>6am</span>
        <span>12pm</span>
        <span>6pm</span>
        <span>12am</span>
      </div>
    </div>
  );
};

const OnboardingPlan = () => {
  const navigate = useNavigate();
  const { state, setFastingPlan } = useOnboarding();

  const selectedPlanName = state.fastingPlan
    ? plans.find(p => p.id === state.fastingPlan)?.title.split(' - ')[0]
    : '';

  return (
    <OnboardingLayout step={3} showBack backPath="/onboarding/goal">
      <div className="max-w-xl w-full">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-2">
            Choose your fasting schedule
          </h1>
          <p className="text-muted-foreground">
            You can change this later based on what works for you
          </p>
        </div>

        {/* Plan Options */}
        <div className="space-y-4 mb-8">
          {plans.map((plan, index) => {
            const isSelected = state.fastingPlan === plan.id;
            return (
              <motion.button
                key={plan.id}
                onClick={() => setFastingPlan(plan.id)}
                className={cn(
                  'relative w-full p-5 rounded-xl border-2 text-left transition-all duration-200',
                  isSelected
                    ? 'border-secondary bg-secondary/10 shadow-emerald-glow'
                    : 'border-border bg-card hover:border-secondary/50'
                )}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.05 * index }}
              >
                {plan.badge && (
                  <span className={cn(
                    'absolute top-3 right-3 px-2 py-0.5 text-xs font-semibold rounded-full',
                    plan.badgeColor
                  )}>
                    {plan.badge}
                  </span>
                )}
                {isSelected && (
                  <span className="absolute top-3 right-3 w-6 h-6 rounded-full bg-secondary flex items-center justify-center">
                    <Check className="w-4 h-4 text-secondary-foreground" />
                  </span>
                )}

                <div className="flex items-baseline gap-3 mb-2">
                  <span className="text-3xl font-bold text-primary">{plan.fastHours}</span>
                  <span className="text-muted-foreground">hrs fasting</span>
                  <span className="text-xl font-semibold text-secondary">{plan.eatHours}</span>
                  <span className="text-muted-foreground">hrs eating</span>
                </div>

                <h3 className="text-lg font-semibold text-foreground mb-1">{plan.title}</h3>
                <p className="text-sm text-muted-foreground mb-1">{plan.description}</p>
                <p className="text-xs text-muted-foreground/70">{plan.example}</p>

                {/* Difficulty */}
                <div className="flex items-center gap-2 mt-3">
                  <span className="text-xs text-muted-foreground">Difficulty:</span>
                  <div className="flex gap-0.5">
                    {[1, 2, 3].map(level => (
                      <Star
                        key={level}
                        className={cn(
                          'w-3.5 h-3.5',
                          level <= plan.difficulty ? 'text-primary fill-primary' : 'text-muted'
                        )}
                      />
                    ))}
                  </div>
                  {plan.warning && (
                    <span className="text-xs text-destructive ml-2">{plan.warning}</span>
                  )}
                </div>

                {/* Recovery */}
                <div className="flex items-center gap-2 mt-2">
                  <span className="text-xs text-muted-foreground">Recovery per fast:</span>
                  <span className="text-sm font-semibold text-primary">$2.00</span>
                </div>

                {/* Timeline */}
                <TimelineBar fastHours={plan.fastHours} eatHours={plan.eatHours} />
              </motion.button>
            );
          })}
        </div>

        {/* Navigation */}
        <div className="flex justify-center">
          <Button
            size="xl"
            className="w-full md:w-80 bg-gradient-gold hover:scale-105 transition-transform shadow-gold-glow disabled:opacity-50 disabled:hover:scale-100"
            disabled={!state.fastingPlan}
            onClick={() => {
              // Skip account creation for OAuth users (they already have an account)
              const hasToken = localStorage.getItem('token');
              navigate(hasToken ? '/onboarding/payment' : '/onboarding/account');
            }}
          >
            {selectedPlanName ? `Continue with ${selectedPlanName}` : 'Continue'}
          </Button>
        </div>
      </div>
    </OnboardingLayout>
  );
};

export default OnboardingPlan;
