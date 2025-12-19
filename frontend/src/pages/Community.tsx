import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Users, MessageSquare, Heart, Trophy, Flame, Search, Siren, Sparkles } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useUser } from '@/hooks/use-user';
import { useLeaderboard } from '@/hooks/use-leaderboard';
import { useSocial } from '@/hooks/use-social';
import { useMyTribes } from '@/hooks/use-tribes';

const Community = () => {
    const navigate = useNavigate();
    const { user } = useUser();
    const { leaderboard, isLoading: isLoadingLeaderboard } = useLeaderboard();
    const { tribes, feed, isLoadingTribes, isLoadingFeed } = useSocial();
    const { data: myTribesData, isLoading: isLoadingMyTribes } = useMyTribes();

    // Helper to get action text for feed items
    const getActionText = (eventType: string, data: Record<string, unknown>) => {
        switch (eventType) {
            case 'fast_completed':
                return 'completed a fast';
            case 'tribe_joined':
                return `joined ${data.tribe_name || 'a tribe'}`;
            case 'sos_flare':
                return 'sent an SOS flare 🆘';
            case 'hype_sent':
                return `sent hype to ${data.recipient_name || 'a tribe member'}`;
            case 'challenge_won':
                return 'won a challenge';
            default:
                return 'had activity';
        }
    };

    return (
        <div className="min-h-screen bg-background">
            {/* Header */}
            <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-40">
                <div className="container mx-auto px-4 py-4 flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <Button variant="ghost" size="icon" onClick={() => navigate('/dashboard')}>
                            <ArrowLeft className="w-5 h-5" />
                        </Button>
                        <div className="flex items-center gap-2">
                            <img src="/fasthero.png" alt="FastingHero" className="w-6 h-6 rounded-lg" />
                            <h1 className="font-bold text-lg text-foreground">Community</h1>
                        </div>
                    </div>
                    <div className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center">
                        <span className="text-sm font-medium text-primary">
                            {user?.name ? user.name.charAt(0).toUpperCase() : 'U'}
                        </span>
                    </div>
                </div>
            </header>

            <main className="container mx-auto px-4 py-6 max-w-4xl space-y-8">
                {/* Tribes CTA Banner */}
                <div className="bg-gradient-to-r from-purple-500/10 to-indigo-600/10 border border-purple-500/20 rounded-2xl p-6">
                    <div className="flex items-center justify-between">
                        <div>
                            <h2 className="font-bold text-lg mb-1">Join or Create a Tribe</h2>
                            <p className="text-sm text-muted-foreground">Connect with like-minded fasters and stay accountable together</p>
                        </div>
                        <Button onClick={() => navigate('/tribes')} size="lg">
                            Explore Tribes →
                        </Button>
                    </div>
                </div>

                {/* Search */}
                <div className="relative">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                    <Input placeholder="Find tribes or people..." className="pl-10 bg-card border-border" />
                </div>

                <div className="grid md:grid-cols-3 gap-6">
                    {/* Main Feed - Left Side */}
                    <div className="md:col-span-2 space-y-6">
                        <h2 className="font-semibold text-foreground flex items-center gap-2">
                            <Flame className="w-5 h-5 text-orange-500" />
                            Activity Feed
                        </h2>

                        <div className="space-y-4">
                            {isLoadingFeed ? (
                                <div className="text-center text-muted-foreground">Loading feed...</div>
                            ) : feed && feed.length > 0 ? (
                                feed.map((item, index) => {
                                    const data = item.data ? JSON.parse(item.data) : {};
                                    const actionText = getActionText(item.event_type, data);
                                    const isSOS = item.event_type.includes('sos');
                                    return (
                                        <motion.div
                                            key={item.id}
                                            className={`bg-card border rounded-xl p-4 ${isSOS ? 'border-red-500/50 bg-red-500/5' : 'border-border'}`}
                                            initial={{ opacity: 0, y: 20 }}
                                            animate={{ opacity: 1, y: 0 }}
                                            transition={{ delay: index * 0.1 }}
                                        >
                                            <div className="flex items-start gap-3">
                                                <div className={`w-10 h-10 rounded-full flex items-center justify-center font-bold ${isSOS ? 'bg-red-500/20 text-red-400' : 'bg-secondary/20 text-secondary'}`}>
                                                    {isSOS ? <Siren className="w-5 h-5" /> : item.user_name.charAt(0)}
                                                </div>
                                                <div className="flex-1">
                                                    <p className="text-sm text-foreground">
                                                        <span className="font-semibold">{item.user_name}</span> {actionText}
                                                    </p>
                                                    <p className="text-xs text-muted-foreground mt-1">{new Date(item.created_at).toLocaleString()}</p>

                                                    <div className="flex items-center gap-4 mt-3">
                                                        {isSOS ? (
                                                            <button className="flex items-center gap-1 text-xs bg-gradient-to-r from-amber-500 to-orange-500 text-white px-3 py-1.5 rounded-full hover:opacity-90 transition-opacity">
                                                                <Sparkles className="w-4 h-4" />
                                                                Send Hype
                                                            </button>
                                                        ) : (
                                                            <>
                                                                <button className="flex items-center gap-1 text-xs text-muted-foreground hover:text-primary transition-colors">
                                                                    <Heart className="w-4 h-4" />
                                                                    {item.likes}
                                                                </button>
                                                                <button className="flex items-center gap-1 text-xs text-muted-foreground hover:text-primary transition-colors">
                                                                    <MessageSquare className="w-4 h-4" />
                                                                    {item.comments}
                                                                </button>
                                                            </>
                                                        )}
                                                    </div>
                                                </div>
                                            </div>
                                        </motion.div>
                                    );
                                })
                            ) : (
                                <div className="text-center text-muted-foreground">No activity yet. Complete your first fast to see events here!</div>
                            )}
                        </div>
                    </div>

                    {/* Sidebar - Right Side */}
                    <div className="space-y-8">
                        {/* My Tribes */}
                        {myTribesData?.tribes && myTribesData.tribes.length > 0 && (
                            <div className="mb-6">
                                <div className="flex items-center justify-between mb-4">
                                    <h2 className="font-semibold text-foreground flex items-center gap-2">
                                        <Sparkles className="w-5 h-5 text-amber-400" />
                                        My Tribes
                                    </h2>
                                </div>
                                <div className="space-y-3">
                                    {myTribesData.tribes.map((tribe) => (
                                        <div
                                            key={tribe.id}
                                            className="bg-gradient-to-r from-primary/10 to-secondary/10 border border-primary/30 rounded-xl p-3 hover:border-primary/50 transition-colors cursor-pointer"
                                            onClick={() => navigate(`/tribes/${tribe.id}`)}
                                        >
                                            <h3 className="font-medium text-foreground">{tribe.name}</h3>
                                            <p className="text-xs text-muted-foreground mb-2">{tribe.description}</p>
                                            <div className="flex items-center gap-1 text-xs text-primary">
                                                <Users className="w-3 h-3" />
                                                {tribe.member_count} members
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        {/* Featured Tribes */}
                        <div>
                            <div className="flex items-center justify-between mb-4">
                                <h2 className="font-semibold text-foreground flex items-center gap-2">
                                    <Users className="w-5 h-5 text-primary" />
                                    Featured Tribes
                                </h2>
                            </div>
                            <div className="space-y-3">
                                {isLoadingTribes ? (
                                    <div className="text-center text-muted-foreground">Loading tribes...</div>
                                ) : tribes && tribes.length > 0 ? (
                                    <>
                                        {tribes.slice(0, 3).map((tribe) => (
                                            <div
                                                key={tribe.id}
                                                className="bg-card border border-border rounded-xl p-3 hover:border-primary/50 transition-colors cursor-pointer"
                                                onClick={() => navigate('/tribes')}
                                            >
                                                <h3 className="font-medium text-foreground">{tribe.name}</h3>
                                                <p className="text-xs text-muted-foreground mb-2">{tribe.description}</p>
                                                <div className="flex items-center gap-1 text-xs text-secondary">
                                                    <Users className="w-3 h-3" />
                                                    {tribe.member_count} members
                                                </div>
                                            </div>
                                        ))}
                                        <Button
                                            variant="outline"
                                            className="w-full text-sm"
                                            onClick={() => navigate('/tribes')}
                                        >
                                            View All Tribes →
                                        </Button>
                                    </>
                                ) : (
                                    <div className="text-center p-4 border border-dashed rounded-xl">
                                        <p className="text-sm text-muted-foreground mb-3">No tribes yet.</p>
                                        <Button
                                            size="sm"
                                            onClick={() => navigate('/tribes')}
                                        >
                                            Create the First One!
                                        </Button>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Leaderboard */}
                        <div>
                            <h2 className="font-semibold text-foreground flex items-center gap-2 mb-4">
                                <Trophy className="w-5 h-5 text-yellow-500" />
                                Top Fasters
                            </h2>
                            <div className="bg-card border border-border rounded-xl overflow-hidden">
                                {isLoadingLeaderboard ? (
                                    <div className="p-4 text-center text-muted-foreground">Loading leaderboard...</div>
                                ) : (
                                    leaderboard?.map((entry, index) => (
                                        <div
                                            key={entry.user_id}
                                            className="flex items-center justify-between p-3 border-b border-border last:border-0"
                                        >
                                            <div className="flex items-center gap-3">
                                                <span className={`w-6 h-6 flex items-center justify-center rounded-full text-xs font-bold ${index === 0 ? 'bg-yellow-500/20 text-yellow-500' :
                                                    index === 1 ? 'bg-gray-400/20 text-gray-400' :
                                                        index === 2 ? 'bg-orange-700/20 text-orange-700' :
                                                            'bg-muted text-muted-foreground'
                                                    }`}>
                                                    {index + 1}
                                                </span>
                                                <span className="text-sm font-medium text-foreground">{entry.user_name || 'Anonymous'}</span>
                                            </div>
                                            <div className="text-xs text-muted-foreground">
                                                {entry.discipline_score.toFixed(0)} DS
                                            </div>
                                        </div>
                                    ))
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default Community;
