import React, { createContext, useContext, useState, ReactNode } from 'react';

export type GoalType = 'weight_loss' | 'metabolic' | 'discipline' | 'longevity' | null;
export type FastingPlan = '16:8' | '18:6' | 'omad' | null;

interface OnboardingState {
  step: number;
  goal: GoalType;
  fastingPlan: FastingPlan;
  completed: boolean;
  email?: string;
  password?: string;
  name?: string;
}

interface OnboardingContextType {
  state: OnboardingState;
  setStep: (step: number) => void;
  setGoal: (goal: GoalType) => void;
  setFastingPlan: (plan: FastingPlan) => void;
  setCredentials: (credentials: { email?: string; password?: string; name?: string }) => void;
  completeOnboarding: () => void;
  resetOnboarding: () => void;
}

const initialState: OnboardingState = {
  step: 1,
  goal: null,
  fastingPlan: null,
  completed: false,
};

const OnboardingContext = createContext<OnboardingContextType | undefined>(undefined);

export const OnboardingProvider = ({ children }: { children: ReactNode }) => {
  const [state, setState] = useState<OnboardingState>(initialState);

  const setStep = (step: number) => setState(prev => ({ ...prev, step }));
  const setGoal = (goal: GoalType) => setState(prev => ({ ...prev, goal }));
  const setFastingPlan = (fastingPlan: FastingPlan) => setState(prev => ({ ...prev, fastingPlan }));
  const setCredentials = (credentials: { email?: string; password?: string; name?: string }) =>
    setState(prev => ({ ...prev, ...credentials }));
  const completeOnboarding = () => setState(prev => ({ ...prev, completed: true }));
  const resetOnboarding = () => setState(initialState);

  return (
    <OnboardingContext.Provider value={{ state, setStep, setGoal, setFastingPlan, setCredentials, completeOnboarding, resetOnboarding }}>
      {children}
    </OnboardingContext.Provider>
  );
};

export const useOnboarding = () => {
  const context = useContext(OnboardingContext);
  if (!context) throw new Error('useOnboarding must be used within OnboardingProvider');
  return context;
};
