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

// Add response interceptor for 401s
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && error.response.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export const authApi = {
    register: (email: string, password: string, referralCode?: string) => api.post<User>('/auth/register', { email, password, referral_code: referralCode }),
    login: (email: string, password: string) => api.post<{ token: string; refresh_token: string }>('/auth/login', { email, password }),
    getProfile: () => api.get<User>('/user/profile'),
};

export const fastingApi = {
    start: (plan_type: string, goal_hours: number, startTime?: string) => api.post<FastingSession>('/fasting/start', { plan_type, goal_hours, start_time: startTime }),
    stop: () => api.post<FastingSession>('/fasting/stop'),
    getCurrent: () => api.get<FastingSession>('/fasting/current'),
};

export const ketoApi = {
    log: (entry: Partial<KetoEntry>) => api.post('/keto/log', entry),
};

export const socialApi = {
    getFeed: () => api.get<SocialPost[]>('/social/feed'),
};

export const telemetryApi = {
    connect: (source: string) => api.post('/telemetry/connect', { source }),
    sync: (source: string) => api.post('/telemetry/sync', { source }),
    getStatus: () => api.get('/telemetry/status'),
    logManual: (type: string, value: number, unit: string) => api.post('/telemetry/manual', { type, value, unit }),
    getMetric: (type: string) => api.get<{ value: number, unit: string, timestamp: string }>('/telemetry/metric', { params: { type } }),
    getWeeklyStats: (type: string) => api.get<{ date: string, day: string, value: number }[]>('/telemetry/weekly', { params: { type } }),
};

export const mealApi = {
    log: (image: string, description: string) => api.post('/meals/', { image, description }),
    list: () => api.get<{ id: string, image: string, description: string, logged_at: string, analysis?: string, is_keto?: boolean, is_authentic?: boolean }[]>('/meals/'),
};

export const cortexApi = {
    chat: (message: string) => api.post<{ response: string }>('/cortex/chat', { message }),
    getInsight: (fastingHours: number) => api.post<{ insight: string }>('/cortex/insight', { fasting_hours: fastingHours }),
};

export interface Recipe {
    id: string;
    title: string;
    description: string;
    ingredients: string[];
    instructions: string[];
    diet: 'vegan' | 'vegetarian' | 'normal';
    is_simple: boolean;
    calories: number;
    carbs: number;
    image: string;
}

export const recipeApi = {
    list: (diet?: string) => api.get<Recipe[]>('/recipes/', { params: { diet } }),
};

export default api;
