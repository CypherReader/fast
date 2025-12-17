import { useQuery } from '@tanstack/react-query';
import { api } from '@/api/client';

export interface BreakFastGuide {
    fast_duration: number;
    meal_type: string;
    portion_size: string;
    recommended_foods: string[];
    foods_to_avoid: string[];
    reintroduction_plan: string;
    hydration_tip: string;
    timing_advice: string;
    ai_guidance: string;
}

export const useBreakFastPlanner = (fastDuration: number, enabled: boolean = false) => {
    const { data, isLoading, error } = useQuery({
        queryKey: ['breakfast-planner', fastDuration],
        queryFn: async () => {
            const response = await api.post<BreakFastGuide>('/cortex/break-fast-guide', {
                fast_duration: fastDuration,
            });
            return response.data;
        },
        enabled,
        staleTime: 5 * 60 * 1000, // Cache for 5 minutes
    });

    return {
        guide: data,
        isLoading,
        error,
    };
};
