import { X } from "lucide-react";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const problems = [
  "No accountability",
  "No skin in the game",
  "Free = no commitment",
  "Just another tracker",
];

export const ProblemSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="py-20 md:py-32 bg-gradient-dark relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="subtle" showOrbs showGrid />

      <div className="container px-4 relative z-10">
        <div className="max-w-4xl mx-auto">
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold mb-12"
          >
            The Problem With Fasting Apps
          </motion.h2>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
            {problems.map((problem, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, x: -30 }}
                animate={isInView ? { opacity: 1, x: 0 } : {}}
                transition={{ duration: 0.5, delay: i * 0.1 }}
                whileHover={{ x: 10, transition: { duration: 0.2 } }}
                className="relative bg-card/80 backdrop-blur rounded-xl p-5 border border-border hover:border-destructive/30 transition-colors"
              >
                <motion.div
                  initial={{ scaleY: 0 }}
                  animate={isInView ? { scaleY: 1 } : {}}
                  transition={{ delay: i * 0.1 + 0.3 }}
                  className="absolute left-0 top-0 bottom-0 w-1 bg-destructive rounded-l-xl origin-top"
                />
                <div className="flex items-center gap-3 pl-3">
                  <motion.div
                    initial={{ rotate: -90, opacity: 0 }}
                    animate={isInView ? { rotate: 0, opacity: 1 } : {}}
                    transition={{ delay: i * 0.1 + 0.2, type: "spring" }}
                  >
                    <X className="w-5 h-5 text-destructive flex-shrink-0" />
                  </motion.div>
                  <span className="text-foreground font-medium">{problem}</span>
                </div>
              </motion.div>
            ))}
          </div>

          <motion.p
            initial={{ opacity: 0 }}
            animate={isInView ? { opacity: 1 } : {}}
            transition={{ delay: 0.6 }}
            className="text-lg"
          >
            <span className="text-primary font-bold">FastingHero</span>
            <span className="text-foreground"> is different.</span>
          </motion.p>
        </div>
      </div>
    </section>
  );
};
