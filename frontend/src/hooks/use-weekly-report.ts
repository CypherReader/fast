import { useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';

export interface WeeklyReport {
    user_id: string;
    week_start: string;
    week_end: string;
    fasts_completed: number;
    average_duration: number;
    total_fasting_hours: number;
    longest_fast: number;
    discipline_change: number;
    best_day: string;
    challenge_day: string;
    goal_achievement_date?: string;
    ai_insights: string;
    predictions: {
        next_week_fasts_estimate: number;
        discipline_trend: string;
        success_probability: number;
    };
    recommendations: string[];
}

export function useWeeklyReport() {
    return useQuery({
        queryKey: ['weekly-report'],
        queryFn: async (): Promise<WeeklyReport> => {
            const response = await api.get('/cortex/weekly-report');
            return response.data;
        },
        staleTime: 1000 * 60 * 30, // Cache for 30 minutes
        retry: 1,
    });
}
