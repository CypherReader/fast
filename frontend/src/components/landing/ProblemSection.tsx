import { motion, useInView } from "framer-motion";
import { AlertCircle, X, TrendingDown } from "lucide-react";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const problems = [
  {
    text: "No accountability - you quit after 3 days",
    icon: X,
  },
  {
    text: "Generic advice that doesn't fit YOUR body",
    icon: X,
  },
  {
    text: "Boring tracking with no real motivation",
    icon: X,
  },
];

export const ProblemSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="py-20 md:py-32 bg-slate-900 dark:bg-gradient-dark relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="accent" showOrbs showGrid />

      <div className="container px-4 relative z-10">
        <div className="max-w-4xl mx-auto text-center">
          {/* Alert badge */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={isInView ? { opacity: 1, scale: 1 } : {}}
            transition={{ duration: 0.5 }}
            className="inline-flex items-center gap-2 mb-6 px-4 py-2 bg-destructive/10 border border-destructive/30 rounded-full"
          >
            <AlertCircle className="w-4 h-4 text-destructive" />
            <span className="text-sm font-semibold text-destructive">Reality Check</span>
          </motion.div>

          {/* Headline */}
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold text-white mb-6"
          >
            Why <span className="text-destructive">97%</span> of Fasting Apps Fail
            <br />
            <span className="text-muted-foreground text-2xl md:text-3xl">(And Why FastingHero Is Different)</span>
          </motion.h2>

          {/* Problem agitation */}
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="text-lg md:text-xl text-gray-300 mb-12 max-w-2xl mx-auto"
          >
            You've tried fasting apps before. They work... for{" "}
            <span className="text-white font-semibold">3 days</span>. Then you're back to square one.
          </motion.p>

          {/* Problems list */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
            {problems.map((problem, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, y: 20 }}
                animate={isInView ? { opacity: 1, y: 0 } : {}}
                transition={{ duration: 0.5, delay: 0.3 + i * 0.1 }}
                className="bg-destructive/5 border border-destructive/20 rounded-xl p-6 text-left"
              >
                <problem.icon className="w-6 h-6 text-destructive mb-3" />
                <p className="text-gray-200 leading-relaxed">{problem.text}</p>
              </motion.div>
            ))}
          </div>

          {/* Solution tease */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.7 }}
            className="relative"
          >
            <div className="absolute top-1/2 left-0 right-0 h-px bg-gradient-to-r from-transparent via-primary/50 to-transparent" />
            <div className="relative inline-flex items-center gap-3 px-6 py-3 bg-slate-900 border-2 border-primary/30 rounded-full">
              <TrendingDown className="w-5 h-5 text-primary" />
              <span className="text-lg font-semibold text-white">Until now.</span>
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
};
