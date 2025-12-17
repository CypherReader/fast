import { CheckCircle, Lock, Shield, ArrowRight, Play, Star, Flame } from "lucide-react";
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
          {/* Urgency badge */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.5 }}
            className="inline-flex items-center gap-2 mb-6 px-4 py-2 bg-destructive/20 border border-destructive/40 rounded-full"
          >
            <Flame className="w-4 h-4 text-destructive" />
            <span className="text-sm font-semibold text-destructive">127 people joined in the last 24 hours</span>
          </motion.div>

          {/* Headline */}
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6"
          >
            <span className="text-gradient-hero">You're One Decision Away</span>
            <br />
            <span className="text-foreground">From Your Goal Weight</span>
          </motion.h2>

          {/* Value reminder */}
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="text-lg text-muted-foreground mb-10"
          >
            AI coach + 12,000 member tribe + proven system = Your transformation starts today
          </motion.p>

          {/* Triple CTA options */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.3 }}
            className="flex flex-col sm:flex-row gap-4 justify-center mb-8"
          >
            {/* Primary CTA */}
            <MagneticButton className="text-lg px-8 py-6">
              Start My Transformation ($4.99/mo)
              <ArrowRight className="w-5 h-5 ml-2" />
            </MagneticButton>

            {/* Secondary CTA */}
            <button className="px-6 py-4 bg-card/50 hover:bg-card border border-border hover:border-secondary/50 rounded-xl transition-all flex items-center justify-center gap-2 text-foreground font-semibold">
              <Play className="w-5 h-5" />
              Watch Demo First
            </button>
          </motion.div>

          {/* Tertiary CTA link */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={isInView ? { opacity: 1 } : {}}
            transition={{ duration: 0.6, delay: 0.4 }}
            className="mb-10"
          >
            <a href="#testimonials" className="text-muted-foreground hover:text-foreground transition-colors text-sm underline-offset-4 hover:underline inline-flex items-center gap-1">
              <Star className="w-4 h-4 text-primary" />
              Read 500+ reviews from real users
            </a>
          </motion.div>

          {/* Trust signals */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={isInView ? { opacity: 1 } : {}}
            transition={{ duration: 0.6, delay: 0.5 }}
            className="flex flex-wrap justify-center gap-6 text-sm text-muted-foreground mb-10"
          >
            {[
              { icon: Lock, text: "Secure Payment" },
              { icon: CheckCircle, text: "Cancel Anytime" },
              { icon: Shield, text: "30-Day Money-Back" },
            ].map((item, i) => (
              <motion.div
                key={i}
                whileHover={{ scale: 1.05 }}
                className="flex items-center gap-2"
              >
                <item.icon className="w-4 h-4 text-secondary" />
                {item.text}
              </motion.div>
            ))}
          </motion.div>

          {/* Live counter */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.6 }}
            className="bg-card/80 backdrop-blur rounded-xl p-6 border border-border"
          >
            <div className="text-sm text-muted-foreground mb-2">Active members right now</div>
            <div className="font-display text-4xl font-bold text-primary">
              <AnimatedCounter value={12847} />
            </div>
            <div className="text-sm text-muted-foreground mt-2">
              And growing by <span className="text-primary font-semibold">~127</span> every day
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
};
