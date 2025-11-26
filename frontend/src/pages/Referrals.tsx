import { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { useAuth } from '@/context/AuthContext';
import { Users, DollarSign, Gift, Share2, ArrowRight } from 'lucide-react';
import { ReferralModal } from '@/features/referral/ReferralModal';

interface ReferralStats {
    total_earned: number;
    count: number;
    referral_code: string;
}

const Referrals = () => {
    const { user } = useAuth();
    const [stats, setStats] = useState<ReferralStats | null>(null);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);

    useEffect(() => {
        // Mock stats fetch
        const fetchStats = async () => {
            try {
                // Simulate API delay
                await new Promise(resolve => setTimeout(resolve, 500));
                setStats({
                    total_earned: 45.00, // Mocked for demo
                    count: 3,
                    referral_code: user?.referral_code || 'FASTHERO2024'
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

    if (loading) return <div className="p-8 text-center text-slate-400">Loading referral data...</div>;

    return (
        <div className="space-y-8 animate-fade-in">
            {/* Header Section */}
            <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
                <div>
                    <h1 className="text-3xl font-bold text-white flex items-center gap-2">
                        <Gift className="w-8 h-8 text-purple-400" />
                        Referral Program
                    </h1>
                    <p className="text-slate-400 mt-1">
                        Invite friends to the Vault. You both get paid.
                    </p>
                </div>
                <Button
                    onClick={() => setShowModal(true)}
                    className="bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white shadow-lg shadow-purple-900/20"
                >
                    <Share2 className="w-4 h-4 mr-2" /> Invite Friends
                </Button>
            </div>

            {/* Stats Grid */}
            <div className="grid gap-4 md:grid-cols-3">
                <Card className="bg-slate-900 border-slate-800">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium text-slate-400">Total Earned</CardTitle>
                        <DollarSign className="h-4 w-4 text-emerald-400" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-3xl font-bold text-white">${stats?.total_earned.toFixed(2)}</div>
                        <p className="text-xs text-emerald-400 mt-1">
                            +12% from last month
                        </p>
                    </CardContent>
                </Card>

                <Card className="bg-slate-900 border-slate-800">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium text-slate-400">Friends Invited</CardTitle>
                        <Users className="h-4 w-4 text-blue-400" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-3xl font-bold text-white">{stats?.count}</div>
                        <p className="text-xs text-slate-500 mt-1">
                            Active Vault members
                        </p>
                    </CardContent>
                </Card>

                <Card className="bg-gradient-to-br from-purple-900/20 to-blue-900/20 border-purple-500/30">
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium text-purple-300">Next Reward</CardTitle>
                        <Gift className="h-4 w-4 text-purple-400" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-xl font-bold text-white mb-2">Free Month</div>
                        <div className="w-full bg-slate-800 rounded-full h-1.5">
                            <div className="bg-purple-500 h-1.5 rounded-full w-[60%]"></div>
                        </div>
                        <p className="text-xs text-purple-300 mt-2">
                            2 more referrals to unlock
                        </p>
                    </CardContent>
                </Card>
            </div>

            {/* How it Works */}
            <Card className="bg-slate-900/50 border-slate-800">
                <CardHeader>
                    <CardTitle className="text-white">How It Works</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid gap-6 md:grid-cols-3">
                        <div className="flex flex-col items-center text-center p-4 bg-slate-900 rounded-lg border border-slate-800">
                            <div className="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center border border-slate-700 mb-3 text-lg font-bold text-white">1</div>
                            <h3 className="font-semibold text-white mb-2">Share Your Link</h3>
                            <p className="text-sm text-slate-400">Send your unique referral link to friends who want to build discipline.</p>
                        </div>
                        <div className="flex flex-col items-center text-center p-4 bg-slate-900 rounded-lg border border-slate-800">
                            <div className="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center border border-slate-700 mb-3 text-lg font-bold text-white">2</div>
                            <h3 className="font-semibold text-white mb-2">They Join the Vault</h3>
                            <p className="text-sm text-slate-400">When they make their first $30 deposit, they get $10 off instantly.</p>
                        </div>
                        <div className="flex flex-col items-center text-center p-4 bg-slate-900 rounded-lg border border-slate-800">
                            <div className="w-10 h-10 rounded-full bg-slate-800 flex items-center justify-center border border-slate-700 mb-3 text-lg font-bold text-white">3</div>
                            <h3 className="font-semibold text-white mb-2">You Get Paid</h3>
                            <p className="text-sm text-slate-400">You receive $10 in Vault credit for every successful referral.</p>
                        </div>
                    </div>
                </CardContent>
            </Card>

            {/* CTA */}
            <div className="bg-gradient-to-r from-purple-900/40 to-blue-900/40 border border-purple-500/20 rounded-xl p-8 text-center">
                <h2 className="text-2xl font-bold text-white mb-4">Ready to grow your tribe?</h2>
                <Button
                    size="lg"
                    onClick={() => setShowModal(true)}
                    className="bg-white text-purple-900 hover:bg-slate-100 font-bold"
                >
                    Start Referring Now <ArrowRight className="w-4 h-4 ml-2" />
                </Button>
            </div>

            <ReferralModal open={showModal} onOpenChange={setShowModal} />
        </div>
    );
};

export default Referrals;
