import { motion, useInView } from "framer-motion";
import { ArrowRight, TrendingDown, Heart, Zap } from "lucide-react";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";
import { MagneticButton } from "./MagneticButton";

const transformations = [
    {
        name: "Sarah M.",
        image: "SM",
        before: "185 lbs",
        after: "152 lbs",
        lost: "33 lbs",
        timeframe: "90 days",
        quote: "Cortex knew exactly when I needed support. The AI coaching made all the difference.",
        icon: Heart,
    },
    {
        name: "Mike R.",
        image: "MR",
        before: "220 lbs",
        after: "195 lbs",
        lost: "25 lbs",
        timeframe: "75 days",
        quote: "I've tried every fasting app. This is the only one that actually worked.",
        icon: Zap,
    },
    {
        name: "Emma L.",
        image: "EL",
        before: "165 lbs",
        after: "145 lbs",
        lost: "20 lbs",
        timeframe: "60 days",
        quote: "The tribe kept me accountable. I wasn't alone in this journey.",
        icon: TrendingDown,
    },
];

export const TransformationShowcase = () => {
    const ref = useRef(null);
    const isInView = useInView(ref, { once: true, margin: "-100px" });

    return (
        <section className="py-20 md:py-32 bg-background relative overflow-hidden" ref={ref}>
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
                            <span className="text-gradient-hero">Real People.</span> Real Results.{" "}
                            <span className="text-primary">Real Stories.</span>
                        </h2>
                        <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
                            These aren't stock photos. These are real FastingHero members who transformed their lives.
                        </p>
                    </motion.div>

                    {/* Transformation cards */}
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12">
                        {transformations.map((person, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0, y: 30 }}
                                animate={isInView ? { opacity: 1, y: 0 } : {}}
                                transition={{ duration: 0.5, delay: i * 0.15 }}
                                whileHover={{ y: -8, transition: { duration: 0.2 } }}
                                className="bg-card/80 backdrop-blur rounded-2xl p-6 border border-border hover:border-secondary/50 transition-all"
                            >
                                {/* Avatar and stats */}
                                <div className="flex items-center gap-4 mb-4">
                                    <div className="w-16 h-16 rounded-full bg-gradient-to-br from-primary/80 to-secondary/80 border-2 border-background flex items-center justify-center text-xl font-bold text-primary-foreground">
                                        {person.image}
                                    </div>
                                    <div>
                                        <h3 className="font-semibold text-foreground text-lg">{person.name}</h3>
                                        <p className="text-sm text-muted-foreground">{person.timeframe}</p>
                                    </div>
                                </div>

                                {/* Weight transformation */}
                                <div className="flex items-center justify-between mb-4 p-4 bg-primary/5 rounded-xl border border-primary/20">
                                    <div className="text-center">
                                        <div className="text-xs text-muted-foreground mb-1">Before</div>
                                        <div className="text-lg font-bold text-foreground">{person.before}</div>
                                    </div>
                                    <ArrowRight className="w-5 h-5 text-primary" />
                                    <div className="text-center">
                                        <div className="text-xs text-muted-foreground mb-1">After</div>
                                        <div className="text-lg font-bold text-primary">{person.after}</div>
                                    </div>
                                </div>

                                {/* Lost badge */}
                                <div className="flex items-center justify-center gap-2 mb-4 py-2 px-4 bg-secondary/10 rounded-full border border-secondary/30">
                                    <person.icon className="w-4 h-4 text-secondary" />
                                    <span className="text-sm font-semibold text-secondary">Lost {person.lost}</span>
                                </div>

                                {/* Quote */}
                                <p className="text-sm text-muted-foreground italic leading-relaxed">"{person.quote}"</p>
                            </motion.div>
                        ))}
                    </div>

                    {/* CTA */}
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={isInView ? { opacity: 1, y: 0 } : {}}
                        transition={{ duration: 0.6, delay: 0.5 }}
                        className="text-center"
                    >
                        <MagneticButton className="text-lg">
                            Start Your Story
                            <ArrowRight className="w-5 h-5 ml-2" />
                        </MagneticButton>
                        <p className="text-sm text-muted-foreground mt-4">
                            Join 12,847 people on their transformation journey
                        </p>
                    </motion.div>
                </div>
            </div>
        </section>
    );
};
