import { useState } from 'react';
import { Heart, Zap, Brain, Clock, CheckCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Textarea } from '@/components/ui/textarea';
import { useCravingHelp } from '@/hooks/use-craving-help';
import { motion, AnimatePresence } from 'framer-motion';

export const CravingHelpButton = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [cravingDesc, setCravingDesc] = useState('');
    const { getCravingHelp, cravingHelpData, isGettingHelp } = useCravingHelp();

    const handleGetHelp = () => {
        if (!cravingDesc.trim()) {
            setCravingDesc('I\'m hungry');
        }
        getCravingHelp(cravingDesc || 'I\'m hungry');
    };

    return (
        <>
            {/* Floating "I'm Hungry" Button */}
            <Button
                onClick={() => setIsOpen(true)}
                variant="destructive"
                size="lg"
                className="fixed bottom-6 right-6 h-14 px-6 shadow-lg z-50 animate-pulse hover:animate-none"
            >
                <Heart className="w-5 h-5 mr-2" />
                I'm Hungry!
            </Button>

            {/* Help Modal */}
            <AnimatePresence>
                {isOpen && (
                    <motion.div
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
                        onClick={() => setIsOpen(false)}
                    >
                        <motion.div
                            initial={{ scale: 0.9, y: 20 }}
                            animate={{ scale: 1, y: 0 }}
                            exit={{ scale: 0.9, y: 20 }}
                            onClick={(e) => e.stopPropagation()}
                            className="w-full max-w-2xl"
                        >
                            <Card>
                                <CardHeader>
                                    <CardTitle className="flex items-center gap-2">
                                        <Brain className="w-6 h-6 text-primary" />
                                        Cortex Emergency Support
                                    </CardTitle>
                                    <CardDescription>
                                        Let's beat this craving together. You're stronger than you think.
                                    </CardDescription>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    {!cravingHelpData ? (
                                        <>
                                            <div>
                                                <label className="text-sm font-medium mb-2 block">
                                                    What are you craving? (optional)
                                                </label>
                                                <Textarea
                                                    value={cravingDesc}
                                                    onChange={(e) => setCravingDesc(e.target.value)}
                                                    placeholder="e.g., pizza, chocolate, anything sweet..."
                                                    className="resize-none"
                                                    rows={2}
                                                />
                                            </div>
                                            <div className="flex gap-2">
                                                <Button
                                                    onClick={handleGetHelp}
                                                    disabled={isGettingHelp}
                                                    className="flex-1"
                                                >
                                                    {isGettingHelp ? 'Getting Help...' : 'Help Me Fight This'}
                                                </Button>
                                                <Button
                                                    onClick={() => setIsOpen(false)}
                                                    variant="outline"
                                                >
                                                    Close
                                                </Button>
                                            </div>
                                        </>
                                    ) : (
                                        <div className="space-y-4">
                                            {/* Immediate Action */}
                                            <div className="p-4 bg-destructive/10 border border-destructive/20 rounded-lg">
                                                <div className="flex items-start gap-2">
                                                    <Zap className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
                                                    <div>
                                                        <h4 className="font-semibold text-destructive mb-1">RIGHT NOW:</h4>
                                                        <p className="text-sm">{cravingHelpData.immediate_action}</p>
                                                    </div>
                                                </div>
                                            </div>

                                            {/* Distraction */}
                                            <div className="p-4 bg-primary/10 border border-primary/20 rounded-lg">
                                                <div className="flex items-start gap-2">
                                                    <Brain className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                                                    <div>
                                                        <h4 className="font-semibold text-primary mb-1">5-Minute Distraction:</h4>
                                                        <p className="text-sm">{cravingHelpData.distraction_idea}</p>
                                                    </div>
                                                </div>
                                            </div>

                                            {/* Science */}
                                            <div className="p-4 bg-secondary/10 border border-secondary/20 rounded-lg">
                                                <div className="flex items-start gap-2">
                                                    <CheckCircle className="w-5 h-5 text-secondary flex-shrink-0 mt-0.5" />
                                                    <div>
                                                        <h4 className="font-semibold text-secondary mb-1">What's Happening:</h4>
                                                        <p className="text-sm">{cravingHelpData.biological_fact}</p>
                                                    </div>
                                                </div>
                                            </div>

                                            {/* Motivation */}
                                            <div className="p-4 bg-accent/10 border border-accent/20 rounded-lg">
                                                <div className="flex items-center justify-center">
                                                    <p className="text-center font-semibold italic text-lg text-accent">
                                                        "{cravingHelpData.motivation}"
                                                    </p>
                                                </div>
                                            </div>

                                            {/* Time Remaining */}
                                            {cravingHelpData.time_remaining && (
                                                <div className="flex items-center justify-center gap-2 text-sm text-muted-foreground">
                                                    <Clock className="w-4 h-4" />
                                                    {cravingHelpData.time_remaining}
                                                </div>
                                            )}

                                            {/* Support Strategies */}
                                            <div className="pt-4 border-t">
                                                <h4 className="text-sm font-medium mb-2">Quick Strategies:</h4>
                                                <div className="flex flex-wrap gap-2">
                                                    {cravingHelpData.support_strategies.map((strategy, index) => (
                                                        <span
                                                            key={index}
                                                            className="px-3 py-1 bg-muted rounded-full text-xs"
                                                        >
                                                            {strategy}
                                                        </span>
                                                    ))}
                                                </div>
                                            </div>

                                            {/* Actions */}
                                            <div className="flex gap-2 pt-2">
                                                <Button
                                                    onClick={() => {
                                                        setCravingDesc('');
                                                        setIsOpen(false);
                                                    }}
                                                    variant="outline"
                                                    className="flex-1"
                                                >
                                                    I've Got This!
                                                </Button>
                                                <Button
                                                    onClick={() => {
                                                        setCravingDesc('');
                                                        getCravingHelp('Still struggling');
                                                    }}
                                                    variant="secondary"
                                                    className="flex-1"
                                                >
                                                    I Need More Help
                                                </Button>
                                            </div>
                                        </div>
                                    )}
                                </CardContent>
                            </Card>
                        </motion.div>
                    </motion.div>
                )}
            </AnimatePresence>
        </>
    );
};
