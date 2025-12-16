import { Check, Star, Sparkles } from "lucide-react";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import { MagneticButton } from "./MagneticButton";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const plans = [
  {
    name: "VAULT",
    price: "$20",
    period: "/mo",
    description: "Fully refundable based on fasting",
    features: [
      "Unlimited fasting tracking",
      "Recover $2 per completed fast",
      "AI coaching (Cortex)",
      "Tribe membership",
      "Monthly auto-refunds",
    ],
    popular: false,
  },
  {
    name: "VAULT PLUS",
    price: "$30",
    period: "/mo",
    description: "Fully refundable ($3 per fast)",
    features: [
      "Everything in Vault",
      "1:1 coaching (2x/mo)",
      "Blood ketone tracking",
      "Extended fasting protocols",
      "Priority support",
    ],
    popular: true,
  },
];

export const PricingSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="py-20 md:py-32 bg-slate-900 dark:bg-gradient-dark relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="accent" showOrbs showGrid showParticles />

      <div className="container px-4 relative z-10">
        <div className="max-w-4xl mx-auto">
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold text-center mb-4 text-white"
          >
            Simple, transparent pricing
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-gray-300 text-center text-lg mb-12"
          >
            Pay to commit. Get it all back if you fast consistently.
          </motion.p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {plans.map((plan, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, y: 30 }}
                animate={isInView ? { opacity: 1, y: 0 } : {}}
                transition={{ duration: 0.5, delay: i * 0.2 }}
                whileHover={{
                  y: -8,
                  rotateY: 2,
                  rotateX: 2,
                  transition: { duration: 0.2 },
                }}
                className={`relative rounded-2xl p-8 transition-all duration-300 backdrop-blur ${plan.popular
                  ? "bg-gradient-to-br from-accent/20 to-purple/10 border-2 border-accent shadow-purple-glow"
                  : "bg-card/80 border border-border hover:border-secondary/30"
                  }`}
              >
                {/* Popular badge */}
                {plan.popular && (
                  <motion.div
                    animate={{ scale: [1, 1.05, 1] }}
                    transition={{ duration: 2, repeat: Infinity }}
                    className="absolute -top-3 left-6 bg-primary text-primary-foreground text-xs font-bold px-3 py-1 rounded-full flex items-center gap-1"
                  >
                    <Star className="w-3 h-3 fill-current" />
                    MOST POPULAR
                  </motion.div>
                )}

                {/* Plan name */}
                <div className="text-sm font-semibold text-gray-400 tracking-wider mb-2">
                  {plan.name}
                </div>

                {/* Price */}
                <div className="flex items-baseline gap-1 mb-2">
                  <motion.span
                    initial={{ scale: 0.5 }}
                    animate={isInView ? { scale: 1 } : {}}
                    transition={{ delay: i * 0.2 + 0.3, type: "spring" }}
                    className="font-display text-5xl font-bold text-white"
                  >
                    {plan.price}
                  </motion.span>
                  <span className="text-gray-300">{plan.period}</span>
                </div>

                {/* Description */}
                <p className="text-gray-300 mb-6">{plan.description}</p>

                {/* Features */}
                <ul className="space-y-3 mb-8">
                  {plan.features.map((feature, j) => (
                    <motion.li
                      key={j}
                      initial={{ opacity: 0, x: -20 }}
                      animate={isInView ? { opacity: 1, x: 0 } : {}}
                      transition={{ delay: i * 0.2 + j * 0.1 }}
                      className="flex items-center gap-3"
                    >
                      <motion.div
                        initial={{ scale: 0 }}
                        animate={isInView ? { scale: 1 } : {}}
                        transition={{ delay: i * 0.2 + j * 0.1, type: "spring" }}
                      >
                        <Check className="w-5 h-5 text-secondary flex-shrink-0" />
                      </motion.div>
                      <span className="text-foreground">{feature}</span>
                    </motion.li>
                  ))}
                </ul>

                {/* CTA */}
                <MagneticButton className="w-full justify-center">
                  <Sparkles className="w-4 h-4" />
                  Start Your Vault
                </MagneticButton>

                {/* Net cost */}
                <p className="text-sm text-muted-foreground mt-4 text-center">
                  Net cost: <span className="text-primary font-semibold">$0</span> if you fast
                  consistently
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
};
