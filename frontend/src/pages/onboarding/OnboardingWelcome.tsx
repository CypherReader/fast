import { motion } from 'framer-motion';
import { Lock, TrendingUp, DollarSign } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import OnboardingLayout from '@/components/onboarding/OnboardingLayout';
import { Button } from '@/components/ui/button';

const valueProps = [
  {
    icon: Lock,
    title: 'Deposit $20',
    description: 'Put your money where your goals are',
    color: 'text-secondary',
  },
  {
    icon: TrendingUp,
    title: 'Fast Consistently',
    description: 'Complete fasts to recover your deposit',
    color: 'text-secondary',
  },
  {
    icon: DollarSign,
    title: 'Get Refunded',
    description: '87% of users pay $0/month',
    color: 'text-primary',
  },
];

const OnboardingWelcome = () => {
  const navigate = useNavigate();

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
          Let's set up your Commitment Vault in 2 minutes
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
