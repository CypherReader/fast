import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

// HIDDEN FOR V2: Vault-related FAQ questions replaced with generic questions
const faqs = [
  {
    question: "How does FastingHero help me succeed?",
    answer:
      "FastingHero combines AI coaching, progress tracking, and community support to keep you motivated. Our Cortex AI provides personalized insights during your fasting journey, helping you understand what's happening in your body at each stage.",
  },
  {
    question: "Is my payment information secure?",
    answer:
      "Absolutely. We use Stripe for payments (the same platform used by Amazon, Google, and Shopify). Your card details are never stored on our servers. All transactions are protected with 256-bit SSL encryption.",
  },
  {
    question: "Can I cancel anytime?",
    answer:
      "Yes, you can cancel your subscription at any time with immediate effect. There are no long-term commitments or cancellation fees.",
  },
  {
    question: "What fasting schedules do you support?",
    answer:
      "We support all popular fasting schedules including 16:8, 18:6, and 23:1 (OMAD). You can choose the plan that fits your lifestyle and goals, and switch between them as needed.",
  },
  {
    question: "Do I need any special equipment?",
    answer:
      "No special equipment needed! FastingHero works as a standalone app. However, if you want to track additional metrics like ketones or glucose, the app can integrate with compatible devices.",
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
