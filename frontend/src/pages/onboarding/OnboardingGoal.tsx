import { motion } from 'framer-motion';
import { Scale, Activity, Target, Heart, Check } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import OnboardingLayout from '@/components/onboarding/OnboardingLayout';
import { Button } from '@/components/ui/button';
import { useOnboarding, GoalType } from '@/contexts/OnboardingContext';
import { cn } from '@/lib/utils';

const goals = [
  {
    id: 'weight_loss' as GoalType,
    icon: Scale,
    title: 'Weight Loss',
    description: 'Shed pounds and feel lighter',
    popular: true,
  },
  {
    id: 'metabolic' as GoalType,
    icon: Activity,
    title: 'Metabolic Health',
    description: 'Improve insulin sensitivity & energy',
    popular: false,
  },
  {
    id: 'discipline' as GoalType,
    icon: Target,
    title: 'Discipline Building',
    description: 'Build mental toughness and willpower',
    popular: false,
  },
  {
    id: 'longevity' as GoalType,
    icon: Heart,
    title: 'Longevity',
    description: 'Activate autophagy for healthy aging',
    popular: false,
  },
];

const OnboardingGoal = () => {
  const navigate = useNavigate();
  const { state, setGoal } = useOnboarding();

  return (
    <OnboardingLayout step={2} showBack backPath="/onboarding">
      <div className="max-w-2xl w-full">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-2">
            What's your primary goal?
          </h1>
          <p className="text-muted-foreground">
            We'll personalize your experience based on this
          </p>
        </div>

        {/* Goal Options */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
          {goals.map((goal, index) => {
            const isSelected = state.goal === goal.id;
            return (
              <motion.button
                key={goal.id}
                onClick={() => setGoal(goal.id)}
                className={cn(
                  'relative p-6 rounded-2xl border-2 text-left transition-all duration-200',
                  isSelected 
                    ? 'border-secondary bg-secondary/10 shadow-emerald-glow' 
                    : 'border-border bg-card hover:border-secondary/50'
                )}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.05 * index }}
                whileHover={{ scale: 1.02 }}
              >
                {goal.popular && (
                  <span className="absolute top-3 right-3 px-2 py-0.5 text-xs font-semibold rounded-full bg-primary text-primary-foreground">
                    Most Popular
                  </span>
                )}
                {isSelected && (
                  <span className="absolute top-3 right-3 w-6 h-6 rounded-full bg-secondary flex items-center justify-center">
                    <Check className="w-4 h-4 text-secondary-foreground" />
                  </span>
                )}
                <goal.icon className={cn(
                  'w-10 h-10 mb-3',
                  isSelected ? 'text-secondary' : 'text-secondary/70'
                )} />
                <h3 className="text-lg font-semibold text-foreground mb-1">{goal.title}</h3>
                <p className="text-sm text-muted-foreground">{goal.description}</p>
              </motion.button>
            );
          })}
        </div>

        {/* Navigation */}
        <div className="flex justify-center">
          <Button
            size="xl"
            className="w-full md:w-80 bg-gradient-gold hover:scale-105 transition-transform shadow-gold-glow disabled:opacity-50 disabled:hover:scale-100"
            disabled={!state.goal}
            onClick={() => navigate('/onboarding/plan')}
          >
            Continue
          </Button>
        </div>
      </div>
    </OnboardingLayout>
  );
};

export default OnboardingGoal;
