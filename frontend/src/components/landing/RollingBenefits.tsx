import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Flame, Brain, Heart, Zap, Scale, Clock } from "lucide-react";

const benefits = [
  { icon: Flame, text: "Burn fat faster", color: "text-primary" },
  { icon: Brain, text: "Boost mental clarity", color: "text-accent" },
  { icon: Heart, text: "Improve heart health", color: "text-destructive" },
  { icon: Zap, text: "Increase energy levels", color: "text-secondary" },
  { icon: Scale, text: "Lose weight naturally", color: "text-primary" },
  { icon: Clock, text: "Trigger autophagy", color: "text-accent" },
];

export const RollingBenefits = () => {
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentIndex((prev) => (prev + 1) % benefits.length);
    }, 2500);
    return () => clearInterval(interval);
  }, []);

  const currentBenefit = benefits[currentIndex];
  const Icon = currentBenefit.icon;

  return (
    <div className="flex items-center justify-center gap-3 h-10">
      <motion.div
        animate={{ scale: [1, 1.2, 1] }}
        transition={{ duration: 2, repeat: Infinity }}
        className="w-2 h-2 rounded-full bg-secondary"
      />
      
      <div className="relative overflow-hidden h-8 min-w-[280px]">
        <AnimatePresence mode="wait">
          <motion.div
            key={currentIndex}
            initial={{ y: 30, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            exit={{ y: -30, opacity: 0 }}
            transition={{ duration: 0.4, ease: "easeOut" }}
            className="absolute inset-0 flex items-center justify-center gap-2"
          >
            <Icon className={`w-5 h-5 ${currentBenefit.color}`} />
            <span className="font-semibold text-foreground text-lg">
              {currentBenefit.text}
            </span>
          </motion.div>
        </AnimatePresence>
      </div>

      {/* Progress dots */}
      <div className="flex gap-1.5">
        {benefits.map((_, i) => (
          <motion.div
            key={i}
            className={`w-1.5 h-1.5 rounded-full transition-colors duration-300 ${
              i === currentIndex ? "bg-primary" : "bg-muted"
            }`}
            animate={i === currentIndex ? { scale: [1, 1.3, 1] } : {}}
            transition={{ duration: 0.3 }}
          />
        ))}
      </div>
    </div>
  );
};
