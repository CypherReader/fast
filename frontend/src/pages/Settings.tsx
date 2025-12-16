import { useState } from 'react';
import { motion } from 'framer-motion';
import { ArrowLeft, Bell, Moon, Globe, Lock, LogOut, Trash2 } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { toast } from 'sonner';

const Settings = () => {
    const navigate = useNavigate();
    const [settings, setSettings] = useState({
        notifications: {
            fastReminders: true,
            milestoneAlerts: true,
            communityUpdates: false,
            emailDigest: true,
        },
        preferences: {
            darkMode: false,
            language: 'en',
        },
    });

    const handleToggle = (category: 'notifications' | 'preferences', key: string) => {
        setSettings(prev => ({
            ...prev,
            [category]: {
                ...prev[category],
                [key]: !prev[category][key as keyof typeof prev[typeof category]],
            },
        }));
        toast.success('Settings updated');
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
        toast.success('Logged out successfully');
    };

    const handleDeleteAccount = () => {
        if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
            // TODO: Call API to delete account
            toast.success('Account deletion initiated');
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
                            <h1 className="font-bold text-lg text-foreground">Settings</h1>
                        </div>
                    </div>
                </div>
            </header>

            <main className="container mx-auto px-4 py-8 max-w-4xl">
                <div className="space-y-6">
                    {/* Notifications */}
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.3 }}
                    >
                        <Card className="p-6">
                            <h3 className="font-semibold text-lg mb-4 flex items-center gap-2">
                                <Bell className="w-5 h-5 text-primary" />
                                Notifications
                            </h3>
                            <div className="space-y-4">
                                <div className="flex items-center justify-between py-3 border-b border-border">
                                    <div>
                                        <Label htmlFor="fast-reminders">Fast Reminders</Label>
                                        <p className="text-sm text-muted-foreground">Get reminders when it's time to start or end your fast</p>
                                    </div>
                                    <Switch
                                        id="fast-reminders"
                                        checked={settings.notifications.fastReminders}
                                        onCheckedChange={() => handleToggle('notifications', 'fastReminders')}
                                    />
                                </div>
                                <div className="flex items-center justify-between py-3 border-b border-border">
                                    <div>
                                        <Label htmlFor="milestone-alerts">Milestone Alerts</Label>
                                        <p className="text-sm text-muted-foreground">Celebrate when you reach fasting milestones</p>
                                    </div>
                                    <Switch
                                        id="milestone-alerts"
                                        checked={settings.notifications.milestoneAlerts}
                                        onCheckedChange={() => handleToggle('notifications', 'milestoneAlerts')}
                                    />
                                </div>
                                <div className="flex items-center justify-between py-3 border-b border-border">
                                    <div>
                                        <Label htmlFor="community-updates">Community Updates</Label>
                                        <p className="text-sm text-muted-foreground">Stay updated with tribe activities and challenges</p>
                                    </div>
                                    <Switch
                                        id="community-updates"
                                        checked={settings.notifications.communityUpdates}
                                        onCheckedChange={() => handleToggle('notifications', 'communityUpdates')}
                                    />
                                </div>
                                <div className="flex items-center justify-between py-3">
                                    <div>
                                        <Label htmlFor="email-digest">Weekly Email Digest</Label>
                                        <p className="text-sm text-muted-foreground">Receive a summary of your progress each week</p>
                                    </div>
                                    <Switch
                                        id="email-digest"
                                        checked={settings.notifications.emailDigest}
                                        onCheckedChange={() => handleToggle('notifications', 'emailDigest')}
                                    />
                                </div>
                            </div>
                        </Card>
                    </motion.div>

                    {/* Preferences */}
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.1, duration: 0.3 }}
                    >
                        <Card className="p-6">
                            <h3 className="font-semibold text-lg mb-4 flex items-center gap-2">
                                <Globe className="w-5 h-5 text-primary" />
                                Preferences
                            </h3>
                            <div className="space-y-4">
                                <div className="flex items-center justify-between py-3 border-b border-border">
                                    <div>
                                        <Label htmlFor="dark-mode">Dark Mode</Label>
                                        <p className="text-sm text-muted-foreground">Use dark theme across the app</p>
                                    </div>
                                    <Switch
                                        id="dark-mode"
                                        checked={settings.preferences.darkMode}
                                        onCheckedChange={() => handleToggle('preferences', 'darkMode')}
                                    />
                                </div>
                                <div className="flex items-center justify-between py-3">
                                    <div>
                                        <Label>Language</Label>
                                        <p className="text-sm text-muted-foreground">Choose your preferred language</p>
                                    </div>
                                    <span className="text-sm font-medium">English</span>
                                </div>
                            </div>
                        </Card>
                    </motion.div>

                    {/* Security */}
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.2, duration: 0.3 }}
                    >
                        <Card className="p-6">
                            <h3 className="font-semibold text-lg mb-4 flex items-center gap-2">
                                <Lock className="w-5 h-5 text-primary" />
                                Security
                            </h3>
                            <div className="space-y-3">
                                <Button variant="outline" className="w-full justify-start">
                                    Change Password
                                </Button>
                                <Button variant="outline" className="w-full justify-start" onClick={handleLogout}>
                                    <LogOut className="w-4 h-4 mr-2" />
                                    Logout
                                </Button>
                            </div>
                        </Card>
                    </motion.div>

                    {/* Danger Zone */}
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: 0.3, duration: 0.3 }}
                    >
                        <Card className="p-6 border-red-200 dark:border-red-900">
                            <h3 className="font-semibold text-lg mb-4 flex items-center gap-2 text-red-600">
                                <Trash2 className="w-5 h-5" />
                                Danger Zone
                            </h3>
                            <div className="space-y-3">
                                <p className="text-sm text-muted-foreground">
                                    Once you delete your account, there is no going back. All your data will be permanently removed.
                                </p>
                                <Button variant="destructive" onClick={handleDeleteAccount}>
                                    Delete Account
                                </Button>
                            </div>
                        </Card>
                    </motion.div>
                </div>
            </main>
        </div>
    );
};

export default Settings;
