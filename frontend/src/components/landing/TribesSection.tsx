import { motion, useInView } from "framer-motion";
import { Users, Heart, Trophy, MessageCircle, Zap, ArrowRight } from "lucide-react";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";
import { MagneticButton } from "./MagneticButton";

const tribeFeatures = [
    {
        icon: Users,
        title: "Fast Together",
        description: "Join or create a tribe of 5-20 fasters. Share your journey with people who get it.",
        color: "text-primary",
    },
    {
        icon: MessageCircle,
        title: "Live Group Chat",
        description: "Need support at 2am? Your tribe is there. Share wins, struggles, and motivation 24/7.",
        color: "text-secondary",
    },
    {
        icon: Trophy,
        title: "Team Challenges",
        description: "Compete in weekly tribe challenges. Earn badges together. Build streaks as a team.",
        color: "text-accent",
    },
    {
        icon: Heart,
        title: "Accountability Partners",
        description: "Get paired with a tribe member who has similar goals. Check in daily, stay committed.",
        color: "text-primary",
    },
];

const stats = [
    { number: "87%", label: "Higher success rate in tribes" },
    { number: "3.2x", label: "Longer streaks with accountability" },
    { number: "12,000+", label: "Active tribe members" },
];

export const TribesSection = () => {
    const ref = useRef(null);
    const isInView = useInView(ref, { once: true, margin: "-100px" });

    return (
        <section className="py-20 md:py-32 bg-slate-900 dark:bg-gradient-dark relative overflow-hidden" ref={ref}>
            <AnimatedSectionBackground variant="accent" showOrbs showGrid showParticles />

            <div className="container px-4 relative z-10">
                <div className="max-w-6xl mx-auto">
                    {/* Section header */}
                    <motion.div
                        initial={{ opacity: 0, y: 30 }}
                        animate={isInView ? { opacity: 1, y: 0 } : {}}
                        transition={{ duration: 0.6 }}
                        className="text-center mb-16"
                    >
                        <div className="inline-flex items-center gap-2 mb-4 px-4 py-2 bg-primary/20 border border-primary/40 rounded-full">
                            <Users className="w-4 h-4 text-primary" />
                            <span className="text-sm font-semibold text-primary">NEW: Tribes Feature</span>
                        </div>
                        <h2 className="text-4xl md:text-5xl lg:text-6xl font-bold text-white mb-6">
                            <span className="text-gradient-hero">You Don't Have to Fast Alone</span>
                        </h2>
                        <p className="text-lg md:text-xl text-gray-300 max-w-3xl mx-auto">
                            Join a tribe and fast as a team. Because the hardest part of fasting isn't hunger—it's doing it alone.
                        </p>
                    </motion.div>

                    {/* Stats bar */}
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={isInView ? { opacity: 1, y: 0 } : {}}
                        transition={{ duration: 0.6, delay: 0.2 }}
                        className="grid grid-cols-3 gap-4 mb-16 p-6 bg-card/30 backdrop-blur rounded-2xl border border-border"
                    >
                        {stats.map((stat, i) => (
                            <div key={i} className="text-center">
                                <div className="text-3xl md:text-4xl font-bold text-primary mb-1">{stat.number}</div>
                                <div className="text-xs md:text-sm text-gray-400">{stat.label}</div>
                            </div>
                        ))}
                    </motion.div>

                    {/* Features grid */}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
                        {tribeFeatures.map((feature, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0, y: 30 }}
                                animate={isInView ? { opacity: 1, y: 0 } : {}}
                                transition={{ duration: 0.5, delay: 0.3 + i * 0.1 }}
                                whileHover={{ y: -8, transition: { duration: 0.2 } }}
                                className="bg-card/50 backdrop-blur rounded-xl p-6 border border-border hover:border-secondary/50 transition-all"
                            >
                                <div className="flex items-start gap-4">
                                    <div className={`w-12 h-12 rounded-xl bg-primary/10 border border-primary/30 flex items-center justify-center flex-shrink-0`}>
                                        <feature.icon className={`w-6 h-6 ${feature.color}`} />
                                    </div>
                                    <div>
                                        <h3 className="text-lg font-semibold text-white mb-2">{feature.title}</h3>
                                        <p className="text-gray-300 leading-relaxed">{feature.description}</p>
                                    </div>
                                </div>
                            </motion.div>
                        ))}
                    </div>

                    {/* Social proof testimonial */}
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={isInView ? { opacity: 1, y: 0 } : {}}
                        transition={{ duration: 0.6, delay: 0.7 }}
                        className="bg-gradient-to-br from-primary/10 to-secondary/10 border border-primary/30 rounded-2xl p-8 mb-8"
                    >
                        <div className="flex items-center gap-3 mb-4">
                            <div className="w-12 h-12 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white font-bold">
                                JW
                            </div>
                            <div>
                                <div className="font-semibold text-white">Jessica W.</div>
                                <div className="text-sm text-gray-400">Lost 28 lbs in 90 days</div>
                            </div>
                        </div>
                        <blockquote className="text-lg text-gray-200 italic leading-relaxed">
                            "My tribe kept me going when I wanted to quit. On day 12, I was ready to give up. One message from my tribe and I pushed through. Now I'm on day 87 and down 28 pounds. I couldn't have done this alone."
                        </blockquote>
                    </motion.div>

                    {/* CTA */}
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={isInView ? { opacity: 1, y: 0 } : {}}
                        transition={{ duration: 0.6, delay: 0.8 }}
                        className="text-center"
                    >
                        <MagneticButton className="text-lg" onClick={() => window.location.href = '/register'}>
                            <Users className="w-5 h-5" />
                            Join a Tribe Today
                            <ArrowRight className="w-5 h-5 ml-2" />
                        </MagneticButton>
                        <p className="text-sm text-gray-400 mt-4">
                            Find your tribe or create one with friends
                        </p>
                    </motion.div>
                </div>
            </div>
        </section>
    );
};
