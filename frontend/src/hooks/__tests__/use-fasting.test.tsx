import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useFasting } from '../use-fasting';
import type { ReactNode } from 'react';

// Mock the API client
vi.mock('@/api/client', () => ({
    api: {
        get: vi.fn(),
        post: vi.fn(),
    },
}));

// Mock the toast hook
vi.mock('@/hooks/use-toast', () => ({
    useToast: () => ({
        toast: vi.fn(),
    }),
}));

// Import the mocked api
import { api } from '@/api/client';

// Create wrapper with QueryClient
const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: {
                retry: false,
            },
        },
    });

    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client={queryClient} > {children} </QueryClientProvider>
    );
};

describe('useFasting', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('currentFast query', () => {
        it('returns current fast when active', async () => {
            const mockFast = {
                id: '123',
                start_time: '2024-01-01T12:00:00Z',
                status: 'active',
                goal_hours: 16,
                plan_type: '16:8',
            };

            vi.mocked(api.get).mockResolvedValue({ data: mockFast });

            const { result } = renderHook(() => useFasting(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isLoading).toBe(false);
            });

            expect(result.current.currentFast).toEqual(mockFast);
        });

        it('returns null when no active fast', async () => {
            const error = { response: { status: 404 } };
            vi.mocked(api.get).mockRejectedValue(error);

            const { result } = renderHook(() => useFasting(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isLoading).toBe(false);
            });

            expect(result.current.currentFast).toBeUndefined();
        });

        it('sets isLoading correctly', () => {
            vi.mocked(api.get).mockImplementation(() => new Promise(() => { }));

            const { result } = renderHook(() => useFasting(), { wrapper: createWrapper() });

            expect(result.current.isLoading).toBe(true);
        });
    });

    describe('startFast mutation', () => {
        it('calls correct API endpoint', async () => {
            vi.mocked(api.get).mockResolvedValue({ data: null });
            vi.mocked(api.post).mockResolvedValue({ data: { id: 'new-fast' } });

            const { result } = renderHook(() => useFasting(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isLoading).toBe(false);
            });

            result.current.startFast({ plan_type: '16:8', goal_hours: 16 });

            await waitFor(() => {
                expect(api.post).toHaveBeenCalledWith('/fasting/start', {
                    plan_type: '16:8',
                    goal_hours: 16,
                });
            });
        });
    });

    describe('stopFast mutation', () => {
        it('calls correct API endpoint', async () => {
            vi.mocked(api.get).mockResolvedValue({ data: { id: 'active-fast' } });
            vi.mocked(api.post).mockResolvedValue({ data: { duration_minutes: 960 } });

            const { result } = renderHook(() => useFasting(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isLoading).toBe(false);
            });

            result.current.stopFast();

            await waitFor(() => {
                expect(api.post).toHaveBeenCalledWith('/fasting/stop');
            });
        });
    });

    describe('hook return values', () => {
        it('returns all expected functions and values', async () => {
            vi.mocked(api.get).mockResolvedValue({ data: null });

            const { result } = renderHook(() => useFasting(), { wrapper: createWrapper() });

            expect(result.current).toHaveProperty('currentFast');
            expect(result.current).toHaveProperty('isLoading');
            expect(result.current).toHaveProperty('startFast');
            expect(result.current).toHaveProperty('stopFast');
            expect(result.current).toHaveProperty('getInsight');
            expect(result.current).toHaveProperty('isStarting');
            expect(result.current).toHaveProperty('isStopping');
        });
    });
});
