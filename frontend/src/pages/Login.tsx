import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { AnimatePresence } from 'framer-motion';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { useToast } from '@/hooks/use-toast';
import { api } from '@/api/client';
import { LoadingTransition } from '@/components/LoadingTransition';

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [showTransition, setShowTransition] = useState(false);
    const navigate = useNavigate();
    const { toast } = useToast();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);

        try {
            const response = await api.post('/auth/login', { email, password });
            const { token, user } = response.data;

            localStorage.setItem('token', token);
            localStorage.setItem('user', JSON.stringify(user));

            toast({
                title: "Welcome back!",
                description: "You have successfully logged in.",
            });

            // Show video transition instead of navigating immediately
            setShowTransition(true);
        } catch (error: unknown) {
            const axiosError = error as { response?: { data?: { error?: string } } };
            toast({
                variant: "destructive",
                title: "Login failed",
                description: axiosError.response?.data?.error || "Invalid credentials. Please try again.",
            });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <>
            <AnimatePresence>
                {showTransition && (
                    <LoadingTransition
                        onComplete={() => navigate('/dashboard')}
                    />
                )}
            </AnimatePresence>

            <div className="min-h-screen flex items-center justify-center bg-background p-4">
                <Card className="w-full max-w-md">
                    <CardHeader>
                        <CardTitle className="text-2xl font-bold text-center">Welcome Back</CardTitle>
                        <CardDescription className="text-center">
                            Enter your credentials to access your vault
                        </CardDescription>
                    </CardHeader>
                    <form onSubmit={handleLogin}>
                        <CardContent className="space-y-4">
                            <div className="space-y-2">
                                <Label htmlFor="email">Email</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="name@example.com"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    required
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="password">Password</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    required
                                />
                            </div>
                        </CardContent>
                        <CardFooter className="flex flex-col gap-4">
                            <Button type="submit" className="w-full" disabled={isLoading}>
                                {isLoading ? "Logging in..." : "Login"}
                            </Button>
                            <Button variant="link" className="text-sm text-muted-foreground" onClick={() => navigate('/onboarding')}>
                                Don't have an account? Start here
                            </Button>
                        </CardFooter>
                    </form>
                </Card>
            </div>
        </>
    );
};

export default Login;
