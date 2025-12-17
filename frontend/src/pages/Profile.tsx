import { useState } from 'react';
import { motion } from 'framer-motion';
import { ArrowLeft, User, Mail, Calendar, Award, Flame, Lock, Edit2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useUser } from '@/hooks/use-user';
import { toast } from 'sonner';

const Profile = () => {
    const navigate = useNavigate();
    const { user } = useUser();
    const [isEditing, setIsEditing] = useState(false);
    const [editData, setEditData] = useState({
        name: user?.name || '',
        email: user?.email || '',
    });

    const handleSave = async () => {
        try {
            // TODO: Call API to update profile
            toast.success('Profile updated successfully!');
            setIsEditing(false);
        } catch (error) {
            toast.error('Failed to update profile');
        }
    };

    const stats = [
        { label: 'Fasts Completed', value: '0', icon: Flame, color: 'text-orange-500' },
        { label: 'Discipline Score', value: '-- / 100', icon: Award, color: 'text-purple-500' },
        { label: 'Current Streak', value: '0 days', icon: Flame, color: 'text-red-500' },
        { label: 'Member Since', value: user?.created_at ? new Date(user.created_at).toLocaleDateString() : '--', icon: Calendar, color: 'text-blue-500' },
    ];

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
                            <h1 className="font-bold text-lg text-foreground">Profile</h1>
                        </div>
                    </div>
                </div>
            </header>

            <main className="container mx-auto px-4 py-8 max-w-4xl">
                {/* Profile Card */}
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ duration: 0.3 }}
                >
                    <Card className="p-8 mb-6">
                        <div className="flex items-start justify-between mb-6">
                            <div className="flex items-center gap-4">
                                <div className="w-20 h-20 rounded-full bg-primary/20 flex items-center justify-center">
                                    <User className="w-10 h-10 text-primary" />
                                </div>
                                <div>
                                    <h2 className="text-2xl font-bold text-foreground">{user?.name || 'User'}</h2>
                                    <p className="text-muted-foreground">{user?.email || 'user@example.com'}</p>
                                </div>
                            </div>
                            <Button
                                variant={isEditing ? 'default' : 'outline'}
                                onClick={() => isEditing ? handleSave() : setIsEditing(true)}
                            >
                                <Edit2 className="w-4 h-4 mr-2" />
                                {isEditing ? 'Save' : 'Edit Profile'}
                            </Button>
                        </div>

                        {isEditing && (
                            <div className="space-y-4 border-t border-border pt-6">
                                <div className="space-y-2">
                                    <Label htmlFor="name">Name</Label>
                                    <Input
                                        id="name"
                                        value={editData.name}
                                        onChange={(e) => setEditData({ ...editData, name: e.target.value })}
                                    />
                                </div>
                                <div className="space-y-2">
                                    <Label htmlFor="email">Email</Label>
                                    <Input
                                        id="email"
                                        type="email"
                                        value={editData.email}
                                        onChange={(e) => setEditData({ ...editData, email: e.target.value })}
                                    />
                                </div>
                            </div>
                        )}
                    </Card>

                    {/* Stats Grid */}
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
                        {stats.map((stat, index) => (
                            <motion.div
                                key={stat.label}
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                transition={{ delay: 0.1 * index }}
                            >
                                <Card className="p-6">
                                    <stat.icon className={`w-8 h-8 ${stat.color} mb-3`} />
                                    <div className="text-2xl font-bold text-foreground mb-1">{stat.value}</div>
                                    <div className="text-sm text-muted-foreground">{stat.label}</div>
                                </Card>
                            </motion.div>
                        ))}
                    </div>

                    {/* Account Details */}
                    <Card className="p-6">
                        <h3 className="font-semibold text-lg mb-4 flex items-center gap-2">
                            <Lock className="w-5 h-5 text-primary" />
                            Account Details
                        </h3>
                        <div className="space-y-3">
                            <div className="flex justify-between items-center py-2 border-b border-border">
                                <span className="text-muted-foreground">Subscription</span>
                                <span className="font-medium">Free Plan</span>
                            </div>
                            {/* HIDDEN FOR V2
                            <div className="flex justify-between items-center py-2 border-b border-border">
                                <span className="text-muted-foreground">Vault Deposit</span>
                                <span className="font-medium">$20.00</span>
                            </div>
                            <div className="flex justify-between items-center py-2">
                                <span className="text-muted-foreground">Vault Balance</span>
                                <span className="font-medium text-green-600">$0.00</span>
                            </div>
                            */}
                        </div>
                    </Card>
                </motion.div>
            </main>
        </div>
    );
};

export default Profile;
