import { Lock, RefreshCw, DollarSign } from "lucide-react";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const steps = [
  {
    number: "1",
    icon: Lock,
    title: "The Deposit",
    description: "Put your money where your goals are",
    detail: "$20 goes into your Vault each month",
    visual: "$20",
    gradient: "from-destructive/20 to-primary/20",
    iconColor: "text-destructive",
  },
  {
    number: "2",
    icon: RefreshCw,
    title: "The Recovery",
    description: "Every completed fast returns $2 of your deposit",
    detail: "10 fasts = full $20 recovered",
    visual: "+$2",
    gradient: "from-background to-secondary/20",
    iconColor: "text-secondary",
  },
  {
    number: "3",
    icon: DollarSign,
    title: "The Refund",
    description: "Get your money back automatically each month",
    detail: "Fast consistently = $0 net cost",
    visual: "$0",
    gradient: "from-primary/20 to-gold-light/20",
    iconColor: "text-primary",
  },
];

export const HowItWorksSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section id="how-it-works" className="py-20 md:py-32 bg-background relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="default" showOrbs showGrid showParticles />

      <div className="container px-4 relative z-10">
        <div className="max-w-5xl mx-auto">
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold text-center mb-4"
          >
            How It Works
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-muted-foreground text-center text-lg mb-16 max-w-2xl mx-auto"
          >
            A simple system that turns your commitment into refunds
          </motion.p>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 relative">
            {/* Animated money flow path */}
            <svg
              className="hidden md:block absolute top-1/2 left-0 right-0 h-2 -translate-y-1/2 z-0 overflow-visible"
              style={{ width: "100%" }}
            >
              <motion.line
                x1="16.5%"
                y1="50%"
                x2="83.5%"
                y2="50%"
                stroke="hsl(var(--border))"
                strokeWidth="2"
                strokeDasharray="8 8"
                initial={{ pathLength: 0 }}
                animate={isInView ? { pathLength: 1 } : {}}
                transition={{ duration: 2, ease: "easeInOut" }}
              />
              {/* Traveling coin */}
              <motion.circle
                r="8"
                fill="hsl(var(--primary))"
                initial={{ cx: "16.5%", cy: "50%" }}
                animate={
                  isInView
                    ? {
                        cx: ["16.5%", "50%", "83.5%", "16.5%"],
                      }
                    : {}
                }
                transition={{
                  duration: 4,
                  repeat: Infinity,
                  ease: "easeInOut",
                  repeatDelay: 1,
                }}
              />
            </svg>

            {steps.map((step, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, y: 30 }}
                animate={isInView ? { opacity: 1, y: 0 } : {}}
                transition={{ duration: 0.5, delay: i * 0.2 }}
                whileHover={{ y: -8, transition: { duration: 0.2 } }}
                className="relative z-10"
              >
                <motion.div
                  className={`bg-gradient-to-br ${step.gradient} backdrop-blur rounded-2xl p-6 border border-border transition-all duration-300`}
                  whileHover={{
                    borderColor: "hsl(var(--primary) / 0.5)",
                    boxShadow: "0 20px 40px -15px hsl(var(--primary) / 0.3)",
                  }}
                >
                  {/* Step number */}
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={isInView ? { scale: 1 } : {}}
                    transition={{ delay: i * 0.2 + 0.3, type: "spring" }}
                    className="absolute -top-3 -left-3 w-10 h-10 rounded-full bg-card border-2 border-border flex items-center justify-center font-display font-bold text-primary"
                  >
                    {step.number}
                  </motion.div>

                  {/* Icon */}
                  <motion.div
                    whileHover={{ rotate: 10, scale: 1.1 }}
                    className={`w-14 h-14 rounded-xl bg-card/80 backdrop-blur flex items-center justify-center mb-4 ${step.iconColor}`}
                  >
                    <step.icon className="w-7 h-7" />
                  </motion.div>

                  {/* Content */}
                  <h3 className="text-xl font-bold mb-2">{step.title}</h3>
                  <p className="text-muted-foreground mb-2">{step.description}</p>
                  <p className="text-sm text-muted-foreground/80 mb-4">{step.detail}</p>

                  {/* Visual with animation */}
                  <motion.div
                    className="font-display text-3xl font-bold text-primary"
                    whileHover={{ scale: 1.1 }}
                  >
                    {step.visual}
                  </motion.div>
                </motion.div>

                {/* Arrow for mobile */}
                {i < steps.length - 1 && (
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={isInView ? { opacity: 1 } : {}}
                    transition={{ delay: i * 0.2 + 0.5 }}
                    className="flex justify-center py-4 md:hidden"
                  >
                    <motion.div
                      animate={{ y: [0, 5, 0] }}
                      transition={{ duration: 1.5, repeat: Infinity }}
                    >
                      <DollarSign className="w-6 h-6 text-primary" />
                    </motion.div>
                  </motion.div>
                )}
              </motion.div>
            ))}
          </div>

          <motion.p
            initial={{ opacity: 0 }}
            animate={isInView ? { opacity: 1 } : {}}
            transition={{ delay: 1 }}
            className="text-center text-muted-foreground mt-12"
          >
            <span className="text-primary font-semibold font-display">87%</span> of users get their
            full deposit back every month
          </motion.p>
        </div>
      </div>
    </section>
  );
};
