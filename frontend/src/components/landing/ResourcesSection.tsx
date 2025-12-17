import { motion, useInView } from "framer-motion";
import { PlayCircle, ChefHat, BookOpen, Download, ArrowRight, Star } from "lucide-react";
import { useRef } from "react";
import { AnimatedSectionBackground } from "./AnimatedSectionBackground";

const videos = [
    {
        title: "Fasting 101: Complete Beginner's Guide",
        duration: "12:34",
        views: "127K",
        thumbnail: "🎥",
    },
    {
        title: "Breaking Your Fast the Right Way",
        duration: "8:45",
        views: "89K",
        thumbnail: "🍽️",
    },
    {
        title: "Common Fasting Mistakes to Avoid",
        duration: "15:20",
        views: "156K",
        thumbnail: "⚠️",
    },
];

const recipes = [
    {
        name: "Keto Bulletproof Coffee",
        type: "Breakfast",
        time: "5 min",
        calories: "230 cal",
        icon: "☕",
        rating: "4.9",
    },
    {
        name: "Avocado Egg Salad",
        type: "Lunch",
        time: "10 min",
        calories: "340 cal",
        icon: "🥑",
        rating: "4.8",
    },
    {
        name: "Garlic Butter Salmon",
        type: "Dinner",
        time: "20 min",
        calories: "420 cal",
        icon: "🐟",
        rating: "5.0",
    },
    {
        name: "Keto Chocolate Fat Bombs",
        type: "Snack",
        time: "15 min",
        calories: "150 cal",
        icon: "🍫",
        rating: "4.9",
    },
];

const resources = [
    {
        icon: BookOpen,
        title: "Fasting Guide Library",
        description: "100+ articles on fasting science, protocols, and tips",
        color: "text-primary",
    },
    {
        icon: ChefHat,
        title: "500+ Keto Recipes",
        description: "Break-fast meals, meal prep, and quick recipes",
        color: "text-secondary",
    },
    {
        icon: PlayCircle,
        title: "Expert Video Library",
        description: "Learn from doctors, nutritionists, and fasting experts",
        color: "text-accent",
    },
];

export const ResourcesSection = () => {
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
                            <span className="text-gradient-hero">Everything You Need</span> to Succeed
                        </h2>
                        <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
                            From expert videos to keto recipes, get access to our complete resource library
                        </p>
                    </motion.div>

                    {/* Resource types */}
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-16">
                        {resources.map((resource, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0, y: 20 }}
                                animate={isInView ? { opacity: 1, y: 0 } : {}}
                                transition={{ duration: 0.5, delay: i * 0.1 }}
                                className="bg-card/80 backdrop-blur rounded-xl p-6 border border-border hover:border-secondary/50 transition-all text-center"
                            >
                                <resource.icon className={`w-12 h-12 ${resource.color} mx-auto mb-4`} />
                                <h3 className="text-lg font-semibold text-foreground mb-2">{resource.title}</h3>
                                <p className="text-sm text-muted-foreground">{resource.description}</p>
                            </motion.div>
                        ))}
                    </div>

                    {/* Videos showcase */}
                    <motion.div
                        initial={{ opacity: 0, y: 30 }}
                        animate={isInView ? { opacity: 1, y: 0 } : {}}
                        transition={{ duration: 0.6, delay: 0.3 }}
                        className="mb-16"
                    >
                        <div className="flex items-center justify-between mb-6">
                            <h3 className="text-2xl md:text-3xl font-bold text-foreground">📺 Popular Videos</h3>
                            <button className="text-sm text-primary hover:underline flex items-center gap-1">
                                View All <ArrowRight className="w-4 h-4" />
                            </button>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                            {videos.map((video, i) => (
                                <motion.div
                                    key={i}
                                    initial={{ opacity: 0, y: 20 }}
                                    animate={isInView ? { opacity: 1, y: 0 } : {}}
                                    transition={{ duration: 0.5, delay: 0.4 + i * 0.1 }}
                                    whileHover={{ y: -4, transition: { duration: 0.2 } }}
                                    className="bg-card/50 backdrop-blur rounded-xl overflow-hidden border border-border hover:border-primary/50 transition-all cursor-pointer group"
                                >
                                    {/* Thumbnail */}
                                    <div className="relative aspect-video bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center">
                                        <span className="text-6xl">{video.thumbnail}</span>
                                        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                                            <PlayCircle className="w-16 h-16 text-white" />
                                        </div>
                                        <div className="absolute bottom-2 right-2 px-2 py-1 bg-black/80 rounded text-xs text-white">
                                            {video.duration}
                                        </div>
                                    </div>

                                    {/* Content */}
                                    <div className="p-4">
                                        <h4 className="font-semibold text-foreground mb-2 line-clamp-2">{video.title}</h4>
                                        <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                            <PlayCircle className="w-3 h-3" />
                                            <span>{video.views} views</span>
                                        </div>
                                    </div>
                                </motion.div>
                            ))}
                        </div>
                    </motion.div>

                    {/* Recipes showcase */}
                    <motion.div
                        initial={{ opacity: 0, y: 30 }}
                        animate={isInView ? { opacity: 1, y: 0 } : {}}
                        transition={{ duration: 0.6, delay: 0.6 }}
                    >
                        <div className="flex items-center justify-between mb-6">
                            <h3 className="text-2xl md:text-3xl font-bold text-foreground">🍳 Keto Recipe Collection</h3>
                            <button className="text-sm text-primary hover:underline flex items-center gap-1">
                                View All <ArrowRight className="w-4 h-4" />
                            </button>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                            {recipes.map((recipe, i) => (
                                <motion.div
                                    key={i}
                                    initial={{ opacity: 0, y: 20 }}
                                    animate={isInView ? { opacity: 1, y: 0 } : {}}
                                    transition={{ duration: 0.5, delay: 0.7 + i * 0.05 }}
                                    whileHover={{ y: -4, transition: { duration: 0.2 } }}
                                    className="bg-card/80 backdrop-blur rounded-xl p-5 border border-border hover:border-secondary/50 transition-all cursor-pointer"
                                >
                                    {/* Icon and rating */}
                                    <div className="flex items-start justify-between mb-3">
                                        <span className="text-4xl">{recipe.icon}</span>
                                        <div className="flex items-center gap-1 px-2 py-1 bg-primary/10 rounded-full">
                                            <Star className="w-3 h-3 text-primary fill-current" />
                                            <span className="text-xs font-semibold text-primary">{recipe.rating}</span>
                                        </div>
                                    </div>

                                    {/* Name and type */}
                                    <h4 className="font-semibold text-foreground mb-1">{recipe.name}</h4>
                                    <div className="text-xs text-muted-foreground mb-3">{recipe.type}</div>

                                    {/* Meta */}
                                    <div className="flex items-center justify-between text-xs text-muted-foreground">
                                        <span>⏱️ {recipe.time}</span>
                                        <span>🔥 {recipe.calories}</span>
                                    </div>

                                    {/* Download button */}
                                    <button className="mt-3 w-full py-2 bg-primary/10 hover:bg-primary/20 border border-primary/30 rounded-lg text-xs font-semibold text-primary transition-all flex items-center justify-center gap-1">
                                        <Download className="w-3 h-3" />
                                        Get Recipe
                                    </button>
                                </motion.div>
                            ))}
                        </div>
                    </motion.div>
                </div>
            </div>
        </section>
    );
};
