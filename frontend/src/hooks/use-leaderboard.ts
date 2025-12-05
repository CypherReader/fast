import { useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';
import { LeaderboardEntry } from '@/api/types';

export const useLeaderboard = () => {
    const { data: leaderboard, isLoading } = useQuery({
        queryKey: ['leaderboard'],
        queryFn: async () => {
            const response = await api.get<LeaderboardEntry[]>('/leaderboard');
            return response.data;
        },
    });

    return {
        leaderboard,
        isLoading,
    };
};
