import posthog from 'posthog-js';

// Initialize PostHog
export const initAnalytics = () => {
    if (import.meta.env.PROD) {
        posthog.init(import.meta.env.VITE_POSTHOG_KEY || 'phc_placeholder_key', {
            api_host: import.meta.env.VITE_POSTHOG_HOST || 'https://app.posthog.com',
            capture_pageview: false, // We'll manually track pageviews for better control in SPAs
        });
    }
};

// Track generic events
export const trackEvent = (eventName: string, properties?: Record<string, any>) => {
    if (import.meta.env.PROD) {
        posthog.capture(eventName, properties);
    } else {
        console.log(`[Analytics] ${eventName}`, properties);
    }
};

// Specific tracking helpers
export const analytics = {
    trackEvent: (eventName: string, properties?: Record<string, any>) => {
        trackEvent(eventName, properties);
    },

    identify: (userId: string, traits?: Record<string, any>) => {
        if (import.meta.env.PROD) {
            posthog.identify(userId, traits);
        } else {
            console.log(`[Analytics] Identify: ${userId}`, traits);
        }
    },

    pageView: (path: string) => {
        if (import.meta.env.PROD) {
            posthog.capture('$pageview', { path });
        } else {
            console.log(`[Analytics] PageView: ${path}`);
        }
    },

    trackPaywallView: (source: string) => {
        trackEvent('paywall_viewed', { source });
    },

    trackSubscriptionStarted: (plan: string) => {
        trackEvent('subscription_started', { plan });
    },

    trackFastStarted: (type: string, duration: number) => {
        trackEvent('fast_started', { type, duration });
    },

    trackFastCompleted: (duration: number, success: boolean) => {
        trackEvent('fast_completed', { duration, success });
    }
};
