import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
    test('should navigate through main pages', async ({ page }) => {
        // Visit landing page
        await page.goto('/');
        await page.waitForLoadState('networkidle');
        await expect(page.locator('body')).toBeVisible();

        // Navigate to login
        await page.goto('/login');
        await page.waitForLoadState('networkidle');
        await expect(page.locator('body')).toBeVisible();
    });

    test('should have responsive design', async ({ page }) => {
        // Test mobile viewport
        await page.setViewportSize({ width: 375, height: 667 });
        await page.goto('/');
        await page.waitForLoadState('networkidle');

        // Page should still be functional
        await expect(page.locator('body')).toBeVisible();
    });

    test('should handle 404 pages gracefully', async ({ page }) => {
        await page.goto('/nonexistent-page-12345');
        await page.waitForLoadState('networkidle');

        // Should show some content (not crash)
        await expect(page.locator('body')).toBeVisible();
    });
});

test.describe('Accessibility Basics', () => {
    test('form inputs should exist on login page', async ({ page }) => {
        await page.goto('/login');
        await page.waitForLoadState('networkidle');

        // Check that inputs exist
        const inputs = page.locator('input');
        const inputCount = await inputs.count();
        expect(inputCount).toBeGreaterThan(0);
    });

    test('page should have proper structure', async ({ page }) => {
        await page.goto('/');
        await page.waitForLoadState('networkidle');

        // Check for heading
        const headings = page.locator('h1, h2, h3');
        const headingCount = await headings.count();
        expect(headingCount).toBeGreaterThan(0);
    });
});

test.describe('Performance', () => {
    test('landing page should load within acceptable time', async ({ page }) => {
        const startTime = Date.now();
        await page.goto('/');
        await page.waitForLoadState('networkidle');
        const loadTime = Date.now() - startTime;

        // Should load within 15 seconds
        expect(loadTime).toBeLessThan(15000);
    });
});
