import { api } from './client';
import type { Tribe, CreateTribeRequest, ListTribesResponse, TribeMember, TribeStats } from './types';

export const tribesAPI = {
    // List all tribes with optional filters
    listTribes: async (params?: {
        search?: string;
        fasting_schedule?: string;
        primary_goal?: string;
        privacy?: string;
        sort_by?: 'newest' | 'popular' | 'active' | 'members';
        limit?: number;
        offset?: number;
    }): Promise<ListTribesResponse> => {
        const response = await api.get<ListTribesResponse>('/tribes', { params });
        return response.data;
    },

    // Get a specific tribe by ID
    getTribe: async (id: string): Promise<Tribe> => {
        const response = await api.get<Tribe>(`/tribes/${id}`);
        return response.data;
    },

    // Create a new tribe
    createTribe: async (data: CreateTribeRequest): Promise<Tribe> => {
        const response = await api.post<Tribe>('/tribes', data);
        return response.data;
    },

    // Join a tribe
    joinTribe: async (id: string): Promise<{ message: string }> => {
        const response = await api.post(`/tribes/${id}/join`);
        return response.data;
    },

    // Leave a tribe
    leaveTribe: async (id: string): Promise<{ message: string }> => {
        const response = await api.post(`/tribes/${id}/leave`);
        return response.data;
    },

    // Get tribe members
    getMembers: async (id: string, params?: { limit?: number; offset?: number }): Promise<{ members: TribeMember[] }> => {
        const response = await api.get(`/tribes/${id}/members`, { params });
        return response.data;
    },

    // Get my tribes
    getMyTribes: async (): Promise<{ tribes: Tribe[] }> => {
        const response = await api.get('/users/me/tribes');
        return response.data;
    },

    // Get tribe stats
    getStats: async (id: string): Promise<TribeStats> => {
        const response = await api.get<TribeStats>(`/tribes/${id}/stats`);
        return response.data;
    },
};
