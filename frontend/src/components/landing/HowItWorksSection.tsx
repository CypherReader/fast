import { motion, useInView } from "framer-motion";
import { Calendar, Zap, TrendingUp, Users, Brain, ArrowRight } from "lucide-react";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";
import { MagneticButton } from "./MagneticButton";

const steps = [
  {
    days: "Day 1-7",
    title: "Choose Your Fasting Window",
    description: "Cortex AI analyzes your lifestyle and builds a personalized fasting plan. No generic templates.",
    icon: Calendar,
    color: "text-primary",
    bgColor: "bg-primary/10",
    borderColor: "border-primary/30",
  },
  {
    days: "Day 8-21",
    title: "Build Your Streak & Join The Tribe",
    description: "Track your progress, get AI coaching when you need it, and connect with 12,000+ fasters.",
    icon: Zap,
    color: "text-secondary",
    bgColor: "bg-secondary/10",
    borderColor: "border-secondary/30",
  },
  {
    days: "Day 22-30",
    title: "See Results & Transform",
    description: "Feel amazing. Look great. Never look back. This is where everything clicks.",
    icon: TrendingUp,
    color: "text-accent",
    bgColor: "bg-accent/10",
    borderColor: "border-accent/30",
  },
];

const features = [
  {
    icon: Brain,
    title: "Cortex AI Coach",
    description: "Get personalized advice at 3am when you're hungry. Cortex knows you.",
  },
  {
    icon: Users,
    title: "12,000+ Member Tribe",
    description: "You're not alone. Share wins, get support, stay accountable.",
  },
  {
    icon: Zap,
    title: "Gamification That Works",
    description: "Turn discipline into dopamine. Streaks, badges, and real progress.",
  },
];

export const HowItWorksSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="py-20 md:py-32 bg-background relative overflow-hidden" id="how-it-works" ref={ref}>
      <AnimatedSectionBackground variant="subtle" showOrbs showGrid />

      <div className="container px-4 relative z-10">
        <div className="max-w-6xl mx-auto">
          {/* Section header */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-center mb-16"
          >
            <h2 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-4">
              <span className="text-gradient-hero">Your First 30 Days</span>
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Here's exactly what happens when you join FastingHero
            </p>
          </motion.div>

          {/* Timeline steps */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
            {steps.map((step, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, y: 30 }}
                animate={isInView ? { opacity: 1, y: 0 } : {}}
                transition={{ duration: 0.5, delay: i * 0.15 }}
                whileHover={{ y: -8, transition: { duration: 0.2 } }}
                className={`relative rounded-2xl p-8 border ${step.borderColor} ${step.bgColor} transition-all`}
              >
                {/* Step number */}
                <div className="absolute -top-4 left-6 px-3 py-1 bg-background border border-border rounded-full">
                  <span className="text-sm font-bold text-muted-foreground">{step.days}</span>
                </div>

                {/* Icon */}
                <div className={`w-16 h-16 rounded-2xl ${step.bgColor} border ${step.borderColor} flex items-center justify-center mb-6`}>
                  <step.icon className={`w-8 h-8 ${step.color}`} />
                </div>

                {/* Content */}
                <h3 className="text-xl font-bold text-foreground mb-3">{step.title}</h3>
                <p className="text-muted-foreground leading-relaxed">{step.description}</p>
              </motion.div>
            ))}
          </div>

          {/* The FastingHero Difference */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.5 }}
            className="mb-12"
          >
            <h3 className="text-3xl md:text-4xl font-bold text-center mb-12">
              <span className="text-gradient-hero">Meet Cortex:</span> Your AI Fasting Coach That Actually Knows You
            </h3>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {features.map((feature, i) => (
                <motion.div
                  key={i}
                  initial={{ opacity: 0, y: 20 }}
                  animate={isInView ? { opacity: 1, y: 0 } : {}}
                  transition={{ duration: 0.5, delay: 0.6 + i * 0.1 }}
                  className="bg-card/50 backdrop-blur rounded-xl p-6 border border-border hover:border-secondary/50 transition-all"
                >
                  <feature.icon className="w-10 h-10 text-secondary mb-4" />
                  <h4 className="text-lg font-semibold text-foreground mb-2">{feature.title}</h4>
                  <p className="text-sm text-muted-foreground">{feature.description}</p>
                </motion.div>
              ))}
            </div>
          </motion.div>

          {/* CTA */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.9 }}
            className="text-center"
          >
            <MagneticButton className="text-lg" onClick={() => window.location.href = '/onboarding'}>
              Begin Day 1
              <ArrowRight className="w-5 h-5 ml-2" />
            </MagneticButton>
            <p className="text-sm text-muted-foreground mt-4">
              30-day money-back guarantee • Cancel anytime
            </p>
          </motion.div>
        </div>
      </div>
    </section>
  );
};
