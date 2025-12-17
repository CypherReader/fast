import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/api/client';
import { useToast } from '@/hooks/use-toast';
import type { Notification as NotificationData } from '@/api/types';

export const useNotifications = () => {
    const queryClient = useQueryClient();
    const { toast } = useToast();

    const { data: notifications, isLoading } = useQuery({
        queryKey: ['notifications'],
        queryFn: async () => {
            const response = await api.get<NotificationData[]>('/notifications/history');
            return response.data;
        },
    });

    const registerTokenMutation = useMutation({
        mutationFn: async (data: { token: string; device_type: string }) => {
            await api.post('/notifications/register-token', data);
        },
        onSuccess: () => {
            toast({
                title: 'Notifications Enabled',
                description: 'You will now receive push notifications.',
            });
        },
        onError: () => {
            toast({
                variant: 'destructive',
                title: 'Error',
                description: 'Failed to enable notifications.',
            });
        },
    });

    const unregisterTokenMutation = useMutation({
        mutationFn: async (token: string) => {
            await api.post('/notifications/unregister-token', { token });
        },
        onSuccess: () => {
            toast({
                title: 'Notifications Disabled',
                description: 'You will no longer receive push notifications.',
            });
        },
        onError: () => {
            toast({
                variant: 'destructive',
                title: 'Error',
                description: 'Failed to disable notifications.',
            });
        },
    });

    // Request browser notification permission
    const requestPermission = async (): Promise<NotificationPermission> => {
        if (!('Notification' in window)) {
            toast({
                variant: 'destructive',
                title: 'Not Supported',
                description: 'Your browser does not support notifications.',
            });
            return 'denied';
        }

        try {
            const permission = await Notification.requestPermission();
            if (permission === 'granted') {
                toast({
                    title: 'Permission Granted',
                    description: 'Notifications are now enabled.',
                });
            }
            return permission;
        } catch (error) {
            toast({
                variant: 'destructive',
                title: 'Error',
                description: 'Failed to request notification permission.',
            });
            return 'denied';
        }
    };

    return {
        notifications,
        isLoading,
        registerToken: registerTokenMutation.mutate,
        unregisterToken: unregisterTokenMutation.mutate,
        isRegistering: registerTokenMutation.isPending,
        isUnregistering: unregisterTokenMutation.isPending,
        requestPermission,
        permissionStatus: typeof window !== 'undefined' && 'Notification' in window ? Notification.permission : 'denied',
        refreshHistory: () => queryClient.invalidateQueries({ queryKey: ['notifications'] }),
    };
};
