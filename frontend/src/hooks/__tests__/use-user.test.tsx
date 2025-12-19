import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useUser } from '../use-user';
import type { ReactNode } from 'react';

// Mock the API client
vi.mock('@/api/client', () => ({
    api: {
        get: vi.fn(),
    },
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
        <QueryClientProvider client={queryClient} > {children} </QueryClientProvider>
    );
};

describe('useUser', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    describe('user profile query', () => {
        it('returns user data when authenticated', async () => {
            const mockUser = {
                id: 'user-123',
                email: 'test@example.com',
                name: 'Test User',
                discipline_index: 75,
                vault_deposit: 20,
                earned_refund: 15,
            };

            vi.mocked(api.get).mockResolvedValue({ data: mockUser });

            const { result } = renderHook(() => useUser(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isLoading).toBe(false);
            });

            expect(result.current.user).toEqual(mockUser);
        });

        it('returns null when not authenticated', async () => {
            const error = { response: { status: 401 } };
            vi.mocked(api.get).mockRejectedValue(error);

            const { result } = renderHook(() => useUser(), { wrapper: createWrapper() });

            await waitFor(() => {
                expect(result.current.isLoading).toBe(false);
            });

            expect(result.current.user).toBeUndefined();
        });

        it('sets isLoading correctly during fetch', () => {
            vi.mocked(api.get).mockImplementation(() => new Promise(() => { }));

            const { result } = renderHook(() => useUser(), { wrapper: createWrapper() });

            expect(result.current.isLoading).toBe(true);
        });
    });

    describe('hook return values', () => {
        it('returns expected properties', async () => {
            vi.mocked(api.get).mockResolvedValue({ data: {} });

            const { result } = renderHook(() => useUser(), { wrapper: createWrapper() });

            expect(result.current).toHaveProperty('user');
            expect(result.current).toHaveProperty('isLoading');
        });
    });
});
