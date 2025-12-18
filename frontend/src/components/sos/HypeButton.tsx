import { useState } from 'react';
import { Flame, Zap, Send, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { useSendHype } from '@/hooks/use-sos-flare';
import { motion, AnimatePresence } from 'framer-motion';

interface HypeButtonProps {
    sosId: string;
    userName?: string;
    onHypeSent?: () => void;
}

const QUICK_EMOJIS = [
    { emoji: '🔥', label: 'Fire' },
    { emoji: '⚡', label: 'Bolt' },
    { emoji: '💪', label: 'Strong' },
    { emoji: '🙌', label: 'Praise' },
    { emoji: '❤️', label: 'Love' },
];

export const HypeButton = ({ sosId, userName, onHypeSent }: HypeButtonProps) => {
    const { mutate: sendHype, isPending, isSuccess } = useSendHype();
    const [selectedEmoji, setSelectedEmoji] = useState<string | null>(null);
    const [message, setMessage] = useState('');
    const [showMessageInput, setShowMessageInput] = useState(false);

    const handleSendHype = (emoji: string) => {
        setSelectedEmoji(emoji);
        sendHype(
            { sosId, emoji, message: message.trim() || undefined },
            {
                onSuccess: () => {
                    onHypeSent?.();
                    setMessage('');
                    setShowMessageInput(false);
                },
            }
        );
    };

    if (isSuccess) {
        return (
            <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                className="text-center p-4 bg-green-500/10 border border-green-500/30 rounded-lg"
            >
                <span className="text-3xl">{selectedEmoji}</span>
                <p className="text-sm text-green-500 font-medium mt-2">
                    Hype sent! You're an amazing ally.
                </p>
            </motion.div>
        );
    }

    return (
        <div className="space-y-4">
            <div className="text-center">
                <p className="text-sm text-muted-foreground mb-2">
                    {userName ? `${userName} needs your support!` : 'A tribe member needs your support!'}
                </p>
            </div>

            {/* Quick Emoji Buttons */}
            <div className="flex justify-center gap-2 flex-wrap">
                {QUICK_EMOJIS.map(({ emoji, label }) => (
                    <Button
                        key={emoji}
                        onClick={() => handleSendHype(emoji)}
                        disabled={isPending}
                        variant="outline"
                        size="lg"
                        className="w-14 h-14 text-2xl hover:scale-110 transition-transform"
                        title={label}
                    >
                        {isPending && selectedEmoji === emoji ? (
                            <Loader2 className="w-5 h-5 animate-spin" />
                        ) : (
                            emoji
                        )}
                    </Button>
                ))}
            </div>

            {/* Toggle message input */}
            <Button
                onClick={() => setShowMessageInput(!showMessageInput)}
                variant="ghost"
                size="sm"
                className="w-full"
            >
                {showMessageInput ? 'Hide message' : 'Add a message'}
            </Button>

            {/* Message Input */}
            <AnimatePresence>
                {showMessageInput && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="overflow-hidden"
                    >
                        <div className="flex gap-2">
                            <Input
                                value={message}
                                onChange={(e) => setMessage(e.target.value)}
                                placeholder="You've got this!"
                                maxLength={100}
                                className="flex-1"
                            />
                            <Button
                                onClick={() => handleSendHype('🔥')}
                                disabled={isPending || !message.trim()}
                            >
                                <Send className="w-4 h-4" />
                            </Button>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Pre-written quick messages */}
            {showMessageInput && (
                <div className="flex flex-wrap gap-2 justify-center">
                    {[
                        "You've got this!",
                        "Stay strong!",
                        "Almost there!",
                        "We believe in you!",
                    ].map((msg) => (
                        <Button
                            key={msg}
                            onClick={() => setMessage(msg)}
                            variant="ghost"
                            size="sm"
                            className="text-xs"
                        >
                            {msg}
                        </Button>
                    ))}
                </div>
            )}
        </div>
    );
};

export default HypeButton;
