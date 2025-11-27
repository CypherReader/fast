import { Star, CheckCircle, Lock, CreditCard, Shield } from "lucide-react";
import { motion, useInView } from "framer-motion";
import { useRef, useState, useEffect } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const testimonials = [
  {
    initials: "SM",
    name: "Sarah M.",
    rating: 5,
    quote: "Lost 22 lbs and never paid a cent. Got every dollar back for 6 months.",
    lost: "22 lbs",
    netCost: "$0 (6 months)",
  },
  {
    initials: "MR",
    name: "Mike R.",
    rating: 5,
    quote: "18h fasts were impossible until I had skin in the game.",
    lost: "14 lbs",
    netCost: "$4/month avg",
  },
  {
    initials: "AK",
    name: "Anna K.",
    rating: 5,
    quote: "30-day streak and full refund every month. Addictive.",
    lost: "19 lbs",
    netCost: "$0 (100% refund)",
  },
];

const trustBadges = [
  { icon: Lock, text: "256-bit SSL Encryption" },
  { icon: CreditCard, text: "Stripe Secured Payments" },
  { icon: CheckCircle, text: "Full Refund Guarantee" },
  { icon: Star, text: "4.8/5 Average Rating" },
];

export const TestimonialsSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });
  const [activeIndex, setActiveIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setActiveIndex((prev) => (prev + 1) % testimonials.length);
    }, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <section className="py-20 md:py-32 bg-gradient-dark relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="accent" showOrbs showGrid />

      <div className="container px-4 relative z-10">
        <div className="max-w-5xl mx-auto">
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold mb-12"
          >
            Does this actually work?
          </motion.h2>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
            {testimonials.map((testimonial, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, y: 30 }}
                animate={isInView ? { opacity: 1, y: 0 } : {}}
                transition={{ duration: 0.5, delay: i * 0.15 }}
                whileHover={{ y: -8, transition: { duration: 0.2 } }}
                className={`bg-card/80 backdrop-blur rounded-xl p-6 border transition-all duration-300 ${
                  activeIndex === i ? "border-secondary/50 scale-105" : "border-border"
                }`}
              >
                {/* Header */}
                <div className="flex items-center justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <motion.div
                      whileHover={{ rotate: 10 }}
                      className="w-10 h-10 rounded-full bg-muted flex items-center justify-center text-sm font-semibold"
                    >
                      {testimonial.initials}
                    </motion.div>
                    <div className="flex items-center gap-1.5">
                      <span className="font-semibold">{testimonial.name}</span>
                      <motion.div
                        animate={{ scale: [1, 1.2, 1] }}
                        transition={{ duration: 3, repeat: Infinity, delay: i * 0.5 }}
                      >
                        <CheckCircle className="w-4 h-4 text-secondary" />
                      </motion.div>
                    </div>
                  </div>
                  <div className="flex">
                    {[...Array(testimonial.rating)].map((_, j) => (
                      <motion.div
                        key={j}
                        initial={{ opacity: 0, scale: 0 }}
                        animate={isInView ? { opacity: 1, scale: 1 } : {}}
                        transition={{ delay: i * 0.15 + j * 0.1 }}
                      >
                        <Star className="w-4 h-4 text-primary fill-primary" />
                      </motion.div>
                    ))}
                  </div>
                </div>

                {/* Quote */}
                <p className="text-foreground mb-6 leading-relaxed">"{testimonial.quote}"</p>

                {/* Stats */}
                <div className="flex justify-between pt-4 border-t border-border">
                  <div>
                    <span className="text-muted-foreground text-sm">Lost: </span>
                    <span className="text-secondary font-semibold">{testimonial.lost}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground text-sm">Net Cost: </span>
                    <span className="text-primary font-semibold font-display">
                      {testimonial.netCost}
                    </span>
                  </div>
                </div>
              </motion.div>
            ))}
          </div>

          {/* Carousel dots for mobile */}
          <div className="flex justify-center gap-2 mb-8 md:hidden">
            {testimonials.map((_, i) => (
              <button
                key={i}
                onClick={() => setActiveIndex(i)}
                className={`w-2 h-2 rounded-full transition-all ${
                  activeIndex === i ? "bg-primary w-6" : "bg-muted"
                }`}
              />
            ))}
          </div>

          {/* Trust badges */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={isInView ? { opacity: 1 } : {}}
            transition={{ delay: 0.6 }}
            className="flex flex-wrap justify-center gap-6 md:gap-10 text-sm text-muted-foreground"
          >
            {trustBadges.map((badge, i) => (
              <motion.div
                key={i}
                whileHover={{ scale: 1.05 }}
                className="flex items-center gap-2"
              >
                <badge.icon className="w-4 h-4 text-secondary" />
                <span>{badge.text}</span>
              </motion.div>
            ))}
          </motion.div>
        </div>
      </div>
    </section>
  );
};
