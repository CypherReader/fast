import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Copy, Check, Share2, MessageCircle } from "lucide-react";
import { useState } from "react";
import { useAuth } from "@/context/AuthContext";

interface ReferralModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

export const ReferralModal = ({ open, onOpenChange }: ReferralModalProps) => {
    const { user } = useAuth();
    const [copied, setCopied] = useState(false);

    const referralCode = user?.referral_code || "FASTHERO2024";
    const referralLink = `${window.location.origin}/register?ref=${referralCode}`;

    const handleCopy = () => {
        navigator.clipboard.writeText(referralLink);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const handleShare = async () => {
        if (navigator.share) {
            try {
                await navigator.share({
                    title: 'Join me on FastingHero',
                    text: 'I\'m earning money by fasting. Join the Vault and get $10 off!',
                    url: referralLink,
                });
            } catch (err) {
                console.error('Error sharing:', err);
            }
        } else {
            handleCopy();
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="bg-slate-950 border-slate-800 text-slate-200 sm:max-w-md">
                <DialogHeader>
                    <DialogTitle className="text-xl font-bold text-center text-white">
                        Invite Friends, Earn Cash
                    </DialogTitle>
                </DialogHeader>

                <div className="space-y-6 py-4">
                    <div className="text-center space-y-2">
                        <div className="bg-purple-900/20 p-4 rounded-full w-20 h-20 mx-auto flex items-center justify-center mb-4">
                            <Share2 className="w-10 h-10 text-purple-400" />
                        </div>
                        <p className="text-slate-400">
                            Share your unique link. When they join the Vault, you both get <span className="text-emerald-400 font-bold">$10</span>.
                        </p>
                    </div>

                    <div className="space-y-2">
                        <label className="text-xs font-medium text-slate-500 uppercase">Your Referral Link</label>
                        <div className="flex gap-2">
                            <Input
                                readOnly
                                value={referralLink}
                                className="bg-slate-900 border-slate-800 text-slate-300 font-mono text-sm"
                            />
                            <Button onClick={handleCopy} size="icon" className="bg-slate-800 hover:bg-slate-700">
                                {copied ? <Check className="w-4 h-4 text-emerald-400" /> : <Copy className="w-4 h-4" />}
                            </Button>
                        </div>
                    </div>

                    <div className="grid grid-cols-2 gap-3">
                        <Button
                            onClick={handleShare}
                            className="w-full bg-purple-600 hover:bg-purple-700 text-white"
                        >
                            <Share2 className="w-4 h-4 mr-2" /> Share
                        </Button>
                        <Button
                            variant="outline"
                            className="w-full border-slate-700 hover:bg-slate-800 text-slate-300"
                            onClick={() => window.open(`https://wa.me/?text=${encodeURIComponent(`Join me on FastingHero! ${referralLink}`)}`, '_blank')}
                        >
                            <MessageCircle className="w-4 h-4 mr-2" /> WhatsApp
                        </Button>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
};
