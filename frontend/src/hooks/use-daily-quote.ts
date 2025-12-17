import { useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';

interface DailyQuoteResponse {
    quote: string;
}

export const useDailyQuote = () => {
    const { data, isLoading, error } = useQuery({
        queryKey: ['daily-quote'],
        queryFn: async () => {
            const response = await api.get<DailyQuoteResponse>('/cortex/daily-quote');
            return response.data;
        },
        staleTime: 24 * 60 * 60 * 1000, // Cache for 24 hours
        refetchOnMount: false,
        refetchOnWindowFocus: false,
    });

    return {
        quote: data?.quote || '',
        isLoading,
        error,
    };
};
