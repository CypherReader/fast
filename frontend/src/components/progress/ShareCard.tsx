import { useRef, useState } from 'react';
import { motion } from 'framer-motion';
import { Share2, Download, Twitter, Facebook, Check } from 'lucide-react';
import { Button } from '@/components/ui/button';
import html2canvas from 'html2canvas';

interface ShareCardProps {
    userName: string;
    streak: number;
    fastingPlan: string;
    totalFasts: number;
    totalHours: number;
}

export const ShareCard = ({
    userName,
    streak,
    fastingPlan,
    totalFasts,
    totalHours,
}: ShareCardProps) => {
    const cardRef = useRef<HTMLDivElement>(null);
    const [isSharing, setIsSharing] = useState(false);
    const [showSuccess, setShowSuccess] = useState(false);

    const generateImage = async () => {
        if (!cardRef.current) return null;

        const canvas = await html2canvas(cardRef.current, {
            backgroundColor: null,
            scale: 2,
            logging: false,
        });

        return canvas.toDataURL('image/png');
    };

    const handleDownload = async () => {
        setIsSharing(true);
        try {
            const dataUrl = await generateImage();
            if (!dataUrl) return;

            const link = document.createElement('a');
            link.download = `fastinghero-progress-${Date.now()}.png`;
            link.href = dataUrl;
            link.click();

            setShowSuccess(true);
            setTimeout(() => setShowSuccess(false), 2000);
        } finally {
            setIsSharing(false);
        }
    };

    const handleShare = async () => {
        setIsSharing(true);
        try {
            const dataUrl = await generateImage();
            if (!dataUrl) return;

            // Convert data URL to blob
            const response = await fetch(dataUrl);
            const blob = await response.blob();
            const file = new File([blob], 'fastinghero-progress.png', { type: 'image/png' });

            if (navigator.share && navigator.canShare({ files: [file] })) {
                await navigator.share({
                    files: [file],
                    title: 'My FastingHero Progress',
                    text: `🔥 ${streak} day streak on FastingHero! ${totalFasts} fasts completed, ${totalHours}+ hours fasted. Join me!`,
                });
            } else {
                // Fallback to download
                handleDownload();
            }
        } catch (error) {
            console.error('Share failed:', error);
        } finally {
            setIsSharing(false);
        }
    };

    const shareToTwitter = () => {
        const text = encodeURIComponent(`🔥 ${streak} day streak on @FastingHero! ${totalFasts} fasts completed, ${totalHours}+ hours fasted. Join the fasting revolution! #FastingHero #IntermittentFasting`);
        window.open(`https://twitter.com/intent/tweet?text=${text}`, '_blank');
    };

    const shareToFacebook = () => {
        // Facebook share URL (image sharing requires backend)
        const url = encodeURIComponent('https://fastinghero.com');
        window.open(`https://www.facebook.com/sharer/sharer.php?u=${url}`, '_blank');
    };

    return (
        <div className="space-y-4">
            {/* Shareable Card */}
            <motion.div
                ref={cardRef}
                className="relative overflow-hidden rounded-2xl p-6"
                style={{
                    background: 'linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%)',
                }}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
            >
                {/* Decorative elements */}
                <div className="absolute -top-10 -right-10 w-40 h-40 bg-white/10 rounded-full blur-3xl" />
                <div className="absolute -bottom-10 -left-10 w-32 h-32 bg-white/10 rounded-full blur-3xl" />

                {/* Content */}
                <div className="relative z-10">
                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                            <div className="w-12 h-12 rounded-full bg-white/20 flex items-center justify-center text-white text-xl font-bold">
                                {userName.charAt(0).toUpperCase()}
                            </div>
                            <div>
                                <p className="text-white font-semibold">{userName}</p>
                                <p className="text-white/70 text-sm">{fastingPlan} Faster</p>
                            </div>
                        </div>
                        <img src="/fasthero.png" alt="FastingHero" className="w-8 h-8 rounded-lg opacity-80" />
                    </div>

                    <div className="text-center py-4">
                        <motion.p
                            className="text-5xl font-bold text-white"
                            initial={{ scale: 0.5 }}
                            animate={{ scale: 1 }}
                            transition={{ type: "spring", damping: 10 }}
                        >
                            {streak} 🔥
                        </motion.p>
                        <p className="text-white/80 text-lg mt-1">Day Streak</p>
                    </div>

                    <div className="grid grid-cols-2 gap-4 mt-4">
                        <div className="bg-white/10 rounded-xl p-3 text-center backdrop-blur-sm">
                            <p className="text-2xl font-bold text-white">{totalFasts}</p>
                            <p className="text-white/70 text-xs">Fasts Completed</p>
                        </div>
                        <div className="bg-white/10 rounded-xl p-3 text-center backdrop-blur-sm">
                            <p className="text-2xl font-bold text-white">{totalHours}+</p>
                            <p className="text-white/70 text-xs">Hours Fasted</p>
                        </div>
                    </div>
                </div>
            </motion.div>

            {/* Share Actions */}
            <div className="flex items-center gap-2">
                <Button
                    onClick={handleDownload}
                    variant="outline"
                    size="sm"
                    className="flex-1 gap-2"
                    disabled={isSharing}
                >
                    {showSuccess ? <Check className="w-4 h-4 text-green-500" /> : <Download className="w-4 h-4" />}
                    {showSuccess ? 'Saved!' : 'Download'}
                </Button>

                <Button
                    onClick={handleShare}
                    size="sm"
                    className="flex-1 gap-2 bg-gradient-to-r from-indigo-500 to-purple-500 hover:from-indigo-600 hover:to-purple-600"
                    disabled={isSharing}
                >
                    <Share2 className="w-4 h-4" />
                    Share
                </Button>

                <Button
                    onClick={shareToTwitter}
                    variant="outline"
                    size="icon"
                    className="shrink-0"
                >
                    <Twitter className="w-4 h-4" />
                </Button>

                <Button
                    onClick={shareToFacebook}
                    variant="outline"
                    size="icon"
                    className="shrink-0"
                >
                    <Facebook className="w-4 h-4" />
                </Button>
            </div>
        </div>
    );
};

export default ShareCard;
