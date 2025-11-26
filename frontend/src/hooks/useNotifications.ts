import { useEffect, useState } from 'react';
import { requestNotificationPermission, onMessageListener } from '../firebase';
import { api } from '../api/client';
import { useToast } from '@/components/ui/use-toast';

export function useNotifications() {
    const [fcmToken, setFcmToken] = useState<string | null>(null);
    const [permissionGranted, setPermissionGranted] = useState(false);
    const { toast } = useToast();

    useEffect(() => {
        const initializeNotifications = async () => {
            try {
                const token = await requestNotificationPermission();
                if (token) {
                    setFcmToken(token);
                    setPermissionGranted(true);

                    // Register token with backend
                    await api.post('/notifications/register-token', {
                        token,
                        device_type: 'web'
                    });

                    console.log('FCM token registered with backend');
                }
            } catch (error) {
                console.error('Error initializing notifications:', error);
            }
        };

        // Only initialize if user is logged in
        const token = localStorage.getItem('token');
        if (token) {
            initializeNotifications();
        }

        // Listen for foreground messages
        onMessageListener()
            .then((payload: any) => {
                console.log('Foreground message received:', payload);

                toast({
                    title: payload.notification?.title || 'New Notification',
                    description: payload.notification?.body || '',
                });
            })
            .catch((err) => console.error('Failed to listen for messages:', err));

        // Cleanup: unregister token on unmount/logout
        return () => {
            if (fcmToken) {
                api.post('/notifications/unregister-token', { token: fcmToken })
                    .catch((error) => console.error('Error unregistering token:', error));
            }
        };
    }, []);

    return { fcmToken, permissionGranted };
}
