import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { tribesAPI } from '../api/tribes';
import type { CreateTribeRequest } from '../api/types';
import { toast } from 'sonner';

// Hook to list tribes with filters
export function useTribes(filters?: {
    search?: string;
    fasting_schedule?: string;
    primary_goal?: string;
    privacy?: string;
    sort_by?: 'newest' | 'popular' | 'active' | 'members';
    limit?: number;
    offset?: number;
}) {
    return useQuery({
        queryKey: ['tribes', filters],
        queryFn: () => tribesAPI.listTribes(filters),
    });
}

// Hook to get a single tribe
export function use Tribe(id: string | undefined) {
    return useQuery({
        queryKey: ['tribe', id],
        queryFn: () => tribesAPI.getTribe(id!),
        enabled: !!id,
    });
}

// Hook to get my tribes
export function useMyTribes() {
    return useQuery({
        queryKey: ['my-tribes'],
        queryFn: () => tribesAPI.getMyTribes(),
    });
}

// Hook to get tribe members
export function useTribeMembers(id: string | undefined, params?: { limit?: number; offset?: number }) {
    return useQuery({
        queryKey: ['tribe-members', id, params],
        queryFn: () => tribesAPI.getMembers(id!, params),
        enabled: !!id,
    });
}

// Hook to get tribe stats
export function useTribeStats(id: string | undefined) {
    return useQuery({
        queryKey: ['tribe-stats', id],
        queryFn: () => tribesAPI.getStats(id!),
        enabled: !!id,
    });
}

// Hook to create a tribe
export function useCreateTribe() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (data: CreateTribeRequest) => tribesAPI.createTribe(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['tribes'] });
            queryClient.invalidateQueries({ queryKey: ['my-tribes'] });
            toast.success('Tribe created successfully!');
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to create tribe');
        },
    });
}

// Hook to join a tribe
export function useJoinTribe() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (id: string) => tribesAPI.joinTribe(id),
        onSuccess: (_, id) => {
            queryClient.invalidateQueries({ queryKey: ['tribe', id] });
            queryClient.invalidateQueries({ queryKey: ['tribes'] });
            queryClient.invalidateQueries({ queryKey: ['my-tribes'] });
            toast.success('Successfully joined tribe!');
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to join tribe');
        },
    });
}

// Hook to leave a tribe
export function useLeaveTribe() {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (id: string) => tribesAPI.leaveTribe(id),
        onSuccess: (_, id) => {
            queryClient.invalidateQueries({ queryKey: ['tribe', id] });
            queryClient.invalidateQueries({ queryKey: ['tribes'] });
            queryClient.invalidateQueries({ queryKey: ['my-tribes'] });
            toast.success('Successfully left tribe');
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to leave tribe');
        },
    });
}
