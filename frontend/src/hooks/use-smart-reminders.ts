import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient } from '@/api/client';
import { useToast } from '@/hooks/use-toast';

// Types
export interface ReminderSettings {
    user_id?: string;
    reminder_fast_start: boolean;
    reminder_fast_end: boolean;
    reminder_hydration: boolean;
    preferred_fast_start_hour: number; // 0-23
    hydration_interval_minutes: number;
}

export interface OptimalFastingWindow {
    suggested_start_time: string;
    suggested_end_time: string;
    suggested_duration_hours: number;
    reasoning: string;
    confidence_score: number;
}

// Hook to get and update reminder settings
export function useReminderSettings() {
    const queryClient = useQueryClient();
    const { toast } = useToast();

    const settingsQuery = useQuery({
        queryKey: ['reminder-settings'],
        queryFn: async (): Promise<ReminderSettings> => {
            const response = await apiClient.get('/user/reminder-settings');
            return response.data;
        },
    });

    const updateMutation = useMutation({
        mutationFn: async (settings: Partial<ReminderSettings>) => {
            const response = await apiClient.put('/user/reminder-settings', settings);
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['reminder-settings'] });
            toast({
                title: 'Settings saved!',
                description: 'Your reminder preferences have been updated.',
            });
        },
        onError: (error: Error) => {
            toast({
                title: 'Error',
                description: error.message || 'Failed to update settings',
                variant: 'destructive',
            });
        },
    });

    return {
        settings: settingsQuery.data,
        isLoading: settingsQuery.isLoading,
        isError: settingsQuery.isError,
        updateSettings: updateMutation.mutate,
        isUpdating: updateMutation.isPending,
    };
}

// Hook to get AI-suggested optimal fasting window
export function useOptimalFastingWindow() {
    return useQuery({
        queryKey: ['optimal-fasting-window'],
        queryFn: async (): Promise<OptimalFastingWindow> => {
            const response = await apiClient.get('/user/optimal-fasting-window');
            return response.data;
        },
        staleTime: 1000 * 60 * 60, // Cache for 1 hour
    });
}

// Helper to format hour to readable time
export function formatHour(hour: number): string {
    const period = hour >= 12 ? 'PM' : 'AM';
    const displayHour = hour % 12 || 12;
    return `${displayHour}:00 ${period}`;
}
