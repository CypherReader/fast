import { Check, Star, Sparkles, ArrowRight, Flame } from "lucide-react";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import { MagneticButton } from "./MagneticButton";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const features = [
  "Cortex AI personal fasting coach",
  "Join 12,000+ member tribe",
  "Unlimited fasting tracking",
  "Progress analytics & insights",
  "Science-backed protocols",
  "Premium support",
];

export const PricingSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="py-20 md:py-32 bg-slate-900 dark:bg-gradient-dark relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="accent" showOrbs showGrid showParticles />

      <div className="container px-4 relative z-10">
        <div className="max-w-4xl mx-auto">
          {/* Section header */}
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold text-center mb-4 text-white"
          >
            One Simple Price. Unlimited Transformation.
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-gray-300 text-center text-lg mb-12 max-w-2xl mx-auto"
          >
            Everything included. No upsells. No hidden fees. Just results.
          </motion.p>

          {/* Pricing card */}
          <motion.div
            initial={{ opacity: 0, y: 30, scale: 0.95 }}
            animate={isInView ? { opacity: 1, y: 0, scale: 1 } : {}}
            transition={{ duration: 0.6, delay: 0.2 }}
            whileHover={{ y: -8, transition: { duration: 0.2 } }}
            className="relative rounded-3xl p-8 md:p-12 bg-gradient-to-br from-accent/20 to-purple/10 border-2 border-accent shadow-purple-glow max-w-2xl mx-auto"
          >
            {/* Popular badge */}
            <motion.div
              animate={{ scale: [1, 1.05, 1] }}
              transition={{ duration: 2, repeat: Infinity }}
              className="absolute -top-4 left-1/2 -translate-x-1/2 bg-primary text-primary-foreground text-sm font-bold px-6 py-2 rounded-full flex items-center gap-2 shadow-lg"
            >
              <Star className="w-4 h-4 fill-current" />
              MOST POPULAR
            </motion.div>

            {/* Urgency badge */}
            <div className="flex justify-end mb-4">
              <div className="flex items-center gap-2 px-3 py-1 bg-destructive/20 border border-destructive/40 rounded-full">
                <Flame className="w-3 h-3 text-destructive" />
                <span className="text-xs font-semibold text-destructive">127 joined today</span>
              </div>
            </div>

            {/* Plan name */}
            <div className="text-center mb-6">
              <div className="text-sm font-semibold text-gray-400 tracking-wider mb-2">PREMIUM</div>

              {/* Price */}
              <div className="flex items-baseline justify-center gap-2 mb-2">
                <motion.span
                  initial={{ scale: 0.5 }}
                  animate={isInView ? { scale: 1 } : {}}
                  transition={{ delay: 0.4, type: "spring" }}
                  className="font-display text-6xl md:text-7xl font-bold text-white"
                >
                  $4.99
                </motion.span>
                <span className="text-2xl text-gray-300">/month</span>
              </div>

              <p className="text-gray-300 text-lg">Full access to FastingHero</p>
            </div>

            {/* Features */}
            <ul className="space-y-4 mb-8">
              {features.map((feature, j) => (
                <motion.li
                  key={j}
                  initial={{ opacity: 0, x: -20 }}
                  animate={isInView ? { opacity: 1, x: 0 } : {}}
                  transition={{ delay: 0.5 + j * 0.1 }}
                  className="flex items-center gap-3"
                >
                  <motion.div
                    initial={{ scale: 0 }}
                    animate={isInView ? { scale: 1 } : {}}
                    transition={{ delay: 0.5 + j * 0.1, type: "spring" }}
                  >
                    <Check className="w-6 h-6 text-secondary flex-shrink-0" />
                  </motion.div>
                  <span className="text-gray-100 text-lg">{feature}</span>
                </motion.li>
              ))}
            </ul>

            {/* CTA */}
            <MagneticButton className="w-full justify-center text-lg py-6" onClick={() => window.location.href = '/register'}>
              <Sparkles className="w-5 h-5" />
              Join for $4.99/month
              <ArrowRight className="w-5 h-5" />
            </MagneticButton>

            {/* Guarantee */}
            <div className="mt-6 p-4 bg-secondary/10 border border-secondary/30 rounded-xl text-center">
              <p className="text-sm font-semibold text-secondary mb-1">💰 30-Day Money-Back Guarantee</p>
              <p className="text-xs text-gray-300">See results in 30 days or get a full refund. No questions asked.</p>
            </div>

            {/* Social proof */}
            <p className="text-center text-sm text-gray-400 mt-4">
              Join <span className="text-primary font-semibold">12,847</span> people transforming their lives
            </p>
          </motion.div>
        </div>
      </div>
    </section>
  );
};
