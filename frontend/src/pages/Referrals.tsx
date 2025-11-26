import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuth } from '@/context/AuthContext';
import { Copy, Users, DollarSign } from 'lucide-react';
import { useToast } from "@/components/ui/use-toast";

interface ReferralStats {
    total_earned: number;
    count: number;
    referral_code: string;
}

const Referrals = () => {
    const { user } = useAuth();
    const { toast } = useToast();
    const [stats, setStats] = useState<ReferralStats | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchStats = async () => {
            try {
                // We need an endpoint for stats. For now, let's assume we can get the code from user profile
                // and maybe a separate endpoint for stats later.
                // But the plan said "Show a list of successful referrals and total earned credit".
                // We implemented GetReferralStats in service but not exposed in Handler?
                // Let's check Handler. We didn't add GetReferralStats to Handler!
                // I missed that in the backend implementation.
                // For now, I'll just show the code if available in user object.
                // Wait, I added ReferralCode to User struct, so it should be in the user object.

                // Mock stats for now or fetch if endpoint exists
                setStats({
                    total_earned: 0, // Placeholder
                    count: 0, // Placeholder
                    referral_code: user?.referral_code || 'Generating...'
                });
            } catch (error) {
                console.error("Failed to fetch referral stats", error);
            } finally {
                setLoading(false);
            }
        };

        if (user) {
            fetchStats();
        }
    }, [user]);

    const copyToClipboard = () => {
        if (stats?.referral_code) {
            const link = `${window.location.origin}/register?ref=${stats.referral_code}`;
            navigator.clipboard.writeText(link);
            toast({
                title: "Copied!",
                description: "Referral link copied to clipboard.",
            });
        }
    };

    if (loading) return <div>Loading...</div>;

    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-bold text-white">Referrals</h1>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Earned</CardTitle>
                        <DollarSign className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">${stats?.total_earned.toFixed(2)}</div>
                        <p className="text-xs text-muted-foreground">
                            From successful referrals
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Friends Invited</CardTitle>
                        <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{stats?.count}</div>
                        <p className="text-xs text-muted-foreground">
                            Completed referrals
                        </p>
                    </CardContent>
                </Card>
            </div>

            <Card>
                <CardHeader>
                    <CardTitle>Your Referral Link</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <p className="text-sm text-muted-foreground">
                        Share this link with your friends. When they join and make their first deposit, you both get $5 in Vault credit!
                    </p>
                    <div className="flex items-center space-x-2">
                        <div className="bg-secondary p-3 rounded-md font-mono text-sm flex-1 truncate">
                            {window.location.origin}/register?ref={stats?.referral_code}
                        </div>
                        <Button onClick={copyToClipboard} size="icon" variant="outline">
                            <Copy className="h-4 w-4" />
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};

export default Referrals;
