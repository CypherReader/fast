import { useState, useEffect } from 'react';
import { useNavigate, Link, useSearchParams } from 'react-router-dom';
import { authApi } from '../api/client';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useFeatureFlag, FEATURE_FLAGS } from '@/lib/flags';
import { analytics } from '@/lib/analytics';

const Onboarding = () => {
    const [step, setStep] = useState(1);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [referralCode, setReferralCode] = useState('');
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();

    const showQuiz = useFeatureFlag(FEATURE_FLAGS.ONBOARDING_QUIZ_FLOW, false);

    useEffect(() => {
        const ref = searchParams.get('ref');
        if (ref) {
            setReferralCode(ref);
        }

        // Track initial view
        analytics.trackEvent('onboarding_start', { flow: showQuiz ? 'quiz' : 'direct' });
    }, [searchParams, showQuiz]);

    const handleRegister = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            analytics.trackEvent('registration_submit');
            await authApi.register(email, password, referralCode);
            analytics.trackEvent('registration_success');
            alert('Registration successful! Please login.');
            navigate('/login');
        } catch (error) {
            analytics.trackEvent('registration_failed');
            alert('Registration failed');
        }
    };

    const handleQuizNext = () => {
        analytics.trackEvent('quiz_step_complete', { step });
        setStep(step + 1);
    };

    if (showQuiz && step < 3) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-slate-950 p-4">
                <Card className="w-full max-w-md bg-slate-900 border-slate-800">
                    <CardHeader>
                        <CardTitle className="text-2xl text-center text-white">
                            {step === 1 ? "What's your primary goal?" : "What's your experience?"}
                        </CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        {step === 1 ? (
                            <div className="grid gap-3">
                                {['Weight Loss', 'Mental Clarity', 'Autophagy', 'Discipline'].map(goal => (
                                    <Button key={goal} variant="outline" className="justify-start h-12 text-lg" onClick={handleQuizNext}>
                                        {goal}
                                    </Button>
                                ))}
                            </div>
                        ) : (
                            <div className="grid gap-3">
                                {['Beginner', 'Intermediate', 'Advanced'].map(level => (
                                    <Button key={level} variant="outline" className="justify-start h-12 text-lg" onClick={handleQuizNext}>
                                        {level}
                                    </Button>
                                ))}
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className="min-h-screen flex items-center justify-center bg-background p-4">
            <Card className="w-full max-w-md">
                <CardHeader>
                    <CardTitle className="text-2xl text-center">Join FastingHero</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleRegister} className="space-y-4">
                        <Input
                            type="email"
                            placeholder="Email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />
                        <Input
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                        <Input
                            type="text"
                            placeholder="Referral Code (Optional)"
                            value={referralCode}
                            onChange={(e) => setReferralCode(e.target.value)}
                        />
                        <Button type="submit" className="w-full">Register</Button>
                        <div className="text-center text-sm">
                            Already have an account? <Link to="/login" className="text-primary hover:underline">Login</Link>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
};

export default Onboarding;
