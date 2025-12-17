import { useState } from 'react';
import { motion } from 'framer-motion';
import { Lock, CreditCard, Shield, CheckCircle, Building2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import OnboardingLayout from '@/components/onboarding/OnboardingLayout';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useOnboarding } from '@/contexts/OnboardingContext';
import { cn } from '@/lib/utils';
import { api } from '@/api/client';

const trustBadges = [
  { icon: Lock, text: 'Secured by Stripe' },
  { icon: Shield, text: '256-bit SSL' },
  { icon: CheckCircle, text: 'PCI Compliant' },
  { icon: Building2, text: 'Bank-level security' },
];

const planNames: Record<string, string> = {
  '16:8': '16:8 Time-Restricted Eating',
  '18:6': '18:6 Extended Fasting',
  'omad': '23:1 One Meal A Day',
};

const goalNames: Record<string, string> = {
  'weight_loss': 'Weight Loss',
  'metabolic': 'Metabolic Health',
  'discipline': 'Discipline Building',
  'longevity': 'Longevity',
};

const OnboardingPayment = () => {
  const navigate = useNavigate();
  const { state } = useOnboarding();

  const [cardNumber, setCardNumber] = useState('');
  const [expiry, setExpiry] = useState('');
  const [cvc, setCvc] = useState('');
  const [zip, setZip] = useState('');
  const [agreeRefund, setAgreeRefund] = useState(false);
  const [agreeTerms, setAgreeTerms] = useState(false);
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState('');

  const formatCardNumber = (value: string) => {
    const v = value.replace(/\s+/g, '').replace(/[^0-9]/gi, '');
    const matches = v.match(/\d{4,16}/g);
    const match = (matches && matches[0]) || '';
    const parts = [];
    for (let i = 0, len = match.length; i < len; i += 4) {
      parts.push(match.substring(i, i + 4));
    }
    return parts.length ? parts.join(' ') : v;
  };

  const formatExpiry = (value: string) => {
    const v = value.replace(/\s+/g, '').replace(/[^0-9]/gi, '');
    if (v.length >= 2) {
      return v.substring(0, 2) + ' / ' + v.substring(2, 4);
    }
    return v;
  };

  const handleSubmit = async () => {
    setError('');

    // Basic validation
    if (cardNumber.replace(/\s/g, '').length < 16) {
      setError('Please enter a valid card number');
      return;
    }
    if (expiry.replace(/\s|\//g, '').length < 4) {
      setError('Please enter a valid expiry date');
      return;
    }
    if (cvc.length < 3) {
      setError('Please enter a valid CVC');
      return;
    }
    if (zip.length < 5) {
      setError('Please enter a valid ZIP code');
      return;
    }

    if (!state.email || !state.password) {
      setError('Account details missing. Please go back and create an account.');
      return;
    }

    setIsProcessing(true);

    try {
      // 1. Register User
      const registerResponse = await api.post<{ token: string }>('/auth/register', {
        email: state.email,
        password: state.password,
        name: state.name,
      });

      // 2. Store Token
      localStorage.setItem('token', registerResponse.data.token);

      // 3. Update Profile (Goal & Plan)
      await api.put('/onboarding/profile', {
        goal: state.goal,
        fasting_plan: state.fastingPlan,
      });

      // 4. Process Subscription Payment (Mock payment, real backend update)
      await api.post('/payments/deposit', {
        amount: 4.99,
        payment_method_id: 'pm_card_visa', // Mock Stripe ID
      });

      navigate('/onboarding/success');
    } catch (err: unknown) {
      console.error('Payment/Registration failed:', err);
      const error = err as { response?: { data?: { error?: string } } };
      setError(error.response?.data?.error || 'Failed to process payment. Please try again.');
    } finally {
      setIsProcessing(false);
    }
  };

  const isFormValid =
    cardNumber.replace(/\s/g, '').length >= 16 &&
    expiry.replace(/\s|\//g, '').length >= 4 &&
    cvc.length >= 3 &&
    zip.length >= 5 &&
    agreeRefund &&
    agreeTerms;

  return (
    <OnboardingLayout step={4} showBack backPath="/onboarding/plan">
      <div className="max-w-4xl w-full">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-2">
            Set Up Payment Method
          </h1>
          <p className="text-muted-foreground">
            Subscribe to FastingHero - $4.99/month
          </p>
        </div>

        <div className="grid md:grid-cols-5 gap-8">
          {/* Payment Form - Left Side */}
          <div className="md:col-span-3 order-2 md:order-1">
            {/* Trust Badges */}
            <div className="flex flex-wrap gap-3 mb-6">
              {trustBadges.map((badge) => (
                <div
                  key={badge.text}
                  className="flex items-center gap-1.5 text-xs text-muted-foreground bg-card px-3 py-1.5 rounded-full border border-border"
                >
                  <badge.icon className="w-3.5 h-3.5" />
                  <span>{badge.text}</span>
                </div>
              ))}
            </div>

            {/* Card Form */}
            <div className="space-y-4">
              {/* Card Number */}
              <div className="space-y-2">
                <Label htmlFor="cardNumber" className="text-muted-foreground">
                  Card number <span className="text-destructive">*</span>
                </Label>
                <div className="relative">
                  <Input
                    id="cardNumber"
                    placeholder="1234 5678 9012 3456"
                    value={cardNumber}
                    onChange={(e) => setCardNumber(formatCardNumber(e.target.value))}
                    maxLength={19}
                    className="bg-card border-border h-12 pl-12 text-foreground placeholder:text-muted-foreground/50"
                  />
                  <CreditCard className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground" />
                </div>
              </div>

              {/* Expiry & CVC */}
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="expiry" className="text-muted-foreground">
                    Expiry date <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="expiry"
                    placeholder="MM / YY"
                    value={expiry}
                    onChange={(e) => setExpiry(formatExpiry(e.target.value))}
                    maxLength={7}
                    className="bg-card border-border h-12 text-foreground placeholder:text-muted-foreground/50"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="cvc" className="text-muted-foreground">
                    CVC <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="cvc"
                    placeholder="123"
                    value={cvc}
                    onChange={(e) => setCvc(e.target.value.replace(/\D/g, '').slice(0, 4))}
                    maxLength={4}
                    className="bg-card border-border h-12 text-foreground placeholder:text-muted-foreground/50"
                  />
                </div>
              </div>

              {/* ZIP Code */}
              <div className="space-y-2">
                <Label htmlFor="zip" className="text-muted-foreground">
                  ZIP code <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="zip"
                  placeholder="12345"
                  value={zip}
                  onChange={(e) => setZip(e.target.value.replace(/\D/g, '').slice(0, 5))}
                  maxLength={5}
                  className="bg-card border-border h-12 text-foreground placeholder:text-muted-foreground/50"
                />
              </div>

              {/* Legal Checkboxes */}
              <div className="space-y-3 pt-4">
                <div className="flex items-start gap-3">
                  <Checkbox
                    id="agreeRefund"
                    checked={agreeRefund}
                    onCheckedChange={(checked) => setAgreeRefund(checked as boolean)}
                    className="mt-0.5"
                  />
                  <label
                    htmlFor="agreeRefund"
                    className="text-sm text-muted-foreground cursor-pointer leading-relaxed"
                  >
                    I understand this is a monthly subscription for $4.99/month
                  </label>
                </div>
                <div className="flex items-start gap-3">
                  <Checkbox
                    id="agreeTerms"
                    checked={agreeTerms}
                    onCheckedChange={(checked) => setAgreeTerms(checked as boolean)}
                    className="mt-0.5"
                  />
                  <label
                    htmlFor="agreeTerms"
                    className="text-sm text-muted-foreground cursor-pointer leading-relaxed"
                  >
                    I agree to the{' '}
                    <a href="#" className="text-secondary underline hover:no-underline">Terms of Service</a>
                    {' '}and{' '}
                    <a href="#" className="text-secondary underline hover:no-underline">Privacy Policy</a>
                  </label>
                </div>
              </div>

              {/* Error Message */}
              {error && (
                <motion.p
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="text-sm text-destructive"
                >
                  {error}
                </motion.p>
              )}

              {/* Submit Button */}
              <Button
                size="xl"
                className={cn(
                  'w-full bg-gradient-gold transition-all shadow-gold-glow',
                  isFormValid && !isProcessing ? 'hover:scale-[1.02]' : 'opacity-50 cursor-not-allowed'
                )}
                disabled={!isFormValid || isProcessing}
                onClick={handleSubmit}
              >
                {isProcessing ? (
                  <span className="flex items-center gap-2">
                    <motion.div
                      className="w-5 h-5 border-2 border-primary-foreground/30 border-t-primary-foreground rounded-full"
                      animate={{ rotate: 360 }}
                      transition={{ duration: 1, repeat: Infinity, ease: 'linear' }}
                    />
                    Processing...
                  </span>
                ) : (
                  'Complete Payment'
                )}
              </Button>

              {/* Money-back Guarantee */}
              <div className="flex items-center gap-3 p-4 bg-secondary/10 border border-secondary/30 rounded-xl">
                <CheckCircle className="w-6 h-6 text-secondary flex-shrink-0" />
                <div>
                  <p className="text-sm font-medium text-foreground">30-Day Money-Back Guarantee</p>
                  <p className="text-xs text-muted-foreground">
                    Not happy? Get a full refund within 30 days, no questions asked.
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Order Summary - Right Side */}
          <div className="md:col-span-2 order-1 md:order-2">
            <div className="bg-card border-2 border-secondary/30 rounded-2xl p-6 sticky top-8">
              <h3 className="text-lg font-semibold text-foreground mb-4">Your Plan Summary</h3>

              <div className="space-y-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Fasting schedule</span>
                  <span className="text-foreground font-medium">
                    {state.fastingPlan ? planNames[state.fastingPlan] || state.fastingPlan : '16:8'}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Goal</span>
                  <span className="text-foreground font-medium">
                    {state.goal ? goalNames[state.goal] || state.goal : 'Weight Loss'}
                  </span>
                </div>
              </div>

              <div className="border-t border-border my-4" />

              <div className="space-y-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Monthly Subscription</span>
                  <span className="text-foreground font-medium">$4.99</span>
                </div>
                {/* HIDDEN FOR V2: Vault recovery system */}
                {/* <div className="flex justify-between">
                  <span className="text-muted-foreground">Potential Recovery</span>
                  <span className="text-secondary font-medium">Up to $20.00</span>
                </div> */}
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Billing Cycle</span>
                  <span className="text-foreground font-medium">Monthly</span>
                </div>
              </div>

              <div className="border-t border-border my-4" />

              <div className="space-y-2">
                <div className="flex justify-between items-baseline">
                  <span className="text-lg font-semibold text-foreground">First charge</span>
                  <span className="text-2xl font-bold text-primary">$4.99</span>
                </div>
                <p className="text-xs text-muted-foreground">
                  Recurring monthly subscription
                </p>
              </div>

              {/* Payment info */}
              <div className="mt-6 p-4 bg-muted/50 rounded-xl text-center">
                <Lock className="w-8 h-8 text-primary mx-auto mb-2" />
                <p className="text-xs text-muted-foreground">
                  Secure subscription powered by Stripe
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </OnboardingLayout>
  );
};

export default OnboardingPayment;
