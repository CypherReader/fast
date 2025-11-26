import posthog from 'posthog-js';
import { useState, useEffect } from 'react';

export const FEATURE_FLAGS = {
    PAYWALL_AGGRESSIVE_MODE: 'paywall_aggressive_mode',
    SHOW_REFERRAL_TEASER: 'show_referral_teaser',
    ONBOARDING_QUIZ_FLOW: 'onboarding_quiz_flow',
};

// Hook to get feature flag value
export const useFeatureFlag = (flagKey: string, defaultValue: boolean = false) => {
    const [isEnabled, setIsEnabled] = useState(defaultValue);

    useEffect(() => {
        if (import.meta.env.PROD) {
            // PostHog automatically handles feature flags
            const updateFlag = () => {
                const value = posthog.isFeatureEnabled(flagKey);
                // isFeatureEnabled returns boolean | undefined, default to false if undefined
                setIsEnabled(!!value);
            };

            // Initial check
            updateFlag();

            // Listen for flag updates
            posthog.onFeatureFlags(updateFlag);
        } else {
            // Local development overrides
            // You can toggle these manually for testing
            const localOverrides: Record<string, boolean> = {
                [FEATURE_FLAGS.PAYWALL_AGGRESSIVE_MODE]: true, // Default to aggressive locally
                [FEATURE_FLAGS.SHOW_REFERRAL_TEASER]: true,
            };
            setIsEnabled(localOverrides[flagKey] ?? defaultValue);
        }
    }, [flagKey, defaultValue]);

    return isEnabled;
};
