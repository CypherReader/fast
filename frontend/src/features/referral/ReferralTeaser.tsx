import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Gift, ArrowRight } from "lucide-react";

interface ReferralTeaserProps {
    onInvite: () => void;
}

export const ReferralTeaser = ({ onInvite }: ReferralTeaserProps) => {
    return (
        <Card className="bg-gradient-to-r from-purple-900/50 to-blue-900/50 border-purple-500/30 mb-6 overflow-hidden relative">
            <div className="absolute top-0 right-0 p-4 opacity-10">
                <Gift className="w-24 h-24 text-purple-400" />
            </div>
            <CardContent className="p-6 flex items-center justify-between relative z-10">
                <div>
                    <h3 className="text-lg font-bold text-white flex items-center gap-2">
                        <Gift className="w-5 h-5 text-purple-400" />
                        Give $10, Get $10
                    </h3>
                    <p className="text-sm text-slate-300 mt-1">
                        Invite a friend to the Vault and you both earn.
                    </p>
                </div>
                <Button
                    onClick={onInvite}
                    className="bg-purple-600 hover:bg-purple-700 text-white font-semibold shadow-lg shadow-purple-900/20"
                >
                    Invite <ArrowRight className="w-4 h-4 ml-2" />
                </Button>
            </CardContent>
        </Card>
    );
};
