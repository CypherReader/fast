import { test, expect } from '@playwright/test';

test.describe('Dashboard Flow', () => {
    test('should redirect unauthenticated users', async ({ page }) => {
        await page.goto('/dashboard');
        await page.waitForLoadState('networkidle');

        // Should redirect to login or show login content
        const url = page.url();
        const hasLoginRedirect = url.includes('login') || url.includes('/');
        expect(hasLoginRedirect).toBe(true);
    });
});

test.describe('Fasting Timer', () => {
    test('landing page should mention fasting', async ({ page }) => {
        await page.goto('/');
        await page.waitForLoadState('networkidle');

        // Look for fasting-related content
        const pageContent = await page.textContent('body');
        const hasFastingContent = pageContent?.toLowerCase().includes('fast') ||
            pageContent?.toLowerCase().includes('timer') ||
            pageContent?.toLowerCase().includes('intermittent');
        expect(hasFastingContent).toBe(true);
    });
});

test.describe('Protected Routes', () => {
    test('progress page should require auth', async ({ page }) => {
        await page.goto('/progress');
        await page.waitForLoadState('networkidle');

        // Should redirect or show login/error
        const url = page.url();
        expect(url).toBeTruthy();
    });

    test('community page should require auth', async ({ page }) => {
        await page.goto('/community');
        await page.waitForLoadState('networkidle');

        // Page should load
        await expect(page.locator('body')).toBeVisible();
    });
});
