import { motion } from "framer-motion";
import { ArrowRight, CheckCircle, Lock, Shield, Star } from "lucide-react";
import { ClockVideoBackground } from "./ClockVideoBackground";
import { FloatingParticles } from "./FloatingParticles";
import { FastingTimerHero } from "./FastingTimerHero";
import { MagneticButton } from "./MagneticButton";

const avatars = ["SM", "MR", "AK", "JD", "LB"];

const trustBadges = [
  { icon: Star, text: "4.9★ Rating" },
  { icon: Shield, text: "30-Day Guarantee" },
  { icon: Lock, text: "Secure Payment" },
];

export const HeroSection = () => {
  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden bg-background">
      {/* Cinematic clock video background */}
      <ClockVideoBackground />
      <FloatingParticles />

      {/* Dark gradient overlay for readability */}
      <div className="absolute inset-0 bg-gradient-to-b from-background/80 via-background/60 to-background" />
      <div className="absolute inset-0 bg-gradient-to-r from-background/90 via-transparent to-background/90" />

      <div className="container relative z-10 px-4 py-20 md:py-32">
        <div className="max-w-6xl mx-auto">
          {/* Two-column layout */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            {/* Left side - Content */}
            <div className="text-center lg:text-left">
              {/* Trust badges */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5 }}
                className="flex flex-wrap justify-center lg:justify-start gap-3 text-xs text-muted-foreground mb-6"
              >
                {trustBadges.map((badge) => (
                  <div
                    key={badge.text}
                    className="flex items-center gap-1.5 bg-card/50 backdrop-blur px-3 py-1.5 rounded-full border border-border/50"
                  >
                    <badge.icon className="w-3 h-3 text-secondary" />
                    {badge.text}
                  </div>
                ))}
              </motion.div>

              {/* Main headline - VIRAL HOOK */}
              <motion.h1
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.7, delay: 0.1 }}
                className="text-4xl md:text-5xl lg:text-6xl xl:text-7xl font-extrabold mb-6 leading-tight"
              >
                <span className="text-gradient-hero">The AI Fasting Coach</span>
                <br />
                <span className="text-foreground">That Helped 12,000+ People</span>
                <br />
                <span className="text-primary">Lose Weight While They Sleep</span>
              </motion.h1>

              {/* Subheadline */}
              <motion.p
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.3 }}
                className="text-lg md:text-xl text-muted-foreground mb-6 max-w-xl mx-auto lg:mx-0"
              >
                Join the fastest-growing fasting movement. AI-powered. Science-backed.{" "}
                <span className="text-primary font-semibold">$4.99/month</span>
              </motion.p>

              {/* Social proof - Avatar stack */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.35 }}
                className="flex items-center justify-center lg:justify-start gap-3 mb-8"
              >
                <div className="flex -space-x-2">
                  {avatars.map((initials, i) => (
                    <motion.div
                      key={i}
                      initial={{ opacity: 0, scale: 0 }}
                      animate={{ opacity: 1, scale: 1 }}
                      transition={{ delay: 0.4 + i * 0.05, type: "spring" }}
                      className="w-8 h-8 rounded-full bg-gradient-to-br from-primary/80 to-secondary/80 border-2 border-background flex items-center justify-center text-xs font-semibold text-primary-foreground"
                    >
                      {initials}
                    </motion.div>
                  ))}
                </div>
                <span className="text-sm text-muted-foreground">
                  <span className="text-primary font-semibold">12,847</span> people transformed
                </span>
              </motion.div>

              {/* Primary CTA */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.4 }}
                className="flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4 mb-6"
              >
                <MagneticButton className="text-lg px-8 py-6" onClick={() => window.location.href = '/register'}>
                  Start My Transformation
                  <ArrowRight className="w-5 h-5 ml-2" />
                </MagneticButton>
                <button className="text-muted-foreground hover:text-foreground transition-colors text-sm underline-offset-4 hover:underline">
                  Watch Demo (2 min)
                </button>
              </motion.div>

              {/* Trust line with checkmarks */}
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.5, delay: 0.5 }}
                className="flex flex-wrap justify-center lg:justify-start gap-4 text-sm text-muted-foreground"
              >
                <div className="flex items-center gap-1.5">
                  <CheckCircle className="w-4 h-4 text-secondary" />
                  <span>No credit card needed</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <CheckCircle className="w-4 h-4 text-secondary" />
                  <span>Cancel anytime</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <CheckCircle className="w-4 h-4 text-secondary" />
                  <span>30-day guarantee</span>
                </div>
              </motion.div>
            </div>

            {/* Right side - Fasting Timer */}
            <motion.div
              initial={{ opacity: 0, scale: 0.9, x: 50 }}
              animate={{ opacity: 1, scale: 1, x: 0 }}
              transition={{ duration: 0.8, delay: 0.3 }}
              className="flex justify-center lg:justify-end"
            >
              <div className="bg-card/80 backdrop-blur-xl rounded-3xl p-6 md:p-8 border border-border/50 shadow-2xl shadow-primary/10">
                <FastingTimerHero />
              </div>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
};
