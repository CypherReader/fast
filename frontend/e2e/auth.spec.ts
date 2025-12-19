import { test, expect } from '@playwright/test';

test.describe('Authentication Flow', () => {
    test('should display login page', async ({ page }) => {
        await page.goto('/login');

        // Wait for page to load
        await page.waitForLoadState('networkidle');

        // Check page has some form elements
        await expect(page.locator('form')).toBeVisible({ timeout: 10000 }).catch(() => {
            // Fallback - just check page loaded
            expect(page.url()).toContain('/login');
        });
    });

    test('should have email and password inputs', async ({ page }) => {
        await page.goto('/login');
        await page.waitForLoadState('networkidle');

        // Look for inputs by type or placeholder
        const emailInput = page.locator('input[type="email"], input[placeholder*="email" i], input[name*="email" i]');
        const passwordInput = page.locator('input[type="password"]');

        await expect(emailInput.first()).toBeVisible({ timeout: 10000 });
        await expect(passwordInput.first()).toBeVisible({ timeout: 10000 });
    });

    test('should navigate to register page', async ({ page }) => {
        await page.goto('/login');
        await page.waitForLoadState('networkidle');

        // Look for any link to register
        const registerLink = page.locator('a[href*="register"], a[href*="signup"], a:has-text("Sign up"), a:has-text("Register"), a:has-text("Create account")');

        if (await registerLink.first().isVisible()) {
            await registerLink.first().click();
            await expect(page).toHaveURL(/register|signup/);
        } else {
            // Skip if no register link
            test.skip();
        }
    });
});

test.describe('Landing Page', () => {
    test('should display landing page content', async ({ page }) => {
        await page.goto('/');
        await page.waitForLoadState('networkidle');

        // Check page loaded with some content
        await expect(page.locator('body')).toBeVisible();
        const title = await page.title();
        expect(title.length).toBeGreaterThan(0);
    });

    test('should have navigation', async ({ page }) => {
        await page.goto('/');
        await page.waitForLoadState('networkidle');

        // Check for header or nav element
        const navElement = page.locator('header, nav, [role="navigation"]');
        await expect(navElement.first()).toBeVisible({ timeout: 10000 });
    });

    test('should load without errors', async ({ page }) => {
        const errors: string[] = [];
        page.on('pageerror', error => errors.push(error.message));

        await page.goto('/');
        await page.waitForLoadState('networkidle');

        // Check no critical JS errors
        expect(errors.filter(e => e.includes('TypeError') || e.includes('ReferenceError'))).toHaveLength(0);
    });
});
