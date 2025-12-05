import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import { useToast } from '@/hooks/use-toast';

export interface Meal {
    id: string;
    user_id: string;
    name: string;
    meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
    calories: number;
    image?: string;
    description?: string;
    logged_at: string;
    analysis?: string;
    is_keto?: boolean;
    is_authentic?: boolean;
}

interface LogMealParams {
    name: string;
    calories: number;
    meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
    image?: string;
    description?: string;
}

export function useMeals() {
    const queryClient = useQueryClient();
    const { toast } = useToast();

    const { data: meals, isLoading } = useQuery<Meal[]>({
        queryKey: ['meals'],
        queryFn: async () => {
            const token = localStorage.getItem('token');
            const response = await axios.get('/api/v1/meals/', {
                headers: { Authorization: `Bearer ${token}` },
            });
            return response.data;
        },
    });

    const logMealMutation = useMutation({
        mutationFn: async (params: LogMealParams) => {
            const token = localStorage.getItem('token');
            const response = await axios.post('/api/v1/meals/', params, {
                headers: { Authorization: `Bearer ${token}` },
            });
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['meals'] });
            toast({
                title: 'Meal Logged',
                description: 'Your meal has been successfully logged.',
            });
        },
        onError: (error: unknown) => {
            const axiosError = error as { response?: { data?: { error?: string } } };
            toast({
                title: 'Error',
                description: axiosError.response?.data?.error || 'Failed to log meal',
                variant: 'destructive',
            });
        },
    });

    return {
        meals: meals || [],
        isLoading,
        logMeal: logMealMutation.mutate,
        isLogging: logMealMutation.isPending,
    };
}
