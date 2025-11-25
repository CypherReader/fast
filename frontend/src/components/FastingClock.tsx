import React from 'react';


interface FastingClockProps {
    elapsedSeconds: number;
    goalHours: number;
    isFasting: boolean;
    stage: string;
}

const FastingClock: React.FC<FastingClockProps> = ({ elapsedSeconds, goalHours, isFasting, stage }) => {
    const targetSeconds = goalHours * 3600;
    const progressPercentage = Math.min((elapsedSeconds / targetSeconds) * 100, 100);

    const hours = Math.floor(elapsedSeconds / 3600);
    const minutes = Math.floor((elapsedSeconds % 3600) / 60);
    const seconds = elapsedSeconds % 60;

    const formatTime = (val: number) => String(val).padStart(2, '0');

    // Circle configuration
    const radius = 120;
    const circumference = 2 * Math.PI * radius;
    const strokeDashoffset = circumference - (progressPercentage / 100) * circumference;

    return (
        <div className="relative flex items-center justify-center w-80 h-80">
            {/* Pulsing Background Glow */}
            {isFasting && (
                <div className="absolute inset-0 rounded-full bg-emerald-500/20 blur-3xl animate-pulse" />
            )}

            {/* SVG Ring */}
            <div className="relative w-72 h-72">
                <svg className="w-full h-full transform -rotate-90 drop-shadow-2xl">
                    {/* Background Circle */}
                    <circle
                        cx="144"
                        cy="144"
                        r={radius}
                        stroke="currentColor"
                        strokeWidth="12"
                        fill="transparent"
                        className="text-slate-800"
                    />

                    {/* Progress Circle */}
                    <circle
                        cx="144"
                        cy="144"
                        r={radius}
                        stroke="url(#gradient)"
                        strokeWidth="12"
                        fill="transparent"
                        strokeDasharray={circumference}
                        strokeDashoffset={strokeDashoffset}
                        strokeLinecap="round"
                        className="transition-all duration-1000 ease-in-out"
                    />

                    {/* Gradient Definition */}
                    <defs>
                        <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                            <stop offset="0%" stopColor="#10b981" />
                            <stop offset="100%" stopColor="#34d399" />
                        </linearGradient>
                    </defs>
                </svg>

                {/* Inner Content */}
                <div className="absolute inset-0 flex flex-col items-center justify-center text-center">
                    {isFasting ? (
                        <>
                            <div className="text-5xl font-mono font-bold tracking-tighter text-white drop-shadow-lg tabular-nums">
                                {formatTime(hours)}:{formatTime(minutes)}
                                <span className="text-2xl text-emerald-400/80 ml-1">:{formatTime(seconds)}</span>
                            </div>
                            <div className="text-sm font-medium text-emerald-400 mt-2 uppercase tracking-widest">
                                {stage}
                            </div>
                            <div className="text-xs text-slate-400 mt-1">
                                Goal: {goalHours}h
                            </div>
                        </>
                    ) : (
                        <div className="flex flex-col items-center animate-float">
                            <span className="text-2xl font-bold text-slate-300">Ready to Fast?</span>
                            <span className="text-xs text-slate-500 mt-2">Select a plan below</span>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default FastingClock;
