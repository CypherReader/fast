import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useProgress } from '../use-progress';
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

import { api } from '@/api/client';

const createWrapper = () => {
    const queryClient = new QueryClient({
        defaultOptions: {
            queries: {
                retry: false,
            },
        },
    });

    return ({ children }: { children: ReactNode }) => (
        <QueryClientProvider client= { queryClient } > { children } </QueryClientProvider>
  );
};

describe('useProgress', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('weight history query', () => {
        it('returns weight history data', async () => {
            const mockWeightHistory = [
                { date: '2024-01-01', weight: 180 },
                { date: '2024-01-02', weight: 179.5 },
            ];

            vi.mocked(api.get).mockImplementation((url: string) => {
                if (url.includes('weight')) {
                    return Promise.resolve({ data: mockWeightHistory });
                }
                return Promise.resolve({ data: [] });
            });

            const { result } = renderHook(() => useProgress(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isWeightLoading).toBe(false);
            });

            expect(result.current.weightHistory).toEqual(mockWeightHistory);
        });

        it('handles empty weight history', async () => {
            vi.mocked(api.get).mockResolvedValue({ data: [] });

            const { result } = renderHook(() => useProgress(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isWeightLoading).toBe(false);
            });

            expect(result.current.weightHistory).toEqual([]);
        });
    });

    describe('hydration query', () => {
        it('returns daily hydration data', async () => {
            const mockHydration = { glasses: 5, goal: 8 };

            vi.mocked(api.get).mockImplementation((url: string) => {
                if (url.includes('hydration')) {
                    return Promise.resolve({ data: mockHydration });
                }
                return Promise.resolve({ data: [] });
            });

            const { result } = renderHook(() => useProgress(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isHydrationLoading).toBe(false);
            });

            expect(result.current.dailyHydration).toEqual(mockHydration);
        });
    });

    describe('logWeight mutation', () => {
        it('calls correct API endpoint with weight data', async () => {
            vi.mocked(api.get).mockResolvedValue({ data: [] });
            vi.mocked(api.post).mockResolvedValue({ data: { success: true } });

            const { result } = renderHook(() => useProgress(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isWeightLoading).toBe(false);
            });

            result.current.logWeight({ weight: 175.5, unit: 'lbs' });

            await waitFor(() => {
                expect(api.post).toHaveBeenCalledWith('/progress/weight', {
                    weight: 175.5,
                    unit: 'lbs',
                });
            });
        });
    });

    describe('logHydration mutation', () => {
        it('calls correct API endpoint with hydration data', async () => {
            vi.mocked(api.get).mockResolvedValue({ data: [] });
            vi.mocked(api.post).mockResolvedValue({ data: { success: true } });

            const { result } = renderHook(() => useProgress(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isWeightLoading).toBe(false);
            });

            result.current.logHydration({ amount: 1, unit: 'glass' });

            await waitFor(() => {
                expect(api.post).toHaveBeenCalledWith('/progress/hydration', {
                    amount: 1,
                    unit: 'glass',
                });
            });
        });
    });

    describe('hook return values', () => {
        it('returns all expected functions and values', async () => {
            vi.mocked(api.get).mockResolvedValue({ data: [] });

            const { result } = renderHook(() => useProgress(), { wrapper: createWrapper() });

            expect(result.current).toHaveProperty('weightHistory');
            expect(result.current).toHaveProperty('dailyHydration');
            expect(result.current).toHaveProperty('logWeight');
            expect(result.current).toHaveProperty('logHydration');
            expect(result.current).toHaveProperty('isWeightLoading');
            expect(result.current).toHaveProperty('isHydrationLoading');
        });
    });
});
