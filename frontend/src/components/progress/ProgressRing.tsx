import { motion } from 'framer-motion';
import { cn } from '@/lib/utils';

interface ProgressRingProps {
    value: number;
    max: number;
    size?: number;
    strokeWidth?: number;
    color: string;
    bgColor?: string;
    icon: React.ReactNode;
    label: string;
    unit: string;
}

export const ProgressRing = ({
    value,
    max,
    size = 80,
    strokeWidth = 6,
    color,
    bgColor = 'rgba(255,255,255,0.1)',
    icon,
    label,
    unit,
}: ProgressRingProps) => {
    const radius = (size - strokeWidth) / 2;
    const circumference = radius * 2 * Math.PI;
    const progress = Math.min(value / max, 1);
    const offset = circumference - progress * circumference;

    return (
        <div className="flex flex-col items-center gap-2">
            <div className="relative" style={{ width: size, height: size }}>
                {/* Background ring */}
                <svg className="absolute inset-0 -rotate-90" width={size} height={size}>
                    <circle
                        cx={size / 2}
                        cy={size / 2}
                        r={radius}
                        fill="none"
                        stroke={bgColor}
                        strokeWidth={strokeWidth}
                    />
                </svg>

                {/* Progress ring */}
                <svg className="absolute inset-0 -rotate-90" width={size} height={size}>
                    <motion.circle
                        cx={size / 2}
                        cy={size / 2}
                        r={radius}
                        fill="none"
                        stroke={color}
                        strokeWidth={strokeWidth}
                        strokeLinecap="round"
                        strokeDasharray={circumference}
                        initial={{ strokeDashoffset: circumference }}
                        animate={{ strokeDashoffset: offset }}
                        transition={{ duration: 1, ease: "easeOut" }}
                    />
                </svg>

                {/* Center content */}
                <div className="absolute inset-0 flex items-center justify-center">
                    <div className={cn("p-2 rounded-full", `bg-opacity-20`)} style={{ backgroundColor: `${color}20` }}>
                        {icon}
                    </div>
                </div>
            </div>

            {/* Label */}
            <div className="text-center">
                <p className="text-sm font-bold text-foreground">
                    {value.toLocaleString()}<span className="text-xs text-muted-foreground">/{max.toLocaleString()}</span>
                </p>
                <p className="text-xs text-muted-foreground">{label}</p>
            </div>
        </div>
    );
};

export default ProgressRing;
