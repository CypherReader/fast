import { useState, useEffect } from 'react';
import { Heart, Zap, Brain, Clock, CheckCircle, Users, Flame, AlertCircle, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Textarea } from '@/components/ui/textarea';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { useSOSFlare, useResolveSOS, useHypeResponses } from '@/hooks/use-sos-flare';
import { motion, AnimatePresence } from 'framer-motion';

export const CravingHelpButton = () => {
    const [isOpen, setIsOpen] = useState(false);
    const [cravingDesc, setCravingDesc] = useState('');
    const { sendSOSFlare, sosData, isSending, reset } = useSOSFlare();
    const resolveSOS = useResolveSOS();

    // Poll for hype responses when we have an active SOS
    const { data: hypeResponses } = useHypeResponses(sosData?.sos_id);

    // Calculate live hype count
    const liveHypeCount = hypeResponses?.length || sosData?.hype_count || 0;

    const handleGetHelp = () => {
        const description = cravingDesc.trim() || "I'm hungry";
        sendSOSFlare(description);
    };

    const handleResolved = (survived: boolean) => {
        if (sosData?.sos_id) {
            resolveSOS.mutate({ sosId: sosData.sos_id, survived });
        }
        reset();
        setCravingDesc('');
        setIsOpen(false);
    };

    // Reset state when modal closes
    useEffect(() => {
        if (!isOpen) {
            // Don't reset immediately to allow animation
            const timer = setTimeout(() => {
                reset();
                setCravingDesc('');
            }, 300);
            return () => clearTimeout(timer);
        }
    }, [isOpen, reset]);

    const aiResponse = sosData?.ai_response;
    const isCooldown = sosData?.status === 'cooldown';

    return (
        <>
            {/* Floating "I'm Hungry" Button */}
            <Button
                onClick={() => setIsOpen(true)}
                variant="destructive"
                size="lg"
                className="fixed bottom-20 md:bottom-6 right-6 h-14 px-6 shadow-lg z-50 animate-pulse hover:animate-none"
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
                            className="w-full max-w-2xl max-h-[90vh] overflow-y-auto"
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
                                    {/* Cooldown Warning */}
                                    {isCooldown && (
                                        <Alert>
                                            <AlertCircle className="h-4 w-4" />
                                            <AlertDescription>
                                                You can only send a tribe SOS once every 24 hours.
                                                AI help is still available!
                                            </AlertDescription>
                                        </Alert>
                                    )}

                                    {!aiResponse && !isSending ? (
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
                                                    disabled={isSending}
                                                    className="flex-1"
                                                >
                                                    {isSending ? (
                                                        <>
                                                            <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                                                            Sending SOS...
                                                        </>
                                                    ) : (
                                                        'Help Me Fight This'
                                                    )}
                                                </Button>
                                                <Button
                                                    onClick={() => setIsOpen(false)}
                                                    variant="outline"
                                                >
                                                    Close
                                                </Button>
                                            </div>
                                        </>
                                    ) : isSending ? (
                                        <div className="flex flex-col items-center justify-center py-8 space-y-4">
                                            <Loader2 className="w-12 h-12 animate-spin text-primary" />
                                            <p className="text-muted-foreground">
                                                Cortex is analyzing your situation...
                                            </p>
                                        </div>
                                    ) : aiResponse && (
                                        <div className="space-y-4">
                                            {/* Tribe Notification Status */}
                                            {sosData?.tribe_notified && (
                                                <motion.div
                                                    initial={{ opacity: 0, y: -10 }}
                                                    animate={{ opacity: 1, y: 0 }}
                                                    className="p-4 bg-purple-500/10 border border-purple-500/30 rounded-lg"
                                                >
                                                    <div className="flex items-center gap-3 mb-2">
                                                        <Users className="w-5 h-5 text-purple-400" />
                                                        <h4 className="font-semibold text-purple-400">
                                                            🆘 SOS Sent to {sosData.allies_notified} Allies!
                                                        </h4>
                                                    </div>
                                                    <p className="text-sm text-muted-foreground mb-3">
                                                        Your tribe has been notified. Help is on the way!
                                                    </p>

                                                    {/* Live Hype Counter */}
                                                    <motion.div
                                                        key={liveHypeCount}
                                                        initial={{ scale: 1.2 }}
                                                        animate={{ scale: 1 }}
                                                        className="flex items-center gap-2"
                                                    >
                                                        <Flame className="w-5 h-5 text-orange-500" />
                                                        <span className="text-lg font-bold text-orange-500">
                                                            {liveHypeCount} Hype Received
                                                        </span>
                                                        {liveHypeCount > 0 && (
                                                            <span className="text-xs text-muted-foreground">
                                                                Your tribe has your back! 🔥
                                                            </span>
                                                        )}
                                                    </motion.div>

                                                    {/* Show recent hype messages */}
                                                    {hypeResponses && hypeResponses.length > 0 && (
                                                        <div className="mt-3 space-y-2">
                                                            {hypeResponses.slice(0, 3).map((hype) => (
                                                                <div
                                                                    key={hype.id}
                                                                    className="flex items-center gap-2 text-sm bg-purple-500/5 rounded-lg p-2"
                                                                >
                                                                    <span>{hype.emoji}</span>
                                                                    <span className="font-medium">
                                                                        {hype.from_name || 'An ally'}
                                                                    </span>
                                                                    {hype.message && (
                                                                        <span className="text-muted-foreground">
                                                                            - "{hype.message}"
                                                                        </span>
                                                                    )}
                                                                </div>
                                                            ))}
                                                        </div>
                                                    )}
                                                </motion.div>
                                            )}

                                            {/* Immediate Action */}
                                            <div className="p-4 bg-destructive/10 border border-destructive/20 rounded-lg">
                                                <div className="flex items-start gap-2">
                                                    <Zap className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
                                                    <div>
                                                        <h4 className="font-semibold text-destructive mb-1">RIGHT NOW:</h4>
                                                        <p className="text-sm">{aiResponse.immediate_action}</p>
                                                    </div>
                                                </div>
                                            </div>

                                            {/* Distraction */}
                                            <div className="p-4 bg-primary/10 border border-primary/20 rounded-lg">
                                                <div className="flex items-start gap-2">
                                                    <Brain className="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
                                                    <div>
                                                        <h4 className="font-semibold text-primary mb-1">5-Minute Distraction:</h4>
                                                        <p className="text-sm">{aiResponse.distraction_idea}</p>
                                                    </div>
                                                </div>
                                            </div>

                                            {/* Science */}
                                            <div className="p-4 bg-secondary/10 border border-secondary/20 rounded-lg">
                                                <div className="flex items-start gap-2">
                                                    <CheckCircle className="w-5 h-5 text-secondary flex-shrink-0 mt-0.5" />
                                                    <div>
                                                        <h4 className="font-semibold text-secondary mb-1">What's Happening:</h4>
                                                        <p className="text-sm">{aiResponse.biological_fact}</p>
                                                    </div>
                                                </div>
                                            </div>

                                            {/* Motivation */}
                                            <div className="p-4 bg-accent/10 border border-accent/20 rounded-lg">
                                                <div className="flex items-center justify-center">
                                                    <p className="text-center font-semibold italic text-lg text-accent">
                                                        "{aiResponse.motivation}"
                                                    </p>
                                                </div>
                                            </div>

                                            {/* Time Remaining */}
                                            {aiResponse.time_remaining && (
                                                <div className="flex items-center justify-center gap-2 text-sm text-muted-foreground">
                                                    <Clock className="w-4 h-4" />
                                                    {aiResponse.time_remaining}
                                                </div>
                                            )}

                                            {/* Support Strategies */}
                                            {aiResponse.support_strategies && aiResponse.support_strategies.length > 0 && (
                                                <div className="pt-4 border-t">
                                                    <h4 className="text-sm font-medium mb-2">Quick Strategies:</h4>
                                                    <div className="flex flex-wrap gap-2">
                                                        {aiResponse.support_strategies.map((strategy, index) => (
                                                            <span
                                                                key={index}
                                                                className="px-3 py-1 bg-muted rounded-full text-xs"
                                                            >
                                                                {strategy}
                                                            </span>
                                                        ))}
                                                    </div>
                                                </div>
                                            )}

                                            {/* Actions */}
                                            <div className="flex gap-2 pt-2">
                                                <Button
                                                    onClick={() => handleResolved(true)}
                                                    variant="default"
                                                    className="flex-1 bg-green-600 hover:bg-green-700"
                                                >
                                                    <CheckCircle className="w-4 h-4 mr-2" />
                                                    I Survived! 💪
                                                </Button>
                                                <Button
                                                    onClick={() => {
                                                        reset();
                                                        setCravingDesc('');
                                                    }}
                                                    variant="secondary"
                                                    className="flex-1"
                                                >
                                                    Still Struggling
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
