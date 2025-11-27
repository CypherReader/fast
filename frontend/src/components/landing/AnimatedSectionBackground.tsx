import { motion, useScroll, useTransform } from "framer-motion";
import { useRef } from "react";

interface AnimatedSectionBackgroundProps {
  variant?: "default" | "accent" | "subtle";
  showGrid?: boolean;
  showOrbs?: boolean;
  showParticles?: boolean;
}

export const AnimatedSectionBackground = ({
  variant = "default",
  showGrid = true,
  showOrbs = true,
  showParticles = false,
}: AnimatedSectionBackgroundProps) => {
  const ref = useRef<HTMLDivElement>(null);
  const { scrollYProgress } = useScroll({
    target: ref,
    offset: ["start end", "end start"],
  });

  // Parallax transforms at different speeds
  const orbY1 = useTransform(scrollYProgress, [0, 1], [100, -100]);
  const orbY2 = useTransform(scrollYProgress, [0, 1], [-50, 150]);
  const orbX1 = useTransform(scrollYProgress, [0, 1], [-30, 30]);
  const orbX2 = useTransform(scrollYProgress, [0, 1], [30, -30]);
  const gridY = useTransform(scrollYProgress, [0, 1], [0, -50]);
  const particleY = useTransform(scrollYProgress, [0, 1], [50, -100]);

  const orbColors = {
    default: ["primary", "secondary"],
    accent: ["accent", "primary"],
    subtle: ["muted", "primary"],
  };

  const colors = orbColors[variant];

  return (
    <div ref={ref} className="absolute inset-0 overflow-hidden pointer-events-none">
      {/* Animated gradient orbs with parallax */}
      {showOrbs && (
        <>
          <motion.div
            className="absolute -top-40 -left-40 w-[500px] h-[500px] rounded-full blur-[120px]"
            style={{ 
              y: orbY1,
              x: orbX1,
              background: `radial-gradient(circle, hsl(var(--${colors[0]}) / 0.15), transparent)` 
            }}
            animate={{
              scale: [1, 1.1, 1],
            }}
            transition={{ duration: 15, repeat: Infinity, ease: "easeInOut" }}
          />
          <motion.div
            className="absolute -bottom-40 -right-40 w-[400px] h-[400px] rounded-full blur-[100px]"
            style={{ 
              y: orbY2,
              x: orbX2,
              background: `radial-gradient(circle, hsl(var(--${colors[1]}) / 0.1), transparent)` 
            }}
            animate={{
              scale: [1, 1.2, 1],
            }}
            transition={{ duration: 12, repeat: Infinity, ease: "easeInOut", delay: 2 }}
          />
          {/* Additional parallax orb */}
          <motion.div
            className="absolute top-1/3 right-1/4 w-[300px] h-[300px] rounded-full blur-[80px]"
            style={{ 
              y: useTransform(scrollYProgress, [0, 1], [-80, 80]),
              background: `radial-gradient(circle, hsl(var(--gold) / 0.08), transparent)` 
            }}
          />
        </>
      )}

      {/* Subtle grid pattern with parallax */}
      {showGrid && (
        <motion.div
          className="absolute inset-0 opacity-[0.02]"
          style={{
            y: gridY,
            backgroundImage: `
              linear-gradient(hsl(var(--foreground)) 1px, transparent 1px),
              linear-gradient(90deg, hsl(var(--foreground)) 1px, transparent 1px)
            `,
            backgroundSize: "60px 60px",
          }}
        />
      )}

      {/* Floating particles with parallax */}
      {showParticles && (
        <>
          {[...Array(12)].map((_, i) => (
            <motion.div
              key={i}
              className="absolute w-1 h-1 rounded-full bg-primary/40"
              style={{
                left: `${5 + i * 8}%`,
                y: particleY,
              }}
              initial={{ y: "100%" }}
              animate={{
                y: "-10%",
                opacity: [0, 0.6, 0],
              }}
              transition={{
                duration: 8 + Math.random() * 4,
                repeat: Infinity,
                delay: i * 1.2,
                ease: "linear",
              }}
            />
          ))}
        </>
      )}

      {/* Subtle scan line effect */}
      <motion.div
        className="absolute inset-0 bg-gradient-to-b from-transparent via-primary/[0.02] to-transparent h-[200px]"
        animate={{ y: ["-100%", "500%"] }}
        transition={{ duration: 8, repeat: Infinity, ease: "linear" }}
      />
    </div>
  );
};
