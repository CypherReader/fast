import { Bell, BellOff, Check } from 'lucide-react';
import { useNotifications } from '@/hooks/use-notifications';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useState } from 'react';

export const NotificationSettings = () => {
    const {
        permissionStatus,
        requestPermission,
        registerToken,
        unregisterToken,
        isRegistering,
        isUnregistering,
        notifications,
        isLoading,
    } = useNotifications();

    const [fcmToken, setFcmToken] = useState<string | null>(null);

    const handleEnableNotifications = async () => {
        // First request browser permission
        const permission = await requestPermission();

        if (permission === 'granted') {
            // In a real implementation, you would get the FCM token from Firebase
            // For now, we'll simulate with a generated token
            const simulatedToken = `fcm_token_${Date.now()}`;
            setFcmToken(simulatedToken);
            registerToken({ token: simulatedToken, device_type: 'web' });
        }
    };

    const handleDisableNotifications = () => {
        if (fcmToken) {
            unregisterToken(fcmToken);
            setFcmToken(null);
        }
    };

    const isEnabled = permissionStatus === 'granted' && fcmToken;

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Bell className="h-5 w-5" />
                    Push Notifications
                </CardTitle>
                <CardDescription>
                    Get notified about your fasting milestones and achievements
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                    <div className="space-y-0.5">
                        <div className="text-sm font-medium">
                            Status: {permissionStatus === 'granted' ? 'Enabled' : permissionStatus === 'denied' ? 'Denied' : 'Not Enabled'}
                        </div>
                        <div className="text-xs text-muted-foreground">
                            {isEnabled ? 'You will receive push notifications' : 'Enable to receive updates'}
                        </div>
                    </div>

                    {isEnabled ? (
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={handleDisableNotifications}
                            disabled={isUnregistering}
                        >
                            <BellOff className="h-4 w-4 mr-2" />
                            Disable
                        </Button>
                    ) : (
                        <Button
                            size="sm"
                            onClick={handleEnableNotifications}
                            disabled={isRegistering || permissionStatus === 'denied'}
                        >
                            <Bell className="h-4 w-4 mr-2" />
                            Enable
                        </Button>
                    )}
                </div>

                {permissionStatus === 'denied' && (
                    <div className="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
                        Notifications are blocked. Please enable them in your browser settings.
                    </div>
                )}

                {/* Recent Notifications */}
                {notifications && notifications.length > 0 && (
                    <div className="space-y-2">
                        <h4 className="text-sm font-medium">Recent Notifications</h4>
                        <div className="space-y-2 max-h-48 overflow-y-auto">
                            {notifications.slice(0, 5).map((notification) => (
                                <div
                                    key={notification.id}
                                    className="flex items-start gap-2 rounded-lg border p-3 text-sm"
                                >
                                    {notification.read ? (
                                        <Check className="h-4 w-4 mt-0.5 text-muted-foreground" />
                                    ) : (
                                        <div className="h-2 w-2 rounded-full bg-primary mt-1.5" />
                                    )}
                                    <div className="flex-1 space-y-1">
                                        <p className="font-medium">{notification.title}</p>
                                        <p className="text-muted-foreground text-xs">{notification.body}</p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}

                {isLoading && (
                    <div className="text-center text-sm text-muted-foreground">
                        Loading notifications...
                    </div>
                )}
            </CardContent>
        </Card>
    );
};
