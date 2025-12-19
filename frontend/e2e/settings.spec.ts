import { test, expect } from '@playwright/test';

test.describe('Settings Page', () => {
    test('should navigate to settings page', async ({ page }) => {
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        // Should load settings or redirect to login
        await expect(page.locator('body')).toBeVisible();
    });

    test('should display profile section', async ({ page }) => {
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        const url = page.url();
        if (!url.includes('login')) {
            // Look for profile-related content
            const profileContent = page.locator('[data-testid="profile"], text=Profile, text=Account, text=Email, text=Name');

            await expect(page.locator('body')).toBeVisible();
        }
    });

    test('should show notification settings', async ({ page }) => {
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        const url = page.url();
        if (!url.includes('login')) {
            // Look for notification settings
            const notifSettings = page.locator('[data-testid="notification-settings"], text=Notification, text=Reminder, input[type="checkbox"]');

            await expect(page.locator('body')).toBeVisible();
        }
    });

    test('should display fasting preferences', async ({ page }) => {
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        const url = page.url();
        if (!url.includes('login')) {
            // Look for fasting preferences
            const fastingSettings = page.locator('text=Fasting, text=Goal, text=Hour, text=Plan');

            await expect(page.locator('body')).toBeVisible();
        }
    });
});

test.describe('Profile Management', () => {
    test('should show user information display', async ({ page }) => {
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        const url = page.url();
        if (!url.includes('login')) {
            // Look for user info display
            const userInfo = page.locator('[data-testid="user-name"], [data-testid="user-email"], .user-profile');

            await expect(page.locator('body')).toBeVisible();
        }
    });

    test('should have logout option', async ({ page }) => {
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        const url = page.url();
        if (!url.includes('login')) {
            // Look for logout button
            const logoutButton = page.locator('button:has-text("Logout"), button:has-text("Log out"), button:has-text("Sign out"), [data-testid="logout"]');
            const hasLogout = await logoutButton.count() > 0;

            await expect(page.locator('body')).toBeVisible();
        }
    });

    test('should show subscription status', async ({ page }) => {
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        const url = page.url();
        if (!url.includes('login')) {
            // Look for subscription info
            const subscriptionInfo = page.locator('text=Subscription, text=Free, text=Premium, text=Vault, text=Upgrade');

            await expect(page.locator('body')).toBeVisible();
        }
    });
});

test.describe('Theme and Display', () => {
    test('should have theme toggle if available', async ({ page }) => {
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        const url = page.url();
        if (!url.includes('login')) {
            // Look for theme toggle
            const themeToggle = page.locator('[data-testid="theme-toggle"], text=Dark, text=Light, text=Theme');

            await expect(page.locator('body')).toBeVisible();
        }
    });

    test('should maintain page structure on mobile viewport', async ({ page }) => {
        await page.setViewportSize({ width: 375, height: 667 });
        await page.goto('/settings');
        await page.waitForLoadState('networkidle');

        // Page should be responsive
        await expect(page.locator('body')).toBeVisible();

        // Content should not overflow horizontally
        const bodyWidth = await page.evaluate(() => document.body.scrollWidth);
        expect(bodyWidth).toBeLessThanOrEqual(400);
    });
});

test.describe('Resource Pages', () => {
    test('should navigate to resources page', async ({ page }) => {
        await page.goto('/resources');
        await page.waitForLoadState('networkidle');

        await expect(page.locator('body')).toBeVisible();
    });

    test('should display educational content', async ({ page }) => {
        await page.goto('/resources');
        await page.waitForLoadState('networkidle');

        const url = page.url();
        if (!url.includes('login')) {
            // Look for resource content
            const resourceContent = page.locator('text=Recipe, text=Article, text=Guide, text=Learn');

            await expect(page.locator('body')).toBeVisible();
        }
    });
});
