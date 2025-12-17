import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/api/client';
import { useToast } from '@/hooks/use-toast';

export interface CravingHelpResponse {
    immediate_action: string;
    distraction_idea: string;
    biological_fact: string;
    motivation: string;
    time_remaining?: string;
    support_strategies: string[];
}

export const useCravingHelp = () => {
    const queryClient = useQueryClient();
    const { toast } = useToast();

    const getCravingHelp = useMutation({
        mutationFn: async (cravingDescription: string) => {
            const response = await api.post<CravingHelpResponse>('/cortex/craving-help', {
                craving_description: cravingDescription
            });
            return response.data;
        },
        onSuccess: () => {
            toast({
                title: "Help is here!",
                description: "Cortex has your back. You've got this!",
            });
        },
        onError: (error: Error) => {
            toast({
                variant: "destructive",
                title: "Error",
                description: error.message || "Failed to get craving help",
            });
        },
    });

    return {
        getCravingHelp: getCravingHelp.mutate,
        cravingHelpData: getCravingHelp.data,
        isGettingHelp: getCravingHelp.isPending,
    };
};
