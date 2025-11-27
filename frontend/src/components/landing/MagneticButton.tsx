import { useRef, useState } from "react";
import { motion } from "framer-motion";
import { cn } from "@/lib/utils";

interface MagneticButtonProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
}

export const MagneticButton = ({ children, className, onClick }: MagneticButtonProps) => {
  const ref = useRef<HTMLButtonElement>(null);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isHovered, setIsHovered] = useState(false);

  const handleMouseMove = (e: React.MouseEvent<HTMLButtonElement>) => {
    if (!ref.current) return;
    
    const rect = ref.current.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    const distanceX = e.clientX - centerX;
    const distanceY = e.clientY - centerY;
    
    // Magnetic pull effect (max 12px movement)
    setPosition({
      x: distanceX * 0.15,
      y: distanceY * 0.15,
    });
  };

  const handleMouseEnter = () => setIsHovered(true);
  const handleMouseLeave = () => {
    setPosition({ x: 0, y: 0 });
    setIsHovered(false);
  };

  return (
    <motion.button
      ref={ref}
      onClick={onClick}
      onMouseMove={handleMouseMove}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
      animate={{ 
        x: position.x, 
        y: position.y,
        scale: isHovered ? 1.05 : 1,
        boxShadow: isHovered 
          ? "0 25px 60px -12px hsl(var(--gold) / 0.7), 0 0 40px hsl(var(--gold) / 0.4), inset 0 1px 0 hsl(var(--gold-light) / 0.5)"
          : "0 10px 30px -10px hsl(var(--gold) / 0.4)"
      }}
      transition={{ type: "spring", stiffness: 200, damping: 15 }}
      whileTap={{ scale: 0.95 }}
      className={cn(
        "relative overflow-hidden text-primary-foreground font-semibold",
        "px-8 py-4 rounded-xl",
        "border border-primary/30",
        className
      )}
    >
      {/* Base gradient background */}
      <motion.div
        className="absolute inset-0 bg-gradient-to-r from-primary via-gold to-primary"
        animate={{
          backgroundPosition: isHovered ? ["0% 50%", "100% 50%"] : "0% 50%",
        }}
        transition={{
          duration: isHovered ? 0.8 : 0,
          ease: "easeInOut",
        }}
        style={{ backgroundSize: "200% 200%" }}
      />
      
      {/* Hover gradient overlay - shifts color */}
      <motion.div
        className="absolute inset-0 bg-gradient-to-r from-gold-light via-primary to-gold-light"
        initial={{ opacity: 0 }}
        animate={{ 
          opacity: isHovered ? 1 : 0,
          backgroundPosition: isHovered ? ["100% 50%", "0% 50%"] : "0% 50%"
        }}
        transition={{ duration: 0.4, ease: "easeOut" }}
        style={{ backgroundSize: "200% 200%" }}
      />

      {/* Shine sweep effect on hover */}
      <motion.div
        className="absolute inset-0 bg-gradient-to-r from-transparent via-white/30 to-transparent"
        initial={{ x: "-100%" }}
        animate={{ x: isHovered ? "100%" : "-100%" }}
        transition={{ duration: 0.6, ease: "easeInOut" }}
      />

      {/* Pulsing glow ring */}
      <motion.div
        className="absolute -inset-1 rounded-xl bg-gradient-to-r from-primary via-gold to-secondary opacity-0 blur-lg"
        animate={{ 
          opacity: isHovered ? [0.5, 0.8, 0.5] : 0,
          scale: isHovered ? [1, 1.05, 1] : 1
        }}
        transition={{ 
          duration: 1.5, 
          repeat: isHovered ? Infinity : 0,
          ease: "easeInOut" 
        }}
      />
      
      {/* Content */}
      <span className="relative z-10 flex items-center gap-2">{children}</span>
    </motion.button>
  );
};
