import { Check, X } from "lucide-react";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const features = [
  { name: "Refundable deposit", fastingHero: "✓ $20/mo", others: false },
  { name: "Financial commitment", fastingHero: true, others: false },
  { name: "Tribe accountability", fastingHero: true, others: "Some" },
  { name: "AI coaching", fastingHero: true, others: "Some" },
  { name: "Success rate*", fastingHero: "87%", others: "42-52%" },
];

export const ComparisonSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="py-20 md:py-32 bg-background relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="default" showOrbs showGrid />

      <div className="container px-4 relative z-10">
        <div className="max-w-3xl mx-auto">
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold text-center mb-4"
          >
            Why FastingHero wins
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-muted-foreground text-center text-lg mb-12"
          >
            The only fasting app with real skin in the game
          </motion.p>

          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="bg-card/80 backdrop-blur rounded-2xl border border-border overflow-hidden"
          >
            {/* Header */}
            <div className="grid grid-cols-3 bg-muted/50">
              <div className="p-4 text-sm font-medium text-muted-foreground">Feature</div>
              <div className="p-4 text-center font-bold text-primary bg-primary/5 border-x border-primary/20">
                FastingHero
              </div>
              <div className="p-4 text-center text-sm font-medium text-muted-foreground">
                Other Apps
              </div>
            </div>

            {/* Rows */}
            {features.map((feature, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, x: -20 }}
                animate={isInView ? { opacity: 1, x: 0 } : {}}
                transition={{ duration: 0.4, delay: 0.3 + i * 0.1 }}
                className="grid grid-cols-3 border-t border-border hover:bg-muted/30 transition-colors"
              >
                <div className="p-4 text-foreground">{feature.name}</div>
                <div className="p-4 flex justify-center items-center bg-primary/5 border-x border-primary/20">
                  {typeof feature.fastingHero === "boolean" ? (
                    <motion.div
                      initial={{ scale: 0 }}
                      animate={isInView ? { scale: 1 } : {}}
                      transition={{ delay: 0.4 + i * 0.1, type: "spring" }}
                    >
                      <Check className="w-5 h-5 text-secondary" />
                    </motion.div>
                  ) : (
                    <span className="font-display font-semibold text-primary">
                      {feature.fastingHero}
                    </span>
                  )}
                </div>
                <div className="p-4 flex justify-center items-center">
                  {typeof feature.others === "boolean" ? (
                    <motion.div
                      initial={{ scale: 0 }}
                      animate={isInView ? { scale: 1 } : {}}
                      transition={{ delay: 0.4 + i * 0.1, type: "spring" }}
                    >
                      <X className="w-5 h-5 text-destructive" />
                    </motion.div>
                  ) : (
                    <span className="text-muted-foreground text-sm">{feature.others}</span>
                  )}
                </div>
              </motion.div>
            ))}
          </motion.div>

          <motion.p
            initial={{ opacity: 0 }}
            animate={isInView ? { opacity: 1 } : {}}
            transition={{ delay: 0.8 }}
            className="text-center text-muted-foreground text-sm mt-6"
          >
            *Based on D30 retention data
          </motion.p>
        </div>
      </div>
    </section>
  );
};
