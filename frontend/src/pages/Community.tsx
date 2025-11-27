import { motion } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Users, MessageSquare, Heart, Trophy, Flame, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useUser } from '@/hooks/use-user';

// Mock Data
const tribes = [
    { id: 1, name: 'OMAD Warriors', members: 1240, description: 'One Meal A Day support group' },
    { id: 2, name: '72h Fast Club', members: 850, description: 'Extended fasting challenges' },
    { id: 3, name: 'Keto Fasters', members: 2100, description: 'Combining Keto with IF' },
];

const feed = [
    { id: 1, user: 'Sarah J.', action: 'completed a 18h fast', time: '2h ago', likes: 12, comments: 3 },
    { id: 2, user: 'Mike T.', action: 'joined the OMAD Warriors', time: '4h ago', likes: 24, comments: 5 },
    { id: 3, user: 'Jessica L.', action: 'reached a 7-day streak!', time: '5h ago', likes: 45, comments: 8 },
];

const leaderboard = [
    { rank: 1, user: 'David K.', score: 98, streak: 45 },
    { rank: 2, user: 'Anna M.', score: 96, streak: 32 },
    { rank: 3, user: 'Tom R.', score: 94, streak: 28 },
];

const Community = () => {
    const navigate = useNavigate();
    const { user } = useUser();

    return (
        <div className="min-h-screen bg-background">
            {/* Header */}
            <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-40">
                <div className="container mx-auto px-4 py-4 flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <Button variant="ghost" size="icon" onClick={() => navigate('/dashboard')}>
                            <ArrowLeft className="w-5 h-5" />
                        </Button>
                        <h1 className="font-bold text-lg text-foreground">Community</h1>
                    </div>
                    <div className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center">
                        <span className="text-sm font-medium text-primary">
                            {user?.name ? user.name.charAt(0).toUpperCase() : 'U'}
                        </span>
                    </div>
                </div>
            </header>

            <main className="container mx-auto px-4 py-6 max-w-4xl space-y-8">
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
                            {feed.map((item, index) => (
                                <motion.div
                                    key={item.id}
                                    className="bg-card border border-border rounded-xl p-4"
                                    initial={{ opacity: 0, y: 20 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    transition={{ delay: index * 0.1 }}
                                >
                                    <div className="flex items-start gap-3">
                                        <div className="w-10 h-10 rounded-full bg-secondary/20 flex items-center justify-center text-secondary font-bold">
                                            {item.user.charAt(0)}
                                        </div>
                                        <div className="flex-1">
                                            <p className="text-sm text-foreground">
                                                <span className="font-semibold">{item.user}</span> {item.action}
                                            </p>
                                            <p className="text-xs text-muted-foreground mt-1">{item.time}</p>

                                            <div className="flex items-center gap-4 mt-3">
                                                <button className="flex items-center gap-1 text-xs text-muted-foreground hover:text-primary transition-colors">
                                                    <Heart className="w-4 h-4" />
                                                    {item.likes}
                                                </button>
                                                <button className="flex items-center gap-1 text-xs text-muted-foreground hover:text-primary transition-colors">
                                                    <MessageSquare className="w-4 h-4" />
                                                    {item.comments}
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                </motion.div>
                            ))}
                        </div>
                    </div>

                    {/* Sidebar - Right Side */}
                    <div className="space-y-8">
                        {/* Featured Tribes */}
                        <div>
                            <h2 className="font-semibold text-foreground flex items-center gap-2 mb-4">
                                <Users className="w-5 h-5 text-primary" />
                                Featured Tribes
                            </h2>
                            <div className="space-y-3">
                                {tribes.map((tribe) => (
                                    <div key={tribe.id} className="bg-card border border-border rounded-xl p-3 hover:border-primary/50 transition-colors cursor-pointer">
                                        <h3 className="font-medium text-foreground">{tribe.name}</h3>
                                        <p className="text-xs text-muted-foreground mb-2">{tribe.description}</p>
                                        <div className="flex items-center gap-1 text-xs text-secondary">
                                            <Users className="w-3 h-3" />
                                            {tribe.members} members
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>

                        {/* Leaderboard */}
                        <div>
                            <h2 className="font-semibold text-foreground flex items-center gap-2 mb-4">
                                <Trophy className="w-5 h-5 text-yellow-500" />
                                Top Fasters
                            </h2>
                            <div className="bg-card border border-border rounded-xl overflow-hidden">
                                {leaderboard.map((entry, index) => (
                                    <div
                                        key={entry.rank}
                                        className="flex items-center justify-between p-3 border-b border-border last:border-0"
                                    >
                                        <div className="flex items-center gap-3">
                                            <span className={`w-6 h-6 flex items-center justify-center rounded-full text-xs font-bold ${index === 0 ? 'bg-yellow-500/20 text-yellow-500' :
                                                    index === 1 ? 'bg-gray-400/20 text-gray-400' :
                                                        'bg-orange-700/20 text-orange-700'
                                                }`}>
                                                {entry.rank}
                                            </span>
                                            <span className="text-sm font-medium text-foreground">{entry.user}</span>
                                        </div>
                                        <div className="text-xs text-muted-foreground">
                                            {entry.streak} day streak
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default Community;
