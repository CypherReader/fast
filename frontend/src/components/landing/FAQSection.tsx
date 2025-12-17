import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const faqs = [
  {
    question: "Will this actually work for me?",
    answer:
      "Yes. FastingHero uses AI to personalize your fasting plan based on YOUR body, YOUR schedule, and YOUR goals. Unlike generic apps, Cortex learns from your patterns and adapts in real-time. Plus, you're joining a tribe of 12,000+ people who've already seen results.",
  },
  {
    question: "I've tried fasting before and failed. Why will this time be different?",
    answer:
      "The difference is Cortex - your AI coach. It knows when you're struggling (like at 3am when you're hungry) and sends personalized support. Plus, our tribe keeps you accountable. You're not alone this time. That's why 87% of our users stick with it past 30 days.",
  },
  {
    question: "What if I fail or can't stick to it?",
    answer:
      "First,, our community support system makes failing much harder than succeeding. Second, we offer a 30-day money-back guarantee. If you're not seeing results or you decide it's not for you, just email us for a full refund. Zero risk.",
  },
  {
    question: "Is $4.99/month actually worth it?",
    answer:
      "Compare: A gym membership is $40/month (and you probably don't go). Personal trainer? $200+/month. Meal delivery service? $300+/month. For $4.99, you get an AI coach, a supportive community, and a proven system. That's less than a coffee at Starbucks.",
  },
  {
    question: "How fast will I see results?",
    answer:
      "Most users notice increased energy within 3-7 days. Weight loss becomes visible around 14-21 days. Significant transformation happens at 30-90 days. Remember: sustainable results take time. FastingHero is built for long-term success, not quick fixes that don't last.",
  },
];

export const FAQSection = () => {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: "-100px" });

  return (
    <section className="py-20 md:py-32 bg-background relative overflow-hidden" ref={ref}>
      <AnimatedSectionBackground variant="subtle" showOrbs showGrid />

      <div className="container px-4 relative z-10">
        <div className="max-w-3xl mx-auto">
          <motion.h2
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
            className="text-4xl md:text-5xl lg:text-6xl font-bold text-center mb-4"
          >
            Your Questions, Answered
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-muted-foreground text-center text-lg mb-12"
          >
            Everything you need to know about FastingHero
          </motion.p>

          <Accordion type="single" collapsible className="space-y-4">
            {faqs.map((faq, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, y: 20 }}
                animate={isInView ? { opacity: 1, y: 0 } : {}}
                transition={{ duration: 0.4, delay: 0.2 + i * 0.1 }}
              >
                <AccordionItem
                  value={`item-${i}`}
                  className="bg-card/80 backdrop-blur rounded-xl border border-border px-6 data-[state=open]:border-secondary/30 transition-colors"
                >
                  <AccordionTrigger className="text-left text-lg font-semibold hover:no-underline py-5 [&[data-state=open]>svg]:text-secondary">
                    {faq.question}
                  </AccordionTrigger>
                  <AccordionContent className="text-muted-foreground pb-5 leading-relaxed">
                    {faq.answer}
                  </AccordionContent>
                </AccordionItem>
              </motion.div>
            ))}
          </Accordion>
        </div>
      </div>
    </section>
  );
};
