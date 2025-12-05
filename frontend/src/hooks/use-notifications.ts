import { useMutation, useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';
import { Notification } from '@/api/types';

export const useNotifications = () => {
    const { data: notifications, isLoading } = useQuery({
        queryKey: ['notifications'],
        queryFn: async () => {
            const response = await api.get<Notification[]>('/notifications/history');
            return response.data;
        },
    });

    const registerTokenMutation = useMutation({
        mutationFn: async (data: { token: string; device_type: string }) => {
            await api.post('/notifications/register-token', data);
        },
    });

    return {
        notifications,
        isLoading,
        registerToken: registerTokenMutation.mutate,
    };
};
