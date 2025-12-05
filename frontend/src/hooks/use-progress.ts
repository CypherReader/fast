import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/api/client';
import { WeightLog, HydrationLog } from '@/api/types';

export const useProgress = () => {
    const queryClient = useQueryClient();

    const { data: weightHistory, isLoading: isWeightLoading } = useQuery({
        queryKey: ['weight-history'],
        queryFn: async () => {
            const response = await api.get<WeightLog[]>('/progress/weight');
            return response.data;
        },
    });

    const { data: dailyHydration, isLoading: isHydrationLoading } = useQuery({
        queryKey: ['daily-hydration'],
        queryFn: async () => {
            const response = await api.get<HydrationLog>('/progress/hydration/daily');
            return response.data;
        },
    });

    const logWeightMutation = useMutation({
        mutationFn: async (data: { weight: number; unit: string }) => {
            const response = await api.post<WeightLog>('/progress/weight', data);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['weight-history'] });
        },
    });

    const logHydrationMutation = useMutation({
        mutationFn: async (data: { amount: number; unit: string }) => {
            const response = await api.post<HydrationLog>('/progress/hydration', data);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['daily-hydration'] });
        },
    });

    return {
        weightHistory,
        dailyHydration,
        isWeightLoading,
        isHydrationLoading,
        logWeight: logWeightMutation.mutate,
        logHydration: logHydrationMutation.mutate,
    };
};
