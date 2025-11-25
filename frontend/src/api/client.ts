import axios from 'axios';
import { User, FastingSession, KetoEntry, SocialPost } from './types';

const API_URL = 'http://localhost:8080/api/v1';

export const api = axios.create({
    baseURL: API_URL,
});

// Add token interceptor
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const authApi = {
    register: (email: string, password: string) => api.post<User>('/auth/register', { email, password }),
    login: (email: string, password: string) => api.post<{ token: string; refresh_token: string }>('/auth/login', { email, password }),
    getProfile: () => api.get<User>('/user/profile'),
};

export const fastingApi = {
    start: (plan_type: string, goal_hours: number) => api.post<FastingSession>('/fasting/start', { plan_type, goal_hours }),
    stop: () => api.post<FastingSession>('/fasting/stop'),
    getCurrent: () => api.get<FastingSession>('/fasting/current'),
};

export const ketoApi = {
    log: (entry: Partial<KetoEntry>) => api.post('/keto/log', entry),
};

export const socialApi = {
    getFeed: () => api.get<SocialPost[]>('/social/feed'),
};

export const cortexApi = {
    chat: (message: string) => api.post<{ response: string }>('/cortex/chat', { message }),
    getInsight: (fastingHours: number) => api.post<{ insight: string }>('/cortex/insight', { fasting_hours: fastingHours }),
};

export default api;
