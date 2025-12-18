import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api } from '@/api/client';
import { useToast } from '@/hooks/use-toast';

export interface CravingHelpData {
    immediate_action: string;
    distraction_idea: string;
    biological_fact: string;
    motivation: string;
    time_remaining?: string;
    support_strategies: string[];
}

export interface SOSFlareResponse {
    sos_id: string;
    ai_response: CravingHelpData | null;
    tribe_notified: boolean;
    allies_notified: number;
    status: 'active' | 'cooldown';
    hype_count: number;
    error?: string;
}

export interface HypeResponse {
    id: string;
    sos_id: string;
    from_user_id: string;
    from_name?: string;
    message?: string;
    emoji: string;
    created_at: string;
}

export const useSOSFlare = () => {
    const queryClient = useQueryClient();
    const { toast } = useToast();

    const sendSOSFlare = useMutation({
        mutationFn: async (cravingDescription: string): Promise<SOSFlareResponse> => {
            const response = await api.post<SOSFlareResponse>('/fasting/sos', {
                craving_description: cravingDescription
            });
            return response.data;
        },
        onSuccess: (data) => {
            if (data.tribe_notified && data.allies_notified > 0) {
                toast({
                    title: "🆘 SOS Sent!",
                    description: `Your tribe (${data.allies_notified} allies) has been notified. Help is on the way!`,
                });
            } else {
                toast({
                    title: "💪 Cortex has your back!",
                    description: "You've got this. Follow the steps to beat this craving.",
                });
            }
            // Invalidate any SOS-related queries
            queryClient.invalidateQueries({ queryKey: ['active-sos'] });
        },
        onError: (error: any) => {
            const errorMessage = error.response?.data?.error || 'Failed to send SOS';
            const isCooldown = error.response?.status === 429;

            toast({
                variant: isCooldown ? "default" : "destructive",
                title: isCooldown ? "⏰ Cooldown Active" : "Error",
                description: isCooldown
                    ? "You can only send one SOS every 24 hours. Take a deep breath!"
                    : errorMessage,
            });
        },
    });

    return {
        sendSOSFlare: sendSOSFlare.mutate,
        sosData: sendSOSFlare.data,
        isSending: sendSOSFlare.isPending,
        error: sendSOSFlare.error,
        reset: sendSOSFlare.reset,
    };
};

// Hook to send hype to a struggling user
export const useSendHype = () => {
    const { toast } = useToast();

    return useMutation({
        mutationFn: async ({ sosId, emoji, message }: { sosId: string; emoji: string; message?: string }) => {
            const response = await api.post(`/sos/${sosId}/hype`, {
                emoji,
                message,
            });
            return response.data;
        },
        onSuccess: () => {
            toast({
                title: "🔥 Hype Sent!",
                description: "Your support has been delivered. You're a great ally!",
            });
        },
        onError: () => {
            toast({
                variant: "destructive",
                title: "Error",
                description: "Failed to send hype. Try again!",
            });
        },
    });
};

// Hook to resolve an SOS (mark as survived or failed)
export const useResolveSOS = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async ({ sosId, survived }: { sosId: string; survived: boolean }) => {
            const response = await api.post(`/sos/${sosId}/resolve`, {
                survived,
            });
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['active-sos'] });
        },
    });
};

// Hook to get hype responses for an SOS
export const useHypeResponses = (sosId: string | undefined) => {
    return useQuery({
        queryKey: ['hype-responses', sosId],
        queryFn: async () => {
            const response = await api.get<{ hypes: HypeResponse[] }>(`/sos/${sosId}/hypes`);
            return response.data.hypes;
        },
        enabled: !!sosId,
        refetchInterval: 5000, // Poll every 5 seconds for live updates
    });
};
