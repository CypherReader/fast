import { useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';

export interface User {
    id: string;
    email: string;
    name?: string;
    is_premium: boolean;
    created_at: string;
    push_notifications_enabled?: boolean;
}

export interface UserStats {
    fasts_completed: number;
    total_fasting_hours: number;
    current_streak: number;
    longest_streak: number;
    vault_balance: number;
    vault_total: number;
}

export const useUser = () => {
    const { data: user, isLoading: isUserLoading } = useQuery({
        queryKey: ['user'],
        queryFn: async () => {
            const response = await api.get<User>('/user/profile');
            return response.data;
        },
    });

    const { data: stats, isLoading: isStatsLoading } = useQuery({
        queryKey: ['user-stats'],
        queryFn: async () => {
            // TODO: Implement /users/stats endpoint in backend or use mock for now if not available
            // For now, we'll try to fetch it, but fallback to mock if 404
            try {
                const response = await api.get<UserStats>('/user/stats');
                return response.data;
            } catch (error) {
                console.warn('Failed to fetch user stats, using defaults');
                return {
                    fasts_completed: 0,
                    total_fasting_hours: 0,
                    current_streak: 0,
                    longest_streak: 0,
                    vault_balance: 0,
                    vault_total: 20,
                };
            }
        },
    });

    return {
        user,
        stats,
        isLoading: isUserLoading || isStatsLoading,
    };
};
