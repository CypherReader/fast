import { Navigate } from "react-router-dom";
import { useAuth } from "@/context/AuthContext";
import { PaywallLockScreen } from "./PaywallLockScreen";
import { VaultPromptDialog } from "./VaultPromptDialog";
import { useEffect } from "react";
import { useFeatureFlag, FEATURE_FLAGS } from "@/lib/flags";
import { analytics } from "@/lib/analytics";

export const ProtectedFastingRoute = ({ children }: { children: React.ReactNode }) => {
    const { user, loading } = useAuth();

    // Feature Flag: Aggressive Mode
    // If true: Show Hard Paywall immediately after 3 fasts (or whatever logic)
    // If false: Show Soft Paywall (Vault Prompt) first
    const isAggressiveMode = useFeatureFlag(FEATURE_FLAGS.PAYWALL_AGGRESSIVE_MODE, true);

    useEffect(() => {
        // Track experiment exposure
        if (!loading && user && !user.is_premium) {
            analytics.trackPaywallView(isAggressiveMode ? 'hard_aggressive' : 'soft_vault');
        }
    }, [loading, user, isAggressiveMode]);

    if (loading) return <div className="p-10 text-center text-slate-500">Loading protocol...</div>;

    if (!user) return <Navigate to="/login" />;

    // Premium users bypass everything
    if (user.is_premium) {
        return <>{children}</>;
    }

    // MOCK LOGIC: Simulate "Free Limit Reached"
    // In a real app, check user.fasts_completed or similar
    const freeFastsCompleted = 3; // Mock
    const isLimitReached = freeFastsCompleted >= 3;

    if (isLimitReached) {
        if (isAggressiveMode) {
            // Hard Paywall - Blocks access completely
            return <PaywallLockScreen onJoinVault={() => window.location.href = '/vault-intro'} />;
        } else {
            // Soft Paywall - Shows dialog but allows access
            // For this implementation, let's say Soft Paywall allows access but shows the dialog
            // We need to manage the dialog state here or pass it down.
            // To keep it simple for the route wrapper:
            // If Soft Mode: Render children BUT also render the VaultPromptDialog on top
            return (
                <>
                    <VaultPromptDialog
                        open={true}
                        onOpenChange={() => { }}
                        onJoinVault={() => window.location.href = '/vault-intro'}
                        onContinueFree={() => { }}
                    />
                    {children}
                </>
            );
        }
    }

    return <>{children}</>;
};
