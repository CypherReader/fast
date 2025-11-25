import React from 'react';
import { ShieldCheck, Info } from 'lucide-react';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";

interface PricingMechanismProps {
    currentPrice: number;
    disciplineScore: number;
}

const PricingMechanism: React.FC<PricingMechanismProps> = ({ currentPrice, disciplineScore }) => {
    return (
        <div className="bg-gradient-to-br from-slate-800 to-slate-900 rounded-xl p-6 border border-slate-700 text-center relative overflow-hidden h-full">
            <div className="absolute top-0 right-0 p-4 opacity-10">
                <ShieldCheck className="w-32 h-32 text-cyan-400" />
            </div>

            <div className="flex items-center justify-center mb-2 space-x-2 relative z-10">
                <h3 className="text-lg font-bold text-white">Dynamic Subscription</h3>
                <Dialog>
                    <DialogTrigger asChild>
                        <button className="text-slate-400 hover:text-cyan-400 transition-colors">
                            <Info className="w-4 h-4" />
                        </button>
                    </DialogTrigger>
                    <DialogContent className="bg-slate-900 border-slate-800 text-slate-200">
                        <DialogHeader>
                            <DialogTitle className="text-xl font-bold text-white flex items-center">
                                <ShieldCheck className="w-5 h-5 mr-2 text-cyan-400" />
                                The Lazy Tax
                            </DialogTitle>
                            <DialogDescription className="text-slate-400 pt-2">
                                We don't charge for features. We charge for lack of discipline.
                            </DialogDescription>
                        </DialogHeader>
                        <div className="space-y-4 text-sm mt-2">
                            <div className="bg-slate-800/50 p-3 rounded-lg border border-slate-700">
                                <p className="font-bold text-white mb-1">Base Price: $50.00 / month</p>
                                <p className="text-slate-400">This is the starting price for everyone.</p>
                            </div>
                            <div className="bg-slate-800/50 p-3 rounded-lg border border-slate-700">
                                <p className="font-bold text-green-400 mb-1">Your Goal: $0.00 / month</p>
                                <p className="text-slate-400">Every completed fast lowers your price. High discipline = Free access.</p>
                            </div>
                            <div className="bg-slate-800/50 p-3 rounded-lg border border-slate-700">
                                <p className="font-bold text-red-400 mb-1">The Penalty</p>
                                <p className="text-slate-400">If you quit a fast early, your discipline score drops, and the price goes back up.</p>
                            </div>
                        </div>
                    </DialogContent>
                </Dialog>
            </div>

            <p className="text-sm text-slate-400 mb-6 relative z-10">Demonstrate discipline. Lower your cost.</p>

            <div className="flex justify-center items-end space-x-2 mb-6 relative z-10">
                <div className="text-4xl font-bold text-white">${currentPrice.toFixed(2)}</div>
                <div className="text-sm text-slate-500 mb-1">/ month</div>
            </div>

            <div className="w-full bg-slate-700 rounded-full h-4 mb-2 overflow-hidden relative z-10">
                <div
                    className="bg-gradient-to-r from-red-500 via-yellow-500 to-green-500 h-full transition-all duration-1000"
                    style={{ width: `${disciplineScore}%` }}
                ></div>
            </div>
            <div className="flex justify-between text-xs text-slate-500 mb-6 relative z-10">
                <span>Lazy Tax ($50)</span>
                <span>Discipline (Free)</span>
            </div>

            <div className="grid grid-cols-3 gap-2 text-xs relative z-10">
                <div className="bg-slate-800 p-2 rounded border border-slate-700">
                    <div className="text-green-400 font-bold">+10 pts</div>
                    <div className="text-slate-400">Streak &gt; 7 Days</div>
                </div>
                <div className="bg-slate-800 p-2 rounded border border-slate-700">
                    <div className="text-green-400 font-bold">+5 pts</div>
                    <div className="text-slate-400">Ketones Verified</div>
                </div>
                <div className="bg-slate-800 p-2 rounded border border-slate-700">
                    <div className="text-green-400 font-bold">+5 pts</div>
                    <div className="text-slate-400">Tribe Support</div>
                </div>
            </div>
        </div>
    );
};

export default PricingMechanism;
