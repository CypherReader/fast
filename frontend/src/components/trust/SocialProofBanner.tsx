import { Avatar, AvatarFallback } from "@/components/ui/avatar";

interface SocialProofBannerProps {
    totalUsers?: number;
    totalEarnedToday?: number;
}

export const SocialProofBanner = ({
    totalUsers = 12453,
    totalEarnedToday = 8420
}: SocialProofBannerProps) => {
    return (
        <div className="flex items-center gap-3 text-sm text-slate-400 mb-6 animate-fade-in">
            {/* Avatar Stack */}
            <div className="flex -space-x-3">
                {[...Array(5)].map((_, i) => (
                    <Avatar key={i} className="border-2 border-slate-950 h-8 w-8">
                        <AvatarFallback className="bg-slate-800 text-xs text-slate-300">
                            {['JS', 'MK', 'AL', 'RW', 'TH'][i]}
                        </AvatarFallback>
                    </Avatar>
                ))}
                <div className="flex items-center justify-center w-8 h-8 rounded-full bg-slate-800 border-2 border-slate-950 text-xs font-semibold text-slate-300">
                    +{(totalUsers - 5).toLocaleString()}
                </div>
            </div>

            {/* Social Proof Text */}
            <div>
                <p>
                    <span className="font-semibold text-white">{totalUsers.toLocaleString()}</span>
                    {' '}people earning with FastingHero
                </p>
                <p className="text-xs">
                    <span className="font-semibold text-emerald-400">${totalEarnedToday.toLocaleString()}</span>
                    {' '}earned today
                </p>
            </div>
        </div>
    );
};
