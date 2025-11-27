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
    question: "What if I don't recover all my money back?",
    answer:
      "You'll still recover based on your discipline. Complete 5 fasts → Recover $10 back. Complete 10 fasts → Recover $20 back. Average user recovers $18.73/month. It's progressive, not all-or-nothing.",
  },
  {
    question: "Is my payment information secure?",
    answer:
      "Absolutely. We use Stripe for payments (the same platform used by Amazon, Google, and Shopify). Your card details are never stored on our servers. All transactions are protected with 256-bit SSL encryption.",
  },
  {
    question: "Can I cancel anytime?",
    answer:
      "Yes, you can cancel your subscription at any time with immediate effect. You'll keep any money you've already recovered for that billing period.",
  },
  {
    question: "What if I miss a few days?",
    answer:
      "No problem! Our recovery system is progressive. Missing a day or two won't wipe out your progress. You recover for every fast you complete, regardless of streaks.",
  },
  {
    question: "Do other fasting apps do this?",
    answer:
      "No. FastingHero is the only fasting app with a commitment-based refund system. Other apps are free (no commitment) or charge you without giving anything back. We're the only app where your money motivates you AND returns to you.",
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
            Questions?
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
