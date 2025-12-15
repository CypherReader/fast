export type GoalType = 'weight_loss' | 'metabolic' | 'discipline' | 'longevity' | null;
export type FastingPlan = '16:8' | '18:6' | 'omad' | null;

export interface OnboardingState {
    step: number;
    goal: GoalType;
    fastingPlan: FastingPlan;
    completed: boolean;
    email?: string;
    password?: string;
    name?: string;
}

export interface OnboardingContextType {
    state: OnboardingState;
    setStep: (step: number) => void;
    setGoal: (goal: GoalType) => void;
    setFastingPlan: (plan: FastingPlan) => void;
    setCredentials: (credentials: { email?: string; password?: string; name?: string }) => void;
    completeOnboarding: () => void;
    resetOnboarding: () => void;
}

export const initialOnboardingState: OnboardingState = {
    step: 1,
    goal: null,
    fastingPlan: null,
    completed: false,
};
