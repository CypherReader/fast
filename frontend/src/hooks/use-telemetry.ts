import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/api/client';
import { useToast } from '@/hooks/use-toast';

export type MetricType = 'weight' | 'water' | 'steps' | 'glucose' | 'ketones';

export interface TelemetryEntry {
    id: string;
    type: MetricType;
    value: number;
    unit: string;
    logged_at: string;
    source: string;
}

export interface WeeklyStats {
    date: string;
    value: number;
}

export const useTelemetry = (type: MetricType) => {
    const queryClient = useQueryClient();
    const { toast } = useToast();

    const { data: latestMetric, isLoading: isLatestLoading } = useQuery({
        queryKey: ['telemetry', type, 'latest'],
        queryFn: async () => {
            try {
                const response = await api.get<TelemetryEntry>(`/telemetry/metric?type=${type}`);
                return response.data;
            } catch (error) {
                return null;
            }
        },
    });

    const { data: history, isLoading: isHistoryLoading } = useQuery({
        queryKey: ['telemetry', type, 'history'],
        queryFn: async () => {
            try {
                // Using weekly stats for the chart for now
                const response = await api.get<WeeklyStats[]>(`/telemetry/weekly?type=${type}`);
                return response.data;
            } catch (error) {
                return [];
            }
        },
    });

    const logMetricMutation = useMutation({
        mutationFn: async (params: { value: number; unit: string }) => {
            const response = await api.post('/telemetry/manual', {
                type,
                value: params.value,
                unit: params.unit,
            });
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['telemetry', type] });
            toast({
                title: "Logged!",
                description: `Your ${type} has been recorded.`,
            });
        },
        onError: () => {
            toast({
                variant: "destructive",
                title: "Error",
                description: "Failed to log data.",
            });
        },
    });

    return {
        latestMetric,
        history,
        isLoading: isLatestLoading || isHistoryLoading,
        logMetric: logMetricMutation.mutate,
        isLogging: logMetricMutation.isPending,
    };
};
