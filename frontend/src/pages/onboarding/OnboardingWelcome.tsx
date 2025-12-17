import { motion } from 'framer-motion';
import { Lock, TrendingUp, DollarSign } from 'lucide-react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useEffect } from 'react';
import OnboardingLayout from '@/components/onboarding/OnboardingLayout';
import { Button } from '@/components/ui/button';

// HIDDEN FOR V2: Vault deposit value props
const valueProps = [
  {
    icon: TrendingUp,
    title: 'Set Your Goals',
    description: 'Choose your personalized fasting plan',
    color: 'text-secondary',
  },
  {
    icon: Lock,
    title: 'Track Progress',
    description: 'Monitor your fasting journey daily',
    color: 'text-secondary',
  },
  {
    icon: DollarSign,
    title: 'Build Discipline',
    description: 'Develop sustainable healthy habits',
    color: 'text-primary',
  },
];

const OnboardingWelcome = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  // Handle OAuth token from URL
  useEffect(() => {
    const token = searchParams.get('token');
    if (token) {
      localStorage.setItem('token', token);
      // Remove token from URL for security
      searchParams.delete('token');
      window.history.replaceState({}, '', `/onboarding?${searchParams.toString()}`);

      // Auto-redirect OAuth users to goal page (skip welcome since account is already created)
      setTimeout(() => navigate('/onboarding/goal'), 100);
    }
  }, [searchParams, navigate]);

  return (
    <OnboardingLayout step={1}>
      <div className="max-w-xl w-full text-center">
        {/* Vault Icon */}
        <motion.div
          className="mb-8 inline-flex items-center justify-center w-28 h-28 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 border border-primary/30"
          animate={{ scale: [1, 1.05, 1] }}
          transition={{ duration: 2, repeat: Infinity, ease: 'easeInOut' }}
        >
          <Lock className="w-14 h-14 text-primary" />
        </motion.div>

        {/* Headlines */}
        <h1 className="text-4xl md:text-5xl font-extrabold text-foreground mb-4">
          Welcome to FastingHero
        </h1>
        <p className="text-lg md:text-xl text-muted-foreground mb-10">
          Let's set up your fasting journey in 2 minutes
        </p>

        {/* Value Props */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-10">
          {valueProps.map((prop, index) => (
            <motion.div
              key={prop.title}
              className="bg-card border border-border rounded-xl p-6 hover:border-secondary/50 hover:scale-[1.02] transition-all duration-200"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 * (index + 1) }}
            >
              <prop.icon className={`w-10 h-10 ${prop.color} mb-3 mx-auto`} />
              <h3 className="font-semibold text-foreground mb-1">{prop.title}</h3>
              <p className="text-sm text-muted-foreground">{prop.description}</p>
            </motion.div>
          ))}
        </div>

        {/* Social proof - Avatar stack */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.35 }}
          className="flex items-center justify-center gap-3 mb-8"
        >
          <div className="flex -space-x-2">
            {[...Array(5)].map((_, i) => (
              <div
                key={i}
                className="w-8 h-8 rounded-full bg-gradient-to-br from-primary/80 to-secondary/80 border-2 border-background flex items-center justify-center text-xs font-semibold text-primary-foreground"
              />
            ))}
          </div>
          <span className="text-sm text-muted-foreground">
            {/* HIDDEN FOR V2: 12,847 people paying $0/month */}
            <span className="text-primary font-semibold">12,847</span> active fasters
          </span>
        </motion.div>

        {/* CTA */}
        <Button
          size="xl"
          className="w-full md:w-80 bg-gradient-gold hover:scale-105 transition-transform shadow-gold-glow"
          onClick={() => navigate('/onboarding/goal')}
        >
          Get Started
        </Button>

        {/* Skip Link */}
        <p className="mt-4">
          <button
            className="text-sm text-muted-foreground hover:text-foreground hover:underline transition-colors"
            onClick={() => navigate('/')}
          >
            I'll set this up later
          </button>
        </p>
      </div>
    </OnboardingLayout>
  );
};

export default OnboardingWelcome;
