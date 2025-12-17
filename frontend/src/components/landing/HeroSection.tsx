import { DollarSign, Lock, CheckCircle, Shield } from "lucide-react";
import { motion } from "framer-motion";
import { ClockVideoBackground } from "./ClockVideoBackground";
import { FloatingParticles } from "./FloatingParticles";
import { FastingTimerHero } from "./FastingTimerHero";
import { MagneticButton } from "./MagneticButton";

const avatars = ["SM", "MR", "AK", "JD", "LB"];

export const HeroSection = () => {
  const scrollToHowItWorks = () => {
    document.getElementById("how-it-works")?.scrollIntoView({ behavior: "smooth" });
  };

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
                className="flex flex-wrap justify-center lg:justify-start gap-4 text-xs text-muted-foreground mb-6"
              >
                <div className="flex items-center gap-1.5 bg-card/50 backdrop-blur px-3 py-1.5 rounded-full border border-border/50">
                  <Lock className="w-3 h-3 text-secondary" />
                  256-bit SSL
                </div>
                <div className="flex items-center gap-1.5 bg-card/50 backdrop-blur px-3 py-1.5 rounded-full border border-border/50">
                  <CheckCircle className="w-3 h-3 text-secondary" />
                  Stripe Secured
                </div>
                <div className="flex items-center gap-1.5 bg-card/50 backdrop-blur px-3 py-1.5 rounded-full border border-border/50">
                  <Shield className="w-3 h-3 text-secondary" />
                  30-Day Guarantee
                </div>
              </motion.div>

              {/* Main headline */}
              <motion.h1
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.7, delay: 0.1 }}
                className="text-4xl md:text-5xl lg:text-6xl xl:text-7xl font-extrabold mb-4 leading-tight"
              >
                <span className="text-gradient-hero">The AI-assisted fasting app</span>
              </motion.h1>

              <motion.h1
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.7, delay: 0.2 }}
                className="text-4xl md:text-5xl lg:text-6xl xl:text-7xl font-extrabold mb-6"
              >
                <span className="text-foreground"></span>
              </motion.h1>

              {/* Subheadline */}
              <motion.p
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.3 }}
                className="text-lg md:text-xl text-muted-foreground mb-6 max-w-xl mx-auto lg:mx-0"
              >
                {/* HIDDEN FOR V2: Deposit $20/month. Fast consistently. Get it all back. */}
                Track your fasting journey. Build discipline. Achieve your goals.
                <br className="hidden md:block" />
                Start your transformation today with{" "}
                <span className="text-primary font-semibold">FastingHero</span>.
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
                  <span className="text-primary font-semibold">12,847</span> people paying $0/month
                </span>
              </motion.div>

              {/* CTAs */}
              <motion.div
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, delay: 0.4 }}
                className="flex flex-col items-center lg:items-start gap-3"
              >
                {/* HIDDEN FOR V2: Start Your Vault - $20 Refundable */}
                <MagneticButton className="text-lg">
                  Get Started Free
                </MagneticButton>
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
