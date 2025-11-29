import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import { api } from '@/api/client';
import { useToast } from '@/hooks/use-toast';

export interface FastingSession {
    id: string;
    start_time: string;
    end_time?: string;
    status: 'active' | 'completed' | 'cancelled';
    goal_hours: number;
    plan_type: string;
    duration_minutes?: number;
}

export const useFasting = () => {
    const queryClient = useQueryClient();
    const { toast } = useToast();

    const { data: currentFast, isLoading } = useQuery({
        queryKey: ['current-fast'],
        queryFn: async () => {
            try {
                const response = await api.get<FastingSession>('/fasting/current');
                console.log('useFasting: /fasting/current response:', response.data);
                return response.data;
            } catch (error: unknown) {
                if (axios.isAxiosError(error) && error.response?.status === 404) return null;
                throw error;
            }
        },
    });

    const startFastMutation = useMutation({
        mutationFn: async (params: { plan_type: string; goal_hours: number; start_time?: string }) => {
            const response = await api.post('/fasting/start', params);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['current-fast'] });
            toast({
                title: "Fast Started",
                description: "Good luck! You've got this.",
            });
        },
        onError: (error: unknown) => {
            const errorMessage = axios.isAxiosError(error)
                ? error.response?.data?.error || "Failed to start fast."
                : "An unexpected error occurred.";
            toast({
                variant: "destructive",
                title: "Error",
                description: errorMessage,
            });
        },
    });

    const stopFastMutation = useMutation({
        mutationFn: async () => {
            const response = await api.post('/fasting/stop');
            return response.data;
        },
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ['current-fast'] });
            toast({
                title: "Fast Completed",
                description: `You fasted for ${data.duration_minutes} minutes!`,
            });
        },
        onError: (error: unknown) => {
            const errorMessage = axios.isAxiosError(error)
                ? error.response?.data?.error || "Failed to stop fast."
                : "An unexpected error occurred.";
            toast({
                variant: "destructive",
                title: "Error",
                description: errorMessage,
            });
        },
    });

    const getInsightMutation = useMutation({
        mutationFn: async (fastingHours: number) => {
            const response = await api.post('/cortex/insight', { fasting_hours: fastingHours });
            return response.data;
        },
    });

    return {
        currentFast,
        isLoading,
        startFast: startFastMutation.mutate,
        stopFast: stopFastMutation.mutate,
        getInsight: getInsightMutation.mutateAsync,
        isInsightLoading: getInsightMutation.isPending,
        isStarting: startFastMutation.isPending,
        isStopping: stopFastMutation.isPending,
    };
};
