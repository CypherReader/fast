import React, { createContext, useContext, useState, ReactNode } from 'react';
import {
  type GoalType,
  type FastingPlan,
  type OnboardingState,
  type OnboardingContextType,
  initialOnboardingState
} from './onboarding-types';

// Re-export types for backward compatibility
export type { GoalType, FastingPlan, OnboardingState, OnboardingContextType };


const OnboardingContext = createContext<OnboardingContextType | undefined>(undefined);

export const OnboardingProvider = ({ children }: { children: ReactNode }) => {
  const [state, setState] = useState<OnboardingState>(initialOnboardingState);

  const setStep = (step: number) => setState(prev => ({ ...prev, step }));
  const setGoal = (goal: GoalType) => setState(prev => ({ ...prev, goal }));
  const setFastingPlan = (fastingPlan: FastingPlan) => setState(prev => ({ ...prev, fastingPlan }));
  const setCredentials = (credentials: { email?: string; password?: string; name?: string }) =>
    setState(prev => ({ ...prev, ...credentials }));
  const completeOnboarding = () => setState(prev => ({ ...prev, completed: true }));
  const resetOnboarding = () => setState(initialOnboardingState);

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
