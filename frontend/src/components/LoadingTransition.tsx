import { useEffect, useRef } from 'react';
import { motion } from 'framer-motion';

interface LoadingTransitionProps {
    onComplete: () => void;
}

export const LoadingTransition = ({ onComplete }: LoadingTransitionProps) => {
    const videoRef = useRef<HTMLVideoElement>(null);

    useEffect(() => {
        const video = videoRef.current;
        if (!video) return;

        const handleEnded = () => {
            onComplete();
        };

        const handleError = () => {
            console.error('Failed to load transition video');
            // Gracefully navigate to dashboard even if video fails
            setTimeout(onComplete, 500);
        };

        video.addEventListener('ended', handleEnded);
        video.addEventListener('error', handleError);

        // Start playing with increased speed
        video.playbackRate = 1.25; // Speed up by 25%
        video.play().catch((error) => {
            console.error('Failed to play transition video:', error);
            setTimeout(onComplete, 500);
        });

        return () => {
            video.removeEventListener('ended', handleEnded);
            video.removeEventListener('error', handleError);
        };
    }, [onComplete]);

    return (
        <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3 }}
            className="fixed inset-0 z-50 flex items-center justify-center bg-background overflow-hidden"
        >
            <video
                ref={videoRef}
                muted
                playsInline
                className="w-full h-full"
                style={{
                    objectFit: 'cover',
                    transform: 'scale(1.5)',
                    transformOrigin: 'center center'
                }}
                src="/Logo_Animation_Fasting_Hero.mp4"
            />
        </motion.div>
    );
};
