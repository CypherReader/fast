import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/api/client';
import { FriendNetwork, Tribe, SocialEvent } from '@/api/types';

export const useSocial = () => {
    const queryClient = useQueryClient();

    const { data: friends, isLoading: isLoadingFriends } = useQuery({
        queryKey: ['friends'],
        queryFn: async () => {
            const response = await api.get<FriendNetwork[]>('/social/friends');
            return response.data;
        },
    });

    const { data: tribes, isLoading: isLoadingTribes } = useQuery({
        queryKey: ['tribes'],
        queryFn: async () => {
            const response = await api.get<Tribe[]>('/tribes');
            return response.data;
        },
    });

    const { data: feed, isLoading: isLoadingFeed } = useQuery({
        queryKey: ['feed'],
        queryFn: async () => {
            const response = await api.get<SocialEvent[]>('/social/feed');
            return response.data;
        },
    });

    const addFriendMutation = useMutation({
        mutationFn: async (friendId: string) => {
            const response = await api.post('/social/friends/add', { friend_id: friendId });
            return response.data;
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['friends'] });
        },
    });

    const createTribeMutation = useMutation({
        mutationFn: async (data: { name: string; description: string; is_public: boolean }) => {
            const response = await api.post<Tribe>('/social/tribes', data);
            return response.data;
        },
        onSuccess: () => {
            // Invalidate tribes list if we had one
        },
    });

    return {
        friends,
        isLoadingFriends,
        tribes,
        isLoadingTribes,
        feed,
        isLoadingFeed,
        addFriend: addFriendMutation.mutate,
        createTribe: createTribeMutation.mutate,
    };
};
