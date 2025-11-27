import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Lock, ArrowDown, Timer, Calendar, CheckCircle, Shield, CreditCard, RefreshCw } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import OnboardingLayout from '@/components/onboarding/OnboardingLayout';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { cn } from '@/lib/utils';

const animationFrames = [
  {
    id: 1,
    title: 'You deposit $20',
    description: 'At the start of each month, your commitment goes into the Vault',
    icon: Lock,
  },
  {
    id: 2,
    title: 'Complete your fasts',
    description: 'Each completed fast recovers $2 of your deposit',
    icon: Timer,
  },
  {
    id: 3,
    title: '10 fasts = Full refund',
    description: 'Complete 10 fasts and get your entire $20 back',
    icon: Calendar,
  },
  {
    id: 4,
    title: 'You control the cost',
    description: 'Disciplined = $0 | Undisciplined = $20',
    icon: CheckCircle,
  },
];

const stats = [
  { value: '87%', label: 'Get full refund', sublabel: 'Most users pay $0' },
  { value: '$1.80', label: 'Avg net cost', sublabel: 'Per user per month' },
  { value: '10', label: 'Fasts needed', sublabel: '~3 per week' },
];

const trustSignals = [
  { icon: Shield, text: 'Stripe-secured payments' },
  { icon: CheckCircle, text: 'Cancel anytime' },
  { icon: CreditCard, text: 'First month pro-rated' },
];

const OnboardingVault = () => {
  const navigate = useNavigate();
  const [currentFrame, setCurrentFrame] = useState(0);
  const [isPlaying, setIsPlaying] = useState(true);
  const [understood, setUnderstood] = useState(false);
  const [moneyAmount, setMoneyAmount] = useState(20);

  // Auto-advance animation frames
  useEffect(() => {
    if (!isPlaying) return;

    const interval = setInterval(() => {
      setCurrentFrame(prev => {
        const next = (prev + 1) % animationFrames.length;
        // Update money amount based on frame
        if (next === 0) setMoneyAmount(20);
        else if (next === 1) setMoneyAmount(18);
        else if (next === 2) setMoneyAmount(0);
        else setMoneyAmount(0);
        return next;
      });
    }, 3000);

    return () => clearInterval(interval);
  }, [isPlaying]);

  const frame = animationFrames[currentFrame];
  const FrameIcon = frame.icon;

  return (
    <OnboardingLayout step={4} showBack backPath="/onboarding/plan">
      <div className="max-w-2xl w-full">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-2">
            Here's How You Get Your Money Back
          </h1>
          <p className="text-muted-foreground">
            The Commitment Vault explained in 30 seconds
          </p>
        </div>

        {/* Animated Explainer */}
        <div className="relative bg-card border border-border rounded-2xl p-8 mb-8 overflow-hidden">
          {/* Background glow effect */}
          <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-transparent to-secondary/5" />

          {/* Money Flow Animation */}
          <div className="relative z-10">
            {/* Vault Visual */}
            <div className="flex flex-col items-center mb-8">
              {/* Money dropping into vault animation */}
              <AnimatePresence mode="wait">
                <motion.div
                  key={currentFrame}
                  initial={{ opacity: 0, y: -20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: 20 }}
                  transition={{ duration: 0.5 }}
                  className="flex flex-col items-center"
                >
                  {/* Animated Icon */}
                  <motion.div
                    className={cn(
                      'w-24 h-24 rounded-full flex items-center justify-center mb-4',
                      currentFrame === 3 ? 'bg-secondary/20' : 'bg-primary/20'
                    )}
                    animate={{ scale: [1, 1.1, 1] }}
                    transition={{ duration: 2, repeat: Infinity }}
                  >
                    <FrameIcon className={cn(
                      'w-12 h-12',
                      currentFrame === 3 ? 'text-secondary' : 'text-primary'
                    )} />
                  </motion.div>

                  {/* Money Counter */}
                  <motion.div
                    className="text-5xl font-bold mb-2"
                    key={moneyAmount}
                    initial={{ scale: 1.2, opacity: 0 }}
                    animate={{ scale: 1, opacity: 1 }}
                  >
                    <span className={moneyAmount === 0 ? 'text-secondary' : 'text-primary'}>
                      ${moneyAmount}
                    </span>
                    <span className="text-muted-foreground text-2xl"> / $20</span>
                  </motion.div>

                  {/* Progress Bar */}
                  <div className="w-64 h-3 bg-muted rounded-full overflow-hidden mb-4">
                    <motion.div
                      className="h-full bg-gradient-to-r from-primary to-secondary"
                      initial={{ width: '100%' }}
                      animate={{ width: `${((20 - moneyAmount) / 20) * 100}%` }}
                      transition={{ duration: 0.5 }}
                    />
                  </div>

                  {/* Frame Title */}
                  <h3 className="text-xl font-semibold text-foreground mb-1">
                    {frame.title}
                  </h3>
                  <p className="text-muted-foreground text-center max-w-sm">
                    {frame.description}
                  </p>
                </motion.div>
              </AnimatePresence>
            </div>

            {/* Frame Indicators */}
            <div className="flex justify-center gap-2 mb-4">
              {animationFrames.map((_, index) => (
                <button
                  key={index}
                  onClick={() => {
                    setCurrentFrame(index);
                    setIsPlaying(false);
                    if (index === 0) setMoneyAmount(20);
                    else if (index === 1) setMoneyAmount(18);
                    else setMoneyAmount(0);
                  }}
                  className={cn(
                    'w-2.5 h-2.5 rounded-full transition-all',
                    index === currentFrame
                      ? 'bg-primary w-8'
                      : 'bg-muted-foreground/30 hover:bg-muted-foreground/50'
                  )}
                />
              ))}
            </div>

            {/* Replay Button */}
            {!isPlaying && (
              <motion.button
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                onClick={() => {
                  setCurrentFrame(0);
                  setMoneyAmount(20);
                  setIsPlaying(true);
                }}
                className="flex items-center gap-2 mx-auto text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                <RefreshCw className="w-4 h-4" />
                Replay animation
              </motion.button>
            )}
          </div>
        </div>

        {/* Key Stats */}
        <div className="grid grid-cols-3 gap-4 mb-8">
          {stats.map((stat, index) => (
            <motion.div
              key={stat.label}
              className="bg-card border border-border rounded-xl p-4 text-center"
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.1 * index }}
            >
              <div className="text-2xl md:text-3xl font-bold text-primary mb-1">
                {stat.value}
              </div>
              <div className="text-sm font-medium text-foreground">{stat.label}</div>
              <div className="text-xs text-muted-foreground">{stat.sublabel}</div>
            </motion.div>
          ))}
        </div>

        {/* Trust Signals */}
        <div className="flex flex-wrap justify-center gap-4 mb-8">
          {trustSignals.map((signal) => (
            <div key={signal.text} className="flex items-center gap-2 text-sm text-muted-foreground">
              <signal.icon className="w-4 h-4 text-secondary" />
              <span>{signal.text}</span>
            </div>
          ))}
        </div>

        {/* Understanding Checkbox */}
        <div className="flex items-start gap-3 p-4 bg-card border border-border rounded-xl mb-6">
          <Checkbox
            id="understand"
            checked={understood}
            onCheckedChange={(checked) => setUnderstood(checked as boolean)}
            className="mt-0.5"
          />
          <label
            htmlFor="understand"
            className="text-sm text-muted-foreground cursor-pointer leading-relaxed"
          >
            I understand my deposit is fully refundable based on my fasting consistency.
            Complete 10 fasts = get all $20 back.
          </label>
        </div>

        {/* Navigation */}
        <div className="flex justify-center">
          <Button
            size="xl"
            className="w-full md:w-80 bg-gradient-gold hover:scale-105 transition-transform shadow-gold-glow disabled:opacity-50 disabled:hover:scale-100"
            disabled={!understood}
            onClick={() => navigate('/onboarding/account')}
          >
            I Understand, Continue
          </Button>
        </div>
      </div>
    </OnboardingLayout>
  );
};

export default OnboardingVault;
