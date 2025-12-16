import { useQuery } from '@tanstack/react-query';
import { api } from '../api/client';

interface FastingInsight {
    hours: number;
    milestone: string;
    insight: string;
    benefits: string[];
    motivation: string;
}

export function useFastingInsight(hours: number | undefined) {
    return useQuery({
        queryKey: ['fasting-insight', hours],
        queryFn: async () => {
            const response = await api.get<FastingInsight>('/fasting/insight', {
                params: { hours },
            });
            return response.data;
        },
        enabled: !!hours && hours > 0,
        staleTime: 1000 * 60 * 60, // Cache for 1 hour
        cacheTime: 1000 * 60 * 60 * 2, // Keep in cache for 2 hours
    });
}
