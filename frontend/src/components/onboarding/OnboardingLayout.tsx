import { ReactNode } from 'react';
import { motion } from 'framer-motion';
import { Progress } from '@/components/ui/progress';
import { ArrowLeft } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

interface OnboardingLayoutProps {
  children: ReactNode;
  step: number;
  totalSteps?: number;
  showBack?: boolean;
  backPath?: string;
}

const OnboardingLayout = ({ 
  children, 
  step, 
  totalSteps = 6, 
  showBack = false,
  backPath 
}: OnboardingLayoutProps) => {
  const navigate = useNavigate();
  const progress = (step / totalSteps) * 100;

  return (
    <div className="min-h-screen bg-background flex flex-col">
      {/* Progress Header */}
      <div className="w-full px-6 py-4">
        <div className="max-w-xl mx-auto">
          <div className="flex items-center gap-4 mb-2">
            {showBack && (
              <button 
                onClick={() => backPath ? navigate(backPath) : navigate(-1)}
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                <ArrowLeft className="w-5 h-5" />
              </button>
            )}
            <span className="text-sm text-muted-foreground">
              Step {step} of {totalSteps}
            </span>
          </div>
          <Progress value={progress} className="h-1 bg-muted" />
        </div>
      </div>

      {/* Main Content */}
      <motion.div 
        className="flex-1 flex flex-col items-center justify-center px-5 py-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
      >
        {children}
      </motion.div>
    </div>
  );
};

export default OnboardingLayout;
