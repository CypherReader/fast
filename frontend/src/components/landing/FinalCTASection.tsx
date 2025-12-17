import { DollarSign, Lock, CheckCircle, Shield, Quote } from "lucide-react";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import { AnimatedCounter } from "./AnimatedCounter";
import { MagneticButton } from "./MagneticButton";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

export const FinalCTASection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="py-20 md:py-32 bg-gradient-dark relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="accent" showOrbs showGrid showParticles />

      {/* Extra spotlight effect */}
      <motion.div
        initial={{ opacity: 0, scale: 0.8 }}
        animate={isInView ? { opacity: 1, scale: 1 } : {}}
        transition={{ duration: 1 }}
        className="absolute top-0 left-1/2 -translate-x-1/2 w-[600px] h-[400px] bg-primary/10 rounded-full blur-[150px] pointer-events-none"
      />

      <div className="container px-4 relative z-10">
        <div className="max-w-3xl mx-auto text-center">
          {/* Headline */}
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6"
          >
            <span className="text-gradient-hero">Start Your</span>
            <br />
            <span className="text-foreground">Transformation Today</span>
          </motion.h2>

          {/* Stats with animated counters */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="grid grid-cols-3 gap-4 md:gap-8 mb-10"
          >
            <div className="text-center">
              <div className="font-display text-2xl md:text-4xl font-bold text-primary mb-1">
                <AnimatedCounter value={12} suffix="k+" />
              </div>
              <div className="text-xs md:text-sm text-muted-foreground">
                Active Users
              </div>
            </div>
            <div className="text-center">
              <div className="font-display text-2xl md:text-4xl font-bold text-secondary mb-1">
                <AnimatedCounter value={18} suffix=" lbs" />
              </div>
              <div className="text-xs md:text-sm text-muted-foreground">Avg Weight Loss</div>
            </div>
            <div className="text-center">
              <div className="font-display text-2xl md:text-4xl font-bold text-accent mb-1">
                $<AnimatedCounter value={4.99} decimals={2} />
              </div>
              <div className="text-xs md:text-sm text-muted-foreground">Per Month</div>
            </div>
          </motion.div>

          {/* CTA */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.3 }}
            className="mb-6"
          >
            <MagneticButton className="text-lg">
              Join 12,847 Active Fasters
            </MagneticButton>
          </motion.div>

          {/* Trust signals */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={isInView ? { opacity: 1 } : {}}
            transition={{ duration: 0.6, delay: 0.4 }}
            className="flex flex-wrap justify-center gap-4 md:gap-6 text-sm text-muted-foreground mb-10"
          >
            {[
              { icon: Lock, text: "Secure Payment" },
              { icon: CheckCircle, text: "Cancel Anytime" },
              { icon: Shield, text: "30-Day Guarantee" },
            ].map((item, i) => (
              <motion.div
                key={i}
                whileHover={{ scale: 1.05 }}
                className="flex items-center gap-1.5"
              >
                <item.icon className="w-4 h-4 text-secondary" />
                {item.text}
              </motion.div>
            ))}
          </motion.div>

          {/* Final testimonial */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.5 }}
            whileHover={{ scale: 1.02 }}
            className="bg-card/80 backdrop-blur rounded-xl p-6 border border-border max-w-lg mx-auto"
          >
            <Quote className="w-8 h-8 text-primary/30 mb-3" />
            <p className="text-foreground italic mb-4">
              "I've tried every fasting app. FastingHero is the only one that worked because of the AI coaching and supportive community."
            </p>
            <div className="flex items-center justify-center gap-2">
              <motion.div
                animate={{ scale: [1, 1.1, 1] }}
                transition={{ duration: 2, repeat: Infinity }}
                className="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-xs font-semibold"
              >
                MR
              </motion.div>
              <span className="font-semibold">Michael R.</span>
              <span className="text-muted-foreground text-sm">(Lost 23 lbs in 3 months)</span>
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
};
